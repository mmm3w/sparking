package wsk

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	clients    = make(map[string]*Client)
	clientLock sync.RWMutex
)

func newClient(conn *websocket.Conn) *Client {
	clientLock.Lock()
	defer clientLock.Unlock()

	client := Client{
		Connect: conn,
		Msg:     make(chan []byte, 10),
		Mark:    "",
		Name:    "",
	}
	client.Mark = client.md5()
	clients[client.Mark] = &client
	fmt.Println("Connect:", "New client")
	return &client
}

func changeName(client *Client, name string) {
	clientLock.Lock()
	defer clientLock.Unlock()
	delete(clients, client.Mark)
	client.Name = name
	clients[name] = client
}

func findClient(name string) *Client {
	clientLock.RLock()
	defer clientLock.RUnlock()
	return clients[name]
}
