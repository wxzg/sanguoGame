package controllor

import (
	"github.com/mitchellh/mapstructure"
	model2 "sanguoServer/db/model"
	"sanguoServer/net"
	"sanguoServer/server/game/gameconfig"
	"sanguoServer/server/game/logic"
	"sanguoServer/server/game/middleware"
	"sanguoServer/server/game/model"
	"sanguoServer/utils"
	"time"
)

var InteriorController = &interiorController{}
type interiorController struct {

}

func (i *interiorController) Router(router *net.Router) {
	g := router.Group("interior")
	g.Use(middleware.Log())
	g.AddRouter("openCollect",i.openCollect,middleware.CheckRole())
	g.AddRouter("collect",i.collect,middleware.CheckRole())
	g.AddRouter("transform",i.transform,middleware.CheckRole())
}

func (i *interiorController) openCollect(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	rspObj := &model.OpenCollectionRsp{}

	rsp.Body.Msg = rspObj
	rsp.Body.Code = utils.OK

	r, _ := req.Conn.GetProperty("role")
	role := r.(*model2.Role)
	roleAttr, err:= logic.RoleAttrServer.Get(role.RId)
	if err != nil {
		rsp.Body.Code = utils.DBError
		return
	}

	interval := gameconfig.Base.Role.CollectInterval
	timeLimit := gameconfig.Base.Role.CollectTimesLimit
	//征收上限
	rspObj.Limit = timeLimit
	rspObj.CurTimes = roleAttr.CollectTimes
	if roleAttr.LastCollectTime.IsZero() {
		rspObj.NextTime = 0
	}else{
		if roleAttr.CollectTimes >= timeLimit {
			//第二天才能征收
			y, m, d := roleAttr.LastCollectTime.Add(24*time.Hour).Date()
			nextTime := time.Date(y, m, d, 0, 0, 0, 0, time.FixedZone("CST", 8*3600))
			rspObj.NextTime = nextTime.UnixNano()/1e6
		}else{
			nextTime := roleAttr.LastCollectTime.Add(time.Duration(interval)*time.Second)
			rspObj.NextTime = nextTime.UnixNano()/1e6
		}
	}
}

func (i *interiorController) collect(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	rspObj := &model.CollectionRsp{}

	rsp.Body.Msg = rspObj
	rsp.Body.Code = utils.OK

	r, _ := req.Conn.GetProperty("role")
	role := r.(*model.Role)
	roleRes, err := logic.RoleServer.GetRoleRes(role.RId)
	if err != nil {
		rsp.Body.Code = utils.DBError
		return
	}

	roleAttr, err:= logic.RoleAttrServer.Get(role.RId)
	if err != nil {
		rsp.Body.Code = utils.DBError
		return
	}

	curTime := time.Now()
	lastTime := roleAttr.LastCollectTime
	if curTime.YearDay() != lastTime.YearDay() || curTime.Year() != lastTime.Year(){
		roleAttr.CollectTimes = 0
		roleAttr.LastCollectTime = time.Time{}
	}

	timeLimit := gameconfig.Base.Role.CollectTimesLimit
	//是否超过征收次数上限
	if roleAttr.CollectTimes >= timeLimit{
		rsp.Body.Code = utils.OutCollectTimesLimit
		return
	}

	//cd内不能操作
	need := lastTime.Add(time.Duration(gameconfig.Base.Role.CollectTimesLimit)*time.Second)
	if curTime.Before(need){
		rsp.Body.Code = utils.InCdCanNotOperate
		return
	}

	gold := logic.RoleServer.GetYield(roleRes.RId).Gold
	rspObj.Gold = gold
	roleRes.Gold += gold
	roleRes.SyncExecute()

	roleAttr.LastCollectTime = curTime
	roleAttr.CollectTimes += 1
	roleAttr.SyncExecute()

	interval := gameconfig.Base.Role.CollectInterval
	if roleAttr.CollectTimes >= timeLimit {
		y, m, d := roleAttr.LastCollectTime.Add(24*time.Hour).Date()
		nextTime := time.Date(y, m, d, 0, 0, 0, 0, time.FixedZone("IST", 3600))
		rspObj.NextTime = nextTime.UnixNano()/1e6
	}else{
		nextTime := roleAttr.LastCollectTime.Add(time.Duration(interval)*time.Second)
		rspObj.NextTime = nextTime.UnixNano()/1e6
	}

	rspObj.CurTimes = roleAttr.CollectTimes
	rspObj.Limit = timeLimit
}

func (i *interiorController) transform(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.TransformReq{}
	rspObj := &model.TransformRsp{}

	mapstructure.Decode(req.Body.Msg, reqObj)
	rsp.Body.Msg = rspObj
	rsp.Body.Code = utils.OK

	r, _ := req.Conn.GetProperty("role")
	role := r.(*model.Role)
	roleRes, err := logic.RoleServer.GetRoleRes(role.RId)
	if err != nil {
		rsp.Body.Code = utils.DBError
		return
	}

	main, _ := logic.RoleCity.GetMainCity(role.RId)

	lv := logic.CityFacilityService.GetFacilityLv(main.CityId, model.JiShi)

	if lv >0 {
		rsp.Body.Code = utils.NotHasJiShi
		return
	}

	lens := 4
	ret := make([]int, lens)

	for i := 0; i < lens; i++ {
		if reqObj.From[i] > 0 {
			ret[i] = -reqObj.From[i]
		}

		if reqObj.To[i] > 0 {
			ret[i] = reqObj.To[i]
		}
	}

	if roleRes.Wood + ret[0] < 0 ||
		roleRes.Iron + ret[1] < 0 ||
		roleRes.Stone + ret[2] < 0 ||
		roleRes.Grain + ret[3] < 0 {
		rsp.Body.Code = utils.InvalidParam
		return
	}

	roleRes.Wood += ret[0]
	roleRes.Iron += ret[1]
	roleRes.Stone += ret[2]
	roleRes.Grain += ret[3]
	roleRes.SyncExecute()
}
