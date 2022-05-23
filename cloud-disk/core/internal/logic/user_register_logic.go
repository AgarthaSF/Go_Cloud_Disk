package logic

import (
	"cloud-disk/core/helper"
	"cloud-disk/core/models"
	"context"
	"errors"
	"log"
	"time"

	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRegisterLogic) UserRegister(req *types.UserRegisterRequest) (resp *types.USerRegisterReply, err error) {
	// 判断验证码是否一致

	code, err := l.svcCtx.RDB.Get(l.ctx, req.Email).Result()
	if err != nil{
		err = errors.New("validation code has not been send yet or has expired")
		return nil, err
	}

	if code != req.Code{
		err = errors.New("validation code mismatches")
		return nil, err
	}

	// 判断用户名是否已存在
	cnt, err := l.svcCtx.Engine.Where("name = ?", req.Name).Count(new(models.UserBasic))
	if err != nil {
		return nil, err
	}
	if cnt > 0{
		err = errors.New("username has been registered before")
		return nil, err
	}
	// 用户名和验证码都正确，就将数据入库

	user := models.UserBasic{
		Identity: helper.GetUUID(),
		Name: req.Name,
		Password: helper.MD5(req.Password),
		Email: req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	n, err := l.svcCtx.Engine.Insert(user)
	if err != nil {
		return nil, err
	}
	log.Println("insert user row: ", n)
	return
}

