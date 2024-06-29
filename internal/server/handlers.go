package server

import (
	"log"
	"net/http"

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
		log.Println("err in get users: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err})
	}

	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

func (s *Server) getUserTasksHandler(ctx *gin.Context) {

}

func (s *Server) createUserHandler(ctx *gin.Context) {

}

func (s *Server) deleteUserHandler(ctx *gin.Context) {

}

func (s *Server) updateUserHandler(ctx *gin.Context) {

}

func (s *Server) startTaskHandler(ctx *gin.Context) {

}

func (s *Server) stopTaskHandler(ctx *gin.Context) {

}
