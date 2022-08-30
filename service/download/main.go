package main

import (
	"fmt"
	"time"

	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
	"github.com/jyyds/filestore/common"
	dbproxy "github.com/jyyds/filestore/service/dbproxy/client"
	cfg "github.com/jyyds/filestore/service/download/config"
	dlProto "github.com/jyyds/filestore/service/download/proto"
	"github.com/jyyds/filestore/service/download/route"
	dlRpc "github.com/jyyds/filestore/service/download/rpc"
	"go-micro.dev/v4"
)

func startRPCService() {
	service := micro.NewService(
		micro.Name("go.micro.service.download"), // 在注册中心中的服务名称
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Flags(common.CustomFlags...),
	)
	service.Init()

	// 初始化dbproxy client
	dbproxy.Init(service)

	dlProto.RegisterDownloadServiceHandler(service.Server(), new(dlRpc.Download))
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

func startAPIService() {
	router := route.Router()
	router.Run(cfg.DownloadServiceHost)
}

func main() {
	// api 服务
	go startAPIService()

	// rpc 服务
	startRPCService()
}
