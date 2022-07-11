package model

import (
	"sanguoServer/server/game/model"
	"sync"
	"time"
)

/**
	User - 用户表
	LoginHistory - 登录历史表
	LoginLast - 记录最后一次登录
	role - 角色表
	role_res - 角色资源表
	MapRoleCity - 玩家城池
	RoleAttribute - 玩家属性
*/

const (
	Login  = iota
	LoginOut
)

type User struct {
	UId      int       `xorm:"uid pk autoincr"` //用户uid - 主键
	Username string    `xorm:"username" validate:"min=4,max=20,regexp=^[a-zA-Z0-9_]*$"` //用户名
	Passcode string    `xorm:"passcode"` 	//密码加密随机数
	Passwd   string    `xorm:"passwd"`		//密码
	Hardware string    `xorm:"hardware"`	//硬件信息
	Status   int       `xorm:"status"`		//用户账号状态。0-默认；1-冻结；2-停号
	Ctime    time.Time `xorm:"ctime"`		//登录时间
	Mtime    time.Time `xorm:"mtime"`
	IsOnline bool      `xorm:"-"`			//是否在线
}

type LoginHistory struct {
	Id       int       `xorm:"id pk autoincr"`
	UId      int       `xorm:"uid"`
	CTime    time.Time `xorm:"ctime"`
	Ip       string    `xorm:"ip"`
	State    int8      `xorm:"state"` 		//登录状态，0登录，1登出
	Hardware string    `xorm:"hardware"`
}

type LoginLast struct {
	Id         int       `xorm:"id pk autoincr"`
	UId        int       `xorm:"uid"` 			// 用户uid
	LoginTime  time.Time `xorm:"login_time"` 	//登录时间
	LogoutTime time.Time `xorm:"logout_time"`	//登出时间
	Ip         string    `xorm:"ip"`			//ip
	Session    string    `xorm:"session"`		//session 会话
	IsLogout   int8      `xorm:"is_logout"`		//是否logout,1:logout，0:login
	Hardware   string    `xorm:"hardware"`		//硬件
}

type Role struct {
	RId			int			`xorm:"rid pk autoincr"`
	UId			int			`xorm:"uid"`
	NickName	string		`xorm:"nick_name" validate:"min=4,max=20,regexp=^[a-zA-Z0-9_]*$"`
	Balance		int			`xorm:"balance"`
	HeadId		int16		`xorm:"headId"`
	Sex			int8		`xorm:"sex"`
	Profile		string		`xorm:"profile"`
	LoginTime   time.Time	`xorm:"login_time"`
	LogoutTime  time.Time	`xorm:"logout_time"`
	CreatedAt	time.Time	`xorm:"created_at"`
}

type RoleRes struct {
	Id        		int    		`xorm:"id pk autoincr"`
	RId       		int    		`xorm:"rid"`
	Wood      		int    		`xorm:"wood"`
	Iron      		int    		`xorm:"iron"`
	Stone     		int    		`xorm:"stone"`
	Grain     		int    		`xorm:"grain"`
	Gold      		int    		`xorm:"gold"`
	Decree    		int    		`xorm:"decree"`	//令牌
}

//玩家属性
type RoleAttribute struct {
	Id              int       		`xorm:"id pk autoincr"`
	RId             int       		`xorm:"rid"`
	UnionId         int       		`xorm:"-"`					//联盟id
	ParentId        int       		`xorm:"parent_id"`			//上级id（被沦陷）
	CollectTimes    int8      		`xorm:"collect_times"`		//征收次数
	LastCollectTime time.Time 		`xorm:"last_collect_time"`	//最后征收的时间
	PosTags			string    		`xorm:"pos_tags"`			//位置标记
	PosTagArray		[]PosTag		`xorm:"-"`
}

type PosTag struct {
	X		int	`json:"x"`
	Y		int	`json:"y"`
	Name string `json:"name"`
}
//玩家城池
type MapRoleCity struct {
	mutex		sync.Mutex	`xorm:"-"`
	CityId		int			`xorm:"cityId pk autoincr"`
	RId			int			`xorm:"rid"`
	Name		string		`xorm:"name" validate:"min=4,max=20,regexp=^[a-zA-Z0-9_]*$"`
	X			int			`xorm:"x"`
	Y			int			`xorm:"y"`
	IsMain		int8		`xorm:"is_main"`
	CurDurable	int			`xorm:"cur_durable"`
	CreatedAt	time.Time	`xorm:"created_at"`
	OccupyTime	time.Time 	`xorm:"occupy_time"`
}

//部队
const (
	MapBuildSysFortress = 50	//系统要塞
	MapBuildSysCity = 51		//系统城市
	MapBuildFortress = 56		//玩家要塞
)
//建筑
type MapRoleBuild struct {
	Id    		int    		`xorm:"id pk autoincr"`
	RId   		int    		`xorm:"rid"`
	Type  		int8   		`xorm:"type"`
	Level		int8   		`xorm:"level"`
	OPLevel		int8		`xorm:"op_level"`	//操作level
	X          	int       	`xorm:"x"`
	Y          	int       	`xorm:"y"`
	Name       	string    	`xorm:"name"`
	Wood       	int       	`xorm:"-"`
	Iron       	int       	`xorm:"-"`
	Stone      	int       	`xorm:"-"`
	Grain      	int       	`xorm:"-"`
	Defender   	int       	`xorm:"-"`
	CurDurable 	int       	`xorm:"cur_durable"`
	MaxDurable 	int       	`xorm:"max_durable"`
	OccupyTime 	time.Time 	`xorm:"occupy_time"`
	EndTime 	time.Time 	`xorm:"end_time"`	//建造或升级完的时间
	GiveUpTime 	int64 		`xorm:"giveUp_time"`
}

const (
	ArmyCmdIdle   		= 0	//空闲
	ArmyCmdAttack 		= 1	//攻击
	ArmyCmdDefend 		= 2	//驻守
	ArmyCmdReclamation 	= 3	//屯垦
	ArmyCmdBack   		= 4 //撤退
	ArmyCmdConscript  	= 5 //征兵
	ArmyCmdTransfer  	= 6 //调动
)

const (
	ArmyStop  		= 0
	ArmyRunning  	= 1
)

//军队
type Army struct {
	Id             		int    		`xorm:"id pk autoincr"`
	RId            		int    		`xorm:"rid"`
	CityId         		int    		`xorm:"cityId"`
	Order          		int8   		`xorm:"order"`
	Generals       		string 		`xorm:"generals"`
	Soldiers       		string 		`xorm:"soldiers"`
	ConscriptTimes 		string 		`xorm:"conscript_times"`	//征兵结束时间，json数组
	ConscriptCnts  		string 		`xorm:"conscript_cnts"`		//征兵数量，json数组
	Cmd                	int8       	`xorm:"cmd"`
	FromX              	int        	`xorm:"from_x"`
	FromY              	int        	`xorm:"from_y"`
	ToX                	int        	`xorm:"to_x"`
	ToY                	int        	`xorm:"to_y"`
	Start              	time.Time  	`json:"-"xorm:"start"`
	End                	time.Time  	`json:"-"xorm:"end"`
	State              	int8       	`xorm:"-"` 				//状态:0:running,1:stop
	GeneralArray       	[]int      	`json:"-" xorm:"-"`
	SoldierArray       	[]int      	`json:"-" xorm:"-"`
	ConscriptTimeArray 	[]int64    	`json:"-" xorm:"-"`
	ConscriptCntArray  	[]int      	`json:"-" xorm:"-"`
	Gens               	[]*General 	`json:"-" xorm:"-"`
	CellX              	int        	`json:"-" xorm:"-"`
	CellY              	int        	`json:"-" xorm:"-"`
}

const (
	GeneralNormal      	= 0 //正常
	GeneralComposeStar 	= 1 //星级合成
	GeneralConvert 		= 2 //转换
)
const SkillLimit = 3

type General struct {
	Id            int       		`xorm:"id pk autoincr"`
	RId           int       		`xorm:"rid"`
	CfgId         int       		`xorm:"cfgId"`
	PhysicalPower int       		`xorm:"physical_power"`
	Level         int8      		`xorm:"level"`
	Exp           int       		`xorm:"exp"`
	Order         int8      		`xorm:"order"`
	CityId        int       		`xorm:"cityId"`
	CreatedAt     time.Time 		`xorm:"created_at"`
	CurArms       int       		`xorm:"arms"`
	HasPrPoint    int       		`xorm:"has_pr_point"`
	UsePrPoint    int       		`xorm:"use_pr_point"`
	AttackDis     int  				`xorm:"attack_distance"`
	ForceAdded    int  				`xorm:"force_added"`
	StrategyAdded int  				`xorm:"strategy_added"`
	DefenseAdded  int  				`xorm:"defense_added"`
	SpeedAdded    int  				`xorm:"speed_added"`
	DestroyAdded  int  				`xorm:"destroy_added"`
	StarLv        int8  			`xorm:"star_lv"`
	Star          int8  			`xorm:"star"`
	ParentId      int  				`xorm:"parentId"`
	Skills		  string			`xorm:"skills"`
	SkillsArray   []*model.GSkill	`xorm:"-"`
	State         int8 				`xorm:"state"`
}


func (r *Role) ToModel() interface{}{
	m := model.Role{}
	m.UId = r.UId
	m.RId = r.RId
	m.Sex = r.Sex
	m.NickName = r.NickName
	m.HeadId = r.HeadId
	m.Balance = r.Balance
	m.Profile = r.Profile
	return m
}

func (r *RoleRes) ToModel() interface{} {
	p := model.RoleRes{}
	p.Gold = r.Gold
	p.Grain = r.Grain
	p.Stone = r.Stone
	p.Iron = r.Iron
	p.Wood = r.Wood
	p.Decree = r.Decree

	p.GoldYield = 1
	p.GrainYield = 1
	p.StoneYield = 1
	p.IronYield = 1
	p.WoodYield = 1
	p.DepotCapacity = 10000
	return p
}

func (m *MapRoleCity) ToModel() interface{}{
	p := model.MapRoleCity{}
	p.X = m.X
	p.Y = m.Y
	p.CityId = m.CityId
	p.UnionId = 0
	p.UnionName = ""
	p.ParentId = 0
	p.MaxDurable = 1000
	p.CurDurable = m.CurDurable
	p.Level = 1
	p.RId = m.RId
	p.Name = m.Name
	p.IsMain = m.IsMain == 1
	p.OccupyTime = m.OccupyTime.UnixNano()/1e6
	return p
}

func (m *MapRoleBuild) ToModel() interface{}{
	p := model.MapRoleBuild{}
	p.RNick = "111"
	p.UnionId = 0
	p.UnionName = ""
	p.ParentId = 0
	p.X = m.X
	p.Y = m.Y
	p.Type = m.Type
	p.RId = m.RId
	p.Name = m.Name

	p.OccupyTime = m.OccupyTime.UnixNano()/1e6
	p.GiveUpTime = m.GiveUpTime*1000
	p.EndTime = m.EndTime.UnixNano()/1e6

	p.CurDurable = m.CurDurable
	p.MaxDurable = m.MaxDurable
	p.Defender = m.Defender
	p.Level = m.Level
	p.OPLevel = m.OPLevel
	return p
}

func (g *General) ToModel() interface{}{
	p := model.General{}
	p.CityId = g.CityId
	p.Order = g.Order
	p.PhysicalPower = g.PhysicalPower
	p.Id = g.Id
	p.CfgId = g.CfgId
	p.Level = g.Level
	p.Exp = g.Exp
	p.CurArms = g.CurArms
	p.HasPrPoint = g.HasPrPoint
	p.UsePrPoint = g.UsePrPoint
	p.AttackDis = g.AttackDis
	p.ForceAdded = g.ForceAdded
	p.StrategyAdded = g.StrategyAdded
	p.DefenseAdded = g.DefenseAdded
	p.SpeedAdded = g.SpeedAdded
	p.DestroyAdded = g.DestroyAdded
	p.StarLv = g.StarLv
	p.Star = g.Star
	p.State = g.State
	p.ParentId = g.ParentId
	p.Skills = g.SkillsArray
	return p
}

func (a *Army) ToModel() interface{}{
	p := model.Army{}
	p.CityId = a.CityId
	p.Id = a.Id
	p.UnionId = 0
	p.Order = a.Order
	p.Generals = a.GeneralArray
	p.Soldiers = a.SoldierArray
	p.ConTimes = a.ConscriptTimeArray
	p.ConCnts = a.ConscriptCntArray
	p.Cmd = a.Cmd
	p.State = a.State
	p.FromX = a.FromX
	p.FromY = a.FromY
	p.ToX = a.ToX
	p.ToY = a.ToY
	p.Start = a.Start.Unix()
	p.End = a.End.Unix()
	return p
}
