package main

import (
	"database/sql"
	"strings"
)

func filterPostedListings(db *sql.DB, listings []Listing) (approvedListings []Listing, badTitleListings []Listing, err error) {
	stmtURL, err := db.Prepare("SELECT id FROM LexBezos.articles WHERE url = ?")
	if err != nil {
		return
	}
	var goodURLListings []Listing
	for _, listing := range listings {
		var id int
		err = stmtURL.QueryRow(listing.Url).Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				err = nil
				goodURLListings = append(goodURLListings, listing)
				continue
			} else {
				return approvedListings, badTitleListings, err
			}
		}
	}

	//We'll get all the news headlines from the past month
	rows, err := db.Query("SELECT title FROM LexBezos.articles WHERE last_seen BETWEEN DATE_SUB(NOW(), INTERVAL 30 DAY) AND NOW()")
	if err != nil {
		return
	}
	defer rows.Close()

	var titleRows []string
	for rows.Next() {
		var title string
		err = rows.Scan(&title)
		if err != nil {
			return
		}
		titleRows = append(titleRows, title)
	}
	if rows.Err() != nil {
		return
	}

	for _, listing := range goodURLListings {
		foundMatch := false
		for _, title := range titleRows {
			//This seems like a good threshold of equality
			if CompareTwoStrings(strings.ToLower(title), strings.ToLower(listing.Title)) > 0.6 {
				foundMatch = true
				break
			}
		}
		if foundMatch == true {
			badTitleListings = append(badTitleListings, listing)
		} else {
			approvedListings = append(approvedListings, listing)
		}
	}

	return
}

func storeListings(db *sql.DB, listings []Listing) error {
	stmtInsert, err := db.Prepare("INSERT INTO LexBezos.articles (url, title, last_seen) VALUES (?, ?, UTC_TIMESTAMP())")
	if err != nil {
		return err
	}
	for i, listing := range listings {
		res, err := stmtInsert.Exec(listing.Url, listing.Title)
		if err != nil {
			return err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		listings[i].DBID = id
	}
	return nil
}

func PopulateTweetTable(db *sql.DB, listings []Listing) error {
	stmtInsert, err := db.Prepare("INSERT INTO LexBezos.tweets (article_id, content, inserted_datetime) VALUES (?, ?, UTC_TIMESTAMP())")
	if err != nil {
		return err
	}
	for _, listing := range listings {
		content := listing.Title + " " + listing.Url
		_, err = stmtInsert.Exec(listing.DBID, content)
		if err != nil {
			return err
		}
	}

	return nil
}
