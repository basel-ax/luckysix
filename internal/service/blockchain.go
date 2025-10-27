package service

import (
	"context"
	"encoding/hex"
	"io"
	"math"
	"math/big"

	"github.com/basel-ax/luckysix/config"
	log "github.com/basel-ax/luckysix/pkg/logger"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// BlockchainService 4 ETH-like blockchains.
type BlockchainService struct {
	log *log.Logger
	cfg *config.Config
}

// NewBlockchainService -.
func NewBlockchainService(l *log.Logger, cfg *config.Config) *BlockchainService {
	return &BlockchainService{
		log: l,
		cfg: cfg,
	}
}

func (service *BlockchainService) GetBlockchainClient(ctx context.Context, blockchainUrl string) (*ethclient.Client, error) {
	url := blockchainUrl
	client, err := ethclient.DialContext(ctx, url)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// based on https://blog.logrocket.com/ethereum-development-using-go-ethereum/ & https://goethereumbook.org/account-balance/
func (service *BlockchainService) GetEthLikeWalletBlance(ctx context.Context, blockchainUrl string, walletAddress string) (*big.Float, error) {
	url := blockchainUrl
	client, err := service.GetBlockchainClient(ctx, url)
	if err != nil {
		return nil, err
	}

	block, errGetLatestBlock := service.GetEthLikeLatestBlock(ctx, client)
	if errGetLatestBlock != nil {
		return nil, errGetLatestBlock
	}
	service.log.Info("Check blockchain Network, last block number: %w", block.Number())

	address := common.HexToAddress(walletAddress)
	balance, errBalance := client.BalanceAt(ctx, address, nil)
	if errBalance != nil {
		return nil, errBalance
	}

	//return result in *big.Int. Exmple - 7115893362963885266, real balance - 7.115893362963885266 MATIC in Polygon
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	floatValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	return floatValue, nil
}

// verify your connection to the ETH-like block node by querying for the current block number of the Polygon blockchain
func (service *BlockchainService) GetEthLikeLatestBlock(ctx context.Context, client *ethclient.Client) (*types.Block, error) {
	block, err := client.BlockByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}
	return block, nil
}

// based on https://github.com/huahuayu/go-transaction-decoder/blob/master/main.go
func (service *BlockchainService) DecodeTxInput(abiReader io.Reader, txInput string) (map[string]interface{}, error) {
	// load contract ABI
	abi, err := abi.JSON(abiReader)
	if err != nil {
		return nil, err
	}

	// decode txInput method signature
	decodedSig, err := hex.DecodeString(txInput[2:10])
	if err != nil {
		return nil, err
	}

	// recover Method from signature and ABI
	method, err := abi.MethodById(decodedSig)
	if err != nil {
		return nil, err
	}

	// decode txInput Payload
	decodedData, err := hex.DecodeString(txInput[10:])
	if err != nil {
		return nil, err
	}

	// unpack method inputs
	inputMap := make(map[string]interface{}, 0)
	err = method.Inputs.UnpackIntoMap(inputMap, decodedData)
	if err != nil {
		return nil, err
	}

	return inputMap, nil
}
package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/alexvec/go-bip39"
	"github.com/basel-ax/luckysix/internal/entity"
	"github.com/basel-ax/luckysix/pkg/logger"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"gorm.io/gorm"
)

// BlockchainService -.
type BlockchainService struct {
	luckyTwoRepo      LuckyTwoRepo
	luckyFiveRepo     LuckyFiveRepo
	luckySixRepo      LuckySixRepo
	walletBalanceRepo WalletBalanceRepo
	logger            logger.Interface
}

// NewBlockchainService -.
func NewBlockchainService(l2r LuckyTwoRepo, l5r LuckyFiveRepo, l6r LuckySixRepo, wbr WalletBalanceRepo, l logger.Interface) *BlockchainService {
	return &BlockchainService{
		luckyTwoRepo:      l2r,
		luckyFiveRepo:     l5r,
		luckySixRepo:      l6r,
		walletBalanceRepo: wbr,
		logger:            l,
	}
}

// ProcessLuckySixCombinations finds unprocessed LuckySix combinations, generates wallets, and stores them.
func (s *BlockchainService) ProcessLuckySixCombinations(ctx context.Context, count int) error {
	s.logger.Info("Starting to process LuckySix combinations for wallet generation")

	lastID, err := s.walletBalanceRepo.GetLastProcessedLuckySixID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get last processed LuckySix ID: %w", err)
	}
	s.logger.Info("Starting processing from LuckySix ID", "last_id", lastID)

	luckySixes, err := s.luckySixRepo.GetBatchStartingFromID(ctx, lastID, count)
	if err != nil {
		return fmt.Errorf("failed to get batch of LuckySixes: %w", err)
	}

	if len(luckySixes) == 0 {
		s.logger.Info("No new LuckySix combinations to process.")
		return nil
	}

	walletBalances := make([]entity.WalletBalance, 0, len(luckySixes))

	for _, ls := range luckySixes {
		mnemonic, err := s.buildMnemonic(ctx, ls)
		if err != nil {
			s.logger.Error(fmt.Errorf("could not build mnemonic for LuckySix ID %d: %w", ls.ID, err))
			continue // Skip to the next one
		}

		address, err := s.deriveAddress(mnemonic)
		if err != nil {
			s.logger.Error(fmt.Errorf("could not derive address for LuckySix ID %d: %w", ls.ID, err))
			continue
		}

		walletBalances = append(walletBalances, entity.WalletBalance{
			LuckySixID: ls.ID,
			Mnemonic:   mnemonic,
			Address:    address,
			Balance:    "0", // Placeholder for now
			Model: gorm.Model{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		})
	}

	if err := s.walletBalanceRepo.StoreBatch(ctx, walletBalances); err != nil {
		return fmt.Errorf("failed to store wallet balances batch: %w", err)
	}

	s.logger.Info("Successfully processed and stored new wallet balances", "count", len(walletBalances))
	return nil
}

func (s *BlockchainService) buildMnemonic(ctx context.Context, ls entity.LuckySix) (string, error) {
	lf, err := s.luckyFiveRepo.GetByID(ctx, ls.LuckyFiveID)
	if err != nil {
		return "", err
	}

	// Collect all six LuckyTwo IDs
	luckyTwoIDs := []uint{
		lf.PairOne, lf.PairTwo, lf.PairThree, lf.PairFour, lf.PairFive,
		ls.LuckyTwoID,
	}

	luckyTwoMap, err := s.luckyTwoRepo.GetByIDs(ctx, luckyTwoIDs)
	if err != nil {
		return "", err
	}

	var words []string
	wordList := bip39.WordList
	// Reconstruct the 12 words in the correct order
	for _, id := range luckyTwoIDs {
		lt, ok := luckyTwoMap[id]
		if !ok {
			return "", fmt.Errorf("missing LuckyTwo data for ID %d", id)
		}
		words = append(words, wordList[lt.WordOne], wordList[lt.WordTwo])
	}

	return strings.Join(words, " "), nil
}

func (s *BlockchainService) deriveAddress(mnemonic string) (string, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return "", err
	}

	// Standard Ethereum derivation path
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		return "", err
	}

	address, err := wallet.Address(account)
	if err != nil {
		return "", err
	}

	return address.Hex(), nil
}
package service

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/alexvec/go-bip39"
	"github.com/basel-ax/luckysix/config"
	"github.com/basel-ax/luckysix/internal/entity"
	"github.com/basel-ax/luckysix/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// BlockchainService -.
type BlockchainService struct {
	luckyTwoRepo      LuckyTwoRepo
	luckyFiveRepo     LuckyFiveRepo
	luckySixRepo      LuckySixRepo
	walletBalanceRepo WalletBalanceRepo
	logger            logger.Interface
	cfg               *config.Config
}

// NewBlockchainService -.
func NewBlockchainService(l2r LuckyTwoRepo, l5r LuckyFiveRepo, l6r LuckySixRepo, wbr WalletBalanceRepo, l logger.Interface, cfg *config.Config) *BlockchainService {
	return &BlockchainService{
		luckyTwoRepo:      l2r,
		luckyFiveRepo:     l5r,
		luckySixRepo:      l6r,
		walletBalanceRepo: wbr,
		logger:            l,
		cfg:               cfg,
	}
}

// ProcessLuckySixCombinations finds unprocessed LuckySix combinations, generates wallets, and stores them.
func (s *BlockchainService) ProcessLuckySixCombinations(ctx context.Context, count int) error {
	s.logger.Info("Starting to process LuckySix combinations for wallet generation")

	lastID, err := s.walletBalanceRepo.GetLastProcessedLuckySixID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get last processed LuckySix ID: %w", err)
	}
	s.logger.Info("Starting processing from LuckySix ID", "last_id", lastID)

	luckySixes, err := s.luckySixRepo.GetBatchStartingFromID(ctx, lastID, count)
	if err != nil {
		return fmt.Errorf("failed to get batch of LuckySixes: %w", err)
	}

	if len(luckySixes) == 0 {
		s.logger.Info("No new LuckySix combinations to process.")
		return nil
	}

	walletBalances := make([]entity.WalletBalance, 0, len(luckySixes))

	for _, ls := range luckySixes {
		mnemonic, err := s.buildMnemonic(ctx, ls)
		if err != nil {
			s.logger.Error(fmt.Errorf("could not build mnemonic for LuckySix ID %d: %w", ls.ID, err))
			continue // Skip to the next one
		}

		address, err := s.deriveAddress(mnemonic)
		if err != nil {
			s.logger.Error(fmt.Errorf("could not derive address for LuckySix ID %d: %w", ls.ID, err))
			continue
		}

		walletBalances = append(walletBalances, entity.WalletBalance{
			LuckySixID: ls.ID,
			Mnemonic:   mnemonic,
			Address:    address,
			Balance:    "0", // Placeholder for now
			Model: gorm.Model{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		})
	}

	if err := s.walletBalanceRepo.StoreBatch(ctx, walletBalances); err != nil {
		return fmt.Errorf("failed to store wallet balances batch: %w", err)
	}

	s.logger.Info("Successfully processed and stored new wallet balances", "count", len(walletBalances))
	return nil
}

func (s *BlockchainService) buildMnemonic(ctx context.Context, ls entity.LuckySix) (string, error) {
	lf, err := s.luckyFiveRepo.GetByID(ctx, ls.LuckyFiveID)
	if err != nil {
		return "", err
	}

	luckyTwoIDs := []uint{
		lf.PairOne, lf.PairTwo, lf.PairThree, lf.PairFour, lf.PairFive,
		ls.LuckyTwoID,
	}

	luckyTwoMap, err := s.luckyTwoRepo.GetByIDs(ctx, luckyTwoIDs)
	if err != nil {
		return "", err
	}

	var words []string
	wordList := bip39.WordList
	for _, id := range luckyTwoIDs {
		lt, ok := luckyTwoMap[id]
		if !ok {
			return "", fmt.Errorf("missing LuckyTwo data for ID %d", id)
		}
		words = append(words, wordList[lt.WordOne], wordList[lt.WordTwo])
	}

	return strings.Join(words, " "), nil
}

func (s *BlockchainService) deriveAddress(mnemonic string) (string, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return "", err
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		return "", err
	}

	address, err := wallet.Address(account)
	if err != nil {
		return "", err
	}

	return address.Hex(), nil
}

// UpdateWalletBalances checks balances for wallets with a zero balance and updates them.
func (s *BlockchainService) UpdateWalletBalances(ctx context.Context, limit int) error {
	s.logger.Info("Starting to update wallet balances")

	wallets, err := s.walletBalanceRepo.GetWalletsForBalanceCheck(ctx, limit)
	if err != nil {
		return fmt.Errorf("failed to get wallets for balance check: %w", err)
	}

	if len(wallets) == 0 {
		s.logger.Info("No wallets with zero balance to update.")
		return nil
	}

	rpcURLs := []string{
		s.cfg.Blockchain.EthMainnetRPC,
		s.cfg.Blockchain.ArbitrumRPC,
		s.cfg.Blockchain.BaseRPC,
		s.cfg.Blockchain.BnbRPC,
	}

	var updatedWallets []entity.WalletBalance

	for _, wallet := range wallets {
		totalBalance := decimal.Zero
		for _, url := range rpcURLs {
			if url == "" {
				continue
			}

			client, err := ethclient.DialContext(ctx, url)
			if err != nil {
				s.logger.Error(fmt.Errorf("failed to connect to blockchain node at %s: %w", url, err))
				continue
			}

			balanceWei, err := client.BalanceAt(ctx, common.HexToAddress(wallet.Address), nil)
			if err != nil {
				s.logger.Error(fmt.Errorf("failed to get balance for address %s from %s: %w", wallet.Address, url, err))
				client.Close()
				continue
			}
			client.Close()

			balanceEther := weiToEther(balanceWei)
			totalBalance = totalBalance.Add(balanceEther)
		}

		if totalBalance.GreaterThan(decimal.Zero) {
			now := time.Now()
			updatedWallet := entity.WalletBalance{
				Model:            gorm.Model{ID: wallet.ID},
				Balance:          totalBalance.String(),
				BalanceUpdatedAt: &now,
			}
			updatedWallets = append(updatedWallets, updatedWallet)
		}
	}

	if len(updatedWallets) > 0 {
		if err := s.walletBalanceRepo.UpdateBalances(ctx, updatedWallets); err != nil {
			return fmt.Errorf("failed to update wallet balances in DB: %w", err)
		}
		s.logger.Info("Successfully updated balances for wallets", "count", len(updatedWallets))
	} else {
		s.logger.Info("Found no wallets with a non-zero balance in this batch.")
	}

	return nil
}

// weiToEther converts a wei value (*big.Int) to an ether value (decimal.Decimal).
func weiToEther(wei *big.Int) decimal.Decimal {
	if wei == nil {
		return decimal.Zero
	}
	weiDecimal := decimal.NewFromBigInt(wei, 0)
	etherMultiplier := decimal.NewFromInt(10).Pow(decimal.NewFromInt(18))
	return weiDecimal.Div(etherMultiplier)
}
