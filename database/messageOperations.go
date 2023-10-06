package database

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func CreateMessageText(db *sql.DB, text string) error {
	query := "INSERT INTO MessagesText (message_text) VALUES (?)"
	_, err := db.Exec(query, text)
	if err != nil {
		return err
	}
	return nil
}

func GetMessageTextByTextId(db *sql.DB, messageTextID int) (string, error) {
	query := "SELECT message_text FROM MessagesText WHERE message_text_id = ?"
	row := db.QueryRow(query, messageTextID)
	var text string

	err := row.Scan(&text)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", err
		}
		return "", err
	}

	return text, nil
}

func LogMessage(db *sql.DB, botID int, userID int, messageID int, httpStatus int) error {
	query := "INSERT INTO Messages (bot_id, user_id, message_text_id, http_status, send_time) VALUES (?, ?, ?, ?, ?)"

	_, err := db.Exec(query, botID, userID, messageID, time.Now(), httpStatus)
	if err != nil {
		return err
	}

	return nil
}

func MessageWasSentEarlier(db *sql.DB, botID int, userID int, messageID int, dateFormat string) (string, error) {
	query := "SELECT send_time FROM Messages WHERE bot_id = ? AND user_id = ? AND message_text_id = ?"
	row := db.QueryRow(query, botID, userID, messageID)
	var sendTime time.Time

	err := row.Scan(&sendTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", err
		}
		return "", err
	}

	return sendTime.Format(dateFormat), nil
}
