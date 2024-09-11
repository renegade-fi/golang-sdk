package client

import (
	"github.com/google/uuid"
	"renegade.fi/golang-sdk/client/api_types"
	"renegade.fi/golang-sdk/wallet"
)

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
