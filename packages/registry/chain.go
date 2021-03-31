//Implementations of wasp/packages/chain.RegistryProvider interface methods
//for wasp/packages/registry.Impl object
package registry

import (
	"fmt"
	"github.com/iotaledger/goshimmer/packages/ledgerstate"
	"github.com/iotaledger/hive.go/kvstore"
	"github.com/iotaledger/wasp/packages/chain"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/dbprovider"
	"github.com/mr-tron/base58"
)

// GetChainRecord reads ChainRecord from registry.
// Returns nil if not found
// Implements wasp/packages/chain.RegistryProvider.GetChainRecord(*coretypes.ChainID)
func (rImplThis *Impl) GetChainRecord(chainID *coretypes.ChainID) (*chain.ChainRecord, error) {
	data, err := rImplThis.dbProvider.GetRegistryPartition().Get(dbKeyChainRecord(chainID))
	if err == kvstore.ErrKeyNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return chain.ChainRecordFromBytes(data)
}

func dbKeyChainRecord(chainID *coretypes.ChainID) []byte {
	return dbprovider.MakeKey(dbprovider.ObjectTypeChainRecord, chainID.Bytes())
}

// Implements wasp/packages/chain.RegistryProvider.SaveChainRecord(*chain.ChainRecord)
func (rImplThis *Impl) SaveChainRecord(rec *chain.ChainRecord) error {
	return rImplThis.dbProvider.GetRegistryPartition().Set(dbKeyChainRecord(&rec.ChainID), rec.Bytes())
}

// Implements wasp/packages/chain.RegistryProvider.UpdateChainRecord(*coretypes.ChainID, func(*chain.ChainRecord) bool)
func (rImplThis *Impl) UpdateChainRecord(chainID *coretypes.ChainID, f func(*chain.ChainRecord) bool) (*chain.ChainRecord, error) {
	rec, err := rImplThis.GetChainRecord(chainID)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, fmt.Errorf("no chain record found for chainID %s", chainID.String())
	}
	if f(rec) {
		err = rImplThis.SaveChainRecord(rec)
		if err != nil {
			return nil, err
		}
	}
	return rec, nil
}

// Implements wasp/packages/chain.RegistryProvider.ActivateChainRecord(*coretypes.ChainID)
func (rImplThis *Impl) ActivateChainRecord(chainID *coretypes.ChainID) (*chain.ChainRecord, error) {
	return rImplThis.UpdateChainRecord(chainID, func(bd *chain.ChainRecord) bool {
		if bd.Active {
			return false
		}
		bd.Active = true
		return true
	})
}

// Implements wasp/packages/chain.RegistryProvider.DeactivateChainRecord(*coretypes.ChainID)
func (rImplThis *Impl) DeactivateChainRecord(chainID *coretypes.ChainID) (*chain.ChainRecord, error) {
	return rImplThis.UpdateChainRecord(chainID, func(bd *chain.ChainRecord) bool {
		if !bd.Active {
			return false
		}
		bd.Active = false
		return true
	})
}

// Implements wasp/packages/chain.RegistryProvider.GetChainRecords()
func (rImplThis *Impl) GetChainRecords() ([]*chain.ChainRecord, error) {
	db := rImplThis.dbProvider.GetRegistryPartition()
	ret := make([]*chain.ChainRecord, 0)

	err := db.Iterate([]byte{dbprovider.ObjectTypeChainRecord}, func(key kvstore.Key, value kvstore.Value) bool {
		if rec, err1 := chain.ChainRecordFromBytes(value); err1 == nil {
			ret = append(ret, rec)
		} else {
			log.Warnf("corrupted chain record with key %s", base58.Encode(key))
		}
		return true
	})
	return ret, err
}

// CommitteeRecordFromRegistry reads CommitteeRecord from registry.
// Returns nil if not found
// Implements wasp/packages/chain.RegistryProvider.GetCommitteeRecord(ledgerstate.Address)
func (rImplThis *Impl) GetCommitteeRecord(addr ledgerstate.Address) (*chain.CommitteeRecord, error) {
	data, err := rImplThis.dbProvider.GetRegistryPartition().Get(dbKeyCommitteeRecord(addr))
	if err == kvstore.ErrKeyNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return chain.CommitteeRecordFromBytes(data)
}

func dbKeyCommitteeRecord(addr ledgerstate.Address) []byte {
	return dbprovider.MakeKey(dbprovider.ObjectTypeCommitteeRecord, addr.Bytes())
}

// Implements wasp/packages/chain.RegistryProvider.SaveCommitteeRecord(rec *chain.CommitteeRecord)
func (rImplThis *Impl) SaveCommitteeRecord(rec *chain.CommitteeRecord) error {
	return rImplThis.dbProvider.GetRegistryPartition().Set(dbKeyCommitteeRecord(rec.Address), rec.Bytes())
}

//TODO: remove this method somehow after the merge
// Implements wasp/packages/chain.RegistryProvider.GetDBProvider()
func (rImplThis *Impl) GetDBProvider() *dbprovider.DBProvider {
	return rImplThis.dbProvider
}
