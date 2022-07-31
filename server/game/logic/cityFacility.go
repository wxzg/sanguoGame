package logic

import (
	"encoding/json"
	"log"
	"sanguoServer/db"
	"sanguoServer/db/model"
	"sanguoServer/server/common"
	"sanguoServer/server/game/gameconfig"
	model1 "sanguoServer/server/game/model"
	"sanguoServer/utils"
	"sync"
	"time"
	"xorm.io/xorm"
)

var CityFacilityService = &cityFacilityService{
	facilities: make(map[int]*model.CityFacility),
	facilitiesByRId: make(map[int][]*model.CityFacility),
}
type cityFacilityService struct {
	mutex sync.RWMutex
	facilities map[int]*model.CityFacility
	facilitiesByRId map[int][]*model.CityFacility	//key:rid
}

func (c *cityFacilityService) TryCreate(cid,rid int, session *xorm.Session) error {

	cf := &model.CityFacility{}
	//查询城市设施
	ok,err := db.Eg.Table(cf).Where("cityId=?", cid).Get(cf)
	if err != nil {
		log.Println(err)
		return common.New(utils.DBError,"数据库错误")
	}
	if ok {
		return nil
	}

	//
	list := gameconfig.FcailityConf.List
	//初始化设施
	fs := make([]model1.Facility,len(list))
	for index,v := range list{
		fac := model1.Facility{
			Name: v.Name,
			Level: 0,
			Type: v.Type,
			UpTime: 0,
		}
		fs[index] = fac
	}
	fsJson,err := json.Marshal(fs)
	if err != nil {
		return common.New(utils.DBError,"转json出错")
	}
	cf.RId = rid
	cf.CityId = cid
	cf.Facilities = string(fsJson)
	if session != nil {
		session.Table(cf).Insert(cf)
	} else {
		db.Eg.Table(cf).Insert(cf)
	}
	return nil
}

func (c *cityFacilityService) GetById(rid int) ([]*model.CityFacility, error){
	cf := make([]*model.CityFacility, 0)
	err := db.Eg.Table(new(model.CityFacility)).Where("rid=?", rid).Find(&cf)

	if err != nil {
		log.Println(err)
		return cf, common.New(utils.DBError, "数据库获取失败")
	}

	return cf, nil
}

func (c *cityFacilityService) GetYield(rid int)  (y model1.Yield) {
	cfs, err := c.GetById(rid)
	if err != nil {
		for _, cf := range cfs {
			for _, f := range cf.Facility() {
				if f.GetLevel() > 0{
					values := gameconfig.FcailityConf.GetValues(f.Type, f.GetLevel())
					additions := gameconfig.FcailityConf.GetAdditions(f.Type)
					for i, aType := range additions {
						if aType == gameconfig.TypeWood {
							y.Wood += values[i]
						}else if aType == gameconfig.TypeGrain {
							y.Grain += values[i]
						}else if aType == gameconfig.TypeIron {
							y.Iron += values[i]
						}else if aType == gameconfig.TypeStone {
							y.Stone += values[i]
						}else if aType == gameconfig.TypeTax {
							y.Gold += values[i]
						}
					}
				}
			}
		}
	}
	return y
}

func (c *cityFacilityService) Get(id int) (*model.CityFacility, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	fac, ok := c.facilities[id]
	if ok {
		return fac,ok
	}
	return nil,ok
}

func (c *cityFacilityService) Load(){
	facs := make([]*model.CityFacility,0)
	err := db.Eg.Find(&facs)
	if err != nil {
		log.Println(" load city_facility table error")
	}
	for _, v:= range facs{
		c.facilities[v.CityId] = v
	}

	for _, cityFacility := range c.facilities {
		rid := cityFacility.RId
		_, ok := c.facilitiesByRId[rid]
		if ok == false {
			c.facilitiesByRId[rid] = make([]*model.CityFacility, 0)
		}
		c.facilitiesByRId[rid] = append(c.facilitiesByRId[rid], cityFacility)
	}
}

func (c *cityFacilityService) UpFacility(rid int, cid int, fType int8) (*model1.Facility, int) {
	c.mutex.RLock()
	fac, ok := c.facilities[cid]
	c.mutex.RUnlock()

	if !ok {
		return nil, utils.CityNotExist
	} else {
		facilities := make([]*model1.Facility, 0)
		var out *model1.Facility
		json.Unmarshal([]byte(fac.Facilities), &facilities)
		for _, f := range facilities{
			if f.Type == fType {
				maxLevel := gameconfig.FcailityConf.MaxLevel(fType)
				if f.CanLv() == false {
					//表示正在升级了
					return nil, utils.UpError
				}else if f.GetLevel() > maxLevel {
					//超过最大等级了，不行
					return nil, utils.UpError
				} else {
					need, ok := gameconfig.FcailityConf.Need(fType, f.GetLevel()+1)

					if !ok {
						return nil, utils.UpError
					}

					code := RoleServer.TryUseNeed(rid, *need)
					if code == utils.OK {
						f.UpTime = time.Now().Unix()
						out = f
						if t, err := json.Marshal(facilities); err == nil {
							fac.Facilities = string(t)
							fac.SyncExecute()
							return out, utils.OK
						} else {
							return nil, code
						}
					}
				}
			}
		}
	}

	return nil, utils.UpError
}

func (c *cityFacilityService) GetFacilityLv(cid int, fType int8) int8{
	cf, ok := c.facilities[cid]
	if ok {
		facility := cf.Facility()
		for _, v := range facility{
			if v.Type == fType{
				return v.GetLevel()
			}
		}
	}

	return 0
}

func (c *cityFacilityService) GetFacility(cid int, fType int8) (model1.Facility, bool) {
	cf ,ok := c.facilities[cid]
	if ok {
		facility := cf.Facility()
		for _,v := range facility{
			if v.Type == fType {
				return v,true
			}
		}
	}
	return model1.Facility{},false
}

func (c *cityFacilityService) GetCost(cid int) int8 {
	//TypeCost
	cf := c.GetByCid(cid)
	facilities := cf.Facility()
	var cost int
	for _,fa := range facilities{
		//计算等级 资源的产出是不同的
		if fa.GetLevel() > 0 {
			values := gameconfig.FcailityConf.GetValues(fa.Type,fa.GetLevel())
			adds := gameconfig.FcailityConf.GetAdditions(fa.Type)
			for i,aType := range adds{
				if aType == gameconfig.TypeCost {
					cost += values[i]
				}
			}
		}
	}
	return int8(cost)
}

func (c *cityFacilityService) GetByCid(cid int) *model.CityFacility {
	f := &model.CityFacility{}
	ok,err := db.Eg.Table(new(model.CityFacility)).Where("cityId=?",cid).Get(f)
	if err != nil {
		log.Println(err)
		return nil
	}
	if ok {
		return f
	}
	return nil
}

func (c *cityFacilityService) GetDepotCapacity(rid int) int {
	cfs, err := c.GetById(rid)
	limit := 0
	if err == nil {
		for _, cf := range cfs {
			for _, f := range cf.Facility() {
				if f.GetLevel() > 0{
					values := gameconfig.FcailityConf.GetValues(f.Type, f.GetLevel())
					additions := gameconfig.FcailityConf.GetAdditions(f.Type)
					for i, aType := range additions {
						if aType == gameconfig.TypeWarehouseLimit {
							limit += values[i]
						}
					}
				}
			}
		}
	}
	return limit
}

func (c *cityFacilityService) GetSoldier(cid int) int {
	cf, ok := c.Get(cid)

	limit := 0

	if ok{
		for _, f := range cf.Facility() {
			if f.GetLevel() > 0{
				val := gameconfig.FcailityConf.GetValues(f.Type,f.GetLevel())
				additions := gameconfig.FcailityConf.GetAdditions(f.Type)

				for i, atype := range additions {
					if atype == gameconfig.TypeSoldierLimit {
						limit += val[i]
					}
				}
			}
		}
	}

	return limit
}

