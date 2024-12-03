package api_types

// ApiToken is a token available on the exchange
type ApiToken struct {
	// The mint (erc20 address) of the token
	Address string `json:"address"`
	// The symbol of the token
	Symbol string `json:"symbol"`
}
