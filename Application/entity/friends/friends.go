package friends

import "time"

type Friend struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	FriendID  string    `json:"friend_id"`
	Alias     string    `json:"alias"`
	Status    string    `json:"status"`
	AddedAt   time.Time `json:"added_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
