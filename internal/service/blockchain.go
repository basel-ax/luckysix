package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/alexvec/go-bip39"
	"github.com/basel-ax/luckysix/config"
	"github.com/basel-ax/luckysix/internal/entity"
	"github.com/basel-ax/luckysix/pkg/logger"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/shopspring/decimal"
	bip39std "github.com/tyler-smith/go-bip39"
	"gorm.io/gorm"
)

func init() {
	config := types.GetConfig()
	config.SetBech32PrefixForAccount("cosmos", "cosmospub")
	config.Seal()
}

// BlockchainService -.
type BlockchainService struct {
	luckyTwoRepo      LuckyTwoRepo
	luckyFiveRepo     LuckyFiveRepo
	luckySixRepo      LuckySixRepo
	walletBalanceRepo WalletBalanceRepo
	telegramSvc       Telegram
	logger            logger.Interface
	cfg               *config.Config
}

// NewBlockchainService -.
func NewBlockchainService(l2r LuckyTwoRepo, l5r LuckyFiveRepo, l6r LuckySixRepo, wbr WalletBalanceRepo, tsv Telegram, l logger.Interface, cfg *config.Config) *BlockchainService {
	return &BlockchainService{
		luckyTwoRepo:      l2r,
		luckyFiveRepo:     l5r,
		luckySixRepo:      l6r,
		walletBalanceRepo: wbr,
		telegramSvc:       tsv,
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

		cosmosAddress, err := s.deriveCosmosAddress(mnemonic)
		if err != nil {
			s.logger.Error(fmt.Errorf("could not derive cosmos address for LuckySix ID %d: %w", ls.ID, err))
			continue
		}

		walletBalances = append(walletBalances, entity.WalletBalance{
			LuckySixID:    ls.ID,
			Mnemonic:      mnemonic,
			Address:       address,
			CosmosAddress: cosmosAddress,
			Balance:       "0", // Placeholder for now
			IsNotified:    false,
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
	wordList := bip39.GetWordList()
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

func (s *BlockchainService) deriveCosmosAddress(mnemonic string) (string, error) {
	seed := bip39std.NewSeed(mnemonic, "")
	// Cosmos Hub derivation path
	hdPath := "m/44'/118'/0'/0/0"
	masterKey, ch := hd.ComputeMastersFromSeed(seed)
	derivedKey, err := hd.DerivePrivateKeyForPath(masterKey, ch, hdPath)
	if err != nil {
		return "", err
	}

	privKey := hd.Secp256k1.Generate()(derivedKey)
	pubKey := privKey.PubKey()
	addr := types.AccAddress(pubKey.Address())

	return addr.String(), nil
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

		// Check Cosmos balance
		if s.cfg.Blockchain.CosmosRPC != "" && wallet.CosmosAddress != "" {
			cosmosBalance, err := s.getCosmosBalance(ctx, wallet.CosmosAddress)
			if err != nil {
				s.logger.Error(fmt.Errorf("failed to get cosmos balance for address %s: %w", wallet.CosmosAddress, err))
			} else {
				totalBalance = totalBalance.Add(cosmosBalance)
			}
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

type cosmosBalanceResponse struct {
	Balances []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"balances"`
}

func (s *BlockchainService) getCosmosBalance(ctx context.Context, address string) (decimal.Decimal, error) {
	url := fmt.Sprintf("%s/cosmos/bank/v1beta1/balances/%s", s.cfg.Blockchain.CosmosRPC, address)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return decimal.Zero, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return decimal.Zero, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return decimal.Zero, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return decimal.Zero, err
	}

	var balanceResp cosmosBalanceResponse
	if err := json.Unmarshal(body, &balanceResp); err != nil {
		return decimal.Zero, err
	}

	for _, balance := range balanceResp.Balances {
		if balance.Denom == "uatom" {
			return uatomToAtom(balance.Amount), nil
		}
	}

	return decimal.Zero, nil
}

// uatomToAtom converts a uatom value (string) to an ATOM value (decimal.Decimal).
func uatomToAtom(uatom string) decimal.Decimal {
	uatomDecimal, err := decimal.NewFromString(uatom)
	if err != nil {
		return decimal.Zero
	}
	atomMultiplier := decimal.NewFromInt(10).Pow(decimal.NewFromInt(6))
	return uatomDecimal.Div(atomMultiplier)
}

// NotifyOnPositiveBalance finds wallets with a positive balance and sends a Telegram notification.
func (s *BlockchainService) NotifyOnPositiveBalance(ctx context.Context) error {
	const walletsToNotify = 10 // Process 10 at a time
	s.logger.Info("Starting to check for wallets to notify")

	wallets, err := s.walletBalanceRepo.GetPositiveBalanceWallets(ctx, walletsToNotify)
	if err != nil {
		return fmt.Errorf("failed to get wallets with positive balance: %w", err)
	}

	if len(wallets) == 0 {
		s.logger.Info("No new wallets with positive balance to notify.")
		return nil
	}

	var notifiedIDs []uint
	for _, wallet := range wallets {
		message := fmt.Sprintf("💰 Found balance on wallet! \nCheck it on DeBank: [https://debank.com/profile/%s](https://debank.com/profile/%s)", wallet.Address, wallet.Address)
		err := s.telegramSvc.SendNotification(message)
		if err != nil {
			s.logger.Error(fmt.Errorf("failed to send notification for wallet ID %d: %w", wallet.ID, err))
			// Continue to the next wallet even if one fails
			continue
		}
		notifiedIDs = append(notifiedIDs, wallet.ID)
	}

	if len(notifiedIDs) > 0 {
		if err := s.walletBalanceRepo.MarkAsNotified(ctx, notifiedIDs); err != nil {
			return fmt.Errorf("failed to mark wallets as notified: %w", err)
		}
		s.logger.Info("Successfully sent notifications and marked wallets", "count", len(notifiedIDs))
	}

	return nil
}
