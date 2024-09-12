package client

import (
	"github.com/google/uuid"
	"renegade.fi/golang-sdk/client/api_types"
	"renegade.fi/golang-sdk/wallet"
)

// placeOrder creates an order via the Renegade API
func (c *RenegadeClient) placeOrder(order *wallet.Order) (*api_types.CreateOrderResponse, error) {
	// Get the back of the queue wallet
	backOfQueueWallet, err := c.GetBackOfQueueWallet()
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

// cancelOrder cancels an order via the Renegade API
func (c *RenegadeClient) cancelOrder(orderId uuid.UUID) (*api_types.CancelOrderResponse, error) {
	// Get the back of the queue wallet
	backOfQueueWallet, err := c.GetBackOfQueueWallet()
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
