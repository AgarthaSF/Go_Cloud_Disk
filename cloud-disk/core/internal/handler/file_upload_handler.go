package handler

import (
	"cloud-disk/core/helper"
	"cloud-disk/core/models"
	"crypto/md5"
	"fmt"
	"net/http"
	"path"

	"cloud-disk/core/internal/logic"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileUploadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// use FormFile to parse the uploaded file to get the file content and file header
		// file header consists of
		//    Filename string
		//    Header   textproto.MIMEHeader
		//    Size     int64
		//    content  []byte
		//    tmpfile  string
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			return
		}

		// judge whether the uploaded file exists
		b := make([]byte, fileHeader.Size)
		_, err = file.Read(b)
		if err != nil {
			return
		}
		// %x means hex output
		hash := fmt.Sprintf("%x", md5.Sum(b))
		rp := new(models.RepositoryPool)
		has, err := svcCtx.Engine.Where("hash = ?", hash).Get(rp)
		if err != nil {
			return
		}
		if has{
			// if the uploaded file exists in COS, then return the fileIdentity directly
			httpx.OkJson(w,&types.FileUploadReply{
				Identity: rp.Identity,
				Ext: rp.Ext,
				Name: rp.Name,
			})
			return
		}

		// if the uploaded file does not exist, then upload it into COS
		cosPath, err := helper.CosUpload(r)
		if err != nil {
			return 
		}

		// pass request to logic
		req.Name = fileHeader.Filename
		req.Ext = path.Ext(fileHeader.Filename)
		req.Size = fileHeader.Size
		req.Hash = hash
		req.Path = cosPath

		l := logic.NewFileUploadLogic(r.Context(), svcCtx)
		resp, err := l.FileUpload(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
