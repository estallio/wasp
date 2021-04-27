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

    // next 3 params define from whom, how much and what type of color the sender wants to get for his/her deposit
    let param_color_recipient = p.get_color(PARAM_COLOR_RECIPIENT);
    let param_amount_recipient = p.get_int64(PARAM_AMOUNT_RECIPIENT);
    let param_address_recipient = p.get_address(PARAM_ADDRESS_RECIPIENT);

    // next 2 parameters are not absolutely necessary, as they could be extracted from balances
    // it was implemented the following way to keep the effort small
    // it would also be possible to exchange multiple different colors this way and adjust the AtomicSwap type to include a "exchange color map"
    let param_color_sender = p.get_color(PARAM_COLOR_SENDER);
    let param_amount_sender = p.get_int64(PARAM_AMOUNT_SENDER);

    // get the sender - this is the one calling this contract/creating the swap
    let param_address_sender = ctx.caller().address();

    // get the parameter how long the atomic swap should be valid - this is only an additional feature
    // and is not needed for this use-case as we can simply cancel the transfer at any time
    let param_duration = p.get_int64(PARAM_DURATION);

    ctx.log("checking params now...");

    // check if all necessary variables are set - only for unfortunate reasons when posting the tx
    ctx.require(param_color_sender.exists(), "missing mandatory sender color");
    ctx.require(param_amount_sender.exists(), "missing mandatory sender amount");
    ctx.require(param_color_recipient.exists(), "missing mandatory recipient color");
    ctx.require(param_amount_recipient.exists(), "missing mandatory recipient amount");
    ctx.require(param_address_recipient.exists(), "missing mandatory recipient id");
    ctx.require(param_swap_id.exists(), "missing swap id");

    ctx.log("checking amount now...");

    // check that the amount of color is sufficient
    let amount = ctx.incoming().balance(&param_color_sender.value());
    ctx.require(amount == param_amount_sender.value(), "transferred balance of color does not match amount parameter");

    // check that no other color is transferred to the smart contract to prevent loosing other tokens
    let color_length = ctx.incoming().colors().length();
    ctx.require(color_length <= 1, "only one color is allowed to be transferred to the contract");

    // we have a little problem in the next step
    //  fist, we have to define how we save an atomic swap object in the state
    //  as we have to access it later, we should have an unique id for the swap
    //  we should get sure it is possible that multiple agents can place multiple atomic swaps simultaneously
    //  1. one possible solution is to save it with the recipient and sender id and the color as hash
    //  and simply restrict the amount of atomic swaps 2 parties are able to start
    //  2. use an array for holding all swaps and simply append the swap objects - there is a upper limit in this case
    //  but we could simply use an swap-array for each 2 agent ids
    //  3. let the user specify a swap id where the atomic swap object is saved
    //  --- for now, we simply use the third option and require a swap id param and don't do the clean up - also it may not be possible anyways

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
        color_sender: param_color_sender.value(),
        color_recipient: param_color_recipient.value(),
        amount_sender: param_amount_sender.value(),
        amount_recipient: param_amount_recipient.value(),
        address_sender: param_address_sender,
        address_recipient: param_address_recipient.value(),
        duration_open: param_duration.value(),
        when_started: ctx.timestamp() / 1000000000,
        // for now, a number is sufficient to model the state, maybe there is any boolean or enum support sometime to model such states
        finished: 0,
    };

    ctx.log("saving atomic swap now...");

    swap.set_value(&atomic_swap.to_bytes());
}

// TODO: maybe also implement an open_swaps view - this is a little bit harder, as dynamic lists are not existing until now
//  2 possibilities:
//    1. save all swaps and keep them in state (expensive)
//    2. run every time over all swaps and find open open swaps by agent_id or save open swap_ids in parallel in an array in a agent_id map (both methods require computational power)

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
    ctx.require(ctx.caller().address() == atomic_swap.address_sender, "only the sender is able to cancel the swap");

    // transfer money back to sender
    transfer(ctx, &atomic_swap.address_sender, &atomic_swap.color_sender, atomic_swap.amount_sender);

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

    // check if the caller of this method is the recipient
    ctx.require(ctx.caller().address() == atomic_swap.address_recipient, "only the recipient is able to finalize the swap");

    // check if atomic swap is still open
    ctx.require(ctx.timestamp() / 1000000000 <= atomic_swap.when_started + atomic_swap.duration_open, "swap is not open anymore");

    // get the balances the recipient has sent to the contract
    let amount = ctx.incoming().balance(&atomic_swap.color_recipient);

    // check if recipient sent enough coins
    ctx.require(amount == atomic_swap.amount_recipient, "swap is not open anymore");

    // check that no other color is transferred to the smart contract to prevent loosing other tokens
    let color_length = ctx.incoming().colors().length();
    ctx.require(color_length <= 1, "only one color is allowed to be transferred to the contract");

    // transfer money to the parties
    transfer(ctx, &atomic_swap.address_recipient, &atomic_swap.color_sender, atomic_swap.amount_sender);
    transfer(ctx, &atomic_swap.address_sender, &atomic_swap.color_recipient, atomic_swap.amount_recipient);

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

// helper method copied from fairauction example
fn transfer(ctx: &ScFuncContext, address: &ScAddress, color: &ScColor, amount: i64) {
    // send back to original Tangle address
    ctx.transfer_to_address(&address, ScTransfers::new(&color, amount));
}
