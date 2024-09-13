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
    baseUrl = "http://testnet.cluster0.renegade.fi:3000"
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
amount := big.NewInt(10000000000000000)  // 2^16
wallet, err = client.Deposit(wbtcMint, amount, privateKey)
```
Note that the amount field is _not_ decimal adjusted; for wBTC -- which has 18 decimals -- this translates to 0.01 wBTC. Tokens and their mint addresses that renegade supports can be found at the following locations:
- [Arbitrum Sepolia](https://github.com/renegade-fi/token-mappings/blob/main/testnet.json)
- [Arbitrum One Mainnet](https://github.com/renegade-fi/token-mappings/blob/main/mainnet.json)

**Note:** It is not required that a wallet contain a balance that capitalizes each of their open orders; open orders without a balance backing them will simply not be matched. Therefore, this step is not strictly a prerequisite to the following step in which we place an order.

### Place an Order
Assuming we wish to sell the wBTC that we deposited in the previous step:
```go
btcMint := "0xa91d929ea161688448f61cb3865a6d948d8bd904"
usdcMint := "0x404b26cd9055b35581c68ba9a2b878cca971b0a7"
amount := big.NewInt(10000000000000000) // Sell the whole balance
order := renegade.NewOrderBuilder().
    WithBaseMintHex(baseMint).
    WithQuoteMintHex(quoteMint).
    WithAmountBigInt(amount).
    WithSide(renegade.OrderSide_SELL).
    Build()

wallet, err = client.PlaceOrder(&order)
```
**Note:** For the moment, all pairs are USDC quoted. E.g. Renegade does not currently support selling wBTC/wETH.

Once the order is placed with a balance to capitalize it, the matching engine will match the order with any counter-flow it finds. 

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
    baseUrl = "http://testnet.cluster0.renegade.fi:3000"
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
    amount := big.NewInt(10000000000000000)  // 2^16
    wallet, err = client.Deposit(wbtcMint, amount, privateKey)
    if err != nil {
        log.Fatal(err)
    }

    // Sell 0.01 wBTC
    usdcMint := "0x404b26cd9055b35581c68ba9a2b878cca971b0a7"
    amount, _ := wallet.GetBalance(wbtcMint)
    order := renegade.NewOrderBuilder().
        WithBaseMintHex(wbtcMint).
        WithQuoteMintHex(usdcMint).
        WithAmountBigInt(amount).
        WithSide(renegade.OrderSide_SELL).
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

## Other Worthwhile Methods
### Cancelling an Order
To cancel an open order, 