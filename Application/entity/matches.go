package entity

import "time"

type Match struct {
	ID       string    `json:"id"`
	Players  []string  `json:"players"`
	Mode     string    `json:"mode"`
	Result   string    `json:"result"`
	PlayedAt time.Time `json:"played_at"`
}
