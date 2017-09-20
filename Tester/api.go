package main

import (
	"time"
)

const (
	StatusOperational string = "OPERATIONAL"
	StatusSlow               = "SLOW"
	StatusError              = "ERROR"
	StatusDown               = "DOWN"
	StatusUnnecessary        = "UNNECESSARY"
)

type Health struct {
	App      string `json:"app"`
	Database string `json:"database"`
	Twitter  string `json:"twitter"`
	Slack    string `json:"slack"`
}

type Tweet struct {
	Id        int64     `json:"id"`
	TwitterId int64     `json:"twitterId,omitempty"`
	Body      string    `json:"body"`
	Link      string    `json:"link"`
	Date      time.Time `json:"date"`
}

func (t *Tweet) UniqueID() int {
	if t.TwitterId != 0 {
		return int(t.TwitterId)
	} else {
		return int(t.Id)
	}
}

type Slack struct {
	Team    string    `json:"team"`
	TweetId int64     `json:"tweetId"`
	IconURL string    `json:"icon_url"`
	Text    string    `json:"text"`
	Date    time.Time `json:"date"`
}
