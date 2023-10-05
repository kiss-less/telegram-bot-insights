package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Bot struct {
	ID          int
	APIKey      string
	Name        string
	Description string
}

func GetBotByID(db *sql.DB, id int) (*Bot, error) {
	query := "SELECT bot_id, api_key, name, description FROM Bots WHERE bot_id = ?"
	row := db.QueryRow(query, id)

	var bot Bot
	err := row.Scan(&bot.ID, &bot.APIKey, &bot.Name, &bot.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &bot, nil
}

func CreateBot(db *sql.DB, id int, apiKey string, name string, desc string) error {
	query := "INSERT INTO Bots (bot_id, api_key, name, description) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, id, apiKey, name, desc)
	if err != nil {
		return err
	}
	return nil
}

func BotExists(db *sql.DB, id int) (bool, error) {
	query := "SELECT COUNT(*) FROM Bots WHERE bot_id = ?"
	var count int
	err := db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func AssociateBotWithUser(db *sql.DB, botID, userID int) error {
	query := "SELECT COUNT(*) FROM BotsUsers WHERE bot_id = ? AND user_id = ?"
	var count int
	err := db.QueryRow(query, botID, userID).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		query := "INSERT INTO BotsUsers (bot_id, user_id) VALUES (?, ?)"
		_, err := db.Exec(query, botID, userID)
		if err != nil {
			return err
		}
		return nil
	} else {
		err := fmt.Errorf("association %v:%v exists", botID, userID)
		return err
	}
}
