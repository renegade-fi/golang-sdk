package main

import (
	"fmt"
	"math/big"

	"github.com/renegade-fi/golang-sdk/client/api_types"
	external_match_client "github.com/renegade-fi/golang-sdk/client/external_match_client"
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
		WithQuoteAmount(api_types.Amount(*quoteAmount)).
		WithSide("Buy").
		WithMinFillSize(api_types.Amount(*minFillSize)).
		Build()
	if err != nil {
		panic(err)
	}

	if err := getQuoteAndSubmitWithReceiver(order, client); err != nil {
		panic(err)
	}
}

// getQuoteAndSubmitWithReceiver gets a quote, assembles it with a separate receiver, then submits
func getQuoteAndSubmitWithReceiver(order *api_types.ApiExternalOrder, client *external_match_client.ExternalMatchClient) error {
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

	// 2. Assemble the bundle with a modified order
	fmt.Println("Assembling bundle...")
	newOrder := order
	newOrder.QuoteAmount = api_types.NewAmount(19_000_000) // Modify to 19 USDC
	receiverAddress := "0xC5fE800A3D92112473e4E811296F194DA7b26BA7"
	options := external_match_client.NewAssembleExternalMatchOptions().
		WithReceiverAddress(&receiverAddress).
		WithUpdatedOrder(newOrder)

	bundle, err := client.AssembleExternalMatchWithOptions(quote, options)
	if err != nil {
		return err
	}

	if bundle == nil {
		fmt.Println("No bundle found")
		return nil
	}

	// 3. Submit the bundle
	fmt.Println("Submitting bundle...")
	if err := common.SubmitBundle(*bundle); err != nil {
		return fmt.Errorf("failed to submit bundle: %w", err)
	}

	fmt.Print("Bundle submitted successfully!\n\n")
	return nil
}
