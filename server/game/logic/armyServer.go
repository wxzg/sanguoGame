package logic

import (
	"log"
	"sanguoServer/db"
	model2 "sanguoServer/db/model"
	"sanguoServer/server/common"
	"sanguoServer/server/game/mapconfig"
	"sanguoServer/server/game/model"
	"sanguoServer/utils"
	"sync"
)

var DefaultArmyService = &ArmyService{
	passByPosArmys: make(map[int]map[int]*model2.Army),
}
type ArmyService struct {
	passBy sync.RWMutex
	passByPosArmys map[int]map[int]*model2.Army //玩家路过位置的军队 posID - armyID
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

func (g *ArmyService) GetArmysByCity(rid,cid int) ([]model.Army,error)  {
	mrs := make([]model2.Army,0)
	mr := &model2.Army{}
	err := db.Eg.Table(mr).Where("rid=? and cityId=?",rid,cid).Find(&mrs)
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

func (a *ArmyService) ScanBlock(roleId int, req *model.ScanBlockReq) []model.Army {
	x := req.X
	y := req.Y
	length := req.Length
	if x < 0 || x >= mapconfig.Mapwidth || y < 0 || y >= mapconfig.MapHeight {
		return nil
	}

	maxX := utils.MinInt(mapconfig.Mapwidth, x+length-1)
	maxY := utils.MinInt(mapconfig.MapHeight, y+length-1)
	out := make([]model.Army, 0)

	a.passBy.RLock()
	for i := x; i <= maxX; i++ {
		for j := y; j <= maxY; j++ {

			posId := mapconfig.ToPosition(i, j)
			armys, ok := a.passByPosArmys[posId]
			if ok {
				//是否在视野范围内
				is := armyIsInView(roleId, i, j)
				if is == false{
					continue
				}
				for _, army := range armys {
					out = append(out, army.ToModel().(model.Army))
				}
			}
		}
	}
	a.passBy.RUnlock()
	return out
}

func (g *ArmyService) Get(rid,cid int) ([]*model2.Army,error)  {
	mrs := make([]*model2.Army,0)
	mr := &model2.Army{}
	err := db.Eg.Table(mr).Where("rid=? and cityId=?",rid,cid).Find(&mrs)
	if err != nil {
		log.Println("军队查询出错",err)
		return nil, common.New(utils.DBError,"军队查询出错")
	}
	for _,v := range mrs{
		v.CheckConscript()
		g.updateGenerals(v)
	}
	return mrs,nil
}

func (g *ArmyService) GetOrCreate(rid int, cid int, order int8) (*model2.Army, error) {
	armys, err := g.Get(rid, cid)
	if err != nil {
		return nil, err
	}

	for _, v := range armys {
		if v.Order == order{
			return v, nil
		}
	}
	//需要创建
	army := &model2.Army{
		RId: rid,
		Order: order,
		CityId: cid,
		Generals: `[0,0,0]`,
		Soldiers: `[0,0,0]`,
		GeneralArray: []int{0,0,0},
		SoldierArray: []int{0,0,0},
		ConscriptCnts: `[0,0,0]`,
		ConscriptTimes: `[0,0,0]`,
		ConscriptCntArray: []int{0,0,0},
		ConscriptTimeArray: []int64{0,0,0},
	}

	city, ok := RoleCity.Get(cid)
	if ok {
		army.FromX = city.X
		army.FromY = city.Y
		army.ToX = city.X
		army.ToY = city.Y
	}
	g.updateGenerals(army)
	_, err = db.Eg.Insert(army)
	if err == nil{
		return army, nil
	}else{
		log.Println("db error", err)
		return nil, err
	}
}

func (a *ArmyService) updateGenerals(armys... *model2.Army) {
	for _, army := range armys {
		army.Gens = make([]*model2.General, 0)
		for _, gid := range army.GeneralArray {
			if gid == 0{
				army.Gens = append(army.Gens, nil)
			}else{
				g, _ := DefaultGeneralService.GetById(gid)
				army.Gens = append(army.Gens, g)
			}
		}
	}
}

func (g *ArmyService) IsRepeat(id int, cfgId int) bool {
	armys , err := g.FindArmys(id)
	if err != nil {
		return true
	}
	for _, army := range armys {
		for _, g := range army.Gens {
			if g != nil {
				if g.CfgId == cfgId && g.CityId != 0{
					return false
				}
			}
		}
	}
	return true
}

func (g *ArmyService) GetByArmyId(id int) (*model2.Army, error) {

	mr := &model2.Army{}
	ok, err := db.Eg.Table(mr).Where("id=?",id).Get(mr)


	if err != nil {
		log.Println("军队查询出错",err)
		return nil, common.New(utils.DBError,"军队查询出错")
	}
	if ok {
		mr.CheckConscript()
		g.updateGenerals(mr)
		return mr,nil
	}

	return nil, common.New(utils.DBError,"军队查询出错")
}

func (g *ArmyService) FindArmys(rid int) ([]*model2.Army, error) {
	mrs := make([]*model2.Army,0)
	mr := &model2.Army{}
	err := db.Eg.Table(mr).Where("rid=?",rid).Find(&mrs)
	if err != nil {
		log.Println("军队查询出错",err)
		return nil, common.New(utils.DBError,"军队查询出错")
	}
	for _,v := range mrs{
		v.CheckConscript()
		g.updateGenerals(v)
	}
	return mrs,nil
}

func (g *ArmyService) GetArmy(cid int, order int8) *model2.Army {
	army := &model2.Army{}
	ok,err := db.Eg.Table(army).Where("cityId=? and a_order=?",cid,order).Get(army)
	if err != nil {
		log.Println("armyService GetArmy err",ok)
		return nil
	}
	if ok {
		//还需要做一步操作  检测一下是否征兵完成
		army.CheckConscript()
		g.updateGenerals(army)
		return army
	}
	return nil
}

func armyIsInView(rid, x, y int) bool {
	//简单点 先设为true
	return true
}


