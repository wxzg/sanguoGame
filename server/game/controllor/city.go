package controllor

import (
	"github.com/mitchellh/mapstructure"
	"sanguoServer/net"
	"sanguoServer/server/game/logic"
	"sanguoServer/server/game/middleware"
	"sanguoServer/server/game/model"
	"sanguoServer/utils"
)

var CityController = &cityController{}
type cityController struct {

}

func (c *cityController) Router(router *net.Router)  {
	g := router.Group("city")
	g.Use(middleware.Log())
	g.AddRouter("facilities",c.facilities,middleware.CheckRole())
	g.AddRouter("upFacility",c.upFacility,middleware.CheckRole())
}

func (c *cityController) facilities(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.FacilitiesReq{}
	rspObj := &model.FacilitiesRsp{}
	mapstructure.Decode(req.Body.Msg, reqObj)
	rsp.Body.Msg = rspObj
	rspObj.CityId = reqObj.CityId
	rsp.Body.Code = utils.OK

	r, _ := req.Conn.GetProperty("role")
	city, ok := logic.RoleCity.Get(reqObj.CityId)
	if ok == false {
		rsp.Body.Code = utils.CityNotExist
		return
	}

	role := r.(*model.Role)
	if city.RId != role.RId {
		rsp.Body.Code = utils.CityNotMe
		return
	}

	f, ok := logic.CityFacilityService.Get(reqObj.CityId)
	if ok == false {
		rsp.Body.Code = utils.CityNotExist
		return
	}

	t := f.Facility()

	rspObj.Facilities = make([]model.Facility, len(t))
	for i, v := range t {
		rspObj.Facilities[i].Name = v.Name
		rspObj.Facilities[i].Level = v.GetLevel()
		rspObj.Facilities[i].Type = v.Type
		rspObj.Facilities[i].UpTime = v.UpTime
	}
}

func (c *cityController) upFacility(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.UpFacilityReq{}
	rspObj := &model.UpFacilityRsp{}
	mapstructure.Decode(req.Body.Msg, reqObj)
	rsp.Body.Msg = rspObj
	rspObj.CityId = reqObj.CityId
	rsp.Body.Code = utils.OK

	r, _ := req.Conn.GetProperty("role")
	city, ok := logic.RoleCity.Get(reqObj.CityId)
	if !ok {
		rsp.Body.Code = utils.CityNotExist
		return
	}

	role := r.(*model.Role)

	if city.RId != role.RId{
		rsp.Body.Code = utils.CityNotMe
		return
	}

	_, ok = logic.CityFacilityService.Get(reqObj.CityId)
	if !ok {
		rsp.Body.Code = utils.CityNotExist
		return
	}

	out, errcode := logic.CityFacilityService.UpFacility(role.RId, reqObj.CityId, reqObj.FType)
	rsp.Body.Code = errcode
	if errcode == utils.OK {
		rspObj.Facility.Level = out.GetLevel()
		rspObj.Facility.Type = out.Type
		rspObj.Facility.Name = out.Name
		rspObj.Facility.UpTime = out.UpTime
		roleRes := logic.RoleServer.Get(role.RId)
		rspObj.RoleRes = roleRes.ToModel().(model.RoleRes)
	}
}


