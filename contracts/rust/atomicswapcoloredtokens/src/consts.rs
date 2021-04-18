// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

#![allow(dead_code)]

use wasmlib::*;

pub const SC_NAME: &str = "atomicswapcoloredtokens";

pub const PARAM_SWAP_ID: &str = "swapId";
pub const PARAM_COLOR_SENDER: &str = "colorSender";
pub const PARAM_COLOR_RECEIVER: &str = "colorReceiver";
pub const PARAM_AMOUNT_SENDER: &str = "amountSender";
pub const PARAM_AMOUNT_RECEIVER: &str = "amountReceiver";
pub const PARAM_ADDRESS_RECEIVER: &str = "addressReceiver";
pub const PARAM_DURATION: &str = "duration";

pub const VAR_ATOMIC_SWAPS: &str = "atomicSwaps";
pub const VAR_SWAP_ID: &str = "swapId";
pub const VAR_COLOR_SENDER: &str = "colorSender";
pub const VAR_COLOR_RECEIVER: &str = "colorReceiver";
pub const VAR_AMOUNT_SENDER: &str = "amountSender";
pub const VAR_AMOUNT_RECEIVER: &str = "amountReceiver";
pub const VAR_ADDRESS_SENDER: &str = "addressSender";
pub const VAR_ADDRESS_RECEIVER: &str = "addressReceiver";
pub const VAR_DURATION_OPEN: &str = "durationOpen";
pub const VAR_WHEN_STARTED: &str = "whenStarted";
pub const VAR_FINISHED: &str = "whenStarted";

pub const FUNC_START_SWAP: &str = "startSwap";
pub const FUNC_CANCEL_SWAP: &str = "cancelSwap";
pub const FUNC_FINALIZE_SWAP: &str = "finalizeSwap";
pub const VIEW_GET_OPEN_SWAPS: &str = "getOpenSwaps";
