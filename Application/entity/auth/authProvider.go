package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	Provider "GameApplication/entity/entityProvider"

	"github.com/google/uuid"
)

const (
	userAuthTable = "user_auth"
	authColumns   = "id, user_id, token, refresh_token, token_expires_at, refresh_token_expires_at, created_at, updated_at"
)

var (
	ErrAuthNotFound        = errors.New("auth not found")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
	ErrInvalidAuth         = errors.New("invalid auth")
)

type scanner interface {
	Scan(dest ...any) error
}

func CreateAuth(ctx context.Context, p Provider.Provider, userID string) (*Auth, error) {
	db, err := p.Conn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, ErrInvalidAuth
	}

	token, err := generateToken()
	if err != nil {
		return nil, err
	}
	refreshToken, err := generateToken()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	a := &Auth{
		ID:                    uuid.New().String(),
		UserID:                userID,
		Token:                 token,
		RefreshToken:          refreshToken,
		TokenExpiresAt:        now.Add(AccessTokenTTL),
		RefreshTokenExpiresAt: now.Add(RefreshTokenTTL),
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	_, err = db.ExecContext(ctx,
		`INSERT INTO `+userAuthTable+`
		   (id, user_id, token, refresh_token, token_expires_at, refresh_token_expires_at, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		a.ID, a.UserID, a.Token, a.RefreshToken, a.TokenExpiresAt, a.RefreshTokenExpiresAt, a.CreatedAt, a.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func GetByRefreshToken(ctx context.Context, p Provider.Provider, refreshToken string) (*Auth, error) {
	db, err := p.Conn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return nil, ErrAuthNotFound
	}

	row := db.QueryRowContext(ctx,
		`SELECT `+authColumns+` FROM `+userAuthTable+` WHERE refresh_token = $1 LIMIT 1`, refreshToken)

	a, err := scanAuth(row)
	if err != nil {
		return nil, err
	}
	if !a.RefreshTokenValid(time.Now().UTC()) {
		return nil, ErrRefreshTokenExpired
	}

	return a, nil
}

func RenewAccessToken(ctx context.Context, p Provider.Provider, a *Auth) (*Auth, error) {
	db, err := p.Conn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if a == nil || strings.TrimSpace(a.ID) == "" {
		return nil, ErrInvalidAuth
	}

	now := time.Now().UTC()
	if !a.RefreshTokenValid(now) {
		return nil, ErrRefreshTokenExpired
	}

	token, err := generateToken()
	if err != nil {
		return nil, err
	}
	tokenExpiresAt := now.Add(AccessTokenTTL)

	result, err := db.ExecContext(ctx,
		`UPDATE `+userAuthTable+`
		    SET token = $2, token_expires_at = $3, updated_at = $4
		  WHERE id = $1 AND refresh_token_expires_at > $4`,
		a.ID, token, tokenExpiresAt, now)
	if err != nil {
		return nil, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, ErrRefreshTokenExpired
	}

	a.Token = token
	a.TokenExpiresAt = tokenExpiresAt
	a.UpdatedAt = now
	return a, nil
}

func scanAuth(row scanner) (*Auth, error) {
	var a Auth
	if err := row.Scan(
		&a.ID,
		&a.UserID,
		&a.Token,
		&a.RefreshToken,
		&a.TokenExpiresAt,
		&a.RefreshTokenExpiresAt,
		&a.CreatedAt,
		&a.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAuthNotFound
		}
		return nil, err
	}
	return &a, nil
}

func generateToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
