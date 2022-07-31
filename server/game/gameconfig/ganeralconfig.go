package gameconfig

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
)

type generals struct {
	Title string 	`json:"title"`
	GArr  []generalDetail    	`json:"list"` //属性
	GMap  map[int]generalDetail
	totalProbability int
}
type generalDetail struct {
	Name     		string	`json:"name"`
	CfgId    		int		`json:"cfgId"`
	Force    		int		`json:"force"`  //武力
	Strategy 		int		`json:"strategy"` //策略
	Defense  		int		`json:"defense"` //防御
	Speed    		int		`json:"speed"` //速度
	Destroy      	int   	`json:"destroy"` //破坏力
	ForceGrow    	int   	`json:"force_grow"`
	StrategyGrow 	int   	`json:"strategy_grow"`
	DefenseGrow  	int   	`json:"defense_grow"`
	SpeedGrow   	int   	`json:"speed_grow"`
	DestroyGrow  	int   	`json:"destroy_grow"`
	Cost         	int8  	`json:"cost"`
	Probability  	int   	`json:"probability"`
	Star         	int8   	`json:"star"`
	Arms         	[]int 	`json:"arms"`
	Camp         	int8  	`json:"camp"`
}

var General = &generals{}
var generalFile = "/conf/gameconfig/general/general.json"

func (g *generals) Load()  {
	//获取当前文件路径
	currentDir, _ := os.Getwd()
	//配置文件位置
	cf := currentDir +generalFile
	//打包后 程序参数加入配置文件路径
	if len(os.Args) > 1 {
		if path := os.Args[1]; path != ""{
			cf = path + generalFile
		}
	}
	data,err := ioutil.ReadFile(cf)
	if err != nil {
		log.Println("武将配置读取失败")
		panic(err)
	}
	err = json.Unmarshal(data, g)
	if err != nil {
		log.Println("武将配置格式定义失败")
		panic(err)
	}
	g.GMap = make(map[int]generalDetail)
	for _,v := range g.GArr{
		g.GMap[v.CfgId] = v
		g.totalProbability += v.Probability
	}
}

//随机武将
func (g *generals) Rand() int {
	//随机一个概率
	rate := rand.Intn(g.totalProbability)
	var current = 0

	for _,v := range g.GArr{
		if rate >= current && rate < current + v.Probability {
			return v.CfgId
		}
		current += v.Probability
	}
	return 0
}

func (g *generals) Cost(cfgId int) int8 {
	c, ok := g.GMap[cfgId]
	if ok {
		return c.Cost
	}else{
		return 0
	}
}
