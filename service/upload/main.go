package main

import (
	"log"

	"github.com/go-micro/plugins/v4/registry/consul"
	cfg "github.com/jyyds/filestore/service/upload/config"
	upProto "github.com/jyyds/filestore/service/upload/proto"
	"github.com/jyyds/filestore/service/upload/route"
	upRpc "github.com/jyyds/filestore/service/upload/rpc"
	"go-micro.dev/v4"
)

func startRpcService() {
	consulReg := consul.NewRegistry()

	service := micro.NewService(
		micro.Name("go.micro.service.upload"),
		micro.Registry(consulReg),
	)
	service.Init()

	upProto.RegisterUploadServiceHandler(service.Server(), new(upRpc.Upload))
	if err := service.Run(); err != nil {
		log.Println(err)
	}

}

func startApiService() {
	router := route.Router()
	router.Run(cfg.UploadServicehost)
}

func main() {
	// rpc 服务
	go startRpcService()

	// api 服务
	startApiService()
}
