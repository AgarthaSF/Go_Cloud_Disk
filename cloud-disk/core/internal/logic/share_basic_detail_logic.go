package logic

import (
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShareBasicDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareBasicDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareBasicDetailLogic {
	return &ShareBasicDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareBasicDetailLogic) ShareBasicDetail(req *types.ShareBasicDetailRequest) (resp *types.ShareBasicDetailReply, err error) {
	// add the click num of the shared record
	_, err = l.svcCtx.Engine.Exec("UPDATE share_basic SET click_num = click_num + 1 WHERE identity = ? ", req.Identity)
	if err != nil {
		return
	}
	resp = new(types.ShareBasicDetailReply)

	// get the detailed information of the shared resource
	_, err = l.svcCtx.Engine.Table("share_basic").
		Select("share_basic.repository_identity, user_repository.name," +
			" repository_pool.ext, repository_pool.size, repository_pool.path").
		Join("left", "repository_pool",
			"share_basic.repository_identity = repository_pool.identity").
		Join("left", "user_repository",
			"share_basic.user_repository_identity = user_repository.identity").
		Where("share_basic.identity = ?", req.Identity).
		Get(resp)

	return
}
