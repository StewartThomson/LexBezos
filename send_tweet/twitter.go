package main

import (
	"github.com/dghubble/go-twitter/twitter"
	"net/http"
)

func SendTweet(httpClient *http.Client, content string) (int64, error) {
	client := twitter.NewClient(httpClient)

	tweet, _, err := client.Statuses.Update(content, nil)
	if err != nil {
		return 0, err
	}

	return tweet.ID, nil
}
