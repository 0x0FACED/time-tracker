package migrations

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

func Up(url string) {
	m, err := migrate.New(
		"file://./migrations/",
		url)
	if err != nil {
		log.Fatalln("failed to create migration: ", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalln("failed to migrate up: ", err)
	}
}
