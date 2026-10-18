package reqres

import (
	match "GameApplication/entity/matches"
	"time"
)

type SignUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

// GET /users/:id/matches
type GetMatchesRequest struct {
	UserID string `json:"user_id" uri:"id" form:"user_id"`
}

type GetMatchesResponse struct {
	UserID string         `json:"user_id"`
	Matches []*match.Match `json:"matches"`
	Count   int            `json:"count"`
}

// GET /matches/:id
type GetMatchDetailsRequest struct {
	MatchID string `json:"match_id" uri:"id" form:"match_id"`
}

type GetMatchDetailsResponse = match.Match

// POST /feedback
type CreateFeedbackRequest struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type CreateFeedbackResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
