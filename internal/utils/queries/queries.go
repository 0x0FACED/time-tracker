package queries

const (
	GetUsers = `SELECT id, passport_number, pass_serie, surname, name, patronymic, address FROM users`

	AddUser = `INSERT INTO users (pass_number, pass_serie, surname, name, patronymic, address) VALUES ($1, $2, $3, $4, $5, $6)`

	DeleteUser = `DELETE FROM users WHERE id = $1`
)
