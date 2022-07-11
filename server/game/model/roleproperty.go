package model

import model2 "sanguoServer/db/model"

type MyRolePropertyReq struct {

}

type MyRolePropertyRsp struct {
	RoleRes  RoleRes        `json:"role_res"`
	MRBuilds []MapRoleBuild `json:"mr_builds"` //角色建筑，包含被占领的基础建筑
	Generals []General      `json:"generals"`
	Citys    []MapRoleCity  `json:"citys"`
	Armys    []Army         `json:"armys"`
}

type MapRoleBuild struct {
	RId        	int    	`json:"rid"`
	RNick      	string 	`json:"RNick"` 		//角色昵称
	Name       	string 	`json:"name"`
	UnionId    	int    	`json:"union_id"`   //联盟id
	UnionName  	string 	`json:"union_name"` //联盟名字
	ParentId   	int    	`json:"parent_id"`  //上级id
	X          	int    	`json:"x"`
	Y          	int    	`json:"y"`
	Type       	int8   	`json:"type"`
	Level      	int8   	`json:"level"`
	OPLevel     int8   	`json:"op_level"`
	CurDurable 	int    	`json:"cur_durable"`
	MaxDurable 	int    	`json:"max_durable"`
	Defender   	int    	`json:"defender"`
	OccupyTime	int64 	`json:"occupy_time"`
	EndTime 	int64 	`json:"end_time"`		//建造完的时间
	GiveUpTime	int64 	`json:"giveUp_time"`	//领地到了这个时间会被放弃
}

type General struct {
	Id        		int     	`json:"id"`
	CfgId     		int			`json:"cfgId"`
	PhysicalPower 	int     	`json:"physical_power"`
	Order     		int8    	`json:"order"`
	Level			int8    	`json:"level"`
	Exp				int			`json:"exp"`
	CityId    		int     	`json:"cityId"`
	CurArms         int     	`json:"curArms"`
	HasPrPoint      int     	`json:"hasPrPoint"`
	UsePrPoint      int     	`json:"usePrPoint"`
	AttackDis       int     	`json:"attack_distance"`
	ForceAdded      int     	`json:"force_added"`
	StrategyAdded   int     	`json:"strategy_added"`
	DefenseAdded    int     	`json:"defense_added"`
	SpeedAdded      int     	`json:"speed_added"`
	DestroyAdded    int     	`json:"destroy_added"`
	StarLv          int8    	`json:"star_lv"`
	Star            int8    	`json:"star"`
	ParentId        int     	`json:"parentId"`
	Skills			[]*model2.GSkill	`json:"skills"`
	State     		int8    	`json:"state"`

}

type MapRoleCity struct {
	CityId     	int    	`json:"cityId"`
	RId        	int    	`json:"rid"`
	Name       	string 	`json:"name"`
	UnionId    	int    	`json:"union_id"` 	//联盟id
	UnionName  	string 	`json:"union_name"`	//联盟名字
	ParentId   	int    	`json:"parent_id"`	//上级id
	X          	int    	`json:"x"`
	Y          	int    	`json:"y"`
	IsMain     	bool   	`json:"is_main"`
	Level      	int8   	`json:"level"`
	CurDurable 	int    	`json:"cur_durable"`
	MaxDurable 	int    	`json:"max_durable"`
	OccupyTime	int64 	`json:"occupy_time"`
}

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