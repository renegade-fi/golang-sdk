package client

import (
	"github.com/renegade-fi/golang-sdk/client/api_types"
	"github.com/renegade-fi/golang-sdk/wallet"
)

// getWallet retrieves a wallet from the relayer
func (c *RenegadeClient) getWallet() (*wallet.Wallet, error) {
	walletID := c.walletSecrets.Id
	path := api_types.BuildGetWalletPath(walletID)

	resp := api_types.GetWalletResponse{}
	err := c.httpClient.GetWithAuth(path, nil /* body */, &resp)
	if err != nil {
		return nil, err
	}

	// Convert the ApiWallet to a Wallet
	wallet, err := resp.Wallet.ToWallet()
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

// getBackOfQueueWallet retrieves the wallet at the back of the processing queue from the relayer
func (c *RenegadeClient) getBackOfQueueWallet() (*wallet.Wallet, error) {
	walletID := c.walletSecrets.Id
	path := api_types.BuildBackOfQueueWalletPath(walletID)

	resp := api_types.GetWalletResponse{}
	err := c.httpClient.GetWithAuth(path, nil /* body */, &resp)
	if err != nil {
		return nil, err
	}

	// Add the root key to the response, the relayer doesn't have it
	rootKey := c.walletSecrets.Keychain.PrivateKeys.SkRoot.ToHexString()
	w := &resp.Wallet
	w.KeyChain.PrivateKeys.SkRoot = &rootKey

	// Convert the ApiWallet to a Wallet
	wallet, err := w.ToWallet()
	if err != nil {
		return nil, err
	}

	return wallet, nil
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
func (c *RenegadeClient) lookupWallet(blocking bool) error {
	walletID := c.walletSecrets.Id
	path := api_types.LookupWalletPath

	// Build the request
	keys, err := new(api_types.ApiPrivateKeychain).
		FromPrivateKeychain(&c.walletSecrets.Keychain.PrivateKeys)
	if err != nil {
		return err
	}
	keys.SkRoot = nil // Omit the root key

	blinderSeed := api_types.ScalarToUintLimbs(c.walletSecrets.BlinderSeed)
	shareSeed := api_types.ScalarToUintLimbs(c.walletSecrets.ShareSeed)
	request := api_types.LookupWalletRequest{
		WalletId:        walletID,
		BlinderSeed:     blinderSeed,
		ShareSeed:       shareSeed,
		PrivateKeychain: *keys,
	}

	// Post to the relayer
	resp := api_types.LookupWalletResponse{}
	err = c.httpClient.PostWithAuth(path, request, &resp)
	if err != nil {
		return err
	}

	// If blocking, wait for the task to complete
	if blocking {
		// Wait for the task to complete
		if err := c.waitForTaskDirect(resp.TaskId); err != nil {
			return err
		}
	}

	return nil
}

// RefreshWallet refreshes the relayer's view of the wallet's state by looking up
// the wallet on-chain.
//
// This method sends a request to the relayer to update its local state with the
// latest on-chain
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
func (c *RenegadeClient) refreshWallet(blocking bool) error {
	walletID := c.walletSecrets.Id
	path := api_types.BuildRefreshWalletPath(walletID)

	resp := api_types.RefreshWalletResponse{}
	err := c.httpClient.PostWithAuth(path, nil, &resp)
	if err != nil {
		return err
	}

	// If blocking, wait for the task to complete
	if blocking {
		// Wait for the task to complete
		if err := c.waitForTask(resp.TaskId); err != nil {
			return err
		}
	}

	return nil
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
func (c *RenegadeClient) createWallet(blocking bool) error {
	// Create a new empty wallet from the base key
	newWallet, err := wallet.NewEmptyWalletFromSecrets(c.walletSecrets)
	if err != nil {
		return err
	}

	apiWallet, err := new(api_types.ApiWallet).FromWallet(newWallet)
	if err != nil {
		return err
	}
	// Omit the root key
	apiWallet.KeyChain.PrivateKeys.SkRoot = nil

	// Post the wallet to the relayer
	blinderSeed := api_types.ScalarToUintLimbs(c.walletSecrets.BlinderSeed)
	request := api_types.CreateWalletRequest{
		Wallet:      *apiWallet,
		BlinderSeed: blinderSeed,
	}
	resp := api_types.CreateWalletResponse{}
	err = c.httpClient.PostWithAuth(api_types.CreateWalletPath, request, &resp)
	if err != nil {
		return err
	}

	// If blocking, wait for the task to complete
	if blocking {
		// Wait for the task to complete
		if err := c.waitForTask(resp.TaskId); err != nil {
			return err
		}
	}

	return nil
}
