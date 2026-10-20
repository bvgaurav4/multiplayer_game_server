package entityProvider

import (
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Provider struct {
	DB string
}

func (p Provider) Conn() (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", p.DB)
	if err != nil {
		return nil, err
	}

	// Configure the connection pool
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}
