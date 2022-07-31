package gameconfig

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type generalBasic struct {
	Title	string    `json:"title"`
	Levels	[]gLevel `json:"levels"`
}
type gLevel struct {
	Level		int8`json:"level"`
	Exp			int `json:"exp"`
	Soldiers	int `json:"soldiers"`
}

var GeneralBasic = &generalBasic{}
var generalBasicFile = "/conf/gameconfig/general/general_basic.json"

func (g *generalBasic) Load()  {
	//获取当前文件路径
	currentDir, _ := os.Getwd()
	//配置文件位置
	cf := currentDir +generalBasicFile
	//打包后 程序参数加入配置文件路径
	if len(os.Args) > 1 {
		if path := os.Args[1]; path != ""{
			cf = path + generalBasicFile
		}
	}
	data,err := ioutil.ReadFile(cf)
	if err != nil {
		log.Println("武将基本配置读取失败")
		panic(err)
	}
	err = json.Unmarshal(data, g)
	if err != nil {
		log.Println("武将基本配置格式定义失败")
		panic(err)
	}
}

func (g *generalBasic) GetLevel(level int8) *gLevel {
	for _, v:= range g.Levels{
		if v.Level == level {
			return &v
		}
	}
	return nil
}