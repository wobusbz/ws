package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var (
	UUID = "40241925309042451117080417735442"
)

func main() {
	http_url := "ws://127.0.0.1:9912/ws"
	wsConn, _, err := websocket.DefaultDialer.Dial(http_url, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	rand.Seed(time.Now().UnixNano())
	wsConn.WriteJSON(map[string]interface{}{
		"cmd":   "hello1",
		"id":    fmt.Sprintf("%d", rand.Int63()),
		"param": "hello world",
	})

	for {
		_, p, _ := wsConn.ReadMessage()
		fmt.Println(string(p[:len(p)-1]))

	}

}
