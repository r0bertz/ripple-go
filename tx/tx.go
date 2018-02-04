package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/r0bertz/ripple/wss"
)

var (
	addr = flag.String("addr", "s2.ripple.com:443", "wss service address")
	tx   = flag.String("tx", "", "ripple transaction")
)

type TxRequest struct {
	Command     string `json:"command"`
	Transaction string `json:"transaction"`
	Binary      bool   `json:"binary"`
}

func NewTxRequest(transaction string) *TxRequest {
	return &TxRequest{
		Command:     "tx",
		Transaction: transaction,
		Binary:      false,
	}
}

func main() {
	flag.Parse()

	if *tx == "" {
		log.Fatal("--tx not set")
	}

	c, _, err := wss.Connect(*addr)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	if err := wss.Send(c, NewTxRequest(*tx)); err != nil {
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
