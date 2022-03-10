package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	Connect *websocket.Conn // 连接
	Msg     chan []byte     // 待发送的数据
	IsAdmit bool
}

func (client *Client) Send(msg []byte) {
	defer recover()
	client.Msg <- msg
}

func (client *Client) readLoop() {
	defer recover()
	defer (func() {
		close(client.Msg)
		fmt.Println("Disconnect:", "Remote")
	})()

	for {
		_, message, err := client.Connect.ReadMessage()
		if err != nil {
			return
		}
		interceptMessage(client, message)
	}
}

func (client *Client) writeLoop() {
	defer recover()
	defer (func() {
		client.Connect.Close()
		fmt.Println("Disconnect:", "Message channel")
	})()

	for {
		message, ok := <-client.Msg
		if !ok {
			return
		}
		client.Connect.WriteMessage(websocket.TextMessage, message)
	}
}
