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
    exports.add_view(VIEW_GET_OPEN_SWAPS, view_get_open_swaps);
}
