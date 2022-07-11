package gameconfig

const (
	TypeDurable   		= 1	//耐久
	TypeCost 			= 2
	TypeArmyTeams 		= 3	//队伍数量
	TypeSpeed			= 4	//速度
	TypeDefense			= 5	//防御
	TypeStrategy		= 6	//谋略
	TypeForce			= 7	//攻击武力
	TypeConscriptTime	= 8 //征兵时间
	TypeReserveLimit 	= 9 //预备役上限
	TypeUnkonw			= 10
	TypeHanAddition 	= 11
	TypeQunAddition		= 12
	TypeWeiAddition 	= 13
	TypeShuAddition 	= 14
	TypeWuAddition		= 15
	TypeDealTaxRate		= 16//交易税率
	TypeWood			= 17
	TypeIron			= 18
	TypeGrain			= 19
	TypeStone			= 20
	TypeTax				= 21//税收
	TypeExtendTimes		= 22//扩建次数
	TypeWarehouseLimit 	= 23//仓库容量
	TypeSoldierLimit 	= 24//带兵数量
	TypeVanguardLimit 	= 25//前锋数量
)

type conditions struct {
	Type  int `json:"type"`
	Level int `json:"level"`
}

type facility struct {
	Title      string       `json:"title"`
	Des        string       `json:"des"`
	Name       string       `json:"name"`
	Type       int8         `json:"type"`
	Additions  []int8       `json:"additions"`
	Conditions []conditions `json:"conditions"`
	Levels     []fLevel     `json:"levels"`
}

type NeedRes struct {
	Decree 		int	`json:"decree"`
	Grain		int `json:"grain"`
	Wood		int `json:"wood"`
	Iron		int `json:"iron"`
	Stone		int `json:"stone"`
	Gold		int	`json:"gold"`
}

type fLevel struct {
	Level  int     `json:"level"`
	Values []int   `json:"values"`
	Need   NeedRes `json:"need"`
	Time   int     `json:"time"`	//升级需要的时间
}

type conf struct {
	Name	string
	Type	int8
}

type facilityConf struct {
	Title		string `json:"title"`
	List 		[]conf  `json:"list"`
	facilitys 	map[int8]*facility
}

var FcailityConf = &facilityConf{}

const facilityFile = "/conf/gameconfig/facility/facility.json"
const facilityPath = "/conf/gameconfig/facility/"