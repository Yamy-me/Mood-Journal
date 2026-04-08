package auth

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)







func HashingPassword(password string) (string, error){
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hash), err
}

func CheckPassword(email string, password string, db *sql.DB) (int, bool) {
    

    var dbPassword string
    err := db.QueryRow("SELECT password FROM user WHERE email = ?", email).Scan(&dbPassword)

    if err == sql.ErrNoRows {
        log.Printf("[ERROR] User not found: %v", email)
        return 0, false
    } else if err != nil {
        log.Printf("[ERROR] DB error: %v", err)
        return 0, false
    }

    var id int
    idScan_EmailScan := db.QueryRow("SELECT id FROM user WHERE email = ?;", email).Scan(&id)
     if idScan_EmailScan == sql.ErrNoRows {
        log.Printf("[ERROR] User not found: %v", email)
        return 0, false
    } else if idScan_EmailScan != nil {
        log.Printf("[ERROR] DB error: %v", err)
        return 0, false
    }

    err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
    return id, err == nil
}