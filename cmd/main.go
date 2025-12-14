package main

import (
	"context"
	"os"

	"access-system-api/internal/cfg"
	"access-system-api/internal/client"
	"access-system-api/internal/handler"
	"access-system-api/internal/repository"
	"access-system-api/internal/router"
	"access-system-api/internal/service"

	"github.com/sirupsen/logrus"
)

func main() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05Z",
	})

	dbCfg, err := cfg.LoadDbCfg()
	if err != nil {
		log.Fatalf("Error while loading db config: %s", err.Error())
	}
	log.Info("DB config loaded successfully")

	db, err := client.ConnectDB(dbCfg)
	if err != nil {
		log.Fatalf("Error while db connection: %s", err.Error())
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Errorf("Error while closing db connection: %s", err.Error())
		}
	}()
	log.Info("DB connection successful")

	embeddingRepo := repository.NewEmbeddingsRepository(db)
	log.Info("Repository initialized successfully")

	embeddingService := service.NewEmbeddingService(embeddingRepo)
	log.Info("Service initialized successfully")

	v1Handler := handler.NewV1Handler(embeddingService, log)
	log.Info("Handler initialized successfully")

	adminHandler := handler.NewAdminHandler(embeddingService, log)
	log.Info("Admin Handler initialized successfully")

	r := router.NewRouter(v1Handler, adminHandler, log)
	r.Run()
	log.Info("Router started successfully")
}
