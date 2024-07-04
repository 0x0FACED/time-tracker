package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"time-tracker/configs"
	"time-tracker/internal/models"
	"time-tracker/internal/utils"
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
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
		p.cfg.Driver, p.cfg.Username, p.cfg.Pass, p.cfg.Host, p.cfg.Port, p.cfg.Name)
}

func (p *Postgres) Connect() error {
	url := p.connectionString()
	db, err := sql.Open(p.cfg.Driver, url)
	if err != nil {
		return errors.New(utils.ErrCantOpenDB)
	}
	if db.Ping() != nil {
		return errors.New(utils.ErrCantPingDB)
	}
	migrations.Up(url)
	p.sql = db
	return nil
}

func (p *Postgres) Disconnect() error {
	// TODO: impl
	panic("not impl")
}

func (p *Postgres) GetUsers(req models.GetUsersRequest) (map[int]models.User, error) {
	query := `
		SELECT id, passport_number, pass_serie, surname, name, patronymic, address 
		FROM users 
		WHERE 1=1
	`
	params := []interface{}{}
	paramCounter := 1

	if req.PassportNumber != "" {
		query += fmt.Sprintf(" AND passport_number = $%d", paramCounter)
		params = append(params, req.PassportNumber)
		paramCounter++
	}
	if req.PassSerie != "" {
		query += fmt.Sprintf(" AND pass_serie = $%d", paramCounter)
		params = append(params, req.PassSerie)
		paramCounter++
	}
	if req.Surname != "" {
		query += fmt.Sprintf(" AND surname = $%d", paramCounter)
		params = append(params, req.Surname)
		paramCounter++
	}
	if req.Name != "" {
		query += fmt.Sprintf(" AND name = $%d", paramCounter)
		params = append(params, req.Name)
		paramCounter++
	}
	if req.Patronymic != "" {
		query += fmt.Sprintf(" AND patronymic = $%d", paramCounter)
		params = append(params, req.Patronymic)
		paramCounter++
	}
	if req.Address != "" {
		query += fmt.Sprintf(" AND address = $%d", paramCounter)
		params = append(params, req.Address)
		paramCounter++
	}
	// pagination
	offset := (req.Page - 1) * req.PageSize
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCounter, paramCounter+1)
	params = append(params, req.PageSize, offset)

	rows, err := p.sql.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf(utils.ErrQuery, query, params, err)
	}
	defer rows.Close()

	users := make(map[int]models.User)
	for rows.Next() {
		var id int
		var surname, name, patronymic, address, passportNumber, passportSerie string
		if err := rows.Scan(&id, &passportNumber, &passportSerie, &surname, &name, &patronymic, &address); err != nil {
			return nil, fmt.Errorf(utils.ErrScanRow, query, params, err)
		}
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

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(utils.ErrRowIteration, query, params, err)
	}
	return users, nil
}

func (p *Postgres) GetUserByID(id int) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, passport_number, pass_serie, surname, name, patronymic, address 
		FROM users 
		WHERE id = $1
	`
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
			id,
			description,
			ROUND(SUM(EXTRACT(EPOCH FROM (end_time - start_time))/3600))::integer as hours,
			ROUND(SUM(EXTRACT(EPOCH FROM (end_time - start_time))/60))::integer as minutes
		FROM tasks
		WHERE user_id = $1 AND start_time >= $2 AND end_time <= $3
		GROUP BY id
		ORDER BY hours DESC, minutes DESC
	`
	rows, err := p.sql.Query(query, req.UserID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, fmt.Errorf(utils.ErrQuery, query, req, err)
	}
	defer rows.Close()

	var worklogs []models.Worklog
	for rows.Next() {
		var worklog models.Worklog
		if err := rows.Scan(&worklog.TaskID, &worklog.Desc, &worklog.Hours, &worklog.Minutes); err != nil {
			return nil, fmt.Errorf(utils.ErrScanRow, query, req, err)
		}
		worklogs = append(worklogs, worklog)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(utils.ErrRowIteration, query, req, err)
	}
	return worklogs, nil
}

func (p *Postgres) AddUser(u *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (passport_number, pass_serie, name, surname, patronymic, address) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := p.sql.Exec(query, u.PassNumber, u.PassSerie, u.Name, u.Surname, u.Patronymic, u.Address)
	if err != nil {
		return nil, fmt.Errorf(utils.ErrQuery, query, u, err)
	}
	return u, nil
}

func (p *Postgres) DeleteUser(id int) error {
	_, err := p.GetUserByID(id)
	if err != nil {
		return fmt.Errorf(utils.ErrQuery, "GetUserByID(id)", id, err)
	}
	tx, err := p.sql.Begin()
	if err != nil {
		return fmt.Errorf(utils.ErrBeginTx, err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	query := `
		DELETE FROM tasks 
		WHERE user_id = $1
	`
	_, err = tx.Exec(query, id)
	if err != nil {
		return fmt.Errorf(utils.ErrQuery, query, id, err)
	}

	query = `
		DELETE FROM users 
		WHERE id = $1
	`
	_, err = tx.Exec(query, id)
	if err != nil {
		return fmt.Errorf(utils.ErrQuery, query, id, err)
	}

	return tx.Commit()
}

func (p *Postgres) UpdateUser(u *models.User) error {
	query := `
		UPDATE users 
		SET surname = $1, name = $2, patronymic = $3, address = $4 
		WHERE id = $5
	`
	_, err := p.sql.Exec(query, u.Surname, u.Name, u.Patronymic, u.Address, u.Id)
	if err != nil {
		return fmt.Errorf(utils.ErrQuery, query, u, err)
	}
	return nil
}

func (p *Postgres) AddStartTask(t *models.Task) error {
	query := `
		INSERT INTO tasks (user_id, description, start_time) 
		VALUES ($1, $2, $3)
	`
	now := time.Now()
	_, err := p.sql.Exec(query, t.UserID, t.Desc, now)
	if err != nil {
		return fmt.Errorf(utils.ErrQuery, query, t, err)
	}
	return nil
}

func (p *Postgres) AddEndTask(id int) error {
	query := `
		UPDATE tasks 
		SET end_time = $1 WHERE id = $2
	`
	now := time.Now()
	_, err := p.sql.Exec(query, now, id)
	if err != nil {
		return fmt.Errorf(utils.ErrQuery, query, id, err)
	}
	return nil
}
