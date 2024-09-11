package client

import (
	"renegade.fi/golang-sdk/client/api_types"
	"renegade.fi/golang-sdk/wallet"
)

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
