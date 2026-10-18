package feedback

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	Provider "GameApplication/entity/entityProvider"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	feedbackTable = "feedback"
)

var ErrInvalidFeedback = errors.New("invalid feedback")

func CreateFeedback(ctx context.Context, p Provider.Provider, f *Feedback) error {
	db, err := p.Conn()
	if err != nil {
		return err
	}
	defer db.Close()

	if f == nil {
		return ErrInvalidFeedback
	}

	f.UserID = strings.TrimSpace(f.UserID)
	f.Message = strings.TrimSpace(f.Message)
	if f.Message == "" {
		return fmt.Errorf("%w: message is required", ErrInvalidFeedback)
	}

	now := time.Now().UTC()
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	if f.CreatedAt.IsZero() {
		f.CreatedAt = now
	}

	_, err = db.ExecContext(ctx,
		`INSERT INTO `+feedbackTable+` (id, user_id, message, created_at)
		 VALUES ($1, $2, $3, $4)`,
		f.ID, f.UserID, f.Message, f.CreatedAt)

	return translateError(err)
}

func translateError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23503":
			return fmt.Errorf("%w: referenced user does not exist", ErrInvalidFeedback)
		}
	}

	return err
}
