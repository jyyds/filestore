package oss

import (
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	cfg "github.com/jyyds/filestore/config"
)

var ossCli *oss.Client

// 创建oss client 对象
func Client() *oss.Client {
	if ossCli != nil {
		return ossCli
	}
	ossCli, err := oss.New(cfg.OSSEndpoint, cfg.OSSAccesskeyID, cfg.OSSAcessKetSecret)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return ossCli
}

// Bucket : 获取bucket存储空间
func Bucket() *oss.Bucket {
	cli := Client()
	if cli != nil {
		bucket, err := cli.Bucket(cfg.OSSBucket)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		return bucket
	}
	return nil
}
