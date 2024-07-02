package storage

import "time-tracker/internal/models"

type Database interface {
	Connect() error
	Disconnect() error

	UserDatabase
	TaskDatabase
}

type UserDatabase interface {
	// get users with filters and pagination
	GetUsers(query string, params ...any) (map[int]models.User, error)
	GetUserByID(id int) (*models.User, error)
	GetUserWorklogs(req *models.GetUserWorklogsRequest) ([]models.Worklog, error)
	AddUser(u *models.User) (*models.User, error)
	DeleteUser(id int) error
	UpdateUser(u *models.User) error
}

type TaskDatabase interface {
	AddStartTask(t *models.Task) error
	AddEndTask(id int) error
}
