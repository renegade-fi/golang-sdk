package external_match_client //nolint:revive

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/renegade-fi/golang-sdk/client"
	"github.com/renegade-fi/golang-sdk/client/api_types"
	"github.com/renegade-fi/golang-sdk/wallet"
)

//nolint:revive
const (
	arbitrumSepoliaBaseUrl        = "https://arbitrum-sepolia.v2.auth-server.renegade.fi"
	arbitrumSepoliaRelayerBaseUrl = "https://arbitrum-sepolia.v2.relayer.renegade.fi"
	arbitrumOneBaseUrl            = "https://arbitrum-one.v2.auth-server.renegade.fi"
	arbitrumOneRelayerBaseUrl     = "https://arbitrum-one.v2.relayer.renegade.fi"
	baseSepoliaBaseUrl            = "https://base-sepolia.v2.auth-server.renegade.fi"
	baseSepoliaRelayerBaseUrl     = "https://base-sepolia.v2.relayer.renegade.fi"
	baseMainnetBaseUrl            = "https://base-mainnet.v2.auth-server.renegade.fi"
	baseMainnetRelayerBaseUrl     = "https://base-mainnet.v2.relayer.renegade.fi"
	ethereumSepoliaBaseUrl        = "https://ethereum-sepolia.v2.auth-server.renegade.fi"
	ethereumSepoliaRelayerBaseUrl = "https://ethereum-sepolia.v2.relayer.renegade.fi"
	apiKeyHeader                  = "X-Renegade-Api-Key" //nolint:gosec
)

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

// NewArbitrumSepoliaExternalMatchClient creates a new ExternalMatchClient for the Arbitrum Sepolia network
func NewArbitrumSepoliaExternalMatchClient(apiKey string, apiSecret *wallet.HmacKey) *ExternalMatchClient {
	return NewExternalMatchClient(arbitrumSepoliaBaseUrl, arbitrumSepoliaRelayerBaseUrl, apiKey, apiSecret)
}

// NewBaseSepoliaExternalMatchClient creates a new ExternalMatchClient for the Base Sepolia network
func NewBaseSepoliaExternalMatchClient(apiKey string, apiSecret *wallet.HmacKey) *ExternalMatchClient {
	return NewExternalMatchClient(baseSepoliaBaseUrl, baseSepoliaRelayerBaseUrl, apiKey, apiSecret)
}

// NewTestnetExternalMatchClient creates a new ExternalMatchClient for the Arbitrum Sepolia network
//
// Deprecated: Use NewArbitrumSepoliaExternalMatchClient instead
func NewTestnetExternalMatchClient(apiKey string, apiSecret *wallet.HmacKey) *ExternalMatchClient {
	return NewArbitrumSepoliaExternalMatchClient(apiKey, apiSecret)
}

// NewArbitrumOneExternalMatchClient creates a new ExternalMatchClient for the Arbitrum One network
func NewArbitrumOneExternalMatchClient(apiKey string, apiSecret *wallet.HmacKey) *ExternalMatchClient {
	return NewExternalMatchClient(arbitrumOneBaseUrl, arbitrumOneRelayerBaseUrl, apiKey, apiSecret)
}

// NewBaseMainnetExternalMatchClient creates a new ExternalMatchClient for the Base Mainnet network
func NewBaseMainnetExternalMatchClient(apiKey string, apiSecret *wallet.HmacKey) *ExternalMatchClient {
	return NewExternalMatchClient(baseMainnetBaseUrl, baseMainnetRelayerBaseUrl, apiKey, apiSecret)
}

// NewMainnetExternalMatchClient creates a new ExternalMatchClient for the Arbitrum One network
//
// Deprecated: Use NewArbitrumOneExternalMatchClient instead
func NewMainnetExternalMatchClient(apiKey string, apiSecret *wallet.HmacKey) *ExternalMatchClient {
	return NewArbitrumOneExternalMatchClient(apiKey, apiSecret)
}

// NewEthereumSepoliaExternalMatchClient creates a new ExternalMatchClient for the Ethereum Sepolia network
func NewEthereumSepoliaExternalMatchClient(apiKey string, apiSecret *wallet.HmacKey) *ExternalMatchClient {
	return NewExternalMatchClient(ethereumSepoliaBaseUrl, ethereumSepoliaRelayerBaseUrl, apiKey, apiSecret)
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

// ----------------------
// | V2 Market Data API |
// ----------------------

// GetMarkets fetches all tradable markets with their prices and fee rates
func (c *ExternalMatchClient) GetMarkets() (*api_types.GetMarketsResponse, error) {
	var response api_types.GetMarketsResponse
	err := c.relayerHttpClient.GetJSON(
		api_types.GetMarketsPath,
		nil, // body
		&response,
	)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetMarketDepth fetches the market depth for a specific token
func (c *ExternalMatchClient) GetMarketDepth(mint string) (*api_types.GetMarketDepthByMintResponse, error) {
	var response api_types.GetMarketDepthByMintResponse
	path := api_types.BuildGetMarketDepthByMintPath(mint)

	headers := make(http.Header)
	headers.Set(apiKeyHeader, c.apiKey)

	err := c.httpClient.GetWithAuthAndHeaders(path, &headers, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetMarketDepthsAllPairs fetches the market depths for all supported pairs
func (c *ExternalMatchClient) GetMarketDepthsAllPairs() (*api_types.GetMarketDepthsResponse, error) {
	var response api_types.GetMarketDepthsResponse

	headers := make(http.Header)
	headers.Set(apiKeyHeader, c.apiKey)

	err := c.httpClient.GetWithAuthAndHeaders(api_types.GetMarketsDepthPath, &headers, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetExchangeMetadata fetches metadata about the Renegade exchange
func (c *ExternalMatchClient) GetExchangeMetadata() (*api_types.ExchangeMetadataResponse, error) {
	var response api_types.ExchangeMetadataResponse

	headers := make(http.Header)
	headers.Set(apiKeyHeader, c.apiKey)

	err := c.httpClient.GetWithAuthAndHeaders(api_types.GetExchangeMetadataPath, &headers, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// --------------------------
// | V2 Quote + Assembly API |
// --------------------------

// GetExternalMatchQuoteV2 requests a v2 quote from the relayer
// returns nil if no match is found
func (c *ExternalMatchClient) GetExternalMatchQuoteV2(
	order *api_types.ApiExternalOrderV2,
) (*SignedExternalQuoteV2, error) {
	return c.GetExternalMatchQuoteWithOptionsV2(order, NewExternalQuoteOptions())
}

// GetExternalMatchQuoteWithOptionsV2 requests a v2 quote with options
func (c *ExternalMatchClient) GetExternalMatchQuoteWithOptionsV2(
	order *api_types.ApiExternalOrderV2,
	options *ExternalQuoteOptions,
) (*SignedExternalQuoteV2, error) {
	requestBody := api_types.ExternalQuoteRequestV2{
		ExternalOrder: *order,
	}

	var response api_types.ExternalQuoteResponseV2
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

	return NewSignedExternalQuoteV2(&response), nil
}

// AssembleExternalQuoteV2 assembles a v2 quote into a malleable match bundle
// returns nil if no match is found
func (c *ExternalMatchClient) AssembleExternalQuoteV2(
	quote *SignedExternalQuoteV2,
) (*MalleableExternalMatchBundle, error) {
	return c.AssembleExternalQuoteWithOptionsV2(quote, NewAssembleExternalMatchOptionsV2())
}

// AssembleExternalQuoteWithOptionsV2 assembles a v2 quote with options
func (c *ExternalMatchClient) AssembleExternalQuoteWithOptionsV2(
	quote *SignedExternalQuoteV2,
	options *AssembleExternalMatchOptionsV2,
) (*MalleableExternalMatchBundle, error) {
	signedQuote := quote.ToApiSignedQuote()
	assemblyOrder := api_types.NewQuotedOrderAssembly(&signedQuote, options.UpdatedOrder)

	requestBody := api_types.AssembleExternalMatchRequestV2{
		DoGasEstimation: options.DoGasEstimation,
		ReceiverAddress: options.ReceiverAddress,
		Order:           assemblyOrder,
	}

	var response api_types.ExternalMatchResponseV2
	path := api_types.AssembleMatchBundleV2Path
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

	return newMalleableExternalMatchBundle(&response), nil
}

// GetExternalMatchBundleV2 requests a v2 match bundle (direct match)
// returns nil if no match is found
func (c *ExternalMatchClient) GetExternalMatchBundleV2(
	order *api_types.ApiExternalOrderV2,
) (*MalleableExternalMatchBundle, error) {
	return c.GetExternalMatchBundleWithOptionsV2(order, NewExternalMatchOptionsV2())
}

// GetExternalMatchBundleWithOptionsV2 requests a v2 match bundle with options
func (c *ExternalMatchClient) GetExternalMatchBundleWithOptionsV2(
	order *api_types.ApiExternalOrderV2,
	options *ExternalMatchOptionsV2,
) (*MalleableExternalMatchBundle, error) {
	assemblyOrder := api_types.NewDirectOrderAssembly(order)

	requestBody := api_types.AssembleExternalMatchRequestV2{
		DoGasEstimation: options.DoGasEstimation,
		ReceiverAddress: options.ReceiverAddress,
		Order:           assemblyOrder,
	}

	var response api_types.ExternalMatchResponseV2
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

	return newMalleableExternalMatchBundle(&response), nil
}

// ----------------------------
// | V1 Deprecated Shim APIs |
// ----------------------------

// GetSupportedTokens requests the list of supported tokens from the relayer
//
// Deprecated: Use GetMarkets instead, which returns all supported tokens along with their current price
func (c *ExternalMatchClient) GetSupportedTokens() ([]api_types.ApiToken, error) {
	resp, err := c.GetMarkets()
	if err != nil {
		return nil, err
	}

	return marketsToSupportedTokens(resp), nil
}

// GetFeeForAsset requests the fees for a given base token
//
// Deprecated: Use GetMarkets instead
func (c *ExternalMatchClient) GetFeeForAsset(addr *string) (*ExternalMatchFee, error) {
	resp, err := c.GetMarkets()
	if err != nil {
		return nil, err
	}

	return marketsToFeeForAsset(resp, *addr)
}

// GetExternalMatchQuote requests a quote from the relayer (v1 shim)
// returns nil if no match is found
func (c *ExternalMatchClient) GetExternalMatchQuote(
	order *api_types.ApiExternalOrder,
) (*api_types.ApiSignedQuote, error) {
	return c.GetExternalMatchQuoteWithOptions(order, NewExternalQuoteOptions())
}

// GetExternalMatchQuoteWithOptions requests a quote with the given options struct (v1 shim)
func (c *ExternalMatchClient) GetExternalMatchQuoteWithOptions(
	order *api_types.ApiExternalOrder,
	options *ExternalQuoteOptions,
) (*api_types.ApiSignedQuote, error) {
	// Convert v1 order to v2
	v2Order := v1OrderToV2(order)

	// Call v2 method
	v2Quote, err := c.GetExternalMatchQuoteWithOptionsV2(&v2Order, options)
	if err != nil {
		return nil, err
	}
	if v2Quote == nil {
		return nil, nil
	}

	// Convert v2 quote to v1
	return v2QuoteToV1(v2Quote, order)
}

// AssembleExternalQuote generates an external match bundle from a signed quote (v1 shim)
func (c *ExternalMatchClient) AssembleExternalQuote(
	quote *api_types.ApiSignedQuote,
) (*ExternalMatchBundle, error) {
	return c.AssembleExternalQuoteWithReceiver(quote, nil /* receiverAddress */)
}

// AssembleExternalQuoteWithReceiver generates an external match bundle from a signed quote (v1 shim)
// returns nil if no match is found
func (c *ExternalMatchClient) AssembleExternalQuoteWithReceiver(
	quote *api_types.ApiSignedQuote,
	receiverAddress *string,
) (*ExternalMatchBundle, error) {
	options := NewAssembleExternalMatchOptions().WithReceiverAddress(receiverAddress)
	return c.AssembleExternalMatchWithOptions(quote, options)
}

// AssembleExternalMatchWithOptions assembles an external quote with the given options struct (v1 shim)
func (c *ExternalMatchClient) AssembleExternalMatchWithOptions(
	quote *api_types.ApiSignedQuote,
	options *AssembleExternalMatchOptions,
) (*ExternalMatchBundle, error) {
	direction := quote.Quote.Order.Side

	// Extract the v2 quote from the v1 quote
	v2Quote, err := v1QuoteToV2(quote)
	if err != nil {
		return nil, err
	}

	// Convert v1 options to v2
	v2Options := v1AssembleOptionsToV2(options, &quote.Quote.Order)

	// Call v2 method
	v2Resp, err := c.AssembleExternalQuoteWithOptionsV2(v2Quote, v2Options)
	if err != nil {
		return nil, err
	}
	if v2Resp == nil {
		return nil, nil
	}

	// Convert v2 response to v1
	return v2MalleableBundleToV1(v2Resp, direction)
}

// GetExternalMatchBundle requests an external match bundle from the relayer (v1 shim)
// returns nil if no match is found
func (c *ExternalMatchClient) GetExternalMatchBundle(
	request *api_types.ApiExternalOrder,
) (*ExternalMatchBundle, error) {
	return c.GetExternalMatchBundleWithReceiver(request, nil /* receiverAddress */)
}

// GetExternalMatchBundleWithReceiver requests an external match bundle from the relayer (v1 shim)
// returns nil if no match is found
func (c *ExternalMatchClient) GetExternalMatchBundleWithReceiver(
	request *api_types.ApiExternalOrder,
	receiverAddress *string,
) (*ExternalMatchBundle, error) {
	options := &ExternalMatchOptions{
		AssembleExternalMatchOptions: AssembleExternalMatchOptions{
			ReceiverAddress:       receiverAddress,
			RequestGasSponsorship: true, // default to true
		},
	}
	return c.GetExternalMatchBundleWithOptions(request, options)
}

// GetExternalMatchBundleWithOptions requests an external match bundle from the relayer with the given options (v1 shim)
// returns nil if no match is found
func (c *ExternalMatchClient) GetExternalMatchBundleWithOptions(
	request *api_types.ApiExternalOrder,
	options *ExternalMatchOptions,
) (*ExternalMatchBundle, error) {
	direction := request.Side

	// Convert v1 order to v2
	v2Order := v1OrderToV2(request)

	// Build v2 options
	v2Options := &ExternalMatchOptionsV2{
		DoGasEstimation:       options.DoGasEstimation,
		ReceiverAddress:       options.ReceiverAddress,
		DisableGasSponsorship: !options.RequestGasSponsorship,
		GasRefundAddress:      options.GasRefundAddress,
	}

	// Call v2 method
	v2Resp, err := c.GetExternalMatchBundleWithOptionsV2(&v2Order, v2Options)
	if err != nil {
		return nil, err
	}
	if v2Resp == nil {
		return nil, nil
	}

	// Convert v2 malleable bundle to v1
	return v2MalleableBundleToV1(v2Resp, direction)
}

// v2MalleableBundleToV1 converts a v2 MalleableExternalMatchBundle to a v1 ExternalMatchBundle.
// This reconstructs the v2 API response from the parsed bundle and delegates to v2ResponseToV1NonMalleable.
func v2MalleableBundleToV1(
	bundle *MalleableExternalMatchBundle,
	direction string,
) (*ExternalMatchBundle, error) {
	// Reconstruct an ApiSettlementTransactionV2 from the parsed SettlementTransaction
	apiSettlementTx := toApiSettlementTransactionV2(bundle.SettlementTx)

	// Build the v2 API response to pass to the conversion function
	v2Resp := &api_types.ExternalMatchResponseV2{
		MatchBundle: api_types.MalleableAtomicMatchApiBundleV2{
			MatchResult:  *bundle.MatchResult,
			FeeRates:     *bundle.FeeRates,
			MaxReceive:   *bundle.MaxReceive,
			MinReceive:   *bundle.MinReceive,
			MaxSend:      *bundle.MaxSend,
			MinSend:      *bundle.MinSend,
			SettlementTx: apiSettlementTx,
			Deadline:     bundle.Deadline,
		},
		GasSponsorshipInfo: bundle.GasSponsorshipInfo,
	}

	return v2ResponseToV1NonMalleable(v2Resp, direction)
}

// ------------------
// | Request Helper |
// ------------------

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
