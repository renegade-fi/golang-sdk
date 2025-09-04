// Package main demonstrates retrieving and submitting an external match bundle directly.
package main

import (
	"fmt"
	"math/big"

	"github.com/renegade-fi/golang-sdk/client/api_types"
	"github.com/renegade-fi/golang-sdk/client/external_match_client"
	"github.com/renegade-fi/golang-sdk/examples/common"
)

func main() {
	client, err := common.CreateArbitrumExternalMatchClient()
	if err != nil {
		panic(err)
	}

	// Fetch token mappings from the relayer
	quoteMint, err := common.FindTokenAddr("USDC", client)
	if err != nil {
		panic(err)
	}
	baseMint, err := common.FindTokenAddr("WETH", client)
	if err != nil {
		panic(err)
	}

	// Create order for 20 USDC worth of WETH
	quoteAmount := new(big.Int).SetUint64(20_000_000) // $20 USDC
	minFillSize := big.NewInt(0)
	order, err := api_types.NewExternalOrderBuilder().
		WithQuoteMint(quoteMint).
		WithBaseMint(baseMint).
		WithQuoteAmount(quoteAmount).
		WithSide("Buy").
		WithMinFillSize(minFillSize).
		Build()
	if err != nil {
		panic(err)
	}

	if err := getMatchAndSubmit(order, client); err != nil {
		panic(err)
	}
}

func getMatchAndSubmit(order *api_types.ApiExternalOrder, client *external_match_client.ExternalMatchClient) error {
	// 1. Rather than quote + assemble as in other examples, we directly request
	// 	a bundle from the relayer
	fmt.Println("Fetching bundle...")
	bundle, err := client.GetExternalMatchBundle(order)
	if err != nil {
		return err
	}

	if bundle == nil {
		fmt.Println("No bundle found")
		return nil
	}

	// 2. Submit the bundle
	fmt.Println("Submitting bundle...")
	if err := common.SubmitBundle(bundle); err != nil {
		return fmt.Errorf("failed to submit bundle: %w", err)
	}

	fmt.Print("Bundle submitted successfully!\n\n")
	return nil
}
