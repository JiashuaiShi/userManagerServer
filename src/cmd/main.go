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
	// 从flag获取配置文件路径
	flag.StringVar(&configPath, "configPath", "./src/conf/config.yml", "input config path")
	flag.Parse()

	// 解析配置文件
	config, err := tools.ParseConf(configPath)
	if err != nil {
		fmt.Printf("parse config failed! error: %+v\n", err)
		return
	}

	// 配置日志库
	logConfig := config.LogConfig
	if err := log.Init(logConfig); err != nil {
		fmt.Printf("log init failed! error: %+v\n", err)
		return
	}

	log.Info(nil, "log init success")

	// 初始化数据库连接
	redisConf := config.RedisConfig

	if err := redis.InitRedis(redisConf); err != nil {
		logField := log.Field{
			"address":  redisConf.Addr,
			"password": redisConf.Password,
			"db":       redisConf.Db,
		}
		log.Error(logField, "redis init failed. error: %+v", err)
		return
	}

	log.Info(nil, "redis init success")

	// 启动gin服务
	port := config.RouterConfig.Port
	router.InitRouter()
	router.RunRouter(port)
}
