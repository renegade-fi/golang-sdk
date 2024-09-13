package client

import (
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/renegade-fi/golang-sdk/client"
	"github.com/renegade-fi/golang-sdk/client/api_types"
	"github.com/renegade-fi/golang-sdk/wallet"
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

// RenegadeClient represents a client for the renegade API
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

// GetWallet retrieves the current wallet state from the relayer.
//
// Returns:
//   - *wallet.Wallet: The retrieved wallet, if successful.
//   - error: An error if the retrieval fails, nil otherwise.
//
// This method sends a GET request to the relayer to fetch the current
// wallet state. It uses the client's wallet ID to construct the API path.
// The retrieved wallet data is converted from the API format to the internal
// wallet.Wallet type before being returned.
func (c *RenegadeClient) GetWallet() (*wallet.Wallet, error) {
	return c.getWallet()
}

// GetBackOfQueueWallet retrieves the wallet at the back of the processing queue from the relayer.
//
// This method sends a GET request to fetch the wallet state after all pending tasks
// in its queue have been processed. It's useful for getting the most up-to-date
// wallet state when there are known pending operations.
//
// Returns:
//   - *wallet.Wallet: The retrieved wallet at the back of the queue, if successful.
//   - error: An error if the retrieval fails, nil otherwise.
//
// The method uses the client's wallet ID to construct the API path and sends
// an authenticated GET request to the relayer.
func (c *RenegadeClient) GetBackOfQueueWallet() (*wallet.Wallet, error) {
	return c.getBackOfQueueWallet()
}

// CheckWallet verifies the wallet's existence in the relayer's state and retrieves it from the blockchain if necessary.
//
// This method first attempts to fetch the wallet from the relayer's local state using GetWallet().
// If successful, it returns the wallet immediately. If the wallet is not found in the local state,
// it initiates a blockchain lookup using LookupWallet() to retrieve the wallet information.
//
// Returns:
//   - *wallet.Wallet: The retrieved wallet, if found either in local state or on-chain.
//   - error: An error if both local retrieval and on-chain lookup fail, nil otherwise.
//
// This method is useful for ensuring that the client has the most up-to-date wallet information,
// especially in scenarios where the wallet might not be synchronized between the relayer and the blockchain.
func (c *RenegadeClient) CheckWallet() (*wallet.Wallet, error) {
	wallet, err := c.GetWallet()
	if err == nil {
		return wallet, nil
	}
	return c.LookupWallet()
}

// LookupWallet looks up a wallet in the relayer from contract state.
//
// This method sends a request to the relayer to retrieve wallet information
// from the blockchain. It uses the client's wallet secrets to construct the request.
//
// Returns:
//   - *api_types.LookupWalletResponse: Contains the wallet ID and task ID if successful.
//   - error: An error if the lookup fails, nil otherwise.
//
// The method constructs a LookupWalletRequest with the wallet ID, blinder seed,
// share seed, and private keychain (excluding the root key). It then sends a POST
// request to the relayer and returns the response.
func (c *RenegadeClient) LookupWallet() (*wallet.Wallet, error) {
	if err := c.lookupWallet(true /* blocking */); err != nil {
		return nil, err
	}
	return c.getWallet()
}

// RefreshWallet refreshes the relayer's view of the wallet's state by looking up the wallet on-chain.
//
// This method sends a request to the relayer to update its local state with the latest on-chain
// information for the wallet associated with the client. It's useful for synchronizing the
// relayer's view with the current blockchain state, especially after on-chain transactions.
//
// Returns:
//   - *api_types.RefreshWalletResponse: Contains the task ID for the refresh operation.
//   - error: An error if the refresh operation fails, nil otherwise.
//
// The method uses the client's wallet ID to construct the API path and sends a POST request
// to the relayer. If successful, it returns the response containing the task ID for tracking
// the refresh operation.
func (c *RenegadeClient) RefreshWallet() (*wallet.Wallet, error) {
	if err := c.refreshWallet(true /* blocking */); err != nil {
		return nil, err
	}
	return c.getWallet()
}

// CreateWallet creates a new wallet derived from the client's wallet secrets.
//
// Returns:
//   - *api_types.CreateWalletResponse: Contains the task ID and wallet ID of the created wallet
//   - error: An error if the wallet creation fails, nil otherwise
//
// The method generates a new Renegade wallet using the client's wallet secrets,
// submits a creation request to the Renegade API, and returns the response.
// This wallet can be used for private transactions within the Renegade network.
func (c *RenegadeClient) CreateWallet() (*wallet.Wallet, error) {
	if err := c.createWallet(true /* blocking */); err != nil {
		return nil, err
	}
	return c.getWallet()
}

// Deposit deposits funds into the wallet associated with the client.
//
// This method initiates a deposit transaction, adding the specified amount of
// a given token (identified by its mint address) to the client's wallet. It
// interacts with the Ethereum blockchain and the Renegade protocol to process
// the deposit.
//
// Parameters:
//   - mint: A pointer to a string representing the token's mint address.
//   - amount: A pointer to a big.Int representing the amount to deposit.
//   - ethPrivateKey: The Ethereum private key used to sign the transaction.
//
// Returns:
//   - *api_types.DepositResponse: Contains information about the deposit transaction,
//     including the task ID and any relevant details from the Renegade protocol.
//   - error: An error if the deposit process fails, nil otherwise.
//
// The method handles the entire deposit flow, including updating the local wallet
// state, approving the Permit2 contract for spending, and submitting the deposit
// request to the Renegade relayer.
func (c *RenegadeClient) Deposit(mint string, amount *big.Int, ethPrivateKey *ecdsa.PrivateKey) (*wallet.Wallet, error) {
	if err := c.deposit(mint, amount, ethPrivateKey, true /* blocking */); err != nil {
		return nil, err
	}
	return c.GetWallet()
}

// Withdraw initiates a withdrawal transaction, removing the specified amount
// of a given token (identified by its mint address) from the client's wallet. It
// interacts with the Ethereum blockchain and the Renegade protocol to process
// the withdrawal.
//
// Parameters:
//   - mint: A pointer to a string representing the token's mint address.
//   - amount: A pointer to a big.Int representing the amount to withdraw.
//   - ethPrivateKey: The Ethereum private key used to sign the transaction.
//
// Returns:
//   - *api_types.WithdrawResponse: Contains information about the withdrawal transaction,
//     including the task ID and any relevant details from the Renegade protocol.
//   - error: An error if the withdrawal process fails, nil otherwise.
func (c *RenegadeClient) Withdraw(mint string, amount *big.Int) (*wallet.Wallet, error) {
	if err := c.withdraw(mint, amount, true /* blocking */); err != nil {
		return nil, err
	}
	return c.GetWallet()
}

// WithdrawToAddress withdraws funds from the wallet to the given address
func (c *RenegadeClient) WithdrawToAddress(mint string, amount *big.Int, destination string) (*wallet.Wallet, error) {
	if err := c.withdrawToAddress(mint, amount, destination, true /* blocking */); err != nil {
		return nil, err
	}
	return c.GetWallet()
}

// PayFees initiates the fee payment process for the wallet.
//
// This method sends a request to the Renegade API to pay any outstanding fees
// associated with the client's wallet. It handles the entire fee payment flow,
// including updating the local wallet state and submitting the fee payment
// request to the Renegade relayer.
//
// Returns:
//   - *wallet.Wallet: An updated wallet object reflecting the new state after fee payment.
//   - error: An error if the fee payment process fails, nil otherwise.
//
// The method waits for the fee payment to be processed before returning the updated wallet.
func (c *RenegadeClient) PayFees() (*wallet.Wallet, error) {
	if err := c.payFees(); err != nil {
		return nil, err
	}

	return c.getBackOfQueueWallet()
}

// PlaceOrder creates an order on the Renegade API.
//
// This method sends a request to the Renegade API to create an order for a specified
// token pair. It uses the client's wallet ID and the provided token details to construct
// the request.
//
// Returns:
//   - *api_types.CreateOrderResponse: Contains the order ID and task ID if successful.
//   - error: An error if the order creation fails, nil otherwise.
func (c *RenegadeClient) PlaceOrder(order *wallet.Order) (*wallet.Wallet, error) {
	if err := c.placeOrder(order, true /* blocking */); err != nil {
		return nil, err
	}
	return c.GetWallet()
}

// CancelOrder cancels an order via the Renegade API.
//
// This method sends a request to the Renegade API to cancel an order for the
// client's wallet. It uses the client's wallet ID and the provided order ID to
// construct the request. The method first retrieves the latest wallet state,
// cancels the order locally, and then sends the update to the API.
//
// Parameters:
//   - orderId: The UUID of the order to cancel.
//
// Returns:
//   - *api_types.CancelOrderResponse: Contains the task ID and the canceled order if successful.
//   - error: An error if the order cancellation fails, nil otherwise.
func (c *RenegadeClient) CancelOrder(orderId uuid.UUID) (*wallet.Wallet, error) {
	if err := c.cancelOrder(orderId, true /* blocking */); err != nil {
		return nil, err
	}
	return c.GetWallet()
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
