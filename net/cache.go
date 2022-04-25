package net

import (
	"sync"
)

var Mgr = &WsMgr{
	userCache: make(map[int]WSConn),
}

type WsMgr struct {
	WLock sync.RWMutex
	userCache map[int]WSConn
}

//用户登录缓存
func (w *WsMgr)UserLogin(conn WSConn, uid int, token string){
	w.WLock.Lock()
	defer w.WLock.Unlock()

	oldConn := w.userCache[uid]

	if oldConn != nil {
		//有用户登录了
		if conn != oldConn {
			//表示登录的连接跟已经存在的登录不一样，由用户抢登录
			oldConn.Push("robLogin",nil)
		}
	}

	w.userCache[uid] = conn
	w.userCache[uid].SetProperty("uid", uid)
	w.userCache[uid].SetProperty("token", token)
}

