package external_match_client //nolint:revive

import (
	"fmt"
	"math/big"
	"strconv"

	geth_common "github.com/ethereum/go-ethereum/common"

	"github.com/renegade-fi/golang-sdk/client/api_types"
)

// ---------------------
// | Order Conversions |
// ---------------------

// v1OrderToV2 converts a v1 ApiExternalOrder to a v2 ApiExternalOrderV2
func v1OrderToV2(order *api_types.ApiExternalOrder) api_types.ApiExternalOrderV2 {
	switch order.Side {
	case "Buy":
		// Buy: input=quote, output=base
		var inputAmount, outputAmount api_types.StringAmount
		var useExactOutput bool
		quoteAmt := (*big.Int)(&order.QuoteAmount)
		baseAmt := (*big.Int)(&order.BaseAmount)
		exactBaseOut := (*big.Int)(&order.ExactBaseAmountOutput)
		exactQuoteOut := (*big.Int)(&order.ExactQuoteAmountOutput)

		if quoteAmt.Sign() != 0 {
			inputAmount = api_types.NewStringAmountFromBigInt(quoteAmt)
			outputAmount = api_types.NewStringAmount(0)
			useExactOutput = false
		} else if baseAmt.Sign() != 0 {
			inputAmount = api_types.NewStringAmount(0)
			outputAmount = api_types.NewStringAmountFromBigInt(baseAmt)
			useExactOutput = false
		} else if exactBaseOut.Sign() != 0 {
			inputAmount = api_types.NewStringAmount(0)
			outputAmount = api_types.NewStringAmountFromBigInt(exactBaseOut)
			useExactOutput = true
		} else {
			inputAmount = api_types.NewStringAmountFromBigInt(exactQuoteOut)
			outputAmount = api_types.NewStringAmount(0)
			useExactOutput = true
		}

		return api_types.ApiExternalOrderV2{
			InputMint:            order.QuoteMint,
			OutputMint:           order.BaseMint,
			InputAmount:          inputAmount,
			OutputAmount:         outputAmount,
			UseExactOutputAmount: useExactOutput,
			MinFillSize:          api_types.NewStringAmountFromBigInt((*big.Int)(&order.MinFillSize)),
		}

	default: // Sell
		// Sell: input=base, output=quote
		var inputAmount, outputAmount api_types.StringAmount
		var useExactOutput bool
		baseAmt := (*big.Int)(&order.BaseAmount)
		quoteAmt := (*big.Int)(&order.QuoteAmount)
		exactQuoteOut := (*big.Int)(&order.ExactQuoteAmountOutput)
		exactBaseOut := (*big.Int)(&order.ExactBaseAmountOutput)

		if baseAmt.Sign() != 0 {
			inputAmount = api_types.NewStringAmountFromBigInt(baseAmt)
			outputAmount = api_types.NewStringAmount(0)
			useExactOutput = false
		} else if quoteAmt.Sign() != 0 {
			inputAmount = api_types.NewStringAmount(0)
			outputAmount = api_types.NewStringAmountFromBigInt(quoteAmt)
			useExactOutput = false
		} else if exactQuoteOut.Sign() != 0 {
			inputAmount = api_types.NewStringAmount(0)
			outputAmount = api_types.NewStringAmountFromBigInt(exactQuoteOut)
			useExactOutput = true
		} else {
			inputAmount = api_types.NewStringAmountFromBigInt(exactBaseOut)
			outputAmount = api_types.NewStringAmount(0)
			useExactOutput = true
		}

		return api_types.ApiExternalOrderV2{
			InputMint:            order.BaseMint,
			OutputMint:           order.QuoteMint,
			InputAmount:          inputAmount,
			OutputAmount:         outputAmount,
			UseExactOutputAmount: useExactOutput,
			MinFillSize:          api_types.NewStringAmountFromBigInt((*big.Int)(&order.MinFillSize)),
		}
	}
}

// ---------------------
// | Quote Conversions |
// ---------------------

// invertPriceString inverts a price string (computes 1/price)
func invertPriceString(priceStr string) (string, error) {
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse price: %w", err)
	}
	if price == 0.0 {
		return "0", nil
	}
	inverted := 1.0 / price
	return strconv.FormatFloat(inverted, 'g', -1, 64), nil
}

// v2QuoteToV1 converts a v2 SignedExternalQuoteV2 to a v1 ApiSignedQuote
func v2QuoteToV1(
	v2Quote *SignedExternalQuoteV2,
	originalOrder *api_types.ApiExternalOrder,
) (*api_types.ApiSignedQuote, error) {
	q := &v2Quote.Quote
	direction := originalOrder.Side

	// Map v2 input/output to v1 base/quote based on direction
	var quoteMint, baseMint string
	var quoteAmount, baseAmount api_types.Amount
	switch direction {
	case "Buy":
		quoteMint = q.MatchResult.InputMint
		baseMint = q.MatchResult.OutputMint
		quoteAmount = api_types.Amount(*q.MatchResult.InputAmount.ToBigInt())
		baseAmount = api_types.Amount(*q.MatchResult.OutputAmount.ToBigInt())
	default: // Sell
		quoteMint = q.MatchResult.OutputMint
		baseMint = q.MatchResult.InputMint
		quoteAmount = api_types.Amount(*q.MatchResult.OutputAmount.ToBigInt())
		baseAmount = api_types.Amount(*q.MatchResult.InputAmount.ToBigInt())
	}

	v1MatchResult := api_types.ApiExternalMatchResult{
		QuoteMint:   quoteMint,
		BaseMint:    baseMint,
		QuoteAmount: quoteAmount,
		BaseAmount:  baseAmount,
		Direction:   direction,
	}

	// Convert price from v2's output/input to v1's quote/base
	var v1Price api_types.TimestampedPrice
	switch direction {
	case "Buy":
		invertedPrice, err := invertPriceString(q.Price.Price)
		if err != nil {
			return nil, err
		}
		v1Price = api_types.TimestampedPrice{
			Price:     invertedPrice,
			Timestamp: q.Price.Timestamp,
		}
	default: // Sell
		v1Price = q.Price
	}

	// Convert v2 send/receive (StringAmount) to v1 (Amount)
	v1Send := api_types.ApiExternalAssetTransfer{
		Mint:   q.Send.Mint,
		Amount: api_types.Amount(*q.Send.Amount.ToBigInt()),
	}
	v1Receive := api_types.ApiExternalAssetTransfer{
		Mint:   q.Receive.Mint,
		Amount: api_types.Amount(*q.Receive.Amount.ToBigInt()),
	}

	// Convert v2 fees (StringAmount) to v1 (Amount)
	v1Fees := api_types.ApiFee{
		RelayerFee:  api_types.Amount(*q.Fees.RelayerFee.ToBigInt()),
		ProtocolFee: api_types.Amount(*q.Fees.ProtocolFee.ToBigInt()),
	}

	v1Quote := api_types.ApiExternalQuote{
		Order:       *originalOrder,
		MatchResult: v1MatchResult,
		Fees:        v1Fees,
		Send:        v1Send,
		Receive:     v1Receive,
		Price:       v1Price,
		Timestamp:   q.Timestamp,
	}

	// Build the gas sponsorship info wrapper
	var gasSponsorshipInfo *api_types.ApiSignedGasSponsorshipInfo
	if v2Quote.GasSponsorshipInfo != nil {
		gasSponsorshipInfo = &api_types.ApiSignedGasSponsorshipInfo{
			GasSponsorshipInfo: *v2Quote.GasSponsorshipInfo,
		}
	}

	// Store the original v2 ApiSignedQuote for round-tripping
	innerV2 := &api_types.ApiSignedQuoteV2{
		Quote:     v2Quote.Quote,
		Signature: v2Quote.Signature,
		Deadline:  v2Quote.Deadline,
	}

	return api_types.NewApiSignedQuote(
		v1Quote,
		v2Quote.Signature,
		v2Quote.Deadline,
		gasSponsorshipInfo,
		innerV2,
	), nil
}

// v1QuoteToV2 extracts the v2 SignedExternalQuoteV2 from a v1 ApiSignedQuote
// for use in the assemble flow
func v1QuoteToV2(v1 *api_types.ApiSignedQuote) (*SignedExternalQuoteV2, error) {
	innerV2 := v1.InnerV2Quote()
	if innerV2 == nil {
		return nil, fmt.Errorf("ApiSignedQuote has no inner v2 quote for round-tripping")
	}

	var gasInfo *api_types.ApiGasSponsorshipInfo
	if v1.GasSponsorshipInfo != nil {
		gasInfo = &v1.GasSponsorshipInfo.GasSponsorshipInfo
	}

	return &SignedExternalQuoteV2{
		Quote:              innerV2.Quote,
		Signature:          innerV2.Signature,
		Deadline:           innerV2.Deadline,
		GasSponsorshipInfo: gasInfo,
	}, nil
}

// -------------------------
// | Response Conversions  |
// -------------------------

// decodeInputAmountFromCalldata reads the input amount from settlement tx calldata bytes 4-36
func decodeInputAmountFromCalldata(tx *api_types.ApiSettlementTransactionV2) (*big.Int, error) {
	// The input field is a hex string; convert to raw bytes
	dataBytes := geth_common.FromHex(tx.Input)
	end := inputAmountOffset + amountCalldataLength

	if len(dataBytes) < end {
		return nil, fmt.Errorf("invalid calldata: too short (len=%d, need=%d)", len(dataBytes), end)
	}

	inputSlice := dataBytes[inputAmountOffset:end]
	return new(big.Int).SetBytes(inputSlice), nil
}

// v2ResponseToV1NonMalleable converts a v2 ExternalMatchResponseV2 to a v1 ExternalMatchBundle
func v2ResponseToV1NonMalleable(
	resp *api_types.ExternalMatchResponseV2,
	direction string,
) (*ExternalMatchBundle, error) {
	matchResult := &resp.MatchBundle.MatchResult
	priceFp := &matchResult.PriceFp

	// Decode the input amount from the calldata
	inputAmount, err := decodeInputAmountFromCalldata(&resp.MatchBundle.SettlementTx)
	if err != nil {
		return nil, fmt.Errorf("failed to decode input amount: %w", err)
	}

	// Compute output amount from price
	outputAmount := priceFp.FloorMulInt(inputAmount)

	// Map v2 input/output to v1 base/quote based on direction
	var quoteMint, baseMint string
	var quoteAmount, baseAmount *big.Int
	switch direction {
	case "Buy":
		quoteMint = matchResult.InputMint
		baseMint = matchResult.OutputMint
		quoteAmount = inputAmount
		baseAmount = outputAmount
	default: // Sell
		quoteMint = matchResult.OutputMint
		baseMint = matchResult.InputMint
		quoteAmount = outputAmount
		baseAmount = inputAmount
	}

	v1MatchResult := &api_types.ApiExternalMatchResult{
		QuoteMint:   quoteMint,
		BaseMint:    baseMint,
		QuoteAmount: api_types.Amount(*quoteAmount),
		BaseAmount:  api_types.Amount(*baseAmount),
		Direction:   direction,
	}

	// Compute fees from fee_rates and output amount
	totalFeeRate := resp.MatchBundle.FeeRates.Total()
	totalFee := totalFeeRate.FloorMulInt(outputAmount)
	relayerFee := resp.MatchBundle.FeeRates.RelayerFeeRate.FloorMulInt(outputAmount)
	protocolFee := new(big.Int).Sub(totalFee, relayerFee)

	v1Fees := &api_types.ApiFee{
		RelayerFee:  api_types.Amount(*relayerFee),
		ProtocolFee: api_types.Amount(*protocolFee),
	}

	// Compute receive/send: fees are subtracted from receive (output) side
	var receive, send *api_types.ApiExternalAssetTransfer
	switch direction {
	case "Buy":
		recvAmount := new(big.Int).Sub(baseAmount, totalFee)
		receive = &api_types.ApiExternalAssetTransfer{
			Mint:   baseMint,
			Amount: api_types.Amount(*recvAmount),
		}
		send = &api_types.ApiExternalAssetTransfer{
			Mint:   quoteMint,
			Amount: api_types.Amount(*quoteAmount),
		}
	default: // Sell
		recvAmount := new(big.Int).Sub(quoteAmount, totalFee)
		receive = &api_types.ApiExternalAssetTransfer{
			Mint:   quoteMint,
			Amount: api_types.Amount(*recvAmount),
		}
		send = &api_types.ApiExternalAssetTransfer{
			Mint:   baseMint,
			Amount: api_types.Amount(*baseAmount),
		}
	}

	gasSponsored := resp.GasSponsorshipInfo != nil

	return &ExternalMatchBundle{
		MatchResult:        v1MatchResult,
		Fees:               v1Fees,
		Receive:            receive,
		Send:               send,
		SettlementTx:       toSettlementTransactionV2(&resp.MatchBundle.SettlementTx),
		GasSponsored:       gasSponsored,
		GasSponsorshipInfo: resp.GasSponsorshipInfo,
	}, nil
}

// --------------------------
// | Options Conversions    |
// --------------------------

// v1AssembleOptionsToV2 converts v1 AssembleExternalMatchOptions to v2 AssembleExternalMatchOptionsV2
func v1AssembleOptionsToV2(
	opts *AssembleExternalMatchOptions,
	originalOrder *api_types.ApiExternalOrder,
) *AssembleExternalMatchOptionsV2 {
	v2Opts := &AssembleExternalMatchOptionsV2{
		DoGasEstimation: opts.DoGasEstimation,
		ReceiverAddress: opts.ReceiverAddress,
	}
	if opts.UpdatedOrder != nil {
		v2Order := v1OrderToV2(opts.UpdatedOrder)
		v2Opts.UpdatedOrder = &v2Order
	} else if originalOrder != nil {
		// No updated order; use the original order converted to v2
		// (not needed for assemble — server uses the quote's order)
	}
	return v2Opts
}

// --------------------------
// | Market Data Conversions |
// --------------------------

// marketsToSupportedTokens converts a GetMarketsResponse to a list of unique ApiToken
func marketsToSupportedTokens(resp *api_types.GetMarketsResponse) []api_types.ApiToken {
	seen := make(map[string]bool)
	var tokens []api_types.ApiToken
	for _, market := range resp.Markets {
		if !seen[market.Base.Address] {
			seen[market.Base.Address] = true
			tokens = append(tokens, market.Base)
		}
		if !seen[market.Quote.Address] {
			seen[market.Quote.Address] = true
			tokens = append(tokens, market.Quote)
		}
	}
	return tokens
}

// marketsToFeeForAsset finds a market containing the given asset and returns the fee rates
// as an ExternalMatchFee (float64 values)
func marketsToFeeForAsset(resp *api_types.GetMarketsResponse, addr string) (*ExternalMatchFee, error) {
	for _, market := range resp.Markets {
		if market.Base.Address == addr || market.Quote.Address == addr {
			return &ExternalMatchFee{
				RelayerFee:  market.ExternalMatchFeeRates.RelayerFeeRate.ToF64(),
				ProtocolFee: market.ExternalMatchFeeRates.ProtocolFeeRate.ToF64(),
			}, nil
		}
	}
	return nil, fmt.Errorf("no market found for asset: %s", addr)
}
