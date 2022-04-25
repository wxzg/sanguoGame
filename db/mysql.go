package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	config "sanguoServer/conf"
	"xorm.io/xorm"
)

var Eg *xorm.Engine

func DBInit(){
	//获取ini-mysql配置
	sqlConfig, err := config.File.GetSection("mysql")
	if err != nil {
		log.Println("数据库配置获取出错",err)
		panic(err)
	}

	// userName:password@tcp(host:port)/dbname?charset=utf8&parseTime=True&loc=Local
	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		sqlConfig["user"],
		sqlConfig["password"],
		sqlConfig["host"],
		sqlConfig["port"],
		sqlConfig["dbname"],
	)
	Eg, err = xorm.NewEngine("mysql", dbConn)
	if err != nil{
		log.Println("数据库连接失败",err)
		panic(err)
	}

	//测试数据库是否连接成功
	err = Eg.Ping()
	if err != nil{
		log.Println("数据库ping失败",err)
		panic(err)
	}
	//最大空闲连接数
	maxIdle := config.File.MustInt("mysql", "max_idle", 2)
	//最大打开连接数
	maxConn := config.File.MustInt("mysql", "max_conn", 10)
	Eg.SetMaxIdleConns(maxIdle)
	Eg.SetMaxOpenConns(maxConn)
	log.Println("数据库初始化成功！")
}