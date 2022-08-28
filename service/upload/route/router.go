package route

import (
	"github.com/gin-gonic/gin"
	"github.com/jyyds/filestore/service/upload/api"
)

// Router : 路由表配置
func Router() *gin.Engine {
	// gin framework, 包括Logger, Recovery
	router := gin.Default()

	// 处理静态资源
	router.Static("/static/", "../../static")

	router.POST("/file/upload", api.DoUploadHandler)
	router.OPTIONS("file/upload", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST,OPTIONS")
		c.Status(204)
	})
	return router
}
