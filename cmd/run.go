package cmd

import (
	"log"
	"time-tracker/configs"
	"time-tracker/internal/server"
	"time-tracker/internal/storage"
)

func Run() {
	cfg, err := configs.Load()
	if err != nil {
		log.Fatalln("cant load cfg file: ", err)
		return
	}
	db := storage.New(cfg.DatabaseConfig)
	s := server.New(cfg.ServerConfig, db)
	err = s.Start()
	if err != nil {
		log.Fatalln("cant start the server: ", err)
		return
	}
}
