package model

import (
	"log"
	"sanguoServer/db"
)

func init(){
	go RoleAttrDao.running()
	go RoleResDao.running()
}

var RoleAttrDao = &roleAttrDao{
	raChan: make(chan *RoleAttribute,100),
}
type roleAttrDao struct {
	raChan chan *RoleAttribute
}

func (d roleAttrDao) Push(r *RoleAttribute) {
	d.raChan <- r
}

func (r *roleAttrDao) running() {
	for  {
		select {
		case ra := <- r.raChan:
			if ra.Id > 0 {
				_,err := db.Eg.Table(ra).ID(ra.Id).Cols(
					"parent_id", "collect_times", "last_collect_time", "pos_tags").Update(ra)
				if err != nil {
					log.Println("roleAttrDao update error",err)
				}
			}
		}
	}
}

func (r *RoleAttribute) SyncExecute() {
	RoleAttrDao.Push(r)
}


var RoleResDao = &roleResDao{
	resChan: make(chan *RoleRes,100),
}
type roleResDao struct {
	resChan  chan *RoleRes
}


func (r *roleResDao) running() {
	for  {
		select {
		case rr := <- r.resChan:
			if rr.Id > 0 {
				_,err := db.Eg.Table(rr).ID(rr.Id).Cols("wood", "iron", "stone",
					"grain", "gold", "decree").Update(rr)
				if err != nil {
					log.Println("roleResDao update error",err)
				}
			}
		}
	}
}

func (r *RoleRes) SyncExecute() {
	//通知修改
	RoleResDao.resChan <- r
}

