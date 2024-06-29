package models

type User struct {
	Id         int    `json:"id" db:"id"`
	PassNumber int    `json:"pass_number" db:"pass_number"`
	PassSerie  int    `json:"pass_serie" db:"pass_serie"`
	Surname    string `json:"surname" db:"surname"`
	Name       string `json:"name" db:"name"`
	Patronymic string `json:"patronymic" db:"patronymic"`
	Address    string `json:"address" db:"address"`
}
