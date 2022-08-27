package handler

import (
	"context"

	"github.com/jyyds/filestore/common"
	cfg "github.com/jyyds/filestore/config"
	dblayer "github.com/jyyds/filestore/db"
	"github.com/jyyds/filestore/service/account/proto"
	"github.com/jyyds/filestore/util"
)

type User struct{}

func (u *User) Signup(ctx context.Context, req *proto.ReqSignup, resp *proto.RespSIgnup) error {

	username := req.Username
	passwd := req.Password

	// 校验用户名密码
	if len(username) < 3 || len(passwd) < 5 {

		resp.Code = common.StatusParamInvalid
		resp.Message = "注册参数无效"
		return nil
	}

	// 对密码进行加盐及取Sha1值加密
	encPasswd := util.Sha1([]byte(passwd + cfg.PasswordSalt))
	// 将用户信息注册到用户表中
	suc := dblayer.UserSignup(username, encPasswd)
	if suc {
		resp.Code = common.StatusOK
		resp.Message = "注册成功"
	} else {
		resp.Code = common.StatusRegisterFailed
		resp.Message = "注册失败"
	}
	return nil
}
