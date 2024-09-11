// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abis

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// IAllowanceTransferAllowanceTransferDetails is an auto generated low-level Go binding around an user-defined struct.
type IAllowanceTransferAllowanceTransferDetails struct {
	From   common.Address
	To     common.Address
	Amount *big.Int
	Token  common.Address
}

// IAllowanceTransferPermitBatch is an auto generated low-level Go binding around an user-defined struct.
type IAllowanceTransferPermitBatch struct {
	Details     []IAllowanceTransferPermitDetails
	Spender     common.Address
	SigDeadline *big.Int
}

// IAllowanceTransferPermitDetails is an auto generated low-level Go binding around an user-defined struct.
type IAllowanceTransferPermitDetails struct {
	Token      common.Address
	Amount     *big.Int
	Expiration *big.Int
	Nonce      *big.Int
}

// IAllowanceTransferPermitSingle is an auto generated low-level Go binding around an user-defined struct.
type IAllowanceTransferPermitSingle struct {
	Details     IAllowanceTransferPermitDetails
	Spender     common.Address
	SigDeadline *big.Int
}

// IAllowanceTransferTokenSpenderPair is an auto generated low-level Go binding around an user-defined struct.
type IAllowanceTransferTokenSpenderPair struct {
	Token   common.Address
	Spender common.Address
}

// ISignatureTransferPermitBatchTransferFrom is an auto generated low-level Go binding around an user-defined struct.
type ISignatureTransferPermitBatchTransferFrom struct {
	Permitted []ISignatureTransferTokenPermissions
	Nonce     *big.Int
	Deadline  *big.Int
}

// ISignatureTransferPermitTransferFrom is an auto generated low-level Go binding around an user-defined struct.
type ISignatureTransferPermitTransferFrom struct {
	Permitted ISignatureTransferTokenPermissions
	Nonce     *big.Int
	Deadline  *big.Int
}

// ISignatureTransferSignatureTransferDetails is an auto generated low-level Go binding around an user-defined struct.
type ISignatureTransferSignatureTransferDetails struct {
	To              common.Address
	RequestedAmount *big.Int
}

// ISignatureTransferTokenPermissions is an auto generated low-level Go binding around an user-defined struct.
type ISignatureTransferTokenPermissions struct {
	Token  common.Address
	Amount *big.Int
}

// AbisMetaData contains all meta data concerning the Abis contract.
var AbisMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"AllowanceExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExcessiveInvalidation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maxAmount\",\"type\":\"uint256\"}],\"name\":\"InvalidAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidContractSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidNonce\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LengthMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"signatureDeadline\",\"type\":\"uint256\"}],\"name\":\"SignatureExpired\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint160\",\"name\":\"amount\",\"type\":\"uint160\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"expiration\",\"type\":\"uint48\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"Lockdown\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"newNonce\",\"type\":\"uint48\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"oldNonce\",\"type\":\"uint48\"}],\"name\":\"NonceInvalidation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint160\",\"name\":\"amount\",\"type\":\"uint160\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"expiration\",\"type\":\"uint48\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"nonce\",\"type\":\"uint48\"}],\"name\":\"Permit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"word\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"mask\",\"type\":\"uint256\"}],\"name\":\"UnorderedNonceInvalidation\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DOMAIN_SEPARATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint160\",\"name\":\"amount\",\"type\":\"uint160\"},{\"internalType\":\"uint48\",\"name\":\"expiration\",\"type\":\"uint48\"},{\"internalType\":\"uint48\",\"name\":\"nonce\",\"type\":\"uint48\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint160\",\"name\":\"amount\",\"type\":\"uint160\"},{\"internalType\":\"uint48\",\"name\":\"expiration\",\"type\":\"uint48\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"newNonce\",\"type\":\"uint48\"}],\"name\":\"invalidateNonces\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"wordPos\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"mask\",\"type\":\"uint256\"}],\"name\":\"invalidateUnorderedNonces\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"internalType\":\"structIAllowanceTransfer.TokenSpenderPair[]\",\"name\":\"approvals\",\"type\":\"tuple[]\"}],\"name\":\"lockdown\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"nonceBitmap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint160\",\"name\":\"amount\",\"type\":\"uint160\"},{\"internalType\":\"uint48\",\"name\":\"expiration\",\"type\":\"uint48\"},{\"internalType\":\"uint48\",\"name\":\"nonce\",\"type\":\"uint48\"}],\"internalType\":\"structIAllowanceTransfer.PermitDetails[]\",\"name\":\"details\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"sigDeadline\",\"type\":\"uint256\"}],\"internalType\":\"structIAllowanceTransfer.PermitBatch\",\"name\":\"permitBatch\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"permit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint160\",\"name\":\"amount\",\"type\":\"uint160\"},{\"internalType\":\"uint48\",\"name\":\"expiration\",\"type\":\"uint48\"},{\"internalType\":\"uint48\",\"name\":\"nonce\",\"type\":\"uint48\"}],\"internalType\":\"structIAllowanceTransfer.PermitDetails\",\"name\":\"details\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"sigDeadline\",\"type\":\"uint256\"}],\"internalType\":\"structIAllowanceTransfer.PermitSingle\",\"name\":\"permitSingle\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"permit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structISignatureTransfer.TokenPermissions\",\"name\":\"permitted\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"internalType\":\"structISignatureTransfer.PermitTransferFrom\",\"name\":\"permit\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"requestedAmount\",\"type\":\"uint256\"}],\"internalType\":\"structISignatureTransfer.SignatureTransferDetails\",\"name\":\"transferDetails\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"permitTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structISignatureTransfer.TokenPermissions[]\",\"name\":\"permitted\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"internalType\":\"structISignatureTransfer.PermitBatchTransferFrom\",\"name\":\"permit\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"requestedAmount\",\"type\":\"uint256\"}],\"internalType\":\"structISignatureTransfer.SignatureTransferDetails[]\",\"name\":\"transferDetails\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"permitTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structISignatureTransfer.TokenPermissions\",\"name\":\"permitted\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"internalType\":\"structISignatureTransfer.PermitTransferFrom\",\"name\":\"permit\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"requestedAmount\",\"type\":\"uint256\"}],\"internalType\":\"structISignatureTransfer.SignatureTransferDetails\",\"name\":\"transferDetails\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"witness\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"witnessTypeString\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"permitWitnessTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structISignatureTransfer.TokenPermissions[]\",\"name\":\"permitted\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"internalType\":\"structISignatureTransfer.PermitBatchTransferFrom\",\"name\":\"permit\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"requestedAmount\",\"type\":\"uint256\"}],\"internalType\":\"structISignatureTransfer.SignatureTransferDetails[]\",\"name\":\"transferDetails\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"witness\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"witnessTypeString\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"permitWitnessTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint160\",\"name\":\"amount\",\"type\":\"uint160\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"internalType\":\"structIAllowanceTransfer.AllowanceTransferDetails[]\",\"name\":\"transferDetails\",\"type\":\"tuple[]\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint160\",\"name\":\"amount\",\"type\":\"uint160\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// AbisABI is the input ABI used to generate the binding from.
// Deprecated: Use AbisMetaData.ABI instead.
var AbisABI = AbisMetaData.ABI

// Abis is an auto generated Go binding around an Ethereum contract.
type Abis struct {
	AbisCaller     // Read-only binding to the contract
	AbisTransactor // Write-only binding to the contract
	AbisFilterer   // Log filterer for contract events
}

// AbisCaller is an auto generated read-only Go binding around an Ethereum contract.
type AbisCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AbisTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AbisTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AbisFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AbisFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AbisSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AbisSession struct {
	Contract     *Abis             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AbisCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AbisCallerSession struct {
	Contract *AbisCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// AbisTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AbisTransactorSession struct {
	Contract     *AbisTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AbisRaw is an auto generated low-level Go binding around an Ethereum contract.
type AbisRaw struct {
	Contract *Abis // Generic contract binding to access the raw methods on
}

// AbisCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AbisCallerRaw struct {
	Contract *AbisCaller // Generic read-only contract binding to access the raw methods on
}

// AbisTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AbisTransactorRaw struct {
	Contract *AbisTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAbis creates a new instance of Abis, bound to a specific deployed contract.
func NewAbis(address common.Address, backend bind.ContractBackend) (*Abis, error) {
	contract, err := bindAbis(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Abis{AbisCaller: AbisCaller{contract: contract}, AbisTransactor: AbisTransactor{contract: contract}, AbisFilterer: AbisFilterer{contract: contract}}, nil
}

// NewAbisCaller creates a new read-only instance of Abis, bound to a specific deployed contract.
func NewAbisCaller(address common.Address, caller bind.ContractCaller) (*AbisCaller, error) {
	contract, err := bindAbis(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AbisCaller{contract: contract}, nil
}

// NewAbisTransactor creates a new write-only instance of Abis, bound to a specific deployed contract.
func NewAbisTransactor(address common.Address, transactor bind.ContractTransactor) (*AbisTransactor, error) {
	contract, err := bindAbis(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AbisTransactor{contract: contract}, nil
}

// NewAbisFilterer creates a new log filterer instance of Abis, bound to a specific deployed contract.
func NewAbisFilterer(address common.Address, filterer bind.ContractFilterer) (*AbisFilterer, error) {
	contract, err := bindAbis(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AbisFilterer{contract: contract}, nil
}

// bindAbis binds a generic wrapper to an already deployed contract.
func bindAbis(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AbisMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Abis *AbisRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Abis.Contract.AbisCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Abis *AbisRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Abis.Contract.AbisTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Abis *AbisRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Abis.Contract.AbisTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Abis *AbisCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Abis.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Abis *AbisTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Abis.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Abis *AbisTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Abis.Contract.contract.Transact(opts, method, params...)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_Abis *AbisCaller) DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Abis.contract.Call(opts, &out, "DOMAIN_SEPARATOR")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_Abis *AbisSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _Abis.Contract.DOMAINSEPARATOR(&_Abis.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_Abis *AbisCallerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _Abis.Contract.DOMAINSEPARATOR(&_Abis.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0x927da105.
//
// Solidity: function allowance(address , address , address ) view returns(uint160 amount, uint48 expiration, uint48 nonce)
func (_Abis *AbisCaller) Allowance(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address, arg2 common.Address) (struct {
	Amount     *big.Int
	Expiration *big.Int
	Nonce      *big.Int
}, error) {
	var out []interface{}
	err := _Abis.contract.Call(opts, &out, "allowance", arg0, arg1, arg2)

	outstruct := new(struct {
		Amount     *big.Int
		Expiration *big.Int
		Nonce      *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Amount = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Expiration = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Nonce = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Allowance is a free data retrieval call binding the contract method 0x927da105.
//
// Solidity: function allowance(address , address , address ) view returns(uint160 amount, uint48 expiration, uint48 nonce)
func (_Abis *AbisSession) Allowance(arg0 common.Address, arg1 common.Address, arg2 common.Address) (struct {
	Amount     *big.Int
	Expiration *big.Int
	Nonce      *big.Int
}, error) {
	return _Abis.Contract.Allowance(&_Abis.CallOpts, arg0, arg1, arg2)
}

// Allowance is a free data retrieval call binding the contract method 0x927da105.
//
// Solidity: function allowance(address , address , address ) view returns(uint160 amount, uint48 expiration, uint48 nonce)
func (_Abis *AbisCallerSession) Allowance(arg0 common.Address, arg1 common.Address, arg2 common.Address) (struct {
	Amount     *big.Int
	Expiration *big.Int
	Nonce      *big.Int
}, error) {
	return _Abis.Contract.Allowance(&_Abis.CallOpts, arg0, arg1, arg2)
}

// NonceBitmap is a free data retrieval call binding the contract method 0x4fe02b44.
//
// Solidity: function nonceBitmap(address , uint256 ) view returns(uint256)
func (_Abis *AbisCaller) NonceBitmap(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Abis.contract.Call(opts, &out, "nonceBitmap", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NonceBitmap is a free data retrieval call binding the contract method 0x4fe02b44.
//
// Solidity: function nonceBitmap(address , uint256 ) view returns(uint256)
func (_Abis *AbisSession) NonceBitmap(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _Abis.Contract.NonceBitmap(&_Abis.CallOpts, arg0, arg1)
}

// NonceBitmap is a free data retrieval call binding the contract method 0x4fe02b44.
//
// Solidity: function nonceBitmap(address , uint256 ) view returns(uint256)
func (_Abis *AbisCallerSession) NonceBitmap(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _Abis.Contract.NonceBitmap(&_Abis.CallOpts, arg0, arg1)
}

// Approve is a paid mutator transaction binding the contract method 0x87517c45.
//
// Solidity: function approve(address token, address spender, uint160 amount, uint48 expiration) returns()
func (_Abis *AbisTransactor) Approve(opts *bind.TransactOpts, token common.Address, spender common.Address, amount *big.Int, expiration *big.Int) (*types.Transaction, error) {
	return _Abis.contract.Transact(opts, "approve", token, spender, amount, expiration)
}

// Approve is a paid mutator transaction binding the contract method 0x87517c45.
//
// Solidity: function approve(address token, address spender, uint160 amount, uint48 expiration) returns()
func (_Abis *AbisSession) Approve(token common.Address, spender common.Address, amount *big.Int, expiration *big.Int) (*types.Transaction, error) {
	return _Abis.Contract.Approve(&_Abis.TransactOpts, token, spender, amount, expiration)
}

// Approve is a paid mutator transaction binding the contract method 0x87517c45.
//
// Solidity: function approve(address token, address spender, uint160 amount, uint48 expiration) returns()
func (_Abis *AbisTransactorSession) Approve(token common.Address, spender common.Address, amount *big.Int, expiration *big.Int) (*types.Transaction, error) {
	return _Abis.Contract.Approve(&_Abis.TransactOpts, token, spender, amount, expiration)
}

// InvalidateNonces is a paid mutator transaction binding the contract method 0x65d9723c.
//
// Solidity: function invalidateNonces(address token, address spender, uint48 newNonce) returns()
func (_Abis *AbisTransactor) InvalidateNonces(opts *bind.TransactOpts, token common.Address, spender common.Address, newNonce *big.Int) (*types.Transaction, error) {
	return _Abis.contract.Transact(opts, "invalidateNonces", token, spender, newNonce)
}

// InvalidateNonces is a paid mutator transaction binding the contract method 0x65d9723c.
//
// Solidity: function invalidateNonces(address token, address spender, uint48 newNonce) returns()
func (_Abis *AbisSession) InvalidateNonces(token common.Address, spender common.Address, newNonce *big.Int) (*types.Transaction, error) {
	return _Abis.Contract.InvalidateNonces(&_Abis.TransactOpts, token, spender, newNonce)
}

// InvalidateNonces is a paid mutator transaction binding the contract method 0x65d9723c.
//
// Solidity: function invalidateNonces(address token, address spender, uint48 newNonce) returns()
func (_Abis *AbisTransactorSession) InvalidateNonces(token common.Address, spender common.Address, newNonce *big.Int) (*types.Transaction, error) {
	return _Abis.Contract.InvalidateNonces(&_Abis.TransactOpts, token, spender, newNonce)
}

// InvalidateUnorderedNonces is a paid mutator transaction binding the contract method 0x3ff9dcb1.
//
// Solidity: function invalidateUnorderedNonces(uint256 wordPos, uint256 mask) returns()
func (_Abis *AbisTransactor) InvalidateUnorderedNonces(opts *bind.TransactOpts, wordPos *big.Int, mask *big.Int) (*types.Transaction, error) {
	return _Abis.contract.Transact(opts, "invalidateUnorderedNonces", wordPos, mask)
}

// InvalidateUnorderedNonces is a paid mutator transaction binding the contract method 0x3ff9dcb1.
//
// Solidity: function invalidateUnorderedNonces(uint256 wordPos, uint256 mask) returns()
func (_Abis *AbisSession) InvalidateUnorderedNonces(wordPos *big.Int, mask *big.Int) (*types.Transaction, error) {
	return _Abis.Contract.InvalidateUnorderedNonces(&_Abis.TransactOpts, wordPos, mask)
}

// InvalidateUnorderedNonces is a paid mutator transaction binding the contract method 0x3ff9dcb1.
//
// Solidity: function invalidateUnorderedNonces(uint256 wordPos, uint256 mask) returns()
func (_Abis *AbisTransactorSession) InvalidateUnorderedNonces(wordPos *big.Int, mask *big.Int) (*types.Transaction, error) {
	return _Abis.Contract.InvalidateUnorderedNonces(&_Abis.TransactOpts, wordPos, mask)
}

// Lockdown is a paid mutator transaction binding the contract method 0xcc53287f.
//
// Solidity: function lockdown((address,address)[] approvals) returns()
func (_Abis *AbisTransactor) Lockdown(opts *bind.TransactOpts, approvals []IAllowanceTransferTokenSpenderPair) (*types.Transaction, error) {
	return _Abis.contract.Transact(opts, "lockdown", approvals)
}

// Lockdown is a paid mutator transaction binding the contract method 0xcc53287f.
//
// Solidity: function lockdown((address,address)[] approvals) returns()
func (_Abis *AbisSession) Lockdown(approvals []IAllowanceTransferTokenSpenderPair) (*types.Transaction, error) {
	return _Abis.Contract.Lockdown(&_Abis.TransactOpts, approvals)
}

// Lockdown is a paid mutator transaction binding the contract method 0xcc53287f.
//
// Solidity: function lockdown((address,address)[] approvals) returns()
func (_Abis *AbisTransactorSession) Lockdown(approvals []IAllowanceTransferTokenSpenderPair) (*types.Transaction, error) {
	return _Abis.Contract.Lockdown(&_Abis.TransactOpts, approvals)
}

// Permit is a paid mutator transaction binding the contract method 0x2a2d80d1.
//
// Solidity: function permit(address owner, ((address,uint160,uint48,uint48)[],address,uint256) permitBatch, bytes signature) returns()
func (_Abis *AbisTransactor) Permit(opts *bind.TransactOpts, owner common.Address, permitBatch IAllowanceTransferPermitBatch, signature []byte) (*types.Transaction, error) {
	return _Abis.contract.Transact(opts, "permit", owner, permitBatch, signature)
}

// Permit is a paid mutator transaction binding the contract method 0x2a2d80d1.
//
// Solidity: function permit(address owner, ((address,uint160,uint48,uint48)[],address,uint256) permitBatch, bytes signature) returns()
func (_Abis *AbisSession) Permit(owner common.Address, permitBatch IAllowanceTransferPermitBatch, signature []byte) (*types.Transaction, error) {
	return _Abis.Contract.Permit(&_Abis.TransactOpts, owner, permitBatch, signature)
}

// Permit is a paid mutator transaction binding the contract method 0x2a2d80d1.
//
// Solidity: function permit(address owner, ((address,uint160,uint48,uint48)[],address,uint256) permitBatch, bytes signature) returns()
func (_Abis *AbisTransactorSession) Permit(owner common.Address, permitBatch IAllowanceTransferPermitBatch, signature []byte) (*types.Transaction, error) {
	return _Abis.Contract.Permit(&_Abis.TransactOpts, owner, permitBatch, signature)
}

// Permit0 is a paid mutator transaction binding the contract method 0x2b67b570.
//
// Solidity: function permit(address owner, ((address,uint160,uint48,uint48),address,uint256) permitSingle, bytes signature) returns()
func (_Abis *AbisTransactor) Permit0(opts *bind.TransactOpts, owner common.Address, permitSingle IAllowanceTransferPermitSingle, signature []byte) (*types.Transaction, error) {
	return _Abis.contract.Transact(opts, "permit0", owner, permitSingle, signature)
}

// Permit0 is a paid mutator transaction binding the contract method 0x2b67b570.
//
// Solidity: function permit(address owner, ((address,uint160,uint48,uint48),address,uint256) permitSingle, bytes signature) returns()
func (_Abis *AbisSession) Permit0(owner common.Address, permitSingle IAllowanceTransferPermitSingle, signature []byte) (*types.Transaction, error) {
	return _Abis.Contract.Permit0(&_Abis.TransactOpts, owner, permitSingle, signature)
}

// Permit0 is a paid mutator transaction binding the contract method 0x2b67b570.
//
// Solidity: function permit(address owner, ((address,uint160,uint48,uint48),address,uint256) permitSingle, bytes signature) returns()
func (_Abis *AbisTransactorSession) Permit0(owner common.Address, permitSingle IAllowanceTransferPermitSingle, signature []byte) (*types.Transaction, error) {
	return _Abis.Contract.Permit0(&_Abis.TransactOpts, owner, permitSingle, signature)
}

// PermitTransferFrom is a paid mutator transaction binding the contract method 0x30f28b7a.
//
// Solidity: function permitTransferFrom(((address,uint256),uint256,uint256) permit, (address,uint256) transferDetails, address owner, bytes signature) returns()
func (_Abis *AbisTransactor) PermitTransferFrom(opts *bind.TransactOpts, permit ISignatureTransferPermitTransferFrom, transferDetails ISignatureTransferSignatureTransferDetails, owner common.Address, signature []byte) (*types.Transaction, error) {
	return _Abis.contract.Transact(opts, "permitTransferFrom", permit, transferDetails, owner, signature)
}

// PermitTransferFrom is a paid mutator transaction binding the contract method 0x30f28b7a.
//
// Solidity: function permitTransferFrom(((address,uint256),uint256,uint256) permit, (address,uint256) transferDetails, address owner, bytes signature) returns()
func (_Abis *AbisSession) PermitTransferFrom(permit ISignatureTransferPermitTransferFrom, transferDetails ISignatureTransferSignatureTransferDetails, owner common.Address, signature []byte) (*types.Transaction, error) {
	return _Abis.Contract.PermitTransferFrom(&_Abis.TransactOpts, permit, transferDetails, owner, signature)
}

// PermitTransferFrom is a paid mutator transaction binding the contract method 0x30f28b7a.
//
// Solidity: function permitTransferFrom(((address,uint256),uint256,uint256) permit, (address,uint256) transferDetails, address owner, bytes signature) returns()
func (_Abis *AbisTransactorSession) PermitTransferFrom(permit ISignatureTransferPermitTransferFrom, transferDetails ISignatureTransferSignatureTransferDetails, owner common.Address, signature []byte) (*types.Transaction, error) {
	return _Abis.Contract.PermitTransferFrom(&_Abis.TransactOpts, permit, transferDetails, owner, signature)
}

// PermitTransferFrom0 is a paid mutator transaction binding the contract method 0xedd9444b.
//
// Solidity: function permitTransferFrom(((address,uint256)[],uint256,uint256) permit, (address,uint256)[] transferDetails, address owner, bytes signature) returns()
func (_Abis *AbisTransactor) PermitTransferFrom0(opts *bind.TransactOpts, permit ISignatureTransferPermitBatchTransferFrom, transferDetails []ISignatureTransferSignatureTransferDetails, owner common.Address, signature []byte) (*types.Transaction, error) {
	return _Abis.contract.Transact(opts, "permitTransferFrom0", permit, transferDetails, owner, signature)
}

// PermitTransferFrom0 is a paid mutator transaction binding the contract method 0xedd9444b.
//
// Solidity: function permitTransferFrom(((address,uint256)[],uint256,uint256) permit, (address,uint256)[] transferDetails, address owner, bytes signature) returns()
func (_Abis *AbisSession) PermitTransferFrom0(permit ISignatureTransferPermitBatchTransferFrom, transferDetails []ISignatureTransferSignatureTransferDetails, owner common.Address, signature []byte) (*types.Transaction, error) {
	return _Abis.Contract.PermitTransferFrom0(&_Abis.TransactOpts, permit, transferDetails, owner, signature)
}

// PermitTransferFrom0 is a paid mutator transaction binding the contract method 0xedd9444b.
//
// Solidity: function permitTransferFrom(((address,uint256)[],uint256,uint256) permit, (address,uint256)[] transferDetails, address owner, bytes signature) returns()
func (_Abis *AbisTransactorSession) PermitTransferFrom0(permit ISignatureTransferPermitBatchTransferFrom, transferDetails []ISignatureTransferSignatureTransferDetails, owner common.Address, signature []byte) (*types.Transaction, error) {
	return _Abis.Contract.PermitTransferFrom0(&_Abis.TransactOpts, permit, transferDetails, owner, signature)
}

// PermitWitnessTransferFrom is a paid mutator transaction binding the contract method 0x137c29fe.
//
// Solidity: function permitWitnessTransferFrom(((address,uint256),uint256,uint256) permit, (address,uint256) transferDetails, address owner, bytes32 witness, string witnessTypeString, bytes signature) returns()
func (_Abis *AbisTransactor) PermitWitnessTransferFrom(opts *bind.TransactOpts, permit ISignatureTransferPermitTransferFrom, transferDetails ISignatureTransferSignatureTransferDetails, owner common.Address, witness [32]byte, witnessTypeString string, signature []byte) (*types.Transaction, error) {
	return _Abis.contract.Transact(opts, "permitWitnessTransferFrom", permit, transferDetails, owner, witness, witnessTypeString, signature)
}

// PermitWitnessTransferFrom is a paid mutator transaction binding the contract method 0x137c29fe.
//
// Solidity: function permitWitnessTransferFrom(((address,uint256),uint256,uint256) permit, (address,uint256) transferDetails, address owner, bytes32 witness, string witnessTypeString, bytes signature) returns()
func (_Abis *AbisSession) PermitWitnessTransferFrom(permit ISignatureTransferPermitTransferFrom, transferDetails ISignatureTransferSignatureTransferDetails, owner common.Address, witness [32]byte, witnessTypeString string, signature []byte) (*types.Transaction, error) {
	return _Abis.Contract.PermitWitnessTransferFrom(&_Abis.TransactOpts, permit, transferDetails, owner, witness, witnessTypeString, signature)
}

// PermitWitnessTransferFrom is a paid mutator transaction binding the contract method 0x137c29fe.
//
// Solidity: function permitWitnessTransferFrom(((address,uint256),uint256,uint256) permit, (address,uint256) transferDetails, address owner, bytes32 witness, string witnessTypeString, bytes signature) returns()
func (_Abis *AbisTransactorSession) PermitWitnessTransferFrom(permit ISignatureTransferPermitTransferFrom, transferDetails ISignatureTransferSignatureTransferDetails, owner common.Address, witness [32]byte, witnessTypeString string, signature []byte) (*types.Transaction, error) {
	return _Abis.Contract.PermitWitnessTransferFrom(&_Abis.TransactOpts, permit, transferDetails, owner, witness, witnessTypeString, signature)
}

// PermitWitnessTransferFrom0 is a paid mutator transaction binding the contract method 0xfe8ec1a7.
//
// Solidity: function permitWitnessTransferFrom(((address,uint256)[],uint256,uint256) permit, (address,uint256)[] transferDetails, address owner, bytes32 witness, string witnessTypeString, bytes signature) returns()
func (_Abis *AbisTransactor) PermitWitnessTransferFrom0(opts *bind.TransactOpts, permit ISignatureTransferPermitBatchTransferFrom, transferDetails []ISignatureTransferSignatureTransferDetails, owner common.Address, witness [32]byte, witnessTypeString string, signature []byte) (*types.Transaction, error) {
	return _Abis.contract.Transact(opts, "permitWitnessTransferFrom0", permit, transferDetails, owner, witness, witnessTypeString, signature)
}

// PermitWitnessTransferFrom0 is a paid mutator transaction binding the contract method 0xfe8ec1a7.
//
// Solidity: function permitWitnessTransferFrom(((address,uint256)[],uint256,uint256) permit, (address,uint256)[] transferDetails, address owner, bytes32 witness, string witnessTypeString, bytes signature) returns()
func (_Abis *AbisSession) PermitWitnessTransferFrom0(permit ISignatureTransferPermitBatchTransferFrom, transferDetails []ISignatureTransferSignatureTransferDetails, owner common.Address, witness [32]byte, witnessTypeString string, signature []byte) (*types.Transaction, error) {
	return _Abis.Contract.PermitWitnessTransferFrom0(&_Abis.TransactOpts, permit, transferDetails, owner, witness, witnessTypeString, signature)
}

// PermitWitnessTransferFrom0 is a paid mutator transaction binding the contract method 0xfe8ec1a7.
//
// Solidity: function permitWitnessTransferFrom(((address,uint256)[],uint256,uint256) permit, (address,uint256)[] transferDetails, address owner, bytes32 witness, string witnessTypeString, bytes signature) returns()
func (_Abis *AbisTransactorSession) PermitWitnessTransferFrom0(permit ISignatureTransferPermitBatchTransferFrom, transferDetails []ISignatureTransferSignatureTransferDetails, owner common.Address, witness [32]byte, witnessTypeString string, signature []byte) (*types.Transaction, error) {
	return _Abis.Contract.PermitWitnessTransferFrom0(&_Abis.TransactOpts, permit, transferDetails, owner, witness, witnessTypeString, signature)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x0d58b1db.
//
// Solidity: function transferFrom((address,address,uint160,address)[] transferDetails) returns()
func (_Abis *AbisTransactor) TransferFrom(opts *bind.TransactOpts, transferDetails []IAllowanceTransferAllowanceTransferDetails) (*types.Transaction, error) {
	return _Abis.contract.Transact(opts, "transferFrom", transferDetails)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x0d58b1db.
//
// Solidity: function transferFrom((address,address,uint160,address)[] transferDetails) returns()
func (_Abis *AbisSession) TransferFrom(transferDetails []IAllowanceTransferAllowanceTransferDetails) (*types.Transaction, error) {
	return _Abis.Contract.TransferFrom(&_Abis.TransactOpts, transferDetails)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x0d58b1db.
//
// Solidity: function transferFrom((address,address,uint160,address)[] transferDetails) returns()
func (_Abis *AbisTransactorSession) TransferFrom(transferDetails []IAllowanceTransferAllowanceTransferDetails) (*types.Transaction, error) {
	return _Abis.Contract.TransferFrom(&_Abis.TransactOpts, transferDetails)
}

// TransferFrom0 is a paid mutator transaction binding the contract method 0x36c78516.
//
// Solidity: function transferFrom(address from, address to, uint160 amount, address token) returns()
func (_Abis *AbisTransactor) TransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int, token common.Address) (*types.Transaction, error) {
	return _Abis.contract.Transact(opts, "transferFrom0", from, to, amount, token)
}

// TransferFrom0 is a paid mutator transaction binding the contract method 0x36c78516.
//
// Solidity: function transferFrom(address from, address to, uint160 amount, address token) returns()
func (_Abis *AbisSession) TransferFrom0(from common.Address, to common.Address, amount *big.Int, token common.Address) (*types.Transaction, error) {
	return _Abis.Contract.TransferFrom0(&_Abis.TransactOpts, from, to, amount, token)
}

// TransferFrom0 is a paid mutator transaction binding the contract method 0x36c78516.
//
// Solidity: function transferFrom(address from, address to, uint160 amount, address token) returns()
func (_Abis *AbisTransactorSession) TransferFrom0(from common.Address, to common.Address, amount *big.Int, token common.Address) (*types.Transaction, error) {
	return _Abis.Contract.TransferFrom0(&_Abis.TransactOpts, from, to, amount, token)
}

// AbisApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Abis contract.
type AbisApprovalIterator struct {
	Event *AbisApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AbisApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AbisApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AbisApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AbisApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AbisApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AbisApproval represents a Approval event raised by the Abis contract.
type AbisApproval struct {
	Owner      common.Address
	Token      common.Address
	Spender    common.Address
	Amount     *big.Int
	Expiration *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0xda9fa7c1b00402c17d0161b249b1ab8bbec047c5a52207b9c112deffd817036b.
//
// Solidity: event Approval(address indexed owner, address indexed token, address indexed spender, uint160 amount, uint48 expiration)
func (_Abis *AbisFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, token []common.Address, spender []common.Address) (*AbisApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Abis.contract.FilterLogs(opts, "Approval", ownerRule, tokenRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &AbisApprovalIterator{contract: _Abis.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0xda9fa7c1b00402c17d0161b249b1ab8bbec047c5a52207b9c112deffd817036b.
//
// Solidity: event Approval(address indexed owner, address indexed token, address indexed spender, uint160 amount, uint48 expiration)
func (_Abis *AbisFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *AbisApproval, owner []common.Address, token []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Abis.contract.WatchLogs(opts, "Approval", ownerRule, tokenRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AbisApproval)
				if err := _Abis.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0xda9fa7c1b00402c17d0161b249b1ab8bbec047c5a52207b9c112deffd817036b.
//
// Solidity: event Approval(address indexed owner, address indexed token, address indexed spender, uint160 amount, uint48 expiration)
func (_Abis *AbisFilterer) ParseApproval(log types.Log) (*AbisApproval, error) {
	event := new(AbisApproval)
	if err := _Abis.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AbisLockdownIterator is returned from FilterLockdown and is used to iterate over the raw logs and unpacked data for Lockdown events raised by the Abis contract.
type AbisLockdownIterator struct {
	Event *AbisLockdown // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AbisLockdownIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AbisLockdown)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AbisLockdown)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AbisLockdownIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AbisLockdownIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AbisLockdown represents a Lockdown event raised by the Abis contract.
type AbisLockdown struct {
	Owner   common.Address
	Token   common.Address
	Spender common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLockdown is a free log retrieval operation binding the contract event 0x89b1add15eff56b3dfe299ad94e01f2b52fbcb80ae1a3baea6ae8c04cb2b98a4.
//
// Solidity: event Lockdown(address indexed owner, address token, address spender)
func (_Abis *AbisFilterer) FilterLockdown(opts *bind.FilterOpts, owner []common.Address) (*AbisLockdownIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _Abis.contract.FilterLogs(opts, "Lockdown", ownerRule)
	if err != nil {
		return nil, err
	}
	return &AbisLockdownIterator{contract: _Abis.contract, event: "Lockdown", logs: logs, sub: sub}, nil
}

// WatchLockdown is a free log subscription operation binding the contract event 0x89b1add15eff56b3dfe299ad94e01f2b52fbcb80ae1a3baea6ae8c04cb2b98a4.
//
// Solidity: event Lockdown(address indexed owner, address token, address spender)
func (_Abis *AbisFilterer) WatchLockdown(opts *bind.WatchOpts, sink chan<- *AbisLockdown, owner []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _Abis.contract.WatchLogs(opts, "Lockdown", ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AbisLockdown)
				if err := _Abis.contract.UnpackLog(event, "Lockdown", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseLockdown is a log parse operation binding the contract event 0x89b1add15eff56b3dfe299ad94e01f2b52fbcb80ae1a3baea6ae8c04cb2b98a4.
//
// Solidity: event Lockdown(address indexed owner, address token, address spender)
func (_Abis *AbisFilterer) ParseLockdown(log types.Log) (*AbisLockdown, error) {
	event := new(AbisLockdown)
	if err := _Abis.contract.UnpackLog(event, "Lockdown", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AbisNonceInvalidationIterator is returned from FilterNonceInvalidation and is used to iterate over the raw logs and unpacked data for NonceInvalidation events raised by the Abis contract.
type AbisNonceInvalidationIterator struct {
	Event *AbisNonceInvalidation // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AbisNonceInvalidationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AbisNonceInvalidation)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AbisNonceInvalidation)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AbisNonceInvalidationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AbisNonceInvalidationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AbisNonceInvalidation represents a NonceInvalidation event raised by the Abis contract.
type AbisNonceInvalidation struct {
	Owner    common.Address
	Token    common.Address
	Spender  common.Address
	NewNonce *big.Int
	OldNonce *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterNonceInvalidation is a free log retrieval operation binding the contract event 0x55eb90d810e1700b35a8e7e25395ff7f2b2259abd7415ca2284dfb1c246418f3.
//
// Solidity: event NonceInvalidation(address indexed owner, address indexed token, address indexed spender, uint48 newNonce, uint48 oldNonce)
func (_Abis *AbisFilterer) FilterNonceInvalidation(opts *bind.FilterOpts, owner []common.Address, token []common.Address, spender []common.Address) (*AbisNonceInvalidationIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Abis.contract.FilterLogs(opts, "NonceInvalidation", ownerRule, tokenRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &AbisNonceInvalidationIterator{contract: _Abis.contract, event: "NonceInvalidation", logs: logs, sub: sub}, nil
}

// WatchNonceInvalidation is a free log subscription operation binding the contract event 0x55eb90d810e1700b35a8e7e25395ff7f2b2259abd7415ca2284dfb1c246418f3.
//
// Solidity: event NonceInvalidation(address indexed owner, address indexed token, address indexed spender, uint48 newNonce, uint48 oldNonce)
func (_Abis *AbisFilterer) WatchNonceInvalidation(opts *bind.WatchOpts, sink chan<- *AbisNonceInvalidation, owner []common.Address, token []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Abis.contract.WatchLogs(opts, "NonceInvalidation", ownerRule, tokenRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AbisNonceInvalidation)
				if err := _Abis.contract.UnpackLog(event, "NonceInvalidation", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNonceInvalidation is a log parse operation binding the contract event 0x55eb90d810e1700b35a8e7e25395ff7f2b2259abd7415ca2284dfb1c246418f3.
//
// Solidity: event NonceInvalidation(address indexed owner, address indexed token, address indexed spender, uint48 newNonce, uint48 oldNonce)
func (_Abis *AbisFilterer) ParseNonceInvalidation(log types.Log) (*AbisNonceInvalidation, error) {
	event := new(AbisNonceInvalidation)
	if err := _Abis.contract.UnpackLog(event, "NonceInvalidation", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AbisPermitIterator is returned from FilterPermit and is used to iterate over the raw logs and unpacked data for Permit events raised by the Abis contract.
type AbisPermitIterator struct {
	Event *AbisPermit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AbisPermitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AbisPermit)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AbisPermit)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AbisPermitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AbisPermitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AbisPermit represents a Permit event raised by the Abis contract.
type AbisPermit struct {
	Owner      common.Address
	Token      common.Address
	Spender    common.Address
	Amount     *big.Int
	Expiration *big.Int
	Nonce      *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterPermit is a free log retrieval operation binding the contract event 0xc6a377bfc4eb120024a8ac08eef205be16b817020812c73223e81d1bdb9708ec.
//
// Solidity: event Permit(address indexed owner, address indexed token, address indexed spender, uint160 amount, uint48 expiration, uint48 nonce)
func (_Abis *AbisFilterer) FilterPermit(opts *bind.FilterOpts, owner []common.Address, token []common.Address, spender []common.Address) (*AbisPermitIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Abis.contract.FilterLogs(opts, "Permit", ownerRule, tokenRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &AbisPermitIterator{contract: _Abis.contract, event: "Permit", logs: logs, sub: sub}, nil
}

// WatchPermit is a free log subscription operation binding the contract event 0xc6a377bfc4eb120024a8ac08eef205be16b817020812c73223e81d1bdb9708ec.
//
// Solidity: event Permit(address indexed owner, address indexed token, address indexed spender, uint160 amount, uint48 expiration, uint48 nonce)
func (_Abis *AbisFilterer) WatchPermit(opts *bind.WatchOpts, sink chan<- *AbisPermit, owner []common.Address, token []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Abis.contract.WatchLogs(opts, "Permit", ownerRule, tokenRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AbisPermit)
				if err := _Abis.contract.UnpackLog(event, "Permit", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePermit is a log parse operation binding the contract event 0xc6a377bfc4eb120024a8ac08eef205be16b817020812c73223e81d1bdb9708ec.
//
// Solidity: event Permit(address indexed owner, address indexed token, address indexed spender, uint160 amount, uint48 expiration, uint48 nonce)
func (_Abis *AbisFilterer) ParsePermit(log types.Log) (*AbisPermit, error) {
	event := new(AbisPermit)
	if err := _Abis.contract.UnpackLog(event, "Permit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AbisUnorderedNonceInvalidationIterator is returned from FilterUnorderedNonceInvalidation and is used to iterate over the raw logs and unpacked data for UnorderedNonceInvalidation events raised by the Abis contract.
type AbisUnorderedNonceInvalidationIterator struct {
	Event *AbisUnorderedNonceInvalidation // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AbisUnorderedNonceInvalidationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AbisUnorderedNonceInvalidation)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AbisUnorderedNonceInvalidation)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AbisUnorderedNonceInvalidationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AbisUnorderedNonceInvalidationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AbisUnorderedNonceInvalidation represents a UnorderedNonceInvalidation event raised by the Abis contract.
type AbisUnorderedNonceInvalidation struct {
	Owner common.Address
	Word  *big.Int
	Mask  *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterUnorderedNonceInvalidation is a free log retrieval operation binding the contract event 0x3704902f963766a4e561bbaab6e6cdc1b1dd12f6e9e99648da8843b3f46b918d.
//
// Solidity: event UnorderedNonceInvalidation(address indexed owner, uint256 word, uint256 mask)
func (_Abis *AbisFilterer) FilterUnorderedNonceInvalidation(opts *bind.FilterOpts, owner []common.Address) (*AbisUnorderedNonceInvalidationIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _Abis.contract.FilterLogs(opts, "UnorderedNonceInvalidation", ownerRule)
	if err != nil {
		return nil, err
	}
	return &AbisUnorderedNonceInvalidationIterator{contract: _Abis.contract, event: "UnorderedNonceInvalidation", logs: logs, sub: sub}, nil
}

// WatchUnorderedNonceInvalidation is a free log subscription operation binding the contract event 0x3704902f963766a4e561bbaab6e6cdc1b1dd12f6e9e99648da8843b3f46b918d.
//
// Solidity: event UnorderedNonceInvalidation(address indexed owner, uint256 word, uint256 mask)
func (_Abis *AbisFilterer) WatchUnorderedNonceInvalidation(opts *bind.WatchOpts, sink chan<- *AbisUnorderedNonceInvalidation, owner []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _Abis.contract.WatchLogs(opts, "UnorderedNonceInvalidation", ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AbisUnorderedNonceInvalidation)
				if err := _Abis.contract.UnpackLog(event, "UnorderedNonceInvalidation", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnorderedNonceInvalidation is a log parse operation binding the contract event 0x3704902f963766a4e561bbaab6e6cdc1b1dd12f6e9e99648da8843b3f46b918d.
//
// Solidity: event UnorderedNonceInvalidation(address indexed owner, uint256 word, uint256 mask)
func (_Abis *AbisFilterer) ParseUnorderedNonceInvalidation(log types.Log) (*AbisUnorderedNonceInvalidation, error) {
	event := new(AbisUnorderedNonceInvalidation)
	if err := _Abis.contract.UnpackLog(event, "UnorderedNonceInvalidation", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
