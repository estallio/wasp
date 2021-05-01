// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

#![allow(dead_code)]

use wasmlib::*;

pub const SC_NAME: &str = "atomicswaperc20sync";

pub const PARAM_SWAP_ID: &str = "swapId";
pub const PARAM_SWAP: &str = "swap";
pub const PARAM_AMOUNT_SENDER: &str = "amountSender";
pub const PARAM_AMOUNT_RECIPIENT: &str = "amountRecipient";
pub const PARAM_AGENT_ID_RECIPIENT: &str = "agentIdRecipient";
pub const PARAM_DURATION: &str = "duration";
pub const PARAM_SC_NAME_SENDER: &str = "scNameSender";
pub const PARAM_SC_NAME_RECIPIENT: &str = "scNameRecipient";
pub const PARAM_KEY_HASH: &str = "keyHash";
pub const PARAM_KEY_SECRET: &str = "keySecret";

pub const VAR_ATOMIC_SWAPS: &str = "atomicSwaps";
pub const VAR_AMOUNT_SENDER: &str = "amountSender";
pub const VAR_AMOUNT_RECIPIENT: &str = "amountRecipient";
pub const VAR_ADDRESS_SENDER: &str = "addressSender";
pub const VAR_ADDRESS_RECIPIENT: &str = "addressRecipient";
pub const VAR_DURATION_OPEN: &str = "durationOpen";
pub const VAR_WHEN_STARTED: &str = "whenStarted";
pub const VAR_FINISHED: &str = "finished";

pub const FUNC_START_SWAP: &str = "startSwap";
pub const FUNC_CANCEL_SWAP: &str = "cancelSwap";
pub const FUNC_FINALIZE_SWAP: &str = "finalizeSwap";

pub const VIEW_GET_SWAP_BY_ID: &str = "getSwapById";
pub const VIEW_GET_SECRET_BY_SWAP_ID: &str = "getSecretBySwapId";

pub const ERC20_FUNC_APPROVE: &str = "approve";
pub const ERC20_FUNC_INIT: &str = "init";
pub const ERC20_FUNC_TRANSFER: &str = "transfer";
pub const ERC20_FUNC_TRANSFER_FROM: &str = "transferFrom";
pub const ERC20_VIEW_ALLOWANCE: &str = "allowance";
pub const ERC20_VIEW_BALANCE_OF: &str = "balanceOf";
pub const ERC20_VIEW_TOTAL_SUPPLY: &str = "totalSupply";

pub const ERC20_PARAM_ACCOUNT: &str = "ac";
pub const ERC20_PARAM_AMOUNT: &str = "am";
pub const ERC20_PARAM_CREATOR: &str = "c";
pub const ERC20_PARAM_DELEGATION: &str = "d";
pub const ERC20_PARAM_RECIPIENT: &str = "r";
pub const ERC20_PARAM_SUPPLY: &str = "s";
