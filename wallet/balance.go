package wallet

import (
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// Balance is a balance in the Renegade system
type Balance struct {
	// Mint is the erc20 address of the balance's asset
	Mint Scalar
	// Amount is the amount of the balance
	Amount Scalar
	// RelayerFeeBalance is the balance due to the relayer in fees
	RelayerFeeBalance Scalar
	// ProtocolFeeBalance is the balance due to the protocol in fees
	ProtocolFeeBalance Scalar
}

// NewEmptyBalance creates a new balance with all zero values
func NewEmptyBalance() Balance {
	return Balance{
		Mint:               Scalar{},
		Amount:             Scalar{},
		RelayerFeeBalance:  Scalar{},
		ProtocolFeeBalance: Scalar{},
	}
}

// NewBalance creates a new balance with the given mint and amount
func NewBalance(mint Scalar, amount Scalar) Balance {
	return NewBalanceBuilder().
		WithMint(mint).
		WithAmount(amount).
		Build()
}

// IsZero returns true if the balance amount and fees are zero
func (b *Balance) IsZero() bool {
	return b.Amount.IsZero() && b.RelayerFeeBalance.IsZero() && b.ProtocolFeeBalance.IsZero()
}

// GetBalance gets the balance for a given mint
func (w *Wallet) GetBalance(mint string) (big.Int, error) {
	mintScalar, err := new(Scalar).FromHexString(mint)
	if err != nil {
		return big.Int{}, err
	}

	idx := w.findMatchingBalance(mintScalar)
	if idx == -1 {
		return big.Int{}, fmt.Errorf("balance not found for mint: %s", mint)
	}

	return *w.Balances[idx].Amount.ToBigInt(), nil
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

// RemoveBalance removes a balance from the wallet
func (w *Wallet) RemoveBalance(balance Balance) error {
	// Find the balance to remove
	idx := w.findMatchingBalance(balance.Mint)
	if idx == -1 {
		return fmt.Errorf("balance not found")
	}

	// Remove the balance
	amt1 := fr.Element(w.Balances[idx].Amount)
	amt2 := fr.Element(balance.Amount)

	if amt1.Cmp(&amt2) < 0 {
		return fmt.Errorf("balance is less than the amount to remove")
	}

	w.Balances[idx].Amount = w.Balances[idx].Amount.Sub(balance.Amount)
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

// BalanceBuilder is a builder for Balance
type BalanceBuilder struct {
	balance Balance
}

// NewBalanceBuilder creates a new BalanceBuilder
func NewBalanceBuilder() *BalanceBuilder {
	return &BalanceBuilder{balance: NewEmptyBalance()}
}

// WithMint sets the Mint
func (bb *BalanceBuilder) WithMint(mint Scalar) *BalanceBuilder {
	bb.balance.Mint = mint
	return bb
}

// WithMintHex sets the Mint from a hex string
func (bb *BalanceBuilder) WithMintHex(hexMint string) *BalanceBuilder {
	mint, err := new(Scalar).FromHexString(hexMint)
	if err != nil {
		panic(err)
	}

	bb.balance.Mint = mint
	return bb
}

// WithAmount sets the Amount
func (bb *BalanceBuilder) WithAmount(amount Scalar) *BalanceBuilder {
	bb.balance.Amount = amount
	return bb
}

// WithAmountBigInt sets the Amount from a big.Int
func (bb *BalanceBuilder) WithAmountBigInt(amount *big.Int) *BalanceBuilder {
	amountScalar := new(Scalar).FromBigInt(amount)
	bb.balance.Amount = amountScalar
	return bb
}

// WithRelayerFeeBalance sets the RelayerFeeBalance
func (bb *BalanceBuilder) WithRelayerFeeBalance(fee Scalar) *BalanceBuilder {
	bb.balance.RelayerFeeBalance = fee
	return bb
}

// WithProtocolFeeBalance sets the ProtocolFeeBalance
func (bb *BalanceBuilder) WithProtocolFeeBalance(fee Scalar) *BalanceBuilder {
	bb.balance.ProtocolFeeBalance = fee
	return bb
}

// Build returns the constructed Balance
func (bb *BalanceBuilder) Build() Balance {
	return bb.balance
}
