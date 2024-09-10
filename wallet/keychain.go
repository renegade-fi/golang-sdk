package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

// HmacKey is a symmetric key for HMAC-SHA256
type HmacKey [32]byte

// ToHexString converts the HMAC key to a hex string
func (k *HmacKey) ToHexString() string {
	return hex.EncodeToString(k[:])
}

// FromHexString converts a hex string to an HMAC key
func (k *HmacKey) FromHexString(hexString string) (HmacKey, error) {
	hexString = preprocessHexString(hexString)
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		return HmacKey{}, err
	}

	if len(bytes) != 32 {
		return HmacKey{}, errors.New("HMAC key must be 32 bytes")
	}

	copy(k[:], bytes)
	return *k, nil
}

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
	if len(xScalars) > 2 || len(yScalars) > 2 {
		return nil, errors.New("public key is not on the curve")
	}

	// Pad xScalars and yScalars to length 2 if needed
	for len(xScalars) < 2 {
		xScalars = append(xScalars, Scalar{})
	}
	for len(yScalars) < 2 {
		yScalars = append(yScalars, Scalar{})
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

// ToHexString converts the public key to a hex string
func (pk *PublicSigningKey) ToHexString() string {
	bytes := secp256k1.S256().Marshal(pk.X, pk.Y)
	return hex.EncodeToString(bytes)
}

// FromHexString converts a hex string to a public key
func (pk *PublicSigningKey) FromHexString(hexString string) (PublicSigningKey, error) {
	hexString = preprocessHexString(hexString)
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		return PublicSigningKey{}, err
	}

	x, y := secp256k1.S256().Unmarshal(bytes)
	pk.X = x
	pk.Y = y
	pk.Curve = secp256k1.S256()
	return *pk, nil
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

// ToHexString converts the private key to a hex string
func (pk *PrivateSigningKey) ToHexString() string {
	return hex.EncodeToString(pk.D.Bytes())
}

// FromHexString converts a hex string to a private key
func (pk *PrivateSigningKey) FromHexString(hexString string) (PrivateSigningKey, error) {
	hexString = preprocessHexString(hexString)
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		return PrivateSigningKey{}, err
	}

	pk.D = new(big.Int).SetBytes(bytes)
	return *pk, nil
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
	Nonce   Scalar
}

// Keychain is a keychain for the API wallet
type Keychain struct {
	PublicKeys  PublicKeychain
	PrivateKeys PrivateKeychain
}

// SkRoot returns the private root key
func (k *Keychain) SkRoot() *PrivateSigningKey {
	return k.PrivateKeys.SkRoot
}

// FeeEncryptionKey is a public encryption key on the Baby Jubjub curve
// We represent the key in coordinate form with scalar values
type FeeEncryptionKey struct {
	X Scalar
	Y Scalar
}

// ToBytes converts the fee encryption key to a byte slice
func (pk *FeeEncryptionKey) ToBytes() []byte {
	xBytes, yBytes := pk.X.LittleEndianBytes(), pk.Y.LittleEndianBytes()
	return append(xBytes[:], yBytes[:]...)
}

// FromBytes converts a byte slice to a fee encryption key
func (pk *FeeEncryptionKey) FromBytes(bytes []byte) error {
	if len(bytes) != 2*fr.Bytes {
		return errors.New("fee encryption key must be 64 bytes")
	}

	var xBytes [fr.Bytes]byte
	var yBytes [fr.Bytes]byte
	copy(xBytes[:], bytes[:fr.Bytes])
	copy(yBytes[:], bytes[fr.Bytes:])
	pk.X.FromLittleEndianBytes(xBytes)
	pk.Y.FromLittleEndianBytes(yBytes)
	return nil
}

// ToHexString converts the fee encryption key to a hex string
func (pk *FeeEncryptionKey) ToHexString() string {
	return hex.EncodeToString(pk.ToBytes())
}

// FromHexString converts a hex string to a fee encryption key
func (pk *FeeEncryptionKey) FromHexString(hexString string) error {
	hexString = preprocessHexString(hexString)
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		return err
	}

	return pk.FromBytes(bytes)
}
