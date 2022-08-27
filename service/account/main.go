package main

import (
	"log"

	"github.com/go-micro/plugins/v4/registry/consul"
	"github.com/jyyds/filestore/service/account/handler"
	"github.com/jyyds/filestore/service/account/proto"
	"go-micro.dev/v4"
)

func main() {

	regConusl := consul.NewRegistry()

	// 创建一个service
	service := micro.NewService(
		micro.Name("go.micro.service.user"),
		micro.Registry(regConusl),
	)
	service.Init()

	proto.RegisterUserServiceHandler(service.Server(), new(handler.User))
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}
