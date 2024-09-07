package wallet

import (
	"crypto/ecdsa"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/google/uuid"
)

// Scalar is a scalar field element from the bn254 curve
type Scalar fr.Element

// OrderSide is an enum for the side of an order
type OrderSide int

const (
	Buy OrderSide = iota
	Sell
)

// Order is an order in the Renegade system
type Order struct {
	BaseAsset      Scalar    `json:"base_asset"`
	QuoteAsset     Scalar    `json:"quote_asset"`
	Amount         Scalar    `json:"amount"`
	Side           OrderSide `json:"side"`
	WorstCasePrice Scalar    `json:"worst_case_price"`
}

// Balance is a balance in the Renegade system
type Balance struct {
	Mint   Scalar `json:"mint"`
	Amount Scalar `json:"amount"`
}

// HmacKey is a symmetric key for HMAC-SHA256
type HmacKey [32]byte

// PrivateKeychain is a private keychain for the API wallet
type PrivateKeychain struct {
	SkRoot       *ecdsa.PrivateKey
	SkMatch      Scalar
	SymmetricKey HmacKey
}

// PublicKeychain is a public keychain for the API wallet
type PublicKeychain struct {
	PkRoot  ecdsa.PublicKey
	PkMatch Scalar
}

// Keychain is a keychain for the API wallet
type Keychain struct {
	PublicKeys  PublicKeychain
	PrivateKeys PrivateKeychain
	Nonce       uint64
}

// Wallet is a wallet in the Renegade system
type Wallet struct {
	ID     uuid.UUID           `json:"id"`
	Orders map[uuid.UUID]Order `json:"orders"`
}
