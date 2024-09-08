package wallet

import (
	"math/big"
	"math/rand"
	"testing"
)

func TestScalarLimbsToBigInt(t *testing.T) {
	// Sample a random big.Int
	limit := new(big.Int).Lsh(big.NewInt(1), 256)
	r := rand.New(rand.NewSource(0))
	randomBigInt := new(big.Int).Rand(r, limit)

	// Convert to scalar limbs and back
	limbs := bigintToScalarLimbs(*randomBigInt)
	recoveredBigInt := scalarLimbsToBigInt(limbs)

	// Assert equality
	if randomBigInt.Cmp(recoveredBigInt) != 0 {
		t.Errorf("Conversion failed: original %v, recovered %v", randomBigInt, recoveredBigInt)
	}
}
