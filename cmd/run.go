package cmd

import (
	"time-tracker/configs"
	"time-tracker/internal/server"
	"time-tracker/internal/storage"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Run() {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	config.OutputPaths = []string{"stdout", "./logs.txt"}

	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.CallerKey = "caller"
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	cfg, err := configs.Load()
	if err != nil {
		sugar.Fatalln("cant load config, error: ", err)
		return
	}
	sugar.Infow("config loaded",
		"config", cfg,
	)
	db := storage.New(cfg.DatabaseConfig)

	err = db.Connect()
	if err != nil {
		sugar.Fatalln("cant connect to db, error: ", err)
		return
	}
	sugar.Infoln("successfully connected to db")
	s := server.New(cfg.ServerConfig, db, sugar)
	sugar.Infoln("starting the server")
	err = s.Start()
	if err != nil {
		sugar.Fatalw("cant start the server",
			"error", err,
		)
		return
	}
}
