package api_types //nolint:revive

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
)

// ---------------------
// | V2 Route Constants |
// ---------------------

//nolint:revive
const (
	// GetMarketsPath is the path for fetching all tradable markets
	GetMarketsPath = "/v2/markets"
	// GetMarketsDepthPath is the path for fetching depth of all markets
	GetMarketsDepthPath = "/v2/markets/depth"
	// GetMarketDepthByMintPath is the path for fetching depth of a specific market
	// Use fmt.Sprintf with the mint address
	GetMarketDepthByMintPath = "/v2/markets/%s/depth"
	// GetQuoteV2Path is the path for requesting a v2 quote
	GetQuoteV2Path = "/v2/external-matches/get-quote"
	// AssembleMatchBundleV2Path is the path for assembling a v2 match bundle
	AssembleMatchBundleV2Path = "/v2/external-matches/assemble-match-bundle"
	// GetExchangeMetadataPath is the path for fetching exchange metadata
	GetExchangeMetadataPath = "/v2/metadata/exchange"
)

// BuildGetMarketDepthByMintPath builds the path for fetching the market depth for a specific mint
func BuildGetMarketDepthByMintPath(mint string) string {
	return fmt.Sprintf(GetMarketDepthByMintPath, mint)
}

// -------------------------
// | Serialization Helpers |
// -------------------------

// StringAmount is a big.Int wrapper that marshals/unmarshals as a quoted JSON string.
// This is needed because v2 wire format uses JSON strings for amounts (e.g. "100")
// while v1's Amount type marshals as bare numbers.
type StringAmount big.Int

// NewStringAmount creates a new StringAmount from an int64
func NewStringAmount(i int64) StringAmount {
	return StringAmount(*big.NewInt(i))
}

// NewStringAmountFromBigInt creates a new StringAmount from a *big.Int
func NewStringAmountFromBigInt(i *big.Int) StringAmount {
	return StringAmount(*i)
}

// ToBigInt converts a StringAmount to a *big.Int
func (a *StringAmount) ToBigInt() *big.Int {
	return (*big.Int)(a)
}

// IsZero returns true if the amount is zero
func (a *StringAmount) IsZero() bool {
	return (*big.Int)(a).Sign() == 0
}

// MarshalJSON marshals the StringAmount as a quoted JSON string
func (a StringAmount) MarshalJSON() ([]byte, error) {
	s := (*big.Int)(&a).String()
	return json.Marshal(s)
}

// UnmarshalJSON unmarshals the StringAmount from a quoted JSON string
func (a *StringAmount) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	i, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return fmt.Errorf("invalid StringAmount: %s", s)
	}
	*a = StringAmount(*i)
	return nil
}

// StringFloat is a float64 wrapper that marshals/unmarshals as a quoted JSON string.
// Used for fields like DepthSide.TotalQuantityUSD.
type StringFloat float64

// MarshalJSON marshals the StringFloat as a quoted JSON string
func (f StringFloat) MarshalJSON() ([]byte, error) {
	s := strconv.FormatFloat(float64(f), 'f', -1, 64)
	return json.Marshal(s)
}

// UnmarshalJSON unmarshals the StringFloat from a quoted JSON string
func (f *StringFloat) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("invalid StringFloat: %s", s)
	}
	*f = StringFloat(val)
	return nil
}

// ----------------
// | FixedPoint   |
// ----------------

// fixedPointPrecisionBits is the number of bits used for fixed-point precision
const fixedPointPrecisionBits = 63

// fixedPointPrecisionShift returns 2^63 as a *big.Int
func fixedPointPrecisionShift() *big.Int {
	shift := new(big.Int).SetUint64(1)
	return shift.Lsh(shift, fixedPointPrecisionBits)
}

// FixedPoint is a fixed-point number with 63-bit precision.
// The value represents the number multiplied by 2^63.
type FixedPoint struct {
	Value *big.Int
}

// NewFixedPoint creates a new FixedPoint from a *big.Int value
func NewFixedPoint(value *big.Int) FixedPoint {
	return FixedPoint{Value: value}
}

// FloorMulInt multiplies this fixed-point by an integer amount and returns the floor.
// Result = (value * amount) / 2^63
func (fp *FixedPoint) FloorMulInt(amount *big.Int) *big.Int {
	product := new(big.Int).Mul(fp.Value, amount)
	return new(big.Int).Div(product, fixedPointPrecisionShift())
}

// CeilDivInt divides an amount by this fixed-point and returns the ceiling.
// Result = ceil(amount * 2^63 / value)
func CeilDivInt(amount *big.Int, fp *FixedPoint) *big.Int {
	numerator := new(big.Int).Mul(amount, fixedPointPrecisionShift())
	quotient := new(big.Int)
	remainder := new(big.Int)
	quotient.DivMod(numerator, fp.Value, remainder)
	if remainder.Sign() > 0 {
		quotient.Add(quotient, big.NewInt(1))
	}
	return quotient
}

// ToF64 converts the fixed-point number to a float64 approximation.
// Result = value / 2^63
func (fp *FixedPoint) ToF64() float64 {
	valueFloat := new(big.Float).SetInt(fp.Value)
	shiftFloat := new(big.Float).SetInt(fixedPointPrecisionShift())
	result := new(big.Float).Quo(valueFloat, shiftFloat)
	f, _ := result.Float64()
	return f
}

// Add adds two fixed-point numbers
func (fp *FixedPoint) Add(other *FixedPoint) FixedPoint {
	sum := new(big.Int).Add(fp.Value, other.Value)
	return FixedPoint{Value: sum}
}

// MarshalJSON serializes the FixedPoint as a quoted decimal string
func (fp FixedPoint) MarshalJSON() ([]byte, error) {
	s := fp.Value.String()
	return json.Marshal(s)
}

// UnmarshalJSON deserializes the FixedPoint from a quoted decimal string
func (fp *FixedPoint) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	i, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return fmt.Errorf("invalid FixedPoint: %s", s)
	}
	fp.Value = i
	return nil
}

// ---------------
// | Order Types |
// ---------------

// ApiExternalOrderV2 is a v2 external order using input/output semantics
type ApiExternalOrderV2 struct { //nolint:revive
	// The mint (erc20 address) of the input token
	InputMint string `json:"input_mint"`
	// The mint (erc20 address) of the output token
	OutputMint string `json:"output_mint"`
	// The input amount
	InputAmount StringAmount `json:"input_amount"`
	// The output amount
	OutputAmount StringAmount `json:"output_amount"`
	// Whether to use exact output amount
	UseExactOutputAmount bool `json:"use_exact_output_amount"`
	// The minimum fill size
	MinFillSize StringAmount `json:"min_fill_size"`
}

// ApiExternalOrderBuilderV2 helps construct ApiExternalOrderV2 with validation
type ApiExternalOrderBuilderV2 struct { //nolint:revive
	order ApiExternalOrderV2
}

// NewExternalOrderBuilderV2 creates a new v2 order builder
func NewExternalOrderBuilderV2() *ApiExternalOrderBuilderV2 {
	return &ApiExternalOrderBuilderV2{
		order: ApiExternalOrderV2{
			InputAmount:  NewStringAmount(0),
			OutputAmount: NewStringAmount(0),
			MinFillSize:  NewStringAmount(0),
		},
	}
}

// WithInputMint sets the input mint
func (b *ApiExternalOrderBuilderV2) WithInputMint(mint string) *ApiExternalOrderBuilderV2 {
	b.order.InputMint = mint
	return b
}

// WithOutputMint sets the output mint
func (b *ApiExternalOrderBuilderV2) WithOutputMint(mint string) *ApiExternalOrderBuilderV2 {
	b.order.OutputMint = mint
	return b
}

// WithInputAmount sets the input amount
func (b *ApiExternalOrderBuilderV2) WithInputAmount(amount StringAmount) *ApiExternalOrderBuilderV2 {
	b.order.InputAmount = amount
	return b
}

// WithOutputAmount sets the output amount
func (b *ApiExternalOrderBuilderV2) WithOutputAmount(amount StringAmount) *ApiExternalOrderBuilderV2 {
	b.order.OutputAmount = amount
	return b
}

// WithExactOutputAmount sets the use exact output amount flag
func (b *ApiExternalOrderBuilderV2) WithExactOutputAmount(exact bool) *ApiExternalOrderBuilderV2 {
	b.order.UseExactOutputAmount = exact
	return b
}

// WithMinFillSize sets the minimum fill size
func (b *ApiExternalOrderBuilderV2) WithMinFillSize(size StringAmount) *ApiExternalOrderBuilderV2 {
	b.order.MinFillSize = size
	return b
}

// Build validates and returns the ApiExternalOrderV2
func (b *ApiExternalOrderBuilderV2) Build() (*ApiExternalOrderV2, error) {
	if b.order.InputMint == "" {
		return nil, errors.New("input mint is required")
	}
	if b.order.OutputMint == "" {
		return nil, errors.New("output mint is required")
	}
	if b.order.InputAmount.IsZero() && b.order.OutputAmount.IsZero() {
		return nil, errors.New("one of input_amount or output_amount must be set")
	}
	return &b.order, nil
}

// ----------------------
// | Match Result Types |
// ----------------------

// ApiTimestampedPriceFp is a timestamped price with full fixed-point precision
type ApiTimestampedPriceFp struct { //nolint:revive
	Price     FixedPoint `json:"price"`
	Timestamp uint64     `json:"timestamp"`
}

// ApiExternalMatchResultV2 is the v2 match result with input/output semantics
type ApiExternalMatchResultV2 struct { //nolint:revive
	InputMint    string                `json:"input_mint"`
	OutputMint   string                `json:"output_mint"`
	InputAmount  StringAmount          `json:"input_amount"`
	OutputAmount StringAmount          `json:"output_amount"`
	PriceFp      ApiTimestampedPriceFp `json:"price_fp"`
}

// ApiBoundedMatchResultV2 is a bounded match result for malleable matches
type ApiBoundedMatchResultV2 struct { //nolint:revive
	InputMint      string       `json:"input_mint"`
	OutputMint     string       `json:"output_mint"`
	PriceFp        FixedPoint   `json:"price_fp"`
	MinInputAmount StringAmount `json:"min_input_amount"`
	MaxInputAmount StringAmount `json:"max_input_amount"`
}

// ---------------
// | Fee Types   |
// ---------------

// FeeTake represents the fee amounts paid to the relayer and protocol
type FeeTake struct {
	RelayerFee  StringAmount `json:"relayer_fee"`
	ProtocolFee StringAmount `json:"protocol_fee"`
}

// Total returns the total fee
func (f *FeeTake) Total() *big.Int {
	return new(big.Int).Add(f.RelayerFee.ToBigInt(), f.ProtocolFee.ToBigInt())
}

// FeeTakeRate represents the fee rates for relayer and protocol
type FeeTakeRate struct {
	RelayerFeeRate  FixedPoint `json:"relayer_fee_rate"`
	ProtocolFeeRate FixedPoint `json:"protocol_fee_rate"`
}

// Total returns the total fee rate
func (f *FeeTakeRate) Total() FixedPoint {
	return f.RelayerFeeRate.Add(&f.ProtocolFeeRate)
}

// --------------------------
// | Asset Transfer (V2)    |
// --------------------------

// ApiExternalAssetTransferV2 represents a v2 asset transfer with string amounts
type ApiExternalAssetTransferV2 struct { //nolint:revive
	Mint   string       `json:"mint"`
	Amount StringAmount `json:"amount"`
}

// ---------------
// | Quote Types |
// ---------------

// ApiExternalQuoteV2 is a v2 quote from the relayer
type ApiExternalQuoteV2 struct { //nolint:revive
	Order       ApiExternalOrderV2         `json:"order"`
	MatchResult ApiExternalMatchResultV2   `json:"match_result"`
	Fees        FeeTake                    `json:"fees"`
	Send        ApiExternalAssetTransferV2 `json:"send"`
	Receive     ApiExternalAssetTransferV2 `json:"receive"`
	Price       TimestampedPrice           `json:"price"`
	Timestamp   uint64                     `json:"timestamp"`
}

// ApiSignedQuoteV2 is a signed v2 quote from the relayer
type ApiSignedQuoteV2 struct { //nolint:revive
	Quote     ApiExternalQuoteV2 `json:"quote"`
	Signature string             `json:"signature"`
	Deadline  uint64             `json:"deadline"`
}

// ----------------
// | Bundle Types |
// ----------------

// ApiSettlementTransactionV2 is the v2 settlement tx format matching alloy's TransactionRequest.
// Uses "input" instead of "data" for the calldata field, and fields are optional.
type ApiSettlementTransactionV2 struct { //nolint:revive
	To    *string `json:"to,omitempty"`
	Input string  `json:"input,omitempty"`
	Value *string `json:"value,omitempty"`
	Gas   *string `json:"gas,omitempty"`
}

// ToV1 converts a v2 settlement tx to the v1 wire format
func (tx *ApiSettlementTransactionV2) ToV1() ApiSettlementTransaction {
	to := ""
	if tx.To != nil {
		to = *tx.To
	}
	value := ""
	if tx.Value != nil {
		value = *tx.Value
	}
	gas := ""
	if tx.Gas != nil {
		gas = *tx.Gas
	}
	return ApiSettlementTransaction{
		To:    to,
		Data:  tx.Input,
		Value: value,
		Gas:   gas,
	}
}

// MalleableAtomicMatchApiBundleV2 contains a malleable match bundle
type MalleableAtomicMatchApiBundleV2 struct { //nolint:revive
	MatchResult  ApiBoundedMatchResultV2    `json:"match_result"`
	FeeRates     FeeTakeRate                `json:"fee_rates"`
	MaxReceive   ApiExternalAssetTransferV2 `json:"max_receive"`
	MinReceive   ApiExternalAssetTransferV2 `json:"min_receive"`
	MaxSend      ApiExternalAssetTransferV2 `json:"max_send"`
	MinSend      ApiExternalAssetTransferV2 `json:"min_send"`
	SettlementTx ApiSettlementTransactionV2 `json:"settlement_tx"`
	Deadline     uint64                     `json:"deadline"`
}

// -------------------------
// | Request/Response Types |
// -------------------------

// ExternalQuoteRequestV2 is the request body for a v2 quote
type ExternalQuoteRequestV2 struct {
	ExternalOrder ApiExternalOrderV2 `json:"external_order"`
}

// ExternalQuoteResponseV2 is the response body for a v2 quote
type ExternalQuoteResponseV2 struct {
	SignedQuote        ApiSignedQuoteV2       `json:"signed_quote"`
	GasSponsorshipInfo *ApiGasSponsorshipInfo `json:"gas_sponsorship_info,omitempty"`
}

// AssemblyType represents the tagged union for the assembly request order field.
// Uses flat struct with omitempty to produce correct JSON for either variant.
type AssemblyType struct {
	Type          string              `json:"type"`                     // "quoted-order" or "direct-order"
	SignedQuote   *ApiSignedQuoteV2   `json:"signed_quote,omitempty"`   // for quoted-order
	UpdatedOrder  *ApiExternalOrderV2 `json:"updated_order,omitempty"`  // for quoted-order (optional)
	ExternalOrder *ApiExternalOrderV2 `json:"external_order,omitempty"` // for direct-order
}

// NewQuotedOrderAssembly creates an AssemblyType for a quoted order
func NewQuotedOrderAssembly(quote *ApiSignedQuoteV2, updatedOrder *ApiExternalOrderV2) AssemblyType {
	return AssemblyType{
		Type:         "quoted-order",
		SignedQuote:  quote,
		UpdatedOrder: updatedOrder,
	}
}

// NewDirectOrderAssembly creates an AssemblyType for a direct order
func NewDirectOrderAssembly(order *ApiExternalOrderV2) AssemblyType {
	return AssemblyType{
		Type:          "direct-order",
		ExternalOrder: order,
	}
}

// AssembleExternalMatchRequestV2 is the request body for a v2 assembly
type AssembleExternalMatchRequestV2 struct {
	DoGasEstimation bool         `json:"do_gas_estimation"`
	ReceiverAddress *string      `json:"receiver_address,omitempty"`
	Order           AssemblyType `json:"order"`
}

// ExternalMatchResponseV2 is the response body for a v2 match
type ExternalMatchResponseV2 struct {
	MatchBundle        MalleableAtomicMatchApiBundleV2 `json:"match_bundle"`
	GasSponsorshipInfo *ApiGasSponsorshipInfo          `json:"gas_sponsorship_info,omitempty"`
}

// ----------------------
// | Market Data Types  |
// ----------------------

// MarketInfo represents information about a tradable market
type MarketInfo struct {
	Base                  ApiToken         `json:"base"`
	Quote                 ApiToken         `json:"quote"`
	Price                 TimestampedPrice `json:"price"`
	InternalMatchFeeRates FeeTakeRate      `json:"internal_match_fee_rates"`
	ExternalMatchFeeRates FeeTakeRate      `json:"external_match_fee_rates"`
}

// DepthSide represents the liquidity depth for one side of a market
type DepthSide struct {
	TotalQuantity    StringAmount `json:"total_quantity"`
	TotalQuantityUSD StringFloat  `json:"total_quantity_usd"`
}

// MarketDepth represents the full depth of a market
type MarketDepth struct {
	Market MarketInfo `json:"market"`
	Buy    DepthSide  `json:"buy"`
	Sell   DepthSide  `json:"sell"`
}

// GetMarketsResponse is the response for the GetMarkets endpoint
type GetMarketsResponse struct {
	Markets []MarketInfo `json:"markets"`
}

// GetMarketDepthByMintResponse is the response for the GetMarketDepthByMint endpoint
type GetMarketDepthByMintResponse struct {
	MarketDepth MarketDepth `json:"market_depth"`
}

// GetMarketDepthsResponse is the response for the GetMarketDepths endpoint
type GetMarketDepthsResponse struct {
	MarketDepths []MarketDepth `json:"market_depths"`
}

// ExchangeMetadataResponse is the response for the GetExchangeMetadata endpoint
type ExchangeMetadataResponse struct {
	ChainID                   uint64     `json:"chain_id"`
	SettlementContractAddress string     `json:"settlement_contract_address"`
	ExecutorAddress           string     `json:"executor_address"`
	RelayerFeeRecipient       string     `json:"relayer_fee_recipient"`
	SupportedTokens           []ApiToken `json:"supported_tokens"`
}
