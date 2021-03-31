// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package consensus_test

// TODO: Tests with corrupted messages.
// TODO: Tests with byzantine messages.
// TODO: Single node down for some time.

import (
	"fmt"
	"testing"
	"time"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	"github.com/iotaledger/hive.go/logger"
	"github.com/iotaledger/wasp/packages/chain"
	"github.com/iotaledger/wasp/packages/chain/chainimpl"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/dbprovider"
	"github.com/iotaledger/wasp/packages/dkg"
	"github.com/iotaledger/wasp/packages/peering"
	"github.com/iotaledger/wasp/packages/registry"
	"github.com/iotaledger/wasp/packages/testutil"
	"github.com/stretchr/testify/require"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing"
	"go.dedis.ch/kyber/v3/util/key"
)

// TestBasic checks if DKG procedure is executed successfully in a common case.
func TestBasic(t *testing.T) {
	log := testutil.NewLogger(t)
	defer log.Sync()
	//
	// Create a fake network and keys for the tests.
	var timeout = 100 * time.Second
	var threshold uint16 = 10
	var peerCount uint16 = 10
	var peerNetIDs []string = make([]string, peerCount)
	var peerPubs []kyber.Point = make([]kyber.Point, len(peerNetIDs))
	var peerSecs []kyber.Scalar = make([]kyber.Scalar, len(peerNetIDs))
	var suite = pairing.NewSuiteBn256() // NOTE: That's from the Pairing Adapter.
	for i := range peerNetIDs {
		peerPair := key.NewKeyPair(suite)
		peerNetIDs[i] = fmt.Sprintf("P%02d", i)
		peerSecs[i] = peerPair.Private
		peerPubs[i] = peerPair.Public
	}
	var peeringNetwork *testutil.PeeringNetwork = testutil.NewPeeringNetwork(
		peerNetIDs, peerPubs, peerSecs, 10000,
		testutil.NewPeeringNetReliable(),
		testutil.WithLevel(log, logger.LevelWarn, false),
	)
	var networkProviders []peering.NetworkProvider = peeringNetwork.NetworkProviders()
	//
	// Initialize the DKG subsystem in each node.
	var dkgNodes []*dkg.Node = make([]*dkg.Node, len(peerNetIDs))
	var dkgRegistryProviders []*testutil.DkgRegistryProvider = make([]*testutil.DkgRegistryProvider, len(peerNetIDs))
	for i := range peerNetIDs {
		dkgRegistryProviders[i] = testutil.NewDkgRegistryProvider(suite)
		dkgNodes[i] = dkg.NewNode(
			peerSecs[i], peerPubs[i], suite, networkProviders[i], dkgRegistryProviders[i],
			testutil.WithLevel(log.With("NetID", peerNetIDs[i]), logger.LevelDebug, false),
		)
	}
	//
	// Initiate the key generation from some client node.
	dkShare, err := dkgNodes[0].GenerateDistributedKey(
		peerNetIDs,
		peerPubs,
		threshold,
		1*time.Second,
		2*time.Second,
		timeout,
	)
	require.Nil(t, err)
	//
	// Initiate the chain
	chainimpl.Init()
	var chains []chain.Chain = make([]chain.Chain, len(peerNetIDs))
	chainRecord := registry.ChainRecord{
		ChainID:        coretypes.ChainID(*dkShare.Address),
		Color:          balance.ColorIOTA,
		CommitteeNodes: peerNetIDs,
		Active:         true,
	}
	for i := range peerNetIDs {
		log.Debugf("XXX chain %v IN", i)
		db := dbprovider.NewInMemoryDBProvider(log)
		registry := registry.NewRegistry(nil, log, db)
		log.Debugf("XXX chain %v CREATE", i)
		//chains[i] = chain.New(&chainRecord, log, networkProviders[i], dkgRegistryProviders[i], registry, func() {
		//nodeconn.Subscribe((address.Address)(chr.ChainID), chr.Color)
		//})
		log.Debugf("XXX chain %v OUT %v %v %v", i, chainRecord, registry, chains[i])
		//require.NotNil(t, chains[i])
	}

}
