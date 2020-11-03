package logic

import (
	"sync"
	"time"
	"ws/server"
)

type LogicLobby struct {
	dieout, shutdown chan bool
	logic            map[string]*BusinessLogic
	lk               sync.Mutex
}

func NewLogicLobby() *LogicLobby {
	return &LogicLobby{logic: make(map[string]*BusinessLogic, 10), dieout: make(chan bool), shutdown: make(chan bool)}
}

var _inctance *LogicLobby

func Instance() *LogicLobby {
	if _inctance == nil {
		_inctance = NewLogicLobby()
	}
	return _inctance
}

func (ll *LogicLobby) Loads() {
	go ll.start()
}

func (ll *LogicLobby) start() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			ll.BroadcastAll(nil)
		}
	}

}

func (ll *LogicLobby) Remove(id string) {
	ll.lk.Lock()
	if v, exist := ll.logic[id]; exist {
		v.Stop(nil)
		ll.logic[id] = nil
		delete(ll.logic, id)
	}
	ll.lk.Unlock()
}

func (ll *LogicLobby) GetShutdown() chan bool {
	return ll.shutdown
}

func (ll *LogicLobby) GetDieout() chan bool {
	return ll.dieout
}

func (ll *LogicLobby) Checkout(id string, start bool) (lg server.Logic) {
	ll.lk.Lock()
	if v, ok := ll.logic[id]; ok {
		lg = v
	} else {
		newLogic := NewBusinessLogic(id, ll)
		ll.logic[id] = newLogic
		if start {
			go newLogic.start()
		}
		lg = newLogic
	}
	ll.lk.Unlock()
	return
}

func (ll *LogicLobby) GetLogic(id string) (lg server.Logic) {
	ll.lk.Lock()
	if v, ok := ll.logic[id]; ok {
		lg = v
	}
	ll.lk.Unlock()
	return
}

func (ll *LogicLobby) BroadcastAll(ic *server.InternalCommand) {
	ll.lk.Lock()
	timeout := time.NewTimer(0)
	for i, _ := range ll.logic {
		if lg := ll.logic[i]; lg != nil {
			timeout.Reset(20 * time.Millisecond)
			select {
			case lg.GetOutGoing() <- server.NewInternalCommand("hello", 1, "wuhuarou"):
			case <-timeout.C:
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
	timeout.Stop()
	ll.lk.Unlock()
}
