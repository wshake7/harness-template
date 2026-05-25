package logic

import (
	"admin/internal/appsvc"
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	"admin/internal/services/orm/models"
	"strings"
)

type StorageFileHandler struct {
	Storage        appsvc.FileStorage
	MaxUploadBytes int64
}

func NewStorageFileHandler(storage appsvc.FileStorage, maxUploadBytes int64) *StorageFileHandler {
	return &StorageFileHandler{
		Storage:        storage,
		MaxUploadBytes: maxUploadBytes,
	}
}

type ReqStorageFileID struct {
	ID uint64 `json:"id" binding:"required" binding_msg:"required=请求错误"`
}

type ReqStorageFilePresigned struct {
	ID             uint64 `json:"id" binding:"required" binding_msg:"required=请求错误"`
	ExpiresSeconds int    `json:"expiresSeconds"`
	Disposition    string `json:"disposition"`
}

type ReqStorageFilePrepareUpload struct {
	OriginalName string `json:"originalName" binding:"required,max=255" binding_msg:"required=文件名不能为空,max=文件名最多255位"`
	ContentType  string `json:"contentType" binding:"max=128" binding_msg:"max=内容类型最多128位"`
	Size         int64  `json:"size" binding:"required" binding_msg:"required=文件大小不能为空"`
	BizType      string `json:"bizType" binding:"max=64" binding_msg:"max=业务类型最多64位"`
	BizID        string `json:"bizID" binding:"max=128" binding_msg:"max=业务ID最多128位"`
	Metadata     string `json:"metadata"`
	Remark       string `json:"remark" binding:"max=255" binding_msg:"max=备注最多255位"`
}

// @Summary 上传文件
// @Tags StorageFile
// @Accept mpfd
// @Produce json
// @Param file formData file true "上传文件"
// @Param bizType formData string false "业务类型"
// @Param bizID formData string false "业务ID"
// @Param metadata formData string false "JSON 元数据"
// @Param remark formData string false "备注"
// @Success 200 {object} res.Response{data=models.FileAsset} "成功"
// @Router /api/storage/file/upload [post]
func (h *StorageFileHandler) Upload(ctx *handler.Ctx) (*models.FileAsset, error) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return nil, res.FailMsg("请选择上传文件")
	}
	if fileHeader.Size <= 0 {
		return nil, res.FailMsg("上传文件不能为空")
	}
	if h.MaxUploadBytes > 0 && fileHeader.Size > h.MaxUploadBytes {
		return nil, res.FailMsg("上传文件超过大小限制")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, res.FailDefault
	}
	defer file.Close()

	return h.Storage.Upload(ctx.Context(), appsvc.UploadInput{
		OperatorID:   ctx.SessionInfo.Id,
		Reader:       file,
		Size:         fileHeader.Size,
		OriginalName: fileHeader.Filename,
		ContentType:  fileHeader.Header.Get("Content-Type"),
		BizType:      ctx.FormValue("bizType"),
		BizID:        ctx.FormValue("bizID"),
		Metadata:     ctx.FormValue("metadata"),
		Remark:       ctx.FormValue("remark"),
	})
}

// @Summary 准备文件直传
// @Tags StorageFile
// @Accept json
// @Produce json
// @Param req body ReqStorageFilePrepareUpload true "上传参数"
// @Success 200 {object} res.Response{data=appsvc.PrepareDirectUploadResult} "成功"
// @Router /api/storage/file/prepareUpload [post]
func (h *StorageFileHandler) PrepareUpload(ctx *handler.Ctx, req *ReqStorageFilePrepareUpload) (*appsvc.PrepareDirectUploadResult, error) {
	return h.Storage.PrepareDirectUpload(ctx.Context(), appsvc.PrepareDirectUploadInput{
		OperatorID:   ctx.SessionInfo.Id,
		OriginalName: req.OriginalName,
		ContentType:  req.ContentType,
		Size:         req.Size,
		BizType:      req.BizType,
		BizID:        req.BizID,
		Metadata:     req.Metadata,
		Remark:       req.Remark,
	})
}

// @Summary 获取文件详情
// @Tags StorageFile
// @Accept json
// @Produce json
// @Param req body ReqStorageFileID true "文件ID"
// @Success 200 {object} res.Response{data=models.FileAsset} "成功"
// @Router /api/storage/file/detail [post]
func (h *StorageFileHandler) Detail(ctx *handler.Ctx, req *ReqStorageFileID) (*models.FileAsset, error) {
	return h.Storage.Detail(ctx.Context(), req.ID)
}

// @Summary 获取文件临时下载链接
// @Tags StorageFile
// @Accept json
// @Produce json
// @Param req body ReqStorageFilePresigned true "签名参数"
// @Success 200 {object} res.Response{data=appsvc.PresignedURLResult} "成功"
// @Router /api/storage/file/presigned [post]
func (h *StorageFileHandler) Presigned(ctx *handler.Ctx, req *ReqStorageFilePresigned) (*appsvc.PresignedURLResult, error) {
	disposition := strings.TrimSpace(req.Disposition)
	if disposition != "" && disposition != "inline" && disposition != "attachment" {
		return nil, res.FailMsg("disposition 参数错误")
	}
	return h.Storage.PresignedURL(ctx.Context(), appsvc.PresignedURLInput{
		ID:             req.ID,
		ExpiresSeconds: req.ExpiresSeconds,
		Disposition:    disposition,
	})
}

// @Summary 完成文件直传
// @Tags StorageFile
// @Accept json
// @Produce json
// @Param req body ReqStorageFileID true "文件ID"
// @Success 200 {object} res.Response{data=models.FileAsset} "成功"
// @Router /api/storage/file/completeUpload [post]
func (h *StorageFileHandler) CompleteUpload(ctx *handler.Ctx, req *ReqStorageFileID) (*models.FileAsset, error) {
	return h.Storage.CompleteDirectUpload(ctx.Context(), appsvc.CompleteDirectUploadInput{
		ID:         req.ID,
		OperatorID: ctx.SessionInfo.Id,
	})
}

// @Summary 删除文件
// @Tags StorageFile
// @Accept json
// @Produce json
// @Param req body ReqStorageFileID true "文件ID"
// @Success 200 {object} res.Response "成功"
// @Router /api/storage/file/del [post]
func (h *StorageFileHandler) Del(ctx *handler.Ctx, req *ReqStorageFileID) error {
	return h.Storage.Delete(ctx.Context(), req.ID, ctx.SessionInfo.Id)
}
