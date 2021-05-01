// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/stretchr/testify/require"
	"testing"
)

// defining erc20 binary - must be copied into test folder
var (
	erc20file = "erc20_bg.wasm"
	erc20swapFile = "atomicswaperc20sync_bg.wasm"
	erc20A = "erc20ScName_A"
	erc20B = "erc20ScName_B"
)

func TestDeploy(t *testing.T) {
	// create new environment
	env := solo.New(t, false, false)


	// create new chain
	chain := env.NewChain(nil, "chain")


	// deploy erc20 first time on chain
	sender := env.NewSignatureSchemeWithFunds()
	senderAgentId := coretypes.NewAgentIDFromAddress(sender.Address())
	err := chain.DeployWasmContract(nil, erc20A, erc20file,
		erc20ParamSupply, 300,
		erc20ParamCreator, senderAgentId,
	)
	require.NoError(t, err)

	_, err = chain.FindContract(erc20A)
	require.NoError(t, err)


	// deploy erc20 second time on chain
	recipient := env.NewSignatureSchemeWithFunds()
	recipientAgentId := coretypes.NewAgentIDFromAddress(recipient.Address())
	err = chain.DeployWasmContract(nil, erc20B, erc20file,
		erc20ParamSupply, 300,
		erc20ParamCreator, recipientAgentId,
	)
	require.NoError(t, err)

	_, err = chain.FindContract(erc20B)
	require.NoError(t, err)


	// deploy swap contract on chain
	err = chain.DeployWasmContract(nil, ScName, erc20swapFile)
	require.NoError(t, err)

	_, err = chain.FindContract(ScName)
	require.NoError(t, err)


	// get agent id of contract C
	contractId := coretypes.NewContractID(chain.ChainID, coretypes.Hn(ScName))
	contractIdAsAgentId := coretypes.NewAgentIDFromContractID(contractId)


	// check the swap id contract, sender and recipient balances
	checkErc20BalanceOnSc(erc20A, chain, senderAgentId, 300)
	checkErc20BalanceOnSc(erc20A, chain, recipientAgentId, 0)
	checkErc20BalanceOnSc(erc20A, chain, contractIdAsAgentId, 0)
	checkErc20BalanceOnSc(erc20B, chain, senderAgentId, 0)
	checkErc20BalanceOnSc(erc20B, chain, recipientAgentId, 300)
	checkErc20BalanceOnSc(erc20B, chain, contractIdAsAgentId, 0)


	// client A allows the swap SC to transfer 100 erc20 tokens to the swap contract
	req := solo.NewCallParams(erc20A, erc20FuncApprove,
		erc20ParamDelegation, contractIdAsAgentId,
		erc20ParamAmount, 100,
	)

	_, err = chain.PostRequestSync(req, sender)
	require.NoError(t, err)


	// client A opens an atomic swap on the contract - this creates the swap and transfers the erc20 tokens
	req = solo.NewCallParams(ScName, FuncStartSwap,
		ParamSwapId, "first-swap-between-colleagues-cancel",
		// on which contract is the erc20 token stored
		ParamScNameSender, erc20A,
		// where is the other erc20 contract
		ParamScNameRecipient, erc20B,
		// which agentId depends to the recipient
		ParamAgentIdRecipient, recipientAgentId,
		// how many erc20 tokens the sender is trading in
		ParamAmountSender, 100,
		// desired amount of erc20 tokens
		ParamAmountRecipient, 100,
		// swap is open for 200 seconds
		ParamDuration, 200,
	)

	_, err = chain.PostRequestSync(req, sender)
	require.NoError(t, err)


	// check the swap id contract, sender and recipient balances
	checkErc20BalanceOnSc(erc20A, chain, senderAgentId, 200)
	checkErc20BalanceOnSc(erc20A, chain, recipientAgentId, 0)
	checkErc20BalanceOnSc(erc20A, chain, contractIdAsAgentId, 100)
	checkErc20BalanceOnSc(erc20B, chain, senderAgentId, 0)
	checkErc20BalanceOnSc(erc20B, chain, recipientAgentId, 300)
	checkErc20BalanceOnSc(erc20B, chain, contractIdAsAgentId, 0)

	// client A cancels the atomic swap on the contract
	req = solo.NewCallParams(ScName, FuncCancelSwap,
		ParamSwapId, "first-swap-between-colleagues-cancel",
	)

	_, err = chain.PostRequestSync(req, sender)
	require.NoError(t, err)


	// check the swap id contract, sender and recipient balances
	checkErc20BalanceOnSc(erc20A, chain, senderAgentId, 300)
	checkErc20BalanceOnSc(erc20A, chain, recipientAgentId, 0)
	checkErc20BalanceOnSc(erc20A, chain, contractIdAsAgentId, 0)
	checkErc20BalanceOnSc(erc20B, chain, senderAgentId, 0)
	checkErc20BalanceOnSc(erc20B, chain, recipientAgentId, 300)
	checkErc20BalanceOnSc(erc20B, chain, contractIdAsAgentId, 0)


	// ************************************************************************************************
	// the state is the same as at the beginning now
	// ************************************************************************************************


	// client A allows the swap SC to transfer 100 erc20 tokens to the swap contract again
	req = solo.NewCallParams(erc20A, erc20FuncApprove,
		erc20ParamDelegation, contractIdAsAgentId,
		erc20ParamAmount, 100,
	)

	_, err = chain.PostRequestSync(req, sender)
	require.NoError(t, err)


	// client A opens an atomic swap on the contract - this creates the swap and transfers the erc20 tokens
	req = solo.NewCallParams(ScName, FuncStartSwap,
		ParamSwapId, "first-swap-between-colleagues",
		// on which contract is the erc20 token stored
		ParamScNameSender, erc20A,
		// where is the other erc20 contract
		ParamScNameRecipient, erc20B,
		// which agentId depends to the recipient
		ParamAgentIdRecipient, recipientAgentId,
		// how many erc20 tokens the sender is trading in
		ParamAmountSender, 100,
		// desired amount of erc20 tokens
		ParamAmountRecipient, 100,
		// swap is open for 200 seconds
		ParamDuration, 200,
	)

	_, err = chain.PostRequestSync(req, sender)
	require.NoError(t, err)


	// check the swap id contract, sender and recipient balances
	checkErc20BalanceOnSc(erc20A, chain, senderAgentId, 200)
	checkErc20BalanceOnSc(erc20A, chain, recipientAgentId, 0)
	checkErc20BalanceOnSc(erc20A, chain, contractIdAsAgentId, 100)
	checkErc20BalanceOnSc(erc20B, chain, senderAgentId, 0)
	checkErc20BalanceOnSc(erc20B, chain, recipientAgentId, 300)
	checkErc20BalanceOnSc(erc20B, chain, contractIdAsAgentId, 0)


	// client B allows the swap SC to transfer 100 erc20 tokens to the swap contract
	req = solo.NewCallParams(erc20B, erc20FuncApprove,
		erc20ParamDelegation, contractIdAsAgentId,
		erc20ParamAmount, 100,
	)

	_, err = chain.PostRequestSync(req, recipient)
	require.NoError(t, err)


	// close swap
	req = solo.NewCallParams(ScName, FuncFinalizeSwap,
		ParamSwapId, "first-swap-between-colleagues",
	)

	_, err = chain.PostRequestSync(req, recipient)
	require.NoError(t, err)


	// check the swap id contract, sender and recipient balances
	checkErc20BalanceOnSc(erc20A, chain, senderAgentId, 200)
	checkErc20BalanceOnSc(erc20A, chain, recipientAgentId, 100)
	checkErc20BalanceOnSc(erc20A, chain, contractIdAsAgentId, 0)
	checkErc20BalanceOnSc(erc20B, chain, senderAgentId, 100)
	checkErc20BalanceOnSc(erc20B, chain, recipientAgentId, 200)
	checkErc20BalanceOnSc(erc20B, chain, contractIdAsAgentId, 0)
}

func checkErc20BalanceOnSc(ScName string, e *solo.Chain, account coretypes.AgentID, amount int64) {
	res, err := e.CallView(ScName, erc20ViewBalanceOf,
		erc20ParamAccount, account,
	)
	require.NoError(e.Env.T, err)
	sup, ok, err := codec.DecodeInt64(res.MustGet(erc20ParamAmount))
	require.NoError(e.Env.T, err)
	require.True(e.Env.T, ok)
	require.EqualValues(e.Env.T, sup, amount)
}
