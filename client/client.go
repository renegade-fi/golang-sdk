package client

import (
	"crypto/ecdsa"
	"encoding/base64"

	"github.com/google/uuid"
	"renegade.fi/golang-sdk/client/api_types"
	"renegade.fi/golang-sdk/wallet"
)

const (
	arbitrumChainId = 42161
)

// Client represents a client for the renegade API
type RenegadeClient struct {
	// chainId is the chain ID of the network
	chainId uint64
	// walletInfo is the information about the wallet necessary to recover it
	walletSecrets *wallet.WalletSecrets
	// httpClient is the HTTP client used to make requests to the renegade API
	httpClient *HttpClient
}

// NewRenegadeClient creates a new Client with the given base URL and auth key
func NewRenegadeClient(baseURL string, ethKey *ecdsa.PrivateKey) (*RenegadeClient, error) {
	return NewRenegadeClientWithChainId(baseURL, ethKey, arbitrumChainId)
}

// NewRenegadeClientWithChainId creates a new Client with the given base URL, auth key, and chain ID
func NewRenegadeClientWithChainId(baseURL string, ethKey *ecdsa.PrivateKey, chainId uint64) (*RenegadeClient, error) {
	walletInfo, err := wallet.DeriveWalletSecrets(ethKey, chainId)
	if err != nil {
		return nil, err
	}

	authKey := walletInfo.Keychain.PrivateKeys.SymmetricKey
	return &RenegadeClient{
		chainId:       chainId,
		walletSecrets: walletInfo,
		httpClient:    NewHttpClient(baseURL, &authKey),
	}, nil
}

// GetWallet retrieves a wallet from the relayer.
//
// Returns:
//   - *api_types.ApiWallet: The retrieved wallet, if successful.
//   - error: An error if the retrieval fails, nil otherwise.
//
// The method uses the client's wallet secrets to construct the API path
// and sends a GET request to the relayer. If successful, it returns the
// wallet data in the ApiWallet format.
func (c *RenegadeClient) GetWallet() (*api_types.ApiWallet, error) {
	walletId := c.walletSecrets.Id
	path := api_types.BuildGetWalletPath(walletId)

	resp := api_types.GetWalletResponse{}
	err := c.httpClient.GetWithAuth(path, nil /* body */, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Wallet, nil
}

// GetBackOfQueueWallet retrieves the wallet at the back of the processing queue from the relayer.
//
// This method sends a GET request to fetch the wallet state after all pending tasks
// in its queue have been processed. It's useful for getting the most up-to-date
// wallet state when there are known pending operations.
//
// Returns:
//   - *api_types.ApiWallet: The retrieved wallet at the back of the queue, if successful.
//   - error: An error if the retrieval fails, nil otherwise.
//
// The method uses the client's wallet ID to construct the API path and sends
// an authenticated GET request to the relayer.
func (c *RenegadeClient) GetBackOfQueueWallet() (*api_types.ApiWallet, error) {
	walletId := c.walletSecrets.Id
	path := api_types.BuildBackOfQueueWalletPath(walletId)

	resp := api_types.GetWalletResponse{}
	err := c.httpClient.GetWithAuth(path, nil /* body */, &resp)
	if err != nil {
		return nil, err
	}

	// Add the root key to the response, the relayer doesn't have it
	rootKey := c.walletSecrets.Keychain.PrivateKeys.SkRoot.ToHexString()
	w := &resp.Wallet
	w.KeyChain.PrivateKeys.SkRoot = &rootKey

	return w, nil
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
func (c *RenegadeClient) LookupWallet() (*api_types.LookupWalletResponse, error) {
	walletId := c.walletSecrets.Id
	path := api_types.LookupWalletPath

	// Build the request
	keys, err := new(api_types.ApiPrivateKeychain).FromPrivateKeychain(&c.walletSecrets.Keychain.PrivateKeys)
	if err != nil {
		return nil, err
	}
	keys.SkRoot = nil // Omit the root key

	blinderSeed := api_types.ScalarToUintLimbs(c.walletSecrets.BlinderSeed)
	shareSeed := api_types.ScalarToUintLimbs(c.walletSecrets.ShareSeed)
	request := api_types.LookupWalletRequest{
		WalletId:        walletId,
		BlinderSeed:     blinderSeed,
		ShareSeed:       shareSeed,
		PrivateKeychain: *keys,
	}

	// Post to the relayer
	resp := api_types.LookupWalletResponse{}
	err = c.httpClient.PostWithAuth(path, request, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
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
func (c *RenegadeClient) RefreshWallet() (*api_types.RefreshWalletResponse, error) {
	walletId := c.walletSecrets.Id
	path := api_types.BuildRefreshWalletPath(walletId)

	resp := api_types.RefreshWalletResponse{}
	err := c.httpClient.PostWithAuth(path, nil, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
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
func (c *RenegadeClient) CreateWallet() (*api_types.CreateWalletResponse, error) {
	// Create a new empty wallet from the base key
	newWallet, err := wallet.NewEmptyWalletFromSecrets(c.walletSecrets)
	if err != nil {
		return nil, err
	}

	apiWallet, err := new(api_types.ApiWallet).FromWallet(newWallet)
	if err != nil {
		return nil, err
	}
	// Omit the root key
	apiWallet.KeyChain.PrivateKeys.SkRoot = nil

	// Post the wallet to the relayer
	request := api_types.CreateWalletRequest{
		Wallet: *apiWallet,
	}
	resp := api_types.CreateWalletResponse{}
	err = c.httpClient.PostWithAuth(api_types.CreateWalletPath, request, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
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
func (c *RenegadeClient) PlaceOrder(order *wallet.Order) (*api_types.CreateOrderResponse, error) {
	// Get the back of the queue wallet
	apiWallet, err := c.GetBackOfQueueWallet()
	if err != nil {
		return nil, err
	}

	// Convert the API wallet to a wallet
	backOfQueueWallet, err := apiWallet.ToWallet()
	if err != nil {
		return nil, err
	}

	// Add the order to the wallet and reblind
	err = backOfQueueWallet.NewOrder(*order)
	if err != nil {
		return nil, err
	}
	backOfQueueWallet.Reblind()

	// Sign the commitment to the new wallet
	auth, err := getWalletUpdateAuth(backOfQueueWallet)
	if err != nil {
		return nil, err
	}

	// Post the order to the relayer
	apiOrder, err := new(api_types.ApiOrder).FromOrder(order)
	if err != nil {
		return nil, err
	}

	req := api_types.CreateOrderRequest{
		Order:                     *apiOrder,
		WalletUpdateAuthorization: *auth,
	}

	walletId := c.walletSecrets.Id
	path := api_types.BuildCreateOrderPath(walletId)
	resp := api_types.CreateOrderResponse{}

	err = c.httpClient.PostWithAuth(path, req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// CancelOrder cancels an order on the Renegade API.
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
func (c *RenegadeClient) CancelOrder(orderId uuid.UUID) (*api_types.CancelOrderResponse, error) {
	// Get the back of the queue wallet
	apiWallet, err := c.GetBackOfQueueWallet()
	if err != nil {
		return nil, err
	}

	// Convert the API wallet to a wallet
	backOfQueueWallet, err := apiWallet.ToWallet()
	if err != nil {
		return nil, err
	}

	// Cancel the order
	err = backOfQueueWallet.CancelOrder(orderId)
	if err != nil {
		return nil, err
	}
	backOfQueueWallet.Reblind()

	// Get the wallet update auth
	auth, err := getWalletUpdateAuth(backOfQueueWallet)
	if err != nil {
		return nil, err
	}

	// Post the order to the relayer
	walletId := c.walletSecrets.Id
	path := api_types.BuildCancelOrderPath(walletId, orderId)
	req := api_types.CancelOrderRequest{
		WalletUpdateAuthorization: *auth,
	}

	resp := api_types.CancelOrderResponse{}
	err = c.httpClient.PostWithAuth(path, req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

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
