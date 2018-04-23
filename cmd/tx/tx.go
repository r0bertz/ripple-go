package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/r0bertz/ripple/data"
	"github.com/r0bertz/ripple/websockets"
)

var (
	addr = flag.String("addr", "wss://s2.ripple.com:443", "wss service address")
	txID = flag.String("tx", "", "ripple transaction")
)

func main() {
	flag.Parse()

	if *txID == "" {
		log.Fatal("--tx not set")
	}

	remote, err := websockets.NewRemote(*addr)
	if err != nil {
		log.Fatal(err)
	}

	th, err := data.NewHash256(*txID)
	if err != nil {
		log.Fatal(err)
	}
	tx, err := remote.Tx(*th)
	if err != nil {
		log.Fatal(err)
	}
	b, _ := json.MarshalIndent(tx, "", "  ")
	fmt.Println(string(b))
}
