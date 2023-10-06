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
	query := "INSERT INTO Messages (bot_id, user_id, message_text_id, send_time, http_status) VALUES (?, ?, ?, ?, ?)"

	_, err := db.Exec(query, botID, userID, messageID, time.Now(), httpStatus)
	if err != nil {
		return err
	}

	return nil
}

func MessageWasSentEarlier(db *sql.DB, botID int, userID int, messageID int) (bool, error) {
	query := "SELECT COUNT(send_time) FROM Messages WHERE bot_id = ? AND user_id = ? AND message_text_id = ?"

	var count int
	err := db.QueryRow(query, botID, userID, messageID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func TotalNumberOfMessagesSentPerBot(db *sql.DB, botID int, messageID int) (int, error) {
	query := "SELECT COUNT(send_time) FROM Messages WHERE bot_id = ? AND message_text_id = ?"

	var count int
	err := db.QueryRow(query, botID, messageID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func TotalNumberOfMessagesPerBotByStatus(db *sql.DB, botID int, messageID int, httpStatus int) (int, error) {
	query := "SELECT COUNT(send_time) FROM Messages WHERE bot_id = ? AND message_text_id = ? AND http_status = ?"

	var count int
	err := db.QueryRow(query, botID, messageID, httpStatus).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func LastMessageSentAndDeliveredTimePerBot(db *sql.DB, botID int, messageID int) (string, error) {
	query := "SELECT send_time FROM Messages WHERE bot_id = ? AND message_text_id = ? AND http_status = 200 ORDER BY send_time DESC LIMIT 1"

	var time string
	err := db.QueryRow(query, botID, messageID).Scan(&time)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", err
		}
		return "", err
	}

	return time, nil
}

func LastMessageSentTimePerBot(db *sql.DB, botID int, messageID int) (string, error) {
	query := "SELECT send_time FROM Messages WHERE bot_id = ? AND message_text_id = ? ORDER BY send_time DESC LIMIT 1"

	var time string
	err := db.QueryRow(query, botID, messageID).Scan(&time)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", err
		}
		return "", err
	}

	return time, nil
}
