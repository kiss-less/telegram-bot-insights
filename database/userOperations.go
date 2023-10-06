package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID        int
	Name      string
	Timestamp string
}

func GetUserByID(db *sql.DB, id int) (*User, error) {
	query := "SELECT user_id, name, timestamp, FROM Users WHERE user_id = ?"
	row := db.QueryRow(query, id)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func CreateUser(db *sql.DB, id int, name string, timestamp string) error {
	query := "INSERT INTO Users (user_id, name, timestamp) VALUES (?, ?, ?)"
	_, err := db.Exec(query, id, name, timestamp)
	if err != nil {
		return err
	}
	return nil
}

func UserExists(db *sql.DB, id int) (bool, error) {
	query := "SELECT COUNT(*) FROM Users WHERE user_id = ?"
	var count int
	err := db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetUsersByBotID(db *sql.DB, botID int) ([]User, error) {
	query := "SELECT user_id, name, timestamp FROM BotsUsers WHERE bot_id = ?"

	rows, err := db.Query(query, botID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Timestamp); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func MakeUserInactiveByBotID(db *sql.DB, botID int, userID int) error {
	query := "UPDATE BotsUsers SET is_active = 0 WHERE bot_id = ? AND user_id = ?"

	_, err := db.Exec(query, botID, userID)
	if err != nil {
		return err
	}

	return nil
}

func UserIsActive(db *sql.DB, userID int, botID int) (bool, error) {
	query := "SELECT COUNT(*) FROM BotsUsers WHERE user_id = ? AND bot_id = ? AND is_active = 1"
	var count int
	err := db.QueryRow(query, userID, botID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
