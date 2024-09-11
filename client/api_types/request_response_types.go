package api_types

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	// GetWalletPath is the path for the GetWallet action
	GetWalletPath = "/v0/wallet/%s"
	// BackOfQueueWalletPath is the path to fetch the wallet after all tasks in its queue have been processed
	BackOfQueueWalletPath = "/v0/wallet/%s/back-of-queue"
	// LookupWalletPath is the path for the LookupWallet action
	LookupWalletPath = "/v0/wallet/lookup"
	// RefreshWalletPath is the path for the RefreshWallet action
	RefreshWalletPath = "/v0/wallet/%s/refresh"
	// CreateWalletPath is the path for the CreateWallet action
	CreateWalletPath = "/v0/wallet"
	// CreateOrderPath is the path for the CreateOrder action
	CreateOrderPath = "/v0/wallet/%s/orders"
)

type ScalarLimbs [secretShareLimbCount]uint32

// WalletUpdateAuthorization encapsulates the client generated authorization for wallet updates
type WalletUpdateAuthorization struct {
	// StatementSig is the signature of the commitment to the new wallet under the client's current root key
	StatementSig *string `json:"statement_sig"`
	// NewRootKey is the root key for the new wallet, if the client prefers to rotate the root key
	NewRootKey *string `json:"new_root_key"`
}

// buildGetWalletPath builds the path for the GetWallet action
func BuildGetWalletPath(walletId uuid.UUID) string {
	return fmt.Sprintf(GetWalletPath, walletId)
}

// buildBackOfQueueWalletPath builds the path for the BackOfQueueWallet action
func BuildBackOfQueueWalletPath(walletId uuid.UUID) string {
	return fmt.Sprintf(BackOfQueueWalletPath, walletId)
}

// buildRefreshWalletPath builds the path for the RefreshWallet action
func BuildRefreshWalletPath(walletId uuid.UUID) string {
	return fmt.Sprintf(RefreshWalletPath, walletId)
}

// buildCreateOrderPath builds the path for the CreateOrder action
func BuildCreateOrderPath(walletId uuid.UUID) string {
	return fmt.Sprintf(CreateOrderPath, walletId)
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

// RefreshWalletResponse is the response body for a RefreshWallet request
type RefreshWalletResponse struct {
	TaskId uuid.UUID `json:"task_id"`
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

// CreateOrderRequest is the request body for the CreateOrder action
type CreateOrderRequest struct {
	Order ApiOrder `json:"order"`
	WalletUpdateAuthorization
}

// CreateOrderResponse is the response body for the CreateOrder action
type CreateOrderResponse struct {
	// Id is the ID of the order that was created
	Id uuid.UUID `json:"id"`
	// TaskId is the ID of the task that was created to update the wallet
	TaskId uuid.UUID `json:"task_id"`
}
