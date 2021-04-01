package chain

import (
	"fmt"

	"github.com/iotaledger/wasp/packages/chain"
	"github.com/iotaledger/wasp/tools/wasp-cli/config"
	"github.com/iotaledger/wasp/tools/wasp-cli/log"
)

func listCmd(args []string) {
	client := config.WaspClient()
	chains, err := client.GetChainRecordList()
	log.Check(err)
	log.Printf("Total %d chain(s) in wasp node %s\n", len(chains), client.BaseURL())
	showChainList(chains)
}

func showChainList(chains []*chain.ChainRecord) {
	header := []string{"chainid", "active"}
	rows := make([][]string, len(chains))
	for i, chain := range chains {
		rows[i] = []string{
			chain.ChainID.String(),
			fmt.Sprintf("%v", chain.Active),
		}
	}
	log.PrintTable(header, rows)
}
