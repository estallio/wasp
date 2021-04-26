// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

use consts::*;
use atomicswapcoloredtokens::*;
use wasmlib::*;

mod consts;
mod atomicswapcoloredtokens;
mod types;

#[no_mangle]
fn on_load() {
    let exports = ScExports::new();
    exports.add_func(FUNC_START_SWAP, func_start_swap);
    exports.add_func(FUNC_CANCEL_SWAP, func_cancel_swap);
    exports.add_func(FUNC_FINALIZE_SWAP, func_finalize_swap);
    exports.add_view(VIEW_GET_OPEN_SWAPS_BY_AGENT_ID, view_get_open_swaps_by_agent_id);
    exports.add_view(VIEW_GET_SWAP_BY_ID, view_get_swap_by_id);
}
