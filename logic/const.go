package logic

import "ws/server"

type Handler func(cmd *server.InternalCommand, g *BusinessLogic) (ic *server.InternalCommand, err error)

// 等待外部注册
var (
	CMap map[string]Handler
)
