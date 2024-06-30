package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"
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
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		p.cfg.DBUsername, p.cfg.DBPass, p.cfg.DBHost, p.cfg.DBPort, p.cfg.DBName)
}

func (p *Postgres) Connect() error {
	url := p.connectionString()
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatalln("cant open db: ", err)
		return err
	}
	if db.Ping() != nil {
		return err
	}
	migrations.Up(url)
	p.sql = db
	log.Println("DB AFTER OPEN: ", p.sql)
	return nil
}

func (p *Postgres) Disconnect() error {
	// TODO: impl
	panic("not impl")
}

func (p *Postgres) GetUsers() (map[int]models.User, error) {
	rows, err := p.sql.Query(queries.GetUsers)
	if err != nil {
		log.Println("err in query: ", err)
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

func (p *Postgres) GetUserByID(id int) (*models.User, error) {
	var user models.User
	query := "SELECT id, passport_number, pass_serie, surname, name, patronymic, address FROM users WHERE id = $1"
	row := p.sql.QueryRow(query, id)
	err := row.Scan(&user.Id, &user.PassNumber, &user.PassSerie, &user.Surname, &user.Name, &user.Patronymic, &user.Address)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *Postgres) GetTasksByUserID(id int) ([]models.Task, error) {
	// TODO: impl
	panic("not impl")
}

func (p *Postgres) AddUser(u *models.User) error {
	_, err := p.sql.Exec(queries.AddUser)
	if err != nil {
		log.Fatalln("cant create user: ", err)
		return err
	}
	return nil
}

func (p *Postgres) DeleteUser(id int) error {
	_, err := p.sql.Exec(queries.DeleteUser, id)
	if err != nil {
		log.Fatalln("cant delete user: ", err)
		return err
	}
	return nil
}

func (p *Postgres) UpdateUser(u *models.User) error {
	// TODO: impl
	panic("not impl")
}

func (p *Postgres) AddStartTask(t *models.Task) error {
	now := time.Now()
	_, err := p.sql.Exec("INSERT INTO tasks (user_id, description, start_time) VALUES ($1, $2, $3)", t.UserID, t.Desc, now)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) AddEndTask(id int) error {
	now := time.Now()
	_, err := p.sql.Exec("UPDATE tasks SET end_time = $1 WHERE id = $2", now, id)
	if err != nil {
		return err
	}
	return nil
}
