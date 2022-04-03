package wsk

import (
	"encoding/json"
)

type NamedMsg struct {
	Name string `json:"name"`
}

type ForwardMsg struct {
	Type    string `json:"type"`
	Target  string `json:"target"`
	Content string `json:"content"`
}

/*****************************************************************/

type OutMsg struct {
	Type    string `json:"type"`
	From    string `json:"from"`
	Content string `json:"content"`
}

/*****************************************************************/
/* {"cmd":"named","data":{"name":""}} */
func named(client *Client, message []byte) {
	data := &NamedMsg{}
	if err := json.Unmarshal(message, data); err != nil {
		return
	}
	if data.Name == "" {
		return
	}

	oldName := client.Name

	if changeName(client, data.Name) {
		online(data.Name, oldName)
		clientList(client)
	}
}

/* {"cmd":"forward","data":{"type":"text","target":"","content":""  }} */
func forward(client *Client, message []byte) {
	if client.Name == "" {
		return
	}

	data := &ForwardMsg{}
	if err := json.Unmarshal(message, data); err != nil {
		return
	}
	if data.Target == "" || data.Target == client.Name {
		return
	}

	result, err := json.Marshal(Resp{
		Event: "forward",
		Content: OutMsg{
			Type:    data.Type,
			From:    client.Name,
			Content: data.Content,
		},
	})

	target := findClient(data.Target)
	if err == nil && target != nil {
		target.Msg <- result
	}
}

/*****************************************************************/
func clientList(client *Client) {
	resp := Resp{
		Event:   "clients",
		Content: allCertifiedClientWithMsg(client.Name),
		// Content: make([]interface{}, 0),
	}
	result, err := json.Marshal(resp)
	if err == nil {
		client.Msg <- result
	}
}

func offline(name string) {
	resp, err := json.Marshal(Resp{
		Event:   "offline",
		Content: name,
	})
	if err == nil {
		for _, v := range allCertifiedClient() {
			if name != v.Name {
				v.Msg <- resp
			}
		}
	}
}

func online(name string, offline string) {

	var offlineResult []byte
	var err error
	if offline != "" {
		offlineResult, err = json.Marshal(Resp{
			Event:   "offline",
			Content: offline,
		})
		if err != nil {
			offlineResult = nil
		}
	} else {
		offlineResult = nil
	}

	for _, v := range allCertifiedClient() {
		if name != v.Name {
			onlineResult, err := json.Marshal(Resp{
				Event:   "online",
				Content: map[string]interface{}{"name": name, "msg": getUnreadMessages(name, v.Name)},
			})
			if err == nil {
				v.Msg <- onlineResult
				if offlineResult != nil {
					v.Msg <- offlineResult
				}
			}
		}
	}
}
