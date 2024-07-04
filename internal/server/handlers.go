package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time-tracker/internal/models"
	"time-tracker/internal/utils"

	"github.com/gin-gonic/gin"
)

func (s *Server) prepareRoutes() {
	s.r.Handle(http.MethodGet, "/users", s.getUsersHandler)
	s.r.Handle(http.MethodGet, "/users/tasks", s.getUserTasksHandler)
	s.r.Handle(http.MethodPost, "/create", s.createUserHandler)
	s.r.Handle(http.MethodDelete, "/users/:id", s.deleteUserHandler)
	s.r.Handle(http.MethodPut, "/users/:id", s.updateUserHandler)
	s.r.Handle(http.MethodPost, "/tasks/start", s.startTaskHandler)
	s.r.Handle(http.MethodPost, "/tasks/:id/stop", s.stopTaskHandler)
}

func (s *Server) getUsersHandler(ctx *gin.Context) {
	var req models.GetUsersRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		s.logger.Errorw("getUsersHandler", "full path", ctx.FullPath())
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "invalid query parameters"})
		return
	}
	s.logger.Debugw("getUsersHandler", "req", req)
	users, err := s.db.GetUsers(req)
	if err != nil {
		s.logger.Errorw("getUsersHandler", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}
	if len(users) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"err": utils.ErrNoUsersFound})
		return
	}

	s.logger.Infoln("successfully executed GetUsers")
	s.logger.Debugw("getUsersHandler", "users: ", users)
	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

func (s *Server) getUserTasksHandler(ctx *gin.Context) {
	var req models.GetUserWorklogsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		s.logger.Errorw("invalid query parameters", "full path", ctx.FullPath())
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "invalid query parameters"})
		return
	}
	s.logger.Debugw("getUserTasksHandler", "req: ", req)
	worklogs, err := s.db.GetUserWorklogs(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "failed to get worklogs"})
		return
	}
	s.logger.Infoln("successfully executed GetUsers")
	s.logger.Debugw("getUserTasksHandler", "worklogs: ", worklogs)
	ctx.JSON(http.StatusOK, gin.H{"worklogs": worklogs})
}

func (s *Server) createUserHandler(ctx *gin.Context) {
	var input struct {
		PassportNumber string `json:"passport_number"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	s.logger.Debugw("createUserHandler", "passport_number", input.PassportNumber)
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
	s.logger.Debugw("createUserHandler", "resp", resp.Body)
	if resp.StatusCode == http.StatusBadRequest {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "bad request"})
		return
	}
	if resp.StatusCode == http.StatusInternalServerError {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}
	if resp.StatusCode == http.StatusNotFound {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "user with these passport data not found in database"})
		return
	}

	var apiResponse models.User
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		s.logger.Errorw("createUserHandler", "cant decode resp", resp.Body)
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
	s.logger.Debugw("createUserHandler", "full user", newUser)
	u, err := s.db.AddUser(newUser)
	if err != nil {
		s.logger.Errorln("cant add user to db, error: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	s.logger.Infow("user successfully added to db", "user", u)
	ctx.JSON(http.StatusOK, u)
}

func (s *Server) deleteUserHandler(ctx *gin.Context) {
	id_param := ctx.Param("id")
	s.logger.Debugw("deleteUserHandler", "id", id_param)
	id, err := strconv.Atoi(id_param)
	if err != nil {
		s.logger.Debugw("deleteUserHandler", "id is not a number", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "id is not a number"})
		return
	}
	err = s.db.DeleteUser(id)
	if err != nil {
		s.logger.Debugw("deleteUserHandler", "cant delete user with err", err)
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
	id_param := ctx.Param("id")
	s.logger.Debugw("updateUserHandler", "id", id_param, "newUser", updateUserInput)
	id, err := strconv.Atoi(id_param)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
	}
	u, err := s.db.GetUserByID(id)
	s.logger.Debugw("getUserByID", "user", u)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Debugln("user not found")
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
	s.logger.Debugln("full new user info", "user", u)
	err = s.db.UpdateUser(u)
	if err != nil {
		s.logger.Errorln("cant update user, error: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	s.logger.Infoln("successfully updated user")
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
	s.logger.Debugw("startTaskHandler", "input_task", input)
	task := &models.Task{
		UserID: input.UserID,
		Desc:   input.Desc,
	}
	err := s.db.AddStartTask(task)
	if err != nil {
		s.logger.Errorln("failed to add task to db, error: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	s.logger.Infoln("successfully start task")
	ctx.JSON(http.StatusOK, gin.H{"res": "task started"})

}

func (s *Server) stopTaskHandler(ctx *gin.Context) {
	id_param := ctx.Param("id")
	s.logger.Debugw("stopTaskHandler", "id", id_param)
	id, err := strconv.Atoi(id_param)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "incorrect task id"})
		return
	}
	err = s.db.AddEndTask(id)
	if err != nil {
		s.logger.Errorln("failed to end task, error: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	s.logger.Infoln("successfully stopped task")
	ctx.JSON(http.StatusOK, gin.H{"res": "task ended"})
}
