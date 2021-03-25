package main

import (
	"log"
	"zego.com/userManageServer/src/service/redis"
	"zego.com/userManageServer/src/service/router"
)

func main() {
	// 读取配置

	// 初始化数据库连接
	err := redis.InitRedis()
	if err != nil {
		log.Println("InitRedis failed")
		return
	}

	// 启动gin服务
	router.RunRouter()
}
