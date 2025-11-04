// Package app configures and runs application.
package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/basel-ax/luckysix/config"
	amqprpc "github.com/basel-ax/luckysix/internal/controller/amqp_rpc"
	"github.com/basel-ax/luckysix/internal/controller/http/v1"
	"github.com/basel-ax/luckysix/internal/service"
	"github.com/basel-ax/luckysix/internal/service/repo"
	"github.com/basel-ax/luckysix/internal/service/webapi"
	"github.com/basel-ax/luckysix/pkg/httpserver"
	"github.com/basel-ax/luckysix/pkg/logger"
	"github.com/basel-ax/luckysix/pkg/postgres"
	"github.com/basel-ax/luckysix/pkg/rabbitmq/rmq_rpc/server"
	"github.com/gin-gonic/gin"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	// Logger
	l := logger.New(context.Background(), cfg.Log.Level, cfg.PARAM.TgBotApi, cfg.PARAM.TgChatId, "", "")

	// Repositories
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	translationRepo := repo.NewTranslationRepo(pg)
	luckyTwoRepo := repo.NewLuckyTwoRepo(pg)
	luckyFiveRepo := repo.NewLuckyFiveRepo(pg)
	luckySixRepo := repo.NewLuckySixRepo(pg)
	walletBalanceRepo := repo.NewWalletBalanceRepo(pg)

	// Services
	telegramService, err := service.NewTelegramService(cfg)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - service.NewTelegramService: %w", err))
	}

	translationService := service.NewTranslationService(
		translationRepo,
		webapi.New(),
	)

	bip39Service := service.NewBip39Service(
		luckyTwoRepo,
		luckyFiveRepo,
		luckySixRepo,
		l,
	)

	blockchainService := service.NewBlockchainService(
		luckyTwoRepo,
		luckyFiveRepo,
		luckySixRepo,
		walletBalanceRepo,
		telegramService,
		l,
		cfg,
	)

	// This is the fix for the compiler error: NewTasksService now receives all its dependencies.
	tasksService := service.NewTasksService(l, bip39Service, blockchainService)

	// Service container
	services := &service.Services{
		Translation: translationService,
		Tasks:       tasksService,
		Bip39:       bip39Service,
		Blockchain:  blockchainService,
		Telegram:    telegramService,
	}

	// AMQP RPC Server
	rmqRouter := amqprpc.NewRouter(services)
	rmqServer, err := server.New(
		cfg.RMQ.URL,
		cfg.RMQ.ServerExchange,
		rmqRouter,
		l,
	)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - rmqServer - server.New: %w", err))
	}

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, services)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Start background tasks
	tasksService.StartTasks(cfg)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err = <-rmqServer.Notify():
		l.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	err = rmqServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - rmqServer.Shutdown: %w", err))
	}
}
