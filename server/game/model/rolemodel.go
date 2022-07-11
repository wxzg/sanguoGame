package model

type EnterServerReq struct{
	Session string `json:"session"`
}

type Role struct {
	RId 		int 	`json:"rid"`
	UId 		int		`json:"uid"`
	NickName 	string 	`json:"nickName"`
	Sex 		int8 	`json:"sex"`
	Balance 	int 	`json:"balance"`
	HeadId 		int16 	`json:"headId"`
	Profile 	string 	`json:"profile"`
}

// 角色拥有的资源
type RoleRes struct {
	Wood			int			`json:"wood"`
	Iron			int			`json:"iron"`
	Stone			int			`json:"stone"`
	Grain			int			`json:"grain"`
	Gold			int			`json:"gold"`
	Decree			int			`json:"decree"`	//令牌
	WoodYield		int			`json:"wood_yield"`
	IronYield		int			`json:"iron_yield"`
	StoneYield		int			`json:"stone_yield"`
	GrainYield		int			`json:"grain_yield"`
	GoldYield		int			`json:"gold_yield"`
	DepotCapacity	int			`json:"depot_capacity"`	//仓库容量
}

type EnterServerRsp struct{
	Role 		Role 		`json:"role"`
	RoleRes 	RoleRes 	`json:"role_res"`
	Time  		int64 		`json:"time"`
	Token		string		`json:"token"`
}
