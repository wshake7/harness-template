package appsvc

import (
	"admin/internal/config"
	"admin/internal/fiberc/res"
	"admin/internal/services/objectstore"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"path"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"orm-crud/gormc/mixin"
)

type storageRuntimeConfig struct {
	engineName              string
	bucket                  string
	objectKeyPrefix         string
	maxUploadBytes          int64
	presignedExpiresSeconds int
}

type fileStorageImpl struct {
	q           *query.Query
	engine      objectstore.Engine
	conf        storageRuntimeConfig
	now         func() time.Time
	newObjectID func() string
}

func NewFileStorage(q *query.Query, conf config.StorageConfig, engine objectstore.Engine) FileStorage {
	if q == nil {
		q = query.Q
	}
	if engine == nil {
		engine = objectstore.Current()
	}
	return &fileStorageImpl{
		q:      q,
		engine: engine,
		conf: storageRuntimeConfig{
			engineName:              strings.ToLower(strings.TrimSpace(conf.Engine)),
			bucket:                  conf.MinIO.Bucket,
			objectKeyPrefix:         conf.ObjectKeyPrefix,
			maxUploadBytes:          conf.MaxUploadBytes,
			presignedExpiresSeconds: conf.PresignedExpiresSeconds,
		},
		now: time.Now,
		newObjectID: func() string {
			return uuid.NewString()
		},
	}
}

func (s *fileStorageImpl) Upload(ctx context.Context, input UploadInput) (*models.FileAsset, error) {
	if input.Reader == nil || input.Size <= 0 {
		return nil, res.FailMsg("上传文件不能为空")
	}
	if s.conf.maxUploadBytes > 0 && input.Size > s.conf.maxUploadBytes {
		return nil, res.FailMsg("上传文件超过大小限制")
	}

	metadata, err := parseStorageMetadata(input.Metadata)
	if err != nil {
		return nil, res.FailMsg("元数据格式错误")
	}

	contentType := strings.TrimSpace(input.ContentType)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	objectKey := objectstore.BuildObjectKey(s.conf.objectKeyPrefix, input.OriginalName, s.now(), s.newObjectID())
	hash := sha256.New()
	reader := io.TeeReader(input.Reader, hash)
	if _, err = s.engine.PutObject(ctx, objectstore.PutInput{
		Bucket:      s.conf.bucket,
		ObjectKey:   objectKey,
		Reader:      reader,
		Size:        input.Size,
		ContentType: contentType,
	}); err != nil {
		return nil, err
	}

	asset := &models.FileAsset{
		OperatorID: mixin.OperatorID{
			CreatedBy: mixin.CreatedBy{CreatedBy: input.OperatorID},
			UpdatedBy: mixin.UpdatedBy{UpdatedBy: input.OperatorID},
		},
		Remark:       mixin.Remark{Remark: input.Remark},
		Engine:       engineNameOrDefault(s.conf.engineName, s.engine),
		Bucket:       s.conf.bucket,
		ObjectKey:    objectKey,
		OriginalName: path.Base(strings.ReplaceAll(input.OriginalName, "\\", "/")),
		ContentType:  contentType,
		Extension:    strings.ToLower(path.Ext(path.Base(strings.ReplaceAll(input.OriginalName, "\\", "/")))),
		Size:         input.Size,
		SHA256:       hex.EncodeToString(hash.Sum(nil)),
		BizType:      strings.TrimSpace(input.BizType),
		BizID:        strings.TrimSpace(input.BizID),
		Metadata:     metadata,
		Status:       models.FileAssetStatusActive,
	}
	if asset.Extension == "." {
		asset.Extension = ""
	}

	if err = s.q.FileAsset.Create(asset); err != nil {
		if removeErr := s.engine.RemoveObject(ctx, s.conf.bucket, objectKey); removeErr != nil && !errors.Is(removeErr, objectstore.ErrObjectNotFound) {
			zap.L().Warn("compensate uploaded object failed", zap.Error(removeErr), zap.String("bucket", s.conf.bucket), zap.String("objectKey", objectKey))
		}
		return nil, err
	}

	return asset, nil
}

func (s *fileStorageImpl) Detail(ctx context.Context, id uint64) (*models.FileAsset, error) {
	item, err := s.q.FileAsset.Where(s.q.FileAsset.ID.Eq(id)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, res.FailMsg("文件不存在")
		}
		return nil, err
	}
	return item, nil
}

func (s *fileStorageImpl) PresignedURL(ctx context.Context, input PresignedURLInput) (*PresignedURLResult, error) {
	if input.Disposition != "" && input.Disposition != "inline" && input.Disposition != "attachment" {
		return nil, res.FailMsg("disposition 参数错误")
	}

	item, err := s.Detail(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	expiry := objectstore.ClampPresignedExpiry(input.ExpiresSeconds, s.conf.presignedExpiresSeconds)
	url, expiresAt, err := s.engine.PresignedGetObject(ctx, objectstore.PresignInput{
		Bucket:      item.Bucket,
		ObjectKey:   item.ObjectKey,
		Expiry:      expiry,
		Disposition: input.Disposition,
	})
	if err != nil {
		return nil, err
	}
	return &PresignedURLResult{URL: url, ExpiresAt: expiresAt}, nil
}

func (s *fileStorageImpl) Delete(ctx context.Context, id uint64, operatorID uint64) error {
	item, err := s.Detail(ctx, id)
	if err != nil {
		return err
	}

	if err = s.engine.RemoveObject(ctx, item.Bucket, item.ObjectKey); err != nil && !errors.Is(err, objectstore.ErrObjectNotFound) {
		return err
	}

	item.Status = models.FileAssetStatusDeleted
	if _, err = s.q.FileAsset.Where(s.q.FileAsset.ID.Eq(id)).UpdateSimple(
		s.q.FileAsset.Status.Value(models.FileAssetStatusDeleted),
		s.q.FileAsset.UpdatedBy.Value(operatorID),
		s.q.FileAsset.DeletedBy.Value(operatorID),
	); err != nil {
		return err
	}
	_, err = s.q.FileAsset.Where(s.q.FileAsset.ID.Eq(id)).Delete(item)
	return err
}

func parseStorageMetadata(raw string) (datatypes.JSONMap, error) {
	if strings.TrimSpace(raw) == "" {
		return datatypes.JSONMap{}, nil
	}

	var metadata datatypes.JSONMap
	if err := sonic.UnmarshalString(raw, &metadata); err != nil {
		return nil, err
	}
	if metadata == nil {
		return datatypes.JSONMap{}, nil
	}
	return metadata, nil
}

func engineNameOrDefault(configName string, engine objectstore.Engine) string {
	if configName != "" {
		return configName
	}
	if engine != nil && engine.Name() != "" {
		return engine.Name()
	}
	return "minio"
}
