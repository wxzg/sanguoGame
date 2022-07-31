package logic

import (
	"sanguoServer/server/game/model"
)

func BeforeInit(){
	model.GetYield = RoleServer.GetYield
	model.GetUnion = RoleAttrServer.GetUnion
}
