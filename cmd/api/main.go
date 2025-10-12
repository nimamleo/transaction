package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"transaction/internal/http"
	userHandler "transaction/internal/http/handler/user"
	"transaction/internal/user/application"
	"transaction/internal/user/infrastructure"
	"transaction/pkg/config"
	"transaction/pkg/logger"
	"transaction/pkg/postgres"
	"transaction/pkg/redis"
	"transaction/pkg/tigerbeetle"
)

func main() {
	cfg := config.Load()

	logger.Init(cfg.Logger)

	pgClient, err := postgres.NewClient(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pgClient.Close()
	logger.GetLogger().Info("PostgreSQL connected")

	redisClient, err := redis.NewClient(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()
	logger.GetLogger().Info("Redis connected")

	tbClient, err := tigerbeetle.NewClient(cfg.TigerBeetle)
	if err != nil {
		log.Fatalf("Failed to connect to TigerBeetle: %v", err)
	}
	defer tbClient.Close()
	logger.GetLogger().Info("TigerBeetle connected")

	userRepo := infrastructure.NewRepository(pgClient.GetDB())
	userService := application.NewService(userRepo)
	userHdlr := userHandler.NewHandler(userService)

	router := http.NewRouter(userHdlr)
	server := http.NewServer(cfg.Server, router)

	go func() {
		logger.GetLogger().Infof("Server starting on port %s", cfg.Server.Port)
		if err := server.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.GetLogger().Info("Shutting down server...")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal(err)
	}
}
