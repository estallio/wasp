// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package registry

import (
	"github.com/iotaledger/hive.go/logger"
	"github.com/iotaledger/wasp/packages/dbprovider"
	"github.com/iotaledger/wasp/packages/tcrypto"
)

// Impl is just a placeholder to implement all interfaces needed by different components.
// Each of the interfaces are implemented in the corresponding file in this package.
// Implements registry.RegistryProvider. Implementation is mostly in chainrecord.go. TODO: tidy it up
type Impl struct {
	suite      tcrypto.Suite
	log        *logger.Logger
	dbProvider *dbprovider.DBProvider
}

// NewRegistry creates new instance of the registry implementation.
func NewRegistry(suite tcrypto.Suite, log *logger.Logger, dbp *dbprovider.DBProvider) *Impl {
	ret := &Impl{
		suite:      suite,
		log:        log.Named("registry"),
		dbProvider: dbp,
	}
	return ret
}

//TODO: remove this method somehow after the merge
func (rImplThis *Impl) GetDBProvider() *dbprovider.DBProvider {
	return rImplThis.dbProvider
}
