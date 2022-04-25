package model

import "time"

/**
	User - 用户表
	LoginHistory - 登录历史表
	LoginLast - 记录最后一次登录
*/

const (
	Login  = iota
	LoginOut
)

type User struct {
	UId      int       `xorm:"uid pk autoincr"` //用户uid - 主键
	Username string    `xorm:"username" validate:"min=4,max=20,regexp=^[a-zA-Z0-9_]*$"` //用户名
	Passcode string    `xorm:"passcode"` 	//密码加密随机数
	Passwd   string    `xorm:"passwd"`		//密码
	Hardware string    `xorm:"hardware"`	//硬件信息
	Status   int       `xorm:"status"`		//用户账号状态。0-默认；1-冻结；2-停号
	Ctime    time.Time `xorm:"ctime"`		//登录时间
	Mtime    time.Time `xorm:"mtime"`
	IsOnline bool      `xorm:"-"`			//是否在线
}

type LoginHistory struct {
	Id       int       `xorm:"id pk autoincr"`
	UId      int       `xorm:"uid"`
	CTime    time.Time `xorm:"ctime"`
	Ip       string    `xorm:"ip"`
	State    int8      `xorm:"state"` 		//登录状态，0登录，1登出
	Hardware string    `xorm:"hardware"`
}

type LoginLast struct {
	Id         int       `xorm:"id pk autoincr"`
	UId        int       `xorm:"uid"` 			// 用户uid
	LoginTime  time.Time `xorm:"login_time"` 	//登录时间
	LogoutTime time.Time `xorm:"logout_time"`	//登出时间
	Ip         string    `xorm:"ip"`			//ip
	Session    string    `xorm:"session"`		//session 会话
	IsLogout   int8      `xorm:"is_logout"`		//是否logout,1:logout，0:login
	Hardware   string    `xorm:"hardware"`		//硬件
}
