package auth

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func HashingPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hash), err
}

func CheckPassword(email string, password string, db *sql.DB) (int, bool) {
	var dbPassword string
	var id int

	err := db.QueryRow("SELECT id, password FROM user WHERE email = ?", email).Scan(&id, &dbPassword)

	if err == sql.ErrNoRows {
		log.Printf("[ERROR] User not found: %v", email)
		return 0, false
	} else if err != nil {
		log.Printf("[ERROR] DB error: %v", err)
		return 0, false
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	return id, err == nil
}
