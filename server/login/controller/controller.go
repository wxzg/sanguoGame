package controller

import (
	"github.com/mitchellh/mapstructure"
	"log"
	"sanguoServer/db"
	"sanguoServer/db/model"
	"sanguoServer/net"
	"sanguoServer/server/login/proto"
	"sanguoServer/utils"
	"time"
)

//新建Account组
var DefaultAccount = &Account{}
type Account struct {

}
//绑定路由
func (a *Account) Router(r *net.Router)  {
	g := r.Group("account")
	g.AddRouter("login",a.login)
}

//登录
func (a *Account) login(req *net.WsMsgReq, rsp *net.WsMsgRsp)  {
	/**
		登录逻辑
		1.前端传来用户名、密码、硬件id
		2.根据用户名查询数据库user表
		3。进行密码比对
		4.保存用户登录记录
		5.保存用户最后一次登录信息
		6.客户端需要一个session， jwt生成一个加密字符串相关的加密算法
		7.再发起用户登录行为时，判断用户是否合法
	*/

	//这两个是ReqBody和ResBody里面msg是数据部分
	loginRes := &proto.LoginRsp{}
	loginReq := &proto.LoginReq{}

	//将map转struct，因为我们不知道JSON是什么类型的，所以接受的时候是map[sring]interface{}类型，现在需要转为struct
	mapstructure.Decode(req.Body.Msg, loginReq)

	//新建一条用户表数据
	user := &model.User{}
	//根据客户端传来的username来查询
	ok,err := db.Eg.Table(user).Where("username=?", loginReq.Username).Get(user)
	if err != nil{
		log.Println("用户表获取失败", err)
		panic(err)
	}

	if !ok {
		log.Println("用户不存在")
		//返回用户不存在Code
		rsp.Body.Code = utils.UserNotExist
		return
	}
	//将客户端发来的密码根据passcode加密以后跟数据库里的密码对比
	pwd := utils.Password(loginReq.Password, user.Passcode)
	if pwd != user.Passwd {
		log.Println("密码错误")
		rsp.Body.Code = utils.PwdIncorrect
		return
	}
	//jwt A.B.C 三部分 A定义加密算法 B定义放入的数据 C部分 根据秘钥+A和B生成加密字符串
	token,_ := utils.Award(user.UId)
	rsp.Body.Code = utils.OK
	loginRes.UId = user.UId
	loginRes.Username = user.Username
	loginRes.Session = token //将生成的token发回去，来表示这次一登录的凭证
	loginRes.Password = ""
	rsp.Body.Msg = loginRes

	//记录到登录历史
	ul := &model.LoginHistory{
		UId: user.UId, //用户uid
		CTime: time.Now(), //登录时间
		Ip: loginReq.Ip,	//登录ip
		Hardware: loginReq.Hardware, //硬件
		State: model.Login,	//登录状态
	}
	//插入
	db.Eg.Table(ul).Insert(ul)

	//记录最后一次登录
	l1 := &model.LoginLast{}
	//先查一下该用户上一次登录记录
	ok, err = db.Eg.Table(l1).Where("uid=?", user.UId).Get(l1)
	if err != nil {
		log.Println("loginLast表获取失败", err)
		panic(err)
	}

	l1.Ip = loginReq.Ip //登录ip
	l1.IsLogout = 0
	l1.Session = token //登录凭证token
	l1.Hardware = loginReq.Hardware
	l1.LoginTime = ul.CTime //最后一次登录时间

	if ok { //有就更新，没有就插入
		db.Eg.Table(l1).Update(l1)
	} else {
		l1.UId = user.UId
		db.Eg.Table(l1).Insert(l1)
	}

	//缓存一下 此用户和当前的ws连接
	net.Mgr.UserLogin(req.Conn,user.UId,token)
}
