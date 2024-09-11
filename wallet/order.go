package wallet

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/google/uuid"
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
	// ID is the id of the order
	Id uuid.UUID `scalar_serialize:"skip"`
	// QuoteMint is the erc20 address of the quote asset
	QuoteMint Scalar
	// BaseMint is the erc20 address of the base asset
	BaseMint Scalar
	// Side is the side of the order
	// 0 for buy, 1 for sell
	Side Scalar
	// Amount is the amount of the order
	Amount Scalar
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

// IsZero returns whether the volume of the order is zero
func (o *Order) IsZero() bool {
	return o.Amount.IsZero()
}

// NewOrder appends an order to the wallet
func (w *Wallet) NewOrder(order Order) error {
	// Find the first order that may be replaced
	if idx := w.findReplaceableOrder(); idx != -1 {
		w.Orders[idx] = order
	} else if len(w.Orders) < MaxOrders {
		w.Orders = append(w.Orders, order)
	} else {
		return fmt.Errorf("wallet already has the maximum number of orders")
	}

	return nil
}

// findReplaceableOrder finds the first order that may be replaced by the new order
// Returns the index of the order to replace, or -1 if no order may be replaced
func (w *Wallet) findReplaceableOrder() int {
	for i, existingOrder := range w.Orders {
		if existingOrder.IsZero() {
			return i
		}
	}

	return -1
}

// CancelOrder cancels an order by ID
func (w *Wallet) CancelOrder(orderId uuid.UUID) error {
	// Find the order to cancel
	idx := w.findOrder(orderId)
	if idx == -1 {
		return fmt.Errorf("order not found")
	}

	// Remove the order and append an empty order to the end
	w.Orders = append(w.Orders[:idx], append(w.Orders[idx+1:], NewEmptyOrder())...)
	return nil
}

// findOrder finds the index of an order with the given ID, or -1 if no order has the given ID
func (w *Wallet) findOrder(orderId uuid.UUID) int {
	for i, order := range w.Orders {
		if order.Id == orderId {
			return i
		}
	}

	return -1
}
