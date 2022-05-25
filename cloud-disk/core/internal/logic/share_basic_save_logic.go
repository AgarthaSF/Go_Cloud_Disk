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

type ShareBasicSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareBasicSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareBasicSaveLogic {
	return &ShareBasicSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareBasicSaveLogic) ShareBasicSave(req *types.ShareBasicSaveRequest, userIdentity string) (resp *types.ShareBasicSaveReply, err error) {

	// get the filename from table user_repository
	ur := &models.UserRepository{
		Identity: req.UserRepositoryIdentity,
	}
	has, err := l.svcCtx.Engine.Get(ur)
	if err != nil {
		return
	}
	if !has {
		err = errors.New("file not found")
		return
	}

	// get the target resource's detailed information
	rp := new(models.RepositoryPool)
	has, err= l.svcCtx.Engine.Where("identity = ?", ur.RepositoryIdentity).Get(rp)
	if err != nil {
		return
	}
	if !has {
		err = errors.New("file not found")
		return
	}

	// check whether the target folder already have the file with same name
	cnt, err := l.svcCtx.Engine.
		Where("user_identity = ? and parent_id = ? and name = ?", userIdentity, req.ParentId, ur.Name).
		Count(new(models.UserRepository))
	if err != nil {
		return
	}
	if cnt > 0 {
		err = errors.New("target folder already has the file with same name")
		return
	}

	// save resource
	data := &models.UserRepository{
		Identity:           helper.GetUUID(),
		UserIdentity:       userIdentity,
		ParentId:           req.ParentId,
		RepositoryIdentity: ur.RepositoryIdentity,
		Ext:                rp.Ext,
		Name:               ur.Name,
	}

	_, err = l.svcCtx.Engine.Insert(data)
	resp = new(types.ShareBasicSaveReply)
	resp.Identity = data.Identity
	return
}
