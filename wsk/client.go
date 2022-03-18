package wsk

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/gorilla/websocket"
)

type Client struct {
	Connect *websocket.Conn
	Msg     chan []byte
	Mark    string
	Name    string
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

func (client *Client) toBytes() []byte {
	var x reflect.SliceHeader
	x.Len = int(unsafe.Sizeof(client))
	x.Cap = int(unsafe.Sizeof(client))
	x.Data = uintptr(unsafe.Pointer(&client))
	return *(*[]byte)(unsafe.Pointer(&x))
}

func (client *Client) md5() string {
	md5h := md5.New()
	md5h.Write(client.toBytes())
	sliceh := md5h.Sum(nil)
	return hex.EncodeToString(sliceh)
}
