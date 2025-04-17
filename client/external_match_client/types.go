package external_match_client //nolint:revive

import (
	"fmt"
	"math/big"
	"strconv"

	geth_common "github.com/ethereum/go-ethereum/common"

	"github.com/renegade-fi/golang-sdk/client/api_types"
)

// --------------------------
// | Request/Response Types |
// --------------------------

// ExternalMatchBundle is the application level analog to the ApiExternalMatchBundle
type ExternalMatchBundle struct {
	MatchResult  *api_types.ApiExternalMatchResult
	Fees         *api_types.ApiFee
	Receive      *api_types.ApiExternalAssetTransfer
	Send         *api_types.ApiExternalAssetTransfer
	SettlementTx *SettlementTransaction
	// Whether the match has received gas sponsorship
	//
	// If `true`, the bundle is routed through a gas rebate contract that
	// refunds the gas used by the match to the configured address
	GasSponsored bool
	// The gas sponsorship info, if the match was sponsored
	GasSponsorshipInfo *api_types.ApiGasSponsorshipInfo
}

// SettlementTransaction is the application level analog to the ApiSettlementTransaction
type SettlementTransaction struct {
	Type  string
	To    geth_common.Address
	Data  []byte
	Value *big.Int
}

// toSettlementTransaction converts an ApiSettlementTransaction to a SettlementTransaction
func toSettlementTransaction(tx *api_types.ApiSettlementTransaction) *SettlementTransaction {
	// Parse a geth address and bytes data from hex strings
	to := geth_common.HexToAddress(tx.To)
	data := geth_common.FromHex(tx.Data)
	valueBytes := geth_common.FromHex(tx.Value)
	value := big.NewInt(0).SetBytes(valueBytes)

	return &SettlementTransaction{
		Type:  tx.Type,
		To:    to,
		Data:  data,
		Value: value,
	}
}

// ExternalMatchFee represents the fees for a given asset in external matches
type ExternalMatchFee struct {
	RelayerFee  float64
	ProtocolFee float64
}

// Total returns the total fee for the asset
func (f *ExternalMatchFee) Total() float64 {
	return f.RelayerFee + f.ProtocolFee
}

// toExternalMatchFee converts an ApiExternalMatchFee to an ExternalMatchFee
func toExternalMatchFee(fee *api_types.ApiExternalMatchFee) (*ExternalMatchFee, error) {
	relayerFee, err := strconv.ParseFloat(fee.RelayerFee, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse relayer fee: %w", err)
	}

	protocolFee, err := strconv.ParseFloat(fee.ProtocolFee, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse protocol fee: %w", err)
	}

	return &ExternalMatchFee{
		RelayerFee:  relayerFee,
		ProtocolFee: protocolFee,
	}, nil
}

// -----------------
// | Options Types |
// -----------------

// ExternalQuoteOptions represents the options for a quote request
type ExternalQuoteOptions struct {
	// DisableGasSponsorship is a flag to disable gas sponsorship for the quote
	//
	// This is subject to rate limit by the auth server, but if approved will refund the gas spent
	// on the settlement tx to the address specified in `GasRefundAddress`, or the associated default
	// if no refund address is specified.
	DisableGasSponsorship bool
	// GasRefundAddress is the address to refund the gas to. If unspecified, then in the case of a
	// native ETH refund, defaults to `tx.origin`, and in the case of an in-kind refund, defaults to
	// the receiver address.
	GasRefundAddress *string
	// RefundNativeEth is a flag to request a receiving the gas sponsorship refund
	// in terms of native ETH, as opposed to the buy-side token ("in-kind" sponsorship).
	RefundNativeEth bool
}

// WithDisableGasSponsorship sets whether to disable gas sponsorship
func (o *ExternalQuoteOptions) WithDisableGasSponsorship(disable bool) *ExternalQuoteOptions {
	o.DisableGasSponsorship = disable
	return o
}

// WithGasRefundAddress sets the gas refund address for the quote options
func (o *ExternalQuoteOptions) WithGasRefundAddress(address *string) *ExternalQuoteOptions {
	o.GasRefundAddress = address
	return o
}

// WithRefundNativeEth sets whether to request a native ETH refund
func (o *ExternalQuoteOptions) WithRefundNativeEth(refundNativeEth bool) *ExternalQuoteOptions {
	o.RefundNativeEth = refundNativeEth
	return o
}

// BuildRequestPath builds the request path for the quote options
func (o *ExternalQuoteOptions) BuildRequestPath() string {
	path := api_types.GetExternalMatchQuotePath
	path += fmt.Sprintf("?%s=%t", api_types.DisableGasSponsorshipParam, o.DisableGasSponsorship)
	path += fmt.Sprintf("&%s=%t", api_types.RefundNativeEthParam, o.RefundNativeEth)
	if o.GasRefundAddress != nil {
		path += fmt.Sprintf("&%s=%s", api_types.GasRefundAddressParam, *o.GasRefundAddress)
	}

	return path
}

// NewExternalQuoteOptions creates a new ExternalQuoteOptions with default values
func NewExternalQuoteOptions() *ExternalQuoteOptions {
	return &ExternalQuoteOptions{
		DisableGasSponsorship: false,
		GasRefundAddress:      nil,
		RefundNativeEth:       false,
	}
}

// AssembleExternalMatchOptions represents the options for an assembly request
type AssembleExternalMatchOptions struct {
	ReceiverAddress *string
	DoGasEstimation bool
	// AllowShared is a flag to allow the assembly of a shared quote
	//
	// If true, the relayer will not enforce exclusive access to the bundle returned in the
	// assemble step. I.e. the relayer may send the same bundle to another client.
	//
	// This affords the client a much higher rate limit
	AllowShared  bool
	UpdatedOrder *api_types.ApiExternalOrder
	// RequestGasSponsorship is a flag to request gas sponsorship for the settlement tx
	//
	// This is subject to rate limit by the auth server, but if approved will refund the gas spent
	// on the settlement tx to the address specified in `GasRefundAddress`. If no refund address is
	// specified, the refund is directed to `tx.origin`
	//
	// Deprecated: Request gas sponsorship when requesting a quote
	RequestGasSponsorship bool
	// GasRefundAddress is the address to refund the gas to
	//
	// This is ignored if `RequestGasSponsorship` is false
	//
	// Deprecated: Request gas sponsorship when requesting a quote
	GasRefundAddress *string
}

// WithReceiverAddress sets the receiver address for the assembly options
func (o *AssembleExternalMatchOptions) WithReceiverAddress(address *string) *AssembleExternalMatchOptions {
	o.ReceiverAddress = address
	return o
}

// WithGasEstimation sets whether to perform gas estimation
func (o *AssembleExternalMatchOptions) WithGasEstimation(estimate bool) *AssembleExternalMatchOptions {
	o.DoGasEstimation = estimate
	return o
}

// WithAllowShared sets whether to allow the assembly of a shared quote
func (o *AssembleExternalMatchOptions) WithAllowShared(allowShared bool) *AssembleExternalMatchOptions {
	o.AllowShared = allowShared
	return o
}

// WithUpdatedOrder sets the updated order for the assembly options
func (o *AssembleExternalMatchOptions) WithUpdatedOrder(order *api_types.ApiExternalOrder) *AssembleExternalMatchOptions {
	o.UpdatedOrder = order
	return o
}

// WithRequestGasSponsorship sets whether to request gas sponsorship
func (o *AssembleExternalMatchOptions) WithRequestGasSponsorship(request bool) *AssembleExternalMatchOptions {
	o.RequestGasSponsorship = request
	return o
}

// WithGasRefundAddress sets the gas refund address for the assembly options
func (o *AssembleExternalMatchOptions) WithGasRefundAddress(address *string) *AssembleExternalMatchOptions {
	o.GasRefundAddress = address
	return o
}

// BuildRequestPath builds the request path for the assembly options
func (o *AssembleExternalMatchOptions) BuildRequestPath() string {
	path := api_types.AssembleExternalQuotePath
	if o.RequestGasSponsorship {
		// We only write this query parameter if it was explicitly set. The
		// expectation of the auth server is that when gas sponsorship is
		// requested at the quote stage, there should be no query parameters
		// at all in the assemble request.
		path += fmt.Sprintf("?%s=%t", api_types.DisableGasSponsorshipParam, !o.RequestGasSponsorship)
	}
	if o.GasRefundAddress != nil {
		path += fmt.Sprintf("&%s=%s", api_types.GasRefundAddressParam, *o.GasRefundAddress)
	}

	return path
}

// NewAssembleExternalMatchOptions creates a new AssembleExternalMatchOptions with default values
func NewAssembleExternalMatchOptions() *AssembleExternalMatchOptions {
	return &AssembleExternalMatchOptions{
		ReceiverAddress: nil,
		DoGasEstimation: false,
		UpdatedOrder:    nil,
	}
}

// ExternalMatchOptions represents the options for an external match request
//
// Deprecated: Use AssembleExternalMatchOptions instead
type ExternalMatchOptions struct {
	AssembleExternalMatchOptions
}

// BuildRequestPath builds the request path for the external match options
func (o *ExternalMatchOptions) BuildRequestPath() string {
	path := api_types.GetExternalMatchBundlePath
	path += fmt.Sprintf("?%s=%t", api_types.DisableGasSponsorshipParam, !o.RequestGasSponsorship)
	if o.GasRefundAddress != nil {
		path += fmt.Sprintf("&%s=%s", api_types.GasRefundAddressParam, *o.GasRefundAddress)
	}
	return path
}

// NewExternalMatchOptions creates a new ExternalMatchOptions with default values
func NewExternalMatchOptions() *ExternalMatchOptions {
	return &ExternalMatchOptions{
		AssembleExternalMatchOptions: *NewAssembleExternalMatchOptions(),
	}
}
