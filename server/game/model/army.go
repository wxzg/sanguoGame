package model

type Army struct {
	Id       int     `json:"id"`
	CityId   int     `json:"cityId"`
	UnionId  int     `json:"union_id"` //联盟id
	Order    int8    `json:"order"`    //第几队，1-5队
	Generals []int   `json:"generals"`
	Soldiers []int   `json:"soldiers"`
	ConTimes []int64 `json:"con_times"`
	ConCnts  []int   `json:"con_cnts"`
	Cmd      int8    `json:"cmd"`
	State    int8    `json:"state"` //状态:0:running,1:stop
	FromX    int     `json:"from_x"`
	FromY    int     `json:"from_y"`
	ToX      int     `json:"to_x"`
	ToY      int     `json:"to_y"`
	Start    int64   `json:"start"`//出征开始时间
	End      int64   `json:"end"`//出征结束时间
}

type ArmyOneRsp struct {
	Army	Army `json:"army"`
}

type ArmyOneReq struct {
	CityId	int  	`json:"cityId"`
	Order	int8	`json:"order"`
}

//派遣队伍
type AssignArmyReq struct {
	ArmyId int  `json:"armyId"` //队伍id
	Cmd    int8 `json:"cmd"`  //命令：0:空闲 1:攻击 2：驻军 3:返回
	X      int  `json:"x"`
	Y      int  `json:"y"`
}

type AssignArmyRsp struct {
	Army Army `json:"army"`
}