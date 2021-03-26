package main

import (
	"github.com/sirupsen/logrus"
	log "zego.com/userManageServer/src/logger"
	"zego.com/userManageServer/src/service/redis"
	"zego.com/userManageServer/src/service/router"
	"zego.com/userManageServer/src/tools"
)

func main() {
	// 读取配置
	path := "../conf/conf.yml"
	logrus.Infoln("读取配置：" + path)
	conf := tools.GetConf(path)
	log.Init()
	//初始化数据库连接
	if err := redis.InitRedis(conf); err != nil {
		log.Info(nil, "InitRedis failed")
		return
	}
	//用户a登陆失败  "user:%s error:%+v",user_id,err
	// 启动gin服务
	router.RunRouter()
}
