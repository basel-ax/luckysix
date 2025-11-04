package app

import (
	"context"
	"fmt"
	"os"

	"github.com/basel-ax/luckysix/internal/config"
	"github.com/basel-ax/luckysix/internal/repository"
	"github.com/basel-ax/luckysix/internal/service"
	"github.com/basel-ax/luckysix/pkg/db"
	"github.com/basel-ax/luckysix/pkg/logger"
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
	pg, err := db.New(a.cfg.PG.URL)
	if err != nil {
		return fmt.Errorf("app - Run - db.New: %w", err)
	}

	blockchainRepo := repository.NewBlockchainRepo(pg)
	luckyTwoRepo := repository.NewLuckyTwoRepo(pg)
	luckyFiveRepo := repository.NewLuckyFiveRepo(pg)
	luckySixRepo := repository.NewLuckySixRepo(pg)
	walletBalanceRepo := repository.NewWalletBalanceRepo(pg)

	// Use cases
	walletService := service.NewWalletService(luckyTwoRepo, luckyFiveRepo)
	notificationService := service.NewLogNotifier()
	luckySixService := service.NewLuckySixService(
		blockchainRepo,
		luckyTwoRepo,
		luckyFiveRepo,
		luckySixRepo,
		walletBalanceRepo,
		walletService,
		notificationService,
	)

	// Run
	l.Info("app - Run - Application started")
	go luckySixService.GenerateAndProcess(context.Background())

	// Shutdown
	// Your graceful shutdown logic here...
	interrupt := make(chan os.Signal, 1)
	<-interrupt
	l.Info("app - Run - Application gracefully shutdown")

	return nil
}
