package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

type user struct {
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	Patronymic string `json:"patronymic"`
	Address    string `json:"address"`
}

func main() {
	r := gin.Default()
	dbURL := "postgres://postgres:postgres@localhost:5432/external_api_db?sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalln("cant connect to db: ", err)
		return
	}

	if db.Ping() != nil {
		log.Fatalln("cant ping db: ", err)
		return
	}
	err = MigrateUp(dbURL)
	if err != nil {
		log.Fatalln("cant migrate up: ", err)
	}
	r.POST("/info", func(ctx *gin.Context) {
		var passport struct {
			PassportNumber string `json:"passport_number"`
		}
		if err := ctx.ShouldBindJSON(&passport); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"err": "bad request"})
			return
		}
		log.Println("PASSPORT: ", passport.PassportNumber)

		number := strings.Fields(passport.PassportNumber)
		query := `
			SELECT name, surname, patronymic, address 
			FROM users
			WHERE passport_number = $1 AND pass_serie = $2
		`
		var u user
		db.QueryRow(query, number[1], number[0]).Scan(u.Name, u.Surname, u.Patronymic, u.Address)
		ctx.JSON(http.StatusOK, u)
	})

	r.Run(":8081")
}

func MigrateUp(dbURL string) error {
	m, err := migrate.New(
		"file://./migrations/",
		dbURL)
	if err != nil {
		log.Fatalln("failed to create migration: ", err)
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalln("failed to migrate up: ", err)
		return err
	}
	return err
}
