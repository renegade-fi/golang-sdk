package client

import (
	"crypto/ecdsa"

	"renegade.fi/golang-sdk/client/api_types"
	"renegade.fi/golang-sdk/wallet"
)

const (
	arbitrumChainId = 42161
)

// Client represents a client for the renegade API
type RenegadeClient struct {
	httpClient *HttpClient
}

// NewRenegadeClient creates a new Client with the given base URL and auth key
func NewRenegadeClient(baseURL string, authKey *wallet.HmacKey) *RenegadeClient {
	return &RenegadeClient{
		httpClient: NewHttpClient(baseURL, authKey),
	}
}

// GetWallet retrieves a wallet from the relayer.
//
// Parameters:
//   - ethKey: *ecdsa.PrivateKey - The Ethereum private key associated with the wallet.
//   - chainId: uint64 - The chain ID of the network.
//
// Returns:
//   - *api_types.ApiWallet: The retrieved wallet, if successful.
//   - error: An error if the retrieval fails, nil otherwise.
//
// The method first derives the wallet ID using the provided Ethereum key and chain ID.
// It then constructs the API path and sends a GET request to the relayer.
// If successful, it returns the wallet data in the ApiWallet format.
func (c *RenegadeClient) GetWallet(ethKey *ecdsa.PrivateKey, chainId uint64) (*api_types.ApiWallet, error) {
	walletId, err := wallet.DeriveWalletID(ethKey, chainId)
	if err != nil {
		return nil, err
	}

	path := api_types.BuildGetWalletPath(walletId)
	resp := api_types.GetWalletResponse{}
	err = c.httpClient.GetWithAuth(path, nil /* body */, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Wallet, nil
}

// CreateWallet creates a new wallet derived with provided Ethereum private key.
//
// Parameters:
//   - ethKey: The Ethereum private key used to create and control the wallet
//
// Returns:
//   - *api_types.CreateWalletResponse: Contains the task ID and wallet ID of the created wallet
//   - error: An error if the wallet creation fails, nil otherwise
//
// The method generates a new Renegade wallet associated with the given Ethereum key,
// submits a creation request to the Renegade API, and returns the response.
// This wallet can be used for private transactions within the Renegade network.
func (c *RenegadeClient) CreateWallet(ethKey *ecdsa.PrivateKey) (*api_types.CreateWalletResponse, error) {
	return c.CreateWalletWithChainId(ethKey, arbitrumChainId)
}

func (c *RenegadeClient) CreateWalletWithChainId(ethKey *ecdsa.PrivateKey, chainId uint64) (*api_types.CreateWalletResponse, error) {
	// Create a new empty wallet from the base key
	newWallet, err := wallet.NewEmptyWallet(ethKey, chainId)
	if err != nil {
		return nil, err
	}

	apiWallet, err := new(api_types.ApiWallet).FromWallet(newWallet)
	if err != nil {
		return nil, err
	}

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
