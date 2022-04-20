package config

import (
	"errors"
	"github.com/Unknwon/goconfig"
	"log"
	"os"
)

//config配置文件的绝对路径
const configFile = "/conf/conf.ini"
//File 使用一个File变量接收读取出来的configFile
var File *goconfig.ConfigFile

func init(){
	//拿到当前的程序的目录
	currentDir ,err := os.Getwd()
	if err != nil{
		panic(err)
	}

	// 拼接形成config.ini文件地址
	configPath := currentDir + configFile

	//如果命令行参数指定了路径则使用命令行的路径
	if len(os.Args)>1{
		path := os.Args[1]
		if path != ""{
			configPath = currentDir+path
		}
	}
	//判断文件是否存在
	if !fileExist(configPath) {
		panic(errors.New("配置文件不存在"))
	}

	//文件系统的读取
	File, err = goconfig.LoadConfigFile(configPath)
	if err != nil {
		log.Fatal("配置文件出错：", err)
	}
}

// fileExist 判断文件是否存在
func fileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || os.IsExist(err)
}
