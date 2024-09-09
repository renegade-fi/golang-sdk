package api_types

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	// GetWalletPath is the path for the GetWallet action
	GetWalletPath = "/v0/wallet/%s"
	// LookupWalletPath is the path for the LookupWallet action
	LookupWalletPath = "/v0/wallet/lookup"
	// CreateWalletPath is the path for the CreateWallet action
	CreateWalletPath = "/v0/wallet"
)

type ScalarLimbs [secretShareLimbCount]uint32

// buildGetWalletPath builds the path for the GetWallet action
func BuildGetWalletPath(walletId uuid.UUID) string {
	return fmt.Sprintf(GetWalletPath, walletId)
}

// GetWalletResponse is the response body for a GetWallet request
type GetWalletResponse struct {
	Wallet ApiWallet `json:"wallet"`
}

// LookupWalletRequest is the request body for the LookupWallet action
type LookupWalletRequest struct {
	WalletId        uuid.UUID          `json:"wallet_id"`
	BlinderSeed     ScalarLimbs        `json:"blinder_seed"`
	ShareSeed       ScalarLimbs        `json:"secret_share_seed"`
	PrivateKeychain ApiPrivateKeychain `json:"private_keychain"`
}

// LookupWalletResponse is the response body for a LookupWallet request
type LookupWalletResponse struct {
	WalletId uuid.UUID `json:"wallet_id"`
	TaskId   uuid.UUID `json:"task_id"`
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
