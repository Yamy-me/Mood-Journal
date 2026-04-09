package entry

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Create(UserID int, Content string, MoodScore int, Sentiment *string) (int64, error) {
	result, err := r.DB.Exec("INSERT INTO entries (user_id, content, mood_score, sentiment) VALUES (?, ?, ?, ?);", UserID, Content, MoodScore, Sentiment)
	if err != nil {
		return 0, fmt.Errorf("[ERROR] Can't create entry: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("[ERROR] Can't get last insert id: %w", err)
	}

	return id, nil
}

func (r *Repository) GetByID(id int, userID int) (*Entry, error) {
    var entry Entry
    
    query := `
        SELECT id, user_id, content, mood_score, sentiment, created_at 
        FROM entries 
        WHERE id = ? AND user_id = ?;`

    err := r.DB.QueryRow(query, id, userID).Scan(
        &entry.ID, 
        &entry.UserID, 
        &entry.Content, 
        &entry.MoodScore, 
        &entry.Sentiment, 
        &entry.CreatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("entry not found: %w", err)
    }
    if err != nil {
        return nil, fmt.Errorf("database error: %w", err)
    }

    return &entry, nil
}

func (r *Repository) GetAllByUser(UserID int) ([]Entry, error) {
	rows, err := r.DB.Query("SELECT id, user_id, content, mood_score, sentiment, created_at FROM entries WHERE user_id = ? ORDER BY created_at DESC;", UserID)
	if err != nil {
		return nil, fmt.Errorf("ERROR Query script: %w", err)
	}
	defer rows.Close()

	var res []Entry
	for rows.Next() {
		var e Entry

		err := rows.Scan(&e.ID, &e.UserID, &e.Content, &e.MoodScore, &e.Sentiment, &e.CreatedAt)

		if err != nil {
			return nil, fmt.Errorf("Error in QueryRows iteration: %w", err)
		}

		res = append(res, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Repository) Delete(ID int, USER_ID int) (error, int64) {
	result, err := r.DB.Exec("DELETE FROM entries WHERE id = ? AND user_id = ?;", ID, USER_ID)
	if err != nil {
		return fmt.Errorf("[ERROR] can't delete: %w", err), -1
	}

	if res, _ := result.RowsAffected(); res == 0 {
		return fmt.Errorf("There's nothing entry with this id"), 0
	}

	return nil, 1
}
func (r *Repository) UpdateSentiment(id int64, sentiment string) error {
	_, err := r.DB.Exec("UPDATE entries SET sentiment = ? WHERE id = ?", sentiment, id)
	if err != nil {
		return fmt.Errorf("[ERROR] Can't update sentiment: %w", err)
	}
	return nil
}