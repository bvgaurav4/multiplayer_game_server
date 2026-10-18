package feedback

import "time"

type Feedback struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
