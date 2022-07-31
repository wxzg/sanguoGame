package model

import "sanguoServer/db"

var GeneralDao = &generalDao{
	aChan: make(chan *General,100),
}
type generalDao struct {
	aChan chan *General
}

func (g *generalDao) running() {
	for  {
		select {
		case gen := <- g.aChan:
			if gen.Id > 0 && gen.RId > 0 {
				db.Eg.Table(gen).ID(gen.Id).Cols(
					"level", "exp", "order", "cityId",
					"physical_power", "star_lv", "has_pr_point",
					"use_pr_point", "force_added", "strategy_added",
					"defense_added", "speed_added", "destroy_added",
					"parentId", "compose_type", "skills", "state").Update(g)
			}
		}
	}
}

func init()  {
	go GeneralDao.running()
}

func (g *General) SyncExecute(){
	GeneralDao.aChan <- g
}