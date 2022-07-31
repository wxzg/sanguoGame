package controllor

import (
	model2 "sanguoServer/db/model"
	"sanguoServer/net"
	"sanguoServer/server/common"
	"sanguoServer/server/game/logic"
	"sanguoServer/server/game/model"
	"sanguoServer/utils"
)

var WarHandler = &warHandler{}
type warHandler struct {

}
func (w *warHandler) InitRouter(r *net.Router)  {
	g := r.Group("war")
	g.AddRouter("report",w.report)
}

func (w *warHandler) report(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	rspObj := &model.WarReportRsp{}
	rsp.Body.Msg = rspObj
	rsp.Body.Code = utils.OK
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name

	role,_ := req.Conn.GetProperty("role")
	r := role.(*model2.Role)

	wReports, err := logic.DefaultWarService.GetWarReports(r.RId)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.List = wReports
	rsp.Body.Msg = rspObj
}
