package logic

import (
	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFolderCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFolderCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFolderCreateLogic {
	return &UserFolderCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFolderCreateLogic) UserFolderCreate(req *types.UserFolderCreateRequest, userIdentity string) (resp *types.UserFolderCreateReply, err error) {

	// check whether the folder name already exists in current directory
	cnt, err := l.svcCtx.Engine.Where("name = ? and parent_id = ?",
		req.Name, req.ParentId).Count(new(models.UserRepository))
	if err != nil {
		return
	}

	if cnt > 0{
		err = errors.New("this folder already exists")
		return
	}

	// create folder

	data := &models.UserRepository{
		Identity:           helper.GetUUID(),
		UserIdentity:       userIdentity,
		ParentId:           req.ParentId,
		Name:               req.Name,
	}

	_, err = l.svcCtx.Engine.Insert(data)
	if err != nil {
		return
	}

	resp = &types.UserFolderCreateReply{}
	resp.Identity = data.Identity
	return
}
