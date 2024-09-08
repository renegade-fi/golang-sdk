package wallet

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// OrderSide is an enum for the side of an order
type OrderSide int

const (
	Buy OrderSide = iota
	Sell
)

func (s *OrderSide) FromScalars(scalars *ScalarIterator) error {
	scalar, err := scalars.Next()
	if err != nil {
		return err
	}

	elt := fr.Element(scalar)
	if !(elt.IsZero() || elt.IsOne()) {
		return fmt.Errorf("invalid OrderSide value: %v", scalar)
	}

	*s = OrderSide(elt.Uint64())
	return nil
}

func (s *OrderSide) ToScalars() ([]Scalar, error) {
	elt := fr.NewElement(uint64(*s))
	return []Scalar{Scalar(elt)}, nil
}

func (s *OrderSide) NumScalars() int {
	return 1
}

// Order is an order in the Renegade system
type Order struct {
	// BaseMint is the erc20 address of the base asset
	BaseMint Scalar
	// QuoteMint is the erc20 address of the quote asset
	QuoteMint Scalar
	// Amount is the amount of the order
	Amount Scalar
	// Side is the side of the order
	// 0 for buy, 1 for sell
	Side Scalar
	// WorstCasePrice is the worst case price of the order
	WorstCasePrice FixedPoint
}

// NewEmptyOrder creates a new order with all zero values
func NewEmptyOrder() Order {
	return Order{
		BaseMint:       Scalar{},
		QuoteMint:      Scalar{},
		Amount:         Scalar{},
		Side:           Scalar{},
		WorstCasePrice: FixedPoint{},
	}
}
