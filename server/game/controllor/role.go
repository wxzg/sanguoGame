package controllor

import (
	"github.com/mitchellh/mapstructure"
	"log"
	"sanguoServer/db"
	model2 "sanguoServer/db/model"
	"sanguoServer/net"
	"sanguoServer/server/common"
	"sanguoServer/server/game/gameconfig"
	"sanguoServer/server/game/logic"
	"sanguoServer/server/game/middleware"
	"sanguoServer/server/game/model"
	"sanguoServer/utils"
	"time"
)

var DefaultRoleRouter = &roleRouter{}
type roleRouter struct{

}

func (r *roleRouter) Router (router *net.Router){
	g := router.Group("role")
	g.Use(middleware.Log())
	g.AddRouter("enterServer", r.enterRole)
	g.AddRouter("myProperty", r.myProperty, middleware.CheckRole())
	g.AddRouter("posTagList", r.posTagList)
	g.AddRouter("create", r.create)
}

func (r *roleRouter) enterRole (req *net.WsMsgReq, rsp *net.WsMsgRsp){
	//进入的游戏的逻辑
	//Session 需要验证是否合法 合法的前提下 可以取出角色id，如果没有就提醒无角色
	//根据角色id 查询角色拥有的资源roleRes ， 如果有资源就返回，没有就初始化资源
	reqObj := &model.EnterServerReq{}
	rspObj := &model.EnterServerRsp{}
	err := mapstructure.Decode(req.Body.Msg, reqObj)
	if err != nil {
		log.Println(err)
		rsp.Body.Code = utils.InvalidParam
		return
	}
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name

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
	//新建事务
	session := db.Eg.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		log.Println("事务开启失败")
		rsp.Body.Code = utils.DBError
		return
	}
	//保存一下事务
	req.Context.Set("session", session)
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
			Wood:          gameconfig.Base.Role.Wood,
			Iron:          gameconfig.Base.Role.Iron,
			Stone:         gameconfig.Base.Role.Stone,
			Grain:         gameconfig.Base.Role.Grain,
			Gold:          gameconfig.Base.Role.Gold,
			Decree:        gameconfig.Base.Role.Decree}
		_, err = session.Table(roleRes).Insert(roleRes)
		if err != nil {
			log.Println("角色资源插入失败",err)
			rsp.Body.Code = utils.DBError
			return
		}

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
	if err := logic.RoleAttrServer.RoleAttributeHandler(rid,session); err != nil {
		session.Rollback()
		rsp.Body.Code = utils.DBError
		rsp.Body.Msg = "角色属性出错"
		return
	}
	//玩家城池
	if err := logic.RoleCity.RoleCityServer(rid,role.NickName, session); err != nil {
		session.Rollback()
		rsp.Body.Code = utils.DBError
		rsp.Body.Msg = "城池出错"
		return
	}
	err = session.Commit()
	if err != nil {
		log.Println("事务提交失败")
		rsp.Body.Code = utils.DBError
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

	rid := role.(*model2.Role).RId
	//查资源
	mrs,err := logic.RoleServer.GetRoleRes(rid)
	if err != nil {
		rsp.Body.Code = utils.DBError
		return
	}
	rspObj.RoleRes = mrs.ToModel().(model.RoleRes)
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

func (r *roleRouter) posTagList(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	rspObj := &model.PosTagListRsp{}
	role, err := req.Conn.GetProperty("role")
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	if err != nil {
		rsp.Body.Code =utils.InvalidParam
		return
	}
	rid := role.(*model2.Role).RId

	rspObj.PosTags, err = logic.RoleAttrServer.GetPosTags(rid)
	if err != nil {
		log.Println(err)
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rsp.Body.Code = utils.OK
	rsp.Body.Msg = rspObj
}

//创建角色
func (r *roleRouter) create(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.CreateRoleReq{}
	rspObj := &model.CreateRoleRsp{}
	mapstructure.Decode(req.Body.Msg,reqObj)

	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	role := &model2.Role{}
	ok,err := db.Eg.Where("uid=?",reqObj.UId).Get(role)
	if err!=nil  {
		rsp.Body.Code = utils.DBError
		return
	}
	if ok {
		rsp.Body.Code = utils.RoleAlreadyCreate
		return
	}
	role.UId = reqObj.UId
	role.Sex = reqObj.Sex
	role.NickName = reqObj.NickName
	role.Balance = 0
	role.HeadId = reqObj.HeadId
	role.CreatedAt = time.Now()
	role.LoginTime = time.Now()
	_,err = db.Eg.InsertOne(role)
	if err!=nil  {
		rsp.Body.Code = utils.DBError
		return
	}
	rspObj.Role = role.ToModel().(model.Role)
	rsp.Body.Code = utils.OK
	rsp.Body.Msg = rspObj
}



