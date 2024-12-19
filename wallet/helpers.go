package wallet

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"

	renegade_crypto "github.com/renegade-fi/golang-sdk/crypto"
)

// HashScalars hashes a slice of scalars using Poseidon2
func HashScalars(scalars []Scalar) Scalar {
	sponge := renegade_crypto.NewPoseidon2Sponge()
	elts := make([]fr.Element, len(scalars))
	for i, scalar := range scalars {
		elts[i] = fr.Element(scalar)
	}

	res := sponge.Hash(elts)
	return Scalar(res)
}
