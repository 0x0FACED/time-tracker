package storage

import "time-tracker/internal/models"

type Database interface {
	Connect() error
	Disconnect() error

	UserDatabase
	TaskDatabase
}

type UserDatabase interface {
	GetUsers() (map[int]models.User, error)
	GetTasksByUserID(id int) ([]models.Task, error)
	AddUser(u *models.User) error
	DeleteUser(id int) error
	UpdateUser(u *models.User) error
}

type TaskDatabase interface {
	AddStartTask(t *models.Task) error
	AddEndTask(t *models.Task) error
}
