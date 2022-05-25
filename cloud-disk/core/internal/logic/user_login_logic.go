package logic

import (
	"cloud-disk/core/define"
	"cloud-disk/core/helper"
	"cloud-disk/core/models"
	"context"
	"errors"

	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginLogic) UserLogin(req *types.LoginRequest) (resp *types.LoginReply, err error) {
	// search current user from database
	user := new(models.UserBasic)
	get, err := l.svcCtx.Engine.Where("name = ? and password = ?", req.Name, helper.MD5(req.Password)).Get(user)
	if err != nil {
		return nil, err
	}

	if !get {
		return nil, errors.New("user name or password is wrong")
	}

	// generate normal token
	token, err := helper.GenerateToken(user.Id, user.Identity, user.Name, define.TokenExpireTime)
	if err != nil {
		return
	}

	// generate the "refresh token" which is used to refresh other token
	refreshToken, err := helper.GenerateToken(user.Id, user.Identity, user.Name, define.RefreshTokenExpireTime)
	if err != nil {
		return
	}

	resp = new(types.LoginReply)
	resp.Token = token
	resp.RefreshToken = refreshToken

	return
}
