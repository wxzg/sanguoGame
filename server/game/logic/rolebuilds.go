package logic

import (
	"log"
	"sanguoServer/db"
	model2 "sanguoServer/db/model"
	"sanguoServer/server/common"
	"sanguoServer/server/game/model"
	"sanguoServer/utils"
)

var DefaultRoleBuildService = &RoleBuildService{}
type RoleBuildService struct {

}


func (r *RoleBuildService) GetBuilds(rid int) ([]model.MapRoleBuild,error)  {
	mrs := make([]model2.MapRoleBuild,0)
	mr := &model2.MapRoleBuild{}
	err := db.Eg.Table(mr).Where("rid=?",rid).Find(&mrs)
	if err != nil {
		log.Println("建筑查询出错",err)
		return nil, common.New(utils.DBError,"建筑查询出错")
	}
	modelMrs := make([]model.MapRoleBuild, 0)
	for _,v := range mrs{
		modelMrs = append(modelMrs,v.ToModel().(model.MapRoleBuild))
	}
	return modelMrs,nil
}