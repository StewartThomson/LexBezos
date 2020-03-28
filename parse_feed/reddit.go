package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Response struct {
	Data Data `json:"data"`
}

type Data struct {
	Listings []Listing `json:"children"`
}

type Listing struct {
	ListingData `json:"data"`
}

type ListingData struct {
	Title string `json:"title"`
	Url   string `json:"url"`
	DBID  int64  `json:"-"`
}

func getListings(client HTTPClient) (listings []Listing, err error) {
	req, _ := http.NewRequest("GET", os.Getenv(EnvFeedUrl), nil)

	//I'm just going to spoof the user agent for now. Will I ever authenticate it properly? Probably not.
	//TODO: Authenticate properly
	req.Header.Set("User-Agent", randomUA())
	logger.Debug("Hitting reddit")
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("On reddit request: %v", err)
		return
	}
	defer resp.Body.Close()

	var redditResp Response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Unable to read response: %v", err)
		return
	}
	logger.Debug("Got: " + string(body))
	err = json.Unmarshal(body, &redditResp)
	if err != nil {
		return
	}

	listings = getJeffListings(redditResp.Data.Listings)
	return
}

func getJeffListings(listings []Listing) []Listing {
	var approvedListings []Listing
	for _, listing := range listings {
		lowerTitle := strings.ToLower(listing.Title)
		if strings.Contains(lowerTitle, "amazon") {
			//Trying to filter out the rainforest... not in the news often but this may show up
			if !strings.Contains(lowerTitle, "burn") &&
				!strings.Contains(lowerTitle, "rainforest") &&
				!strings.Contains(lowerTitle, "tribe") {
				approvedListings = append(approvedListings, listing)
			}
		} else if strings.Contains(lowerTitle, "bezos") || strings.Contains(lowerTitle, "amzn") {
			approvedListings = append(approvedListings, listing)
		}
	}

	return approvedListings
}
