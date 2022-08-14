package mq

import "github.com/jyyds/filestore/common"

// 转移队列中消息载体的结构体格式
type TransferData struct {
	FileHash      string
	CurLocation   string
	DestLocation  string
	DestStoreType common.StoreType
}
