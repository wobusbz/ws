package logic

import (
	"errors"
	"sync"
	"time"
	"ws/server"
)

type BusinessLogic struct {
	ID     string
	Lobby  *LogicLobby
	logout chan string
	//***** 加载数据 **** //
	// TODO DB

	me                 sync.Mutex
	lastRecvTime       int64 //最近一次收到消息的时间
	lastSaveTime       int64 // 最近一次保存的时间
	incoming, outgoing chan *server.InternalCommand
}

func NewBusinessLogic(id string, lb *LogicLobby) *BusinessLogic {
	return &BusinessLogic{ID: id, Lobby: lb, incoming: make(chan *server.InternalCommand, 1000), outgoing: make(chan *server.InternalCommand, 1000)}
}

func (lb *BusinessLogic) start() {

	for {
		select {
		case msg := <-lb.incoming:
			if msg == nil {
				lb.Lobby.Remove(lb.ID) // 锁
				return
			}
			ic, err := lb.Dispatch(msg) // 处理客户端上行指令
			if err != nil {
				lb.Stop(err)
				lb.Lobby.Remove(lb.ID)
				return
			}
			if ic != nil {
				lb.GetOutGoing() <- ic
			}
		}
	}

}

func (lb *BusinessLogic) Stop(err error) {
	if err != nil {
		ic := server.NewInternalCommand("logout", 999, err)
		timeout := time.NewTimer(20 * time.Millisecond)
		select {
		case lb.GetOutGoing() <- ic:
		case <-timeout.C:
		}
		timeout.Stop()
	}
}

func (lb *BusinessLogic) GetIncoming() chan *server.InternalCommand {
	return lb.incoming
}

func (lb *BusinessLogic) GetOutGoing() chan *server.InternalCommand {
	return lb.outgoing
}

func (lb *BusinessLogic) GetLogout() chan string {
	return lb.logout
}

func (lb *BusinessLogic) Dispatch(cmd *server.InternalCommand) (ic *server.InternalCommand, err error) {
	handler, ok := CMap[cmd.CMD]
	if !ok {
		return nil, errors.New("up cmdindex not found")
	}
	ic, err = handler(cmd, lb)
	return
}
