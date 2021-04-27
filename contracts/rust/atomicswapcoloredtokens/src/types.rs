// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

use wasmlib::*;

//@formatter:off
pub struct AtomicSwap {
    pub color_sender:       ScColor,
    pub color_recipient:    ScColor,
    pub amount_sender:      i64,
    pub amount_recipient:   i64,
    pub address_sender:     ScAddress,
    pub address_recipient:  ScAddress,
    pub duration_open:      i64,
    pub when_started:       i64,
    pub finished:           i64,
}
//@formatter:on

impl AtomicSwap {
    pub fn from_bytes(bytes: &[u8]) -> AtomicSwap {
        let mut decode = BytesDecoder::new(bytes);
        AtomicSwap {
            color_sender: decode.color(),
            color_recipient: decode.color(),
            amount_sender: decode.int64(),
            amount_recipient: decode.int64(),
            address_sender: decode.address(),
            address_recipient: decode.address(),
            duration_open: decode.int64(),
            when_started: decode.int64(),
            finished: decode.int64()
        }
    }

    pub fn to_bytes(&self) -> Vec<u8> {
        let mut encode = BytesEncoder::new();
        encode.color(&self.color_sender);
        encode.color(&self.color_recipient);
        encode.int64(self.amount_sender);
        encode.int64(self.amount_recipient);
        encode.address(&self.address_sender);
        encode.address(&self.address_recipient);
        encode.int64(self.duration_open);
        encode.int64(self.when_started);
        encode.int64(self.finished);
        return encode.data();
    }
}
