package logic

import (
	"log"
	"sanguoServer/db"
	model2 "sanguoServer/db/model"
	"sanguoServer/server/game/model"
)

type roleServer struct{

}

var RoleServer = &roleServer{}

func (r *roleServer) GetRoleRes (rid int) (model.RoleRes, error) {
	role := &model2.RoleRes{}
	ok, err := db.Eg.Table(role).Where("rid=?", rid).Get(role)
	if err != nil {
		log.Println("查询角色资源出错", err)
		return model.RoleRes{}, err
	}

	if ok {
		return role.ToModel().(model.RoleRes), nil
	}

	return model.RoleRes{},nil
}
