package queries

const (
	GetUsers = `
		SELECT id, passport_number, pass_serie, name, surname, patronymic, address 
		FROM users
	`

	AddUser = `
		INSERT INTO users (passport_number, pass_serie, name, surname, patronymic, address) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	UpdateUser = `
		UPDATE users 
		SET surname = $1, name = $2, patronymic = $3, address = $4 
		WHERE id = $5
	`
	DeleteUser = `
		DELETE FROM users 
		WHERE id = $1
	`

	DeleteUserTasks = `
		DELETE FROM tasks 
		WHERE user_id = $1
	`

	GetUserByID = `
		SELECT id, passport_number, pass_serie, surname, name, patronymic, address 
		FROM users 
		WHERE id = $1
	`

	GetUserWorklogs = `
		SELECT
			id,
			SUM(EXTRACT(EPOCH FROM (end_time - start_time))/3600) as hours,
			SUM(EXTRACT(EPOCH FROM (end_time - start_time))/60) as minutes
		FROM tasks
		WHERE user_id = $1 AND start_time >= $2 AND end_time <= $3
		GROUP BY id
		ORDER BY hours DESC, minutes DESC
	`

	AddStartTask = `
		INSERT INTO tasks (user_id, description, start_time) 
		VALUES ($1, $2, $3)
	`

	AddEndTask = `
		UPDATE tasks 
		SET end_time = $1 WHERE id = $2
	`
)
