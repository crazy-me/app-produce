package server

import (
	"fmt"
	"github.com/crazy-me/apps-produce/client"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type ws struct {
	Host      string
	Port      int
	ClientMap map[string]*client.Users
	Lock      sync.RWMutex
	BroadCast chan string
}

func (ws *ws) Start() {
	var (
		conn *websocket.Conn
		err  error
	)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			HandshakeTimeout: time.Second * 3,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err = upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrader.Upgrade err:", err)
			return
		}

		// 记录新的链接
		user := client.NewClient(conn)
		ws.Lock.Lock()
		defer ws.Lock.Unlock()
		ws.ClientMap[user.Addr] = user
		go ws.wsHandel(user)
	})
	go ws.healthy()

	err = http.ListenAndServe(fmt.Sprintf("%s:%d", ws.Host, ws.Port), nil)
	if err != nil {
		log.Println("http.ListenAndServe err:", err)
		return
	}
}

func (ws *ws) wsHandel(user *client.Users) {
	var (
		data []byte
		err  error
	)
	// 读取消息
	for {
		if _, data, err = user.Conn.ReadMessage(); err != nil {
			log.Println("conn.ReadMessage err:", err)
			delete(ws.ClientMap, user.Addr)
			user.Conn.Close()
			break
		}
		user.C <- string(data)
	}
}

// healthy 向所有客户端广播健康检测
func (ws *ws) healthy() {
	timer := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-timer.C:
			for _, v := range ws.ClientMap {
				v.Conn.WriteMessage(websocket.TextMessage, []byte("healthy..."))
			}
		}
	}
}

// New 初始化
func New(host string, port int) (wsServer *ws) {
	wsServer = &ws{
		Host:      host,
		Port:      port,
		ClientMap: make(map[string]*client.Users),
		BroadCast: make(chan string),
	}
	return
}
