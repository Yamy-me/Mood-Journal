package auth

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type InputSignUp struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}


func NewInputSignUp(name string, email string, password string) *InputSignUp{
	return &InputSignUp{
		Name: name,
		Email: email,
		Password: password,
	}
}


func (a *InputSignUp) validate(db *sql.DB) error {
    var exists bool
    err := db.QueryRow("SELECT 1 FROM user WHERE email = ?", a.Email).Scan(&exists)

    if err == nil {
        return fmt.Errorf("user with this email already exists")
    }

    if err != sql.ErrNoRows {
        return fmt.Errorf("[ERROR] DB error: %w", err)
    }

    return nil
}


func (a *InputSignUp) SignUp(database *sql.DB) error{
	errValidate := a.validate(database)

	if errValidate != nil{
		return fmt.Errorf("[ERROR] SignUp VALIDATE error: %v", errValidate)
	}
	hashedPassword, errHash := HashingPassword(a.Password)
	if errHash != nil{
		return fmt.Errorf("[ERROR] SignUp Hashing error: %v", errHash)
	}
	_, errDB := database.Exec("INSERT INTO user (name, email, password, created_at, updated_at, streak_days) VALUES (?, ?, ?, ?, ?, 0)", a.Name, a.Email, hashedPassword, time.Now(), time.Now().Format("2006-01-02 15:04:05"))

	if errDB != nil{
		return fmt.Errorf("[ERROR] SignUp DATABASE error: %v", errDB)
	}

	return nil
}