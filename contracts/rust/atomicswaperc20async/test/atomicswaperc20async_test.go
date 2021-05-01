// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// defining erc20 binary - must be copied into test folder
var (
	erc20file = "erc20_bg.wasm"
	erc20swapFile = "atomicswaperc20async_bg.wasm"
	erc20A = "erc20ScName_A"
	erc20B = "erc20ScName_B"
)

func TestDeploy(t *testing.T) {
	// create new environment
	env := solo.New(t, false, false)


	// create new chains
	chainA := env.NewChain(nil, "chain_A")
	chainB := env.NewChain(nil, "chain_B")


	// deploy erc20 first time on chain A
	sender := env.NewSignatureSchemeWithFunds()
	senderAgentId := coretypes.NewAgentIDFromAddress(sender.Address())
	err := chainA.DeployWasmContract(nil, erc20A, erc20file,
		erc20ParamSupply, 300,
		erc20ParamCreator, senderAgentId,
	)
	require.NoError(t, err)

	_, err = chainA.FindContract(erc20A)
	require.NoError(t, err)


	// deploy erc20 second time on chain B
	recipient := env.NewSignatureSchemeWithFunds()
	recipientAgentId := coretypes.NewAgentIDFromAddress(recipient.Address())
	err = chainB.DeployWasmContract(nil, erc20B, erc20file,
		erc20ParamSupply, 300,
		erc20ParamCreator, recipientAgentId,
	)
	require.NoError(t, err)

	_, err = chainB.FindContract(erc20B)
	require.NoError(t, err)


	// deploy swap contract on chain A
	err = chainA.DeployWasmContract(nil, ScName, erc20swapFile)
	require.NoError(t, err)

	_, err = chainA.FindContract(ScName)
	require.NoError(t, err)


	// deploy swap contract on chain B
	err = chainB.DeployWasmContract(nil, ScName, erc20swapFile)
	require.NoError(t, err)

	_, err = chainB.FindContract(ScName)
	require.NoError(t, err)


	// define secret locking key
	secretString := "secret_key"
	secretBytes := []byte(secretString)
	secretHash := hashing.HashDataBlake2b(secretBytes)


	// get agent id of swap contract on chain A
	contractAId := coretypes.NewContractID(chainA.ChainID, coretypes.Hn(ScName))
	contractAIdAsAgentId := coretypes.NewAgentIDFromContractID(contractAId)


	// get agent id of swap contract on chain B
	contractBId := coretypes.NewContractID(chainB.ChainID, coretypes.Hn(ScName))
	contractBIdAsAgentId := coretypes.NewAgentIDFromContractID(contractBId)


	// check the swap id contract, sender and recipient balances
	checkErc20BalanceOnSc(erc20A, chainA, senderAgentId, 300)
	checkErc20BalanceOnSc(erc20A, chainA, recipientAgentId, 0)
	checkErc20BalanceOnSc(erc20A, chainA, contractAIdAsAgentId, 0)
	checkErc20BalanceOnSc(erc20B, chainB, senderAgentId, 0)
	checkErc20BalanceOnSc(erc20B, chainB, recipientAgentId, 300)
	checkErc20BalanceOnSc(erc20B, chainB, contractBIdAsAgentId, 0)


	// client A allows the swap SC on chain A to transfer 100 erc20 tokens to the swap contract
	req := solo.NewCallParams(erc20A, erc20FuncApprove,
		erc20ParamDelegation, contractAIdAsAgentId,
		erc20ParamAmount, 100,
	)

	_, err = chainA.PostRequestSync(req, sender)
	require.NoError(t, err)


	// client A opens an atomic swap on the contract - this creates the swap and transfers the erc20 tokens
	req = solo.NewCallParams(ScName, FuncStartSwap,
		ParamSwapId, "first-swap-between-colleagues-cancel",
		ParamKeyHash, secretHash,
		// on which contract is the erc20 token stored
		ParamScNameSender, erc20A,
		// which agentId depends to the recipient
		ParamAgentIdRecipient, recipientAgentId,
		// how many erc20 tokens the sender is trading in
		ParamAmountSender, 100,
		// swap is open for 200 seconds
		ParamDuration, 200,
	)

	_, err = chainA.PostRequestSync(req, sender)
	require.NoError(t, err)


	// check the swap id contract, sender and recipient balances
	checkErc20BalanceOnSc(erc20A, chainA, senderAgentId, 200)
	checkErc20BalanceOnSc(erc20A, chainA, recipientAgentId, 0)
	checkErc20BalanceOnSc(erc20A, chainA, contractAIdAsAgentId, 100)
	checkErc20BalanceOnSc(erc20B, chainB, senderAgentId, 0)
	checkErc20BalanceOnSc(erc20B, chainB, recipientAgentId, 300)
	checkErc20BalanceOnSc(erc20B, chainB, contractBIdAsAgentId, 0)


	// atomic swap must have been open for at least 200 seconds (greater than 200s)
	// to get the funds back, fast forward in time to +201 seconds
	chainA.Env.AdvanceClockBy(201 * time.Second)


	// client A cancels the atomic swap on the contract
	req = solo.NewCallParams(ScName, FuncCancelSwap,
		ParamSwapId, "first-swap-between-colleagues-cancel",
	)

	_, err = chainA.PostRequestSync(req, sender)
	require.NoError(t, err)


	// check the swap id contract, sender and recipient balances
	checkErc20BalanceOnSc(erc20A, chainA, senderAgentId, 300)
	checkErc20BalanceOnSc(erc20A, chainA, recipientAgentId, 0)
	checkErc20BalanceOnSc(erc20A, chainA, contractAIdAsAgentId, 0)
	checkErc20BalanceOnSc(erc20B, chainB, senderAgentId, 0)
	checkErc20BalanceOnSc(erc20B, chainB, recipientAgentId, 300)
	checkErc20BalanceOnSc(erc20B, chainB, contractBIdAsAgentId, 0)


	// ************************************************************************************************
	// the state is the same as at the beginning now
	// ************************************************************************************************


	// client A allows the swap SC on chain A to transfer 100 erc20 tokens to the swap contract
	req = solo.NewCallParams(erc20A, erc20FuncApprove,
		erc20ParamDelegation, contractAIdAsAgentId,
		erc20ParamAmount, 100,
	)

	_, err = chainA.PostRequestSync(req, sender)
	require.NoError(t, err)


	// client A opens an atomic swap on the contract - this creates the swap and transfers the erc20 tokens
	req = solo.NewCallParams(ScName, FuncStartSwap,
		ParamSwapId, "first-swap-between-colleagues",
		ParamKeyHash, secretHash,
		// on which contract is the erc20 token stored
		ParamScNameSender, erc20A,
		// which agentId depends to the recipient
		ParamAgentIdRecipient, recipientAgentId,
		// how many erc20 tokens the sender is trading in
		ParamAmountSender, 100,
		// swap is open for 200 seconds
		ParamDuration, 200,
	)

	_, err = chainA.PostRequestSync(req, sender)
	require.NoError(t, err)


	// check the swap id contract, sender and recipient balances
	checkErc20BalanceOnSc(erc20A, chainA, senderAgentId, 200)
	checkErc20BalanceOnSc(erc20A, chainA, recipientAgentId, 0)
	checkErc20BalanceOnSc(erc20A, chainA, contractAIdAsAgentId, 100)
	checkErc20BalanceOnSc(erc20B, chainB, senderAgentId, 0)
	checkErc20BalanceOnSc(erc20B, chainB, recipientAgentId, 300)
	checkErc20BalanceOnSc(erc20B, chainB, contractBIdAsAgentId, 0)


	// client B allows the swap SC to transfer 100 erc20 tokens to the swap contract
	req = solo.NewCallParams(erc20B, erc20FuncApprove,
		erc20ParamDelegation, contractBIdAsAgentId,
		erc20ParamAmount, 100,
	)

	_, err = chainB.PostRequestSync(req, recipient)
	require.NoError(t, err)


	// client B opens an atomic swap on the contract - this creates the swap and transfers the erc20 tokens
	req = solo.NewCallParams(ScName, FuncStartSwap,
		ParamSwapId, "first-swap-between-colleagues",
		ParamKeyHash, secretHash,
		// on which contract is the erc20 token stored
		ParamScNameSender, erc20B,
		// which agentId depends to the recipient
		ParamAgentIdRecipient, senderAgentId,
		// how many erc20 tokens the sender is trading in
		ParamAmountSender, 100,
		// swap is open for 200 seconds
		ParamDuration, 200,
	)

	_, err = chainB.PostRequestSync(req, recipient)
	require.NoError(t, err)


	// check the swap id contract, sender and recipient balances
	checkErc20BalanceOnSc(erc20A, chainA, senderAgentId, 200)
	checkErc20BalanceOnSc(erc20A, chainA, recipientAgentId, 0)
	checkErc20BalanceOnSc(erc20A, chainA, contractAIdAsAgentId, 100)
	checkErc20BalanceOnSc(erc20B, chainB, senderAgentId, 0)
	checkErc20BalanceOnSc(erc20B, chainB, recipientAgentId, 200)
	checkErc20BalanceOnSc(erc20B, chainB, contractBIdAsAgentId, 100)


	// A can now close swap on B
	req = solo.NewCallParams(ScName, FuncFinalizeSwap,
		ParamSwapId, "first-swap-between-colleagues",
		ParamKeySecret, secretString,
	)

	_, err = chainB.PostRequestSync(req, sender)
	require.NoError(t, err)


	// B can now close swap on A
	req = solo.NewCallParams(ScName, FuncFinalizeSwap,
		ParamSwapId, "first-swap-between-colleagues",
		ParamKeySecret, secretString,
	)

	_, err = chainA.PostRequestSync(req, recipient)
	require.NoError(t, err)


	// check the swap id contract, sender and recipient balances
	checkErc20BalanceOnSc(erc20A, chainA, senderAgentId, 200)
	checkErc20BalanceOnSc(erc20A, chainA, recipientAgentId, 100)
	checkErc20BalanceOnSc(erc20A, chainA, contractAIdAsAgentId, 0)
	checkErc20BalanceOnSc(erc20B, chainB, senderAgentId, 100)
	checkErc20BalanceOnSc(erc20B, chainB, recipientAgentId, 200)
	checkErc20BalanceOnSc(erc20B, chainB, contractBIdAsAgentId, 0)
}

func checkErc20BalanceOnSc(ScName string, e *solo.Chain, account coretypes.AgentID, amount int64) {
	res, err := e.CallView(ScName, erc20ViewBalanceOf,
		erc20ParamAccount, account,
	)
	require.NoError(e.Env.T, err)
	sup, ok, err := codec.DecodeInt64(res.MustGet(erc20ParamAmount))
	require.NoError(e.Env.T, err)
	require.True(e.Env.T, ok)
	require.EqualValues(e.Env.T, amount, sup)
}
