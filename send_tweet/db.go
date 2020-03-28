package main

import "database/sql"

func GetTweet(db *sql.DB) (int64, string, error) {
	var content string
	var id int64
	err := db.QueryRow("SELECT id, content FROM LexBezos.tweets WHERE tweeted_datetime IS NULL ORDER BY inserted_datetime LIMIT 1").Scan(&id, &content)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", nil
		} else {
			return 0, "", err
		}
	}
	return id, content, nil
}

func UpdateTwitInfo(db *sql.DB, id, twitID int64) error {
	stmt, err := db.Prepare("UPDATE LexBezos.tweets SET tweeted_datetime = UTC_TIMESTAMP, twitter_id = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(twitID, id)
	if err != nil {
		return err
	}
	return nil
}
