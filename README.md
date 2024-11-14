# golang-sdk
The Golang sdk for building Renegade clients

## Basic Use
### Building a Client
The majority of the client's goes through the `RenegadeClient` type. The client encapsulates the environment specific config information. You can create a client as follows:
```go
package renegade_demo

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	renegade_client "github.com/renegade-fi/golang-sdk/client/renegade_client"
	renegade_wallet "github.com/renegade-fi/golang-sdk/wallet"
)

const (
    // baseUrl is the location at which to dial the relayer
    baseUrl = "https://testnet.cluster0.renegade.fi:3000"
)

func main() {
	// Practically speaking you will bring your own Arbitrum private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new renegade client
	client, err := renegade_client.NewSepoliaRenegadeClient(baseUrl, privateKey)
	if err != nil {
		log.Fatal(err)
	}
}
```
Behind the scenes, this method deterministically derives a Renegade wallet from your Ethereum keypair. The client will refer to this derived wallet for all further operations.

### Creating a wallet
For first time users of the protocol, the first step is to create a _new_ Renegade wallet. This is as simple as:
```go
wallet, err := client.CreateWallet()
```

### Looking up A Wallet
When reconnecting to a relayer after some time, it is worth checking that the relayer has indexed your wallet from its on-chain storage. This can be done as:
```go
wallet, err := client.CheckWallet()
```
This method will check for the configured wallet in the relayer's state. If not found, the client will instruct the relayer to find the wallet on-chain.

### Deposit Funds
Suppose we want to sell Bitcoin (wBTC) in the darkpool. The first step is to deposit from your configured arbitrum address:
```go
wbtcMint := "0xa91d929ea161688448f61cb3865a6d948d8bd904"
amount := big.NewInt(1000000)  // 10^6
wallet, err = client.Deposit(wbtcMint, amount, privateKey)
```
Note that the amount field is _not_ decimal adjusted; for wBTC -- which has 8 decimals on mainnet -- this translates to 0.01 wBTC. Tokens and their mint addresses that renegade supports can be found at the following locations:
- [Arbitrum Sepolia](https://github.com/renegade-fi/token-mappings/blob/main/testnet.json)
- [Arbitrum One Mainnet](https://github.com/renegade-fi/token-mappings/blob/main/mainnet.json)

**Note:** It is not required that a wallet contain a balance that capitalizes each of their open orders; open orders without a balance backing them will simply not be matched. Therefore, this step is not strictly a prerequisite to the following step in which we place an order.

### Place an Order
Assuming we wish to sell the wBTC that we deposited in the previous step:
```go
btcMint := "0xa91d929ea161688448f61cb3865a6d948d8bd904"
usdcMint := "0x404b26cd9055b35581c68ba9a2b878cca971b0a7"
amount, _ := wallet.GetBalance(btcMint) // Sell the whole balance
order := renegade.NewOrderBuilder().
    WithBaseMintHex(baseMint).
    WithQuoteMintHex(quoteMint).
    WithAmountBigInt(amount).
    WithSide(renegade.OrderSide_SELL).
    Build()

wallet, err = client.PlaceOrder(&order)
```
**Note:** For the moment, all pairs are USDC quoted. E.g. Renegade does not currently support selling wBTC/wETH.

Once the order is placed with a balance to capitalize it, the matching engine will match the order with counter-flow that it finds. 

### Pay Fees and Withdraw
Suppose the order above matched, and your wallet now holds roughly 600000000 USDC ($600 decimal adjusted), which you wish to withdraw back to your Arbitrum wallet. The first step is to pay fees. The Renegade protocol requires that all relayer and protocol fees are paid out before any balance is withdrawn. 

The following snippet pays fees for the wallet then withdraws the entire USDC balance:
```go
wallet, err = client.PayFees()
if err != nil { log.Fatal(err) }

usdcBalance, err := wallet.GetBalance(usdcMint)
wallet, err = client.Withdraw(usdcMint, usdcBalance)
```

### Putting it Together
```go
package test

import (
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	renegade_client "github.com/renegade-fi/golang-sdk/client/renegade_client"
	renegade_wallet "github.com/renegade-fi/golang-sdk/wallet"
)

const (
	// baseUrl is the location at which to dial the relayer
	baseUrl = "https://testnet.cluster0.renegade.fi:3000"
)

func main() {
	// Practically speaking you will bring your own Arbitrum private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new renegade client
	client, err := renegade_client.NewSepoliaRenegadeClient(baseUrl, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Lookup your Renegade wallet (you should create one if not already done)
	wallet, err := client.CheckWallet()
	if err != nil {
		log.Fatal(err)
	}

	// Deposit 0.01 wBTC
	wbtcMint := "0xa91d929ea161688448f61cb3865a6d948d8bd904"
	amount := big.NewInt(1000000) // 2^16
	wallet, err = client.Deposit(wbtcMint, amount, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Sell 0.01 wBTC
	usdcMint := "0x404b26cd9055b35581c68ba9a2b878cca971b0a7"
	amount, _ = wallet.GetBalance(wbtcMint)
	order := renegade_wallet.NewOrderBuilder().
		WithBaseMintHex(wbtcMint).
		WithQuoteMintHex(usdcMint).
		WithAmountBigInt(amount).
		WithSide(renegade_wallet.OrderSide_SELL).
		Build()

	wallet, err = client.PlaceOrder(&order)
	if err != nil {
		log.Fatal(err)
	}

	// ... Matching Engine Matches Order ... //

	// Pay fees and withdraw
	wallet, err = client.PayFees()
	if err != nil {
		log.Fatal(err)
	}

	usdcBalance, _ := wallet.GetBalance(usdcMint)
	wallet, err = client.Withdraw(usdcMint, usdcBalance)
	if err != nil {
		log.Fatal(err)
	}
}

```

## Other Methods and Notes
### Cancelling an Order
```go
orderId := wallet.Orders[0].Id
wallet, err := client.CancelOrder(orderId)
```

### Reading Balances and Orders
To get the non-empty balances and orders on a wallet:
```go
balances := wallet.GetNonzeroBalances()
orders := wallet.GetNonzeroOrders()
```

The types on these balance fields are `wallet.Scalar` types. These represent values in our zero-knowledge proof system, but can sometimes be difficult to work with otherwise. For that reason, the scalar type implements a few methods that convert to/from more ergonomic types.
```go
amount := wallet.Balances[0].Amount // type `Scalar`
amtBigint := amount.ToBigInt() // type `*big.Int`
amtHexString := amount.ToHexString() // type `string`
```

### `MAX_BALANCES` and `MAX_ORDERS`
Because our system encodes all its computation in zero-knowledge "circuits", the size of each wallet must be known ahead of time. To this end, we impose the restriction that each wallet has at most `MAX_BALANCES = 10` balances, and `MAX_ORDERS = 4` orders. 

The SDK and the relayer will prevent you from allocating more balances and orders than are allowed.

# External (Atomic) Matching
We also allow for matches to be generated _externally_; meaning generated as a match between a Renegade user -- with state committed into the darkpool -- and an external user, with no state in the darkpool.

To generate an external match, a client may request an `ExternalMatchBundle` from the relayer. This type contains:
- The result of the match, including the amount and mint (erc20 address) of each token in the match. This can be a partial match; the external order may not be fully filled.
- A transaction that the client can submit on-chain to settle the match.

When the protocol receives such a transaction, it will update the internal party's state to reflect the match, and settle any obligations to the external party via ERC20 transfers.

As such, the external party must approve the darkpool contract to spend the tokens it _sells_ to the internal party before the transaction can be successfully submitted.

An example of how to use this functionality is below:
```go
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
	apiKey          = "..." // Issued by Renegade
	apiSecret       = "..." // Issued by Renegade
	quoteMint       = "0xdf8d259c04020562717557f2b5a3cf28e92707d1" // USDC on Arbitrum Sepolia
	baseMint        = "0xc3414a7ef14aaaa9c4522dfc00a4e66e74e9c25a" // wETH on Arbitrum Sepolia
	rpcUrl          = "..." // replace with your RPC URL
)

func getEthClient() (*ethclient.Client, error) {
	return ethclient.Dial(rpcUrl)
}

func getPrivateKey() (*ecdsa.PrivateKey, error) {
	privKeyHex := os.Getenv("PRIVATE_KEY")
	if privKeyHex == "" {
		return nil, fmt.Errorf("PRIVATE_KEY environment variable not set")
	}
	return crypto.HexToECDSA(privKeyHex)
}

func main() {
	// ... Token Approvals to the Darkpool Contract ... //

	// Build a client
	apiSecretKey, err := new(wallet.HmacKey).FromBase64String(apiSecret)
	if err != nil {
		panic(err)
	}
	externalMatchClient := external_match_client.NewTestnetExternalMatchClient(apiKey, &apiSecretKey)
	if err != nil {
		panic(err)
	}

	// Request an external match
	amount := new(big.Int).SetUint64(1000000000000000000) // 1 wETH
	minFillSize := big.NewInt(0)
	order, _ := api_types.NewExternalOrderBuilder().
		WithQuoteMint(quoteMint).
		WithBaseMint(baseMint).
		// Note that `WithQuoteAmount` can be used to specify the volume denominated in the quote token
		WithBaseAmount(api_types.Amount(*amount)).
		WithSide("Sell").
		WithMinFillSize(api_types.Amount(*minFillSize)).
		Build()
	externalMatchBundle, err := externalMatchClient.GetExternalMatchBundle(&order)
	if err != nil {
		panic(err)
	}

	if externalMatchBundle == nil {
		fmt.Println("No match found")
		return
	}

	// Submit the bundle to the sequencer
	if err := submitBundle(*externalMatchBundle); err != nil {
		panic(err)
	}
}

// submitBundle forwards an external match bundle to the sequencer
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
		ChainID:   big.NewInt(421614), // Sepolia chain ID
		Nonce:     nonce,
		GasTipCap: gasPrice,
		GasFeeCap: new(big.Int).Mul(gasPrice, big.NewInt(2)),
		Gas:       uint64(10000000),
		To:        &bundle.SettlementTx.To,
		Value:     bundle.SettlementTx.Value,
		Data:      []byte(bundle.SettlementTx.Data),
	})

	// Sign and send transaction
	signer := types.LatestSignerForChainID(big.NewInt(421614 /* arbitrum sepolia */))
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
```

### Supported Tokens
Renegade supports a specific set of tokens for external matches. These can be found at:
- [Testnet (Arbitrum Sepolia)](https://github.com/renegade-fi/token-mappings/blob/main/testnet.json)
- [Mainnet (Arbitrum One)](https://github.com/renegade-fi/token-mappings/blob/main/mainnet.json)

*Note:* For external matches, Renegade supports swapping native ETH directly. To do so, specify the `baseMint` as `0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE`.

In testnet, we use a set of mock ERC20s that match the mainnet tokens. For convenience while testing, you can use the Renegade faucet to fund your wallet with testnet tokens.
This is most easily accessed through the API using the curl request below
```curl
curl --request POST \
  --url https://testnet.trade.renegade.fi/api/faucet \
  --header 'Content-Type: application/json' \
  --data '{
  "tokens": [
    {
      "ticker": "WETH",
      "amount": "1"
    },
    {
      "ticker": "USDC",
      "amount": "10000"
    }
  ],
  "address": "<ADDRESS>"
}'
```

Note that the `amount` fields here are decimal adjusted, e.g. 1 WETH here is 10^18 wei.
