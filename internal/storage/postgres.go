package storage

import (
	"database/sql"
	_ "database/sql"
	"fmt"
	"time-tracker/configs"
)

type Postgres struct {
	sql *sql.DB
	cfg configs.DatabaseConfig
}

func New(cfg configs.DatabaseConfig) *Postgres {
	return &Postgres{
		cfg: cfg,
	}
}

func (p *Postgres) connectionString() string {
	return fmt.Sprintf("%s:%s@%s:%s/%s?sslmode=disable",
		p.cfg.DBUsername, p.cfg.DBPass, p.cfg.DBHost, p.cfg.DBPort, p.cfg.DBName)
}

func (p *Postgres) Connect() error {
	db, err := sql.Open("postgres", p.connectionString())
	if err != nil {
		return err
	}
	if db.Ping() != nil {
		return err
	}
	p.sql = db
	return nil
}

func (p *Postgres) Disconnect() error {
	// TODO: impl
	panic("not impl")
}
