// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

use wasmlib::*;

use crate::*;
use crate::types::*;

pub fn func_start_swap(ctx: &ScFuncContext) {
    let p = ctx.params();

    // get the id where the swap should be accessible
    // currently a string is fine as map-key, maybe a hash or simple bytes suite better at a later point
    let param_swap_id = p.get_string(PARAM_SWAP_ID);

    // next 2 params define where the erc20 contracts are accessible
    let param_sc_name_sender = p.get_string(PARAM_SC_NAME_SENDER);
    let param_sc_name_recipient = p.get_string(PARAM_SC_NAME_RECIPIENT);

    // next 2 params define the amount of tokens that should be exchanged
    let param_amount_sender = p.get_int64(PARAM_AMOUNT_SENDER);
    let param_amount_recipient = p.get_int64(PARAM_AMOUNT_RECIPIENT);

    // get the agent ids of sender and recipient
    let param_agent_id_recipient = p.get_agent_id(PARAM_AGENT_ID_RECIPIENT);
    let param_agent_id_sender = ctx.caller().address().as_agent_id();

    //
    let param_duration = p.get_int64(PARAM_DURATION);

    ctx.log("checking params now...");

    // check if all necessary variables are set - only for unfortunate reasons when posting the tx
    ctx.require(param_amount_sender.exists(), "missing mandatory sender amount");
    ctx.require(param_amount_recipient.exists(), "missing mandatory recipient amount");
    ctx.require(param_swap_id.exists(), "missing swap id");
    ctx.require(param_sc_name_sender.exists(), "missing name of sender erc20 smart contract");
    ctx.require(param_sc_name_recipient.exists(), "missing name of recipient erc20 smart contract");
    ctx.require(param_agent_id_recipient.exists(), "missing agent id of recipient");

    ctx.log("checking amount now...");

    // check that the amount of tokens is sufficient
    let allowance_params = ScMutableMap::new();
    allowance_params.get_agent_id(ERC20_PARAM_ACCOUNT).set_value(&param_agent_id_sender);
    allowance_params.get_agent_id(ERC20_PARAM_DELEGATION).set_value(&ctx.contract_id().as_agent_id());

    let allowance_result_map = ctx.call(
        ScHname::new(param_sc_name_sender.value().as_str()),
        ScHname::new(ERC20_VIEW_ALLOWANCE),
        Some(allowance_params),
        None
    );

    ctx.require(allowance_result_map.get_int64(ERC20_PARAM_AMOUNT).value() >= param_amount_sender.value(), "contract is not allowed to transfer the specified amount of tokens");

    // prepare transfer of erc20 tokens
    let transfer_from_params = ScMutableMap::new();
    transfer_from_params.get_agent_id(ERC20_PARAM_ACCOUNT).set_value(&param_agent_id_sender);
    transfer_from_params.get_agent_id(ERC20_PARAM_RECIPIENT).set_value(&ctx.contract_id().as_agent_id());
    transfer_from_params.get_int64(ERC20_PARAM_AMOUNT).set_value(param_amount_sender.value());

    // transfer tokens - errors are thrown in the erc20 contract
    ctx.call(
        ScHname::new(param_sc_name_sender.value().as_str()),
        ScHname::new(ERC20_FUNC_TRANSFER_FROM),
        Some(transfer_from_params),
        None
    );

    // get the state
    let state: ScMutableMap = ctx.state();
    // get the atomic swap map
    let atomic_swaps = state.get_map(VAR_ATOMIC_SWAPS);
    // get the swap with swap_id
    let swap = atomic_swaps.get_bytes(&param_swap_id.value());

    ctx.log("creating swap now...");

    // swap id already busy
    ctx.require(!swap.exists(), "swap id already exists");

    // create an atomic swap object to save it in our register
    let atomic_swap = AtomicSwap {
        erc20_sc_name_sender: param_sc_name_sender.value(),
        erc20_sc_name_recipient: param_sc_name_recipient.value(),
        amount_sender: param_amount_sender.value(),
        amount_recipient: param_amount_recipient.value(),
        agent_id_recipient: param_agent_id_recipient.value(),
        agent_id_sender: param_agent_id_sender,
        duration_open: param_duration.value(),
        when_started: ctx.timestamp() / 1000000000,
        // for now, a number is sufficient to model the state, maybe there is any boolean or enum support sometime to model such states
        finished: 0,
    };

    ctx.log("saving atomic swap now...");

    swap.set_value(&atomic_swap.to_bytes());
}

pub fn func_cancel_swap(ctx: &ScFuncContext) {
    let p = ctx.params();

    // get the id where the swap should be accessible
    let param_swap_id = p.get_string(PARAM_SWAP_ID);

    ctx.require(param_swap_id.exists(), "missing mandatory swap id");

    // get the state
    let state: ScMutableMap = ctx.state();
    // get the atomic swap map
    let atomic_swaps = state.get_map(VAR_ATOMIC_SWAPS);
    // get the swap with swap_id
    let swap = atomic_swaps.get_bytes(&param_swap_id.value());

    // check if swap id exists
    ctx.require(swap.exists(), "swap id does not exists");

    // parse atomic swap
    let mut atomic_swap = AtomicSwap::from_bytes(&swap.value());

    // check if atomic swap is already finished
    ctx.require(atomic_swap.finished == 0, "swap is already finished");

    // check if the caller of this method is the sender
    ctx.require(ctx.caller().address().as_agent_id() == atomic_swap.agent_id_sender, "only the sender is able to cancel the swap");

    // transfer erc20 tokens from this contract back to sender's agent id
    let transfer_from_params = ScMutableMap::new();
    transfer_from_params.get_agent_id(ERC20_PARAM_ACCOUNT).set_value(&atomic_swap.agent_id_sender);
    transfer_from_params.get_int64(ERC20_PARAM_AMOUNT).set_value(atomic_swap.amount_sender);

    // call transfer contract
    ctx.call(
        ScHname::new(atomic_swap.erc20_sc_name_sender.as_str()),
        ScHname::new(ERC20_FUNC_TRANSFER),
        Some(transfer_from_params),
        None
    );

    // set the atomic swap to completed
    atomic_swap.finished = 1;

    // the auction object should be deleted here, unfortunately, the current implementation of the vm does not support it so we simply set the atomic swap to finished
    swap.set_value(&atomic_swap.to_bytes());
}

pub fn func_finalize_swap(ctx: &ScFuncContext) {
    let p = ctx.params();

    // get the id where the swap should be accessible
    let param_swap_id = p.get_string(PARAM_SWAP_ID);

    ctx.require(param_swap_id.exists(), "missing mandatory swap id");

    // get the state
    let state: ScMutableMap = ctx.state();
    // get the atomic swap map
    let atomic_swaps = state.get_map(VAR_ATOMIC_SWAPS);
    // get the swap with swap_id
    let swap = atomic_swaps.get_bytes(&param_swap_id.value());

    // check if swap id exists
    ctx.require(swap.exists(), "swap id does not exist");

    // parse atomic swap
    let mut atomic_swap = AtomicSwap::from_bytes(&swap.value());

    // check if atomic swap is already finished
    ctx.require(atomic_swap.finished == 0, "swap is already finished");

    // check if this contract is allowed to transfer the amount of tokens from receiver's erc20 contract
    let allowance_params = ScMutableMap::new();
    allowance_params.get_agent_id(ERC20_PARAM_ACCOUNT).set_value(&atomic_swap.agent_id_recipient);
    allowance_params.get_agent_id(ERC20_PARAM_DELEGATION).set_value(&ctx.contract_id().as_agent_id());

    // get the allowance value from the erc20 contract by calling the view
    let allowance_result_map = ctx.call(
        ScHname::new(atomic_swap.erc20_sc_name_recipient.as_str()),
        ScHname::new(ERC20_VIEW_ALLOWANCE),
        Some(allowance_params),
        None
    );

    // check if allowance amount is enough
    ctx.require(allowance_result_map.get_int64(ERC20_PARAM_AMOUNT).value() >= atomic_swap.amount_recipient, "contract is not allowed to transfer the specified amount of tokens");

    // check if atomic swap is still open
    ctx.require(ctx.timestamp() / 1000000000 <= atomic_swap.when_started + atomic_swap.duration_open, "swap is not open anymore");

    // transfer erc20 tokens from recipient's erc20 contract to this contract
    let transfer_from_params = ScMutableMap::new();
    transfer_from_params.get_agent_id(ERC20_PARAM_ACCOUNT).set_value(&atomic_swap.agent_id_recipient);
    transfer_from_params.get_agent_id(ERC20_PARAM_RECIPIENT).set_value(&ctx.contract_id().as_agent_id());
    transfer_from_params.get_int64(ERC20_PARAM_AMOUNT).set_value(atomic_swap.amount_recipient);

    // call transfer contract
    ctx.call(
        ScHname::new(atomic_swap.erc20_sc_name_recipient.as_str()),
        ScHname::new(ERC20_FUNC_TRANSFER_FROM),
        Some(transfer_from_params),
        None
    );

    // transfer erc20 tokens from this contract to sender's agent id
    let transfer_from_params_sender = ScMutableMap::new();
    transfer_from_params_sender.get_agent_id(ERC20_PARAM_ACCOUNT).set_value(&atomic_swap.agent_id_sender);
    transfer_from_params_sender.get_int64(ERC20_PARAM_AMOUNT).set_value(atomic_swap.amount_recipient);

    // call transfer contract
    ctx.call(
        ScHname::new(atomic_swap.erc20_sc_name_recipient.as_str()),
        ScHname::new(ERC20_FUNC_TRANSFER),
        Some(transfer_from_params_sender),
        None
    );

    // transfer erc20 tokens from this contract to recipient's agent id
    let transfer_from_params_recipient = ScMutableMap::new();
    transfer_from_params_recipient.get_agent_id(ERC20_PARAM_ACCOUNT).set_value(&atomic_swap.agent_id_recipient);
    transfer_from_params_recipient.get_int64(ERC20_PARAM_AMOUNT).set_value(atomic_swap.amount_sender);

    // call transfer contract
    ctx.call(
        ScHname::new(atomic_swap.erc20_sc_name_sender.as_str()),
        ScHname::new(ERC20_FUNC_TRANSFER),
        Some(transfer_from_params_recipient),
        None
    );

    // set the atomic swap to completed
    atomic_swap.finished = 1;

    // the auction object should be deleted here, unfortunately, the current implementation of the vm does not support it so we simply set the atomic swap to finished
    swap.set_value(&atomic_swap.to_bytes());
}

pub fn view_get_swap_by_id(ctx: &ScViewContext) {
    let p = ctx.params();

    let param_swap_id = p.get_string(PARAM_SWAP_ID);

    ctx.require(param_swap_id.exists(), "missing mandatory swap id");

    // get the state
    let state = ctx.state();
    // get the atomic swap map
    let atomic_swaps = state.get_map(VAR_ATOMIC_SWAPS);
    // get the swap with swap_id
    let swap = atomic_swaps.get_bytes(&param_swap_id.value());

    // check if swap id exists
    ctx.require(swap.exists(), "swap id does not exist");

    ctx.results().get_bytes(PARAM_SWAP).set_value(&swap.value());
}
