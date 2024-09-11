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
	// CancelOrderPath is the path for the CancelOrder action
	CancelOrderPath = "/v0/wallet/%s/orders/%s/cancel"
	// DepositPath is the path for the Deposit action
	DepositPath = "/v0/wallet/%s/balances/deposit"
)

type ScalarLimbs [secretShareLimbCount]uint32

// WalletUpdateAuthorization encapsulates the client generated authorization for wallet updates
type WalletUpdateAuthorization struct {
	// StatementSig is the signature of the commitment to the new wallet under the client's current root key
	StatementSig *string `json:"statement_sig"`
	// NewRootKey is the root key for the new wallet, if the client prefers to rotate the root key
	NewRootKey *string `json:"new_root_key"`
}

// BuildGetWalletPath builds the path for the GetWallet action
func BuildGetWalletPath(walletId uuid.UUID) string {
	return fmt.Sprintf(GetWalletPath, walletId)
}

// BuildBackOfQueueWalletPath builds the path for the BackOfQueueWallet action
func BuildBackOfQueueWalletPath(walletId uuid.UUID) string {
	return fmt.Sprintf(BackOfQueueWalletPath, walletId)
}

// BuildRefreshWalletPath builds the path for the RefreshWallet action
func BuildRefreshWalletPath(walletId uuid.UUID) string {
	return fmt.Sprintf(RefreshWalletPath, walletId)
}

// BuildCreateOrderPath builds the path for the CreateOrder action
func BuildCreateOrderPath(walletId uuid.UUID) string {
	return fmt.Sprintf(CreateOrderPath, walletId)
}

// BuildCancelOrderPath builds the path for the CancelOrder action
func BuildCancelOrderPath(walletId uuid.UUID, orderId uuid.UUID) string {
	return fmt.Sprintf(CancelOrderPath, walletId, orderId)
}

// BuildDepositPath builds the path for the Deposit action
func BuildDepositPath(walletId uuid.UUID) string {
	return fmt.Sprintf(DepositPath, walletId)
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

// CancelOrderRequest is the request body for the CancelOrder action
type CancelOrderRequest struct {
	WalletUpdateAuthorization
}

// CancelOrderResponse is the response body for the CancelOrder action
type CancelOrderResponse struct {
	// TaskId is the ID of the task that was created to update the wallet
	TaskId uuid.UUID `json:"task_id"`
	// Order is the order that was canceled
	Order ApiOrder `json:"order"`
}

// DepositRequest is the request body for the Deposit action
type DepositRequest struct {
	// FromAddr is the address to deposit from
	FromAddr string `json:"from_addr"`
	// Mint is the mint of the token to deposit
	Mint string `json:"mint"`
	// Amount is the amount of the token to deposit
	Amount string `json:"amount"`
	// WalletUpdateAuthorization is the authorization for the wallet update
	WalletUpdateAuthorization
	// PermitNonce is the nonce used in the associated Permit2 permit
	PermitNonce string `json:"permit_nonce"`
	// PermitDeadline is the deadline used in the associated Permit2 permit
	PermitDeadline string `json:"permit_deadline"`
	// PermitSignature is the signature over the associated Permit2 permit,
	// allowing the contract to guarantee that the deposit is sourced from
	// the correct account
	PermitSignature string `json:"permit_signature"`
}

// DepositResponse is the response body for the Deposit action
type DepositResponse struct {
	// TaskId is the ID of the task that was created to update the wallet
	TaskId uuid.UUID `json:"task_id"`
}
