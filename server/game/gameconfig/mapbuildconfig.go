package gameconfig

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var MapBuilder = &mapBuildConf{

}

type ConfigRsp struct {
	Confs []Conf
}

type ConfigReq struct {

}

type Conf struct {
	Type		int8		`json:"type"`
	Level		int8		`json:"level"`
	Name		string		`json:"name"`
	Wood		int			`json:"Wood"`
	Iron		int			`json:"iron"`
	Stone		int			`json:"stone"`
	Grain		int			`json:"grain"`
	Durable		int			`json:"durable"`	//耐久
	Defender	int			`json:"defender"`	//防御等级
}

type cfg struct {
	Type     int8   `json:"type"`
	Name     string `json:"name"`
	Level    int8   `json:"level"`
	Grain    int    `json:"grain"`
	Wood     int    `json:"wood"`
	Iron     int    `json:"iron"`
	Stone    int    `json:"stone"`
	Durable  int    `json:"durable"`
	Defender int    `json:"defender"`
}


type mapBuildConf struct {
	Title   string   `json:"title"`
	Cfg		[]cfg `json:"cfg"`
	cfgMap  map[int8][]cfg
}

const (
	mapBuilderConfFile = "/conf/gameconfig/map_build.json"
)

//获取地图配置信息
func (m *mapBuildConf) Load(){
	current, _ := os.Getwd()
	path := current + mapBuilderConfFile

	if len(os.Args) > 1 {
		if  os.Args[1] != "" {
			path = os.Args[1] + mapBuilderConfFile
		}
	}

	data, err := ioutil.ReadFile(path)
	if err != nil{
		log.Println("地图资源读取失败")
		panic(err)
	}

	err = json.Unmarshal(data, m)

	if err != nil{
		log.Println("地图资源JSON转换失败")
		panic(err)
	}

	for _, v := range m.Cfg{
		_, ok := m.cfgMap[v.Type]
		if !ok {
			m.cfgMap[v.Type] = make([]cfg,0)
		} else {
			m.cfgMap[v.Type] = append(m.cfgMap[v.Type])
		}
	}

}

// 加载角色建筑配置
func(m *mapBuildConf) BuildConfig(buildType, level int8) *cfg{
	cfgs := m.cfgMap[buildType]

	for _,v := range cfgs{
		if v.Level == level{
			return &v
		}
	}

	return nil
}
