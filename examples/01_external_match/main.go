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
	darkpoolAddress = "0x9af58f1ff20ab22e819e40b57ffd784d115a9ef5"
	chainId         = 421614 // Testnet
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

	// You can fetch token mappings from the relayer using the client
	quoteMint, err := findTokenAddr("USDC", externalMatchClient)
	if err != nil {
		panic(err)
	}
	baseMint, err := findTokenAddr("WETH", externalMatchClient)
	if err != nil {
		panic(err)
	}

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
	quote, err := client.GetExternalMatchQuote(order)
	if err != nil {
		return err
	}

	if quote == nil {
		fmt.Println("No quote found")
		return nil
	}

	// ... Check if the quote is acceptable ... //

	// 2. Assemble the bundle
	fmt.Println("Assembling bundle...")
	bundle, err := client.AssembleExternalQuote(quote)
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

	fmt.Print("Bundle submitted successfully!\n\n")
	return nil
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
		ChainID:   big.NewInt(chainId), // Sepolia chain ID
		Nonce:     nonce,
		GasTipCap: gasPrice,                                  // Use suggested gas price as tip cap
		GasFeeCap: new(big.Int).Mul(gasPrice, big.NewInt(2)), // Fee cap at 2x gas price
		Gas:       uint64(10_000_000),                        // Gas limit
		To:        &bundle.SettlementTx.To,                   // Contract address
		Value:     bundle.SettlementTx.Value,                 // No ETH transfer
		Data:      []byte(bundle.SettlementTx.Data),          // Contract call data
	})

	// Sign and send transaction
	signer := types.LatestSignerForChainID(big.NewInt(chainId))
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

func getRpcUrl() string {
	rpcUrl := os.Getenv("RPC_URL")
	if rpcUrl == "" {
		panic("RPC_URL environment variable not set")
	}
	return rpcUrl
}

func getEthClient() (*ethclient.Client, error) {
	return ethclient.Dial(getRpcUrl())
}

func getPrivateKey() (*ecdsa.PrivateKey, error) {
	privKeyHex := os.Getenv("PKEY")
	if privKeyHex == "" {
		return nil, fmt.Errorf("PKEY environment variable not set")
	}

	return crypto.HexToECDSA(privKeyHex)
}

// findTokenAddr fetches the address of a token from the relayer
func findTokenAddr(symbol string, client *external_match_client.ExternalMatchClient) (string, error) {
	// Fetch the list of supported tokens from the relayer
	tokens, err := client.GetSupportedTokens()
	if err != nil {
		return "", err
	}

	// Find the token with the matching symbol
	for _, token := range tokens {
		if token.Symbol == symbol {
			return token.Address, nil
		}
	}

	return "", fmt.Errorf("token not found")
}
