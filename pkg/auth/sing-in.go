package auth

import (
	"database/sql"
	"fmt"
)

type InputSignIn struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func NewSignInInput(email string, password string) *InputSignIn {
	return &InputSignIn{
		Email:    email,
		Password: password,
	}
}

func (a *InputSignIn) SignIn(db *sql.DB) (bool, error){
	if !CheckPassword(a.Email, a.Password, db){
		return false, fmt.Errorf("Invalid password")
	}

	return true, nil
}