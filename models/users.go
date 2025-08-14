package models

import (
	"database/sql"
	"log"
)

const UsersSchema = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(24) UNIQUE NOT NULL,
		password TEXT NOT NULL);`

type User struct {
	Id       int
	Username string
	Password string
}

func UserCreate(db *sql.DB, username string, password string) error {
	_, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, password)
	return err
}

func UserQueryById(db *sql.DB, id string) User {
	var user User

	row := db.QueryRow("SELECT * FROM users WHERE id=?", id)
	err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Password,
	)

	if err != nil {
		log.Println(err.Error())
		return User{}
	}

	return user
}

func UserQueryByUsernameFull(db *sql.DB, username string) User {
	var user User

	row := db.QueryRow("SELECT * FROM users WHERE username=?", username)
	err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Password,
	)

	if err != nil {
		log.Println(err.Error())
		return User{}
	}

	return user
}
