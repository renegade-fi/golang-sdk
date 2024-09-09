package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/stretchr/testify/assert"
)

func TestScalarToFromHex(t *testing.T) {
	// Generate a random scalar
	randomScalar, err := RandomScalar()
	if err != nil {
		t.Fatalf("Failed to generate random scalar: %v", err)
	}

	// Convert to hex string and back
	hexString := randomScalar.ToHexString()
	var newScalar Scalar
	_, err = newScalar.FromHexString(hexString)
	assert.NoError(t, err)
	assert.Equal(t, randomScalar, newScalar)
}

func TestScalarToBigIntAndBack(t *testing.T) {
	// Generate a random scalar
	randomScalar, err := RandomScalar()
	assert.NoError(t, err, "Failed to generate random scalar")

	// Convert to and from big.Int
	bigInt := randomScalar.ToBigInt()
	var newScalar Scalar
	newScalar.FromBigInt(bigInt)

	// Compare the original and new scalar
	assert.Equal(t, randomScalar, newScalar, "Scalar -> BigInt -> Scalar conversion failed")
}

func TestNewEmptyWallet(t *testing.T) {
	// Generate a random private key
	privateKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	assert.NoError(t, err, "Failed to generate random private key")

	// Create a new empty wallet and ensure it doesn't error
	_, err = NewEmptyWallet(privateKey, 1 /* chainId */)
	assert.NoError(t, err, "Failed to create new empty wallet")
}
