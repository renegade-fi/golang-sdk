package client

import (
	"github.com/google/uuid"
	"github.com/renegade-fi/golang-sdk/client/api_types"
	"github.com/renegade-fi/golang-sdk/wallet"
)

// placeOrder creates an order via the Renegade API
func (c *RenegadeClient) placeOrder(order *wallet.Order, blocking bool) error {
	// Get the back of the queue wallet
	backOfQueueWallet, err := c.GetBackOfQueueWallet()
	if err != nil {
		return err
	}

	// Add the order to the wallet and reblind
	err = backOfQueueWallet.NewOrder(*order)
	if err != nil {
		return err
	}
	err = backOfQueueWallet.Reblind()
	if err != nil {
		return err
	}

	// Sign the commitment to the new wallet
	auth, err := getWalletUpdateAuth(backOfQueueWallet)
	if err != nil {
		return err
	}

	// Post the order to the relayer
	apiOrder, err := new(api_types.ApiOrder).FromOrder(order)
	if err != nil {
		return err
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
		return err
	}

	// If blocking, wait for the task to complete
	if blocking {
		if err := c.waitForTask(resp.TaskId); err != nil {
			return err
		}
	}

	return nil
}

// cancelOrder cancels an order via the Renegade API
func (c *RenegadeClient) cancelOrder(orderId uuid.UUID, blocking bool) error {
	// Get the back of the queue wallet
	backOfQueueWallet, err := c.GetBackOfQueueWallet()
	if err != nil {
		return err
	}

	// Cancel the order
	err = backOfQueueWallet.CancelOrder(orderId)
	if err != nil {
		return err
	}
	err = backOfQueueWallet.Reblind()
	if err != nil {
		return err
	}

	// Get the wallet update auth
	auth, err := getWalletUpdateAuth(backOfQueueWallet)
	if err != nil {
		return err
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
		return err
	}

	// If blocking, wait for the task to complete
	if blocking {
		if err := c.waitForTask(resp.TaskId); err != nil {
			return err
		}
	}

	return nil
}
