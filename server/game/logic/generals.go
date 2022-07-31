package logic

import (
	"encoding/json"
	"log"
	"sanguoServer/db"
	model2 "sanguoServer/db/model"
	"sanguoServer/server/common"
	"sanguoServer/server/game/gameconfig"
	"sanguoServer/server/game/model"
	"sanguoServer/utils"
	"sync"
	"time"
)

var DefaultGeneralService = &GeneralService{
	genByRole: make(map[int][]*model2.General),
	genByGId: make(map[int]*model2.General),
}
type GeneralService struct {
	mutex     sync.RWMutex
	genByRole map[int][]*model2.General
	genByGId  map[int]*model2.General
}


func (g *GeneralService) GetGenerals(rid int) ([]model.General,error)  {
	//存武将的数组
	mrs := make([]*model2.General,0)
	//武将表以及用于接收数据的数组
	mr := &model2.General{}
	//查询武将表
	err := db.Eg.Table(mr).Where("rid=?",rid).Find(&mrs)
	if err != nil {
		log.Println("武将查询出错",err)
		return nil, common.New(utils.DBError,"武将查询出错")
	}

	//如果没有武将，就随机初始化三个武将
	if len(mrs) <= 0 {
		//随机三个武将 做为初始武将
		var count = 0
		for  {
			if count >= 3 {
				break
			}
			cfgId := gameconfig.General.Rand()
			if cfgId != 0 {
				gen,err := g.NewGeneral(cfgId,rid,0)
				if err != nil {
					log.Println("生成武将出错",err)
					continue
				}
				mrs = append(mrs,gen)
				count++
			}
		}

	}

	//把武将数据库结构转换成返还给前端的结构
	modelMrs := make([]model.General,0)
	for _,v := range mrs{
		modelMrs = append(modelMrs,v.ToModel().(model.General))
	}
	return modelMrs,nil
}

const (
	GeneralNormal      	= 0 //正常
	GeneralComposeStar 	= 1 //星级合成
	GeneralConvert 		= 2 //转换
)


func (g *GeneralService) NewGeneral(cfgId int, rid int, level int8) (*model2.General, error) {
	//取出配置id
	cfg := gameconfig.General.GMap[cfgId]
	//初始化技能槽
	sa := make([]*model.GSkill,3)
	ss,_ := json.Marshal(sa)
	ge := &model2.General{
		PhysicalPower: gameconfig.Base.General.PhysicalPowerLimit,
		RId: rid,
		CfgId: cfg.CfgId,
		Order: 0,
		CityId: 0,
		Level: level,
		CreatedAt: time.Now(),
		CurArms: cfg.Arms[0],
		HasPrPoint: 0,
		UsePrPoint: 0,
		AttackDis: 0,
		ForceAdded: 0,
		StrategyAdded: 0,
		DefenseAdded: 0,
		SpeedAdded: 0,
		DestroyAdded: 0,
		Star: cfg.Star,
		StarLv: 0,
		ParentId: 0,
		SkillsArray: sa,
		Skills: string(ss),
		State: GeneralNormal,
	}
	_,err := db.Eg.Table(ge).Insert(ge)
	if err != nil {
		return nil, err
	}
	return ge , nil
}

//抽卡
func (g *GeneralService) Draw(rid int, nums int) []model.General {
	mrs := make([]*model2.General, 0)
	for i := 0; i < nums; i++ {
		cfgId := gameconfig.General.Rand()
		gen, _ := g.NewGeneral(cfgId, rid, 0)
		mrs = append(mrs, gen)
	}

	modelMrs := make([]model.General,0)
	for _, v := range mrs {
		modelMrs = append(modelMrs,v.ToModel().(model.General))
	}

	return modelMrs
}

func (g *GeneralService) GetById(id int) (*model2.General, bool) {
	gen := &model2.General{}
	ok,err := db.Eg.Table(new(model2.General)).Where("id=? and state=?",id,model2.GeneralNormal).Get(gen)
	if err != nil {
		log.Println(err)
		return nil,false
	}
	if ok {
		return gen,true
	}
	return nil,false
}

func (g *GeneralService) Load(){

	err := db.Eg.Table(model2.General{}).Where("state=?",
		model2.GeneralNormal).Find(g.genByGId)

	if err != nil {
		log.Println(err)
		return
	}

	for _, v := range g.genByGId {
		if _, ok := g.genByRole[v.RId]; ok==false {
			g.genByRole[v.RId] = make([]*model2.General, 0)
		}
		g.genByRole[v.RId] = append(g.genByRole[v.RId], v)
	}

	//开启一个协程定期恢复体力
	go g.updatePhysicalPower()
}

func (g *GeneralService) updatePhysicalPower() {
	//体力上限
	limit := gameconfig.Base.General.PhysicalPowerLimit
	//恢复体力
	recoverCnt := gameconfig.Base.General.RecoveryPhysicalPower
	for true {
		//循环每小时恢复一次体力
		time.Sleep(1*time.Hour)
		g.mutex.RLock()
		for _, gen := range g.genByGId {
			//只要没超过上限
			if gen.PhysicalPower < limit{
				//恢复体力
				gen.PhysicalPower = utils.MinInt(limit, gen.PhysicalPower+recoverCnt)
				//更新数据库
				gen.SyncExecute()
			}
		}
		g.mutex.RUnlock()
	}
}