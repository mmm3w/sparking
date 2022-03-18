package wsk

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var isInited = false

func Init() {
	if isInited {
		return
	}
	isInited = true

	initHandler() //注册消息处理
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
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
