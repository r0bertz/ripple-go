package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gorilla/websocket"
	"github.com/r0bertz/ripple/tx"
	"github.com/r0bertz/ripple/wss"
)

var (
	addr       = flag.String("addr", "s2.ripple.com:443", "WSS service address")
	acct       = flag.String("account", "", "Ripple address")
	checkpoint = flag.String("checkpoint", "checkpoint", "File contains the last examined transaction id")
	resume     = flag.Bool("resume", false, "If true, resume from checkpoint; else from latest transaction.")
	lastN      = flag.Int64("last_n", -1, "Show the last N transactions only")
)

type accountInfoRequest struct {
	Command     string `json:"command"`
	Account     string `json:"account"`
	Strict      bool   `json:"strict"`
	LedgerIndex string `json:"ledger_index"`
	Queue       bool   `json:"queue"`
}

func newAccountInfoRequest(account string) *accountInfoRequest {
	return &accountInfoRequest{
		Command:     "account_info",
		Account:     account,
		Strict:      true,
		LedgerIndex: "current",
		Queue:       true,
	}
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
func previousTxnIDAffectsAccountRoot(c *websocket.Conn) (string, error) {
	i, err := wss.Receive(c)
	if err != nil {
		return "", err
	}
	m := i.(map[string]interface{})
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
				return s["PreviousTxnID"].(string), nil
			}
		}
	}
	b, _ := json.MarshalIndent(m, "", "  ")
	return "", fmt.Errorf("No previous tx id: %s", b)
}

// account_data: map[
//	Balance:2154803734620
//	LedgerEntryType:AccountRoot
//	OwnerCount:10
//	PreviousTxnLgrSeq:3.6147383e+07
//	index:A3AA57D945E845DF258BE00D4800D0372E6292C61B06AA897C09E3D15B2DCE26
//	Account:rspwpmBx2BhveK3Maoj29dNiSwCjZ2Vf6H
//	PreviousTxnID:C17C9F1144CE4900A313AB5FE724712A53DF62F6FF488ACFC12371D08F8F3FED
//	Sequence:3544
//	Flags:0
// ]
func previousTxnIDInAccountData(c *websocket.Conn) (string, error) {
	i, err := wss.Receive(c)
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

	c, _, err := wss.Connect(*addr)
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
			log.Fatalf("Empty checkpoint file %s", *checkpoint)
		}
		txID = string(content)
	} else {
		if err := wss.Send(c, newAccountInfoRequest(*acct)); err != nil {
			log.Fatal(err)
		}

		if txID, err = previousTxnIDInAccountData(c); err != nil {
			log.Fatal(err)
		}
	}
	for {
		if err := wss.Send(c, tx.NewRequest(txID)); err != nil {
			log.Fatal(err)
		}
		txID, err = previousTxnIDAffectsAccountRoot(c)
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
