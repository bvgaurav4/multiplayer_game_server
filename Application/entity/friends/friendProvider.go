package friends

import (
	"context"
	"strings"

	Provider "GameApplication/entity/entityProvider"

	"github.com/jmoiron/sqlx"
)

const (
	friendsTable  = "friends"
	friendColumns = "id, user_id, friend_id, alias, status, added_at, updated_at"
)

func GetFriends(ctx context.Context, p Provider.Provider, userID string) ([]*Friend, error) {
	db, err := p.Conn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	userID = strings.TrimSpace(userID)

	var rows *sqlx.Rows
	if userID == "" {
		rows, err = db.QueryxContext(ctx,
			`SELECT `+friendColumns+` FROM `+friendsTable+`
			 ORDER BY added_at DESC`)
	} else {
		rows, err = db.QueryxContext(ctx,
			`SELECT `+friendColumns+` FROM `+friendsTable+`
			 WHERE user_id = $1
			 ORDER BY added_at DESC`, userID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*Friend, 0)
	for rows.Next() {
		f := &Friend{}
		if err := rows.Scan(
			&f.ID,
			&f.UserID,
			&f.FriendID,
			&f.Alias,
			&f.Status,
			&f.AddedAt,
			&f.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, f)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}
