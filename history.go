package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

var (
	addr       = flag.String("addr", "s2.ripple.com:443", "WSS service address")
	acct       = flag.String("account", "", "Ripple address")
	checkpoint = flag.String("checkpoint", "checkpoint", "File contains the last examined transaction id")
	resume     = flag.Bool("resume", false, "If true, resume from checkpoint; else from latest transaction.")
	lastN      = flag.Int64("last_n", -1, "Show the last N transactions only")
)

type AccountInfoRequest struct {
	Command     string `json:"command"`
	Account     string `json:"account"`
	Strict      bool   `json:"strict"`
	LedgerIndex string `json:"ledger_index"`
	Queue       bool   `json:"queue"`
}

func NewAccountInfoRequest(account string) *AccountInfoRequest {
	return &AccountInfoRequest{
		Command:     "account_info",
		Account:     account,
		Strict:      true,
		LedgerIndex: "current",
		Queue:       true,
	}
}

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
		return err
	}
	return c.WriteMessage(websocket.TextMessage, message)
}

func receive(c *websocket.Conn) (interface{}, error) {
	var i interface{}
	if err := c.ReadJSON(&i); err != nil {
		return i, err
	}
	return i, nil
}

// "meta": {
// 	"AffectedNodes": [
// 		{
// 			"ModifiedNode": {
//				"LedgerEntryType": "AccountRoot",
// 				"LedgerIndex": "A3AA57D945E845DF258BE00D4800D0372E6292C61B06AA897C09E3D15B2DCE26",
// 				"PreviousFields": {
// 					"Balance": "2160273137054",
// 					"OwnerCount": 7,
// 					"Sequence": 3551
// 				},
// 				"PreviousTxnID": "578A72EA96B5661410374039060A4FF0C8A3280809F6EBF6B07FB1CAB95F9A1B",
// 				"PreviousTxnLgrSeq": 36226594
// 			}
// 		},
//       ],
// }
func previousTxnIdAffectsAccountRoot(c *websocket.Conn) (string, error) {
	i, err := receive(c)
	if err != nil {
		return "", err
	}
	m := i.(map[string]interface{})
	b, _ := json.MarshalIndent(m, "", "  ")
	fmt.Printf("%s\n", b)
	result := m["result"].(map[string]interface{})
	meta := result["meta"].(map[string]interface{})
	affectedNodes := meta["AffectedNodes"].([]interface{})
	for _, v := range affectedNodes {
		n := v.(map[string]interface{})
		if m, ok := n["ModifiedNode"]; ok {
			s := m.(map[string]interface{})
			m, ok := s["FinalFields"]
			if !ok {
				continue
			}
			ff := m.(map[string]interface{})
			if s["LedgerEntryType"].(string) == "AccountRoot" && ff["Account"] == *acct {
				fmt.Println(ff["Balance"])
				return s["PreviousTxnID"].(string), nil
			}
		}
	}
	return "", fmt.Errorf("No previous tx id: %s", b)
}

// account_data: map[Balance:2154803734620 LedgerEntryType:AccountRoot OwnerCount:10 PreviousTxnLgrSeq:3.6147383e+07 index:A3AA57D945E845DF258BE00D4800D0372E6292C61B06AA897C09E3D15B2DCE26 Account:rspwpmBx2BhveK3Maoj29dNiSwCjZ2Vf6H PreviousTxnID:C17C9F1144CE4900A313AB5FE724712A53DF62F6FF488ACFC12371D08F8F3FED Sequence:3544 Flags:0]
func previousTxnIdInAccountData(c *websocket.Conn) (string, error) {
	i, err := receive(c)
	if err != nil {
		return "", err
	}
	m := i.(map[string]interface{})
	result := m["result"].(map[string]interface{})
	accountData := result["account_data"].(map[string]interface{})
	return accountData["PreviousTxnID"].(string), nil
}

func main() {
	flag.Parse()

	if *acct == "" {
		log.Fatal("--account no set")
	}

	u := url.URL{Scheme: "wss", Host: *addr}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	var (
		txID  string
		count int64
	)
	if *resume {
		content, err := ioutil.ReadFile(*checkpoint)
		if err != nil {
			log.Fatal(err)
		}
		if len(content) == 0 {
			log.Fatal("Emtpy ", *checkpoint)
		}
		txID = string(content)
	} else {
		if err := send(c, NewAccountInfoRequest(*acct)); err != nil {
			log.Fatal(err)
		}

		if txID, err = previousTxnIdInAccountData(c); err != nil {
			log.Fatal(err)
		}
	}
	for {
		if err := send(c, NewTxRequest(txID)); err != nil {
			log.Fatal(err)
		}
		txID, err := previousTxnIdAffectsAccountRoot(c)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(txID)
		if err := ioutil.WriteFile(*checkpoint, []byte(txID), 0644); err != nil {
			log.Fatal(err)
		}
		count++
		if *lastN > 0 && count > *lastN {
			break
		}
	}
}
