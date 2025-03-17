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

	// Create order for 20 USDC worth of WETH
	quoteAmount := new(big.Int).SetUint64(20_000_000) // $20 USDC
	minFillSize := big.NewInt(0)
	order, err := api_types.NewExternalOrderBuilder().
		WithQuoteMint(quoteMint).
		WithBaseMint(baseMint).
		WithQuoteAmount(api_types.Amount(*quoteAmount)).
		WithSide("Sell").
		WithMinFillSize(api_types.Amount(*minFillSize)).
		Build()
	if err != nil {
		panic(err)
	}

	if err := getQuoteAndSubmitWithGasSponsorship(order, client); err != nil {
		panic(err)
	}
}

// getQuoteAndSubmitWithGasSponsorship gets a quote with gas sponsorship, assembles it, then submits
func getQuoteAndSubmitWithGasSponsorship(
	order *api_types.ApiExternalOrder,
	client *external_match_client.ExternalMatchClient,
) error {
	// 1. Get a quote from the relayer, explicitly requesting in-kind gas sponsorship.
	fmt.Println("Getting quote with gas sponsorship...")
	// When you leave the `refundAddress` unset, the in-kind refund is directed to the receiver address.
	// This is equivalent to the trade itself having a better price, so the price in the quote will be
	// updated to reflect this
	quoteOptions := external_match_client.NewExternalQuoteOptions().
		WithDisableGasSponsorship(false). // Note that this is false by default
		WithRefundNativeEth(false)        // Note that this is false by default

	quote, err := client.GetExternalMatchQuoteWithOptions(order, quoteOptions)
	if err != nil {
		return err
	}

	if quote == nil {
		fmt.Println("No quote found")
		return nil
	}

	// 2. Assemble the bundle
	fmt.Println("Assembling bundle...")
	receiverAddress := "0xC5fE800A3D92112473e4E811296F194DA7b26BA7"
	assemblyOptions := external_match_client.NewAssembleExternalMatchOptions().
		WithReceiverAddress(&receiverAddress)

	bundle, err := client.AssembleExternalMatchWithOptions(quote, assemblyOptions)
	if err != nil {
		return err
	}

	if bundle == nil {
		fmt.Println("No bundle found")
		return nil
	}

	if !bundle.GasSponsored {
		fmt.Println("Bundle was not sponsored, abandoning...")
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
