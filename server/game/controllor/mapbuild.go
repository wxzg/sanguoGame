package controllor

import (
	"github.com/mitchellh/mapstructure"
	model2 "sanguoServer/db/model"
	"sanguoServer/net"
	"sanguoServer/server/game/gameconfig"
	"sanguoServer/server/game/logic"
	"sanguoServer/server/game/model"
	"sanguoServer/utils"
)

type mapBuildHandler struct {

}

var DeafultMapBuildHandler = &mapBuildHandler{}

func (h mapBuildHandler) Router(router *net.Router) {
	n := router.Group("nationMap")
	n.AddRouter("config", h.configHandler)
	n.AddRouter("scanBlock", h.scanBlock)
}

func (h mapBuildHandler) configHandler(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqobj := &gameconfig.ConfigReq{}
	rspobj := &gameconfig.ConfigRsp{}

	_ = mapstructure.Decode(req.Body.Msg, reqobj)

	rsp.Body.Msg = rspobj
	rsp.Body.Name = req.Body.Name
	rsp.Body.Code = utils.OK

	m := gameconfig.MapBuilder.Cfg
	rspobj.Confs = make([]gameconfig.Conf, len(m))
	for index, v := range m{
		rspobj.Confs[index].Type = v.Type
		rspobj.Confs[index].Name = v.Name
		rspobj.Confs[index].Level = v.Level
		rspobj.Confs[index].Defender = v.Defender
		rspobj.Confs[index].Durable = v.Durable
		rspobj.Confs[index].Grain = v.Grain
		rspobj.Confs[index].Iron = v.Iron
		rspobj.Confs[index].Stone = v.Stone
		rspobj.Confs[index].Wood = v.Wood
	}
}

func (n *mapBuildHandler) scanBlock(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.ScanBlockReq{}
	rspObj := &model.ScanRsp{}
	mapstructure.Decode(req.Body.Msg, reqObj)
	rsp.Body.Msg = rspObj
	rsp.Body.Code = utils.OK
	//取出角色信息
	r, _ := req.Conn.GetProperty("role")
	role := r.(*model2.Role)
	//扫描建筑
	rm := logic.DefaultRoleBuildService.ScanBlock(reqObj)
	rspObj.MRBuilds = rm
	//扫描城池
	rc := logic.RoleCity.ScanBlock(reqObj)
	rspObj.MCBuilds = rc
	//扫描军队
	armys := logic.DefaultArmyService.ScanBlock(role.RId,reqObj)
	rspObj.Armys = armys

}

