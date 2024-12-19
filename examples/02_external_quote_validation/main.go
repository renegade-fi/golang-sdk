// Package main is an example of how to use the Renegade SDK to get an external
// quote, validate it, and submit it to the sequencer.
package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/renegade-fi/golang-sdk/client/api_types"
	external_match_client "github.com/renegade-fi/golang-sdk/client/external_match_client"
	"github.com/renegade-fi/golang-sdk/wallet"
)

const (
	quoteMint       = "0xdf8d259c04020562717557f2b5a3cf28e92707d1" // USDC
	baseMint        = "0xc3414a7ef14aaaa9c4522dfc00a4e66e74e9c25a" // WETH
	darkpoolAddress = "0x9af58f1ff20ab22e819e40b57ffd784d115a9ef5"
	chainID         = 421614 // Testnet
)

func main() {
	// ... Token Approvals to Darkpool ... //

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

	externalMatchClient := external_match_client.NewTestnetExternalMatchClient(apiKey, &apiSecretKey)

	// Request an external match
	// We can denominate the order size in either the quote or base token with
	// `WithQuoteAmount` or `WithBaseAmount` respectively.
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

	if err := getQuoteAndSubmit(order, externalMatchClient); err != nil {
		panic(err)
	}
}

// getQuoteAndSubmit gets a quote, assembled is, then submits the bundle
func getQuoteAndSubmit(order *api_types.ApiExternalOrder, client *external_match_client.ExternalMatchClient) error {
	// 1. Get a quote from the relayer
	fmt.Println("Getting quote...")
	signedQuote, err := client.GetExternalMatchQuote(order)
	if err != nil {
		return err
	}

	if signedQuote == nil {
		fmt.Println("No quote found")
		return nil
	}

	if !validateQuote(&signedQuote.Quote) {
		fmt.Println("Quote is not acceptable")
		return nil
	}

	// 2. Assemble the bundle
	fmt.Println("Assembling bundle...")
	bundle, err := client.AssembleExternalQuote(signedQuote)
	if err != nil {
		return err
	}

	if bundle == nil {
		fmt.Println("No bundle found")
		return nil
	}

	// 3. Submit the bundle
	fmt.Println("Submitting bundle...")
	if err := submitBundle(*bundle); err != nil {
		return err
	}

	fmt.Println("Bundle submitted successfully!")
	return nil
}

// validateQuote validates a quote before submitting it
func validateQuote(quote *api_types.ApiExternalQuote) bool {
	minFillSize := api_types.NewAmount(1000000000000000) // 0.001 WETH
	maxFees := api_types.NewAmount(10000000000000)       // 0.0001 WETH

	recv := quote.Receive.Amount
	fees := quote.Fees.Total()

	if recv.Cmp(minFillSize) < 0 {
		fmt.Printf("Quote fill size is less than minimum fill size (%s < %s)\n", recv.String(), minFillSize.String())
		return false
	}

	if fees.Cmp(maxFees) > 0 {
		fmt.Printf("Quote fees are greater than the maximum allowed fees (%s > %s)\n", fees.String(), maxFees.String())
		return false
	}

	return true
}

// submitBundle submits the bundle to the sequencer
func submitBundle(bundle external_match_client.ExternalMatchBundle) error {
	// Initialize eth client
	ethClient, err := getEthClient()
	if err != nil {
		panic(err)
	}

	privateKey, err := getPrivateKey()
	if err != nil {
		panic(err)
	}

	// Send the transaction to the sequencer
	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}

	nonce, err := ethClient.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(privateKey.PublicKey))
	if err != nil {
		panic(err)
	}

	ethTx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   big.NewInt(chainID), // Sepolia chain ID
		Nonce:     nonce,
		GasTipCap: gasPrice,                                  // Use suggested gas price as tip cap
		GasFeeCap: new(big.Int).Mul(gasPrice, big.NewInt(2)), // Fee cap at 2x gas price
		Gas:       uint64(10_000_000),                        // Gas limit
		To:        &bundle.SettlementTx.To,                   // Contract address
		Value:     bundle.SettlementTx.Value,                 // No ETH transfer
		Data:      []byte(bundle.SettlementTx.Data),          // Contract call data
	})

	// Sign and send transaction
	signer := types.LatestSignerForChainID(big.NewInt(chainID))
	signedTx, err := types.SignTx(ethTx, signer, privateKey)
	if err != nil {
		panic(err)
	}

	err = ethClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Transaction submitted! Hash: %s\n", signedTx.Hash().Hex())
	return nil
}

// -----------
// | Helpers |
// -----------

func getRpcURL() string { //nolint:revive
	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		panic("RPC_URL environment variable not set")
	}
	return rpcURL
}

func getEthClient() (*ethclient.Client, error) {
	return ethclient.Dial(getRpcURL())
}

func getPrivateKey() (*ecdsa.PrivateKey, error) {
	privKeyHex := os.Getenv("PKEY")
	if privKeyHex == "" {
		return nil, fmt.Errorf("PKEY environment variable not set")
	}

	return crypto.HexToECDSA(privKeyHex)
}
