package server

import (
	"errors"

	"github.com/gorilla/websocket"
)

type Servers struct {
	client *Client
	Lobby  LogicLobby
}

func NewServers() *Servers {
	return &Servers{
		newClinet(),
		nil,
	}
}

var _instance *Servers

func Instance() *Servers {
	if _instance == nil {
		_instance = NewServers()
	}
	return _instance
}

func (s *Servers) ListenBinary(conn *websocket.Conn) {
	client := NewSession(s, s.Lobby)
	client.Conn = conn
	client.reader()
}

func (s *Servers) KillClient(id string) bool {
	client := s.client.getClient(id)
	if client == nil {
		return false
	}
	client.CloseNet(errors.New("Servers::KillClient - client killed"))
	s.client.delClient(client)
	return true
}
