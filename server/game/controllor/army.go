package controllor

import (
	"github.com/mitchellh/mapstructure"
	model2 "sanguoServer/db/model"
	"sanguoServer/net"
	"sanguoServer/server/common"
	"sanguoServer/server/game/gameconfig"
	"sanguoServer/server/game/logic"
	"sanguoServer/server/game/model"
	"sanguoServer/utils"
)

var ArmyHandler = &armyHandler{}
type armyHandler struct {

}
func (gh *armyHandler) Router(r *net.Router)  {
	g := r.Group("army")
	g.AddRouter("myList",gh.myList)
	//配置武将
	g.AddRouter("dispose",gh.dispose)
	g.AddRouter("myOne",gh.myOne)
	g.AddRouter("assign",gh.assign)
}

func (gh *armyHandler) myList(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.ArmyListReq{}
	_ = mapstructure.Decode(req.Body.Msg, reqObj)
	rspObj := &model.ArmyListRsp{}
	rsp.Body.Msg = rspObj
	rsp.Body.Code = utils.OK
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name

	role, _ := req.Conn.GetProperty("role")
	rid := role.(*model2.Role).RId
	arms , err := logic.DefaultArmyService.GetArmysByCity(rid, reqObj.CityId)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.Armys = arms
	rspObj.CityId = reqObj.CityId
	rsp.Body.Msg = rspObj
}

func (gh *armyHandler) dispose(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	//注意武将在初始化的时候必须是1级的
	reqobj := &model.DisposeReq{}
	rspobj := &model.DisposeRsp{}
	mapstructure.Decode(req.Body.Msg, reqobj)
	rsp.Body.Msg = rspobj
	rsp.Body.Code = utils.OK

	r,_ := req.Conn.GetProperty("role")
	role := r.(*model.Role)

	if reqobj.Order <=0 || reqobj.Order >5 || reqobj.Position < -1 || reqobj.Position > 2{
		rsp.Body.Code = utils.InvalidParam
		return
	}

	city, ok := logic.RoleCity.Get(reqobj.CityId)

	if !ok {
		rsp.Body.Code = utils.CityNotExist
		return
	}

	if city.RId != role.RId {
		rsp.Body.Code = utils.CityNotMe
		return
	}

	jc, ok := logic.CityFacilityService.GetFacility(city.CityId, model.JiaoChang)
	if !ok || jc.GetLevel() < reqobj.Order{
		rsp.Body.Code = utils.GeneralNotFound
		return
	}

	newGen, ok := logic.DefaultGeneralService.GetById(reqobj.GeneralId)
	if !ok {
		rsp.Body.Code = utils.GeneralNotFound
		return
	}

	if newGen.RId != role.RId {
		rsp.Body.Code = utils.GeneralNotMe
		return
	}

	army, err := logic.DefaultArmyService.GetOrCreate(role.RId, reqobj.CityId, reqobj.Order)
	if err != nil {
		rsp.Body.Code = utils.DBError
		return
	}

	if army.FromX != city.X || army.FromY != city.Y {
		rsp.Body.Code = utils.ArmyIsOutside
		return
	}

	//下阵
	if reqobj.Position == -1{
		for pos, g := range army.Gens {
			if g != nil && g.Id == newGen.Id{

				//征兵中不能下阵
				if army.PositionCanModify(pos) == false{
					if army.Cmd == model2.ArmyCmdConscript{
						rsp.Body.Code = utils.GeneralBusy
					}else{
						rsp.Body.Code = utils.ArmyBusy
					}
					return
				}

				army.GeneralArray[pos] = 0
				army.SoldierArray[pos] = 0
				army.Gens[pos] = nil
				army.SyncExecute()
				break
			}
		}
		newGen.Order = 0
		newGen.CityId = 0
		newGen.SyncExecute()
	}else{
		//征兵中不能上阵
		if army.PositionCanModify(reqobj.Position) == false{
			if army.Cmd == model2.ArmyCmdConscript{
				rsp.Body.Code = utils.GeneralBusy
			}else{
				rsp.Body.Code = utils.ArmyBusy
			}
			return
		}
		if newGen.CityId != 0{
			rsp.Body.Code = utils.GeneralBusy
			return
		}

		if logic.DefaultArmyService.IsRepeat(role.RId, newGen.CfgId) == false{
			rsp.Body.Code = utils.GeneralRepeat
			return
		}

		//判断是否能配前锋
		tst, ok := logic.CityFacilityService.GetFacility(city.CityId, model.TongShuaiTing)
		if reqobj.Position == 2 && ( ok == false || tst.GetLevel() < reqobj.Order) {
			rsp.Body.Code = utils.TongShuaiNotEnough
			return
		}

		//判断cost
		cost := gameconfig.General.Cost(newGen.CfgId)
		for i, g := range army.Gens {
			if g == nil || i == reqobj.Position {
				continue
			}
			cost += gameconfig.General.Cost(g.CfgId)
		}

		if logic.RoleCity.GetCityCost(city.CityId) < cost{
			rsp.Body.Code = utils.CostNotEnough
			return
		}

		oldG := army.Gens[reqobj.Position]
		if oldG != nil {
			//旧的下阵
			oldG.CityId = 0
			oldG.Order = 0
			oldG.SyncExecute()
		}

		//新的上阵
		army.GeneralArray[reqobj.Position] = reqobj.GeneralId
		army.Gens[reqobj.Position] = newGen
		army.SoldierArray[reqobj.Position] = 0

		newGen.Order = reqobj.Order
		newGen.CityId = reqobj.CityId
		newGen.SyncExecute()
	}

	army.FromX = city.X
	army.FromY = city.Y
	army.SyncExecute()
	//队伍
	rspobj.Army = army.ToModel().(model.Army)
}

//查询某一个部队详情
func (gh *armyHandler) myOne(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.ArmyOneReq{}
	rspObj := &model.ArmyOneRsp{}
	mapstructure.Decode(req.Body.Msg, reqObj)
	rsp.Body.Msg = rspObj
	rsp.Body.Code = utils.OK

	//角色
	r ,_ := req.Conn.GetProperty("role")
	role := r.(*model.Role)

	//获取角色城池
	city, ok := logic.RoleCity.Get(reqObj.CityId)
	if !ok {
		rsp.Body.Code = utils.CityNotExist
		return
	}

	//还要判断城池是不是当前角色的
	if role.RId != city.RId {
		rsp.Body.Code = utils.CityNotMe
		return
	}

	army := logic.DefaultArmyService.GetArmy(reqObj.CityId, reqObj.Order)
	rspObj.Army = army.ToModel().(model.Army)
}

//派遣部队
func (gh *armyHandler) assign(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.AssignArmyReq{}
	rspObj := &model.AssignArmyRsp{}
	mapstructure.Decode(req.Body.Msg, reqObj)
	rsp.Body.Msg = rspObj
	rsp.Body.Code = utils.OK

	r, _ := req.Conn.GetProperty("role")
	role := r.(*model.Role)

	army,ok := logic.DefaultArmyService.GetByArmyId(reqObj.ArmyId)
	if ok != nil{
		rsp.Body.Code = utils.ArmyNotFound
		return
	}

	if role.RId != army.RId{
		rsp.Body.Code = utils.ArmyNotMe
		return
	}

	switch reqObj.Cmd {
	case model2.ArmyCmdBack:
		rsp.Body.Code = gh.back(army) //撤退
	case model2.ArmyCmdAttack:
		rsp.Body.Code = gh.attack(army) //攻击
	case model2.ArmyCmdDefend:
		rsp.Body.Code = gh.defend(army) //驻守
	case model2.ArmyCmdReclamation:
		rsp.Body.Code = gh.reclamation(army) //屯
	case model2.ArmyCmdTransfer:
		rsp.Body.Code = gh.transfer(army) //调动
	}

	rspObj.Army = army.ToModel().(model.Army)
}

