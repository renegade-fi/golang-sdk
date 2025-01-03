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

	// Fetch token mappings from the relayer
	quoteMint, err := findTokenAddr("USDC", externalMatchClient)
	if err != nil {
		panic(err)
	}
	baseMint, err := findTokenAddr("WETH", externalMatchClient)
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

	if err := getQuoteAndSubmitWithReceiver(order, externalMatchClient); err != nil {
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

	// 2. Assemble the bundle with a separate receiver address
	fmt.Println("Assembling bundle...")
	newOrder := order
	newOrder.QuoteAmount = api_types.NewAmount(19_000_000)
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
	if err := submitBundle(*bundle); err != nil {
		return err
	}

	fmt.Print("Bundle submitted successfully!\n\n")
	return nil
}

// submitBundle submits the bundle to the sequencer
func submitBundle(bundle external_match_client.ExternalMatchBundle) error {
	ethClient, err := getEthClient()
	if err != nil {
		panic(err)
	}

	privateKey, err := getPrivateKey()
	if err != nil {
		panic(err)
	}

	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}

	nonce, err := ethClient.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(privateKey.PublicKey))
	if err != nil {
		panic(err)
	}

	ethTx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   big.NewInt(chainId),
		Nonce:     nonce,
		GasTipCap: gasPrice,
		GasFeeCap: new(big.Int).Mul(gasPrice, big.NewInt(2)),
		Gas:       uint64(10_000_000),
		To:        &bundle.SettlementTx.To,
		Value:     bundle.SettlementTx.Value,
		Data:      []byte(bundle.SettlementTx.Data),
	})

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
