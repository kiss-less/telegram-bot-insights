package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbPath string, debugFlag bool) (*sql.DB, error) {

	if debugFlag {
		log.Println("Tables initialization started")
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Error connecting to DB: %v\n", err)
	}

	if debugFlag {
		log.Println("Initializing Bots table")
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS Bots (
			bot_id INTEGER PRIMARY KEY,
			api_key TEXT UNIQUE,
			name TEXT,
			description TEXT
		)
	`)
	if err != nil {
		log.Fatalf("Could not create Bots table: %v\n", err)
	}

	if debugFlag {
		log.Println("Initializing Users table")
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS Users (
			user_id INTEGER PRIMARY KEY,
			name TEXT,
			timestamp TEXT
		)
	`)
	if err != nil {
		log.Fatalf("Could not create Users table: %v\n", err)
	}

	if debugFlag {
		log.Println("Initializing BotsUsers table")
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS BotsUsers (
			bot_user_id INTEGER PRIMARY KEY,
			bot_id INTEGER,
			user_id INTEGER,
			is_active BOOLEAN,
			FOREIGN KEY (bot_id) REFERENCES Bots(bot_id),
			FOREIGN KEY (user_id) REFERENCES Users(user_id)
		)
	`)
	if err != nil {
		log.Fatalf("Could not create BotsUsers table: %v\n", err)
	}

	if debugFlag {
		log.Println("Initializing MessagesText table")
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS MessagesText (
			message_text_id INTEGER PRIMARY KEY,
			message_text TEXT UNIQUE
		)
	`)
	if err != nil {
		log.Fatalf("Could not create MessagesText table: %v\n", err)
	}

	if debugFlag {
		log.Println("Initializing Messages table")
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS Messages (
			message_id INTEGER PRIMARY KEY,
			bot_id INTEGER,
			user_id INTEGER,
			message_text_id INTEGER,
			http_status INTEGER,
			send_time TIMESTAMP,
			FOREIGN KEY (bot_id) REFERENCES Bots(bot_id),
			FOREIGN KEY (user_id) REFERENCES Users(user_id)
			FOREIGN KEY (message_text_id) REFERENCES MessagesText(message_text_id)
		)
	`)
	if err != nil {
		log.Fatalf("Could not create Messages table: %v\n", err)
	}

	if debugFlag {
		log.Println("Tables have been initialized")
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}
