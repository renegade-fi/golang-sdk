package wallet

import (
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

// HmacKey is a symmetric key for HMAC-SHA256
type HmacKey [32]byte

// PublicSigningKey is a verification key over the secp256k1 curve
type PublicSigningKey ecdsa.PublicKey

func bigintToScalarLimbs(b big.Int) []Scalar {
	localB := new(big.Int).Set(&b) // Create a local copy
	scalarMod := fr.Modulus()
	var limbs []Scalar
	for localB.Sign() != 0 {
		word := new(big.Int).Mod(localB, scalarMod)
		elt := fr.Element{}
		elt.SetBigInt(word)

		limbs = append(limbs, Scalar(elt))
		localB.Div(localB, scalarMod)
	}
	return limbs
}

func scalarLimbsToBigInt(limbs []Scalar) *big.Int {
	scalarMod := fr.Modulus()
	b := new(big.Int)
	for i := len(limbs) - 1; i >= 0; i-- {
		elt := fr.Element(limbs[i])

		var eltBigint big.Int
		elt.BigInt(&eltBigint)
		b.Add(b, &eltBigint)

		if i > 0 {
			b.Mul(b, scalarMod)
		}
	}
	return b
}

func (pk *PublicSigningKey) ToScalars() ([]Scalar, error) {
	xScalars := bigintToScalarLimbs(*pk.X)
	yScalars := bigintToScalarLimbs(*pk.Y)
	if len(xScalars) != 2 || len(yScalars) != 2 {
		return nil, errors.New("public key is not on the curve")
	}

	return []Scalar{xScalars[0], xScalars[1], yScalars[0], yScalars[1]}, nil
}

func (pk *PublicSigningKey) FromScalars(scalars *ScalarIterator) error {
	xScalars := make([]Scalar, 2)
	yScalars := make([]Scalar, 2)

	for i := 0; i < 2; i++ {
		x, err := scalars.Next()
		if err != nil {
			return err
		}
		xScalars[i] = x
	}

	for i := 0; i < 2; i++ {
		y, err := scalars.Next()
		if err != nil {
			return err
		}
		yScalars[i] = y
	}

	pk.X = scalarLimbsToBigInt(xScalars)
	pk.Y = scalarLimbsToBigInt(yScalars)
	pk.Curve = secp256k1.S256()
	return nil
}

func (pk *PublicSigningKey) NumScalars() int {
	return 4
}

type PrivateSigningKey ecdsa.PrivateKey

func (pk *PrivateSigningKey) ToScalars() ([]Scalar, error) {
	limbs := bigintToScalarLimbs(*pk.D)
	return limbs, nil
}

func (pk *PrivateSigningKey) FromScalars(scalars *ScalarIterator) error {
	d, err := scalars.Next()
	if err != nil {
		return err
	}
	pk.D = scalarLimbsToBigInt([]Scalar{d})
	return nil
}

func (pk *PrivateSigningKey) NumScalars() int {
	return 2
}

// PrivateKeychain is a private keychain for the API wallet
type PrivateKeychain struct {
	SkRoot       *PrivateSigningKey
	SkMatch      Scalar
	SymmetricKey HmacKey
}

// PublicKeychain is a public keychain for the API wallet
type PublicKeychain struct {
	PkRoot  PublicSigningKey
	PkMatch Scalar
	Nonce   Uint64
}

// Keychain is a keychain for the API wallet
type Keychain struct {
	PublicKeys  PublicKeychain
	PrivateKeys PrivateKeychain
}

// FeeEncryptionKey is a public encryption key on the Baby Jubjub curve
// We represent the key in coordinate form with scalar values
type FeeEncryptionKey struct {
	X Scalar
	Y Scalar
}
