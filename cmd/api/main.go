package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	accountApp "transaction/internal/account/application"
	accountDomain "transaction/internal/account/domain"
	accountInfra "transaction/internal/account/infrastructure"
	"transaction/internal/http"
	accountHandler "transaction/internal/http/handler/account"
	userHandler "transaction/internal/http/handler/user"
	"transaction/internal/user/application"
	"transaction/internal/user/infrastructure"
	"transaction/pkg/config"
	"transaction/pkg/logger"
	"transaction/pkg/migrator"
	"transaction/pkg/postgres"
	"transaction/pkg/redis"
	"transaction/pkg/tigerbeetle"
)

func main() {
	cfg := config.Load()

	logger.Init(cfg.Logger)

	if cfg.Migration.Enabled {
		m := migrator.New(cfg.Database)

		switch cfg.Migration.Direction {
		case "up":
			logger.GetLogger().Info("Running migrations up...")
			if err := m.Up(); err != nil {
				log.Fatalf("Migration failed: %v", err)
			}
			logger.GetLogger().Info("Migrations completed")

		case "down":
			logger.GetLogger().Info("Running migrations down...")
			if err := m.Down(); err != nil {
				log.Fatalf("Migration rollback failed: %v", err)
			}
			logger.GetLogger().Info("Migrations rolled back")
		}
	}

	pgClient, err := postgres.NewClient(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pgClient.Close()
	logger.GetLogger().Info("PostgreSQL connected")

	redisCacheClient, err := redis.NewClient(cfg.Redis, 0)
	if err != nil {
		log.Fatalf("Failed to connect to Redis Cache: %v", err)
	}
	defer redisCacheClient.Close()
	logger.GetLogger().Info("Redis Cache connected (DB: 0)")

	redisLockClient, err := redis.NewClient(cfg.Redis, 1)
	if err != nil {
		log.Fatalf("Failed to connect to Redis Lock: %v", err)
	}
	defer redisLockClient.Close()
	logger.GetLogger().Info("Redis Lock connected (DB: 1)")

	tbClient, err := tigerbeetle.NewClient(cfg.TigerBeetle)
	if err != nil {
		log.Fatalf("Failed to connect to TigerBeetle: %v", err)
	}
	defer tbClient.Close()
	logger.GetLogger().Info("TigerBeetle connected")

	userRepo := infrastructure.NewRepository(pgClient.GetDB())
	apiKeyRepo := infrastructure.NewAPIKeyRepository(pgClient.GetDB())
	userService := application.NewService(userRepo, apiKeyRepo)
	userHdlr := userHandler.NewHandler(userService)

	accountRepo := accountInfra.NewAccountRepository(pgClient.GetDB())
	accountLedger := accountInfra.NewLedger(tbClient)
	accountCache := accountInfra.NewAccountCache(redisCacheClient.GetClient())
	accountService := accountApp.NewService(accountRepo, accountLedger, accountCache)
	accountHdlr := accountHandler.NewHandler(accountService)

	ctx := context.Background()
	if err := accountService.InitializeSystemAccount(ctx, accountDomain.USD, 100000000); err != nil {
		log.Fatalf("Failed to initialize system account: %v", err)
	}
	logger.GetLogger().Info("System account initialized")

	router := http.NewRouter(userHdlr, accountHdlr, userService)
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
