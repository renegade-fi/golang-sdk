package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"

	renegade_crypto "renegade.fi/golang-sdk/crypto"
)

const (
	// numScalarsWalletShare is the number of scalars in a wallet share
	numScalarsWalletShare = 70
	// maxBalances is the maximum number of balances in a wallet
	MaxBalances = 10
	// maxOrders is the maximum number of orders in a wallet
	MaxOrders = 4
)

// Scalar is a scalar field element from the bn254 curve
type Scalar fr.Element

// RandomScalar generates a random scalar
func RandomScalar() (Scalar, error) {
	var res fr.Element
	_, err := res.SetRandom()
	if err != nil {
		return Scalar{}, err
	}

	return Scalar(res), nil

}

// IsZero returns whether the scalar is zero
func (s *Scalar) IsZero() bool {
	return (*fr.Element)(s).IsZero()
}

// IsOne returns whether the scalar is one
func (s *Scalar) IsOne() bool {
	return (*fr.Element)(s).IsOne()
}

// Uint64 returns the scalar as a uint64
func (s *Scalar) Uint64() uint64 {
	return (*fr.Element)(s).Uint64()
}

// SetUint64 sets the scalar from a uint64
func (s *Scalar) SetUint64(val uint64) {
	(*fr.Element)(s).SetUint64(val)
}

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

// Bytes returns the bytes representation of the scalar in big-endian order
func (s *Scalar) Bytes() [fr.Bytes]byte {
	return (*fr.Element)(s).Bytes()
}

// FromBytes sets the scalar from a big-endian byte slice
func (s *Scalar) FromBytes(bytes [fr.Bytes]byte) {
	(*fr.Element)(s).SetBytes(bytes[:])
}

// HexString returns the hex string representation of the scalar
func (s *Scalar) ToHexString() string {
	bytes := s.Bytes()
	return hex.EncodeToString(bytes[:])
}

// FromHexString sets the scalar from a hex string
func (s *Scalar) FromHexString(hexString string) (Scalar, error) {
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		return Scalar{}, err
	}

	var fixedBytes [fr.Bytes]byte
	copy(fixedBytes[fr.Bytes-len(bytes):], bytes)
	s.FromBytes(fixedBytes)

	return *s, nil
}

// ToBigInt converts the scalar to a big.Int
func (s *Scalar) ToBigInt() *big.Int {
	var res big.Int
	(*fr.Element)(s).BigInt(&res)
	return &res
}

// FromBigInt sets the scalar from a big.Int
func (s *Scalar) FromBigInt(i *big.Int) Scalar {
	(*fr.Element)(s).SetBigInt(i)
	return *s
}

// WalletSecrets contains the information about a wallet necessary to recover it
type WalletSecrets struct {
	// Id is the UUID of the wallet
	Id uuid.UUID
	// Keychain is the keychain used to manage the wallet
	Keychain *Keychain
	// BlinderSeed is the seed of the CSPRNG used to generate blinders and blinder shares
	BlinderSeed Scalar
	// ShareSeed is the seed of the CSPRNG used to generate wallet secret shares
	ShareSeed Scalar
}

// DeriveWalletSecrets derives the wallet secrets from the given Ethereum private key
func DeriveWalletSecrets(ethKey *ecdsa.PrivateKey, chainId uint64) (*WalletSecrets, error) {
	walletId, err := DeriveWalletID(ethKey, chainId)
	if err != nil {
		return nil, err
	}

	keychain, err := DeriveKeychain(ethKey, chainId)
	if err != nil {
		return nil, err
	}

	blinderSeed, shareSeed, err := DeriveWalletSeeds(ethKey, chainId)
	if err != nil {
		return nil, err
	}

	return &WalletSecrets{
		Id:          walletId,
		Keychain:    keychain,
		BlinderSeed: blinderSeed,
		ShareSeed:   shareSeed,
	}, nil
}

// WalletShare represents a secret share of a wallet, containing only the elements of a wallet that are stored on-chain
type WalletShare struct {
	// Balances are the balances of the wallet
	Balances [MaxBalances]Balance
	// Orders are the orders of the wallet
	Orders [MaxOrders]Order
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

// CombineShares combines two wallet shares into a single wallet share
func CombineShares(publicShare WalletShare, privateShare WalletShare, blinder Scalar) (WalletShare, error) {
	publicScalars, err := ToScalarsRecursive(&publicShare)
	if err != nil {
		return WalletShare{}, err
	}

	privateScalars, err := ToScalarsRecursive(&privateShare)
	if err != nil {
		return WalletShare{}, err
	}

	combinedScalars := make([]Scalar, len(publicScalars))
	for i := range publicScalars {
		tmp := publicScalars[i].Add(privateScalars[i])
		combinedScalars[i] = tmp.Sub(blinder)
	}

	combined := WalletShare{}
	err = FromScalarsRecursive(&combined, NewScalarIterator(combinedScalars))
	if err != nil {
		return WalletShare{}, err
	}

	return combined, nil
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
	secrets, err := DeriveWalletSecrets(privateKey, chainId)
	if err != nil {
		return nil, err
	}

	return NewEmptyWalletFromSecrets(secrets)
}

// NewEmptyWalletFromSecrets creates a new wallet from the given wallet secrets
func NewEmptyWalletFromSecrets(secrets *WalletSecrets) (*Wallet, error) {
	walletId := secrets.Id
	keychain := secrets.Keychain

	// Setup a wallet with empty shares
	emptyShare, err := EmptyWalletShare(keychain.PublicKeys)
	if err != nil {
		return nil, err
	}

	// Reblind the wallet
	blinder, blinderPrivateShare := walletBlinderFromSeed(secrets.BlinderSeed)
	privateShareScalars := walletSharesFromStream(secrets.ShareSeed)

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

// GetShareCommitment returns a Poseidon hash commitment of the wallet's shares
func (w *Wallet) GetShareCommitment() (Scalar, error) {
	privateCommitment, err := w.GetPrivateShareCommitment()
	if err != nil {
		return Scalar{}, err
	}

	// Hash in the public shares
	publicShares, err := ToScalarsRecursive(&w.BlindedPublicShares)
	if err != nil {
		return Scalar{}, err
	}

	// Create a hash input that is the privateCommitment concatenated with the publicShares
	hashInput := append([]Scalar{privateCommitment}, publicShares...)
	return HashScalars(hashInput), nil
}

// GetPrivateShareCommitment returns a Poseidon hash commitment of the wallet's private share
func (w *Wallet) GetPrivateShareCommitment() (Scalar, error) {
	privateShares, err := ToScalarsRecursive(&w.PrivateShares)
	if err != nil {
		return Scalar{}, err
	}

	return HashScalars(privateShares), nil
}

// SignCommitment signs the given commitment using the private root key
func (w *Wallet) SignCommitment(commitment Scalar) ([]byte, error) {
	privateRootKey := w.Keychain.SkRoot()
	signKey := ecdsa.PrivateKey(*privateRootKey)

	commBytes := commitment.ToBigInt().Bytes()
	digest := crypto.Keccak256(commBytes)
	sig, err := crypto.Sign(digest, &signKey)
	if err != nil {
		return nil, err
	}

	return sig, nil
}

// ReblindWallet reblinds the wallet, sampling new secret shares and blinders from the CSPRNGs
func (w *Wallet) Reblind() error {
	privateShares, err := ToScalarsRecursive(&w.PrivateShares)
	if err != nil {
		return err
	}

	// Sample new private shares from the CSPRNG, using the last existing private share as the seed
	// And sample a new blinder using the old blinder private share as the seed
	newPrivateShares := walletSharesFromStream(privateShares[len(privateShares)-1])
	newBlinder, newBlinderPrivateShare := walletBlinderFromSeed(w.PrivateShares.Blinder)

	// Split the new private shares into a private and public share
	existingShare, err := w.getExistingWalletShare()
	if err != nil {
		return err
	}

	publicShare, privateShare, err := existingShare.SplitPublicPrivate(newPrivateShares, newBlinder)
	privateShare.Blinder = newBlinderPrivateShare
	if err != nil {
		return err
	}

	w.PrivateShares = privateShare
	w.BlindedPublicShares = publicShare
	w.Blinder = newBlinder
	return nil
}

// getExistingWalletShare combines the existing public and private shares into a single wallet share
func (w *Wallet) getExistingWalletShare() (WalletShare, error) {
	ws := new(WalletShare)

	// Deep copy Balances
	for i, balance := range w.Balances {
		if i >= MaxBalances {
			break
		}
		ws.Balances[i] = balance
	}

	// Deep copy Orders
	for i, order := range w.Orders {
		if i >= MaxOrders {
			break
		}
		ws.Orders[i] = order
	}

	// These are likely value types, so simple assignment should be fine
	ws.Keys = w.Keychain.PublicKeys
	ws.MatchFee = w.MatchFee
	ws.ManagingCluster = w.ManagingCluster
	ws.Blinder = w.Blinder

	return *ws, nil
}

// NewOrder appends an order to the wallet
func (w *Wallet) NewOrder(order Order) error {
	// Find the first order that may be replaced
	if idx := w.findReplaceableOrder(); idx != -1 {
		w.Orders[idx] = order
	} else if len(w.Orders) < MaxOrders {
		w.Orders = append(w.Orders, order)
	} else {
		return fmt.Errorf("wallet already has the maximum number of orders")
	}

	return nil
}

// findReplaceableOrder finds the first order that may be replaced by the new order
// Returns the index of the order to replace, or -1 if no order may be replaced
func (w *Wallet) findReplaceableOrder() int {
	for i, existingOrder := range w.Orders {
		if existingOrder.IsZero() {
			return i
		}
	}

	return -1
}

// AddBalance appends a balance to the wallet
func (w *Wallet) AddBalance(balance Balance) error {
	// Find an existing balance for the mint if one exists
	if idx := w.findMatchingBalance(balance.Mint); idx != -1 {
		w.Balances[idx].Amount = w.Balances[idx].Amount.Add(balance.Amount)
		return nil
	}

	// If the balance is not found, try to append one
	if idx := w.findReplaceableBalance(); idx != -1 {
		w.Balances[idx] = balance
	} else if len(w.Balances) < MaxBalances {
		w.Balances = append(w.Balances, balance)
	} else {
		return fmt.Errorf("wallet already has the maximum number of balances")
	}

	return nil
}

// findMatchingBalance finds the index of a balance with the given mint, or -1 if no balance has the given mint
func (w *Wallet) findMatchingBalance(mint Scalar) int {
	for i, balance := range w.Balances {
		if balance.Mint == mint {
			return i
		}
	}

	return -1
}

// findReplaceableBalance finds the first balance that may be replaced, returning the index of the balance, or -1 if no balance may be replaced
func (w *Wallet) findReplaceableBalance() int {
	for i, balance := range w.Balances {
		if balance.IsZero() {
			return i
		}
	}

	return -1
}
