package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/1boombacks1/testTaskInfotecs/internal/config"
	v1 "github.com/1boombacks1/testTaskInfotecs/internal/controller/http/v1"
	"github.com/1boombacks1/testTaskInfotecs/internal/repo"
	"github.com/1boombacks1/testTaskInfotecs/internal/service"
	_ "github.com/1boombacks1/testTaskInfotecs/migrations"
	uuidgen "github.com/1boombacks1/testTaskInfotecs/pkg/UUIDGen"
	"github.com/1boombacks1/testTaskInfotecs/pkg/postgres"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func Run() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config err: %s", err)
	}

	SetLogrus(cfg.Log.Level)

	log.Info("Initializing postgres...")
	pg, err := postgres.New(cfg.Postgres.DSN)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - pgdb.NewServices: %w", err))
	}
	defer pg.Close()

	log.Info("Initializing repositories...")
	repositories := repo.NewRepos(pg)

	log.Info("Initializing services...")
	deps := service.ServicesDependencies{
		Repos: repositories,
		IDGen: &uuidgen.DefaultGenerator{},
	}
	srvc := service.NewService(deps)

	log.Info("Initializing handlers and routes...")
	router := fiber.New()
	v1.NewRouter(router, srvc)

	serverNotify := make(chan error)
	go func() {
		log.Info("Starting fiber server...")
		serverNotify <- router.Listen(fmt.Sprintf(":%s", cfg.Server.Port))
		close(serverNotify)
	}()

	// Waiting signal
	log.Info("Configuring graceful shutdown...")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-quit:
		log.Info("app - Main - signal: " + s.String())
	case err = <-serverNotify:
		log.Error(fmt.Errorf("app - Main - server notify: %w", err))
	}

	// Graceful shutdown
	log.Info("Gracefully shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = router.ShutdownWithContext(ctx); err != nil {
		log.Error(fmt.Errorf("app - Main - app.ShutdownWithContext: %w", err))
	}
}
