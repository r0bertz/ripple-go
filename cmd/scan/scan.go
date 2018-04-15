package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/r0bertz/ripple/data"
	"github.com/r0bertz/ripple/websockets"
)

var (
	addr        = flag.String("addr", "wss://s2.ripple.com:443", "wss service address")
	ledgerIndex = flag.Uint("ledger_index", 0, "starting ledger index. scan backwards.")
)

func main() {
	flag.Parse()

	if *ledgerIndex == 0 {
		log.Fatal("--ledger_index not set")
	}
	remote, err := websockets.NewRemote(*addr)
	if err != nil {
		log.Fatal(err)
	}
	for ; *ledgerIndex > 0; *ledgerIndex-- {
		fmt.Printf("ledger: %d\n", *ledgerIndex)
		result, err := remote.Ledger(*ledgerIndex, true)
		if err != nil {
			log.Fatal(err)
		}
		for _, tmx := range result.Ledger.Transactions {
			account := tmx.GetBase().Account
			for _, n := range tmx.MetaData.AffectedNodes {
				node, final, _, s := n.AffectedNode()
				if node.LedgerEntryType == data.RIPPLE_STATE && s == data.Created {
					e := final.(*data.RippleState)
					if !e.HighLimit.Issuer.Equals(account) {
						fmt.Println(tmx.GetBase().Hash.String())
						return
					}
				}
			}
		}
	}
}
