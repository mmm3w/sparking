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

func removeClient(name string) {
	clientLock.Lock()
	defer clientLock.Unlock()
	delete(clients, name)
}

func changeName(client *Client, name string) bool {
	clientLock.Lock()
	defer clientLock.Unlock()
	if clients[name] != nil {
		return false
	}
	delete(clients, client.Mark)
	client.Name = name
	clients[name] = client
	return true
}

func findClient(name string) *Client {
	clientLock.RLock()
	defer clientLock.RUnlock()
	return clients[name]
}

// func allCertifiedClientName(name string) map[string]interface{} {
// 	clientLock.RLock()
// 	defer clientLock.RUnlock()
// 	data := make(map[string][]interface{})
// 	for _, v := range clients {
// 		if v.Name != "" && v.Name != name {
// 			data[v.Name] = getUnreadMessages(v.Name, name)
// 		}
// 	}
// 	return data
// }

func allCertifiedClientWithMsg(name string) []interface{} {
	clientLock.RLock()
	defer clientLock.RUnlock()
	data := make([]interface{}, 0)
	for _, v := range clients {
		if v.Name != "" && v.Name != name {
			data = append(data, map[string]interface{}{
				"name": v.Name,
				"msg":  getUnreadMessages(v.Name, name),
			})
		}
	}
	return data
}

func allCertifiedClient() []*Client {
	clientLock.RLock()
	defer clientLock.RUnlock()
	data := make([]*Client, 0)
	for _, v := range clients {
		if v.Name != "" {
			data = append(data, v)
		}
	}
	return data
}
