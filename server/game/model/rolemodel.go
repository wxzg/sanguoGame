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

type PosTagListRsp struct {
	PosTags	[]PosTag	`json:"pos_tags"`
}

//坐标
type PosTag struct {
	X		int	`json:"x"`
	Y		int	`json:"y"`
	Name string `json:"name"`
}

//武将
type MyGeneralRsp struct {
	Generals []General `json:"generals"`
}

//我的军队
type ArmyListReq struct {
	CityId	int  `json:"cityId"`
}

type ArmyListRsp struct {
	CityId	int   `json:"cityId"`
	Armys	[]Army `json:"armys"`
}

//我的战报
type WarReportReq struct {

}

type WarReportRsp struct {
	List	[]WarReport `json:"list"`
}

type WarReport struct {
	Id             		int    		`json:"id"`
	AttackRid      		int    		`json:"a_rid"`
	DefenseRid     		int    		`json:"d_rid"`
	BegAttackArmy  		string 		`json:"b_a_army"`
	BegDefenseArmy 		string 		`json:"b_d_army"`
	EndAttackArmy  		string 		`json:"e_a_army"`
	EndDefenseArmy 		string 		`json:"e_d_army"`
	BegAttackGeneral  	string    	`json:"b_a_general"`
	BegDefenseGeneral 	string    	`json:"b_d_general"`
	EndAttackGeneral  	string    	`json:"e_a_general"`
	EndDefenseGeneral 	string    	`json:"e_d_general"`
	Result				int      	`json:"result"`	//0失败，1打平，2胜利
	Rounds				string		`json:"rounds"` //回合
	AttackIsRead   		bool   		`json:"a_is_read"`
	DefenseIsRead  		bool   		`json:"d_is_read"`
	DestroyDurable 		int    		`json:"destroy"`
	Occupy         		int    		`json:"occupy"`
	X              		int    		`json:"x"`
	Y              		int    		`json:"y"`
	CTime          		int  		`json:"ctime"`
}

//技能
type SkillListReq struct {

}

type Skill struct {
	Id             	int 	`json:"id"`
	CfgId          	int 	`json:"cfgId"`
	Generals 		[]int 	`json:"generals"`
}


type SkillListRsp struct {
	List []Skill `json:"list"`
}

//加载角色建筑
type ScanBlockReq struct {
	X 		int    		`json:"x"`
	Y 		int    		`json:"y"`
	Length	int			`json:"length"`
}

type ScanRsp struct {
	MRBuilds []MapRoleBuild `json:"mr_builds"` //角色建筑，包含被占领的基础建筑
	MCBuilds []MapRoleCity  `json:"mc_builds"` //角色城市
	Armys 	[]Army         	`json:"armys"`     //军队
}


//创建角色
type CreateRoleReq struct {
	UId			int		`json:"uid"`
	NickName 	string	`json:"nickName"`
	Sex			int8	`json:"sex"`
	SId			int		`json:"sid"`
	HeadId		int16	`json:"headId"`
}

type CreateRoleRsp struct {
	Role Role `json:"role"`
}

//征收查询
type OpenCollectionReq struct {
}

type OpenCollectionRsp struct {
	Limit    int8  `json:"limit"`
	CurTimes int8  `json:"cur_times"`
	NextTime int64 `json:"next_time"`
}

//征收
type CollectionReq struct {
}

type CollectionRsp struct {
	Gold		int		`json:"gold"`
	Limit    	int8  	`json:"limit"`
	CurTimes 	int8  	`json:"cur_times"`
	NextTime 	int64 `json:"next_time"`
}

type Yield struct {
	Wood  int
	Iron  int
	Stone int
	Grain int
	Gold  int
}

var GetYield func(rid int) Yield
var  GetUnion func(rid int) int

