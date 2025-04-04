package api_types //nolint:revive

import (
	"fmt"

	"github.com/google/uuid"
)

//nolint:revive
const (
	// --- Orderbook Endpoints --- //
	// GetSupportedTokensPath is the path for the GetSupportedTokens action
	GetSupportedTokensPath = "/v0/supported-tokens"
	// GetFeeForAssetPath is the path for the GetFeeForAsset action
	GetFeeForAssetPath = "/v0/order_book/external-match-fee"

	// --- Wallet Endpoints --- //
	// GetWalletPath is the path for the GetWallet action
	GetWalletPath = "/v0/wallet/%s"
	// BackOfQueueWalletPath is the path to fetch the wallet after all tasks
	// in its queue have been processed
	BackOfQueueWalletPath = "/v0/wallet/%s/back-of-queue"
	// LookupWalletPath is the path for the LookupWallet action
	LookupWalletPath = "/v0/wallet/lookup" //nolint:gosec
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
	// WithdrawPath is the path for the Withdraw action
	WithdrawPath = "/v0/wallet/%s/balances/%s/withdraw"
	// PayFeesPath is the path to enqueue tasks to pay wallet fees
	PayFeesPath = "/v0/wallet/%s/pay-fees"
	// TaskStatusPath is the path to fetch the status of a task
	TaskStatusPath = "/v0/tasks/%s"
	// TaskHistoryPath is the path to fetch the task history for a wallet
	TaskHistoryPath = "/v0/wallet/%s/task-history"

	// --- External Match Endpoints --- //
	// GetExternalMatchBundlePath is the path to fetch an external match bundle
	GetExternalMatchBundlePath = "/v0/matching-engine/request-external-match"
	// GetExternalMatchQuotePath is the path to fetch an external match quote
	GetExternalMatchQuotePath = "/v0/matching-engine/quote"
	// AssembleExternalQuotePath is the path to assemble a quote into a settlement transaction
	AssembleExternalQuotePath = "/v0/matching-engine/assemble-external-match"

	// --- External Match Query Params --- //
	// DisableGasSponsorshipParam is the query param used to disable gas sponsorship
	DisableGasSponsorshipParam = "disable_gas_sponsorship"
	// GasRefundAddressParam is the query param used to specify the gas refund address
	GasRefundAddressParam = "refund_address"
	// RefundNativeEthParam is the query param used to specify whether to refund the gas in native ETH
	RefundNativeEthParam = "refund_native_eth"
)

// ScalarLimbs is an array of uint32 limbs
type ScalarLimbs [secretShareLimbCount]uint32

// WalletUpdateAuthorization encapsulates the client generated authorization for wallet updates
type WalletUpdateAuthorization struct {
	// StatementSig is the signature of the commitment to the new wallet under
	// the client's current root key
	StatementSig *string `json:"statement_sig"`
	// NewRootKey is the root key for the new wallet, if the client prefers to rotate the root key
	NewRootKey *string `json:"new_root_key"`
}

// BuildGetFeeForAssetPath builds the path for the GetFeeForAsset action
func BuildGetFeeForAssetPath(mint string) string {
	return fmt.Sprintf("%s?mint=%s", GetFeeForAssetPath, mint)
}

// BuildGetWalletPath builds the path for the GetWallet action
func BuildGetWalletPath(walletID uuid.UUID) string {
	return fmt.Sprintf(GetWalletPath, walletID)
}

// BuildBackOfQueueWalletPath builds the path for the BackOfQueueWallet action
func BuildBackOfQueueWalletPath(walletID uuid.UUID) string {
	return fmt.Sprintf(BackOfQueueWalletPath, walletID)
}

// BuildRefreshWalletPath builds the path for the RefreshWallet action
func BuildRefreshWalletPath(walletID uuid.UUID) string {
	return fmt.Sprintf(RefreshWalletPath, walletID)
}

// BuildCreateOrderPath builds the path for the CreateOrder action
func BuildCreateOrderPath(walletID uuid.UUID) string {
	return fmt.Sprintf(CreateOrderPath, walletID)
}

// BuildCancelOrderPath builds the path for the CancelOrder action
func BuildCancelOrderPath(walletID uuid.UUID, orderID uuid.UUID) string {
	return fmt.Sprintf(CancelOrderPath, walletID, orderID)
}

// BuildDepositPath builds the path for the Deposit action
func BuildDepositPath(walletID uuid.UUID) string {
	return fmt.Sprintf(DepositPath, walletID)
}

// BuildWithdrawPath builds the path for the Withdraw action
func BuildWithdrawPath(walletID uuid.UUID, mint string) string {
	return fmt.Sprintf(WithdrawPath, walletID, mint)
}

// BuildPayFeesPath builds the path for the PayFees action
func BuildPayFeesPath(walletID uuid.UUID) string {
	return fmt.Sprintf(PayFeesPath, walletID)
}

// BuildTaskStatusPath builds the path for the TaskStatus action
func BuildTaskStatusPath(taskID uuid.UUID) string {
	return fmt.Sprintf(TaskStatusPath, taskID)
}

// BuildTaskHistoryPath builds the path for the TaskHistory action
func BuildTaskHistoryPath(walletID uuid.UUID) string {
	return fmt.Sprintf(TaskHistoryPath, walletID)
}

// -----------------------
// | Orderbook Endpoints |
// -----------------------

// GetSupportedTokensResponse is the response body for the GetSupportedTokens request
type GetSupportedTokensResponse struct {
	Tokens []ApiToken `json:"tokens"`
}

// --------------------
// | Wallet Endpoints |
// --------------------

// GetWalletResponse is the response body for a GetWallet request
type GetWalletResponse struct {
	Wallet ApiWallet `json:"wallet"`
}

// LookupWalletRequest is the request body for the LookupWallet action
type LookupWalletRequest struct {
	WalletId        uuid.UUID          `json:"wallet_id"` //nolint:revive
	BlinderSeed     ScalarLimbs        `json:"blinder_seed"`
	ShareSeed       ScalarLimbs        `json:"secret_share_seed"`
	PrivateKeychain ApiPrivateKeychain `json:"private_keychain"`
}

// LookupWalletResponse is the response body for a LookupWallet request
type LookupWalletResponse struct {
	WalletId uuid.UUID `json:"wallet_id"` //nolint:revive
	TaskId   uuid.UUID `json:"task_id"`   //nolint:revive
}

// RefreshWalletResponse is the response body for a RefreshWallet request
type RefreshWalletResponse struct {
	TaskId uuid.UUID `json:"task_id"` //nolint:revive
}

// CreateWalletRequest is the request body for the CreateWallet action
type CreateWalletRequest struct {
	Wallet      ApiWallet   `json:"wallet"`
	BlinderSeed ScalarLimbs `json:"blinder_seed"`
}

// CreateWalletResponse is the response body for the CreateWallet action
type CreateWalletResponse struct {
	TaskId   uuid.UUID `json:"task_id"`   //nolint:revive
	WalletId uuid.UUID `json:"wallet_id"` //nolint:revive
}

// CreateOrderRequest is the request body for the CreateOrder action
type CreateOrderRequest struct {
	Order ApiOrder `json:"order"`
	WalletUpdateAuthorization
}

// CreateOrderResponse is the response body for the CreateOrder action
type CreateOrderResponse struct {
	// Id is the ID of the order that was created
	Id uuid.UUID `json:"id"` //nolint:revive
	// TaskId is the ID of the task that was created to update the wallet
	TaskId uuid.UUID `json:"task_id"` //nolint:revive
}

// CancelOrderRequest is the request body for the CancelOrder action
type CancelOrderRequest struct {
	WalletUpdateAuthorization
}

// CancelOrderResponse is the response body for the CancelOrder action
type CancelOrderResponse struct {
	// TaskId is the ID of the task that was created to update the wallet
	TaskId uuid.UUID `json:"task_id"` //nolint:revive
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
	TaskId uuid.UUID `json:"task_id"` //nolint:revive
}

// WithdrawRequest is the request body for the Withdraw action
type WithdrawRequest struct {
	// DestinationAddr is the address to withdraw to
	DestinationAddr string `json:"destination_addr"`
	// Amount is the amount of the token to withdraw
	Amount string `json:"amount"`
	// ExternalTransferSig is a signature of the external transfer to authorize
	// the withdrawal and location
	ExternalTransferSig *string `json:"external_transfer_sig"`
	// WalletUpdateAuthorization is the authorization for the wallet update
	WalletUpdateAuthorization
}

// WithdrawResponse is the response body for the Withdraw action
type WithdrawResponse struct {
	// TaskId is the ID of the task that was created to update the wallet
	TaskId uuid.UUID `json:"task_id"` //nolint:revive
}

// PayFeesResponse is the response body for the PayFees action
type PayFeesResponse struct {
	// TaskIds are the IDs of the tasks that were created to pay the fees
	TaskIds []uuid.UUID `json:"task_ids"` //nolint:revive
}

// ApiTaskStatus is the status of a running task
// ApiTaskStatus represents the status of a task
type ApiTaskStatus struct { //nolint:revive
	// ID is the identifier of the task
	ID uuid.UUID `json:"id"`
	// Description is the description of the task
	Description string `json:"description"`
	// State is the current state of the task
	State string `json:"state"`
	// Committed indicates whether the task has already committed
	Committed bool `json:"committed"`
}

// TaskResponse is the response body for the Task endpoint
type TaskResponse struct {
	// Status is the current status of the task
	Status ApiTaskStatus `json:"status"`
}

// ApiHistoricalTask represents a historical task
type ApiHistoricalTask struct { //nolint:revive
	// ID is the identifier of the task
	Id uuid.UUID `json:"id"` //nolint:revive
	// State is the current state of the task
	State string `json:"state"`
	// CreatedAt is the timestamp when the task was created
	CreatedAt uint64 `json:"created_at"`
}

// TaskHistoryResponse is the response body for the TaskHistory endpoint
type TaskHistoryResponse struct {
	// Tasks is the list of tasks in the queue
	Tasks []ApiHistoricalTask `json:"tasks"`
}

// ----------------------------
// | External Match Endpoints |
// ----------------------------

// ExternalMatchRequest is a request to generate an external match
type ExternalMatchRequest struct {
	ExternalOrder   ApiExternalOrder `json:"external_order"`
	DoGasEstimation bool             `json:"do_gas_estimation"`
	// ReceiverAddress is the address to receive the settlement,
	// i.e. the address to which the darkpool will send tokens
	ReceiverAddress *string `json:"receiver_address,omitempty"`
}

// ExternalMatchResponse is the response body for the ExternalMatch action
type ExternalMatchResponse struct {
	Bundle       ApiExternalMatchBundle `json:"match_bundle"`
	GasSponsored bool                   `json:"is_sponsored"`
	// The gas sponsorship info, if the match was sponsored
	GasSponsorshipInfo *ApiGasSponsorshipInfo `json:"gas_sponsorship_info,omitempty"`
}

// ExternalQuoteRequest is a request to fetch an external match quote
type ExternalQuoteRequest struct {
	ExternalOrder ApiExternalOrder `json:"external_order"`
}

// ExternalQuoteResponse is the response body for the ExternalQuote action
type ExternalQuoteResponse struct {
	Quote SignedQuoteResponse `json:"signed_quote"`
	// The signed gas sponsorship info, if sponsorship was requested
	GasSponsorshipInfo *ApiSignedGasSponsorshipInfo `json:"gas_sponsorship_info,omitempty"`
}

// AssembleExternalQuoteRequest is a request to assemble an external match quote
// into a settlement transaction
type AssembleExternalQuoteRequest struct {
	Quote           SignedQuoteResponse `json:"signed_quote"`
	DoGasEstimation bool                `json:"do_gas_estimation"`
	// ReceiverAddress is the address to receive the settlement,
	// i.e. the address to which the darkpool will send tokens
	ReceiverAddress *string `json:"receiver_address,omitempty"`
	// UpdatedOrder is the order to use for the assembly, if different from the quote
	UpdatedOrder *ApiExternalOrder `json:"updated_order,omitempty"`
}

// SignedQuoteResponse represents the shape of a signed quote payload directly returned by
// the auth server's API
type SignedQuoteResponse struct {
	Quote     ApiExternalQuote `json:"quote"`
	Signature string           `json:"signature"`
}
