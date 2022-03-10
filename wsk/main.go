package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	//协议升级
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true }, // 解决跨域问题
	}).Upgrade(w, r, nil)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	//构建client
	client := newClient(conn)
	//开启读写循环
	go client.readLoop()
	go client.writeLoop()
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Index Page")
}

func main() {

	initHandler() //注册消息处理

	http.HandleFunc("/socket", wsHandler)
	http.HandleFunc("/", home)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		fmt.Println(err)
	}
}
