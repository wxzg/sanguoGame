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
	"time"
)

var GeneralHandler = &generalHandler{}
type generalHandler struct {

}
func (gh *generalHandler) Router(r *net.Router)  {
	g := r.Group("general")
	g.AddRouter("myGenerals",gh.myGenerals)
	g.AddRouter("drawGeneral",gh.drawGeneral)
	//征兵
	g.AddRouter("conscript",gh.conscript)
}

func (gh *generalHandler) myGenerals(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	//查询角色拥有的武将，如果没有武将就初始化随机三个武将
	rspObj := &model.MyGeneralRsp{}
	role,_ := req.Conn.GetProperty("role")
	//取出角色id
	rid := role.(*model2.Role).RId
	rsp.Body.Msg = rspObj
	rsp.Body.Code = utils.OK
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name

	gs, err := logic.DefaultGeneralService.GetGenerals(rid)

	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}

	rspObj.Generals = gs
	rsp.Body.Msg = rspObj
}

func (gh *generalHandler) drawGeneral(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	//抽武将的步骤大概是
	//1.先计算抽卡要花费的金钱，然后判断钱够不够
	//2.抽卡次数 + 已有武将 -> 卡池是否足够
	//3.随机生成武将
	//4.扣除金币

	reqObj := &model.DrawGeneralReq{}
	rspObj := &model.DrawGeneralRsp{}
	mapstructure.Decode(req.Body.Msg, reqObj)
	rsp.Body.Code = utils.OK
	rsp.Body.Msg = rspObj

	role, _ := req.Conn.GetProperty("role")
	rid := role.(*model2.Role).RId

	cost := gameconfig.Base.General.DrawGeneralCost * reqObj.DrawTimes
	//判断下金币够不够
	if !logic.RoleServer.IsEnoughGold(rid,cost) {
		rsp.Body.Code = utils.GoldNotEnough
		return
	}
	limit := gameconfig.Base.General.Limit

	gs,err := logic.DefaultGeneralService.GetGenerals(rid)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	if len(gs) + reqObj.DrawTimes > limit {
		rsp.Body.Code = utils.OutGeneralLimit
		return
	}
	mgs := logic.DefaultGeneralService.Draw(rid,reqObj.DrawTimes)
	logic.RoleServer.CostGold(rid,cost)
	rspObj.Generals = mgs
}

func (gh *generalHandler) conscript(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.ConscriptReq{}
	mapstructure.Decode(req.Body.Msg,reqObj)
	rspObj := &model.ConscriptRsp{}

	rsp.Body.Code = utils.OK
	rsp.Body.Msg = rspObj

	//检查参数
	if len(reqObj.Cnts) != 3 || reqObj.ArmyId <=0 {
		if reqObj.Cnts[0] < 0 || reqObj.Cnts[1] < 0 || reqObj.Cnts[2] < 0 {
			rsp.Body.Code = utils.InvalidParam
			return
		}
	}

	//查询角色
	r, _ := req.Conn.GetProperty("role")
	role := r.(*model.Role)

	var army *model2.Army
	//查询军队
	army, err := logic.DefaultArmyService.GetByArmyId(reqObj.ArmyId)

	if err != nil{
		rsp.Body.Code = utils.ArmyNotFound
		return
	}

	if role.RId != army.RId {
		rsp.Body.Code = utils.ArmyNotMe
		return
	}
	//判断位置是否可以征兵
	for pos,v := range reqObj.Cnts{
		if v > 0 {
			if army.Gens[pos] == nil{
				rsp.Body.Code = utils.InvalidParam
				return
			}
			if !army.PositionCanModify(pos) {
				rsp.Body.Code = utils.GeneralBusy
				return
			}
		}
	}
	//募兵所
	level := logic.CityFacilityService.GetFacilityLv(army.CityId,model.MBS)
	if level <= 0 {
		rsp.Body.Code = utils.BuildMBSNotFound
		return
	}
	//是否征兵超限制 根据武将的等级和设施的加成 计算征兵上限
	for i,g := range army.Gens{
		if g == nil {
			continue
		}
		lv := gameconfig.GeneralBasic.GetLevel(g.Level)
		if lv == nil {
			rsp.Body.Code = utils.InvalidParam
			return
		}
		add := logic.CityFacilityService.GetSoldier(army.CityId)
		if lv.Soldiers + add < reqObj.Cnts[i] + army.SoldierArray[i] {
			rsp.Body.Code = utils.OutArmyLimit
			return
		}
	}

	//开始征兵 计算消耗资源
	var total int
	for _,v := range reqObj.Cnts {
		total += v
	}
	needRes := gameconfig.NeedRes{
		Decree: 0,
		Gold: total*gameconfig.Base.ConScript.CostGold,
		Wood: total*gameconfig.Base.ConScript.CostWood,
		Iron: total*gameconfig.Base.ConScript.CostIron,
		Grain: total*gameconfig.Base.ConScript.CostGrain,
		Stone: total*gameconfig.Base.ConScript.CostStone,
	}
	//判断是否有足够的资源
	code := logic.RoleServer.TryUseNeed(role.RId,needRes)
	if code != utils.OK {
		rsp.Body.Code = code
		return
	}
	//更新部队配置
	for i,_ := range army.SoldierArray{
		var curTime = time.Now().Unix()
		if reqObj.Cnts[i] > 0 {
			army.ConscriptCntArray[i] = reqObj.Cnts[i]
			army.ConscriptTimeArray[i] = int64(reqObj.Cnts[i] * gameconfig.Base.ConScript.CostTime) + curTime -2
		}
	}
	army.Cmd =model2.ArmyCmdConscript
	army.SyncExecute()
	rspObj.Army = army.ToModel().(model.Army)
	if res,err := logic.RoleServer.GetRoleRes(role.RId);err != nil{
		rspObj.RoleRes = res.ToModel().(model.RoleRes)
	}
}