package api_types

import (
	"math/big"

	"renegade.fi/golang-sdk/wallet"
)

// scalarToUintLimbs converts a scalar to an array of uint32 limbs
func scalarToUintLimbs(s wallet.Scalar) [secretShareLimbCount]uint32 {
	bigint := s.ToBigInt()
	limbs := [secretShareLimbCount]uint32{}
	for i := 0; i < secretShareLimbCount; i++ {
		if bigint.BitLen() > 0 {
			limbs[i] = uint32(bigint.Uint64())
			bigint.Rsh(bigint, 32)
		} else {
			break
		}
	}

	return limbs
}

// scalarFromUintLimbs converts an array of uint32 limbs to a scalar
func scalarFromUintLimbs(limbs [secretShareLimbCount]uint32) wallet.Scalar {
	bigint := new(big.Int)
	for i := secretShareLimbCount - 1; i >= 0; i-- {
		bigint.Lsh(bigint, 32)
		bigint.Or(bigint, big.NewInt(int64(limbs[i])))
	}
	return new(wallet.Scalar).FromBigInt(bigint)
}
