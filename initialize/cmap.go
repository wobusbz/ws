package initialize

import (
	"ws/logic"
	"ws/server"
)

var CMap map[string]logic.Handler = make(map[string]logic.Handler, 255)

func init() {
	CMap["hello"] = func(cmd *server.InternalCommand, g *logic.BusinessLogic) (ic *server.InternalCommand, err error) {
		ic = server.NewInternalCommand(cmd.CMD, 1, map[string]interface{}{
			"name": "wuhaurou",
			"addr": "cq",
		})
		return
	}

	CMap["hello1"] = func(cmd *server.InternalCommand, g *logic.BusinessLogic) (ic *server.InternalCommand, err error) {
		ic = server.NewInternalCommand(cmd.CMD, 1, map[string]interface{}{
			"name": "wuhaurou1",
			"addr": "cq1",
		})
		return
	}
}
