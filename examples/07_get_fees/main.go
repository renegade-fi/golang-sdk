// This example demonstrates how to get fees for a given asset
package main

import (
	"fmt"

	"github.com/renegade-fi/golang-sdk/examples/common"
)

func main() {
	client, err := common.CreateExternalMatchClient()
	if err != nil {
		panic(err)
	}

	// Get WETH address
	wethAddr, err := common.FindTokenAddr("WETH", client)
	if err != nil {
		panic(err)
	}

	// Get fees for WETH
	fees, err := client.GetFeeForAsset(&wethAddr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Fees for WETH:\n")
	fmt.Printf("Relayer Fee: %v\n", fees.RelayerFee)
	fmt.Printf("Protocol Fee: %v\n", fees.ProtocolFee)
	fmt.Printf("Total Fee: %v\n", fees.Total())
}
