package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/r0bertz/ripple-go/tx"
	"github.com/r0bertz/ripple-go/wss"
)

var (
	addr = flag.String("addr", "s2.ripple.com:443", "wss service address")
	txID = flag.String("tx", "", "ripple transaction")
)

func main() {
	flag.Parse()

	if *txID == "" {
		log.Fatal("--tx not set")
	}

	c, _, err := wss.Connect(*addr)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	if err := wss.Send(c, tx.NewRequest(*txID)); err != nil {
		log.Fatal(err)
	}
	i, err := wss.Receive(c)
	if err != nil {
		log.Fatal(err)
	}
	m := i.(map[string]interface{})
	b, _ := json.MarshalIndent(m, "", "  ")
	fmt.Println(string(b))
}
