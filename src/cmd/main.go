package main

import (
	"flag"
	"fmt"
	log "zego.com/userManageServer/src/logger"
	"zego.com/userManageServer/src/service/redis"
	"zego.com/userManageServer/src/service/router"
	"zego.com/userManageServer/src/tools"
)

var configPath string

func main() {
	// 读取配置文件
	flag.StringVar(&configPath, "configPath", "./src/conf/config.yml", "input config path")
	flag.Parse()

	config, err := tools.ParseConf(configPath)
	if err != nil {
		fmt.Printf("parse config failed! error: %+v\n", err)
		return
	}
	// 配置日志库
	logConfig := config.LogConfig
	if err := log.Init(logConfig); err != nil {

	}
	// 初始化数据库连接
	redisConf := config.RedisConfig
	if err := redis.InitRedis(redisConf); err != nil {
		log.Error(log.Field{
			"address":  redisConf.Addr,
			"password": redisConf.Password,
			"DbNo.":    redisConf.Db,
		},
			"redis init failed. error: %+v", err)
		//return
	} else {
		log.Info(nil, "redis init success")
	}

	// 启动gin服务
	router.RunRouter()
}
