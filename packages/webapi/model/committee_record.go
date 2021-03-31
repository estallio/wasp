package model

import (
	"github.com/iotaledger/wasp/packages/chain"
)

type CommitteeRecord struct {
	Address        Address  `swagger:"desc(Committee address (base58-encoded))"`
	CommitteeNodes []string `swagger:"desc(List of committee nodes (network IDs))"`
}

func NewCommitteeRecord(bd *chain.CommitteeRecord) *CommitteeRecord {
	return &CommitteeRecord{
		Address:        NewAddress(bd.Address),
		CommitteeNodes: bd.CommitteeNodes,
	}
}

func (bd *CommitteeRecord) Record() *chain.CommitteeRecord {
	return &chain.CommitteeRecord{
		Address:        bd.Address.Address(),
		CommitteeNodes: bd.CommitteeNodes,
	}
}
