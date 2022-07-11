package game

import (
	"sanguoServer/db"
	"sanguoServer/net"
	"sanguoServer/server/game/controllor"
	"sanguoServer/server/game/gameconfig"
)

var Router = &net.Router{}
func Init(){
	db.DBInit()
	gameconfig.Base.Load()
	gameconfig.MapBuilder.Load()
	initRouter()
}

func initRouter() {
	controllor.DefaultRoleRouter.Router(Router)
	controllor.DeafultMapBuildHandler.Router(Router)
}

