package main

import "time"

type Tweet struct {
	Id        int64     `json:"id"`
	Body      string    `json:"body"`
	Link      string    `json:"link,omitempty"`
	Date      time.Time `json:"date"`
}

type Slack struct {
	Team    string    `json:"team,omitempty"`
	TweetId int64     `json:"tweetId"`
	IconURL string    `json:"icon_url,omitempty"`
	Text    string    `json:"text,omitempty"`
	Date    time.Time `json:"date"`
}
