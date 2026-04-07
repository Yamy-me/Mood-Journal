package entity

import "time"

type User struct {
	ID           int64  `json:"-" db:"id"`
	Email        string `json:"email" db:"email"`
	Name         string `json:"name" db:"name"`
	PasswordHash string `json:"-" db:"password_hash"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`

	StreakDays  int        `json:"streak_days" db:"streak_days"`
}