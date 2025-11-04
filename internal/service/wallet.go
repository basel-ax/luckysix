package service

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/alexvec/go-bip39"
	"github.com/basel-ax/luckysix/internal/entity"
	"github.com/basel-ax/luckysix/internal/repository"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

var wordList []string

func init() {
	wordList = bip39.GetWordList()
}

// WalletService -.
type WalletService struct {
	luckyTwoRepo  repository.LuckyTwo
	luckyFiveRepo repository.LuckyFive
}

// NewWalletService -.
func NewWalletService(luckyTwoRepo repository.LuckyTwo, luckyFiveRepo repository.LuckyFive) *WalletService {
	return &WalletService{
		luckyTwoRepo:  luckyTwoRepo,
		luckyFiveRepo: luckyFiveRepo,
	}
}

// GenerateMnemonic -.
func (s *WalletService) GenerateMnemonic(ctx context.Context, luckySix *entity.LuckySix) (string, error) {
	luckyFive, err := s.luckyFiveRepo.GetByID(ctx, luckySix.LuckyFiveID)
	if err != nil {
		return "", fmt.Errorf("error getting lucky five by id: %w", err)
	}

	luckyTwo, err := s.luckyTwoRepo.GetByID(ctx, luckySix.LuckyTwoID)
	if err != nil {
		return "", fmt.Errorf("error getting lucky two by id: %w", err)
	}

	pairOne, err := s.luckyTwoRepo.GetByID(ctx, luckyFive.PairOne)
	if err != nil {
		return "", err
	}
	pairTwo, err := s.luckyTwoRepo.GetByID(ctx, luckyFive.PairTwo)
	if err != nil {
		return "", err
	}
	pairThree, err := s.luckyTwoRepo.GetByID(ctx, luckyFive.PairThree)
	if err != nil {
		return "", err
	}
	pairFour, err := s.luckyTwoRepo.GetByID(ctx, luckyFive.PairFour)
	if err != nil {
		return "", err
	}
	pairFive, err := s.luckyTwoRepo.GetByID(ctx, luckyFive.PairFive)
	if err != nil {
		return "", err
	}

	words := []string{
		wordList[pairOne.WordOne-1],
		wordList[pairOne.WordTwo-1],
		wordList[pairTwo.WordOne-1],
		wordList[pairTwo.WordTwo-1],
		wordList[pairThree.WordOne-1],
		wordList[pairThree.WordTwo-1],
		wordList[pairFour.WordOne-1],
		wordList[pairFour.WordTwo-1],
		wordList[pairFive.WordOne-1],
		wordList[pairFive.WordTwo-1],
		wordList[luckyTwo.WordOne-1],
		wordList[luckyTwo.WordTwo-1],
	}

	return strings.Join(words, " "), nil
}

// GenerateWallet -.
func (s *WalletService) GenerateWallet(ctx context.Context, mnemonic string) (string, string, error) {

	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return "", "", err
	}

	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		return "", "", err
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		return "", "", err
	}

	address := account.Address.Hex()

	// Cosmos address
	cosmosPath := sdk.GetConfig().GetFullBIP44Path()
	params := hd.NewFundraiserParams(0, sdk.CoinType, 0)

	hdPath := params.String() + "/" + cosmosPath
	master, ch := hd.ComputeMastersFromSeed(seed)
	derivedP, err := hd.DerivePrivateKeyForPath(master, ch, hdPath)
	if err != nil {
		return "", "", err
	}
	cosmosAddress := sdk.AccAddress(derivedP[:20]).String()

	return address, cosmosAddress, nil
}

// CheckBalance -.
func (s *WalletService) CheckBalance(ctx context.Context, address string, blockchains []*entity.Blockchain) (string, error) {
	totalBalance := big.NewInt(0)

	for _, b := range blockchains {
		client, err := ethclient.Dial(b.NodeHTTPProvider)
		if err != nil {
			return "", fmt.Errorf("error connecting to blockchain %s: %w", b.Title, err)
		}

		balance, err := client.BalanceAt(ctx, common.HexToAddress(address), nil)
		if err != nil {
			return "", fmt.Errorf("error checking balance on blockchain %s: %w", b.Title, err)
		}
		totalBalance.Add(totalBalance, balance)
	}

	return totalBalance.String(), nil
}
