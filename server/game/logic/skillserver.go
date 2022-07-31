package logic

import (
	"log"
	"sanguoServer/db"
	model2 "sanguoServer/db/model"
	"sanguoServer/server/common"
	"sanguoServer/server/game/model"
	"sanguoServer/utils"
)

var DefaultSkillService = &SkillService{}
type SkillService struct {

}


func (s *SkillService) GetSkills(rid int) ([]model.Skill,error)  {
	mrs := make([]model2.Skill,0)
	mr := &model2.Skill{}
	err := db.Eg.Table(mr).Where("rid=?",rid).Find(&mrs)
	if err != nil {
		log.Println("技能查询出错",err)
		return nil, common.New(utils.DBError,"技能查询出错")
	}
	modelMrs := make([]model.Skill,0)
	for _,v := range mrs{
		modelMrs = append(modelMrs,v.ToModel().(model.Skill))
	}
	return modelMrs,nil
}