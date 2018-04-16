package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/r0bertz/ripple/data"
	"github.com/r0bertz/ripple/websockets"
)

var (
	addr    = flag.String("addr", "wss://s2.ripple.com:443", "wss service address")
	account = flag.String("account", "", "ripple account")
)

func main() {
	flag.Parse()

	if *account == "" {
		log.Fatal("--account not set")
	}
	acct, err := data.NewAccountFromAddress(*account)
	if err != nil {
		log.Fatal(err)
	}
	remote, err := websockets.NewRemote(*addr)
	if err != nil {
		log.Fatal(err)
	}
	for t := range remote.AccountTx(*acct, 10, -1, -1) {
		fmt.Println(t.GetBase().Hash)
	}
}
