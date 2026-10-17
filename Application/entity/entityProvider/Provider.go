package entityProvider

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Provider struct {
	DB string
}

func (p Provider) Conn() (*sqlx.DB, error) {
	return sqlx.Connect("pgx", p.DB)
}