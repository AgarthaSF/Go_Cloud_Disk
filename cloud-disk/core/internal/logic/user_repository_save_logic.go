package logic

import (
	"cloud-disk/core/helper"
	"cloud-disk/core/models"
	"context"
	"time"

	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRepositorySaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRepositorySaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRepositorySaveLogic {
	return &UserRepositorySaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRepositorySaveLogic) UserRepositorySave(req *types.UserRepositorySaveRequest, userIdentity string) (resp *types.UserRepositorySaveReply, err error) {

	ur := models.UserRepository{
		Identity:           helper.GetUUID(),
		UserIdentity:       userIdentity,
		ParentId:           req.ParentId,
		RepositoryIdentity: req.Repository_Identity,
		Ext:                req.Ext,
		Name:               req.Name,
		CreatedAt:          time.Time{},
		UpdatedAt:          time.Time{},
		DeletedAt:          time.Time{},
	}

	_, err = l.svcCtx.Engine.Insert(ur)
	if err != nil {
		return nil, err
	}
	return
}
