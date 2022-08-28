package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-micro/plugins/v4/registry/consul"
	"github.com/jyyds/filestore/service/account/proto"
	"go-micro.dev/v4"
)

var (
	userCli proto.UserService
)

func init() {
	consulReg := consul.NewRegistry()
	service := micro.NewService(
		// micro.Name("go.micro.api.user"),
		micro.Registry(consulReg),
	)

	service.Init()

	userCli = proto.NewUserService("go.micro.service.user", service.Client())

}

// SignupHandler : 响应注册页面
func SignupHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signup.html")
}

// DoSignupHandler : 处理注册post请求
func DoSignupHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	passwd := c.Request.FormValue("password")

	resp, err := userCli.Signup(context.TODO(), &proto.ReqSignup{
		Username: username,
		Password: passwd,
	})
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": resp.Code,
		"msg":  resp.Message,
	})
}
