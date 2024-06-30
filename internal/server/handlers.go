package server

import (
	"log"
	"net/http"
	"strconv"
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
	s.r.Handle(http.MethodPost, "/tasks/stop", s.stopTaskHandler)
}

func (s *Server) getUsersHandler(ctx *gin.Context) {
	users, err := s.db.GetUsers()
	if err != nil {
		log.Println("get users err: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}
	if len(users) == 0 {
		log.Println(errors.ErrNoUsersFound)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": errors.ErrNoUsersFound})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

func (s *Server) getUserTasksHandler(ctx *gin.Context) {

}

func (s *Server) createUserHandler(ctx *gin.Context) {
	var input struct {
		passportNumber string `json:"pass_number"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// fetching fata from external api //
	// .......
	// creating new user with pass_number and fetched data from external api
	// .......

	// TODO: refactor
	err := s.db.AddUser(&models.User{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err})
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
		ctx.JSON(http.StatusBadRequest, gin.H{"res": "cant delete", "err": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"res": "succesfully deleted"})
}

func (s *Server) updateUserHandler(ctx *gin.Context) {

}

func (s *Server) startTaskHandler(ctx *gin.Context) {

}

func (s *Server) stopTaskHandler(ctx *gin.Context) {

}
