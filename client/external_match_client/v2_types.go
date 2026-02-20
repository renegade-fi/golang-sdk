package external_match_client //nolint:revive

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/renegade-fi/golang-sdk/client/api_types"
)

// NativeAssetAddr is the sentinel address for native ETH
const NativeAssetAddr = "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"

// inputAmountOffset is the byte offset of the input amount in settlement tx calldata
// (after the 4-byte function selector)
const inputAmountOffset = 4

// amountCalldataLength is the length of a uint256 in calldata (32 bytes)
const amountCalldataLength = 32

// --------------------------
// | SignedExternalQuoteV2 |
// --------------------------

// SignedExternalQuoteV2 is the application-level v2 signed quote
type SignedExternalQuoteV2 struct {
	Quote              api_types.ApiExternalQuoteV2
	Signature          string
	Deadline           uint64
	GasSponsorshipInfo *api_types.ApiGasSponsorshipInfo
}

// NewSignedExternalQuoteV2 creates a SignedExternalQuoteV2 from an API response
func NewSignedExternalQuoteV2(resp *api_types.ExternalQuoteResponseV2) *SignedExternalQuoteV2 {
	return &SignedExternalQuoteV2{
		Quote:              resp.SignedQuote.Quote,
		Signature:          resp.SignedQuote.Signature,
		Deadline:           resp.SignedQuote.Deadline,
		GasSponsorshipInfo: resp.GasSponsorshipInfo,
	}
}

// MatchResult returns the match result from the quote
func (q *SignedExternalQuoteV2) MatchResult() api_types.ApiExternalMatchResultV2 {
	return q.Quote.MatchResult
}

// Fees returns the fees from the quote
func (q *SignedExternalQuoteV2) Fees() api_types.FeeTake {
	return q.Quote.Fees
}

// ReceiveAmount returns the receive transfer from the quote
func (q *SignedExternalQuoteV2) ReceiveAmount() api_types.ApiExternalAssetTransferV2 {
	return q.Quote.Receive
}

// SendAmount returns the send transfer from the quote
func (q *SignedExternalQuoteV2) SendAmount() api_types.ApiExternalAssetTransferV2 {
	return q.Quote.Send
}

// ToApiSignedQuote converts to the API wire format (without gas info)
func (q *SignedExternalQuoteV2) ToApiSignedQuote() api_types.ApiSignedQuoteV2 {
	return api_types.ApiSignedQuoteV2{
		Quote:     q.Quote,
		Signature: q.Signature,
		Deadline:  q.Deadline,
	}
}

// ---------------------------------
// | MalleableExternalMatchBundle |
// ---------------------------------

// MalleableExternalMatchBundle is the application-level v2 match bundle
// with support for malleable (bounded) input amounts
type MalleableExternalMatchBundle struct {
	MatchResult        *api_types.ApiBoundedMatchResultV2
	FeeRates           *api_types.FeeTakeRate
	MaxReceive         *api_types.ApiExternalAssetTransferV2
	MinReceive         *api_types.ApiExternalAssetTransferV2
	MaxSend            *api_types.ApiExternalAssetTransferV2
	MinSend            *api_types.ApiExternalAssetTransferV2
	SettlementTx       *SettlementTransaction
	Deadline           uint64
	GasSponsorshipInfo *api_types.ApiGasSponsorshipInfo
	inputAmount        *big.Int // set by SetInputAmount; unexported
}

// newMalleableExternalMatchBundle creates a MalleableExternalMatchBundle from an API response
func newMalleableExternalMatchBundle(resp *api_types.ExternalMatchResponseV2) *MalleableExternalMatchBundle {
	return &MalleableExternalMatchBundle{
		MatchResult:        &resp.MatchBundle.MatchResult,
		FeeRates:           &resp.MatchBundle.FeeRates,
		MaxReceive:         &resp.MatchBundle.MaxReceive,
		MinReceive:         &resp.MatchBundle.MinReceive,
		MaxSend:            &resp.MatchBundle.MaxSend,
		MinSend:            &resp.MatchBundle.MinSend,
		SettlementTx:       toSettlementTransactionV2(&resp.MatchBundle.SettlementTx),
		Deadline:           resp.MatchBundle.Deadline,
		GasSponsorshipInfo: resp.GasSponsorshipInfo,
	}
}

// InputBounds returns the (min, max) input amount bounds
func (b *MalleableExternalMatchBundle) InputBounds() (min, max *big.Int) {
	return b.MatchResult.MinInputAmount.ToBigInt(), b.MatchResult.MaxInputAmount.ToBigInt()
}

// OutputBounds returns the (min, max) output amount bounds
// Computed from the price and input bounds
func (b *MalleableExternalMatchBundle) OutputBounds() (min, max *big.Int) {
	minInput, maxInput := b.InputBounds()
	price := &b.MatchResult.PriceFp

	minOutput := price.FloorMulInt(minInput)
	maxOutput := price.FloorMulInt(maxInput)

	return minOutput, maxOutput
}

// currentInputAmount returns the currently set input amount, defaulting to max
func (b *MalleableExternalMatchBundle) currentInputAmount() *big.Int {
	if b.inputAmount != nil {
		return b.inputAmount
	}
	return b.MatchResult.MaxInputAmount.ToBigInt()
}

// outputAmount computes the output amount at the given input
func (b *MalleableExternalMatchBundle) outputAmount(inputAmount *big.Int) *big.Int {
	return b.MatchResult.PriceFp.FloorMulInt(inputAmount)
}

// computeReceiveAmount computes the receive (output) amount net of fees
func (b *MalleableExternalMatchBundle) computeReceiveAmount(inputAmount *big.Int) *big.Int {
	preSponsoredAmt := b.outputAmount(inputAmount)

	// Subtract fees
	totalFeeRate := b.FeeRates.Total()
	totalFeeAmount := totalFeeRate.FloorMulInt(preSponsoredAmt)
	preSponsoredAmt = new(big.Int).Sub(preSponsoredAmt, totalFeeAmount)

	// Add gas sponsorship refund if in-kind (not native ETH)
	if b.GasSponsorshipInfo != nil && !b.GasSponsorshipInfo.RefundNativeETH {
		refund := (*big.Int)(&b.GasSponsorshipInfo.RefundAmount)
		preSponsoredAmt = new(big.Int).Add(preSponsoredAmt, refund)
	}

	return preSponsoredAmt
}

// ReceiveAmount returns the receive amount at the currently set input amount
func (b *MalleableExternalMatchBundle) ReceiveAmount() *big.Int {
	return b.computeReceiveAmount(b.currentInputAmount())
}

// ReceiveAmountAtInput returns the receive amount at a specific input amount
func (b *MalleableExternalMatchBundle) ReceiveAmountAtInput(inputAmount *big.Int) *big.Int {
	return b.computeReceiveAmount(inputAmount)
}

// SendAmount returns the current send amount
func (b *MalleableExternalMatchBundle) SendAmount() *big.Int {
	return b.currentInputAmount()
}

// isNativeEthSell returns whether the trade is a native ETH sell
func (b *MalleableExternalMatchBundle) isNativeEthSell() bool {
	return strings.EqualFold(b.MatchResult.InputMint, NativeAssetAddr)
}

// checkInputAmount validates that the input amount is within bounds
func (b *MalleableExternalMatchBundle) checkInputAmount(inputAmount *big.Int) error {
	minInput, maxInput := b.InputBounds()
	if inputAmount.Cmp(minInput) < 0 || inputAmount.Cmp(maxInput) > 0 {
		return fmt.Errorf("invalid input amount: must be between %s and %s, got %s",
			minInput.String(), maxInput.String(), inputAmount.String())
	}
	return nil
}

// SetInputAmount sets the input amount, modifies the settlement tx calldata,
// and returns the resulting receive amount.
// The amount must be within the input bounds.
func (b *MalleableExternalMatchBundle) SetInputAmount(amount *big.Int) (*big.Int, error) {
	if err := b.checkInputAmount(amount); err != nil {
		return nil, err
	}

	// Modify the calldata
	b.setInputAmountCalldata(amount)

	// Store the input amount
	b.inputAmount = new(big.Int).Set(amount)

	return b.ReceiveAmount(), nil
}

// setInputAmountCalldata writes the input amount into the settlement tx calldata
func (b *MalleableExternalMatchBundle) setInputAmountCalldata(inputAmount *big.Int) {
	data := b.SettlementTx.Data
	end := inputAmountOffset + amountCalldataLength

	// ABI-encode the amount as a 32-byte big-endian uint256
	amountBytes := inputAmount.Bytes()
	encoded := make([]byte, amountCalldataLength)
	// Left-pad with zeros
	copy(encoded[amountCalldataLength-len(amountBytes):], amountBytes)

	// Write into calldata at bytes 4-36
	copy(data[inputAmountOffset:end], encoded)

	// If native ETH sell, also update the tx value
	if b.isNativeEthSell() {
		b.SettlementTx.Value = new(big.Int).Set(inputAmount)
	}
}

// GetSettlementTx returns the parsed settlement transaction
func (b *MalleableExternalMatchBundle) GetSettlementTx() *SettlementTransaction {
	return b.SettlementTx
}

// -----------------
// | Options Types |
// -----------------

// AssembleExternalMatchOptionsV2 represents options for a v2 assembly request
type AssembleExternalMatchOptionsV2 struct {
	DoGasEstimation bool
	ReceiverAddress *string
	UpdatedOrder    *api_types.ApiExternalOrderV2
}

// NewAssembleExternalMatchOptionsV2 creates default v2 assembly options
func NewAssembleExternalMatchOptionsV2() *AssembleExternalMatchOptionsV2 {
	return &AssembleExternalMatchOptionsV2{}
}

// WithGasEstimation sets the gas estimation flag
func (o *AssembleExternalMatchOptionsV2) WithGasEstimation(estimate bool) *AssembleExternalMatchOptionsV2 {
	o.DoGasEstimation = estimate
	return o
}

// WithReceiverAddress sets the receiver address
func (o *AssembleExternalMatchOptionsV2) WithReceiverAddress(address *string) *AssembleExternalMatchOptionsV2 {
	o.ReceiverAddress = address
	return o
}

// WithUpdatedOrder sets the updated order
func (o *AssembleExternalMatchOptionsV2) WithUpdatedOrder(order *api_types.ApiExternalOrderV2) *AssembleExternalMatchOptionsV2 {
	o.UpdatedOrder = order
	return o
}

// ExternalMatchOptionsV2 represents options for a v2 direct match request
type ExternalMatchOptionsV2 struct {
	DoGasEstimation       bool
	ReceiverAddress       *string
	DisableGasSponsorship bool
	GasRefundAddress      *string
	RefundNativeEth       bool
}

// NewExternalMatchOptionsV2 creates default v2 match options
func NewExternalMatchOptionsV2() *ExternalMatchOptionsV2 {
	return &ExternalMatchOptionsV2{}
}

// WithGasEstimation sets the gas estimation flag
func (o *ExternalMatchOptionsV2) WithGasEstimation(estimate bool) *ExternalMatchOptionsV2 {
	o.DoGasEstimation = estimate
	return o
}

// WithReceiverAddress sets the receiver address
func (o *ExternalMatchOptionsV2) WithReceiverAddress(address *string) *ExternalMatchOptionsV2 {
	o.ReceiverAddress = address
	return o
}

// WithDisableGasSponsorship disables gas sponsorship
func (o *ExternalMatchOptionsV2) WithDisableGasSponsorship(disable bool) *ExternalMatchOptionsV2 {
	o.DisableGasSponsorship = disable
	return o
}

// WithGasRefundAddress sets the gas refund address
func (o *ExternalMatchOptionsV2) WithGasRefundAddress(address *string) *ExternalMatchOptionsV2 {
	o.GasRefundAddress = address
	return o
}

// WithRefundNativeEth sets whether to refund in native ETH
func (o *ExternalMatchOptionsV2) WithRefundNativeEth(refund bool) *ExternalMatchOptionsV2 {
	o.RefundNativeEth = refund
	return o
}

// BuildRequestPath builds the request path for the v2 match options
func (o *ExternalMatchOptionsV2) BuildRequestPath() string {
	path := api_types.AssembleMatchBundleV2Path
	path += fmt.Sprintf("?%s=%t", api_types.DisableGasSponsorshipParam, o.DisableGasSponsorship)
	if o.RefundNativeEth {
		path += fmt.Sprintf("&%s=%t", api_types.RefundNativeEthParam, o.RefundNativeEth)
	}
	if o.GasRefundAddress != nil {
		path += fmt.Sprintf("&%s=%s", api_types.GasRefundAddressParam, *o.GasRefundAddress)
	}
	return path
}
