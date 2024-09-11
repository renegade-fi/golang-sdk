package client

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"renegade.fi/golang-sdk/abis"
)

type PermitWitnessTransferFrom struct {
	Permitted abis.ISignatureTransferTokenPermissions
	Spender   common.Address
	Nonce     *big.Int
	Deadline  *big.Int
	Witness   *DepositWitness
}

// DepositWitness is the witness for the permit
type DepositWitness struct {
	// PkRoot is the root of the public key serialized as u256 values
	PkRoot [4]*big.Int
}

const PERMIT2_EIP712_DOMAIN_NAME = "Permit2"

type EIP712Domain struct {
	Name              string
	ChainId           *big.Int
	VerifyingContract common.Address
}

// ConstructEIP712Domain constructs an EIP712Domain
func ConstructEIP712Domain(chainId *big.Int, verifyingContract common.Address) EIP712Domain {
	return EIP712Domain{
		Name:              PERMIT2_EIP712_DOMAIN_NAME,
		ChainId:           chainId,
		VerifyingContract: verifyingContract,
	}
}

// Hash hashes the EIP712Domain
func (domain EIP712Domain) Hash() common.Hash {
	typeHash := crypto.Keccak256(
		[]byte("EIP712Domain(string name,uint256 chainId,address verifyingContract)"),
	)

	return crypto.Keccak256Hash(
		typeHash,
		crypto.Keccak256([]byte(domain.Name)),
		common.LeftPadBytes(domain.ChainId.Bytes(), 32),
		common.LeftPadBytes(domain.VerifyingContract.Bytes(), 32),
	)
}

// getPermitSigningHash gets the eip712 hash of the permit
func getPermitSigningHash(permit PermitWitnessTransferFrom, domain EIP712Domain) (common.Hash, error) {
	// EIP-712 type hashes
	permitTypeHash := crypto.Keccak256(
		[]byte("PermitWitnessTransferFrom(TokenPermissions permitted,address spender,uint256 nonce,uint256 deadline,DepositWitness witness)DepositWitness(uint256[4] pkRoot)TokenPermissions(address token,uint256 amount)"),
	)

	// Hash TokenPermissions
	tokenPermissionsHash := crypto.Keccak256(
		crypto.Keccak256([]byte("TokenPermissions(address token,uint256 amount)")),
		common.LeftPadBytes(permit.Permitted.Token.Bytes(), 32),
		common.LeftPadBytes(permit.Permitted.Amount.Bytes(), 32),
	)

	// Construct the struct hash
	witnessHash := hashPermit2Witness(permit.Witness)
	structHash := crypto.Keccak256(
		permitTypeHash,
		tokenPermissionsHash,
		common.LeftPadBytes(permit.Spender.Bytes(), 32),
		common.LeftPadBytes(permit.Nonce.Bytes(), 32),
		common.LeftPadBytes(permit.Deadline.Bytes(), 32),
		witnessHash,
	)

	// Compute the final hash
	return crypto.Keccak256Hash(
		[]byte("\x19\x01"),
		domain.Hash().Bytes(),
		structHash,
	), nil
}

// hashPermit2Witness hashes the DepositWitness struct
func hashPermit2Witness(permit *DepositWitness) []byte {
	permitTypeHash := crypto.Keccak256(
		[]byte("DepositWitness(uint256[4] pkRoot)"),
	)

	// Hash the array of uint256 values
	pkRootHash := crypto.Keccak256(
		common.LeftPadBytes(permit.PkRoot[0].Bytes(), 32),
		common.LeftPadBytes(permit.PkRoot[1].Bytes(), 32),
		common.LeftPadBytes(permit.PkRoot[2].Bytes(), 32),
		common.LeftPadBytes(permit.PkRoot[3].Bytes(), 32),
	)

	witnessHash := crypto.Keccak256(
		permitTypeHash,
		pkRootHash,
	)

	return witnessHash
}
