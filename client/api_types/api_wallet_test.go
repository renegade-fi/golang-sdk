package api_types

import (
	"crypto/ecdsa"
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/stretchr/testify/assert"
	"renegade.fi/golang-sdk/wallet"
)

func TestApiWalletConversion(t *testing.T) {
	key, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	assert.NoError(t, err)
	originalWallet, err := wallet.NewEmptyWallet(key, 0 /* chainId */)
	assert.NoError(t, err)

	// Convert to API wallet
	apiWallet, err := new(ApiWallet).FromWallet(originalWallet)
	assert.NoError(t, err)

	// Convert back to wallet
	recoveredWallet, err := apiWallet.ToWallet()
	assert.NoError(t, err)

	// Check that the recovered wallet is the same as the original wallet
	assert.Equal(t, originalWallet, recoveredWallet)
}
