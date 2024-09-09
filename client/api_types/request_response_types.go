package api_types

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	// GetWalletPath is the path for the GetWallet action
	GetWalletPath = "/v0/wallet/%s"
	// CreateWalletPath is the path for the CreateWallet action
	CreateWalletPath = "/v0/wallet"
)

// buildGetWalletPath builds the path for the GetWallet action
func BuildGetWalletPath(walletId uuid.UUID) string {
	return fmt.Sprintf(GetWalletPath, walletId)
}

// GetWalletResponse is the response body for a GetWallet request
type GetWalletResponse struct {
	Wallet ApiWallet `json:"wallet"`
}

// CreateWalletRequest is the request body for the CreateWallet action
type CreateWalletRequest struct {
	Wallet ApiWallet `json:"wallet"`
}

// CreateWalletResponse is the response body for the CreateWallet action
type CreateWalletResponse struct {
	TaskId   uuid.UUID `json:"task_id"`
	WalletId uuid.UUID `json:"wallet_id"`
}
