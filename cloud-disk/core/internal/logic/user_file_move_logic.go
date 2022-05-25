package logic

import (
	"cloud-disk/core/models"
	"context"
	"errors"

	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileMoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileMoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileMoveLogic {
	return &UserFileMoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileMoveLogic) UserFileMove(req *types.UserFileMoveRequest, userIdentity string) (resp *types.UserFileMoveReply, err error) {

	// check whether the target folder exists
	parent := new(models.UserRepository)
	has, err := l.svcCtx.Engine.Where("identity = ? and user_identity = ?", req.ParentIdentity, userIdentity).Get(parent)
	if err != nil {
		return
	}
	if !has{
		err = errors.New("this folder does not exist")
		return
	}

	// get the source file name
	src := new(models.UserRepository)
	has, err = l.svcCtx.Engine.Where("identity = ? and user_identity = ?",  req.Identity, userIdentity).Get(src)
	if err != nil {
		return
	}

	// check whether the target folder already have the file with same name
	cnt, err := l.svcCtx.Engine.Where("parent_id = ? and user_identity = ? and name = ?", parent.Id, userIdentity, src.Name).
		Limit(1,0).
		Count(new(models.UserRepository))
	if err != nil {
		return
	}
	if cnt > 0{
		err = errors.New("target folder already has the file with same name")
		return
	}

	// update
	_, err = l.svcCtx.Engine.Where("identity = ? and user_identity = ?", req.Identity, userIdentity).
	Update(models.UserRepository{
		ParentId: int64(parent.Id),
	})

	return
}
