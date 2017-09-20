package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"github.com/nightlyone/lockfile"
	"bufio"
)

const tweetsPerPage = 20
const slowDuration = time.Second * 5
const downDuration = time.Second * 15

var health Health
var messages []Tweet
var sentTweets map[int64]bool
var tweetsCh chan Tweet

const persistentFile = "/storage/tweets.json"
const lockFile = "/storage/tweets.lock"

func receiveTweets() {
	last := time.Now().AddDate(0, -1, 0)

	for {
		time.Sleep(time.Second)
		now := time.Now().AddDate(0, 1, 0)
		tweets, duration, err := RequestTweets(last, now)
		if err != nil {
			println("[receiveTweets]", err.Error())
			if duration > downDuration {
				health.Twitter = StatusDown
			} else {
				health.Twitter = StatusError
			}
			continue
		}
		if duration > slowDuration {
			health.Twitter = StatusSlow
		} else {
			health.Twitter = StatusOperational
		}

		for _, tweet := range tweets {
			tweetsCh <- tweet
		}
	}
}

func sendTweet(tweet Tweet) {
	if sentTweets[tweet.Id] {
		return
	}

	slack := Slack{
		TweetId: tweet.Id,
		Text:    tweet.Body,
		IconURL: tweet.Link,
		Date:    tweet.Date,
		Team:    "our-team",
	}

	for {
		println("PushSlack", slack.TweetId)
		duration, err := PushSlack(slack)
		if err == nil {
			if duration > slowDuration {
				health.Slack = StatusSlow
			} else {
				health.Slack = StatusOperational
			}
			break
		}

		println("[push]", err.Error())
		if duration > downDuration {
			health.Slack = StatusDown
		} else {
			health.Slack = StatusError
		}
		time.Sleep(time.Second)
	}

	sentTweets[tweet.Id] = true
}

func lockOnFile() (lock lockfile.Lockfile, err error) {
	lock, err = lockfile.New(lockFile)
	if err != nil {
		return
	}

	for {
		err = lock.TryLock()
		if err == nil {
			break
		}
		log.Println("Wating for lock file...", err)
		time.Sleep(time.Microsecond * 100)
	}
	return
}

func readTweets() (tweets []Tweet, err error) {
	defer func() {
		log.Println("Reading tweets", len(tweets), err)
	}()

	lock, err := lockOnFile()
	if err != nil {
		return
	}
	defer lock.Unlock()

	f, err := os.Open(persistentFile)
	if err != nil {
		return
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			break
		}

		var tweet Tweet
		err = json.Unmarshal(line, &tweet)
		if err != nil {
			continue
		}

		tweets = append(tweets, tweet)
	}
	return
}

func saveTweet(tweet Tweet) (err error) {
	defer func() {
		log.Println("Saving tweet", tweet.Id, err)
	}()

	data, err := json.Marshal(tweet)
	if err != nil {
		return
	}
	data = append(data, '\n')

	lock, err := lockOnFile()
	if err != nil {
		return
	}
	defer lock.Unlock()

	f, err := os.OpenFile(persistentFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer f.Close()
	_, err = f.Write(data)
	return
}

func processTweets() {
	for tweet := range tweetsCh {
		messages = append(messages, tweet)
		sendTweet(tweet)
		saveTweet(tweet)
	}
}

func handleLatest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "408 Method not allowed", 408)
		return
	}

	n, err := strconv.Atoi(r.FormValue("n"))
	if err != nil {
		http.Error(w, "400 "+err.Error(), 400)
		return
	}

	messages, _ := readTweets()
	var latest []Tweet
	if n*tweetsPerPage < len(messages) {
		latest = messages[n*tweetsPerPage:]
	}
	if len(latest) > tweetsPerPage {
		latest = latest[:tweetsPerPage]
	}

	data, err := json.Marshal(latest)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 408)
		return
	}

	data, err := json.Marshal(&health)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func main() {
	transport := http.DefaultTransport.(*http.Transport)
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	transport.ResponseHeaderTimeout = time.Second * 15

	tweetsCh = make(chan Tweet, 100)
	sentTweets = make(map[int64]bool)

	go receiveTweets()
	go processTweets()

	health.App = StatusOperational
	health.Twitter = StatusOperational
	health.Slack = StatusOperational

	log.SetFlags(log.Ltime)
	log.SetOutput(os.Stderr)
	http.HandleFunc("/latest", handleLatest)
	http.HandleFunc("/health", handleHealth)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
