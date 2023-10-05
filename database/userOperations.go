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
