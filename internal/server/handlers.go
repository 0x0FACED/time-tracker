package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time-tracker/internal/models"
	"time-tracker/internal/utils/errors"

	"github.com/gin-gonic/gin"
)

func (s *Server) prepareRoutes() {
	s.r.Handle(http.MethodGet, "/users", s.getUsersHandler)
	s.r.Handle(http.MethodGet, "/users/:id/tasks", s.getUserTasksHandler)
	s.r.Handle(http.MethodPost, "/users", s.createUserHandler)
	s.r.Handle(http.MethodDelete, "/users/:id", s.deleteUserHandler)
	s.r.Handle(http.MethodPut, "/users/:id", s.updateUserHandler)
	s.r.Handle(http.MethodPost, "/tasks/start", s.startTaskHandler)
	s.r.Handle(http.MethodPost, "/tasks/:id/stop", s.stopTaskHandler)
}

func (s *Server) getUsersHandler(ctx *gin.Context) {
	var req models.GetUsersRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "invalid query parameters"})
		return
	}
	query := "SELECT id, passport_number, pass_serie, surname, name, patronymic, address FROM users WHERE 1=1"
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
	users, err := s.db.GetUsers(query, params...)
	if err != nil {
		log.Println("get users err: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}
	if len(users) == 0 {
		log.Println(errors.ErrNoUsersFound)
		ctx.JSON(http.StatusNotFound, gin.H{"err": errors.ErrNoUsersFound})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

func (s *Server) getUserTasksHandler(ctx *gin.Context) {

}

func (s *Server) createUserHandler(ctx *gin.Context) {
	var input struct {
		PassportNumber string `json:"passport_number"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	log.Println("PASSPORT: ", input.PassportNumber)

	reqBody, err := json.Marshal(map[string]string{
		"passport_number": input.PassportNumber,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "failed to marshal request"})
		return
	}
	// external api call (just for example)
	resp, err := http.Post("http://localhost:8081/info", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "failed to call external API"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("external API returned error: %s", body)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "failed to get user info from external API"})
		return
	}

	var apiResponse models.User
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "failed to decode API response"})
		return
	}
	passport := strings.Fields(input.PassportNumber)
	newUser := &models.User{
		PassNumber: passport[1],
		PassSerie:  passport[0],
		Surname:    apiResponse.Surname,
		Name:       apiResponse.Name,
		Patronymic: apiResponse.Patronymic,
		Address:    apiResponse.Address,
	}

	err = s.db.AddUser(newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"res": "created"})
}

func (s *Server) deleteUserHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Fatalln("id is not a number: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "id is not a number"})
		return
	}
	err = s.db.DeleteUser(idInt)
	if err != nil {
		log.Println("cant delete user: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"res": "cant delete", "err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"res": "successfully deleted"})
}

func (s *Server) updateUserHandler(ctx *gin.Context) {
	var updateUserInput struct {
		Surname    *string `json:"surname,omitempty"`
		Name       *string `json:"name,omitempty"`
		Patronymic *string `json:"patronymic,omitempty"`
		Address    *string `json:"address,omitempty"`
	}
	if err := ctx.ShouldBindJSON(&updateUserInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
	}
	u, err := s.db.GetUserByID(idInt)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"err": "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "server error"})
		log.Println(err)
		return
	}

	if updateUserInput.Surname != nil {
		u.Surname = *updateUserInput.Surname
	}
	if updateUserInput.Name != nil {
		u.Name = *updateUserInput.Name
	}
	if updateUserInput.Patronymic != nil {
		u.Patronymic = *updateUserInput.Patronymic
	}
	if updateUserInput.Address != nil {
		u.Address = *updateUserInput.Address
	}

	err = s.db.UpdateUser(u)
	if err != nil {
		log.Fatalln("cant update user: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"res": "successfully updated"})
}

func (s *Server) startTaskHandler(ctx *gin.Context) {
	var input struct {
		UserID int    `json:"user_id"`
		Desc   string `json:"description"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	task := &models.Task{
		UserID: input.UserID,
		Desc:   input.Desc,
	}
	err := s.db.AddStartTask(task)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "cant start task, try again"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"res": "task started"})

}

func (s *Server) stopTaskHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "incorrect task id"})
		return
	}
	err = s.db.AddEndTask(idInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"res": "task ended"})
}
