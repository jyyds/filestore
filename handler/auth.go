package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jyyds/filestore/util"
)

// HTTP 连接器
func HTTPInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {

		username := c.Request.FormValue("username")
		token := c.Request.FormValue("token")

		if len(username) < 3 || !IsTokenValid(token) {
			c.Abort()
			resp := util.NewRespMsg(
				500,
				"token无效",
				nil,
			)
			c.JSON(http.StatusOK, resp)
			return
		}
		c.Next()
	}
}
