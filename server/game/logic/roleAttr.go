package logic

import (
	"log"
	"sanguoServer/db"
	model2 "sanguoServer/db/model"
)

type roleAttrServer struct{

}

var RoleAttrServer = &roleAttrServer{}

func (r *roleAttrServer) RoleAttributeHandler (rid int) error{
	rr := &model2.RoleAttribute{}
	ok, err := db.Eg.Table(rr).Where("rid=?", rid).Get(rr)
	if err != nil {
		log.Println("玩家属性查询出错：",err)
		return err
	}

	if !ok {
		rr.RId = rid
		rr.ParentId = 0
		rr.UnionId = 0
		_, err = db.Eg.Table(rr).Insert(rr)
		if err != nil {
			log.Println("玩家属性插入出错",err)
			return err
		}
	}

	return nil
}