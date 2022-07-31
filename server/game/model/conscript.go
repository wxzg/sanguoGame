package model

//征兵
type ConscriptReq struct {
	ArmyId		int  	`json:"armyId"`		//队伍id
	Cnts		[]int	`json:"cnts"`		//征兵人数 [20,20,0]
}

type ConscriptRsp struct {
	Army    Army    `json:"army"`
	RoleRes RoleRes `json:"role_res"`
}

