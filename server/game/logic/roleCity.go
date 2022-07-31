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
	"sync"
	"time"
	"xorm.io/xorm"
)

//玩家城池
type roleCity struct{
	roleCity map[int]*model2.MapRoleCity
	posRC map[int]*model2.MapRoleCity
	roleRC map[int][]*model2.MapRoleCity
	mutex sync.RWMutex
}


var RoleCity = &roleCity{
	posRC: make(map[int]*model2.MapRoleCity),
	roleCity: make(map[int]*model2.MapRoleCity),
	roleRC: make(map[int][]*model2.MapRoleCity),
}

func (r *roleCity) RoleCityServer (rid int,NickName string,session *xorm.Session) error{
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
				_, err = session.Table(city).Insert(city)

				if err != nil {
					log.Println("插入玩家城池失败",err)
					return err
				}

				//初始化城池设施·
				//TODO
				if err = CityFacilityService.TryCreate(city.CityId, rid, session); err != nil{
					log.Println("城池出错",err)
					return common.New(err.(*common.MyError).Code(), err.Error())
				}

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

func (r *roleCity) Load(){
	dbCity := make(map[int]*model2.MapRoleCity)
	err := db.Eg.Find(dbCity)
	if err != nil {
		log.Println("RoleCityService load role_city table error")
		return
	}

	//转成posCity、roleCity
	for _,v := range dbCity{
		posId := mapconfig.ToPosition(v.X,v.Y)
		r.posRC[posId] = v
		_,ok := r.roleRC[v.RId]
		if !ok {
			r.roleRC[v.RId] = make([]*model2.MapRoleCity,0)
		}
		r.roleRC[v.RId] = append(r.roleRC[v.RId],v)
	}
	//耐久度计算 后续做

}

func (r *roleCity) ScanBlock(req *model.ScanBlockReq) []model.MapRoleCity {
	x := req.X
	y := req.Y
	length := req.Length
	if x < 0 || x >= mapconfig.Mapwidth || y < 0 || y >= mapconfig.MapHeight {
		return nil
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	maxX := utils.MinInt(mapconfig.Mapwidth, x+length-1)
	maxY := utils.MinInt(mapconfig.MapHeight, y+length-1)

	rb := make([]model.MapRoleCity, 0)
	for i := x; i <= maxX; i++ {
		for j := y; j <= maxY; j++ {
			posId := mapconfig.ToPosition(i, j)
			v, ok := r.posRC[posId]
			if ok {
				rb = append(rb, v.ToModel().(model.MapRoleCity))
			}
		}
	}
	return rb
}

func (r *roleCity) IsCanBuild(x int, y int) bool {
	confs := gameconfig.MapRes.Confs
	pIndex := mapconfig.ToPosition(x,y)
	_,ok := confs[pIndex]
	if !ok {
		return false
	}

	//城池 1范围内 不能超过边界
	if x + 1 >= mapconfig.Mapwidth || y+1 >= mapconfig.MapHeight  || y-1 < 0 || x-1<0{
		return false
	}
	sysBuild := gameconfig.MapRes.SysBuild
	//系统城池的5格内 不能创建玩家城池
	for _,v := range sysBuild{
		if v.Type == gameconfig.MapBuildSysCity {
			if x >= v.X-5 &&
				x <= v.X+5 &&
				y >= v.Y-5 &&
				y <= v.Y+5{
				return false
			}
		}
	}

	//玩家城池的5格内 也不能创建城池
	for i := x-5;i<=x+5;i++ {
		for j := y-5;j<=y+5;j++ {
			posId := mapconfig.ToPosition(i,j)
			_,ok := r.posRC[posId]
			if ok {
				return false
			}
		}
	}
	return true
}

func (r *roleCity) Get(id int) (*model2.MapRoleCity, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	city,ok := r.roleCity[id]
	if ok {
		return city,ok
	}
	return nil,ok
}

func (r *roleCity) GetMainCity(rid int) (*model2.MapRoleCity, bool) {
	rcs, ok := r.roleRC[rid]
	if !ok {
		return nil, false
	}

	for _, v := range rcs {
		if v.IsMain == 1{
			return v, true
		}
	}

	return nil,false
}

func (r *roleCity) GetCityCost(cid int) int8 {
	return CityFacilityService.GetCost(cid) + gameconfig.Base.City.Cost
}