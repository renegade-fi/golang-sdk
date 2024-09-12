package wallet

import (
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/google/uuid"
)

// OrderSide is an enum for the side of an order
type OrderSide int

const (
	Buy OrderSide = iota
	Sell
)
const OrderSide_BUY = 0
const OrderSide_SELL = 1

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

// OrderBuilder is a builder for Order
type OrderBuilder struct {
	order Order
}

// NewOrderBuilder creates a new OrderBuilder
func NewOrderBuilder() *OrderBuilder {
	return &OrderBuilder{order: NewEmptyOrder()}
}

// WithId sets the Id
func (ob *OrderBuilder) WithId(id uuid.UUID) *OrderBuilder {
	ob.order.Id = id
	return ob
}

// WithQuoteMint sets the QuoteMint
func (ob *OrderBuilder) WithQuoteMint(quoteMint Scalar) *OrderBuilder {
	ob.order.QuoteMint = quoteMint
	return ob
}

// WithQuoteMintHex sets the QuoteMint from a hex string
func (ob *OrderBuilder) WithQuoteMintHex(hexQuoteMint string) *OrderBuilder {
	quoteMint, err := new(Scalar).FromHexString(hexQuoteMint)
	if err != nil {
		panic(err)
	}
	ob.order.QuoteMint = quoteMint
	return ob
}

// WithBaseMint sets the BaseMint
func (ob *OrderBuilder) WithBaseMint(baseMint Scalar) *OrderBuilder {
	ob.order.BaseMint = baseMint
	return ob
}

// WithBaseMintHex sets the BaseMint from a hex string
func (ob *OrderBuilder) WithBaseMintHex(hexBaseMint string) *OrderBuilder {
	baseMint, err := new(Scalar).FromHexString(hexBaseMint)
	if err != nil {
		panic(err)
	}
	ob.order.BaseMint = baseMint
	return ob
}

// WithSide sets the Side
func (ob *OrderBuilder) WithSide(side OrderSide) *OrderBuilder {
	sideScalar, _ := side.ToScalars()
	ob.order.Side = sideScalar[0]
	return ob
}

// WithAmount sets the Amount
func (ob *OrderBuilder) WithAmount(amount Scalar) *OrderBuilder {
	ob.order.Amount = amount
	return ob
}

// WithAmountBigInt sets the Amount from a big.Int
func (ob *OrderBuilder) WithAmountBigInt(amount *big.Int) *OrderBuilder {
	ob.order.Amount = new(Scalar).FromBigInt(amount)
	return ob
}

// WithWorstCasePrice sets the WorstCasePrice
func (ob *OrderBuilder) WithWorstCasePrice(price FixedPoint) *OrderBuilder {
	ob.order.WorstCasePrice = price
	return ob
}

// Build returns the constructed Order
func (ob *OrderBuilder) Build() Order {
	return ob.order
}

// Update NewEmptyOrder to use the builder
func NewEmptyOrder() Order {
	id := uuid.New()
	return Order{
		Id:             id,
		QuoteMint:      Scalar{},
		BaseMint:       Scalar{},
		Side:           Scalar{},
		Amount:         Scalar{},
		WorstCasePrice: FixedPoint{},
	}
}

// Add a new function to create an order with some default values
func NewOrder(quoteMint, baseMint Scalar, side OrderSide, amount Scalar, worstCasePrice FixedPoint) Order {
	return NewOrderBuilder().
		WithId(uuid.New()).
		WithQuoteMint(quoteMint).
		WithBaseMint(baseMint).
		WithSide(side).
		WithAmount(amount).
		WithWorstCasePrice(worstCasePrice).
		Build()
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
