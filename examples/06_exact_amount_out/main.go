// Package main provides an example of how to get a quote for an exact amount out
package main

import (
	"fmt"
	"math/big"

	"github.com/renegade-fi/golang-sdk/client/api_types"
	external_match_client "github.com/renegade-fi/golang-sdk/client/external_match_client"
	"github.com/renegade-fi/golang-sdk/examples/common"
)

func main() {
	client, err := common.CreateExternalMatchClient()
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

	// Specify an exact output amount of the quote token
	quoteAmountOut := new(big.Int).SetUint64(2_123_456)
	order, err := api_types.NewExternalOrderBuilder().
		WithQuoteMint(quoteMint).
		WithBaseMint(baseMint).
		WithExactQuoteAmountOutput(api_types.Amount(*quoteAmountOut)).
		WithSide("Sell").
		Build()
	if err != nil {
		panic(err)
	}

	if err := getQuoteWithExactAmount(order, client); err != nil {
		panic(err)
	}
}

// getQuoteWithExactAmount gets a quote and prints the details
func getQuoteWithExactAmount(order *api_types.ApiExternalOrder, client *external_match_client.ExternalMatchClient) error {
	// 1. Get a quote from the relayer
	fmt.Println("Getting quote...")
	quote, err := client.GetExternalMatchQuote(order)
	if err != nil {
		return err
	}

	if quote == nil {
		fmt.Println("No quote found")
		return nil
	}

	// Print the quote details
	fmt.Printf("Quote found!\n")
	fmt.Printf("You will send: %v %s\n", quote.Quote.Send.Amount, quote.Quote.Send.Mint)
	fmt.Printf("You will receive (net of fees): %v %s\n", quote.Quote.Receive.Amount, quote.Quote.Receive.Mint)
	fmt.Printf("Total fees: %v\n", quote.Quote.Fees.Total())

	// You can now assemble the quote and submit a bundle, see `01_external_match` for an example
	return nil
}

// -----------
// | Helpers |
// -----------

func findTokenAddr(symbol string, client *external_match_client.ExternalMatchClient) (string, error) {
	tokens, err := client.GetSupportedTokens()
	if err != nil {
		return "", err
	}

	for _, token := range tokens {
		if token.Symbol == symbol {
			return token.Address, nil
		}
	}

	return "", fmt.Errorf("token not found")
}
