package initialize

import (
	"net/http"
	"ws/logic"
	"ws/server"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool { return true },
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	server.Instance().ListenBinary(conn)
}

func StartServer() *server.Servers {
	logic.CMap = CMap
	logic.Instance().Loads()
	svr := server.Instance()
	svr.Lobby = logic.Instance()
	return svr
}
