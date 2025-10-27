// Package Tasks implements Golang Job Scheduling
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/basel-ax/luckysix/config"
	"github.com/jasonlvhit/gocron"

	log "github.com/basel-ax/luckysix/pkg/logger"
	"github.com/basel-ax/luckysix/pkg/rabbitmq/rmq_rpc/client"
)

type TasksService struct {
	log        *log.Logger
	bip39      Bip39
	blockchain Blockchain
}

func NewTasksService(l *log.Logger, bip39 Bip39, blockchain Blockchain) *TasksService {
	return &TasksService{
		log:        l,
		bip39:      bip39,
		blockchain: blockchain,
	}
}

func (service *TasksService) StartTasks(cfg *config.Config) {
	service.log.Info("StartTasks")

	gocron.Clear()

	//gocron.Every(1).Minute().From(gocron.NextTick()).Do(service.EveryMinuteTask, cfg)
	gocron.Every(10).Minute().From(gocron.NextTick()).Do(service.CheckWalletBalances)
	//gocron.Every(10).Minute().From(gocron.NextTick()).Do(service.EveryTenMinuteTask, cfg)
	gocron.Every(24).Hours().From(gocron.NextTick()).Do(service.EveryDayTask, cfg)

	<-gocron.Start()
}

func (service *TasksService) EveryDayTask(cfg *config.Config) {
	service.log.Info("Start EveryDayTask")

	//ctx := context.Background()
	/*

	 */

	service.log.Info("End EveryDayTask")
}

func (service *TasksService) EveryMinuteTask(cfg *config.Config) {
	service.log.Info("Start EveryMinuteTasks")

	service.log.Info("End everyMinuteTasks")
}

func (service *TasksService) EveryTenMinuteTask(cfg *config.Config) {
	service.log.Info("Start EveryTenMinuteTask")

	//service.Flow(cfg);

	service.log.Info("End everyMinuteTasks")
}

func (service *TasksService) CheckRabbit(cfg *config.Config, ctx context.Context) {
	//Test RabbitMQ
	rmqClient, err := client.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, cfg.RMQ.ClientExchange)
	if err != nil {
		service.log.Fatal("RabbitMQ RPC Client - init error - client.New")
	}
	defer func() {
		err = rmqClient.Shutdown()
		if err != nil {
			service.log.Fatal("RabbitMQ RPC Client - shutdown error - rmqClient.RemoteCall", err)
		}
	}()
	var answer struct {
		Message string `json:"message"`
	}

	//TODO: fix
	/*
		"message":"RabbitMQ RPC Client - remote call error - rmqClient.RemoteCall(checkRabbit)%!(EXTRA *fmt.wrapError=rmq_rpc client - Client - RemoteCall - json.Unmarshal: json: cannot unmarshal object into Go value of type string)"}
	*/
	err = rmqClient.RemoteCall("CheckRabbit", nil, &answer)
	if err != nil {
		service.log.Fatal("RabbitMQ RPC Client - remote call error - rmqClient.RemoteCall(checkRabbit)", err)
	}

	if len(answer.Message) > 0 {
		service.log.Info(answer.Message)
	} else {
		service.log.Fatal("CheckRabbit answer is empty")
	}
}


// Task execute in RabbitMQ controller
func (service *TasksService) CheckRabbitTask() string {
	msg := fmt.Sprintf("CheckRabbitTask, time = %s", time.Now())

	return msg
}

// GenerateLuckyTwo starts the generation of Luckytwo combinations in the background.
func (service *TasksService) GenerateLuckyTwo() {
	service.log.Info("Task started: GenerateLuckyTwo")
	go func() {
		err := service.bip39.GenerateAndStoreLuckyTwoCombinations(context.Background())
		if err != nil {
			service.log.Error(fmt.Errorf("GenerateLuckyTwo task failed: %w", err))
		} else {
			service.log.Info("Task finished successfully: GenerateLuckyTwo")
		}
	}()
}

// ProcessWallets finds unprocessed LuckySix combinations and generates wallet data for them.
func (service *TasksService) ProcessWallets() {
	const combinationsToProcess = 100 // Process 100 at a time to avoid long-running tasks

	service.log.Info("Task started: ProcessWallets", "count", combinationsToProcess)
	go func() {
		err := service.blockchain.ProcessLuckySixCombinations(context.Background(), combinationsToProcess)
		if err != nil {
			service.log.Error(fmt.Errorf("ProcessWallets task failed: %w", err))
		} else {
			service.log.Info("Task finished successfully: ProcessWallets")
		}
	}()
}

// GenerateLuckyFive starts the generation of a batch of LuckyFive combinations.
func (service *TasksService) GenerateLuckyFive() {
	const combinationsToGenerate = 10000 // Define how many to generate per task run

	service.log.Info("Task started: GenerateLuckyFive", "count", combinationsToGenerate)
	go func() {
		err := service.bip39.GenerateAndStoreLuckyFiveCombinations(context.Background(), combinationsToGenerate)
		if err != nil {
			service.log.Error(fmt.Errorf("GenerateLuckyFive task failed: %w", err))
		} else {
			service.log.Info("Task finished successfully: GenerateLuckyFive")
		}
	}()
}

// GenerateLuckySix starts the generation of a batch of LuckySix combinations.
func (service *TasksService) GenerateLuckySix() {
	const combinationsToGenerate = 100000 // Define how many to generate per task run

	service.log.Info("Task started: GenerateLuckySix", "count", combinationsToGenerate)
	go func() {
		err := service.bip39.GenerateAndStoreLuckySixCombinations(context.Background(), combinationsToGenerate)
		if err != nil {
			service.log.Error(fmt.Errorf("GenerateLuckySix task failed: %w", err))
		} else {
			service.log.Info("Task finished successfully: GenerateLuckySix")
		}
	}()
}

// CheckWalletBalances starts the process of checking and updating wallet balances.
func (service *TasksService) CheckWalletBalances() {
	const walletsToCheck = 50 // Check 50 wallets per run to avoid hitting RPC rate limits

	service.log.Info("Task started: CheckWalletBalances", "count", walletsToCheck)
	go func() {
		err := service.blockchain.UpdateWalletBalances(context.Background(), walletsToCheck)
		if err != nil {
			service.log.Error(fmt.Errorf("CheckWalletBalances task failed: %w", err))
		} else {
			service.log.Info("Task finished successfully: CheckWalletBalances")
		}
	}()
}
