package wallet

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
