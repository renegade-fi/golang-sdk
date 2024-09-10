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

func TestScalarToFromLittleEndianBytes(t *testing.T) {
	// Generate a random scalar
	randomScalar, err := RandomScalar()
	assert.NoError(t, err, "Failed to generate random scalar")

	// Convert to little-endian bytes and back
	littleEndianBytes := randomScalar.LittleEndianBytes()
	var newScalar Scalar
	newScalar.FromLittleEndianBytes(littleEndianBytes)

	// Compare the original and new scalar
	assert.Equal(t, randomScalar, newScalar, "Scalar -> LittleEndianBytes -> Scalar conversion failed")
}

func TestNewEmptyWallet(t *testing.T) {
	// Generate a random private key
	privateKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	assert.NoError(t, err, "Failed to generate random private key")

	// Create a new empty wallet and ensure it doesn't error
	_, err = NewEmptyWallet(privateKey, 1 /* chainId */)
	assert.NoError(t, err, "Failed to create new empty wallet")
}

func TestWalletReblind(t *testing.T) {
	// Generate a random private key
	privateKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	assert.NoError(t, err, "Failed to generate random private key")

	// Create a new empty wallet
	wallet, err := NewEmptyWallet(privateKey, 1 /* chainId */)
	assert.NoError(t, err, "Failed to create new empty wallet")

	// Add a balance
	balance := Balance{
		Mint:               Scalar{1},
		Amount:             Scalar{2},
		RelayerFeeBalance:  Scalar{3},
		ProtocolFeeBalance: Scalar{4},
	}
	err = wallet.AddBalance(balance)
	assert.NoError(t, err, "Failed to add new balance")

	// Add an order
	order := Order{
		BaseMint:       Scalar{2},
		QuoteMint:      Scalar{3},
		Amount:         Scalar{4},
		Side:           Scalar{0}, // Buy
		WorstCasePrice: FixedPoint{Repr: Scalar{0}},
	}
	err = wallet.NewOrder(order)
	assert.NoError(t, err, "Failed to add new order")

	// Reblind the wallet
	err = wallet.Reblind()
	assert.NoError(t, err, "Failed to reblind wallet")

	// Combine the public and private shares
	walletShare, err := CombineShares(wallet.BlindedPublicShares, wallet.PrivateShares, wallet.Blinder)
	assert.NoError(t, err, "Failed to get existing wallet share")

	// Check if the balance is correctly represented
	assert.Equal(t, balance, walletShare.Balances[0], "Balance not correctly represented after reblinding")

	// Check if the order is correctly represented
	assert.Equal(t, order.BaseMint, walletShare.Orders[0].BaseMint, "Order BaseMint not correctly represented after reblinding")
	assert.Equal(t, order.QuoteMint, walletShare.Orders[0].QuoteMint, "Order QuoteMint not correctly represented after reblinding")
	assert.Equal(t, order.Amount, walletShare.Orders[0].Amount, "Order Amount not correctly represented after reblinding")
	assert.Equal(t, order.Side, walletShare.Orders[0].Side, "Order Side not correctly represented after reblinding")
	assert.Equal(t, order.WorstCasePrice, walletShare.Orders[0].WorstCasePrice, "Order WorstCasePrice not correctly represented after reblinding")
}
