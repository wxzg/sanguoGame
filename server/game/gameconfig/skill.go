package gameconfig

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var Skill skill

type skill struct {
	skills []SkillConf
	skillConfMap map[int]SkillConf
	outline outline
}


type trigger struct {
	Type int    `json:"type"`
	Des  string `json:"des"`
}

type triggerType struct {
	Des  string     `json:"des"`
	List [] trigger `json:"list"`
}

type effect struct {
	Type   int    `json:"type"`
	Des    string `json:"des"`
	IsRate bool   `json:"isRate"`
}

type effectType struct {
	Des  string    `json:"des"`
	List [] effect `json:"list"`
}

type target struct {
	Type   int    `json:"type"`
	Des    string `json:"des"`
}

type targetType struct {
	Des  string    `json:"des"`
	List [] target `json:"list"`
}


type outline struct {
	TriggerType triggerType `json:"trigger_type"`	//触发类型
	EffectType  effectType  `json:"effect_type"`	//效果类型
	TargetType  targetType  `json:"target_type"`	//目标类型
}


type level struct {
	Probability int   `json:"probability"`  //发动概率
	EffectValue []int `json:"effect_value"` //效果值
	EffectRound []int `json:"effect_round"` //效果持续回合数
}

type SkillConf struct {
	CfgId		  int	  `json:"cfgId"`
	Name          string  `json:"name"`
	Trigger       int     `json:"trigger"` 			//发起类型
	Target        int     `json:"target"`  			//目标类型
	Des           string  `json:"des"`
	Limit         int     `json:"limit"`          	//可以被武将装备上限
	Arms          []int   `json:"arms"`           	//可以装备的兵种
	IncludeEffect []int   `json:"include_effect"` 	//技能包括的效果
	Levels        []level `json:"levels"`
}

const skillFile = "/conf/gameconfig/skill/skill_outline.json"
const skillPath = "/conf/gameconfig/skill/"
func (s *skill) Load()  {
	s.skills = make([]SkillConf, 0)
	s.skillConfMap = make(map[int]SkillConf)
	//获取当前文件路径
	currentDir, _ := os.Getwd()
	//配置文件位置
	cf := currentDir +skillFile
	cp := currentDir + skillPath
	//打包后 程序参数加入配置文件路径
	if len(os.Args) > 1 {
		if path := os.Args[1]; path != ""{
			cf = path + skillFile
			cp = path + skillPath
		}
	}
	data,err := ioutil.ReadFile(cf)
	if err != nil {
		log.Println("技能配置读取失败")
		panic(err)
	}
	err = json.Unmarshal(data, &s.outline)
	if err != nil {
		log.Println("技能配置格式定义失败")
		panic(err)
	}
	files, err := ioutil.ReadDir(cp)
	if err != nil {
		log.Println("技能文件读取失败")
		panic(err)
	}
	for _, v := range files{
		if v.IsDir() {
			//读取目录下的配置
			name := v.Name()
			dirFile := cp + name
			skillFiles, err := ioutil.ReadDir(dirFile)
			if err != nil {
				log.Println(name+"技能文件读取失败")
				panic(err)
			}
			for _, sv := range skillFiles{
				if sv.IsDir() {
					continue
				}
				fileJson := path.Join(dirFile,sv.Name())
				conf := SkillConf{}
				data,err := ioutil.ReadFile(fileJson)
				if err != nil {
					log.Println(name+"技能文件格式错误")
					panic(err)
				}
				err = json.Unmarshal(data,&conf)
				if err != nil {
					log.Println(name+"技能文件格式错误")
					panic(err)
				}
				s.skills = append(s.skills,conf)
				s.skillConfMap[conf.CfgId] = conf
			}
		}
	}
	log.Println(s)
}
