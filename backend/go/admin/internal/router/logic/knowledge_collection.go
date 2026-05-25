package logic

import (
	"errors"

	"admin/internal/appsvc"
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"

	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc"
	"orm-crud/gormc/mixin"

	"go.uber.org/zap"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type KnowledgeCollectionHandler struct {
	Q       *query.Query
	indexer appsvc.KnowledgeDocumentIndexer
}

func NewKnowledgeCollectionHandler(q *query.Query, indexer appsvc.KnowledgeDocumentIndexer) *KnowledgeCollectionHandler {
	return &KnowledgeCollectionHandler{Q: q, indexer: indexer}
}

type RespKnowledgeCollection struct {
	models.KnowledgeCollection
	CanWrite  bool `json:"canWrite"`
	CanDelete bool `json:"canDelete"`
}

type ReqKnowledgeCollectionCreate struct {
	CollectionName string `json:"collectionName" change:"集合名称" binding:"required,max=128" binding_msg:"required=集合名称不能为空,max=集合名称最多128位"`
	DisplayName    string `json:"displayName" change:"显示名称" binding:"required,max=255" binding_msg:"required=显示名称不能为空,max=显示名称最多255位"`
	Description    string `json:"description" change:"描述" binding:"max=512" binding_msg:"max=描述最多512位"`
	MetricType     string `json:"metricType" change:"度量类型" binding:"max=32" binding_msg:"max=度量类型最多32位"`
	IndexType      string `json:"indexType" change:"索引类型" binding:"max=32" binding_msg:"max=索引类型最多32位"`
	IsEnabled      bool   `json:"isEnabled" change:"启用状态"`
	Remark         string `json:"remark" change:"备注" binding:"max=255" binding_msg:"max=备注最多255位"`
}

type ReqKnowledgeCollectionUpdate struct {
	ID             uint64  `json:"id" binding:"required" binding_msg:"required=请求错误"`
	CollectionName *string `json:"collectionName" change:"集合名称" binding:"omitempty,max=128" binding_msg:"max=集合名称最多128位"`
	DisplayName    *string `json:"displayName" change:"显示名称" binding:"omitempty,max=255" binding_msg:"max=显示名称最多255位"`
	Description    *string `json:"description" change:"描述" binding:"omitempty,max=512" binding_msg:"max=描述最多512位"`
	MetricType     *string `json:"metricType" change:"度量类型" binding:"omitempty,max=32" binding_msg:"max=度量类型最多32位"`
	IndexType      *string `json:"indexType" change:"索引类型" binding:"omitempty,max=32" binding_msg:"max=索引类型最多32位"`
	IsEnabled      *bool   `json:"isEnabled" change:"启用状态"`
	Remark         *string `json:"remark" change:"备注" binding:"omitempty,max=255" binding_msg:"max=备注最多255位"`
}

type ReqKnowledgeCollectionID struct {
	ID uint64 `json:"id" binding:"required" binding_msg:"required=请求错误"`
}

// @Summary 获取知识库集合分页列表
// @Tags KnowledgeCollection
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[RespKnowledgeCollection]} "成功"
// @Router /api/knowledge/collection/list [post]
func (h *KnowledgeCollectionHandler) List(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[RespKnowledgeCollection], error) {
	if req.GetOrderBy() == "" {
		orderBy := "id desc"
		req.OrderBy = &orderBy
	}

	pagination, err := h.Q.KnowledgeCollection.PageWithPaging(req)
	if err != nil {
		ctx.L().Error("query knowledge collection list fail", zap.Error(err))
		return nil, res.FailDefault
	}

	items := make([]*RespKnowledgeCollection, 0, len(pagination.Items))
	for _, item := range pagination.Items {
		if item == nil {
			continue
		}
		items = append(items, &RespKnowledgeCollection{
			KnowledgeCollection: *item,
			CanWrite:            true,
			CanDelete:           true,
		})
	}
	return &gormc.PagingResult[RespKnowledgeCollection]{Items: items, Total: pagination.Total}, nil
}

// @Summary 获取知识库集合详情
// @Tags KnowledgeCollection
// @Accept json
// @Produce json
// @Param req body ReqKnowledgeCollectionID true "集合ID"
// @Success 200 {object} res.Response{data=models.KnowledgeCollection} "成功"
// @Router /api/knowledge/collection/detail [post]
func (h *KnowledgeCollectionHandler) Detail(ctx *handler.Ctx, req *ReqKnowledgeCollectionID) (*models.KnowledgeCollection, error) {
	collection := h.Q.KnowledgeCollection
	item, err := collection.Where(collection.ID.Eq(req.ID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, res.FailMsg("集合不存在")
		}
		ctx.L().Error("query knowledge collection detail fail", zap.Error(err), zap.Uint64("id", req.ID))
		return nil, res.FailDefault
	}
	return item, nil
}

// @Summary 创建知识库集合
// @Tags KnowledgeCollection
// @Accept json
// @Produce json
// @Param req body ReqKnowledgeCollectionCreate true "创建参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/knowledge/collection/create [post]
func (h *KnowledgeCollectionHandler) Create(ctx *handler.Ctx, req *ReqKnowledgeCollectionCreate) error {
	operationID := ctx.SessionInfo.Id

	metricType := req.MetricType
	if metricType == "" {
		metricType = "COSINE"
	}
	indexType := req.IndexType
	if indexType == "" {
		indexType = "auto"
	}

	err := h.Q.KnowledgeCollection.Create(&models.KnowledgeCollection{
		OperatorID: mixin.OperatorID{
			CreatedBy: mixin.CreatedBy{CreatedBy: operationID},
			UpdatedBy: mixin.UpdatedBy{UpdatedBy: operationID},
		},
		IsEnabled:      mixin.IsEnabled{IsEnabled: req.IsEnabled},
		Remark:         mixin.Remark{Remark: req.Remark},
		CollectionName: req.CollectionName,
		DisplayName:    req.DisplayName,
		MetricType:     metricType,
		IndexType:      indexType,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return res.FailMsg("集合名称已存在")
		}
		ctx.L().Error("create knowledge collection fail", zap.Error(err))
		return res.FailDefault
	}

	if h.indexer != nil {
		collection := &models.KnowledgeCollection{
			CollectionName: req.CollectionName,
			MetricType:     metricType,
			IndexType:      indexType,
		}
		if err := h.indexer.EnsureCollection(ctx, collection); err != nil {
			ctx.L().Error("ensure milvus collection fail", zap.Error(err), zap.String("collectionName", req.CollectionName))
		}
	}

	return nil
}

// @Summary 更新知识库集合
// @Tags KnowledgeCollection
// @Accept json
// @Produce json
// @Param req body ReqKnowledgeCollectionUpdate true "更新参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/knowledge/collection/update [post]
func (h *KnowledgeCollectionHandler) Update(ctx *handler.Ctx, req *ReqKnowledgeCollectionUpdate) error {
	operationID := ctx.SessionInfo.Id
	collection := h.Q.KnowledgeCollection

	exprs := []field.AssignExpr{collection.UpdatedBy.Value(operationID)}
	query.ExprAppendSelf(&exprs, req.CollectionName, collection.CollectionName.Value)
	query.ExprAppendSelf(&exprs, req.DisplayName, collection.DisplayName.Value)
	query.ExprAppendSelf(&exprs, req.MetricType, collection.MetricType.Value)
	query.ExprAppendSelf(&exprs, req.IndexType, collection.IndexType.Value)
	query.ExprAppendSelf(&exprs, req.IsEnabled, collection.IsEnabled.Value)
	query.ExprAppendSelf(&exprs, req.Remark, collection.Remark.Value)

	info, err := collection.Where(collection.ID.Eq(req.ID)).UpdateSimple(exprs...)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return res.FailMsg("集合名称已存在")
		}
		ctx.L().Error("update knowledge collection fail", zap.Error(err), zap.Uint64("id", req.ID))
		return res.FailDefault
	}
	if info.RowsAffected == 0 {
		return res.FailMsg("集合不存在")
	}
	return nil
}

// @Summary 删除知识库集合
// @Tags KnowledgeCollection
// @Accept json
// @Produce json
// @Param req body ReqKnowledgeCollectionID true "集合ID"
// @Success 200 {object} res.Response "成功"
// @Router /api/knowledge/collection/del [post]
func (h *KnowledgeCollectionHandler) Del(ctx *handler.Ctx, req *ReqKnowledgeCollectionID) error {
	collectionModel := h.Q.KnowledgeCollection

	record, err := collectionModel.Where(collectionModel.ID.Eq(req.ID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res.FailMsg("集合不存在")
		}
		ctx.L().Error("query knowledge collection for delete fail", zap.Error(err), zap.Uint64("id", req.ID))
		return res.FailDefault
	}

	if h.indexer != nil {
		if err := h.indexer.DropCollection(ctx, record.CollectionName); err != nil {
			ctx.L().Error("drop milvus collection fail", zap.Error(err), zap.String("collectionName", record.CollectionName))
		}
	}

	info, err := collectionModel.Where(collectionModel.ID.Eq(req.ID)).Delete()
	if err != nil {
		ctx.L().Error("delete knowledge collection fail", zap.Error(err), zap.Uint64("id", req.ID))
		return res.FailDefault
	}
	if info.RowsAffected == 0 {
		return res.FailMsg("集合不存在")
	}
	return nil
}
