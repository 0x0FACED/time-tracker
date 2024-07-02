package server

import (
	"time-tracker/configs"
	"time-tracker/internal/storage"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	cfg    configs.ServerConfig
	r      *gin.Engine
	db     storage.Database
	logger *zap.SugaredLogger
}

func New(cfg configs.ServerConfig, db storage.Database, log *zap.SugaredLogger) *Server {
	return &Server{
		r:      gin.Default(),
		cfg:    cfg,
		db:     db,
		logger: log,
	}
}

func (s *Server) Start() error {
	s.prepare()
	addr := s.cfg.Host + ":" + s.cfg.Port
	s.logger.Infow("server ready to start", "addr", addr)
	s.r.Run(addr)
	return nil
}

func (s *Server) prepare() {
	s.prepareRoutes()
	// ... //
}
