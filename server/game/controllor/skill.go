package controllor

import (
	model2 "sanguoServer/db/model"
	"sanguoServer/net"
	"sanguoServer/server/common"
	"sanguoServer/server/game/logic"
	"sanguoServer/server/game/model"
	"sanguoServer/utils"
)

var SkillHandler = &skillHandler{}
type skillHandler struct {

}
func (sh *skillHandler) Router(r *net.Router)  {
	g := r.Group("skill")
	g.AddRouter("list",sh.list)
}

func (sh *skillHandler) list(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	rspObj := &model.SkillListRsp{}
	rsp.Body.Msg = rspObj
	rsp.Body.Code = utils.OK
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name

	role,_ := req.Conn.GetProperty("role")
	r := role.(*model2.Role)
	skills , err := logic.DefaultSkillService.GetSkills(r.RId)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.List = skills
	rsp.Body.Msg = rspObj
}