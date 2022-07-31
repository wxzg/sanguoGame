package logic

import (
	"log"
	"sanguoServer/db"
	model2 "sanguoServer/db/model"
	"sanguoServer/server/common"
	"sanguoServer/server/game/gameconfig"
	"sanguoServer/server/game/model"
	"sanguoServer/utils"
	"time"
)

type roleServer struct{
	roleRes map[int]*model2.RoleRes
}

var RoleServer = &roleServer{
	roleRes: make(map[int]*model2.RoleRes),
}

func (r *roleServer) GetRoleRes (rid int) (*model2.RoleRes, error) {
	role := &model2.RoleRes{}
	ok, err := db.Eg.Table(role).Where("rid=?", rid).Get(role)
	if err != nil {
		log.Println("查询角色资源出错", err)
		return nil, common.New(utils.DBError,"db error")
	}

	if ok {
		return role, nil
	}

	return nil,common.New(utils.DBError,"db error")
}

func (r *roleServer) GetYield (rid int) (y model.Yield) {
	//产量 = 建筑 + 城市固定收益
	rbYield := DefaultRoleBuildService.GetYield(rid)
	rcYield := CityFacilityService.GetYield(rid)

	y.Gold = rbYield.Gold + rcYield.Gold + gameconfig.Base.Role.GoldYield
	y.Stone = rbYield.Stone + rcYield.Stone + gameconfig.Base.Role.StoneYield
	y.Iron = rbYield.Iron + rcYield.Iron + gameconfig.Base.Role.IronYield
	y.Grain = rbYield.Grain + rcYield.Grain + gameconfig.Base.Role.GrainYield
	y.Wood = rbYield.Wood + rcYield.Wood +  gameconfig.Base.Role.WoodYield

	return y
}

func (r *roleServer) Get(rid int) *model2.Role {
	role := &model2.Role{}
	ok, err := db.Eg.Where("rid=?", rid).Get(role)
	if err != nil {
		log.Println("角色查询出错", err)
		return nil
	}

	if !ok {
		return nil
	}

	return role
}

func (r *roleServer) IsEnoughGold(rid int, cost int) bool {
	rr, _ := r.GetRoleRes(rid)
	return rr.Gold >= cost
}

func (r *roleServer) CostGold(rid int, cost int) {
	rr, _ := r.GetRoleRes(rid)
	if rr.Gold >= cost {
		rr.Gold -= cost
		rr.SyncExecute()
	}
}

func (r *roleServer) TryUseNeed(rid int, res gameconfig.NeedRes) int {
	rr, err := r.GetRoleRes(rid)


	if err != nil {
		if res.Decree <= rr.Decree && res.Grain <= rr.Grain && res.Stone <=rr.Stone &&
			res.Wood <= rr.Wood && res.Iron <= rr.Iron && res.Gold <= rr.Gold{
			rr.Decree -= res.Decree
			rr.Iron -= res.Iron
			rr.Wood -= res.Wood
			rr.Stone -= res.Stone
			rr.Grain -= res.Grain
			rr.Gold -= res.Gold

			rr.SyncExecute()
			return utils.OK
		}else{
			if res.Decree >rr.Decree{
				return utils.DecreeNotEnough
			}else {
				return utils.ResNotEnough
			}
		}
	} else {
		return utils.RoleNotExist
	}
}

func (r *roleServer) Load(){
	rr := make([]*model2.RoleRes, 0)
	err := db.Eg.Find(&rr)
	if err != nil {
		log.Println("load role_res error")
	}
	
	for _, v := range rr {
		r.roleRes[v.RId] = v
	}
	
	go r.produce()
}

func (r *roleServer) produce() {
	index := 1
	for true {
		t := gameconfig.Base.Role.RecoveryTime
		time.Sleep(time.Duration(t) * time.Second)

		for _, v := range r.roleRes {
			//加判断是因为爆仓了，资源不无故减少
			capacity := GetDepotCapacity(v.RId)
			yield := model.GetYield(v.RId)
			if v.Wood < capacity{
				v.Wood += utils.MinInt(yield.Wood/6, capacity)
			}

			if v.Iron < capacity{
				v.Iron += utils.MinInt(yield.Iron/6, capacity)
			}

			if v.Stone < capacity{
				v.Stone += utils.MinInt(yield.Stone/6, capacity)
			}

			if v.Grain < capacity{
				v.Grain += utils.MinInt(yield.Grain/6, capacity)
			}

			if v.Gold < capacity{
				v.Grain += utils.MinInt(yield.Grain/6, capacity)
			}

			if index%6 == 0{
				if v.Decree < gameconfig.Base.Role.DecreeLimit{
					v.Decree+=1
				}
			}
			v.SyncExecute()
		}
		index++
	}
}

func GetDepotCapacity(rid int) int{
	return CityFacilityService.GetDepotCapacity(rid) + gameconfig.Base.Role.DepotCapacity
}