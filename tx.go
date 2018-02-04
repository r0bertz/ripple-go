package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
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

func send(c *websocket.Conn, r interface{}) error {
	message, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
	}
	return c.WriteMessage(websocket.TextMessage, message)
}

func receive(c *websocket.Conn) {
	var i interface{}
	if err := c.ReadJSON(&i); err != nil {
		log.Fatal("ReadJSON failed: ", err)
	}
	m := i.(map[string]interface{})
	b, _ := json.MarshalIndent(m, "", "  ")
	fmt.Println(string(b))
}

func main() {
	flag.Parse()

	if *tx == "" {
		log.Fatal("--tx not set")
	}

	u := url.URL{Scheme: "wss", Host: *addr}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	if err := send(c, NewTxRequest(*tx)); err != nil {
		log.Fatal(err)
	}
	receive(c)
}
