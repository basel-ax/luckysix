package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/basel-ax/luckysix/internal/entity"
	"github.com/basel-ax/luckysix/internal/repository"
)

// LuckySixService -.
type LuckySixService struct {
	blockchainRepo    repository.Blockchain
	luckyTwoRepo      repository.LuckyTwo
	luckyFiveRepo     repository.LuckyFive
	luckySixRepo      repository.LuckySix
	walletBalanceRepo repository.WalletBalance
	walletService     Wallet
	notification      Notification
}

// NewLuckySixService -.
func NewLuckySixService(
	blockchainRepo repository.Blockchain,
	luckyTwoRepo repository.LuckyTwo,
	luckyFiveRepo repository.LuckyFive,
	luckySixRepo repository.LuckySix,
	walletBalanceRepo repository.WalletBalance,
	walletService Wallet,
	notification Notification,
) *LuckySixService {
	return &LuckySixService{
		blockchainRepo:    blockchainRepo,
		luckyTwoRepo:      luckyTwoRepo,
		luckyFiveRepo:     luckyFiveRepo,
		luckySixRepo:      luckySixRepo,
		walletBalanceRepo: walletBalanceRepo,
		walletService:     walletService,
		notification:      notification,
	}
}

// GenerateAndProcess -.
func (s *LuckySixService) GenerateAndProcess(ctx context.Context) {
	// For this example, we'll process all LuckySix combinations.
	// In a real application, you might want to process them in batches.
	luckySixes, err := s.luckySixRepo.GetAll(ctx)
	if err != nil {
		log.Printf("Error getting all lucky sixes: %v", err)
		return
	}

	var wg sync.WaitGroup
	for _, ls := range luckySixes {
		wg.Add(1)
		go func(ls *entity.LuckySix) {
			defer wg.Done()
			s.processLuckySix(ctx, ls)
		}(ls)
	}
	wg.Wait()
}

func (s *LuckySixService) processLuckySix(ctx context.Context, ls *entity.LuckySix) {
	mnemonic, err := s.walletService.GenerateMnemonic(ctx, ls)
	if err != nil {
		log.Printf("Error generating mnemonic for lucky six %d: %v", ls.ID, err)
		return
	}

	address, cosmosAddress, err := s.walletService.GenerateWallet(ctx, mnemonic)
	if err != nil {
		log.Printf("Error generating wallet for lucky six %d: %v", ls.ID, err)
		return
	}

	// For this example, we'll check all enabled blockchains.
	// In a real application, you might want to be more selective.
	blockchains, err := s.blockchainRepo.GetEnabledBlockchains(ctx)
	if err != nil {
		log.Printf("Error getting enabled blockchains: %v", err)
		return
	}

	balance, err := s.walletService.CheckBalance(ctx, address, blockchains)
	if err != nil {
		log.Printf("Error checking balance for lucky six %d: %v", ls.ID, err)
		return
	}

	now := time.Now()
	wb := &entity.WalletBalance{
		LuckySixID:       ls.ID,
		Mnemonic:         mnemonic,
		Address:          address,
		CosmosAddress:    cosmosAddress,
		Balance:          balance,
		BalanceUpdatedAt: &now,
	}

	err = s.walletBalanceRepo.Create(ctx, wb)
	if err != nil {
		log.Printf("Error creating wallet balance for lucky six %d: %v", ls.ID, err)
		return
	}

	if balance != "0" {
		message := fmt.Sprintf("Found wallet with balance! Mnemonic: %s, Address: %s, Balance: %s", mnemonic, address, balance)
		err := s.notification.Notify(ctx, message)
		if err != nil {
			log.Printf("Error sending notification for lucky six %d: %v", ls.ID, err)
		} else {
			wb.IsNotified = true
			err := s.walletBalanceRepo.Update(ctx, wb)
			if err != nil {
				log.Printf("Error updating wallet balance notification status for lucky six %d: %v", ls.ID, err)
			}
		}
	}
}
