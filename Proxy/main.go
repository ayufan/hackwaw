package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
	"strings"

	"github.com/coreos/pkg/flagutil"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"mime"
	"bytes"
)

var listenAddress = flag.String("listen-address", ":8080", "Application listen address")
var consumerKey = flag.String("twitter-consumer-key", "", "Twitter Consumer Key")
var consumerSecret = flag.String("twitter-consumer-secret", "", "Twitter Consumer Secret")
var accessToken = flag.String("twitter-access-token", "", "Twitter Access Token")
var accessSecret = flag.String("twitter-access-secret", "", "Twitter Access Secret")
var slackHookURL = flag.String("slack-hook-url", "", "Slack Push address")
var twitterTrackFilter = flag.String("twitter-track-filter", "", "Twitter Track Filter")
var twitterLanguageFilter = flag.String("twitter-language-filter", "", "Twitter Language Filter")
var twitterLocationFilter = flag.String("twitter-location-filter", "", "Twitter Location Filter")
var tweetsPerPage = flag.Int("tweets-per-page", 10, "Number of tweets to send per-page")
var maxTweets = flag.Int("max-tweets", 1000, "How many tweets store in memory")

var tweets []*Tweet = []*Tweet{}

func createTwitterClient() (*twitter.Client, error) {
	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		return nil, errors.New("consumer key/secret and Access token/secret required")
	}

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	return twitter.NewClient(httpClient), nil
}

func insertTweet(tweet *Tweet) {
	idx := sort.Search(len(tweets), func(idx int) bool {
		return tweets[idx].Date.After(tweet.Date)
	})
	if idx < len(tweets) {
		newTweets := make([]*Tweet, len(tweets)+1)
		copy(newTweets[:idx], tweets[:idx])
		copy(newTweets[idx+1:], tweets[idx:])
		newTweets[idx] = tweet
		tweets = newTweets
	} else {
		tweets = append(tweets, tweet)
	}
	if len(tweets) > *maxTweets {
		tweets = tweets[len(tweets)-*maxTweets:]
	}
}

func processTweet(tweet *twitter.Tweet) {
	createdAt, err := time.Parse(time.RubyDate, tweet.CreatedAt)
	if err != nil {
		log.Println("[tweet]", err)
		return
	}

	ourTweet := &Tweet{
		Id:   tweet.ID,
		Body: tweet.Text,
		Date: createdAt,
	}
	insertTweet(ourTweet)
}

func splitCommaSeparatedString(in string) []string {
	if in == "" {
		return nil
	}
	return strings.Split(in, ",")
}

func receiveTweets() (stream *twitter.Stream, err error) {
	client, err := createTwitterClient()
	if err != nil {
		return
	}

	// Use Twitter Stream
	filterParams := &twitter.StreamFilterParams{
		StallWarnings: twitter.Bool(true),
		Track:         splitCommaSeparatedString(*twitterTrackFilter),
		Language:      splitCommaSeparatedString(*twitterLanguageFilter),
		Locations:     splitCommaSeparatedString(*twitterLocationFilter),
	}
	if len(filterParams.Track) == 0 {
		return nil, errors.New("Define --twitter-track-filter")
	}

	stream, err = client.Streams.Filter(filterParams)
	if err != nil {
		return
	}

	// Handle Twitter Stream
	demux := twitter.NewSwitchDemux()
	demux.Tweet = processTweet

	log.Println("[twitter] starting receiving...")
	go func() {
		demux.HandleChan(stream.Messages)
		log.Println("[twitter] finished")
	}()
	return
}

func filterTweets(fromDate, toDate time.Time) []*Tweet {
	fromDate = fromDate.Add(-time.Nanosecond)
	allTweets := tweets
	lowerBound := sort.Search(len(allTweets), func(idx int) bool {
		return fromDate.Before(allTweets[idx].Date)
	})
	upperBound := sort.Search(len(allTweets)-lowerBound, func(idx int) bool {
		return allTweets[idx+lowerBound].Date.After(toDate)
	})
	return allTweets[lowerBound : lowerBound+upperBound]
}

func paginateTweets(tweets []*Tweet, page, perPage int) (paged []*Tweet, pages int) {
	if page < 0 {
		page = 0
	}
	if perPage <= 0 {
		perPage = *tweetsPerPage
	}

	idx := page * perPage
	if idx > len(tweets) {
		idx = len(tweets)
	}

	maxIdx := (page + 1) * perPage
	if maxIdx > len(tweets) {
		maxIdx = len(tweets)
	}

	paged = tweets[idx:maxIdx]
	pages = (len(tweets) + perPage - 1) / perPage
	return
}

func getTweets(fromDate, toDate time.Time, page, perPage int) (data []byte, pages int, err error) {
	ourTweets := filterTweets(fromDate, toDate)
	ourTweets, pages = paginateTweets(ourTweets, page, perPage)
	data, err = json.Marshal(ourTweets)
	return
}

func formValueOrDefault(r *http.Request, key, defaultValue string) (value string) {
	value = r.FormValue(key)
	if value == "" {
		value = defaultValue
	}
	return
}

func handleTweets(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "408 Method Not Allowed", 408)
		return
	}

	fromDate, err := time.Parse(time.RFC3339, formValueOrDefault(r, "from", time.Now().Add(-time.Minute).Format(time.RFC3339)))
	if err != nil {
		http.Error(w, "400 'from': "+err.Error(), 400)
		return
	}

	toDate, err := time.Parse(time.RFC3339, formValueOrDefault(r, "to", time.Now().Format(time.RFC3339)))
	if err != nil {
		http.Error(w, "400 'to': "+err.Error(), 400)
		return
	}

	page, err := strconv.Atoi(formValueOrDefault(r, "page", "0"))
	if err != nil {
		http.Error(w, "400 'page': "+err.Error(), 400)
		return
	}

	perPage, err := strconv.Atoi(formValueOrDefault(r, "per-page", strconv.Itoa(*tweetsPerPage)))
	if err != nil {
		http.Error(w, "400 'page': "+err.Error(), 400)
		return
	}

	data, pages, err := getTweets(fromDate, toDate, page, perPage)
	if err != nil {
		http.Error(w, "500 "+err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Pages", strconv.Itoa(pages))
	w.Write(data)
}

func handlePush(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "408 Method Not Allowed", 408)
		return
	}

	mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		http.Error(w, "400 Error parsing Content-Type: " + err.Error(), 400)
		return
	}
	if mt != "application/json" {
		http.Error(w, "400 Not application/json: " + mt, 400)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "400 "+err.Error(), 400)
		return
	}

	var slack Slack
	err = json.Unmarshal(data, &slack)
	if err != nil {
		http.Error(w, "400 "+err.Error(), 400)
		return
	}

	payload := struct {
		Username string `json:"username"`
		IconURL string `json:"icon_url"`
		Text string `json:"text"`
	}{
		Username: slack.Team,
		IconURL: slack.IconURL,
		Text: fmt.Sprintf("Tweet with id %v, at %v with message %v",
			slack.TweetId, slack.Date, slack.Text),
	}
	payloadData, err := json.Marshal(&payload)
	if err != nil {
		http.Error(w, "500 "+err.Error(), 500)
		return
	}

	resp, err := http.Post(*slackHookURL, "application/json", bytes.NewReader(payloadData))
	if err != nil {
		http.Error(w, "400 "+err.Error(), 400)
		return
	}
	defer resp.Body.Close()

	// Reply back with Slack response
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main() {
	err := flagutil.SetFlagsFromEnv(flag.CommandLine, "APP")
	if err != nil {
		log.Fatalln(err)
	}
	flag.Parse()

	stream, err := receiveTweets()
	if err != nil {
		log.Fatalln(err)
	}
	defer stream.Stop()

	log.Println("Starting server...")
	log.SetFlags(log.Ltime | log.Lshortfile)
	log.SetOutput(os.Stderr)
	http.HandleFunc("/tweets", handleTweets)
	http.HandleFunc("/push", handlePush)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
