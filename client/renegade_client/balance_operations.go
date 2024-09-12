package client

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"renegade.fi/golang-sdk/abis"
	"renegade.fi/golang-sdk/client/api_types"
	"renegade.fi/golang-sdk/wallet"
)

// Deposit deposits funds into the wallet
func (c *RenegadeClient) Deposit(mint *string, amount *big.Int, ethPrivateKey *ecdsa.PrivateKey) (*api_types.DepositResponse, error) {
	// Get the back of the queue wallet
	apiWallet, err := c.GetBackOfQueueWallet()
	if err != nil {
		return nil, err
	}

	// Convert the API wallet to a wallet
	backOfQueueWallet, err := apiWallet.ToWallet()
	if err != nil {
		return nil, err
	}

	// Add the balance to the wallet
	bal := wallet.NewBalanceBuilder().WithMintHex(*mint).WithAmountBigInt(amount).Build()
	err = backOfQueueWallet.AddBalance(bal)
	if err != nil {
		return nil, err
	}
	backOfQueueWallet.Reblind()

	// Approve Permit2 contract to spend the deposited amount
	req, err := c.setupDeposit(*mint, amount, ethPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to setup deposit: %w", err)
	}

	// Get the wallet update auth
	auth, err := getWalletUpdateAuth(backOfQueueWallet)
	if err != nil {
		return nil, err
	}
	req.WalletUpdateAuthorization = *auth

	// Post the deposit to the relayer
	walletId := c.walletSecrets.Id
	path := api_types.BuildDepositPath(walletId)

	resp := api_types.DepositResponse{}
	err = c.httpClient.PostWithAuth(path, req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to post deposit request: %w", err)
	}

	return &resp, nil
}

// setupDeposit sets up the deposit request, this includes approving the Permit2 contract, and generating the witness and signature
func (c *RenegadeClient) setupDeposit(mint string, amount *big.Int, ethPrivateKey *ecdsa.PrivateKey) (*api_types.DepositRequest, error) {
	// Approve the Permit2 contract to spend the balance
	err := c.approvePermit2Deposit(mint, amount, ethPrivateKey)
	if err != nil {
		return nil, err
	}

	// Generate the witness and signature for the permit
	witness, signature, err := c.generatePermit2Signature(mint, amount, ethPrivateKey)
	if err != nil {
		return nil, err
	}

	// Create the deposit request
	fromAddr := crypto.PubkeyToAddress(ethPrivateKey.PublicKey).Hex()
	sig := base64.RawStdEncoding.EncodeToString(signature)

	return &api_types.DepositRequest{
		FromAddr:        fromAddr,
		Mint:            mint,
		Amount:          amount.String(),
		PermitNonce:     witness.Nonce.String(),
		PermitDeadline:  witness.Deadline.String(),
		PermitSignature: sig,
	}, nil
}

// Withdraw withdraws funds from the wallet to the address for the given private key
func (c *RenegadeClient) Withdraw(mint string, amount *big.Int, ethPrivateKey *ecdsa.PrivateKey) (*api_types.WithdrawResponse, error) {
	addr := hex.EncodeToString(crypto.PubkeyToAddress(ethPrivateKey.PublicKey).Bytes())
	return c.WithdrawToAddress(mint, amount, &addr)
}

// WithdrawToAddress withdraws funds from the wallet to the given address
func (c *RenegadeClient) WithdrawToAddress(mint string, amount *big.Int, destination *string) (*api_types.WithdrawResponse, error) {
	// Get the back of the queue wallet
	apiWallet, err := c.GetBackOfQueueWallet()
	if err != nil {
		return nil, err
	}

	// Convert the API wallet to a wallet
	backOfQueueWallet, err := apiWallet.ToWallet()
	if err != nil {
		return nil, err
	}

	// Remove the balance from the wallet
	bal := wallet.NewBalanceBuilder().WithMintHex(mint).WithAmountBigInt(amount).Build()
	err = backOfQueueWallet.RemoveBalance(bal)
	if err != nil {
		return nil, err
	}
	backOfQueueWallet.Reblind()

	// Get the wallet update auth
	auth, err := getWalletUpdateAuth(backOfQueueWallet)
	if err != nil {
		return nil, err
	}

	// Get the external transfer signature
	// Construct the external transfer signature

	externalTransferSig, err := c.generateWithdrawalSignature(mint, amount, destination)
	if err != nil {
		return nil, fmt.Errorf("failed to generate external transfer signature: %w", err)
	}

	// Create the withdraw request
	req := &api_types.WithdrawRequest{
		DestinationAddr:           *destination,
		Amount:                    amount.String(),
		ExternalTransferSig:       externalTransferSig,
		WalletUpdateAuthorization: *auth,
	}

	// Post the request to the relayer
	path := api_types.BuildWithdrawPath(c.walletSecrets.Id, mint)
	var resp api_types.WithdrawResponse
	err = c.httpClient.PostWithAuth(path, req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to post withdraw request: %w", err)
	}

	return &resp, nil
}

// --- Helpers --- //

// approvePermit2Deposit approves the Permit2 contract to spend the deposited amount
func (c *RenegadeClient) approvePermit2Deposit(mint string, amount *big.Int, ethPrivateKey *ecdsa.PrivateKey) error {
	// Create an RPC client
	rpcClient, err := c.createRpcClient()
	if err != nil {
		return fmt.Errorf("failed to create RPC client: %w", err)
	}

	// Create a transactor
	auth, err := c.createTransactor(ethPrivateKey)
	if err != nil {
		return err
	}

	// Get the ERC20 contract
	erc20Contract, err := abis.NewContracts(common.HexToAddress(mint), rpcClient)
	if err != nil {
		return fmt.Errorf("failed to create ERC20 contract: %w", err)
	}

	// Check the existing balance
	bal, err := erc20Contract.BalanceOf(&bind.CallOpts{}, auth.From)
	if err != nil {
		return fmt.Errorf("failed to get balance: %w", err)
	}

	if bal.Cmp(amount) < 0 {
		return fmt.Errorf("insufficient balance for deposit: have %s, need %s", bal.String(), amount.String())
	}

	// Check existing allowance
	// If allowance is sufficient, no need for a new approval
	permit2Addr := common.HexToAddress(c.chainConfig.Permit2Address)
	allowance, err := erc20Contract.Allowance(&bind.CallOpts{}, auth.From, permit2Addr)
	if err != nil {
		return fmt.Errorf("failed to get allowance: %w", err)
	}

	if allowance.Cmp(amount) >= 0 {
		log.Printf("Existing allowance (%s) is sufficient for the deposit amount (%s)", allowance.String(), amount.String())
		return nil
	}

	// Approve the Permit2 contract to spend the balance
	log.Printf("Existing allowance (%s) is insufficient. Approving Permit2 contract to spend %s tokens", allowance.String(), amount.String())
	tx, err := erc20Contract.Approve(auth, permit2Addr, amount)
	if err != nil {
		return fmt.Errorf("failed to approve Permit2 contract: %w", err)
	}

	receipt, err := bind.WaitMined(context.Background(), rpcClient, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for approval transaction: %w", err)
	}
	log.Printf("Approval transaction hash: %s", receipt.TxHash.Hex())

	return nil
}

// generatePermit2Signature generates a Permit2 signature for the deposit
func (c *RenegadeClient) generatePermit2Signature(mint string, amount *big.Int, ethPrivateKey *ecdsa.PrivateKey) (*PermitWitnessTransferFrom, []byte, error) {
	// Construct the EIP712 domain
	permit2Address := common.HexToAddress(c.chainConfig.Permit2Address)
	chainId := big.NewInt(int64(c.chainConfig.ChainID))
	domain := ConstructEIP712Domain(chainId, permit2Address)

	// Create the TokenPermissions struct
	tokenPermissions := abis.ISignatureTransferTokenPermissions{
		Token:  common.HexToAddress(mint),
		Amount: amount,
	}

	// Generate nonce and deadline
	nonce, err := randomU256()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate nonce: %w", err)
	}
	deadline := new(big.Int).SetUint64(^uint64(0))

	// Generate a random witness (replace this with actual witness generation if needed)
	witness, err := c.getPermitWitness()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate witness: %w", err)
	}

	// Create the PermitWitnessTransferFrom struct
	permitWitnessTransferFrom := PermitWitnessTransferFrom{
		Permitted: tokenPermissions,
		Spender:   common.HexToAddress(c.chainConfig.DarkpoolAddress),
		Nonce:     nonce,
		Deadline:  deadline,
		Witness:   witness,
	}

	// Generate the signing hash
	signingHash, err := getPermitSigningHash(permitWitnessTransferFrom, domain)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get signing hash: %w", err)
	}

	// Sign the hash
	signature, err := crypto.Sign(signingHash.Bytes(), ethPrivateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to sign permit: %w", err)
	}

	// Add 27 to the last byte of the signature, we expect the bitcoin style replay protection
	signature[len(signature)-1] += 27
	return &permitWitnessTransferFrom, signature, nil
}

// generateWithdrawalSignature generates a signature for the withdrawal
func (c *RenegadeClient) generateWithdrawalSignature(mint string, amount *big.Int, destination *string) (*string, error) {
	rootKey := ecdsa.PrivateKey(*c.walletSecrets.Keychain.SkRoot())
	sigBytes, err := postcardSerializeTransfer(mint, amount, destination)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize transfer: %w", err)
	}

	// Hash and sign
	digest := crypto.Keccak256(sigBytes)
	signature, err := crypto.Sign(digest, &rootKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign withdrawal: %w", err)
	}

	sig := base64.RawStdEncoding.EncodeToString(signature)
	return &sig, nil
}

// randomU256 generates a random 256-bit unsigned integer
func randomU256() (*big.Int, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	nonceBig := new(big.Int).SetBytes(randomBytes)
	return nonceBig, nil
}

// getPermitWitness generates a witness for the permit
func (c *RenegadeClient) getPermitWitness() (*DepositWitness, error) {
	pkRoot := c.walletSecrets.Keychain.PublicKeys.PkRoot
	scalars, err := wallet.ToScalarsRecursive(&pkRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to convert pkRoot to scalars: %w", err)
	}

	// Convert the scalars to big.Ints
	rootValues := [4]*big.Int{
		scalars[0].ToBigInt(),
		scalars[1].ToBigInt(),
		scalars[2].ToBigInt(),
		scalars[3].ToBigInt(),
	}

	return &DepositWitness{
		PkRoot: rootValues,
	}, nil
}
