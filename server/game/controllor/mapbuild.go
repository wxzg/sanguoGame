package controllor

import (
	"github.com/mitchellh/mapstructure"
	"sanguoServer/net"
	"sanguoServer/server/game/gameconfig"
	"sanguoServer/utils"
)

type mapBuildHandler struct {

}

func (h mapBuildHandler) Router(router *net.Router) {
	n := router.Group("nationMap")
	n.AddRouter("config", h.configHandler)
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

var DeafultMapBuildHandler = &mapBuildHandler{}

