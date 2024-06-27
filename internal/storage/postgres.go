package storage

import (
	"database/sql"
	_ "database/sql"
	"time-tracker/configs"
)

type Postgres struct {
	sql *sql.DB
	cfg configs.DatabaseConfig
}

func (p *Postgres) Connect() error {
	// TODO: impl
	panic("not impl")
}

func (p *Postgres) Disconnect() error {
	// TODO: impl
	panic("not impl")
}
