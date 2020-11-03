package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type SocketClient struct {
	ID      string
	cmd     InternalCommand
	Conn    *websocket.Conn
	servers *Servers
	lobby   LogicLobby
	logic   Logic

	recvToWriter chan int
}

func NewSession(s *Servers, lobby LogicLobby) *SocketClient {
	return &SocketClient{servers: s, lobby: lobby, recvToWriter: make(chan int, 1)}
}

func (sc *SocketClient) reader() {
	for {
		_, bufs, err := sc.Conn.ReadMessage()
		if err != nil {
			return
		}
		if err := sc.readMessage(bufs); err != nil {
			return
		}
	}
}

func (sc *SocketClient) readMessage(buf []byte) error {
	ts1 := time.Now().UnixNano()
	var ic = new(InternalCommand)
	if err := json.Unmarshal(buf, &ic); err != nil {
		return err
	}
	if sc.ID == "" {
		sc.ID = ic.ID
		sc.servers.KillClient(sc.ID)
		sc.servers.client.addClient(sc)
		sc.logic = sc.lobby.Checkout(sc.ID, true)
		go sc.writer()
	}

	lastRecvTime := time.Now().UnixNano()
	ic.DFlag = new(DownFlag)
	ic.DFlag.Ts1 = ts1
	ic.DFlag.Ts2 = lastRecvTime
	timeout := time.NewTimer(10 * time.Microsecond)
	select {
	case sc.logic.GetIncoming() <- ic:
	case <-timeout.C:
		fmt.Printf("writer quit (logout) %s, %s\n", sc.Conn.RemoteAddr(), ic.ID)
	}
	timeout.Stop()
	return nil
}

func (sc *SocketClient) writer() {
	for {
		select {
		case ic := <-sc.logic.GetOutGoing():
			if err := sc.sendData(ic); err != nil {
				sc.CloseNet2writer(err)
				return
			}
		case <-sc.recvToWriter: //recver通知writer关闭
			fmt.Printf("writer quit %s, %s \n", sc.Conn.RemoteAddr(), sc.ID)
			return
		}
	}
}

func (sc *SocketClient) sendData(ic *InternalCommand) error {
	if ic == nil {
		return errors.New("writer shutdown closing network routine")
	}
	if ic.CMD == "logout" {
		err, ok := ic.DFlag.Date.(error)
		if !ok {
			err = errors.New("writer shutdown closing network routine ... logout")
		}
		fmt.Println(err)
		return err
	}

	client := sc.servers.client.getClient(sc.ID)
	if client == nil {
		return errors.New("client write not found")
	}
	return client.Conn.WriteJSON(ic.DFlag)
}

func (sc *SocketClient) CloseNet(e error) {
	sc.Conn.Close()
	timeout := time.NewTimer(20 * time.Millisecond)
	select {
	case sc.recvToWriter <- 0:
	case <-timeout.C:
	}
	timeout.Stop()
	sc.servers.client.delClient(sc)
}

func (c *SocketClient) CloseNet2writer(e error) {
	c.Conn.Close()
	c.servers.client.delClient(c)
}
