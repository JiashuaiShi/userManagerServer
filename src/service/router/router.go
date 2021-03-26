package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
	log "zego.com/userManageServer/src/logger"
)

func RunRouter() {
	router := gin.Default()

	router.Use(
		cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "OPTIONS", "DELETE"},
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			MaxAge:           12 * time.Hour,
			AllowCredentials: true,
		}))

	routerGroup := router.Group("/user")

	routerGroup.POST("/add_user", addHandle)
	routerGroup.GET("/get_user", getHandle)
	routerGroup.GET("/get_userlist", getListHandle)
	routerGroup.DELETE("/del_user/:id", deleteHandle)

	// 端口可配置？
	_ = router.Run(":9001")
	log.Info(nil, "router run at port 9001")
}
