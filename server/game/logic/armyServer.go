package logic

import (
	"log"
	"sanguoServer/db"
	model2 "sanguoServer/db/model"
	"sanguoServer/server/common"
	"sanguoServer/server/game/model"
	"sanguoServer/utils"
)

var DefaultArmyService = &ArmyService{}
type ArmyService struct {

}

func (g *ArmyService)GetArmys(rid int) ([]model.Army, error){
	mrs := make([]model2.Army, 0)
	mr := &model2.Army{}

	err := db.Eg.Table(mr).Where("rid=?", rid).Find(&mrs)
	if err != nil {
		log.Println("军队查询出错",err)
		return nil, common.New(utils.DBError,"军队查询出错")
	}
	modelMrs := make([]model.Army,0)
	for _,v := range mrs{
		modelMrs = append(modelMrs,v.ToModel().(model.Army))
	}
	return modelMrs,nil
}