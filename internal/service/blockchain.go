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
