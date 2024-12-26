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

---

# External (Atomic) Matching
We also allow for matches to be generated _externally_; meaning generated as a match between a Renegade user -- with state committed into the darkpool -- and an external user, with no state in the darkpool.

To generate an external match, a client may request an `ExternalMatchBundle` from the relayer. This type contains:
- The result of the match, including the amount and mint (erc20 address) of each token in the match. This can be a partial match; the external order may not be fully filled.
- A transaction that the client can submit on-chain to settle the match.

When the protocol receives such a transaction, it will update the internal party's state to reflect the match, and settle any obligations to the external party via ERC20 transfers.

As such, the external party must approve the darkpool contract to spend the tokens it _sells_ to the internal party before the transaction can be successfully submitted.

### Generating an External Match

Generating an external match breaks down into three steps:
1. Fetch a quote for the order.
2. If the quote is acceptable, assemble the quote into a **bundle**. Bundles contain a transaction that may be used to settle the trade on-chain.
3. Submit the settlement transaction on-chain.

### Example
A full example can be found in [`examples/01_external_match/main.go`](examples/01_external_match/main.go).

<details>
<summary>Example Code</summary>

```go
// ... See `examples/01_external_match/main.go` for the prelude ... //

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

	fmt.Println("Bundle submitted successfully!\n")
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
```

</details>

## Bundle Structure
The *quote* returned by the relayer for an external match has the following structure:
- `Order`: The original external order
- `MatchResult`: The result of the match, including:
- `Fees`: The fees for the match
    - `RelayerFee`: The fee paid to the relayer
    - `ProtocolFee`: The fee paid to the protocol
- `Receive`: The asset transfer the external party will receive, *after fees are deducted*.
    - `Mint`: The token address
    - `Amount`: The amount to receive
- `Send`: The asset transfer the external party needs to send. No fees are charged on the send transfer.  (same fields as `Receive`) 
- `Price`: The price used for the match
- `Timestamp`: The timestamp of the quote

When assembled into a bundle (returned from `AssembleExternalQuote` or `GetExternalMatchBundle`), the structure is as follows:
- `MatchResult`: The final match result
- `Fees`: The fees to be paid
- `Receive`: The asset transfer the external party will receive
- `Send`: The asset transfer the external party needs to send
- `SettlementTx`: The transaction to submit on-chain
    - `Type`: The transaction type
    - `To`: The contract address
    - `Data`: The calldata
    - `Value`: The ETH value to send

See example [`02_external_quote_validation`](examples/02_external_quote_validation/main.go) for an example of using these fields to validate a quote before submitting it.

## Supported Tokens
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
