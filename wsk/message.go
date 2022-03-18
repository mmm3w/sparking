package wsk

import (
	"encoding/json"
	"fmt"
	"sync"
)

type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

type MessageHandler func(client *Client, message []byte)

var (
	handler     = make(map[string]MessageHandler)
	handlerLock sync.RWMutex
)

func initHandler() {
	handlerLock.Lock()
	defer handlerLock.Unlock()

	handler["named"] = named     //客户端别名
	handler["forward"] = forward //转发消息
}

func getHandler(t string) MessageHandler {
	handlerLock.RLock()
	defer handlerLock.RUnlock()
	return handler[t]
}

func interceptMessage(client *Client, message []byte) {
	defer recover()

	data := &Message{}

	err := json.Unmarshal(message, data)
	if err != nil {
		fmt.Println("Data error:", err)
		return
	}

	content, err := json.Marshal(data.Content)
	if err != nil {
		fmt.Println("Content error:", err)
		return
	}

	h := getHandler(data.Type)
	if h != nil {
		h(client, content)
	} else {
		fmt.Println("Not found handler:", data.Type)
	}
}
