package app

import (
	"fmt"
	"os"

	"github.com//basel-ax/luckysix/internal/config"
	"github.com//basel-ax/luckysix/pkg/db"
	"github.com//basel-ax/luckysix/pkg/logger"
)

type App struct {
	cfg *config.Config
}

func New() (*App, error) {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	return &App{cfg: cfg}, nil
}

func (a *App) Run() error {
	l := logger.New(a.cfg.Log.Level)

	// Repository
	pg, err := db.New(a.cfg.PG.URL, db.MaxPoolSize(a.cfg.PG.MaxPoolSize))
	if err != nil {
		return fmt.Errorf("app - Run - db.New: %w", err)
	}
	defer pg.Close()

	// Use cases
	// Your services here...

	// Run
	// Your services startup here...
	l.Info("app - Run - Application started")

	// Shutdown
	// Your graceful shutdown logic here...
	interrupt := make(chan os.Signal, 1)
	<-interrupt
	l.Info("app - Run - Application gracefully shutdown")

	return nil
}
