package router

import (
	"github.com/gin-gonic/gin"
	"time"
	"github.com/gin-contrib/cors"
)

func RunRouter()  {
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

	_ = router.Run(":9001")
}




