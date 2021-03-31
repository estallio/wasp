// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package chain

import (
	"github.com/iotaledger/goshimmer/packages/ledgerstate"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/dbprovider"
)

// RegistryProvider stands for a partial registry interface, needed for chain package.
// It should be implemented by wasp/packages/registry.Impl
type RegistryProvider interface {
	SaveChainRecord(bd *ChainRecord) error
	GetChainRecord(chainID *coretypes.ChainID) (*ChainRecord, error)
	UpdateChainRecord(chainID *coretypes.ChainID, f func(*ChainRecord) bool) (*ChainRecord, error)
	ActivateChainRecord(chainID *coretypes.ChainID) (*ChainRecord, error)
	DeactivateChainRecord(chainID *coretypes.ChainID) (*ChainRecord, error)
	GetChainRecords() ([]*ChainRecord, error)

	GetCommitteeRecord(addr ledgerstate.Address) (*CommitteeRecord, error)
	SaveCommitteeRecord(rec *CommitteeRecord) error

	GetDBProvider() *dbprovider.DBProvider //TODO: remove this method somehow after the merge
}
