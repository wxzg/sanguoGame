package logic

import (
	"log"
	"sanguoServer/db"
	model2 "sanguoServer/db/model"
	"sanguoServer/server/common"
	"sanguoServer/server/game/model"
	"sanguoServer/utils"
	"sync"
)

type coalitionService struct{
	mutex  sync.RWMutex
	unions map[int]*model2.Coalition
}

var CoalitionService = &coalitionService{
	unions: make(map[int]*model2.Coalition),
}

func (c *coalitionService) Load()  {
	rr := make([]*model2.Coalition, 0)
	err := db.Eg.Table(new(model2.Coalition)).Where("state=?",model2.UnionRunning).Find(&rr)

	if err != nil {
		log.Println("coalitionService load err", err)
	}

	for _, v := range rr{
		c.unions[v.Id] = v
	}
}

func (c *coalitionService) List() ([]model.Union, error) {
	r := make([]model.Union, 0)
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	for _, coalition := range c.unions {
		union := coalition.ToModel()
		//盟主和副盟主的信息
		mainInfo := make([]model.Major, 0)
		//查询盟主和副盟主的信息存进去
		if role := RoleServer.Get(coalition.Chairman); role != nil {
			m := model.Major{
				Name: role.NickName,
				RId: role.RId,
				Title: model.UnionChairman,
			}
			mainInfo = append(mainInfo, m)
		}

		if role := RoleServer.Get(coalition.ViceChairman); role != nil {
			m := model.Major{
				Name: role.NickName,
				RId: role.RId,
				Title: model.UnionChairman,
			}
			mainInfo = append(mainInfo, m)
		}
		union.Major = mainInfo
		r = append(r, union)
	}

	return r, nil
}

func (c *coalitionService) ListCoalition() ([]*model2.Coalition, error) {
	r := make([]*model2.Coalition, 0)
	for _, coalition := range c.unions {
		r = append(r, coalition)
	}

	return r, nil
}

func (c *coalitionService) Get(id int) (model.Union, error) {
	c.mutex.RUnlock()
	defer c.mutex.RUnlock()
	coalition, ok := c.unions[id]
	if ok {
		union := coalition.ToModel()
		//盟主和副盟主信息
		main := make([]model.Major, 0)
		if role := RoleServer.Get(coalition.Chairman);role != nil{
			m := model.Major{Name: role.NickName, RId: role.RId, Title: model.UnionChairman}
			main = append(main, m)
		}
		if role := RoleServer.Get(coalition.ViceChairman);role != nil{
			m := model.Major{Name: role.NickName, RId: role.RId, Title: model.UnionChairman}
			main = append(main, m)
		}
		union.Major = main
		return union,nil
	}
	return model.Union{}, nil
}

func (c *coalitionService) GetCoalition(id int) *model2.Coalition {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	coal, ok := c.unions[id]
	if ok {
		return coal
	}

	return nil
}

func (c *coalitionService) GetListApply(unionId int, state int) ([]model.ApplyItem, error) {
	applys := make([]model2.CoalitionApply,0)
	err := db.Eg.Table(new(model2.CoalitionApply)).
		Where("union_id=? and state=?",unionId,state).
		Find(&applys)
	if err != nil {
		log.Println("coalitionService GetListApply find error",err)
		return nil,common.New(utils.DBError,"数据库错误")
	}
	ais := make([]model.ApplyItem,0)
	for _,v := range applys{
		var ai model.ApplyItem
		ai.Id = v.Id
		role := RoleServer.Get(v.RId)
		ai.NickName = role.NickName
		ai.RId = role.RId
		ais = append(ais,ai)
	}
	return ais,nil
}