package client

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// postcardSerializeTransfer serializes a withdrawal transfer in the format expected by the renegade contracts:
//
//	https://github.com/renegade-fi/renegade-contracts/blob/main/contracts-common/src/types.rs#L204
func postcardSerializeTransfer(mint string, amount *big.Int, destination string) ([]byte, error) {
	// Serialize the destination address as a 20 byte array
	destinationBytes, err := postcardSerializeAddress(destination)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize destination address: %w", err)
	}

	// Serialize the mint as a 20 byte array
	mintBytes, err := postcardSerializeAddress(mint)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize mint: %w", err)
	}

	// Serialize the amount as a u256
	amountBytes, err := postcardSerializeU256(amount)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize amount: %w", err)
	}

	// Append all the bytes together
	transferBytes := append(destinationBytes, mintBytes...)
	transferBytes = append(transferBytes, amountBytes...)
	transferBytes = append(transferBytes, 1) // withdraw flag

	return transferBytes, nil
}

// postcardSerializeAddress serializes an address to the format expected by the renegade contracts
func postcardSerializeAddress(address string) ([]byte, error) {
	// Remove '0x' prefix if present
	if len(address) >= 2 && address[:2] == "0x" {
		address = address[2:]
	}

	addressBytes := common.Hex2Bytes(address)
	addressBytesPadded := common.LeftPadBytes(addressBytes, 20)
	if len(addressBytesPadded) != 20 {
		return nil, fmt.Errorf("address must be 20 bytes, got %d bytes", len(addressBytesPadded))
	}

	return addressBytesPadded, nil
}

// postcardSerializeU256 serializes a u256 to the format expected by the renegade contracts
func postcardSerializeU256(val *big.Int) ([]byte, error) {
	const nLimbs = 4

	// Clone the big.Int to avoid modifying the original
	val = new(big.Int).Set(val)
	if val.BitLen() > 256 {
		return nil, fmt.Errorf("value exceeds 256 bits")
	}

	// Initialize result with full length
	result := make([]byte, nLimbs*binary.MaxVarintLen64)
	cursor := 0
	for i := 0; i < nLimbs; i++ {
		limb := val.Uint64()
		val.Rsh(val, 64)

		n := binary.PutUvarint(result[cursor:], limb)
		cursor += n
	}

	// Trim any unused bytes
	return result[:cursor], nil
}
