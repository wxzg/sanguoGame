package gameconfig

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sanguoServer/server/game/logic"
)

type NationalMap struct {
	MId			int		`xorm:"mid"`
	X			int		`xorm:"x"`
	Y			int		`xorm:"y"`
	Type		int8	`xorm:"type"`
	Level		int8	`xorm:"level"`
}

type mapRes struct {
	Confs map[int]NationalMap
	SysBuild map[int]NationalMap
}

type mapData struct {
	Width	int 			`json:"w"`
	Height	int				`json:"h"`
	List	[][]int			`json:"list"`
}

const (
	MapBuildSysFortress = 50	//系统要塞
	MapBuildSysCity = 51		//系统城市
	MapBuildFortress = 56		//玩家要塞
)

var MapRes = &mapRes{
	Confs: make(map[int]NationalMap),
	SysBuild: make(map[int]NationalMap),
}
const mapFile = "conf/gameconfig/map.json"

func (m *mapRes) Load() {

	currentDir, _ := os.Getwd()
	//配置文件位置
	cf := currentDir +mapFile
	//打包后 程序参数加入配置文件路径
	if len(os.Args) > 1 {
		if path := os.Args[1]; path != ""{
			cf = path + mapFile
		}
	}
	data,err := ioutil.ReadFile(cf)
	if err != nil {
		log.Println("地图读取失败")
		panic(err)
	}
	mapd := &mapData{}
	err = json.Unmarshal(data, mapd)
	if err != nil {
		log.Println("地图格式定义失败")
		panic(err)
	}

	logic.Mapwidth = mapd.Width
	logic.MapHeight = mapd.Height

	for index, v := range mapd.List{
		//类型，就是这块土地上面是个啥玩意
		t := int8(v[0])
		//等级
		l := int8(v[1])
		nm := NationalMap{
			X : index % logic.Mapwidth,
			Y : index / logic.MapHeight,
			MId: index,
			Type: t,
			Level: l,
		}

		m.Confs[index] = nm
		if t == MapBuildFortress || t == MapBuildSysFortress{
			m.SysBuild[index] = nm
		}
	}
}