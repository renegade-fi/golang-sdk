package api_types

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/google/uuid"

	"renegade.fi/golang-sdk/wallet"
)

// The number of u32 limbs in the serialized form of a secret share
const secretShareLimbCount = 8 // 256 bits

type Amount big.Int

func (a *Amount) String() string {
	return (*big.Int)(a).String()
}

func (a *Amount) MarshalJSON() ([]byte, error) {
	return []byte(a.String()), nil
}

func (a *Amount) SetString(s string, base int) error {
	i, ok := new(big.Int).SetString(s, base)
	if !ok {
		return fmt.Errorf("invalid number: %s", s)
	}
	*a = Amount(*i)
	return nil
}

func (a *Amount) UnmarshalJSON(b []byte) error {
	return a.SetString(string(b), 10)
}

// orderSideFromScalar converts a wallet.Scalar to an order side
func orderSideFromScalar(s wallet.Scalar) (string, error) {
	if s.IsZero() {
		return "Buy", nil
	} else if s.IsOne() {
		return "Sell", nil
	}

	return "", fmt.Errorf("invalid order side: %s", s.ToHexString())
}

// orderSideToScalar converts an order side to a wallet.Scalar
func orderSideToScalar(side string) (wallet.Scalar, error) {
	if side == "Buy" {
		return wallet.Scalar(fr.NewElement(0)), nil
	} else if side == "Sell" {
		return wallet.Scalar(fr.NewElement(1)), nil
	}

	return wallet.Scalar{}, fmt.Errorf("invalid order side: %s", side)
}

// ApiOrder is an order in a Renegade wallet
type ApiOrder struct {
	// The id of the order
	Id uuid.UUID `json:"id"`
	// The mint (erc20 address) of the base asset
	// As a hex string
	BaseMint string `json:"base_mint"`
	// The mint (erc20 address) of the quote asset
	// As a hex string
	QuoteMint string `json:"quote_mint"`
	// The amount of the base asset to buy/sell
	Amount Amount `json:"amount"`
	// The side of the order
	Side string `json:"side"`
	// The worst case price to execute the order at
	// The serialized form of this is the `Scalar` representation of the fixed point,
	// i.e. if a fixed point value represents `r`, this value is `floor(r << PRECISION)`
	WorstCasePrice string `json:"worst_case_price"`
}

// FromOrder converts a wallet.Order to an ApiOrder
func (a *ApiOrder) FromOrder(o *wallet.Order) (*ApiOrder, error) {
	a.Id = o.Id
	a.BaseMint = o.BaseMint.ToHexString()
	a.QuoteMint = o.QuoteMint.ToHexString()
	a.Amount = Amount(*o.Amount.ToBigInt())
	side, err := orderSideFromScalar(o.Side)
	if err != nil {
		return nil, err
	}

	a.Side = side
	a.WorstCasePrice = o.WorstCasePrice.ToReprDecimalString()

	return a, nil
}

// ToOrder converts an ApiOrder to a wallet.Order
func (a *ApiOrder) ToOrder(o *wallet.Order) error {
	o.Id = a.Id
	if _, err := o.BaseMint.FromHexString(a.BaseMint); err != nil {
		return err
	}
	if _, err := o.QuoteMint.FromHexString(a.QuoteMint); err != nil {
		return err
	}

	o.Amount = new(wallet.Scalar).FromBigInt((*big.Int)(&a.Amount))
	side, err := orderSideToScalar(a.Side)
	if err != nil {
		return err
	}

	o.Side = side
	return nil
}

// ApiBalance is a balance in a Renegade wallet
type ApiBalance struct {
	// The mint (erc20 address) of the asset
	Mint string `json:"mint"`
	// The amount of the asset
	Amount Amount `json:"amount"`
	// The amount of this balance owed to the managing relayer cluster
	RelayerFeeBalance Amount `json:"relayer_fee_balance"`
	// The amount of this balance owed to the protocol
	ProtocolFeeBalance Amount `json:"protocol_fee_balance"`
}

// FromBalance converts a wallet.Balance to an ApiBalance
func (a *ApiBalance) FromBalance(b *wallet.Balance) error {
	a.Mint = b.Mint.ToHexString()
	a.Amount = Amount(*b.Amount.ToBigInt())
	a.RelayerFeeBalance = Amount(*b.RelayerFeeBalance.ToBigInt())
	a.ProtocolFeeBalance = Amount(*b.ProtocolFeeBalance.ToBigInt())

	return nil
}

// ToBalance converts an ApiBalance to a wallet.Balance
func (a *ApiBalance) ToBalance(b *wallet.Balance) error {
	if _, err := b.Mint.FromHexString(a.Mint); err != nil {
		return err
	}
	b.Amount = new(wallet.Scalar).FromBigInt((*big.Int)(&a.Amount))
	b.RelayerFeeBalance = new(wallet.Scalar).FromBigInt((*big.Int)(&a.RelayerFeeBalance))
	b.ProtocolFeeBalance = new(wallet.Scalar).FromBigInt((*big.Int)(&a.ProtocolFeeBalance))

	return nil
}

// ApiPublicKeychain is a public keychain in the Renegade system
type ApiPublicKeychain struct {
	// The public root key of the wallet
	// As a hex string
	PkRoot string `json:"pk_root"`
	// The public match key of the wallet
	// As a hex string
	PkMatch string `json:"pk_match"`
}

func (a *ApiPublicKeychain) FromPublicKeychain(pk *wallet.PublicKeychain) error {
	a.PkRoot = pk.PkRoot.ToHexString()
	a.PkMatch = pk.PkMatch.ToHexString()

	return nil
}

func (a *ApiPublicKeychain) ToPublicKeychain() (*wallet.PublicKeychain, error) {
	pkRoot, err := new(wallet.PublicSigningKey).FromHexString(a.PkRoot)
	if err != nil {
		return nil, err
	}
	pkMatch, err := new(wallet.Scalar).FromHexString(a.PkMatch)
	if err != nil {
		return nil, err
	}

	return &wallet.PublicKeychain{
		PkRoot:  pkRoot,
		PkMatch: pkMatch,
	}, nil
}

// ApiPrivateKeychain represents a private keychain for the API wallet
type ApiPrivateKeychain struct {
	// The private root key of the wallet
	// As a hex string, optional
	SkRoot *string `json:"sk_root,omitempty"`
	// The private match key of the wallet
	// As a hex string
	SkMatch string `json:"sk_match"`
	// The symmetric key of the wallet
	// As a hex string
	SymmetricKey string `json:"symmetric_key"`
}

// FromPrivateKeychain converts a wallet.PrivateKeychain to an ApiPrivateKeychain
func (a *ApiPrivateKeychain) FromPrivateKeychain(pk *wallet.PrivateKeychain) (*ApiPrivateKeychain, error) {
	if pk.SkRoot != nil {
		skRootHex := pk.SkRoot.ToHexString()
		a.SkRoot = &skRootHex
	}

	a.SkMatch = pk.SkMatch.ToHexString()
	a.SymmetricKey = pk.SymmetricKey.ToHexString()

	return a, nil
}

// ToPrivateKeychain converts an ApiPrivateKeychain to a wallet.PrivateKeychain
func (a *ApiPrivateKeychain) ToPrivateKeychain() (*wallet.PrivateKeychain, error) {
	skRoot, err := new(wallet.PrivateSigningKey).FromHexString(*a.SkRoot)
	if err != nil {
		return nil, err
	}
	skMatch, err := new(wallet.Scalar).FromHexString(a.SkMatch)
	if err != nil {
		return nil, err
	}
	symmetricKey, err := new(wallet.HmacKey).FromHexString(a.SymmetricKey)
	if err != nil {
		return nil, err
	}

	return &wallet.PrivateKeychain{
		SkRoot:       &skRoot,
		SkMatch:      skMatch,
		SymmetricKey: symmetricKey,
	}, nil
}

// ApiKeychain represents a keychain API type that maintains all keys as hex strings
type ApiKeychain struct {
	// The public keychain
	PublicKeys ApiPublicKeychain `json:"public_keys"`
	// The private keychain
	PrivateKeys ApiPrivateKeychain `json:"private_keys"`
	// The nonce of the keychain
	Nonce uint64 `json:"nonce"`
}

// FromKeychain converts a wallet.Keychain to an ApiKeychain
func (a *ApiKeychain) FromKeychain(k *wallet.Keychain) (*ApiKeychain, error) {
	if err := a.PublicKeys.FromPublicKeychain(&k.PublicKeys); err != nil {
		return nil, err
	}
	if _, err := a.PrivateKeys.FromPrivateKeychain(&k.PrivateKeys); err != nil {
		return nil, err
	}
	a.Nonce = k.PublicKeys.Nonce.Uint64()
	return a, nil
}

// ToKeychain converts an ApiKeychain to a wallet.Keychain
func (a *ApiKeychain) ToKeychain() (*wallet.Keychain, error) {
	publicKeys, err := a.PublicKeys.ToPublicKeychain()
	if err != nil {
		return nil, err
	}
	publicKeys.Nonce.SetUint64(a.Nonce)

	privateKeys, err := a.PrivateKeys.ToPrivateKeychain()
	if err != nil {
		return nil, err
	}
	privateKeys.SkRoot.PublicKey = ecdsa.PublicKey(publicKeys.PkRoot)

	return &wallet.Keychain{
		PublicKeys:  *publicKeys,
		PrivateKeys: *privateKeys,
	}, nil
}

// ApiWallet is a wallet in the Renegade system
type ApiWallet struct {
	// Identifier
	Id uuid.UUID `json:"id"`
	// The orders maintained by this wallet
	Orders []ApiOrder `json:"orders"`
	// The balances maintained by the wallet to cover orders
	Balances []ApiBalance `json:"balances"`
	// The keys that authenticate wallet access
	KeyChain ApiKeychain `json:"key_chain"`
	// The managing cluster's public key
	// The public encryption key of the cluster that may collect relayer fees
	// on this wallet
	ManagingCluster string `json:"managing_cluster"`
	// The take rate at which the managing cluster may collect relayer fees on
	// a match
	MatchFee string `json:"match_fee"`
	// The public secret shares of the wallet
	BlindedPublicShares [][secretShareLimbCount]uint32 `json:"blinded_public_shares"`
	// The private secret shares of the wallet
	PrivateShares [][secretShareLimbCount]uint32 `json:"private_shares"`
	// The wallet blinder, used to blind wallet secret shares
	Blinder [secretShareLimbCount]uint32 `json:"blinder"`
}

func (a *ApiWallet) FromWallet(w *wallet.Wallet) (*ApiWallet, error) {
	a.Id = w.ID

	// Convert orders
	a.Orders = make([]ApiOrder, len(w.Orders))
	for _, order := range w.Orders {
		var apiOrder ApiOrder
		if _, err := apiOrder.FromOrder(&order); err != nil {
			return nil, err
		}
		a.Orders = append(a.Orders, apiOrder)
	}

	// Convert balances
	a.Balances = make([]ApiBalance, len(w.Balances))
	for _, balance := range w.Balances {
		var apiBalance ApiBalance
		if err := apiBalance.FromBalance(&balance); err != nil {
			return nil, err
		}
		a.Balances = append(a.Balances, apiBalance)
	}

	// Convert keychain, managing cluster, and match fee
	if _, err := a.KeyChain.FromKeychain(w.Keychain); err != nil {
		return nil, err
	}
	a.ManagingCluster = w.ManagingCluster.ToHexString()
	a.MatchFee = w.MatchFee.ToReprDecimalString()

	// Convert the public shares
	publicShares, err := wallet.ToScalarsRecursive(&w.BlindedPublicShares)
	if err != nil {
		return nil, err
	}

	for _, share := range publicShares {
		a.BlindedPublicShares = append(a.BlindedPublicShares, ScalarToUintLimbs(share))
	}

	// Convert the private shares
	privateShares, err := wallet.ToScalarsRecursive(&w.PrivateShares)
	if err != nil {
		return nil, err
	}

	for _, share := range privateShares {
		a.PrivateShares = append(a.PrivateShares, ScalarToUintLimbs(share))
	}

	// Convert the blinder
	a.Blinder = ScalarToUintLimbs(w.Blinder)
	return a, nil
}

// ToWallet converts an ApiWallet to a Wallet
func (a *ApiWallet) ToWallet() (*wallet.Wallet, error) {
	w := &wallet.Wallet{}

	// Convert ID
	w.ID = a.Id

	// Convert orders
	w.Orders = make([]wallet.Order, len(a.Orders))
	for i, apiOrder := range a.Orders {
		if err := apiOrder.ToOrder(&w.Orders[i]); err != nil {
			return nil, err
		}
	}

	// Convert balances
	w.Balances = make([]wallet.Balance, len(a.Balances))
	for i, apiBalance := range a.Balances {
		if err := apiBalance.ToBalance(&w.Balances[i]); err != nil {
			return nil, err
		}
	}

	// Convert keychain, managing cluster, and match fee
	keychain, err := a.KeyChain.ToKeychain()
	if err != nil {
		return nil, err
	}
	w.Keychain = keychain
	if err := w.ManagingCluster.FromHexString(a.ManagingCluster); err != nil {
		return nil, err
	}
	if _, err := w.MatchFee.FromReprDecimalString(a.MatchFee); err != nil {
		return nil, err
	}

	// Convert the public shares
	publicShares := make([]wallet.Scalar, len(a.BlindedPublicShares))
	for i, limbs := range a.BlindedPublicShares {
		publicShares[i] = ScalarFromUintLimbs(limbs)
	}
	if err := wallet.FromScalarsRecursive(&w.BlindedPublicShares, wallet.NewScalarIterator(publicShares)); err != nil {
		return nil, err
	}

	// Convert the private shares
	privateShares := make([]wallet.Scalar, len(a.PrivateShares))
	for i, limbs := range a.PrivateShares {
		privateShares[i] = ScalarFromUintLimbs(limbs)
	}
	if err := wallet.FromScalarsRecursive(&w.PrivateShares, wallet.NewScalarIterator(privateShares)); err != nil {
		return nil, err
	}

	// Convert the blinder
	w.Blinder = ScalarFromUintLimbs(a.Blinder)
	return w, nil
}
