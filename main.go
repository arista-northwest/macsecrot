package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aristanetworks/goeapi"
)

/* eAPI JSON reference for Go structures
{
    "interfaces": {
        "Ethernet3/1/1": {
            "portId": 455,
            "details": {
                "oldKeyTransmitting": true,
                "keyServerPriority": 4,
                "oldKeyReceiving": true,
                "oldKeyMsgId": "3e48b33fa2dae48c539b8c3c",
                "fipsPostStatus": "none",
                "oldKeyMsgNum": 769,
                "latestKeyMsgId": "",
                "latestKeyReceiving": false,
                "latestKeyTransmitting": false,
                "traffic": "Protected",
                "localSsci": "00000001",
                "latestKeyMsgNum": 0,
                "sessionReKeyPeriod": 1800.0
            },
            "address": "28:99:3a:82:68:62",
            "controlledPort": true,
            "keyNum": 769,
            "keyMsgId": "3e48b33fa2dae48c539b8c3c"
        },
        "Ethernet3/1/2": {
            "portId": 456,
            "details": {
                "oldKeyTransmitting": true,
                "keyServerPriority": 4,
                "oldKeyReceiving": true,
                "oldKeyMsgId": "1a7dd1cf2c3362a48e21f261",
                "fipsPostStatus": "none",
                "oldKeyMsgNum": 769,
                "latestKeyMsgId": "",
                "latestKeyReceiving": false,
                "latestKeyTransmitting": false,
                "traffic": "Protected",
                "localSsci": "00000001",
                "latestKeyMsgNum": 0,
                "sessionReKeyPeriod": 1800.0
            },
            "address": "28:99:3a:82:68:5c",
            "controlledPort": true,
            "keyNum": 769,
            "keyMsgId": "1a7dd1cf2c3362a48e21f261"
        }
    }
}
*/

// type MacSecInterfaces struct {
// 	Interfaces map[string]Interface `json:"interfaces"`
// }
//
// type Interface struct {
// 	PortId         int    `json:"portId"`
// 	Details        Detail `json:"details"`
// 	Address        string `json:"address"`
// 	ControlledPort bool   `json:"controlledPort"`
// 	KeyNum         int    `json:"keyNum"`
// 	KeyMsgId       string `json:"keyMsgId"`
// }
//
// type Detail struct {
// 	OldKeyTransmitting    bool    `json:"oldKeyTransmitting"`
// 	KeyServerPriority     int     `json:"keyServerPriority"`
// 	OldKeyReceiving       bool    `json:"oldKeyReceiving"`
// 	OldKeyMsgId           string  `json:"oldKeyMsgId"`
// 	FipsPostStatus        string  `json:"fipsPostStatus"`
// 	OldKeyMsgNum          int     `json:"oldKeyMsgNum"`
// 	LatestKeyMsgId        string  `json:"latestKeyMsgId"`
// 	LatestKeyReceiving    bool    `json:"latestKeyReceiving"`
// 	LatestKeyTransmitting bool    `json:"latestKeyTransmitting"`
// 	Traffic               string  `json:"traffic"`
// 	LocalSsci             string  `json:"localSsci"`
// 	LatestKeyMsgNum       int     `json:"latestKeyMsgNum"`
// 	SessionReKeyPeriod    float64 `json:"sessionReKeyPeriod"`
// }
//
// func (s *MacSecInterfaces) GetCmd() string {
// 	return "show mac security interface detail"
// }
//
// /*
// #show mac security status | json
// {
//     "securedInterfaces": 2,
//     "licenseEnabled": true,
//     "fipsMode": false,
//     "activeProfiles": 1,
//     "delayProtection": false
// }
// */
//
// type MacSecStatus struct {
// 	SecuredInterfaces int  `json:"securedInterfaces"`
// 	LicenseEnabled    bool `json:"licenseEnabled"`
// 	FipsMode          bool `json:"fipsMode"`
// 	ActiveProfiles    int  `json:"activeProfiles"`
// 	DelayProtection   bool `json:"delayProtection"`
// }
//
// func (s *MacSecStatus) GetCmd() string {
// 	return "show mac security status"
// }

// func main() {
// 	node, err := goeapi.Connect("http", "tg219", "admin", "", 80)
// 	if err != nil {
// 		panic(err)
// 	}
// 	node.EnableAuthentication("")
//
// 	// Sample show commands...
// 	handle, _ := node.GetHandle("json")
//
// 	msec_intfs := &MacSecInterfaces{}
// 	msec_status := &MacSecStatus{}
// 	handle.AddCommand(msec_intfs)
// 	handle.AddCommand(msec_status)
//
// 	if err := handle.Call(); err != nil {
// 		panic(err)
// 	}
//
// 	fmt.Print("\n\n## MAC Sec Interfaces ##\n")
// 	for k, v := range msec_intfs.Interfaces {
// 		fmt.Printf("Interface: %s\n", k)
// 		fmt.Printf("  Address: %s\n", v.Address)
// 		fmt.Printf("  KeyNum: %d\n", v.KeyNum)
// 	}
//
// 	// Sample "show running-config section mac security"
// 	sect, err := node.GetSection("mac security", "")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Print("\n\n## MAC Sec Configuration ##\n")
// 	fmt.Printf("%s\n", sect)
// }

type Message struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func handler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	// case http.MethodGet:
	// 	// TODO: Maybe get macsec status??
	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		var msg Message
		err = json.Unmarshal(body, &msg)

		log.Println(msg)
		node, err := goeapi.Connect("socket", "localhost", "admin", "", 0)
		if err != nil {
			panic(err)
		}
		node.EnableAuthentication("") // no enable secret

		ok := node.Config("interface "+msg.Name, "description "+msg.Description)
		if !ok {
			panic("not ok")
		}
	default:
		// Nothing to see here...
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":50001", nil))
}
