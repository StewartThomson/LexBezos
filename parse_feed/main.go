package main

import (
	"database/sql"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const EnvFeedUrl = "LB_FEED_URL"
const EnvDBHost = "LB_DB_HOST"
const EnvDBUser = "LB_DB_USER"
const EnvDBPass = "LB_DB_PASS"

type input struct {
	Debug bool `json:"debug"`
}

var logger *zap.SugaredLogger

func handler(request input) error {
	var err error
	var l *zap.Logger
	if request.Debug == true {
		l, err = zap.NewDevelopment()
	} else {
		l, err = zap.NewProduction()
	}
	if err != nil {
		return err
	}
	logger = l.Sugar()
	defer logger.Sync()

	logger.Debug("Starting function")

	client := &http.Client{}

	listings, err := getListings(client)
	if err != nil {
		return err
	}

	for i := range listings {
		clean, err := cleanURL(listings[i].Url)
		if err != nil {
			return err
		}
		listings[i].Url = clean
	}

	db, err := sql.Open("mysql", os.Getenv(EnvDBUser)+":"+os.Getenv(EnvDBPass)+"@tcp("+os.Getenv(EnvDBHost)+")/")
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	listings, badListings, err := filterPostedListings(db, listings)
	if err != nil {
		return err
	}
	//We'll store the bad ones as well, so that we can keep their URLs for ez checking in the DB next time
	err = storeListings(db, badListings)
	if err != nil {
		return err
	}
	err = storeListings(db, listings)
	if err != nil {
		return err
	}

	listings = modifyTitles(listings)

	err = PopulateTweetTable(db, listings)
	if err != nil {
		return err
	}
	return nil
}

func modifyTitles(listings []Listing) (modified []Listing) {
	replaceMap := map[string]string{
		"Jeff":   "Lex",
		"jeff":   "Lex",
		"Bezos'": "Luthor's",
		"bezos'": "Luthor's",
		"Bezos":  "Luthor",
		"bezos":  "Luthor",
		"Amazon": "LexCorp",
		"amazon": "LexCorp",
		"amzn":   "lxcr",
		"AMZN":   "LXCR",
	}
	for _, listing := range listings {
		for oldWord, newWord := range replaceMap {
			listing.Title = strings.ReplaceAll(listing.Title, oldWord, newWord)
		}
		modified = append(modified, listing)
	}

	return
}

func cleanURL(inputUrl string) (string, error) {
	u, err := url.Parse(inputUrl)
	if err != nil {
		return "", err
	}
	u.RawQuery = ""
	return u.String(), nil
}

func randomUA() string {
	userAgents := []string{
		"PostmanRuntime/7.24.0", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36", "Mozilla/5.0 (Windows NT 5.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36",
		"Mozilla/5.0 (Windows NT 5.1; rv:7.0.1) Gecko/20100101 Firefox/7.0.1", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:40.0) Gecko/20100101 Firefox/40.1", "Mozilla/5.0 (iPhone; CPU iPhone OS 12_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.1 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 13_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.4 Mobile/15E148 Safari/604.1",
	}

	rand.Seed(time.Now().Unix())
	return userAgents[rand.Intn(len(userAgents))]
}

func main() {
	lambda.Start(handler)
}
