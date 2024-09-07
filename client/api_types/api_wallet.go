package api_wallet

import (
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/google/uuid"
)

// The number of u32 limbs in the serialized form of a secret share
const secretShareLimbCount = 8 // 256 bits

type Scalar fr.Element
type Amount big.Int

func (a *Amount) String() string {
	return (*big.Int)(a).String()
}

func (a *Amount) MarshalJSON() ([]byte, error) {
	return []byte(a.String()), nil
}

func (a *Amount) SetString(s string, base int) error {
	i, ok := new(big.Int).SetString(s, base)
	if !ok {
		return fmt.Errorf("invalid number: %s", s)
	}
	*a = Amount(*i)
	return nil
}

func (a *Amount) UnmarshalJSON(b []byte) error {
	return a.SetString(string(b), 10)
}

// OrderSide represents the side of an order (buy or sell)
type OrderSide uint8

const (
	Buy OrderSide = iota
	Sell
)

// OrderType represents the type of an order (midpoint or limit)
type OrderType uint8

const (
	Midpoint OrderType = iota
	Limit
)

// Order is an order in a Renegade wallet
type Order struct {
	// The id of the order
	Id uuid.UUID `json:"id"`
	// The mint (erc20 address) of the base asset
	// As a hex string
	BaseMint string `json:"base_mint"`
	// The mint (erc20 address) of the quote asset
	// As a hex string
	QuoteMint string `json:"quote_mint"`
	// The amount of the base asset to buy/sell
	Amount Amount `json:"amount"`
	// The side of the order
	Side string `json:"side"`
	// The worst case price to execute the order at
	// The serialized form of this is the `Scalar` representation of the fixed point,
	// i.e. if a fixed point value represents `r`, this value is `floor(r << PRECISION)`
	WorstCasePrice string `json:"worst_case_price"`
}

// Balance is a balance in a Renegade wallet
type Balance struct {
	// The mint (erc20 address) of the asset
	Mint string `json:"mint"`
	// The amount of the asset
	Amount Amount `json:"amount"`
	// The amount of this balance owed to the managing relayer cluster
	RelayerFeeBalance Amount `json:"relayer_fee_balance"`
	// The amount of this balance owed to the protocol
	ProtocolFeeBalance Amount `json:"protocol_fee_balance"`
}

// PublicKeychain is a public keychain in the Renegade system
type PublicKeychain struct {
	// The public root key of the wallet
	// As a hex string
	PkRoot string `json:"pk_root"`
	// The public match key of the wallet
	// As a hex string
	PkMatch string `json:"pk_match"`
}

// ApiPrivateKeychain represents a private keychain for the API wallet
type PrivateKeychain struct {
	// The private root key of the wallet
	// As a hex string, optional
	SkRoot *string `json:"sk_root,omitempty"`
	// The private match key of the wallet
	// As a hex string
	SkMatch string `json:"sk_match"`
	// The symmetric key of the wallet
	// As a hex string
	SymmetricKey string `json:"symmetric_key"`
}

// Keychain represents a keychain API type that maintains all keys as hex strings
type Keychain struct {
	// The public keychain
	PublicKeys PublicKeychain `json:"public_keys"`
	// The private keychain
	PrivateKeys PrivateKeychain `json:"private_keys"`
	// The nonce of the keychain
	Nonce uint64 `json:"nonce"`
}

// Wallet is a wallet in the Renegade system
type Wallet struct {
	// Identifier
	Id uuid.UUID `json:"id"`
	// The orders maintained by this wallet
	Orders []Order `json:"orders"`
	// The balances maintained by the wallet to cover orders
	Balances []Balance `json:"balances"`
	// The keys that authenticate wallet access
	KeyChain Keychain `json:"key_chain"`
	// The managing cluster's public key
	// The public encryption key of the cluster that may collect relayer fees
	// on this wallet
	ManagingCluster string `json:"managing_cluster"`
	// The take rate at which the managing cluster may collect relayer fees on
	// a match
	MatchFee string `json:"match_fee"`
	// The public secret shares of the wallet
	BlindedPublicShares [][secretShareLimbCount]uint32 `json:"blinded_public_shares"`
	// The private secret shares of the wallet
	PrivateShares [][secretShareLimbCount]uint32 `json:"private_shares"`
	// The wallet blinder, used to blind wallet secret shares
	Blinder [secretShareLimbCount]uint32 `json:"blinder"`
}
