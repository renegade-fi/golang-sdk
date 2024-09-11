package client

import (
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"renegade.fi/golang-sdk/client"
	"renegade.fi/golang-sdk/client/api_types"
	"renegade.fi/golang-sdk/wallet"
)

// ChainConfig represents the configuration for a specific chain
type ChainConfig struct {
	// ChainID is the chain ID of the chain
	ChainID uint64
	// Permit2Address is the address of the Permit2 contract
	Permit2Address string
	// DarkpoolAddress is the address of the Darkpool contract
	DarkpoolAddress string
	// EthereumRpcUrl is the URL of the Ethereum RPC
	EthereumRpcUrl string
}

var (
	ArbitrumOneConfig = ChainConfig{
		ChainID:         42161,
		Permit2Address:  "0x000000000022D473030F116dDEE9F6B43aC78BA3",
		DarkpoolAddress: "0x30bd8eab29181f790d7e495786d4b96d7afdc518",
		EthereumRpcUrl:  "https://arb1.arbitrum.io/rpc",
	}

	ArbitrumSepoliaConfig = ChainConfig{
		ChainID:         421614,
		Permit2Address:  "0x9458198bcc289c42e460cb8ca143e5854f734442",
		DarkpoolAddress: "0x9af58f1ff20ab22e819e40b57ffd784d115a9ef5",
		EthereumRpcUrl:  "https://sepolia-rollup.arbitrum.io/rpc",
	}
)

// Client represents a client for the renegade API
type RenegadeClient struct {
	chainConfig   ChainConfig
	walletSecrets *wallet.WalletSecrets
	httpClient    *client.HttpClient
}

// NewRenegadeClient creates a new Client with the given base URL and auth key
func NewRenegadeClient(baseURL string, ethKey *ecdsa.PrivateKey) (*RenegadeClient, error) {
	return NewRenegadeClientWithConfig(baseURL, ethKey, ArbitrumOneConfig)
}

// NewSepoliaRenegadeClient creates a new Client with the given base URL and auth key
func NewSepoliaRenegadeClient(baseURL string, ethKey *ecdsa.PrivateKey) (*RenegadeClient, error) {
	return NewRenegadeClientWithConfig(baseURL, ethKey, ArbitrumSepoliaConfig)
}

// NewRenegadeClientWithConfig creates a new Client with the given base URL, auth key, and chain config
func NewRenegadeClientWithConfig(baseURL string, ethKey *ecdsa.PrivateKey, config ChainConfig) (*RenegadeClient, error) {
	walletInfo, err := wallet.DeriveWalletSecrets(ethKey, config.ChainID)
	if err != nil {
		return nil, err
	}

	authKey := walletInfo.Keychain.PrivateKeys.SymmetricKey
	return &RenegadeClient{
		chainConfig:   config,
		walletSecrets: walletInfo,
		httpClient:    client.NewHttpClient(baseURL, &authKey),
	}, nil
}

// --- Helpers --- //

// getWalletUpdateAuth gets the wallet update authorization for the given wallet
func getWalletUpdateAuth(wallet *wallet.Wallet) (*api_types.WalletUpdateAuthorization, error) {
	// Compute the commitment to the new wallet
	commitment, err := wallet.GetShareCommitment()
	if err != nil {
		return nil, err
	}

	// Sign the commitment with skRoot
	signature, err := wallet.SignCommitment(commitment)
	if err != nil {
		return nil, err
	}

	// base64 encode the signature without padding
	signatureStr := base64.RawStdEncoding.EncodeToString(signature)
	return &api_types.WalletUpdateAuthorization{
		StatementSig: &signatureStr,
	}, nil
}

// createRpcClient creates a new RPC client
func (c *RenegadeClient) createRpcClient() (*ethclient.Client, error) {
	return ethclient.Dial(c.chainConfig.EthereumRpcUrl)
}

// createTransactor creates a new transactor with the given private key and chain ID
func (c *RenegadeClient) createTransactor(ethPrivateKey *ecdsa.PrivateKey) (*bind.TransactOpts, error) {
	chainID := big.NewInt(int64(c.chainConfig.ChainID))
	auth, err := bind.NewKeyedTransactorWithChainID(ethPrivateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}
	return auth, nil
}
