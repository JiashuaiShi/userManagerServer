package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
	log "zego.com/userManageServer/src/logger"
)

var router *gin.Engine

func InitRouter() {
	router := gin.Default()

	// Albert提供解决客户端跨域访问的解决方案
	router.Use(
		cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "OPTIONS", "DELETE"},
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			MaxAge:           12 * time.Hour,
			AllowCredentials: true,
		}))

	// 使用post方法实现增删改查
	routerGroup := router.Group("/user")
	routerGroup.POST("/add_user", addHandle)
	routerGroup.GET("/get_user", getHandle)
	routerGroup.GET("/get_userlist", getListHandle)
	routerGroup.DELETE("/del_user/:id", deleteHandle)
}

// 指定端口号启动服务
func RunRouter(port string) {
	_ = router.Run(port)
	log.Info(nil, "router run at port 9001")
}
