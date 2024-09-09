package api_types

import (
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/stretchr/testify/assert"
	"renegade.fi/golang-sdk/wallet"
)

func TestScalarUintLimbConversion(t *testing.T) {
	scalar, err := wallet.RandomScalar()
	assert.NoError(t, err)

	// Convert to and from uint32 limbs
	limbs := scalarToUintLimbs(scalar)
	recoveredScalar := scalarFromUintLimbs(limbs)

	// Assert equality
	assert.Equal(t, scalar, recoveredScalar, "Recovered scalar should match original")

	// Test with zero
	zeroScalar := wallet.Scalar(fr.NewElement(0))
	zeroLimbs := scalarToUintLimbs(zeroScalar)
	recoveredZeroScalar := scalarFromUintLimbs(zeroLimbs)
	assert.Equal(t, zeroScalar, recoveredZeroScalar, "Recovered zero scalar should match original")
}