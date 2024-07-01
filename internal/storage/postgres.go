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

func (p *Postgres) GetUsers(query string, params ...any) (map[int]models.User, error) {
	rows, err := p.sql.Query(query, params...)
	if err != nil {
		log.Println("err in query: ", err)
		return nil, err
	}
	defer rows.Close()

	users := make(map[int]models.User)
	for rows.Next() {
		var id int
		var surname, name, patronymic, address, passportNumber, passportSerie string
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
	err := row.Scan(&user.Id, &user.PassNumber, &user.PassSerie, &user.Name, &user.Surname, &user.Patronymic, &user.Address)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *Postgres) GetUserWorklogs(req *models.GetUserWorklogsRequest) ([]models.Worklog, error) {
	query := `
        SELECT
            task_id,
            SUM(EXTRACT(EPOCH FROM (end_time - start_time))/3600) as hours,
            SUM(EXTRACT(EPOCH FROM (end_time - start_time))/60) as minutes
        FROM tasks
        WHERE user_id = $1 AND start_time >= $2 AND end_time <= $3
        GROUP BY task_id
        ORDER BY hours DESC, minutes DESC
    `

	rows, err := p.sql.Query(query, req.UserID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var worklogs []models.Worklog
	for rows.Next() {
		var worklog models.Worklog
		if err := rows.Scan(&worklog.TaskID, &worklog.Hours, &worklog.Minutes); err != nil {
			return nil, err
		}
		worklogs = append(worklogs, worklog)
	}

	return worklogs, nil
}

func (p *Postgres) AddUser(u *models.User) error {
	_, err := p.sql.Exec(queries.AddUser, u.PassNumber, u.PassSerie, u.Name, u.Surname, u.Patronymic, u.Address)
	if err != nil {
		log.Fatalln("cant create user: ", err)
		return err
	}
	return nil
}

// TODO: transactions bcz when delete user -> delete all his tasks. must be transaction
func (p *Postgres) DeleteUser(id int) error {
	_, err := p.sql.Exec(queries.DeleteUser, id)
	if err != nil {
		log.Fatalln("cant delete user: ", err)
		return err
	}
	return nil
}

func (p *Postgres) UpdateUser(u *models.User) error {
	_, err := p.sql.Exec("UPDATE users SET surname = $1, name = $2, patronymic = $3, address = $4 WHERE id = $5", u.Surname, u.Name, u.Patronymic, u.Address, u.Id)
	if err != nil {
		log.Println("cant update user: ", err)
		return err
	}
	return nil
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
