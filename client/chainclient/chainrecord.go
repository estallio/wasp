package chainclient

import (
	"github.com/iotaledger/wasp/packages/chain"
)

// GetChainRecord fetches the chain's Record
func (c *Client) GetChainRecord() (*chain.ChainRecord, error) {
	return c.WaspClient.GetChainRecord(c.ChainID)
}
