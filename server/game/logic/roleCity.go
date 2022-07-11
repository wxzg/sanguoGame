package logic

import (
	"log"
	"math/rand"
	"sanguoServer/db"
	model2 "sanguoServer/db/model"
	"sanguoServer/server/common"
	"sanguoServer/server/game"
	"sanguoServer/server/game/model"
	"sanguoServer/utils"
	"time"
)

//玩家城池
type roleCity struct{

}

var Mapwidth = 200
var MapHeight = 200
var RoleCity = &roleCity{}

func (r *roleCity) RoleCityServer (rid int,NickName string) error{
	city := &model2.MapRoleCity{}
	ok, err := db.Eg.Table(city).Where("rid=?", rid).Get(city)

	if err != nil {
		log.Println("玩家城池查询失败")
		return err
	}

	if !ok {
		city.X = rand.Intn(Mapwidth)
		city.Y = rand.Intn(MapHeight)
		city.RId = rid
		city.CreatedAt = time.Now()
		city.Name = NickName
		city.CurDurable = game.Base.City.Durable
		city.IsMain = 1
		_, err = db.Eg.Table(city).Insert(city)

		if err != nil {
			log.Println("插入玩家城池失败",err)
			return err
		}
	}

	return nil
}

func (r *roleCity) GetCitys(rid int) ([]model.MapRoleCity,error)  {
	mrs := make([]model2.MapRoleCity,0)
	mr := &model2.MapRoleCity{}
	err := db.Eg.Table(mr).Where("rid=?",rid).Find(&mrs)
	if err != nil {
		log.Println("城池查询出错",err)
		return nil, common.New(utils.DBError,"城池查询出错")
	}
	modelMrs := make([]model.MapRoleCity,len(mrs))
	for _,v := range mrs{
		modelMrs = append(modelMrs,v.ToModel().(model.MapRoleCity))
	}
	return modelMrs,nil
}