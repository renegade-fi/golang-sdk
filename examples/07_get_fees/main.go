// This example demonstrates how to get fees for a given asset
package main

import (
	"fmt"
	"os"

	external_match_client "github.com/renegade-fi/golang-sdk/client/external_match_client"
	"github.com/renegade-fi/golang-sdk/wallet"
)

func main() {
	// Get API credentials from environment
	apiKey := os.Getenv("EXTERNAL_MATCH_KEY")
	apiSecret := os.Getenv("EXTERNAL_MATCH_SECRET")
	if apiKey == "" || apiSecret == "" {
		panic("EXTERNAL_MATCH_KEY and EXTERNAL_MATCH_SECRET must be set")
	}

	apiSecretKey, err := new(wallet.HmacKey).FromBase64String(apiSecret)
	if err != nil {
		panic(err)
	}

	// Get fees for WETH
	externalMatchClient := external_match_client.NewTestnetExternalMatchClient(apiKey, &apiSecretKey)

	mint := "0xc3414a7ef14aaaa9c4522dfc00a4e66e74e9c25a" // Testnet WETH
	fees, err := externalMatchClient.GetFeeForAsset(&mint)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Fees for WETH:\n")
	fmt.Printf("Relayer Fee: %v\n", fees.RelayerFee)
	fmt.Printf("Protocol Fee: %v\n", fees.ProtocolFee)
	fmt.Printf("Total Fee: %v\n", fees.Total())
}
