package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"golang.org/x/crypto/sha3"
)

// derivationKeyMessage is the message that is signed to derive the derivation key
// From which all other keys can be derived
const derivationKeyMessage = "Unlock your Renegade Wallet on chain ID:"

// rootKeyMessage is the message that is signed to derive the root key
const rootKeyMessage = "root key"

// symmetricKeyMessage is the message that is signed to derive the symmetric key
const symmetricKeyMessage = "symmetric key"

// matchKeyMessage is the message that is signed to derive the match key
const matchKeyMessage = "match key"

// DeriveKeychain derives the keychain from the private key
func DeriveKeychain(pkey *ecdsa.PrivateKey, chainId uint64) (*ecdsa.PrivateKey, error) {
	// Create the derivation key
	derivationKey, err := createDerivationKey(pkey, chainId)
	if err != nil {
		return nil, err
	}

	// Derive the root key
	rootKey, err := deriveRootKey(derivationKey)
	if err != nil {
		return nil, err
	}

	// Derive the symmetric key
	symmetricKey, err := deriveSymmetricKey(derivationKey)
	if err != nil {
		return nil, err
	}
	fmt.Println("symmetricKey", symmetricKey)

	// Derive the match key
	matchKey, err := deriveMatchKey(derivationKey)
	if err != nil {
		return nil, err
	}
	fmt.Println("matchKey", matchKey)

	return rootKey, nil
}

// createDerivationKey creates a new private key from the signature
func createDerivationKey(pkey *ecdsa.PrivateKey, chainId uint64) (*ecdsa.PrivateKey, error) {
	message := []byte(fmt.Sprintf("%s%d", derivationKeyMessage, chainId))
	keyBytes, err := getExtendedSigBytes(message, pkey)
	if err != nil {
		return nil, err
	}

	derivedKey, err := secpKeyFromBytes(keyBytes)
	if err != nil {
		return nil, err
	}

	return derivedKey, nil
}

// deriveRootKey derives the `sk_root` key from the derivation key
func deriveRootKey(derivationKey *ecdsa.PrivateKey) (*ecdsa.PrivateKey, error) {
	message := []byte(rootKeyMessage)
	keyBytes, err := getExtendedSigBytes(message, derivationKey)
	if err != nil {
		return nil, err
	}

	rootKey, err := secpKeyFromBytes(keyBytes)
	if err != nil {
		return nil, err
	}

	return rootKey, nil
}

// deriveSymmetricKey derives the symmetric key from the derivation key
func deriveSymmetricKey(rootKey *ecdsa.PrivateKey) ([]byte, error) {
	message := []byte(symmetricKeyMessage)
	return getSigBytes(rootKey, message)
}

// deriveMatchKey derives the secret match key from the derivation key
func deriveMatchKey(derivationKey *ecdsa.PrivateKey) (Scalar, error) {
	message := []byte(matchKeyMessage)
	return deriveScalar(message, derivationKey)
}

// secpKeyFromBytes creates a secp256k1 private key from a byte slice
func secpKeyFromBytes(b []byte) (*ecdsa.PrivateKey, error) {
	if len(b) != 64 {
		return nil, fmt.Errorf("secpKeyFromBytes: input must be 64 bytes, extend before using")
	}

	// Reduce the extended signature to the secp256k1 scalar field
	curve := secp256k1.S256()
	reduced := new(big.Int).SetBytes(b)
	reduced.Mod(reduced, curve.Params().N)

	// Create a new private key
	derivedKey := new(ecdsa.PrivateKey)
	derivedKey.PublicKey.Curve = curve
	derivedKey.D = reduced
	derivedKey.PublicKey.X, derivedKey.PublicKey.Y = curve.ScalarBaseMult(reduced.Bytes())

	return derivedKey, nil

}

// deriveScalar derives a bn254 scalar from a message
func deriveScalar(message []byte, pkey *ecdsa.PrivateKey) (Scalar, error) {
	bytes, err := getExtendedSigBytes(message, pkey)
	if err != nil {
		return Scalar{}, err
	}

	var scalar fr.Element
	scalar.SetBytes(bytes)
	return Scalar(scalar), nil
}

// getSigBytes signs the message and returns a keccak256 hash of the signature
func getSigBytes(pkey *ecdsa.PrivateKey, message []byte) ([]byte, error) {
	signature, err := signMessage(pkey, message)
	if err != nil {
		return nil, err
	}

	return keccak256(signature), nil
}

// getExtendedSigBytes signs the message and extends the signature to 64 bytes
func getExtendedSigBytes(message []byte, pkey *ecdsa.PrivateKey) ([]byte, error) {
	sigBytes, err := getSigBytes(pkey, message)
	if err != nil {
		return nil, err
	}

	return extendTo64Bytes(sigBytes)
}

// extendTo64Bytes extends the byte slice to 64 bytes by keccak256 hashing
// We use this method to reduce into a 256 bit field with sufficient entropy
func extendTo64Bytes(b []byte) ([]byte, error) {
	if len(b) != 32 {
		fmt.Println("n_bytes", len(b))
		return nil, fmt.Errorf("extendTo64Bytes: input must be 32 bytes")
	}

	// Copy in the original bytes
	extended := make([]byte, 64)
	copy(extended[:len(b)], b)

	// Hash the original bytes to get the top 64 bytes
	topBytes := keccak256(b)
	copy(extended[len(b):], topBytes[:64-len(b)])

	return extended, nil
}

// Keccak256 calculates the Keccak256 hash of the input data
func keccak256(data []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	return hash.Sum(nil)
}

func signMessage(pkey *ecdsa.PrivateKey, message []byte) ([]byte, error) {
	return ecdsa.SignASN1(rand.Reader, pkey, message)
}
