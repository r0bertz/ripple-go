package wss

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

func Send(c *websocket.Conn, r interface{}) error {
	message, err := json.Marshal(r)
	if err != nil {
		return err
	}
	return c.WriteMessage(websocket.TextMessage, message)
}

func Receive(c *websocket.Conn) (interface{}, error) {
	var i interface{}
	if err := c.ReadJSON(&i); err != nil {
		return i, err
	}
	return i, nil
}

func Connect(addr string) (*websocket.Conn, *http.Response, error) {
	u := url.URL{Scheme: "wss", Host: addr}
	log.Printf("connecting to %s", u.String())
	return websocket.DefaultDialer.Dial(u.String(), nil)
}
