package wallet

import "fmt"

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

// IsZero returns true if the balance amount and fees are zero
func (b *Balance) IsZero() bool {
	return b.Amount.IsZero() && b.RelayerFeeBalance.IsZero() && b.ProtocolFeeBalance.IsZero()
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
