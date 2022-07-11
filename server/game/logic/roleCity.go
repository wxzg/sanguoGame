package logic

import (
	"log"
	"math/rand"
	"sanguoServer/db"
	model2 "sanguoServer/db/model"
	"sanguoServer/server/common"
	"sanguoServer/server/game/gameconfig"
	"sanguoServer/server/game/mapconfig"
	"sanguoServer/server/game/model"
	"sanguoServer/utils"
	"time"
)

//玩家城池
type roleCity struct{

}


var RoleCity = &roleCity{}

func (r *roleCity) RoleCityServer (rid int,NickName string) error{
	city := &model2.MapRoleCity{}
	ok, err := db.Eg.Table(city).Where("rid=?", rid).Get(city)

	if err != nil {
		log.Println("玩家城池查询失败")
		return err
	}

	if !ok {
		for {

			city.X = rand.Intn(mapconfig.Mapwidth)
			city.Y = rand.Intn(mapconfig.MapHeight)

			//这里还要判断下城池可不可以创建
			if IsCanBuild(city.X, city.Y){

				city.RId = rid
				city.CreatedAt = time.Now()
				city.Name = NickName
				city.CurDurable = gameconfig.Base.City.Durable
				city.IsMain = 1
				_, err = db.Eg.Table(city).Insert(city)

				if err != nil {
					log.Println("插入玩家城池失败",err)
					return err
				}

				//初始化城池设施·
				//TODO
				break
			}
		}

	}

	return nil
}

func IsCanBuild(x int, y int) bool {
	sysBuild := gameconfig.MapRes.SysBuild
	confs := gameconfig.MapRes.Confs
	pIndex := mapconfig.ToPosition(x, y)
	_, ok := confs[pIndex]
	if !ok {
		return false
	}

	for _, v := range sysBuild{
		if v.Type == gameconfig.MapBuildSysCity{
			//5格内不能有玩家城池
			if x >= v.X-5 &&
				x<=v.X + 5 &&
				y >= v.Y-5 &&
				y <= v.Y + 5{
				return false
			}
		}
	}
	return true
}

func (r *roleCity) GetCitys(rid int) ([]model.MapRoleCity,error)  {
	mrs := make([]model2.MapRoleCity,0)
	mr := &model2.MapRoleCity{}
	err := db.Eg.Table(mr).Where("rid=?",rid).Find(&mrs)
	if err != nil {
		log.Println("城池查询出错",err)
		return nil, common.New(utils.DBError,"城池查询出错")
	}
	modelMrs := make([]model.MapRoleCity,0)
	for _,v := range mrs{
		modelMrs = append(modelMrs,v.ToModel().(model.MapRoleCity))
	}
	return modelMrs,nil
}