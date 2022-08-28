package route

import (
	"github.com/gin-gonic/gin"
	"github.com/jyyds/filestore/service/apigw/handler"
)

// 网关api
func Router() *gin.Engine {
	router := gin.Default()

	router.Static("/static", "../../static")

	router.GET("/user/signup", handler.SignupHandler)
	router.POST("/user/signup", handler.DoSignupHandler)

	return router

}
