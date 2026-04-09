package entry

import "time"

type Entry struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	MoodScore int       `json:"mood_score" db:"mood_score"`
	Sentiment *string   `json:"sentiment,omitempty" db:"sentiment"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
