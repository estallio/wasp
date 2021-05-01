// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

use wasmlib::*;

//@formatter:off
pub struct AtomicSwap {
    // hash of the key_secret that is used to "lock" the tokens - this is passed at swap start
    pub key_hash:               ScHash,
    // secret to redeem the tokens (secret of key_hash) - this is shared at a later point to redeem the locked tokens
    pub key_secret:             String,

    // name of the erc20 contract where the sender stores his/her tokens
    pub erc20_sc_name_sender:   String,

    pub amount_sender:          i64,

    pub agent_id_sender:        ScAgentId,
    pub agent_id_recipient:     ScAgentId,

    pub duration_open:          i64,
    pub when_started:           i64,

    // TODO: could also be used as state to indicate started/cancelled/successful
    pub finished:               i64,
}
//@formatter:on

impl AtomicSwap {
    pub fn from_bytes(bytes: &[u8]) -> AtomicSwap {
        let mut decode = BytesDecoder::new(bytes);
        AtomicSwap {
            key_hash: decode.hash(),
            key_secret: decode.string(),
            erc20_sc_name_sender: decode.string(),
            amount_sender: decode.int64(),
            agent_id_sender: decode.agent_id(),
            agent_id_recipient: decode.agent_id(),
            duration_open: decode.int64(),
            when_started: decode.int64(),
            finished: decode.int64()
        }
    }

    pub fn to_bytes(&self) -> Vec<u8> {
        let mut encode = BytesEncoder::new();
        encode.hash(&self.key_hash);
        encode.string(&self.key_secret);
        encode.string(&self.erc20_sc_name_sender);
        encode.int64(self.amount_sender);
        encode.agent_id(&self.agent_id_sender);
        encode.agent_id(&self.agent_id_recipient);
        encode.int64(self.duration_open);
        encode.int64(self.when_started);
        encode.int64(self.finished);
        return encode.data();
    }
}
