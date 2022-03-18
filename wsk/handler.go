package wsk

import (
	"encoding/json"
)

type NamedMsg struct {
	Name string `json:"name"`
}

type ForwardMsg struct {
	Type    string `json:"type"`    //消息类型
	Target  string `json:"target"`  //转发目标
	Content string `json:"content"` //转发内容
}

func named(client *Client, message []byte) {
	data := &NamedMsg{}

	if err := json.Unmarshal(message, data); err != nil {
		return
	}
	changeName(client, data.Name)
}

func forward(client *Client, message []byte) {
	data := &ForwardMsg{}
	if err := json.Unmarshal(message, data); err != nil {
		return
	}

	if data.Type == "text" {
		target := findClient(data.Target)
		if target != nil {
			target.Msg <- []byte(data.Content)
		}
	}
}
