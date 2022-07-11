package logic

import (
	"log"
	"sanguoServer/db"
	model2 "sanguoServer/db/model"
	"sanguoServer/server/common"
	"sanguoServer/server/game/model"
	"sanguoServer/utils"
)

var DefaultGeneralService = &GeneralService{}
type GeneralService struct {

}


func (g *GeneralService) GetGenerals(rid int) ([]model.General,error)  {
	mrs := make([]model2.General,0)
	mr := &model2.General{}
	err := db.Eg.Table(mr).Where("rid=?",rid).Find(&mrs)
	if err != nil {
		log.Println("武将查询出错",err)
		return nil, common.New(utils.DBError,"武将查询出错")
	}
	modelMrs := make([]model.General,0)
	for _,v := range mrs{
		modelMrs = append(modelMrs,v.ToModel().(model.General))
	}
	return modelMrs,nil
}