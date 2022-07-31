package controllor

import (
	"github.com/mitchellh/mapstructure"
	model2 "sanguoServer/db/model"
	"sanguoServer/net"
	"sanguoServer/server/common"
	"sanguoServer/server/game/logic"
	"sanguoServer/server/game/middleware"
	"sanguoServer/server/game/model"
	"sanguoServer/utils"
)

var CoalitionController = &coalitionController{
}
type coalitionController struct {

}

func (c *coalitionController) Router(router *net.Router) {
	g := router.Group("union")
	g.Use(middleware.Log())
	g.AddRouter("list",c.list,middleware.CheckRole())
	g.AddRouter("info",c.info,middleware.CheckRole())
	g.AddRouter("applyList",c.applyList,middleware.CheckRole())
}

func (c *coalitionController)list(req *net.WsMsgReq, rsp *net.WsMsgRsp){
	rspObj := &model.ListRsp{}

	rsp.Body.Name = req.Body.Name
	rsp.Body.Code = utils.OK
	uns, err := logic.CoalitionService.List()
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.List = uns
	rsp.Body.Msg = rspObj
}

func (c *coalitionController) info(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.InfoReq{}
	mapstructure.Decode(req.Body.Msg, reqObj)
	rspObj := &model.InfoRsp{}
	rsp.Body.Msg = rspObj
	rsp.Body.Code = utils.OK
	//通过联盟ID去获取信息
	un, err := logic.CoalitionService.Get(reqObj.Id)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.Info = un
	rspObj.Id = reqObj.Id
	rsp.Body.Msg = rspObj
}

func (c *coalitionController) applyList(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	// 根据联盟id去查询申请列表， rid是申请人，只有副盟主和盟主才能看到申请列表
	reqObj := &model.ApplyReq{}
	mapstructure.Decode(req.Body.Msg, reqObj)
	rspObj := &model.ApplyRsp{}
	rsp.Body.Code = utils.OK
	rsp.Body.Msg = rspObj

	r, _ := req.Conn.GetProperty("role")
	//拿到当前角色
	role := r.(*model2.Role)
	//查询联盟
	un := logic.CoalitionService.GetCoalition(reqObj.Id)
	if un == nil {
		rsp.Body.Code = utils.DBError
		return
	}
	if un.Chairman != role.RId && un.ViceChairman != role.RId {
		rspObj.Id = reqObj.Id
		rspObj.Applys = make([]model.ApplyItem,0)
		return
	}

	ais,err := logic.CoalitionService.GetListApply(reqObj.Id,0)
	if err != nil {
		rsp.Body.Code = utils.DBError
		return
	}
	rspObj.Id = reqObj.Id
	rspObj.Applys = ais
}