package external_match_client

import (
	"net/http"

	geth_common "github.com/ethereum/go-ethereum/common"
	"github.com/renegade-fi/golang-sdk/client"
	"github.com/renegade-fi/golang-sdk/client/api_types"
	"github.com/renegade-fi/golang-sdk/wallet"
)

const (
	testnetBaseUrl = "https://testnet.auth-server.renegade.fi:3000"
	mainnetBaseUrl = "https://mainnet.auth-server.renegade.fi:3000"
	apiKeyHeader   = "X-Renegade-Api-Key"
)

// ExternalMatchBundle is the application level analog to the ApiExternalMatchBundle
type ExternalMatchBundle struct {
	MatchResult  *api_types.ApiExternalMatchResult
	SettlementTx *SettlementTransaction
}

// SettlementTransaction is the application level analog to the ApiSettlementTransaction
type SettlementTransaction struct {
	Type string
	To   geth_common.Address
	Data []byte
}

// toSettlementTransaction converts an ApiSettlementTransaction to a SettlementTransaction
func toSettlementTransaction(tx *api_types.ApiSettlementTransaction) *SettlementTransaction {
	// Parse a geth address and bytes data from hex strings
	to := geth_common.HexToAddress(tx.To)
	data := geth_common.FromHex(tx.Data)

	return &SettlementTransaction{
		Type: tx.Type,
		To:   to,
		Data: data,
	}
}

// ExternalMatchClient represents a client for the external match API
//
// This client can be used to request external match bundles from a relayer.
// The relayer will return a match and a transaction to submit on-chain
type ExternalMatchClient struct {
	apiKey     string
	httpClient *client.HttpClient
}

// NewTestnetExternalMatchClient creates a new ExternalMatchClient for the testnet
func NewTestnetExternalMatchClient(apiKey string, apiSecret *wallet.HmacKey) *ExternalMatchClient {
	return NewExternalMatchClient(testnetBaseUrl, apiKey, apiSecret)
}

// NewMainnetExternalMatchClient creates a new ExternalMatchClient for the mainnet
func NewMainnetExternalMatchClient(apiKey string, apiSecret *wallet.HmacKey) *ExternalMatchClient {
	return NewExternalMatchClient(mainnetBaseUrl, apiKey, apiSecret)
}

// NewExternalMatchClient creates a new ExternalMatchClient with the given base URL, api key, and api secret
func NewExternalMatchClient(baseURL string, apiKey string, apiSecret *wallet.HmacKey) *ExternalMatchClient {
	return &ExternalMatchClient{
		apiKey:     apiKey,
		httpClient: client.NewHttpClient(baseURL, apiSecret),
	}
}

// GetExternalMatchBundle requests an external match bundle from the relayer
func (c *ExternalMatchClient) GetExternalMatchBundle(request *api_types.ApiExternalOrder) (*ExternalMatchBundle, error) {
	// Construct a request
	requestBody := api_types.ExternalMatchRequest{
		ExternalOrder: *request,
	}

	path := api_types.GetExternalMatchBundlePath
	headers := make(http.Header)
	headers.Set(apiKeyHeader, c.apiKey)

	// Send the request
	response := api_types.ExternalMatchResponse{}
	if err := c.httpClient.PostWithAuthAndHeaders(path, &headers, requestBody, &response); err != nil {
		return nil, err
	}

	// Convert into the application level type
	return &ExternalMatchBundle{
		MatchResult:  &response.Bundle.MatchResult,
		SettlementTx: toSettlementTransaction(&response.Bundle.SettlementTx),
	}, nil
}
