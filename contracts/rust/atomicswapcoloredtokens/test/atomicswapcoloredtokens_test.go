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

	chain.Env.AssertAddressBalance(sender.Address(), balance.ColorIOTA, solo.Saldo - 10)
	chain.Env.AssertAddressBalance(sender.Address(), senderColor, 10)

	// same for receiver
	var receiver = chain.Env.NewSignatureSchemeWithFunds()
	receiverColor, err := chain.Env.MintTokens(receiver, 10)
	require.NoError(t, err)

	chain.Env.AssertAddressBalance(receiver.Address(), balance.ColorIOTA, solo.Saldo - 10)
	chain.Env.AssertAddressBalance(receiver.Address(), receiverColor, 10)

	// both accounts now have 10 tokens of different colors and want to exchange them

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

	// sender should have no tokens of senderColor now
	chain.Env.AssertAddressBalance(sender.Address(), senderColor, 0)

	// cancel swap by sender
	req = solo.NewCallParams(ScName, FuncCancelSwap,
		ParamSwapId, "first-swap-between-colleagues",
	)

	_, err = chain.PostRequestSync(req, sender)
	require.NoError(t, err)

	// sender should now have the coins back
	chain.Env.AssertAddressBalance(sender.Address(), senderColor, 10)

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

	// claim the tokens by the receiver
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
}

/*
var auctioneer signaturescheme.SignatureScheme
var tokenColor balance.Color

func setupTest(t *testing.T) *solo.Chain {
	chain := common.StartChainAndDeployWasmContractByName(t, ScName)

	// set up auctioneer account and mint some tokens to auction off
	auctioneer = chain.Env.NewSignatureSchemeWithFunds()
	newColor, err := chain.Env.MintTokens(auctioneer, 10)
	require.NoError(t, err)
	chain.Env.AssertAddressBalance(auctioneer.Address(), balance.ColorIOTA, solo.Saldo-10)
	chain.Env.AssertAddressBalance(auctioneer.Address(), newColor, 10)
	tokenColor = newColor

	// start auction
	req := solo.NewCallParams(ScName, FuncStartAuction,
		ParamColor, tokenColor,
		ParamMinimumBid, 500,
		ParamDescription, "Cool tokens for sale!",
	).WithTransfers(map[balance.Color]int64{
		balance.ColorIOTA: 25, // deposit, must be >=minimum*margin
		tokenColor:        10, // the tokens to auction
	})
	_, err = chain.PostRequestSync(req, auctioneer)
	require.NoError(t, err)
	return chain
}

func TestFaStartAuction(t *testing.T) {
	chain := setupTest(t)

	// note 1 iota should be stuck in the delayed finalize_auction
	chain.AssertAccountBalance(common.ContractAccount, balance.ColorIOTA, 25-1)
	chain.AssertAccountBalance(common.ContractAccount, tokenColor, 10)

	// auctioneer sent 25 deposit + 10 tokenColor + used 1 for request
	chain.Env.AssertAddressBalance(auctioneer.Address(), balance.ColorIOTA, solo.Saldo-35-1)
	// 1 used for request was sent back to auctioneer's account on chain
	account := coretypes.NewAgentIDFromSigScheme(auctioneer)
	chain.AssertAccountBalance(account, balance.ColorIOTA, 1)
}

func TestFaAuctionInfo(t *testing.T) {
	chain := setupTest(t)

	res, err := chain.CallView(
		ScName, ViewGetInfo,
		ParamColor, tokenColor,
	)
	require.NoError(t, err)
	account := coretypes.NewAgentIDFromSigScheme(auctioneer)
	requireAgent(t, res, VarCreator, account)
	requireInt64(t, res, VarBidders, 0)
}

func TestFaNoBids(t *testing.T) {
	chain := setupTest(t)

	// wait for finalize_auction
	chain.Env.AdvanceClockBy(61 * time.Minute)
	chain.WaitForEmptyBacklog()

	res, err := chain.CallView(
		ScName, ViewGetInfo,
		ParamColor, tokenColor,
	)
	require.NoError(t, err)
	requireInt64(t, res, VarBidders, 0)
}

func TestFaOneBidTooLow(t *testing.T) {
	chain := setupTest(t)

	req := solo.NewCallParams(ScName, FuncPlaceBid,
		ParamColor, tokenColor,
	).WithTransfer(balance.ColorIOTA, 100)
	_, err := chain.PostRequestSync(req, auctioneer)
	require.Error(t, err)

	// wait for finalize_auction
	chain.Env.AdvanceClockBy(61 * time.Minute)
	chain.WaitForEmptyBacklog()

	res, err := chain.CallView(
		ScName, ViewGetInfo,
		ParamColor, tokenColor,
	)
	require.NoError(t, err)
	requireInt64(t, res, VarHighestBid, -1)
	requireInt64(t, res, VarBidders, 0)
}

func TestFaOneBid(t *testing.T) {
	chain := setupTest(t)

	bidder := chain.Env.NewSignatureSchemeWithFunds()
	req := solo.NewCallParams(ScName, FuncPlaceBid,
		ParamColor, tokenColor,
	).WithTransfer(balance.ColorIOTA, 500)
	_, err := chain.PostRequestSync(req, bidder)
	require.NoError(t, err)

	// wait for finalize_auction
	chain.Env.AdvanceClockBy(61 * time.Minute)
	chain.WaitForEmptyBacklog()

	res, err := chain.CallView(
		ScName, ViewGetInfo,
		ParamColor, tokenColor,
	)
	require.NoError(t, err)
	requireInt64(t, res, VarBidders, 1)
	requireInt64(t, res, VarHighestBid, 500)
	requireAgent(t, res, VarHighestBidder, coretypes.NewAgentIDFromSigScheme(bidder))
}

func requireAgent(t *testing.T, res dict.Dict, key string, expected coretypes.AgentID) {
	actual, exists, err := codec.DecodeAgentID(res.MustGet(kv.Key(key)))
	require.NoError(t, err)
	require.True(t, exists)
	require.EqualValues(t, expected, actual)
}

func requireInt64(t *testing.T, res dict.Dict, key string, expected int64) {
	actual, exists, err := codec.DecodeInt64(res.MustGet(kv.Key(key)))
	require.NoError(t, err)
	require.True(t, exists)
	require.EqualValues(t, expected, actual)
}

func requireString(t *testing.T, res dict.Dict, key string, expected string) {
	actual, exists, err := codec.DecodeString(res.MustGet(kv.Key(key)))
	require.NoError(t, err)
	require.True(t, exists)
	require.EqualValues(t, expected, actual)
}
*/
