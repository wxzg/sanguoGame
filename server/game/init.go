package game

import (
	"sanguoServer/db"
	"sanguoServer/net"
	"sanguoServer/server/game/controllor"
	"sanguoServer/server/game/gameconfig"
	"sanguoServer/server/game/logic"
)

var Router = &net.Router{}
func Init(){
	db.DBInit()
	gameconfig.Base.Load()
	gameconfig.MapBuilder.Load()
	gameconfig.General.Load()
	gameconfig.MapRes.Load()
	gameconfig.FcailityConf.Load()
	gameconfig.Skill.Load()
	logic.DefaultRoleBuildService.Load()
	logic.RoleCity.Load()
	logic.RoleAttrServer.Load()
	logic.CoalitionService.Load()
	logic.CityFacilityService.Load()
	logic.RoleServer.Load()
	initRouter()
	logic.BeforeInit()
}

func initRouter() {
	controllor.DefaultRoleRouter.Router(Router)
	controllor.DeafultMapBuildHandler.Router(Router)
	controllor.GeneralHandler.Router(Router)
	controllor.ArmyHandler.Router(Router)
	controllor.SkillHandler.Router(Router)
	controllor.InteriorController.Router(Router)
	controllor.WarHandler.InitRouter(Router)
}

