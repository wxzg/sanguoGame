package game

import (
	"sanguoServer/db"
	"sanguoServer/net"
	"sanguoServer/server/game/controllor"
	"sanguoServer/server/game/gameconfig"
)

var Router = &net.Router{}
var Base = &gameconfig.Basic{}
func Init(){
	db.DBInit()
	Base.Load()
	gameconfig.MapBuilder.Load()
	initRouter()
}

func initRouter() {
	controllor.DefaultRoleRouter.Router(Router)
	controllor.DeafultMapBuildHandler.Router(Router)
}

