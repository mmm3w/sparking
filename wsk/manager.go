package main

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	clients              = make(map[string]*Client) //client管理
	clientLock           sync.RWMutex
	clientAdmitState     = make(map[*Client]bool)
	clientAdmitStateLock sync.RWMutex
)

func newClient(conn *websocket.Conn) *Client {
	clientAdmitStateLock.Lock()
	defer clientAdmitStateLock.Unlock()
	client := Client{
		Connect: conn,
		Msg:     make(chan []byte, 10),
		IsAdmit: false,
	}
	clientAdmitState[&client] = false
	fmt.Println("Connect:", "New client")
	return &client
}

func admitClient(client *Client, flag string) {
	clientLock.Lock()
	defer clientLock.Unlock()

	// client[]
}
