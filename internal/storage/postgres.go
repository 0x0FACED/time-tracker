package storage

import (
	"database/sql"
	_ "database/sql"
	"fmt"
	"log"
	"time-tracker/configs"
	"time-tracker/internal/models"
	"time-tracker/internal/utils/queries"
	"time-tracker/migrations"

	_ "github.com/lib/pq"
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

func (p *Postgres) GetUsers() (map[int]models.User, error) {
	rows, err := p.sql.Query("SELECT id, passport_number, passport_serie, surname, name, patronymic, address FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make(map[int]models.User)
	for rows.Next() {
		var id, passportNumber, passportSerie int
		var surname, name, patronymic, address string
		rows.Scan(&id, &passportNumber, &passportSerie, &surname, &name, &patronymic, &address)
		user := models.User{
			Id:         id,
			PassNumber: passportNumber,
			PassSerie:  passportSerie,
			Surname:    surname,
			Patronymic: patronymic,
			Address:    address,
		}
		users[id] = user
	}
	return users, nil
}

func (p *Postgres) GetTasksByUserID(id int) ([]models.Task, error) {
	// TODO: impl
	panic("not impl")
}

func (p *Postgres) AddUser(u *models.User) error {
	// TODO: impl
	panic("not impl")
}

func (p *Postgres) DeleteUser(id int) error {
	// TODO: impl
	panic("not impl")
}

func (p *Postgres) UpdateUser(u *models.User) error {
	// TODO: impl
	panic("not impl")
}

func (p *Postgres) AddStartTask(t *models.Task) error {
	// TODO: impl
	panic("not impl")
}

func (p *Postgres) AddEndTask(t *models.Task) error {
	// TODO: impl
	panic("not impl")
}
