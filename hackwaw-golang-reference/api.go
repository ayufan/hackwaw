package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
	"log"
	"mime"
)

type Status string

const (
	StatusOperational Status = "OPERATIONAL"
	StatusSlow               = "SLOW"
	StatusError              = "ERROR"
	StatusDown               = "DOWN"
	StatusUnnecessary        = "UNNECESSARY"
)

type Health struct {
	App      Status `json:"app,omitempty"`
	Database Status `json:"database,omitempty"`
	Twitter  Status `json:"twitter,omitempty"`
	Slack    Status `json:"slack,omitempty"`
}

type Tweet struct {
	Id        int64     `json:"id"`
	Body      string    `json:"body,omitempty"`
	Link      string    `json:"link,omitempty"`
	Date      time.Time `json:"date"`
}

type Slack struct {
	Team    string    `json:"team"`
	TweetId int64     `json:"tweetId"`
	IconURL string    `json:"icon_url,omitempty"`
	Text    string    `json:"text,omitempty"`
	Date    time.Time `json:"date,omitempty"`
}

func RequestTweets(from time.Time, to time.Time) (messages []Tweet, duration time.Duration, err error) {
	started := time.Now()
	defer func() {
		duration = time.Since(started)
	}()

	twitterURL := os.Getenv("TWITTER_URL")
	if twitterURL == "" {
		err = errors.New("Missing TWITTER_URL")
		return
	}

	url, err := url.Parse(twitterURL + "/tweets")
	if err != nil {
		return
	}

	query := url.Query()
	query.Add("from", from.Format(time.RFC3339Nano))
	query.Add("to", to.Format(time.RFC3339Nano))
	url.RawQuery = query.Encode()

	log.Println("GET", url)

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return
	}
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		err = fmt.Errorf("Invalid status code: %d, got: %d", 200, res.StatusCode)
		return
	}

	contentType, _, err := mime.ParseMediaType(res.Header.Get("Content-Type"))
	if err != nil {
		return
	}
	if contentType != "application/json" {
		err = fmt.Errorf("Not application/json: %v", contentType)
		return
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	log.Println("GET", url, string(data))

	err = json.Unmarshal(data, &messages)
	return
}

func PushSlack(message Slack) (duration time.Duration, err error) {
	started := time.Now()
	defer func() {
		duration = time.Since(started)
	}()

	slackURL := os.Getenv("SLACK_URL")
	if slackURL == "" {
		err = errors.New("Missing SLACK_URL")
		return
	}

	data, err := json.Marshal(&message)
	if err != nil {
		return
	}

	log.Println("POST", slackURL, string(data))

	req, err := http.NewRequest("POST", slackURL + "/push", bytes.NewReader(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		err = fmt.Errorf("Expected status code: %d, got: %d", 200, res.StatusCode)
		return
	}
	return
}
