package model

import (
	"sanguoServer/server/game/gameconfig"
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
	PosTagArray		[]model.PosTag		`xorm:"-"`
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



//技能
type Skill struct {
	Id             	int    `xorm:"id pk autoincr"`
	RId          	int    `xorm:"rid"`
	CfgId          	int    `xorm:"cfgId"`
	BelongGenerals 	string `xorm:"belong_generals"`
	Generals 		[]int  `xorm:"-"`
}

func NewSkill(rid int, cfgId int) *Skill{
	return &Skill{
		CfgId: cfgId,
		RId: rid,
		Generals: []int{},
		BelongGenerals: "[]",
	}
}

func (s *Skill) ToModel() interface{}{
	p := model.Skill{}
	p.Id = s.Id
	p.CfgId = s.CfgId
	p.Generals = s.Generals
	return p
}

//战报
type WarReport struct {
	Id                	int    		`xorm:"id pk autoincr"`
	AttackRid         	int    		`xorm:"a_rid"`
	DefenseRid        	int    		`xorm:"d_rid"`
	BegAttackArmy     	string 		`xorm:"b_a_army"`
	BegDefenseArmy    	string 		`xorm:"b_d_army"`
	EndAttackArmy     	string 		`xorm:"e_a_army"`
	EndDefenseArmy    	string 		`xorm:"e_d_army"`
	BegAttackGeneral  	string 		`xorm:"b_a_general"`
	BegDefenseGeneral 	string 		`xorm:"b_d_general"`
	EndAttackGeneral  	string 		`xorm:"e_a_general"`
	EndDefenseGeneral 	string    	`xorm:"e_d_general"`
	Result				int      	`xorm:"result"`	//0失败，1打平，2胜利
	Rounds				string		`xorm:"rounds"` //回合
	AttackIsRead      	bool      	`xorm:"a_is_read"`
	DefenseIsRead     	bool      	`xorm:"d_is_read"`
	DestroyDurable    	int       	`xorm:"destroy"`
	Occupy            	int       	`xorm:"occupy"`
	X                 	int       	`xorm:"x"`
	Y                 	int       	`xorm:"y"`
	CTime             	time.Time 	`xorm:"ctime"`
}


func (w *WarReport) ToModel() interface{}{
	p := model.WarReport{}
	p.CTime = int(w.CTime.UnixNano() / 1e6)
	p.Id = w.Id
	p.AttackRid = w.AttackRid
	p.DefenseRid = w.DefenseRid
	p.BegAttackArmy = w.BegAttackArmy
	p.BegDefenseArmy = w.BegDefenseArmy
	p.EndAttackArmy = w.EndAttackArmy
	p.EndDefenseArmy = w.EndDefenseArmy
	p.BegAttackGeneral = w.BegAttackGeneral
	p.BegDefenseGeneral = w.BegDefenseGeneral
	p.EndAttackGeneral = w.EndAttackGeneral
	p.EndDefenseGeneral = w.EndDefenseGeneral
	p.Result = w.Result
	p.Rounds = w.Rounds
	p.AttackIsRead = w.AttackIsRead
	p.DefenseIsRead = w.DefenseIsRead
	p.DestroyDurable = w.DestroyDurable
	p.Occupy = w.Occupy
	p.X = w.X
	p.X = w.X
	return p
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
	p.UnionId = model.GetUnion(m.RId)
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


func (m *MapRoleBuild) Init() {
	if cfg := gameconfig.MapBuilder.BuildConfig(m.Type, m.Level); cfg != nil {
		m.Name = cfg.Name
		m.Level = cfg.Level
		m.Type = cfg.Type
		m.Wood = cfg.Wood
		m.Iron = cfg.Iron
		m.Stone = cfg.Stone
		m.Grain = cfg.Grain
		m.MaxDurable = cfg.Durable
		m.CurDurable = cfg.Durable
		m.Defender = cfg.Defender
	}
}

func (m *MapRoleBuild) IsSysCity() bool {
	return m.Type == gameconfig.MapBuildSysCity
}
func (m* MapRoleBuild) IsSysFortress() bool  {
	return m.Type == gameconfig.MapBuildSysFortress
}
func (m* MapRoleBuild) IsRoleFortress() bool  {
	return m.Type == gameconfig.MapBuildFortress
}


