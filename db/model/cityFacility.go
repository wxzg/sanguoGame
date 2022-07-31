package model

import (
	"encoding/json"
	"log"
	"sanguoServer/db"
	"sanguoServer/server/game/model"
)

//城市设施
type CityFacility struct {
	Id         int    `xorm:"id pk autoincr"`
	RId        int    `xorm:"rid"`
	CityId     int    `xorm:"cityId"`
	Facilities string `xorm:"facilities"`
}

func (f *CityFacility) Facility() []model.Facility{
	facilities := make([]model.Facility, 0)
	json.Unmarshal([]byte(f.Facilities), &facilities)
	return facilities
}

func (f *CityFacility) SyncExecute() {
	CityFacilityDao.cfChan <- f
}

var CityFacilityDao = &cityFacilityDao{
	cfChan: make(chan *CityFacility),
}
type cityFacilityDao struct {
	cfChan chan *CityFacility
}

func (cf *cityFacilityDao) running() {
	for true {
		select {
		case c := <- cf.cfChan:
			if c.Id >0 {
				_, err := db.Eg.Table(c).ID(c.Id).Cols("facilities").Update(c)
				if err != nil{
					log.Println("db error", err)
				}
			}else{
				log.Println("update CityFacility fail, because id <= 0")
			}
		}
	}
}

func init()  {
	go CityFacilityDao.running()
}
