package api_types //nolint:revive

import (
	"math/big"

	"github.com/renegade-fi/golang-sdk/wallet"
)

// ScalarToUintLimbs converts a scalar to an array of uint32 limbs
func ScalarToUintLimbs(s wallet.Scalar) ScalarLimbs {
	bigint := s.ToBigInt()
	limbs := [secretShareLimbCount]uint32{}
	for i := 0; i < secretShareLimbCount; i++ {
		if bigint.BitLen() > 0 {
			limbs[i] = uint32(bigint.Uint64()) //nolint:gosec
			bigint.Rsh(bigint, 32)
		} else {
			break
		}
	}

	return limbs
}

// ScalarFromUintLimbs converts an array of uint32 limbs to a scalar
func ScalarFromUintLimbs(limbs ScalarLimbs) wallet.Scalar {
	bigint := new(big.Int)
	for i := secretShareLimbCount - 1; i >= 0; i-- {
		bigint.Lsh(bigint, 32)
		bigint.Or(bigint, big.NewInt(int64(limbs[i])))
	}
	return new(wallet.Scalar).FromBigInt(bigint)
}
