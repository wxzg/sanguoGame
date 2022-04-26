package gate

import (
	"sanguoServer/net"
	"sanguoServer/server/gate/gatecontrollor"
)


var Router = &net.Router{}
func InitGate(){
	initRouter()
}

func initRouter(){
	gatecontrollor.GateHandler.Router(Router)
}
