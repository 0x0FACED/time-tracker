package storage

type Database interface {
	Connect() error
	Disconnect() error
}
