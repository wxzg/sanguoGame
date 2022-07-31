package logic

import (
	"log"
	"sanguoServer/db"
	model2 "sanguoServer/db/model"
	"sanguoServer/server/common"
	"sanguoServer/server/game/gameconfig"
	"sanguoServer/server/game/mapconfig"
	"sanguoServer/server/game/model"
	"sanguoServer/utils"
	"sync"
)

var DefaultRoleBuildService = &RoleBuildService{
	posRB: make(map[int]*model2.MapRoleBuild),
	roleRB: make(map[int][]*model2.MapRoleBuild),
}
type RoleBuildService struct {
	posRB map[int]*model2.MapRoleBuild
	roleRB map[int][]*model2.MapRoleBuild
	mutex sync.RWMutex
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

func (r *RoleBuildService) Load(){
	//加载系统建筑到数据库
	total, err := db.Eg.Where("type=? or type=?", gameconfig.MapBuildSysCity, gameconfig.MapBuildSysFortress).Count(new(model2.MapRoleBuild))
	if err != nil {
		log.Println("查询出错", err)
		return
	}

	if total != int64(len(gameconfig.MapRes.SysBuild)) {
		//如果对不上，就需要将系统建筑存入数据库，先删除后插入
		db.Eg.Where("type=? or type=?",gameconfig.MapBuildSysCity, gameconfig.MapBuildFortress).Delete(new(model2.MapRoleBuild))

		for _, v:= range gameconfig.MapRes.SysBuild{
			build := model2.MapRoleBuild{
				RId: 0,
				Type: v.Type,
				Level: v.Level,
				X: v.X,
				Y:v.Y,
			}
			//初始化一下
			build.Init()
			//插入到数据库
			db.Eg.Insert(&build)
		}
	}

	//查询所有的角色建筑
	dbRb := make(map[int]*model2.MapRoleBuild)
	db.Eg.Find(dbRb)
	//将其转换为 角色id - 建筑 位置 - 建筑
	for _, v := range dbRb {
		v.Init()
		pos := mapconfig.ToPosition(v.X,v.Y)
		r.posRB[pos] = v
		_, ok := r.roleRB[v.RId]
		if !ok {
			r.roleRB[v.RId] = make([]*model2.MapRoleBuild,0)
		} else {
			r.roleRB[v.RId] = append(r.roleRB[v.RId], v)
		}
	}
}

func (r *RoleBuildService) ScanBlock(req *model.ScanBlockReq) []model.MapRoleBuild{
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

	rb := make([]model.MapRoleBuild, 0)
	for i := x; i <= maxX; i++ {
		for j := y; j <= maxY; j++ {
			posId := mapconfig.ToPosition(i, j)
			v, ok := r.posRB[posId]
			if ok && (v.RId != 0 || v.IsSysCity() || v.IsSysFortress()) {
				rb = append(rb, v.ToModel().(model.MapRoleBuild))
			}
		}
	}
	return rb
}

func (r *RoleBuildService) GetYield(rid int)  (y model.Yield) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	builds, ok := r.roleRB[rid]
	if ok {
		for _, b := range builds {
			y.Iron += b.Iron
			y.Wood += b.Wood
			y.Grain += b.Grain
			y.Stone += b.Grain
		}
	}
	return
}
