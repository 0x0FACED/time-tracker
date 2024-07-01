package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type user struct {
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	Patronymic string `json:"patronymic"`
	Address    string `json:"address"`
}

func main() {
	r := gin.Default()

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
		if number[0] == "1234" && number[1] == "567890" {
			ctx.JSON(http.StatusOK, user{
				Surname:    "Иванов",
				Name:       "Иван",
				Patronymic: "Иванович",
				Address:    "город. Москва, ул. Ленина, д. 5, кв. 1",
			})
		} else {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		}
	})

	r.Run(":8081")
}
