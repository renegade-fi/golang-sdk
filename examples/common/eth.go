package common

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	external_match_client "github.com/renegade-fi/golang-sdk/client/external_match_client"
)

const (
	// ChainID is the chain ID for the testnet
	ArbitrumSepoliaChainID = 421614
	BaseSepoliaChainID     = 84532
)

// SubmitBundle submits the bundle to the Arbitrum Sepolia network
func SubmitBundle(bundle external_match_client.ExternalMatchBundle) error {
	return SubmitBundleWithChainID(bundle, ArbitrumSepoliaChainID)
}

// SubmitBundle submits the bundle with the given chain ID
func SubmitBundleWithChainID(bundle external_match_client.ExternalMatchBundle, chainID int64) error {
	ethClient, err := GetEthClient()
	if err != nil {
		return fmt.Errorf("failed to create eth client: %w", err)
	}

	privateKey, err := GetPrivateKey()
	if err != nil {
		return fmt.Errorf("failed to get private key: %w", err)
	}

	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get gas price: %w", err)
	}

	nonce, err := ethClient.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(privateKey.PublicKey))
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}

	ethTx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   big.NewInt(chainID),
		Nonce:     nonce,
		GasTipCap: gasPrice,
		GasFeeCap: new(big.Int).Mul(gasPrice, big.NewInt(2)),
		Gas:       uint64(10_000_000),
		To:        &bundle.SettlementTx.To,
		Value:     bundle.SettlementTx.Value,
		Data:      []byte(bundle.SettlementTx.Data),
	})

	signer := types.LatestSignerForChainID(big.NewInt(chainID))
	signedTx, err := types.SignTx(ethTx, signer, privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	err = ethClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	fmt.Printf("Transaction submitted! Hash: %s\n", signedTx.Hash().Hex())
	return nil
}

// GetEthClient creates a new Ethereum client
func GetEthClient() (*ethclient.Client, error) {
	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		return nil, fmt.Errorf("RPC_URL environment variable not set")
	}
	return ethclient.Dial(rpcURL)
}

// GetPrivateKey gets the private key from environment variables
func GetPrivateKey() (*ecdsa.PrivateKey, error) {
	privKeyHex := os.Getenv("PKEY")
	if privKeyHex == "" {
		return nil, fmt.Errorf("PKEY environment variable not set")
	}

	return crypto.HexToECDSA(privKeyHex)
}
