package controllor

import (
	"github.com/mitchellh/mapstructure"
	"log"
	"sanguoServer/db"
	model2 "sanguoServer/db/model"
	"sanguoServer/net"
	"sanguoServer/server/common"
	"sanguoServer/server/game"
	"sanguoServer/server/game/logic"
	"sanguoServer/server/game/model"
	"sanguoServer/utils"
	"time"
)

var DefaultRoleRouter = &roleRouter{}
type roleRouter struct{

}

func (r *roleRouter) Router (router *net.Router){
	g := router.Group("role")

	g.AddRouter("enterServer", r.enterRole)
	g.AddRouter("myProperty", r.myProperty)
}

func (r *roleRouter) enterRole (req *net.WsMsgReq, rsp *net.WsMsgRsp){
	//进入的游戏的逻辑
	//Session 需要验证是否合法 合法的前提下 可以取出角色id，如果没有就提醒无角色
	//根据角色id 查询角色拥有的资源roleRes ， 如果有资源就返回，没有就初始化资源
	reqObj := &model.EnterServerReq{}
	rspObj := &model.EnterServerRsp{}
	err := mapstructure.Decode(reqObj, req.Body.Msg)
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	if err != nil {
		rsp.Body.Code = utils.InvalidParam
		return
	}

	token := reqObj.Session
	_, claim, err := utils.ParseToken(token)
	if err != nil {
		rsp.Body.Code = utils.SessionInvalid
		return
	}
	//用户id
	uid := claim.Uid
	//根据用户id 查找角色
	role := &model2.Role{}
	ok, err := db.Eg.Table(role).Where("uid=?", uid).Get(role)
	if err != nil {
		log.Println("角色查询失败",err)
		rsp.Body.Code = utils.DBError
		return
	}

	if !ok {
		rsp.Body.Code = utils.RoleNotExist
		return
	}

	rid := role.RId
	rsp.Body.Code = utils.OK
	//将数据库类型的role转换为要还给前端的role类型
	rspObj.Role = role.ToModel().(model.Role)
	//查询资源
	roleRes := &model2.RoleRes{}
	ok, err = db.Eg.Table(roleRes).Where("rid=?",rid).Get(roleRes)
	if err != nil {
		rsp.Body.Code = utils.DBError
		return
	}
	//说明资源不存在， 加载初始资源
	if !ok {
		roleRes = &model2.RoleRes{
			RId: role.RId,
			Wood:          game.Base.Role.Wood,
			Iron:          game.Base.Role.Iron,
			Stone:         game.Base.Role.Stone,
			Grain:         game.Base.Role.Grain,
			Gold:          game.Base.Role.Gold,
			Decree:        game.Base.Role.Decree}
	}

	rspObj.RoleRes = roleRes.ToModel().(model.RoleRes)
	//转换成Token
	rspObj.Token, err = utils.Award(rid)
	if err != nil {
		log.Println("生成token出错", err)
		rsp.Body.Code = utils.SessionInvalid
	}
	rspObj.Time = time.Now().UnixNano()/1e6
	//把角色存一下
	req.Conn.SetProperty("role", role)
	//查询角色属性
	if err := logic.RoleAttrServer.RoleAttributeHandler(rid); err != nil {
		rsp.Body.Code = utils.DBError
		rsp.Body.Msg = "角色属性出错"
		return
	}
	//玩家城池
	if err := logic.RoleCity.RoleCityServer(rid,role.NickName); err != nil {
		rsp.Body.Code = utils.DBError
		rsp.Body.Msg = "城池出错"
		return
	}

	rsp.Body.Msg = rspObj
	rsp.Body.Code = utils.OK
}

func (r *roleRouter) myProperty(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.MyRolePropertyReq{}
	rspObj := &model.MyRolePropertyRsp{}
	mapstructure.Decode(req.Body.Msg, reqObj)
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	//根据角色id去查询军队、资源、建筑、城池、武将
	role, err := req.Conn.GetProperty("role")
	if err != nil {
		log.Println("session invalid",err)
		rsp.Body.Code = utils.SessionInvalid
		return
	}

	rid := role.(*model.Role).RId
	//查资源
	rspObj.RoleRes,err = logic.RoleServer.GetRoleRes(rid)
	if err != nil {
		rsp.Body.Code = utils.DBError
		return
	}

	//城池
	rspObj.Citys,err = logic.RoleCity.GetCitys(rid)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}

	//建筑
	rspObj.MRBuilds,err = logic.DefaultRoleBuildService.GetBuilds(rid)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}

	//军队
	rspObj.Armys,err = logic.DefaultArmyService.GetArmys(rid)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}

	//武将
	rspObj.Generals,err = logic.DefaultGeneralService.GetGenerals(rid)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}

	rsp.Body.Msg = rspObj
}
