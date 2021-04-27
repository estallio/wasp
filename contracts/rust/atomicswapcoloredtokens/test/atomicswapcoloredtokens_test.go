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

 // set up recipient account and mint some tokens to start atomic swap
 var recipient = chain.Env.NewSignatureSchemeWithFunds()
 recipientColor, err := chain.Env.MintTokens(recipient, 10)
 require.NoError(t, err)

 // get sure recipient got the right balances
 chain.Env.AssertAddressBalance(recipient.Address(), balance.ColorIOTA, solo.Saldo - 10)
 chain.Env.AssertAddressBalance(recipient.Address(), recipientColor, 10)

 // ************************************************************************************************
 // both accounts now have 10 tokens of different colors and want to exchange them
 // ************************************************************************************************

 // prepare sender's start swap call
 req := solo.NewCallParams(ScName, FuncStartSwap,
  ParamSwapId, "first-swap-between-colleagues",
  // who is the recipient
  ParamAddressRecipient, recipient.Address(),
  // sender wants following coins from recipient...
  ParamColorRecipient, recipientColor,
  ParamAmountRecipient, 10,
  // ...and gives following coins in exchange
  ParamColorSender, senderColor,
  ParamAmountSender, 10,
  // swap is open for 200 seconds
  ParamDuration, 200,
 ).WithTransfers(map[balance.Color]int64{
  senderColor: 10,
 })

 // send open swap call
 _, err = chain.PostRequestSync(req, sender)
 require.NoError(t, err)

 // sender should have 0 tokens of senderColor now
 chain.Env.AssertAddressBalance(sender.Address(), senderColor, 0)

 // prepare sender's cancel swap call
 req = solo.NewCallParams(ScName, FuncCancelSwap,
  ParamSwapId, "first-swap-between-colleagues",
 )

 // send cancel swap call
 _, err = chain.PostRequestSync(req, sender)
 require.NoError(t, err)

 // sender should now have the coins back
 chain.Env.AssertAddressBalance(sender.Address(), senderColor, 10)

 // ************************************************************************************************
 // except for the single tokens that are necessary for every SC-call, the state is the same
 // as at the beginning now
 // ************************************************************************************************

 // now deposit the colored coins again
 // prepare sender's start swap call
 req = solo.NewCallParams(ScName, FuncStartSwap,
  ParamSwapId, "second-swap-between-colleagues",
  // who is the recipient
  ParamAddressRecipient, recipient.Address(),
  // sender wants following coins from recipient...
  ParamColorRecipient, recipientColor,
  ParamAmountRecipient, 10,
  // ...and gives following coins in exchange
  ParamColorSender, senderColor,
  ParamAmountSender, 10,
  // swap is open for 200 seconds
  ParamDuration, 200,
 ).WithTransfers(map[balance.Color]int64{
  senderColor: 10,
 })

 // send sender's swap call
 _, err = chain.PostRequestSync(req, sender)
 require.NoError(t, err)

 // recipient now claims the tokens
 req = solo.NewCallParams(ScName, FuncFinalizeSwap,
  ParamSwapId, "second-swap-between-colleagues",
 ).WithTransfers(map[balance.Color]int64{
  recipientColor: 10,
 })

 // send token claim call
 _, err = chain.PostRequestSync(req, recipient)
 require.NoError(t, err)

 // sender should now have the recipient color and vice versa
 chain.Env.AssertAddressBalance(sender.Address(), recipientColor, 10)
 chain.Env.AssertAddressBalance(recipient.Address(), senderColor, 10)

 // ************************************************************************************************
 // sender and recipient have now exchanged their colored tokens
 // ************************************************************************************************
}
