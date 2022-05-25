package logic

import (
	"cloud-disk/core/define"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type UserFileListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileListLogic {
	return &UserFileListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileListLogic) UserFileList(req *types.UserFileListRequest, userIdentity string) (resp *types.UserFileListReply, err error) {

	uf := make([]*types.UserFile, 0)
	var cnt int64
	resp = new(types.UserFileListReply)

	// get pagination parameter
	size := req.Size
	if size == 0 {
		size = define.PageSize
	}
	page := req.Page
	if page == 0 {
		page = 1
	}

	// calculate the offset
	offset := (page - 1) * size

	// kind of complex sql query

	l.svcCtx.Engine.ShowSQL(true)

	err = l.svcCtx.Engine.Table("user_repository").
		Where("parent_id = ? and user_identity = ?", req.Id, userIdentity).
		Select("user_repository.id, user_repository.identity, user_repository.repository_identity, "+
			"user_repository.ext, user_repository.name, repository_pool.path, repository_pool.size").
		Join("LEFT", "repository_pool", "user_repository.repository_identity = repository_pool.identity").
		Where("user_repository.deleted_at = ? or user_repository.deleted_at is null", time.Time{}.Format(define.DateTime)).
		Limit(size, offset).Find(&uf)
	if err != nil {
		return
	}

	// get file total num
	cnt, err = l.svcCtx.Engine.Where("parent_id = ? and user_identity = ?",
		req.Id, userIdentity).Count(new(models.UserRepository))
	if err != nil {
		return
	}

	resp.List = uf
	resp.Count = cnt
	return
}
