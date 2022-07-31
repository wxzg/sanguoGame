package logic

import (
	"encoding/json"
	"errors"
	"log"
	"sanguoServer/db"
	model2 "sanguoServer/db/model"
	"sanguoServer/server/common"
	"sanguoServer/server/game/model"
	"sanguoServer/utils"
	"sync"
	"xorm.io/xorm"
)

type roleAttrServer struct{
	mutex sync.RWMutex
	attrs map[int]*model2.RoleAttribute
}

var RoleAttrServer = &roleAttrServer{
	attrs: make(map[int]*model2.RoleAttribute,0),
}

func (r *roleAttrServer) RoleAttributeHandler (rid int, session *xorm.Session) error{
	rr := &model2.RoleAttribute{}
	ok, err := db.Eg.Table(rr).Where("rid=?", rid).Get(rr)
	if err != nil {
		log.Println("玩家属性查询出错：",err)
		return err
	}

	if !ok {
		rr.RId = rid
		rr.ParentId = 0
		rr.UnionId = 0
		rr.PosTags = ""

		if session != nil {
			_, err = session.Table(rr).Insert(rr)
		} else {
			_, err = db.Eg.Table(rr).Insert(rr)
		}

		if err != nil {
			log.Println("玩家属性插入出错",err)
			return err
		}
	}

	r.mutex.Lock()
	r.attrs[rid] = rr
	r.mutex.Unlock()

	return nil
}

func (r *roleAttrServer) GetPosTags(rid int) ([]model.PosTag,error) {
	var err error
	ra, ok := r.attrs[rid]
	if !ok {
		ra = &model2.RoleAttribute{}
		//查询posTag
		ok ,err = db.Eg.Table(ra).Where("rid=?",rid).Get(ra)

		if err != nil {
			log.Println("GetTagList", err)
			return nil, common.New(utils.DBError, "posTag查询出错")
		}

		r.mutex.Lock()
		r.attrs[rid] = ra
		r.mutex.Unlock()
	}


	posTags := make([]model.PosTag,0)
	if ok {
		//如果成功的话取出来的是JSON格式的数据需要转一下
		if ra.PosTags != ""{
			err = json.Unmarshal([]byte(ra.PosTags),posTags)
		}
	}

	return posTags, nil
}

func (r *roleAttrServer) Load(){
	//加载
	t := make(map[int]*model2.RoleAttribute)
	err := db.Eg.Find(t)
	if err != nil {
		log.Println(err)
	}

	//获取联盟id
	for _, v := range t {
		r.attrs[v.RId] = v
	}

	l, _ := CoalitionService.ListCoalition()

	for _,c := range l {
		for _, rid := range c.MemberArray {
			attr, ok := r.attrs[rid]
			if ok {
				attr.UnionId = c.Id
			}
		}
	}
}

func (r *roleAttrServer) Get(rid int) (*model2.RoleAttribute, error){
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	ra, ok := r.attrs[rid]
	if ok {
		return ra,nil
	}else {
		ra = &model2.RoleAttribute{}
		ok , err :=  db.Eg.Table(ra).Where("rid=?",rid).Get(ra)
		if err != nil {
			return nil, errors.New("角色属性获取失败")
		}

		if ok {
			return ra, nil
		}
	}
	return nil, errors.New("角色属性获取失败")
}

func (r *roleAttrServer) GetUnion(rid int) int {
	return r.attrs[rid].UnionId
}

