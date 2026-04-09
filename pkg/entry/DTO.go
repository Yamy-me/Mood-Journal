package entry

// POST /entries Создание
type CreateEntryInput struct {
	Content   string `json:"content" binding:"required,min=1,max=1000"`
	MoodScore int    `json:"mood_score" binding:"required,gte=1,lte=10"`
}

// Response Ответ
type EntryResponse struct {
	ID        int64   `json:"id"`
	Content   string  `json:"content"`
	MoodScore int     `json:"mood_score"`
	Sentiment *string `json:"sentiment,omitempty"`
	CreatedAt string  `json:"created_at"`
}

// GET /entries Фильтры
type GetEntriesQuery struct {
	Limit  int    `form:"limit" binding:"omitempty,gte=1,lte=100"`
	Offset int    `form:"offset" binding:"omitempty,gte=0"`
	From   string `form:"from" binding:"omitempty,datetime=2006-01-02"`
	To     string `form:"to" binding:"omitempty,datetime=2006-01-02"`
}
