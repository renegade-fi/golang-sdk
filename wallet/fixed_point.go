package wallet

import (
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// precisionBits is the number of bits of precision in the fixed point number
const precisionBits = 63

// FixedPoint is a fixed point number with a scalar representation
// The scalar represents the value `floor(repr >> 2^PRECISION)`
// For our purposes, the precision is 63 bits
type FixedPoint struct {
	// Repr is the integral representation of the fixed point number
	Repr Scalar
}

// ZeroFixedPoint is the fixed point number 0
func ZeroFixedPoint() FixedPoint {
	return NewFixedPoint(Scalar{})
}

// NewFixedPoint creates a new fixed point number from a scalar representation
func NewFixedPoint(repr Scalar) FixedPoint {
	return FixedPoint{Repr: repr}
}

// FixedPointFromFloat creates a new fixed point number from a float
func FixedPointFromFloat(f float64) FixedPoint {
	bigF := big.NewFloat(f)
	// Shift left by precisionBits
	bigF.Mul(bigF, big.NewFloat(1<<precisionBits))

	// Floor the value
	intF, _ := bigF.Int(nil)

	// Convert to Scalar
	var elt fr.Element
	elt.SetBigInt(intF)

	return FixedPoint{Repr: Scalar(elt)}
}

// ToFloat converts a fixed point number to a float
func (fp FixedPoint) ToFloat() float64 {
	// Convert the repr to a bigF
	bigF := big.NewFloat(0)
	elt := fr.Element(fp.Repr)
	var intF big.Int
	bigF.SetInt(elt.BigInt(&intF))

	// Shift right by precisionBits
	bigF.Quo(bigF, big.NewFloat(1<<precisionBits))

	f, _ := bigF.Float64()
	return f
}

// ToReprDecimalString converts a fixed point number to the base10 string representation of its `repr`
func (fp FixedPoint) ToReprDecimalString() string {
	reprBigint := fp.Repr.ToBigInt()

	// Convert to string with 10 decimal places
	return reprBigint.Text(10 /* base */)
}

// FromReprDecimalString creates a new fixed point number from a decimal string
func (fp *FixedPoint) FromReprDecimalString(s string) (FixedPoint, error) {
	reprBigint, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return FixedPoint{}, fmt.Errorf("failed to convert decimal string to big.Int")
	}

	repr := new(Scalar).FromBigInt(reprBigint)
	fp.Repr = repr
	return *fp, nil
}
