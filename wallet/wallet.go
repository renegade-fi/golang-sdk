package wallet

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/google/uuid"

	renegade_crypto "renegade.fi/golang-sdk/crypto"
)

const (
	// numScalarsWalletShare is the number of scalars in a wallet share
	numScalarsWalletShare = 70
	// maxBalances is the maximum number of balances in a wallet
	maxBalances = 10
	// maxOrders is the maximum number of orders in a wallet
	maxOrders = 4
)

// Scalar is a scalar field element from the bn254 curve
type Scalar fr.Element

// Add adds two scalars
func (s *Scalar) Add(other Scalar) Scalar {
	var result fr.Element
	fr1 := fr.Element(*s)
	fr2 := fr.Element(other)
	result.Add(&fr1, &fr2)

	return Scalar(result)
}

// Sub subtracts two scalars
func (s *Scalar) Sub(other Scalar) Scalar {
	var result fr.Element
	fr1 := fr.Element(*s)
	fr2 := fr.Element(other)
	result.Sub(&fr1, &fr2)

	return Scalar(result)
}

// WalletShare represents a secret share of a wallet, containing only the elements of a wallet that are stored on-chain
type WalletShare struct {
	// Balances are the balances of the wallet
	Balances [maxBalances]Balance
	// Orders are the orders of the wallet
	Orders [maxOrders]Order
	// Keys are the public keys of the wallet
	Keys PublicKeychain
	// MatchFee is the fee that the wallet pays to the cluster that matches its orders
	MatchFee FixedPoint
	// ManagingCluster is the public encryption key of the cluster that
	// receives fees for matching orders in the wallet
	ManagingCluster FeeEncryptionKey
	// Blinder is the additive blinder applied to all secret shares to make an adequately determined
	// algebraic system on the shares impossible, even when one knows the underlying value
	Blinder Scalar
}

// EmptyWalletShare creates a new wallet share with all zero values
func EmptyWalletShare(publicKeys PublicKeychain) (WalletShare, error) {
	// Create a slice of scalars with all zero values
	scalars := make([]Scalar, numScalarsWalletShare)
	for i := range scalars {
		scalars[i] = Scalar{}
	}

	// Deserialize a wallet share from the scalars
	share := WalletShare{}
	err := FromScalarsRecursive(&share, NewScalarIterator(scalars))
	if err != nil {
		return WalletShare{}, err
	}

	share.Keys = publicKeys
	return share, nil
}

// SplitIntoShares splits a wallet share into two shares using the given private shares and blinder
func (ws *WalletShare) SplitPublicPrivate(privateShares []Scalar, blinder Scalar) (WalletShare, WalletShare, error) {
	// Serialize the wallet share into scalars
	scalars, err := ToScalarsRecursive(ws)
	if err != nil {
		return WalletShare{}, WalletShare{}, err
	}

	// The shares should be the same length as the scalars
	if len(privateShares) != len(scalars) {
		return WalletShare{}, WalletShare{}, fmt.Errorf("private shares and scalars have different lengths")
	}

	// Subtract the private shares from the scalars to get the public shares
	// Then blind the public shares with the blinder
	publicShares := make([]Scalar, len(scalars))
	for i := range privateShares {
		publicShares[i] = scalars[i].Sub(privateShares[i])
		publicShares[i] = publicShares[i].Add(blinder)
	}

	// Blind the public shares additively with the blinder
	// Deserialize the shares from the scalars
	privateShare := WalletShare{}
	err = FromScalarsRecursive(&privateShare, NewScalarIterator(privateShares))
	if err != nil {
		return WalletShare{}, WalletShare{}, err
	}

	publicShare := WalletShare{}
	err = FromScalarsRecursive(&publicShare, NewScalarIterator(publicShares))
	if err != nil {
		return WalletShare{}, WalletShare{}, err
	}

	return privateShare, publicShare, nil
}

// Wallet is a wallet in the Renegade system
type Wallet struct {
	ID                  uuid.UUID
	Orders              []Order
	Balances            []Balance
	Keychain            *Keychain
	ManagingCluster     FeeEncryptionKey
	MatchFee            FixedPoint
	BlindedPublicShares WalletShare
	PrivateShares       WalletShare
	Blinder             Scalar
}

// NewEmptyWallet creates a new empty wallet
func NewEmptyWallet(privateKey *ecdsa.PrivateKey, chainId uint64) (*Wallet, error) {
	walletId, err := DeriveWalletID(privateKey, chainId)
	if err != nil {
		return nil, err
	}

	// Derive the keychain
	keychain, err := DeriveKeychain(privateKey, chainId)
	if err != nil {
		return nil, err
	}

	// Derive the seeds for secret shares and blinders
	blinderSeed, shareSeed, err := DeriveWalletSeeds(privateKey, chainId)
	if err != nil {
		return nil, err
	}

	// Setup a wallet with empty shares
	emptyShare, err := EmptyWalletShare(keychain.PublicKeys)
	if err != nil {
		return nil, err
	}

	// Reblind the wallet
	blinder, blinderPrivateShare := walletBlinderFromSeed(blinderSeed)
	privateShareScalars := walletSharesFromStream(shareSeed)
	privateShare, publicShare, err := emptyShare.SplitPublicPrivate(privateShareScalars, blinder)
	if err != nil {
		return nil, err
	}

	privateShare.Blinder = blinderPrivateShare
	publicShare.Blinder = blinder.Sub(blinderPrivateShare)

	return &Wallet{
		ID:       walletId,
		Orders:   make([]Order, 0),
		Balances: make([]Balance, 0),
		Keychain: keychain,
		// The managing relayer will set the following two fields
		ManagingCluster:     emptyShare.ManagingCluster,
		MatchFee:            emptyShare.MatchFee,
		BlindedPublicShares: publicShare,
		PrivateShares:       privateShare,
		Blinder:             blinder,
	}, nil
}

// walletSharesFromStream generates numScalarsWalletShare scalars from a CSPRNG seeded with the given scalar
func walletSharesFromStream(seed Scalar) []Scalar {
	// Create a poseidon CSPRNG from the seed and generate numScalarsWalletShare scalars
	csprng := renegade_crypto.NewPoseidonCSPRNG(fr.Element(seed))
	elts := csprng.NextN(numScalarsWalletShare)

	// Wrap the elements in a slice of scalars
	scalars := make([]Scalar, len(elts))
	for i, elt := range elts {
		scalars[i] = Scalar(elt)
	}

	return scalars
}

// walletBlinderFromSeed generates a wallet blinder and blinder private share from a CSPRNG seeded with the given scalar
func walletBlinderFromSeed(seed Scalar) (Scalar, Scalar) {
	csprng := renegade_crypto.NewPoseidonCSPRNG(fr.Element(seed))

	// Generate a blinder and blinder private share
	blinder := Scalar(csprng.Next())
	blinderPrivateShare := Scalar(csprng.Next())

	return blinder, blinderPrivateShare
}
