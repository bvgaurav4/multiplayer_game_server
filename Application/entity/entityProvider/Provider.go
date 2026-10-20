package entityProvider

import (
	Config "GameApplication/config"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Provider struct {
	Config *Config.Config
	DB     *sqlx.DB
}

func (p *Provider) Conn() (*sqlx.DB, error) {
	if p.DB == nil {
		return nil, fmt.Errorf("database is not connected")
	}

	return p.DB, nil
}
func (p *Provider) Connect() error {
	fmt.Print("connectin to ", p)
	db, err := sqlx.Connect("pgx", p.Config.Database.DSN())
	if err != nil {
		return err
	}

	// Configure the connection pool
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(30 * time.Minute)
	p.DB = db
	return nil
}
