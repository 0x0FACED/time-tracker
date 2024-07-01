package models

type User struct {
	Id         int    `json:"id" db:"id"`
	PassNumber string `json:"passport_number" db:"passport_number"`
	PassSerie  string `json:"pass_serie" db:"pass_serie"`
	Surname    string `json:"surname" db:"surname"`
	Name       string `json:"name" db:"name"`
	Patronymic string `json:"patronymic" db:"patronymic"`
	Address    string `json:"address" db:"address"`
}

type GetUsersRequest struct {
	PassportNumber string `form:"passport_number"`
	PassSerie      string `form:"pass_serie"`
	Surname        string `form:"surname"`
	Name           string `form:"name"`
	Patronymic     string `form:"patronymic"`
	Address        string `form:"address"`
	Page           int    `form:"page" binding:"required"`
	PageSize       int    `form:"page_size" binding:"required"`
}
