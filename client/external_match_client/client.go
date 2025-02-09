package external_match_client //nolint:revive

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	geth_common "github.com/ethereum/go-ethereum/common"

	"github.com/renegade-fi/golang-sdk/client"
	"github.com/renegade-fi/golang-sdk/client/api_types"
	"github.com/renegade-fi/golang-sdk/wallet"
)

//nolint:revive
const (
	testnetBaseUrl        = "https://testnet.auth-server.renegade.fi:3000"
	testnetRelayerBaseUrl = "https://testnet.cluster0.renegade.fi:3000"
	mainnetBaseUrl        = "https://mainnet.auth-server.renegade.fi:3000"
	mainnetRelayerBaseUrl = "https://mainnet.cluster0.renegade.fi:3000"
	apiKeyHeader          = "X-Renegade-Api-Key" //nolint:gosec
)

// -----------------
// | Request Types |
// -----------------

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

// AssembleExternalMatchOptions represents the options for an assembly request
type AssembleExternalMatchOptions struct {
	ReceiverAddress *string
	DoGasEstimation bool
	UpdatedOrder    *api_types.ApiExternalOrder
	// RequestGasSponsorship is a flag to request gas sponsorship for the settlement tx
	//
	// This is subject to rate limit by the auth server, but if approved will refund the gas spent
	// on the settlement tx to the address specified in `GasRefundAddress`. If no refund address is
	// specified, the refund is directed to `tx.origin`
	RequestGasSponsorship bool
	// GasRefundAddress is the address to refund the gas to
	//
	// This is ignored if `RequestGasSponsorship` is false
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
	path += fmt.Sprintf("?%s=%t", api_types.RequestGasSponsorshipParam, o.RequestGasSponsorship)
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
	path += fmt.Sprintf("?%s=%t", api_types.RequestGasSponsorshipParam, o.RequestGasSponsorship)
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

// -------------------------
// | Client Implementation |
// -------------------------

// ExternalMatchClient represents a client for the external match API
//
// This client can be used to request external match bundles from a relayer.
// The relayer will return a match and a transaction to submit on-chain
type ExternalMatchClient struct {
	apiKey            string
	httpClient        *client.HttpClient
	relayerHttpClient *client.HttpClient //nolint:revive
}

// NewTestnetExternalMatchClient creates a new ExternalMatchClient for the testnet
func NewTestnetExternalMatchClient(apiKey string, apiSecret *wallet.HmacKey) *ExternalMatchClient {
	return NewExternalMatchClient(testnetBaseUrl, testnetRelayerBaseUrl, apiKey, apiSecret)
}

// NewMainnetExternalMatchClient creates a new ExternalMatchClient for the mainnet
func NewMainnetExternalMatchClient(apiKey string, apiSecret *wallet.HmacKey) *ExternalMatchClient {
	return NewExternalMatchClient(mainnetBaseUrl, mainnetRelayerBaseUrl, apiKey, apiSecret)
}

// NewExternalMatchClient creates a new ExternalMatchClient with the given base
// URL, api key, and api secret
func NewExternalMatchClient(
	baseURL string,
	relayerBaseURL string,
	apiKey string,
	apiSecret *wallet.HmacKey,
) *ExternalMatchClient {
	return &ExternalMatchClient{
		apiKey:            apiKey,
		httpClient:        client.NewHttpClient(baseURL, apiSecret),
		relayerHttpClient: client.NewHttpClient(relayerBaseURL, apiSecret),
	}
}

// GetSupportedTokens requests the list of supported tokens from the relayer
func (c *ExternalMatchClient) GetSupportedTokens() ([]api_types.ApiToken, error) {
	var response api_types.GetSupportedTokensResponse
	err := c.relayerHttpClient.GetJSON(
		api_types.GetSupportedTokensPath,
		nil, // body
		&response,
	)
	if err != nil {
		return nil, err
	}

	return response.Tokens, nil
}

// GetExternalMatchQuote requests a quote from the relayer
// returns nil if no match is found
func (c *ExternalMatchClient) GetExternalMatchQuote(
	order *api_types.ApiExternalOrder,
) (*api_types.ApiSignedQuote, error) {
	requestBody := api_types.ExternalQuoteRequest{
		ExternalOrder: *order,
	}

	var response api_types.ExternalQuoteResponse
	success, err := c.doExternalMatchRequest(
		api_types.GetExternalMatchQuotePath,
		requestBody,
		&response,
	)
	if err != nil {
		return nil, err
	}
	if !success {
		return nil, nil
	}

	return &response.Quote, nil
}

// AssembleExternalQuote generates an external match bundle from a signed quote
func (c *ExternalMatchClient) AssembleExternalQuote(
	quote *api_types.ApiSignedQuote,
) (*ExternalMatchBundle, error) {
	return c.AssembleExternalQuoteWithReceiver(quote, nil /* receiverAddress */)
}

// AssembleExternalQuoteWithReceiver generates an external match bundle from a signed quote
// returns nil if no match is found
func (c *ExternalMatchClient) AssembleExternalQuoteWithReceiver(
	quote *api_types.ApiSignedQuote,
	receiverAddress *string,
) (*ExternalMatchBundle, error) {
	options := NewAssembleExternalMatchOptions().WithReceiverAddress(receiverAddress)
	return c.AssembleExternalMatchWithOptions(quote, options)
}

// AssembleExternalMatchWithOptions assembles an external quote with the given options struct
func (c *ExternalMatchClient) AssembleExternalMatchWithOptions(
	quote *api_types.ApiSignedQuote,
	options *AssembleExternalMatchOptions,
) (*ExternalMatchBundle, error) {
	requestBody := api_types.AssembleExternalQuoteRequest{
		Quote:           *quote,
		ReceiverAddress: options.ReceiverAddress,
		DoGasEstimation: options.DoGasEstimation,
		UpdatedOrder:    options.UpdatedOrder,
	}

	var response api_types.ExternalMatchResponse
	path := options.BuildRequestPath()
	success, err := c.doExternalMatchRequest(
		path,
		requestBody,
		&response,
	)
	if err != nil {
		return nil, err
	}
	if !success {
		return nil, nil
	}

	return &ExternalMatchBundle{
		MatchResult:  &response.Bundle.MatchResult,
		Fees:         &response.Bundle.Fees,
		Receive:      &response.Bundle.Receive,
		Send:         &response.Bundle.Send,
		SettlementTx: toSettlementTransaction(&response.Bundle.SettlementTx),
		GasSponsored: response.GasSponsored,
	}, nil
}

// GetExternalMatchBundle requests an external match bundle from the relayer
// returns nil if no match is found
//
// Deprecated: Use the quote + assemble methods instead
func (c *ExternalMatchClient) GetExternalMatchBundle(
	request *api_types.ApiExternalOrder,
) (*ExternalMatchBundle, error) {
	return c.GetExternalMatchBundleWithReceiver(request, nil /* receiverAddress */)
}

// GetExternalMatchBundleWithReceiver requests an external match bundle from the relayer
// returns nil if no match is found
//
// Deprecated: Use the quote + assemble methods instead
func (c *ExternalMatchClient) GetExternalMatchBundleWithReceiver(
	request *api_types.ApiExternalOrder,
	receiverAddress *string,
) (*ExternalMatchBundle, error) {
	options := &ExternalMatchOptions{
		AssembleExternalMatchOptions: AssembleExternalMatchOptions{
			ReceiverAddress: receiverAddress,
		},
	}
	return c.GetExternalMatchBundleWithOptions(request, options)
}

// GetExternalMatchBundleWithOptions requests an external match bundle from the relayer with the given options
// returns nil if no match is found
//
// Deprecated: Use the quote + assemble methods instead
func (c *ExternalMatchClient) GetExternalMatchBundleWithOptions(
	request *api_types.ApiExternalOrder,
	options *ExternalMatchOptions,
) (*ExternalMatchBundle, error) {
	requestBody := api_types.ExternalMatchRequest{
		ExternalOrder:   *request,
		ReceiverAddress: options.ReceiverAddress,
	}

	var response api_types.ExternalMatchResponse
	path := options.BuildRequestPath()
	success, err := c.doExternalMatchRequest(
		path,
		requestBody,
		&response,
	)
	if err != nil {
		return nil, err
	}
	if !success {
		return nil, nil
	}

	return &ExternalMatchBundle{
		MatchResult:  &response.Bundle.MatchResult,
		SettlementTx: toSettlementTransaction(&response.Bundle.SettlementTx),
		GasSponsored: response.GasSponsored,
	}, nil
}

// doExternalMatchRequest handles an external match request
// returns false if the response was NO_CONTENT or if unmarshaling failed
func (c *ExternalMatchClient) doExternalMatchRequest(
	path string,
	request interface{},
	response interface{},
) (bool, error) {
	headers := make(http.Header)
	headers.Set(apiKeyHeader, c.apiKey)

	// Send the request
	statusCode, respBody, err := c.httpClient.PostWithAuthRaw(path, &headers, request)
	if err != nil {
		return false, err
	}

	// Check the status code
	if statusCode < 200 || statusCode >= 300 {
		return false, fmt.Errorf("unexpected status code: %d, body: %s", statusCode, string(respBody))
	} else if statusCode == http.StatusNoContent {
		return false, nil
	}

	// Unmarshal the response
	if err := json.Unmarshal(respBody, response); err != nil {
		return false, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return true, nil
}
