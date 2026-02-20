// Package common contains common functions for the examples
package common

import (
	"fmt"
	"os"

	external_match_client "github.com/renegade-fi/golang-sdk/client/external_match_client"
	"github.com/renegade-fi/golang-sdk/wallet"
)

// CreateArbitrumExternalMatchClient creates a new external match client using environment variables
func CreateArbitrumExternalMatchClient() (*external_match_client.ExternalMatchClient, error) {
	apiKey := os.Getenv("EXTERNAL_MATCH_KEY")
	apiSecret := os.Getenv("EXTERNAL_MATCH_SECRET")
	if apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("EXTERNAL_MATCH_KEY and EXTERNAL_MATCH_SECRET must be set")
	}

	apiSecretKey, err := new(wallet.HmacKey).FromBase64String(apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API secret: %w", err)
	}

	return external_match_client.NewArbitrumSepoliaExternalMatchClient(apiKey, &apiSecretKey), nil
}

// CreateBaseExternalMatchClient creates a new external match client for the Base network
func CreateBaseExternalMatchClient() (*external_match_client.ExternalMatchClient, error) {
	apiKey := os.Getenv("EXTERNAL_MATCH_KEY")
	apiSecret := os.Getenv("EXTERNAL_MATCH_SECRET")
	if apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("EXTERNAL_MATCH_KEY and EXTERNAL_MATCH_SECRET must be set")
	}

	apiSecretKey, err := new(wallet.HmacKey).FromBase64String(apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API secret: %w", err)
	}

	return external_match_client.NewBaseSepoliaExternalMatchClient(apiKey, &apiSecretKey), nil
}

// testnetTokenAddresses contains fallback token addresses for Arbitrum Sepolia
var testnetTokenAddresses = map[string]string{
	"USDC": "0xdf8d259c04020562717557f2b5a3cf28e92707d1",
	"WETH": "0xc3414a7ef14aaaa9c4522dfc00a4e66e74e9c25a",
}

// FindTokenAddr fetches the address of a token from the relayer,
// falling back to hardcoded testnet addresses if the API is unavailable
func FindTokenAddr(symbol string, client *external_match_client.ExternalMatchClient) (string, error) {
	tokens, err := client.GetSupportedTokens()
	if err != nil {
		// Fallback to hardcoded testnet addresses
		if addr, ok := testnetTokenAddresses[symbol]; ok {
			fmt.Printf("Warning: GetSupportedTokens failed (%v), using hardcoded address for %s\n", err, symbol)
			return addr, nil
		}
		return "", err
	}

	for _, token := range tokens {
		if token.Symbol == symbol {
			return token.Address, nil
		}
	}

	return "", fmt.Errorf("token %s not found", symbol)
}
