package client

import "github.com/gorilla/websocket"

type Users struct {
	Addr string
	C    chan string
	Conn *websocket.Conn
}

// NewClient 初始化客户端
func NewClient(conn *websocket.Conn) (client *Users) {
	client = &Users{
		Addr: conn.RemoteAddr().String(),
		C:    make(chan string),
		Conn: conn,
	}

	go client.ListerUserMsg()
	return
}

// ListerUserMsg 监听用户消息
func (client *Users) ListerUserMsg() {
	for {
		select {
		case msg := <-client.C: // 收到用户的消息
			client.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
		}
	}
}
