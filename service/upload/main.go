package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
	"github.com/jyyds/filestore/common"
	"github.com/jyyds/filestore/mq"
	dbproxy "github.com/jyyds/filestore/service/dbproxy/client"
	cfg "github.com/jyyds/filestore/service/upload/config"
	upProto "github.com/jyyds/filestore/service/upload/proto"
	"github.com/jyyds/filestore/service/upload/route"
	upRpc "github.com/jyyds/filestore/service/upload/rpc"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4"
)

func startRPCService() {
	service := micro.NewService(
		micro.Name("go.micro.service.upload"), // 服务名称
		micro.RegisterTTL(time.Second*10),     // TTL指定从上一次心跳间隔起，超过这个时间服务会被服务发现移除
		micro.RegisterInterval(time.Second*5), // 让服务在指定时间内重新注册，保持TTL获取的注册时间有效
		micro.Flags(common.CustomFlags...),
	)
	service.Init(
		micro.Action(func(c *cli.Context) error {
			// 检查是否指定mqhost
			mqhost := c.String("mqhost")
			if len(mqhost) > 0 {
				log.Println("custom mq address: " + mqhost)
				mq.UpdateRabbitHost(mqhost)
			}
			return nil
		}),
	)

	// 初始化dbproxy client
	dbproxy.Init(service)
	// 初始化mq client
	mq.Init()

	upProto.RegisterUploadServiceHandler(service.Server(), new(upRpc.Upload))
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

func startApiService() {
	router := route.Router()
	router.Run(cfg.UploadServicehost)
}

func main() {
	// rpc 服务
	go startRPCService()

	// api 服务
	startApiService()
}
