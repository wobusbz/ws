package server

type Logic interface {
	GetIncoming() chan *InternalCommand
	GetOutGoing() chan *InternalCommand
	GetLogout() chan string
}

type LogicLobby interface {
	GetShutdown() chan bool
	GetDieout() chan bool
	Checkout(string, bool) Logic
	Remove(string)
	GetLogic(string) Logic
	BroadcastAll(*InternalCommand)
}

type DownFlag struct {
	Status   int         `json:"status"`
	Date     interface{} `json:"date,omitempty"`
	Ts1, Ts2 int64       `json:"-"`
}

type InternalCommand struct {
	ID    string      `json:"id,omitempty"`
	CMD   string      `json:"cmd,omitempty"`
	Param interface{} `json:"param,omitempty"`
	DFlag *DownFlag   `json:"dFlag,omitempty"`
}

func NewInternalCommand(cmd string, status int, p interface{}) *InternalCommand {
	return &InternalCommand{CMD: cmd, DFlag: &DownFlag{Status: status, Date: p}}
}
