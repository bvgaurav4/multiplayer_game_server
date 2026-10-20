package auth

import "time"

const (
	AccessTokenTTL  = time.Hour
	RefreshTokenTTL = 6 * time.Hour
)

type Auth struct {
	ID                    string    `json:"id"`
	UserID                string    `json:"user_id"`
	Token                 string    `json:"token"`
	RefreshToken          string    `json:"refresh_token"`
	TokenExpiresAt        time.Time `json:"token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

func (a *Auth) AccessTokenValid(now time.Time) bool {
	return a != nil && a.TokenExpiresAt.After(now)
}

func (a *Auth) RefreshTokenValid(now time.Time) bool {
	return a != nil && a.RefreshTokenExpiresAt.After(now)
}
