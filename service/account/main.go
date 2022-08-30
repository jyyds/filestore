package main

import (
	"log"
	"time"

	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
	"github.com/jyyds/filestore/common"
	"github.com/jyyds/filestore/service/account/handler"
	proto "github.com/jyyds/filestore/service/account/proto"
	dbproxy "github.com/jyyds/filestore/service/dbproxy/client"
	"go-micro.dev/v4"
)

func main() {

	service := micro.NewService(
		// service := k8s.NewService(
		micro.Name("go.micro.service.user"),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Flags(common.CustomFlags...),
	)

	// 初始化service, 解析命令行参数等
	service.Init()

	// 初始化dbproxy client
	dbproxy.Init(service)

	proto.RegisterUserServiceHandler(service.Server(), new(handler.User))
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}
