package logic

import (
	"log"
	"sanguoServer/db"
	"sanguoServer/db/model"
	"sanguoServer/server/common"
	"sanguoServer/server/web/webmodel"
	"sanguoServer/utils"
	"time"
)

var DefaultAccountLogic *AccountLogic = &AccountLogic{}

type AccountLogic struct {

}

func (l *AccountLogic) Register (rq *webmodel.RegisterReq) error{
	/**
		登录逻辑的业务处理
	*/
	//拿到用户名
	username := rq.Username
	user := &model.User{}
	ok,err := db.Eg.Table(user).Where("username=?",username).Get(user)
	if err != nil {
		log.Println("注册查询失败",err)
		return common.New(utils.DBError,"数据库异常")
	}
	if ok {
		//有数据 提示用户已存在
		return common.New(utils.UserExist,"用户已存在")
	}else {
		user.Mtime = time.Now()
		user.Ctime = time.Now()
		user.Username = rq.Username
		user.Passcode = utils.RandSeq(6)
		user.Passwd = utils.Password(rq.Password,user.Passcode)
		user.Hardware = rq.Hardware
		_,err := db.Eg.Table(user).Insert(user)
		if err != nil {
			log.Println("注册插入失败",err)
			return common.New(utils.DBError,"数据库异常")
		}
		return nil
	}
}
