package logic

import (
	"log"
	"sanguoServer/db"
	model2 "sanguoServer/db/model"
	"sanguoServer/server/common"
	"sanguoServer/server/game/model"
	"sanguoServer/utils"
)

var DefaultWarService = &WarService{}
type WarService struct {

}


func (w *WarService) GetWarReports(rid int) ([]model.WarReport,error)  {
	mrs := make([]model2.WarReport,0)
	mr := &model2.WarReport{}
	err := db.Eg.Table(mr).
		Where("a_rid=? or d_rid",rid,rid).
		Desc("ctime").
		Limit(30,0).
		Find(&mrs)
	if err != nil {
		log.Println("战报查询出错",err)
		return nil, common.New(utils.DBError,"战报查询出错")
	}
	modelMrs := make([]model.WarReport,0)
	for _,v := range mrs{
		modelMrs = append(modelMrs,v.ToModel().(model.WarReport))
	}
	return modelMrs,nil
}
