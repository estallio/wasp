package model

import (
	"github.com/iotaledger/wasp/packages/chain"
)

type CommitteeRecord struct {
	Address Address  `swagger:"desc(Committee address (base58-encoded))"`
	Nodes   []string `swagger:"desc(List of committee nodes (network IDs))"`
}

func NewCommitteeRecord(bd *chain.CommitteeRecord) *CommitteeRecord {
	return &CommitteeRecord{
		Address: NewAddress(bd.Address),
		Nodes:   bd.Nodes,
	}
}

func (bd *CommitteeRecord) Record() *chain.CommitteeRecord {
	return &chain.CommitteeRecord{
		Address: bd.Address.Address(),
		Nodes:   bd.Nodes,
	}
}
