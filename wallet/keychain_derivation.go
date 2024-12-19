package wallet

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/google/uuid"

	renegade_crypto "github.com/renegade-fi/golang-sdk/crypto"
)

const (
	// derivationKeyMessage is the message that is signed to derive the derivation key
	// From which all other keys can be derived
	derivationKeyMessage = "Unlock your Renegade Wallet on chain ID:"

	// rootKeyMessage is the message that is signed to derive the root key
	rootKeyMessage = "root key"

	// symmetricKeyMessage is the message that is signed to derive the symmetric key
	symmetricKeyMessage = "symmetric key"

	// matchKeyMessage is the message that is signed to derive the match key
	matchKeyMessage = "match key"

	// blinderSeedMessage is the message used to derive the blinder stream seed
	blinderSeedMessage = "blinder seed"

	// shareSeedMessage is the message used to derive the secret share stream seed
	shareSeedMessage = "share seed"

	// walletIdMessage is the message used to derive the wallet ID
	walletIDMessage = "wallet id"

	// walletIdNumBytes is the number of bytes in the wallet ID
	walletIDNumBytes = 16
)

// DeriveKeychain derives the keychain from the private key
func DeriveKeychain(pkey *ecdsa.PrivateKey, chainID uint64) (*Keychain, error) {
	// Create the derivation key
	derivationKey, err := createDerivationKey(pkey, chainID)
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

	// Derive the match key
	matchKey, err := deriveMatchKey(derivationKey)
	if err != nil {
		return nil, err
	}

	keychain := createKeychain(rootKey, matchKey, symmetricKey)
	return keychain, nil
}

// DeriveWalletSeeds derives the blinder and secret share seeds from the derivation key
func DeriveWalletSeeds(privateKey *ecdsa.PrivateKey, chainID uint64) (
	blinderSeed,
	shareSeed Scalar,
	err error,
) {
	// Create the derivation key
	derivationKey, err := createDerivationKey(privateKey, chainID)
	if err != nil {
		return Scalar{}, Scalar{}, err
	}

	blinderSeed, err = deriveScalar([]byte(blinderSeedMessage), derivationKey)
	if err != nil {
		return Scalar{}, Scalar{}, err
	}

	shareSeed, err = deriveScalar([]byte(shareSeedMessage), derivationKey)
	if err != nil {
		return Scalar{}, Scalar{}, err
	}

	return blinderSeed, shareSeed, nil
}

// DeriveWalletID derives the wallet ID from the private key
func DeriveWalletID(privateKey *ecdsa.PrivateKey, chainID uint64) (uuid.UUID, error) {
	// Create the derivation key
	derivationKey, err := createDerivationKey(privateKey, chainID)
	if err != nil {
		return uuid.Nil, err
	}

	// Derive the wallet ID
	walletIDBytes, err := getExtendedSigBytes([]byte(walletIDMessage), derivationKey)
	if err != nil {
		return uuid.Nil, err
	}

	// Convert the bytes to a UUID
	walletID, err := uuid.FromBytes(walletIDBytes[:walletIDNumBytes])
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create UUID from bytes: %v", err)
	}

	return walletID, nil

}

// createKeychain creates a new keychain from the private keys
func createKeychain(skRoot *ecdsa.PrivateKey, skMatch Scalar, symmetricKey HmacKey) *Keychain {
	privateKeys := PrivateKeychain{
		SkRoot:       (*PrivateSigningKey)(skRoot),
		SkMatch:      skMatch,
		SymmetricKey: symmetricKey,
	}

	pkRoot := PublicSigningKey(skRoot.PublicKey)
	pkMatch := publicMatchKey(skMatch)
	publicKeys := PublicKeychain{
		PkRoot:  pkRoot,
		PkMatch: pkMatch,
		Nonce:   Scalar{},
	}

	return &Keychain{
		PublicKeys:  publicKeys,
		PrivateKeys: privateKeys,
	}
}

// createDerivationKey creates a new private key from the signature
func createDerivationKey(pkey *ecdsa.PrivateKey, chainID uint64) (*ecdsa.PrivateKey, error) {
	message := []byte(fmt.Sprintf("%s%d", derivationKeyMessage, chainID))
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
func deriveSymmetricKey(rootKey *ecdsa.PrivateKey) (HmacKey, error) {
	message := []byte(symmetricKeyMessage)
	bytes, err := getSigBytes(rootKey, message)
	if err != nil {
		return HmacKey{}, err
	}

	return HmacKey(bytes), nil
}

// deriveMatchKey derives the secret match key from the derivation key
func deriveMatchKey(derivationKey *ecdsa.PrivateKey) (Scalar, error) {
	message := []byte(matchKeyMessage)
	return deriveScalar(message, derivationKey)
}

// publicMatchKey derives the public match key from the private match key
func publicMatchKey(skMatch Scalar) Scalar {
	sponge := renegade_crypto.NewPoseidon2Sponge()
	res := sponge.Hash([]fr.Element{fr.Element(skMatch)})
	return Scalar(res)
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

// getSigBytes signs the message and returns a Keccak256 hash of the signature
func getSigBytes(pkey *ecdsa.PrivateKey, message []byte) ([]byte, error) {
	signature, err := signMessage(pkey, message)
	if err != nil {
		return nil, err
	}

	return crypto.Keccak256(signature), nil
}

// getExtendedSigBytes signs the message and extends the signature to 64 bytes
func getExtendedSigBytes(message []byte, pkey *ecdsa.PrivateKey) ([]byte, error) {
	sigBytes, err := getSigBytes(pkey, message)
	if err != nil {
		return nil, err
	}

	return extendTo64Bytes(sigBytes)
}

// extendTo64Bytes extends the byte slice to 64 bytes by Keccak256 hashing
// We use this method to reduce into a 256 bit field with sufficient entropy
func extendTo64Bytes(b []byte) ([]byte, error) {
	if len(b) != 32 {
		return nil, fmt.Errorf("extendTo64Bytes: input must be 32 bytes")
	}

	// Copy in the original bytes
	extended := make([]byte, 64)
	copy(extended[:len(b)], b)

	// Hash the original bytes to get the top 64 bytes
	topBytes := crypto.Keccak256(b)
	copy(extended[len(b):], topBytes[:64-len(b)])

	return extended, nil
}

// signMessage signs a message deterministically using RFC 6979
func signMessage(pkey *ecdsa.PrivateKey, message []byte) ([]byte, error) {
	// Hash the message using Keccak256
	messageHash := crypto.Keccak256(message)

	// Sign the hash deterministically
	signature, err := crypto.Sign(messageHash, pkey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %v", err)
	}

	return signature, nil
}
