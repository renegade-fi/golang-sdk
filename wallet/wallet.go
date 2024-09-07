package wallet

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/google/uuid"
)

type Scalar fr.Element

type OrderSide int

const (
	Buy OrderSide = iota
	Sell
)

type Order struct {
	BaseAsset      Scalar    `json:"base_asset"`
	QuoteAsset     Scalar    `json:"quote_asset"`
	Amount         Scalar    `json:"amount"`
	Side           OrderSide `json:"side"`
	WorstCasePrice Scalar    `json:"worst_case_price"`
}

type Balance struct {
	Mint   Scalar `json:"mint"`
	Amount Scalar `json:"amount"`
}

type Wallet struct {
	ID     uuid.UUID           `json:"id"`
	Orders map[uuid.UUID]Order `json:"orders"`
}
