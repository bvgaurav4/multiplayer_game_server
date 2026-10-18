package matches

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	Provider "GameApplication/entity/entityProvider"

	"github.com/jmoiron/sqlx"
)

const (
	matchesTable = "matches"
	playedTable  = "played"
	matchColumns = "id, mode, result, played_at"
)

var ErrMatchNotFound = errors.New("match not found")

type scanner interface {
	Scan(dest ...any) error
}

func GetMatchesByUserID(ctx context.Context, p Provider.Provider, userID string) ([]*Match, error) {
	db, err := p.Conn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, ErrMatchNotFound
	}

	rows, err := db.QueryxContext(ctx,
		`SELECT m.id, m.mode, m.result, m.played_at
		   FROM `+matchesTable+` m
		   JOIN `+playedTable+` p ON p.match_id = m.id
		  WHERE p.user_id = $1
		  ORDER BY m.played_at DESC`, userID)
	if err != nil {
		return nil, translateError(err)
	}
	defer rows.Close()

	list := make([]*Match, 0)
	ids := make([]string, 0)
	for rows.Next() {
		m, err := scanMatch(rows)
		if err != nil {
			return nil, translateError(err)
		}
		list = append(list, m)
		ids = append(ids, m.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, translateError(err)
	}
	if len(list) == 0 {
		return list, nil
	}

	players, err := playersByMatchIDs(ctx, db, ids)
	if err != nil {
		return nil, err
	}
	for _, m := range list {
		m.Players = players[m.ID]
		if m.Players == nil {
			m.Players = []string{}
		}
	}

	return list, nil
}

func GetMatchByID(ctx context.Context, p Provider.Provider, matchID string) (*Match, error) {
	db, err := p.Conn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	matchID = strings.TrimSpace(matchID)
	if matchID == "" {
		return nil, ErrMatchNotFound
	}

	row := db.QueryRowContext(ctx,
		`SELECT `+matchColumns+` FROM `+matchesTable+` WHERE id = $1 LIMIT 1`, matchID)

	m, err := scanMatch(row)
	if err != nil {
		return nil, err
	}

	players, err := playersByMatchIDs(ctx, db, []string{matchID})
	if err != nil {
		return nil, err
	}
	m.Players = players[matchID]
	if m.Players == nil {
		m.Players = []string{}
	}

	return m, nil
}

func playersByMatchIDs(ctx context.Context, db *sqlx.DB, matchIDs []string) (map[string][]string, error) {
	placeholders := make([]string, 0, len(matchIDs))
	args := make([]any, 0, len(matchIDs))
	for i, id := range matchIDs {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		args = append(args, id)
	}

	rows, err := db.QueryxContext(ctx,
		`SELECT match_id, user_id FROM `+playedTable+`
		  WHERE match_id IN (`+strings.Join(placeholders, ", ")+`)
		  ORDER BY user_id`, args...)
	if err != nil {
		return nil, translateError(err)
	}
	defer rows.Close()

	players := make(map[string][]string, len(matchIDs))
	for rows.Next() {
		var matchID, userID string
		if err := rows.Scan(&matchID, &userID); err != nil {
			return nil, translateError(err)
		}
		players[matchID] = append(players[matchID], userID)
	}
	if err := rows.Err(); err != nil {
		return nil, translateError(err)
	}

	return players, nil
}

func scanMatch(row scanner) (*Match, error) {
	var m Match
	if err := row.Scan(
		&m.ID,
		&m.Mode,
		&m.Result,
		&m.PlayedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMatchNotFound
		}
		return nil, err
	}
	return &m, nil
}

func translateError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrMatchNotFound
	}
	return err
}
