package server

import (
	"time-tracker/configs"
	"time-tracker/internal/storage"

	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg configs.ServerConfig
	r   *gin.Engine
	db  storage.Database
}

func New(cfg configs.ServerConfig, db storage.Database) *Server {
	return &Server{
		r:   gin.Default(),
		cfg: cfg,
		db:  db,
	}
}

func (s *Server) Start() error {
	s.prepare()
	addr := s.cfg.Host + ":" + s.cfg.Port
	s.r.Run(addr)
	return nil
}

func (s *Server) prepare() {
	s.prepareRoutes()
	// ... //
}
