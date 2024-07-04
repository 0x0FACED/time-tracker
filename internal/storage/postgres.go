package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"time-tracker/configs"
	"time-tracker/internal/models"
	"time-tracker/internal/utils"
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
	row := p.sql.QueryRow(queries.GetUserByID, id)
	err := row.Scan(&user.Id, &user.PassNumber, &user.PassSerie, &user.Name, &user.Surname, &user.Patronymic, &user.Address)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *Postgres) GetUserWorklogs(req *models.GetUserWorklogsRequest) ([]models.Worklog, error) {
	rows, err := p.sql.Query(queries.GetUserWorklogs, req.UserID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, fmt.Errorf(utils.ErrQuery, queries.GetUserWorklogs, req, err)
	}
	defer rows.Close()

	var worklogs []models.Worklog
	for rows.Next() {
		var worklog models.Worklog
		if err := rows.Scan(&worklog.TaskID, &worklog.Hours, &worklog.Minutes); err != nil {
			return nil, fmt.Errorf(utils.ErrScanRow, queries.GetUserWorklogs, req, err)
		}
		worklogs = append(worklogs, worklog)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(utils.ErrRowIteration, queries.GetUserWorklogs, req, err)
	}
	return worklogs, nil
}

func (p *Postgres) AddUser(u *models.User) (*models.User, error) {
	_, err := p.sql.Exec(queries.AddUser, u.PassNumber, u.PassSerie, u.Name, u.Surname, u.Patronymic, u.Address)
	if err != nil {
		return nil, fmt.Errorf(utils.ErrQuery, queries.AddUser, u, err)
	}
	return u, nil
}

func (p *Postgres) DeleteUser(id int) error {
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

	_, err = tx.Exec(queries.DeleteUserTasks, id)
	if err != nil {
		return fmt.Errorf(utils.ErrQuery, queries.DeleteUserTasks, id, err)
	}

	_, err = tx.Exec(queries.DeleteUser, id)
	if err != nil {
		return fmt.Errorf(utils.ErrQuery, queries.DeleteUser, id, err)
	}

	return tx.Commit()
}

func (p *Postgres) UpdateUser(u *models.User) error {
	_, err := p.sql.Exec(queries.UpdateUser, u.Surname, u.Name, u.Patronymic, u.Address, u.Id)
	if err != nil {
		return fmt.Errorf(utils.ErrQuery, queries.UpdateUser, u, err)
	}
	return nil
}

func (p *Postgres) AddStartTask(t *models.Task) error {
	now := time.Now()
	_, err := p.sql.Exec(queries.AddStartTask, t.UserID, t.Desc, now)
	if err != nil {
		return fmt.Errorf(utils.ErrQuery, queries.AddStartTask, t, err)
	}
	return nil
}

func (p *Postgres) AddEndTask(id int) error {
	now := time.Now()
	_, err := p.sql.Exec(queries.AddEndTask, now, id)
	if err != nil {
		return fmt.Errorf(utils.ErrQuery, queries.AddEndTask, id, err)
	}
	return nil
}
