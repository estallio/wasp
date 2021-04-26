// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	"github.com/iotaledger/wasp/contracts/common"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDeploy(t *testing.T) {
	// deploy contract
	chain := common.StartChainAndDeployWasmContractByName(t, ScName)

	// search contract and get sure it is available
	_, err := chain.FindContract(ScName)
	require.NoError(t, err)

	// set up sender account and mint some tokens to start atomic swap
	var sender = chain.Env.NewSignatureSchemeWithFunds()
	senderColor, err := chain.Env.MintTokens(sender, 10)
	require.NoError(t, err)

	// get sure sender got the right balances
	chain.Env.AssertAddressBalance(sender.Address(), balance.ColorIOTA, solo.Saldo - 10)
	chain.Env.AssertAddressBalance(sender.Address(), senderColor, 10)

	// same for receiver
	var receiver = chain.Env.NewSignatureSchemeWithFunds()
	receiverColor, err := chain.Env.MintTokens(receiver, 10)
	require.NoError(t, err)

	chain.Env.AssertAddressBalance(receiver.Address(), balance.ColorIOTA, solo.Saldo - 10)
	chain.Env.AssertAddressBalance(receiver.Address(), receiverColor, 10)

	// ************************************************************************************************
	// both accounts now have 10 tokens of different colors and want to exchange them
	// ************************************************************************************************

	// start swap by sender
	req := solo.NewCallParams(ScName, FuncStartSwap,
		ParamSwapId, "first-swap-between-colleagues",
		// sender wants following coins from receiver
		ParamAddressReceiver, receiver.Address(),
		ParamColorReceiver, receiverColor,
		ParamAmountReceiver, 10,
		// ... and gives following coins to the swap
		ParamColorSender, senderColor,
		ParamAmountSender, 10,
		// swap is open for 200 seconds
		ParamDuration, 200,
	).WithTransfers(map[balance.Color]int64{
		senderColor: 10,
	})

	_, err = chain.PostRequestSync(req, sender)
	require.NoError(t, err)

	// sender should have 0 tokens of senderColor now
	chain.Env.AssertAddressBalance(sender.Address(), senderColor, 0)

	// cancel swap by sender
	req = solo.NewCallParams(ScName, FuncCancelSwap,
		ParamSwapId, "first-swap-between-colleagues",
	)

	_, err = chain.PostRequestSync(req, sender)
	require.NoError(t, err)

	// sender should now have the coins back
	chain.Env.AssertAddressBalance(sender.Address(), senderColor, 10)

	// ************************************************************************************************
	// except for the single tokens that are necessary for every SC-call, the state is the same
	// as at the beginning now
	// ************************************************************************************************

	// now deposit the colored coins again
	// start swap by sender
	req = solo.NewCallParams(ScName, FuncStartSwap,
		ParamSwapId, "second-swap-between-colleagues",
		// sender wants following coins from receiver
		ParamAddressReceiver, receiver.Address(),
		ParamColorReceiver, receiverColor,
		ParamAmountReceiver, 10,
		// ... and gives following coins to the swap
		ParamColorSender, senderColor,
		ParamAmountSender, 10,
		// swap is open for 200 seconds
		ParamDuration, 200,
	).WithTransfers(map[balance.Color]int64{
		senderColor: 10,
	})

	_, err = chain.PostRequestSync(req, sender)
	require.NoError(t, err)

	// receiver now claims the tokens
	req = solo.NewCallParams(ScName, FuncFinalizeSwap,
		ParamSwapId, "second-swap-between-colleagues",
	).WithTransfers(map[balance.Color]int64{
		receiverColor: 10,
	})

	_, err = chain.PostRequestSync(req, receiver)
	require.NoError(t, err)

	// sender should now have the receiver color and vice versa
	chain.Env.AssertAddressBalance(sender.Address(), receiverColor, 10)
	chain.Env.AssertAddressBalance(receiver.Address(), senderColor, 10)

	// ************************************************************************************************
	// sender and receiver have now exchanged their colored tokens
	// ************************************************************************************************
}
