package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	Provider "GameApplication/entity/entityProvider"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	usersTable   = "users"
	userColumns  = "id, username, email, password, status, created_at, updated_at"
	userStatuses = "online, offline, in_match, away"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidUser       = errors.New("invalid user")
	ErrNoFieldsToUpdate  = errors.New("no fields to update")
)

type scanner interface {
	Scan(dest ...any) error
}

func GetUserWithUsernameOrEmail(ctx context.Context, p Provider.Provider, identifier string) (*User, error) {
	db, err := p.Conn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	value := normalize(identifier)
	if value == "" {
		return nil, ErrUserNotFound
	}

	row := db.QueryRowContext(ctx,
		`SELECT `+userColumns+` FROM `+usersTable+`
		 WHERE LOWER(username) = $1 OR LOWER(email) = $1
		 LIMIT 1`, value)

	u, err := scanUser(row)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func GetUserByID(ctx context.Context, p Provider.Provider, id string) (*User, error) {
	db, err := p.Conn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if strings.TrimSpace(id) == "" {
		return nil, ErrUserNotFound
	}

	row := db.QueryRowContext(ctx,
		`SELECT `+userColumns+` FROM `+usersTable+` WHERE id = $1 LIMIT 1`, id)

	u, err := scanUser(row)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func CreateUser(ctx context.Context, p Provider.Provider, u *User) error {
	db, err := p.Conn()
	if err != nil {
		return err
	}
	defer db.Close()

	if u == nil {
		return ErrInvalidUser
	}

	u.Username = strings.TrimSpace(u.Username)
	u.Email = normalize(u.Email)
	if u.Username == "" || u.Email == "" {
		return fmt.Errorf("%w: username and email are required", ErrInvalidUser)
	}
	if u.Status == "" {
		u.Status = "offline"
	}
	if !validStatus(u.Status) {
		return fmt.Errorf("%w: status must be one of %s", ErrInvalidUser, userStatuses)
	}

	now := time.Now().UTC()
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	u.UpdatedAt = now

	_, err = db.ExecContext(ctx,
		`INSERT INTO `+usersTable+`
		   (id, username, email, password, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		u.ID, u.Username, u.Email, u.Password, u.Status, u.CreatedAt, u.UpdatedAt)

	return translateError(err)
}

func UpdateUser(ctx context.Context, p Provider.Provider, u *User) error {
	db, err := p.Conn()
	if err != nil {
		return err
	}
	defer db.Close()

	if u == nil || strings.TrimSpace(u.ID) == "" {
		return ErrUserNotFound
	}

	username := strings.TrimSpace(u.Username)
	email := normalize(u.Email)
	if username == "" && email == "" && u.Password == "" && u.Status == "" {
		return ErrNoFieldsToUpdate
	}
	if u.Status != "" && !validStatus(u.Status) {
		return fmt.Errorf("%w: status must be one of %s", ErrInvalidUser, userStatuses)
	}

	updatedAt := u.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now().UTC()
	}

	result, err := db.ExecContext(ctx,
		`UPDATE `+usersTable+`
		    SET username = COALESCE(NULLIF($2, ''), username),
		        email    = COALESCE(NULLIF($3, ''), email),
		        password = COALESCE(NULLIF($4, ''), password),
		        status   = COALESCE(NULLIF($5, ''), status),
		        updated_at = $6
		  WHERE id = $1`,
		u.ID, username, email, u.Password, u.Status, updatedAt)
	if err != nil {
		return translateError(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrUserNotFound
	}

	u.Username = username
	u.Email = email
	u.UpdatedAt = updatedAt
	return nil
}

func GetUsersByUsernames(ctx context.Context, p Provider.Provider, usernames []string) ([]User, error) {
	db, err := p.Conn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	names := uniqueNormalized(usernames)
	if len(names) == 0 {
		return []User{}, nil
	}

	placeholders := make([]string, 0, len(names))
	args := make([]any, 0, len(names))
	for i, name := range names {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		args = append(args, name)
	}

	rows, err := db.QueryxContext(ctx,
		`SELECT `+userColumns+` FROM `+usersTable+`
		 WHERE LOWER(username) IN (`+strings.Join(placeholders, ", ")+`)
		 ORDER BY LOWER(username)`, args...)
	if err != nil {
		return nil, translateError(err)
	}
	defer rows.Close()

	users := make([]User, 0, len(names))
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, translateError(err)
		}
		users = append(users, *u)
	}
	if err := rows.Err(); err != nil {
		return nil, translateError(err)
	}

	return users, nil
}

func scanUser(row scanner) (*User, error) {
	var u User
	if err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.Status,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func translateError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrUserNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return fmt.Errorf("%w: %s", ErrUserAlreadyExists, pgErr.ConstraintName)
		case "23503":
			return fmt.Errorf("%w: referenced record does not exist", ErrInvalidUser)
		}
	}

	return err
}

func normalize(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func uniqueNormalized(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	unique := make([]string, 0, len(values))
	for _, value := range values {
		normalized := normalize(value)
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		unique = append(unique, normalized)
	}
	return unique
}

func validStatus(status string) bool {
	for _, allowed := range strings.Split(userStatuses, ", ") {
		if status == allowed {
			return true
		}
	}
	return false
}
