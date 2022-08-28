package rpc

import (
	"context"

	"github.com/jyyds/filestore/service/upload/config"
	"github.com/jyyds/filestore/service/upload/proto"
)

// 用于实现UploadServiceHandler接口
type Upload struct{}

// 用于获取上传的入口地址
func (u *Upload) UploadEntry(ctx context.Context, req *proto.ReqEntry, resp *proto.RespEntry) error {
	resp.Entry = config.UploadEntry
	return nil
}
