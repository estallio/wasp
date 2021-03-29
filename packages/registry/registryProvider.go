// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package registry //TODO: move the module to packages/chain

import (
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/dbprovider"
)

// RegistryProvider stands for a partial registry interface, needed for chain package.
// It should be implemented by registry.Impl
type RegistryProvider interface {
	SaveChainRecord(bd *ChainRecord) error
	GetChainRecord(chainID *coretypes.ChainID) (*ChainRecord, error)
	UpdateChainRecord(chainID *coretypes.ChainID, f func(*ChainRecord) bool) (*ChainRecord, error)
	ActivateChainRecord(chainID *coretypes.ChainID) (*ChainRecord, error)
	DeactivateChainRecord(chainID *coretypes.ChainID) (*ChainRecord, error)
	GetChainRecords() ([]*ChainRecord, error)

	GetDBProvider() *dbprovider.DBProvider //TODO: remove this method somehow after the merge
}
