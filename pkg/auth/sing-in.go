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

func (a *InputSignIn) SignIn(db *sql.DB) (string, error){
	id, check := CheckPassword(a.Email, a.Password, db)

	if !check{
		return "", fmt.Errorf("Invalid login:password")
	}

	jwtKey, err := GenerateJWTkey(a.Email, id)

	if err != nil{
		return "", fmt.Errorf("[ERROR] Sign in jwt-Generate ERROR: %v", err)
	}

	return jwtKey, nil
}