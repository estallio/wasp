package multiclient

import (
	"github.com/iotaledger/wasp/client"
	"github.com/iotaledger/wasp/packages/chain"
)

// PutChainRecord calls PutChainRecord in all wasp nodes
func (m *MultiClient) PutCommitteeRecord(bd *chain.CommitteeRecord) error {
	return m.Do(func(i int, w *client.WaspClient) error {
		return w.PutCommitteeRecord(bd)
	})
}
