package model

import (
	"sanguoServer/server/game/gameconfig"
	"time"
)

type Facility struct {
	Name         string `json:"name"`
	Level int8   `json:"level"` 		//等级，外部读的时候不能直接读，要用GetLevel
	Type         int8   `json:"type"`
	UpTime       int64  `json:"up_time"`	//升级的时间戳，0表示该等级已经升级完成了
}

type FacilitiesReq struct {
	CityId		int		`json:"cityId"`
}

type FacilitiesRsp struct {
	CityId		int        `json:"cityId"`
	Facilities	[]Facility `json:"facilities"`
}

//升级设施
type UpFacilityReq struct {
	CityId int 	`json:"cityId"` //升级的建筑id
	FType  int8	`json:"fType"`	//建筑类型
}

type UpFacilityRsp struct {
	CityId   int      `json:"cityId"`
	Facility Facility `json:"facility"`
	RoleRes  RoleRes  `json:"role_res"`
}


func (f *Facility) GetLevel() int8 {
	if f.UpTime > 0 {
		cur := time.Now().Unix()
		cost := gameconfig.FcailityConf.CostTime(f.Type, f.Level+1)
		if cur >= f.UpTime + int64(cost) {
			f.Level += 1
			f.UpTime = 0
		}
	}
	return f.Level
}

func (f *Facility) CanLv() bool {
	f.GetLevel()
	return f.UpTime == 0
}

