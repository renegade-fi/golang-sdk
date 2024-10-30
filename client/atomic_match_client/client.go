package atomic_match_client

import (
	"net/http"

	"github.com/renegade-fi/golang-sdk/client"
	"github.com/renegade-fi/golang-sdk/client/api_types"
	"github.com/renegade-fi/golang-sdk/wallet"
)

const (
	testnetBaseUrl = "https://testnet.auth-server.renegade.fi:3000"
	mainnetBaseUrl = "https://mainnet.auth-server.renegade.fi:3000"
	apiKeyHeader   = "X-Renegade-Api-Key"
)

// AtomicMatchClient represents a client for the atomic match API
//
// This client can be used to request atomic match bundles from a relayer.
// The relayer will return a match and a transaction to submit on-chain
type AtomicMatchClient struct {
	apiKey     string
	httpClient *client.HttpClient
}

// NewTestnetAtomicMatchClient creates a new AtomicMatchClient for the testnet
func NewTestnetAtomicMatchClient(apiKey string, apiSecret *wallet.HmacKey) *AtomicMatchClient {
	return NewAtomicMatchClient(testnetBaseUrl, apiKey, apiSecret)
}

// NewMainnetAtomicMatchClient creates a new AtomicMatchClient for the mainnet
func NewMainnetAtomicMatchClient(apiKey string, apiSecret *wallet.HmacKey) *AtomicMatchClient {
	return NewAtomicMatchClient(mainnetBaseUrl, apiKey, apiSecret)
}

// NewAtomicMatchClient creates a new AtomicMatchClient with the given base URL, api key, and api secret
func NewAtomicMatchClient(baseURL string, apiKey string, apiSecret *wallet.HmacKey) *AtomicMatchClient {
	return &AtomicMatchClient{
		apiKey:     apiKey,
		httpClient: client.NewHttpClient(baseURL, apiSecret),
	}
}

// GetAtomicMatchBundle requests an atomic match bundle from the relayer
func (c *AtomicMatchClient) GetAtomicMatchBundle(request *api_types.ApiExternalOrder) (*api_types.AtomicMatchBundle, error) {
	// Construct a request
	requestBody := api_types.ExternalMatchRequest{
		ExternalOrder: *request,
	}

	path := api_types.GetAtomicMatchBundlePath
	headers := make(http.Header)
	headers.Set(apiKeyHeader, c.apiKey)

	// Send the request
	response := api_types.ExternalMatchResponse{}
	if err := c.httpClient.PostWithAuthAndHeaders(path, &headers, requestBody, &response); err != nil {
		return nil, err
	}

	return &response.Bundle, nil
}
