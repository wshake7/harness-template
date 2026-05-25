package logic

import (
	"admin/internal/appsvc"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"admin/internal/services/temporaljob"
	knowledgedocument "admin/internal/workflows/knowledge_document"

	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc"
	"orm-crud/gormc/mixin"

	"github.com/bytedance/sonic"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type KnowledgeDocumentHandler struct {
	Q                 *query.Query
	Indexer           appsvc.KnowledgeDocumentIndexer
	Temporal          appsvc.TemporalService
	TemporalTaskQueue string
}

func NewKnowledgeDocumentHandler(q *query.Query, indexer appsvc.KnowledgeDocumentIndexer, temporal appsvc.TemporalService, temporalTaskQueue string) *KnowledgeDocumentHandler {
	if indexer == nil {
		indexer = noopKnowledgeDocumentIndexer{}
	}
	return &KnowledgeDocumentHandler{Q: q, Indexer: indexer, Temporal: temporal, TemporalTaskQueue: temporalTaskQueue}
}

type RespKnowledgeDocument struct {
	models.KnowledgeDocument
	CanWrite  bool `json:"canWrite"`
	CanDelete bool `json:"canDelete"`
}

type ReqKnowledgeDocumentCreate struct {
	CollectionID uint64 `json:"collectionID" change:"所属集合" binding:"required" binding_msg:"required=请选择所属集合"`
	DocumentID   string `json:"documentID" change:"文档ID" binding:"required,max=255" binding_msg:"required=文档ID不能为空,max=文档ID最多255位"`
	Title        string `json:"title" change:"标题" binding:"required,max=512" binding_msg:"required=标题不能为空,max=标题最多512位"`
	Content      string `json:"content" change:"内容" binding:"required" binding_msg:"required=内容不能为空"`
	ContentType  string `json:"contentType" change:"内容类型" binding:"max=64" binding_msg:"max=内容类型最多64位"`
	Source       string `json:"source" change:"来源" binding:"max=512" binding_msg:"max=来源最多512位"`
	ChunkIndex   int    `json:"chunkIndex" change:"分块索引"`
	TotalChunks  int    `json:"totalChunks" change:"总分块数"`
	Metadata     string `json:"metadata" change:"元数据"`
	IsEnabled    bool   `json:"isEnabled" change:"启用状态"`
	Remark       string `json:"remark" change:"备注" binding:"max=255" binding_msg:"max=备注最多255位"`
}

type ReqKnowledgeDocumentUpdate struct {
	ID           uint64  `json:"id" binding:"required" binding_msg:"required=请求错误"`
	CollectionID *uint64 `json:"collectionID" change:"所属集合"`
	DocumentID   *string `json:"documentID" change:"文档ID" binding:"omitempty,max=255" binding_msg:"max=文档ID最多255位"`
	Title        *string `json:"title" change:"标题" binding:"omitempty,max=512" binding_msg:"max=标题最多512位"`
	Content      *string `json:"content" change:"内容"`
	ContentType  *string `json:"contentType" change:"内容类型" binding:"omitempty,max=64" binding_msg:"max=内容类型最多64位"`
	Source       *string `json:"source" change:"来源" binding:"omitempty,max=512" binding_msg:"max=来源最多512位"`
	ChunkIndex   *int    `json:"chunkIndex" change:"分块索引"`
	TotalChunks  *int    `json:"totalChunks" change:"总分块数"`
	Metadata     *string `json:"metadata" change:"元数据"`
	VectorStatus *string `json:"vectorStatus" change:"向量状态" binding:"omitempty,max=32" binding_msg:"max=向量状态最多32位"`
	IsEnabled    *bool   `json:"isEnabled" change:"启用状态"`
	Remark       *string `json:"remark" change:"备注" binding:"omitempty,max=255" binding_msg:"max=备注最多255位"`
}

type ReqKnowledgeDocumentID struct {
	ID uint64 `json:"id" binding:"required" binding_msg:"required=请求错误"`
}

type ReqKnowledgeDocumentImportFile struct {
	CollectionID uint64 `json:"collectionID" binding:"required" binding_msg:"required=请选择所属集合"`
	FileAssetID  uint64 `json:"fileAssetID" binding:"required" binding_msg:"required=请选择文件"`
	Title        string `json:"title" binding:"max=512" binding_msg:"max=标题最多512位"`
	ContentType  string `json:"contentType" binding:"max=64" binding_msg:"max=内容类型最多64位"`
	Remark       string `json:"remark" binding:"max=255" binding_msg:"max=备注最多255位"`
}

type ReqKnowledgeDocumentListByCollection struct {
	CollectionID uint64  `json:"collectionID" binding:"required" binding_msg:"required=请选择集合"`
	Page         *uint32 `json:"page,omitempty"`
	PageSize     *uint32 `json:"pageSize,omitempty"`
	OrderBy      *string `json:"orderBy,omitempty"`
	Query        *string `json:"query,omitempty"`
}

// @Summary 获取知识库文档分页列表
// @Tags KnowledgeDocument
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[RespKnowledgeDocument]} "成功"
// @Router /api/knowledge/document/list [post]
func (h *KnowledgeDocumentHandler) List(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[RespKnowledgeDocument], error) {
	if req.GetOrderBy() == "" {
		orderBy := "id desc"
		req.OrderBy = &orderBy
	}

	pagination, err := h.Q.KnowledgeDocument.PageWithPaging(req)
	if err != nil {
		ctx.L().Error("query knowledge document list fail", zap.Error(err))
		return nil, res.FailDefault
	}

	items := make([]*RespKnowledgeDocument, 0, len(pagination.Items))
	for _, item := range pagination.Items {
		if item == nil {
			continue
		}
		items = append(items, &RespKnowledgeDocument{
			KnowledgeDocument: *item,
			CanWrite:          true,
			CanDelete:         true,
		})
	}
	return &gormc.PagingResult[RespKnowledgeDocument]{Items: items, Total: pagination.Total}, nil
}

// @Summary 按集合获取知识库文档分页列表
// @Tags KnowledgeDocument
// @Accept json
// @Produce json
// @Param req body ReqKnowledgeDocumentListByCollection true "查询参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[RespKnowledgeDocument]} "成功"
// @Router /api/knowledge/document/listByCollection [post]
func (h *KnowledgeDocumentHandler) ListByCollection(ctx *handler.Ctx, req *ReqKnowledgeDocumentListByCollection) (*gormc.PagingResult[RespKnowledgeDocument], error) {
	pagingReq := &v1.PagingRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		OrderBy:  req.OrderBy,
	}
	// 先设置 FilteringType（query string 通过 oneof 赋值）
	if req.Query != nil && *req.Query != "" {
		q := *req.Query
		pagingReq.FilteringType = &v1.PagingRequest_Query{Query: q}
	}

	if pagingReq.GetOrderBy() == "" {
		orderBy := "id desc"
		pagingReq.OrderBy = &orderBy
	}

	// 把 collectionID 条件合并到 query filter
	collectionFilter := map[string]any{"collectionID": req.CollectionID}
	current := strings.TrimSpace(pagingReq.GetQuery())
	if current == "" {
		data, _ := sonic.MarshalString(collectionFilter)
		pagingReq.FilteringType = &v1.PagingRequest_Query{Query: data}
	} else {
		extraData, _ := sonic.MarshalString(collectionFilter)
		pagingReq.FilteringType = &v1.PagingRequest_Query{Query: fmt.Sprintf(`{"$and":[%s,%s]}`, current, extraData)}
	}

	pagination, err := h.Q.KnowledgeDocument.PageWithPaging(pagingReq)
	if err != nil {
		ctx.L().Error("query knowledge document by collection fail", zap.Error(err), zap.Uint64("collectionID", req.CollectionID))
		return nil, res.FailDefault
	}

	items := make([]*RespKnowledgeDocument, 0, len(pagination.Items))
	for _, item := range pagination.Items {
		if item == nil {
			continue
		}
		items = append(items, &RespKnowledgeDocument{
			KnowledgeDocument: *item,
			CanWrite:          true,
			CanDelete:         true,
		})
	}
	return &gormc.PagingResult[RespKnowledgeDocument]{Items: items, Total: pagination.Total}, nil
}

// @Summary 获取知识库文档详情
// @Tags KnowledgeDocument
// @Accept json
// @Produce json
// @Param req body ReqKnowledgeDocumentID true "文档ID"
// @Success 200 {object} res.Response{data=models.KnowledgeDocument} "成功"
// @Router /api/knowledge/document/detail [post]
func (h *KnowledgeDocumentHandler) Detail(ctx *handler.Ctx, req *ReqKnowledgeDocumentID) (*models.KnowledgeDocument, error) {
	doc := h.Q.KnowledgeDocument
	item, err := doc.Where(doc.ID.Eq(req.ID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, res.FailMsg("文档不存在")
		}
		ctx.L().Error("query knowledge document detail fail", zap.Error(err), zap.Uint64("id", req.ID))
		return nil, res.FailDefault
	}
	return item, nil
}

// @Summary 创建知识库文档
// @Tags KnowledgeDocument
// @Accept json
// @Produce json
// @Param req body ReqKnowledgeDocumentCreate true "创建参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/knowledge/document/create [post]
func (h *KnowledgeDocumentHandler) Create(ctx *handler.Ctx, req *ReqKnowledgeDocumentCreate) error {
	operationID := ctx.SessionInfo.Id

	contentType := req.ContentType
	if contentType == "" {
		contentType = "text"
	}

	var metadata datatypes.JSONMap
	if req.Metadata != "" {
		if err := sonic.UnmarshalString(req.Metadata, &metadata); err != nil {
			return res.FailMsg("元数据格式错误")
		}
	}
	if metadata == nil {
		metadata = datatypes.JSONMap{}
	}

	doc := &models.KnowledgeDocument{
		OperatorID: mixin.OperatorID{
			CreatedBy: mixin.CreatedBy{CreatedBy: operationID},
			UpdatedBy: mixin.UpdatedBy{UpdatedBy: operationID},
		},
		IsEnabled:    mixin.IsEnabled{IsEnabled: req.IsEnabled},
		Remark:       mixin.Remark{Remark: req.Remark},
		CollectionID: req.CollectionID,
		DocumentID:   req.DocumentID,
		Title:        req.Title,
		Content:      req.Content,
		ContentType:  contentType,
		Source:       req.Source,
		ChunkIndex:   req.ChunkIndex,
		TotalChunks:  req.TotalChunks,
		VectorStatus: "pending",
		Metadata:     metadata,
	}

	err := h.Q.KnowledgeDocument.Create(doc)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return res.FailMsg("文档ID已存在")
		}
		ctx.L().Error("create knowledge document fail", zap.Error(err))
		return res.FailDefault
	}
	return h.Indexer.IndexDocument(ctx.Context(), doc.ID)
}

// @Summary 更新知识库文档
// @Tags KnowledgeDocument
// @Accept json
// @Produce json
// @Param req body ReqKnowledgeDocumentUpdate true "更新参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/knowledge/document/update [post]
func (h *KnowledgeDocumentHandler) Update(ctx *handler.Ctx, req *ReqKnowledgeDocumentUpdate) error {
	operationID := ctx.SessionInfo.Id
	doc := h.Q.KnowledgeDocument
	before, err := doc.Where(doc.ID.Eq(req.ID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res.FailMsg("文档不存在")
		}
		ctx.L().Error("query knowledge document before update fail", zap.Error(err), zap.Uint64("id", req.ID))
		return res.FailDefault
	}

	exprs := []field.AssignExpr{doc.UpdatedBy.Value(operationID)}
	query.ExprAppendSelf(&exprs, req.CollectionID, doc.CollectionID.Value)
	query.ExprAppendSelf(&exprs, req.DocumentID, doc.DocumentID.Value)
	query.ExprAppendSelf(&exprs, req.Title, doc.Title.Value)
	query.ExprAppendSelf(&exprs, req.Content, doc.Content.Value)
	query.ExprAppendSelf(&exprs, req.ContentType, doc.ContentType.Value)
	query.ExprAppendSelf(&exprs, req.Source, doc.Source.Value)
	query.ExprAppendSelf(&exprs, req.ChunkIndex, doc.ChunkIndex.Value)
	query.ExprAppendSelf(&exprs, req.TotalChunks, doc.TotalChunks.Value)
	query.ExprAppendSelf(&exprs, req.VectorStatus, doc.VectorStatus.Value)
	query.ExprAppendSelf(&exprs, req.IsEnabled, doc.IsEnabled.Value)
	query.ExprAppendSelf(&exprs, req.Remark, doc.Remark.Value)

	if req.Metadata != nil {
		var metadata datatypes.JSONMap
		if *req.Metadata != "" {
			if err := sonic.UnmarshalString(*req.Metadata, &metadata); err != nil {
				return res.FailMsg("元数据格式错误")
			}
		}
		if metadata == nil {
			metadata = datatypes.JSONMap{}
		}
		exprs = append(exprs, doc.Metadata.Value(metadata))
	}

	shouldReindex := shouldReindexDocument(before, req)

	info, err := doc.Where(doc.ID.Eq(req.ID)).UpdateSimple(exprs...)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return res.FailMsg("文档ID已存在")
		}
		ctx.L().Error("update knowledge document fail", zap.Error(err), zap.Uint64("id", req.ID))
		return res.FailDefault
	}
	if info.RowsAffected == 0 {
		return res.FailMsg("文档不存在")
	}
	if shouldReindex {
		return h.Indexer.IndexDocument(ctx.Context(), req.ID)
	}
	return nil
}

// @Summary 删除知识库文档
// @Tags KnowledgeDocument
// @Accept json
// @Produce json
// @Param req body ReqKnowledgeDocumentID true "文档ID"
// @Success 200 {object} res.Response "成功"
// @Router /api/knowledge/document/del [post]
func (h *KnowledgeDocumentHandler) Del(ctx *handler.Ctx, req *ReqKnowledgeDocumentID) error {
	doc := h.Q.KnowledgeDocument
	if err := h.Indexer.DeleteDocumentVectors(ctx.Context(), req.ID); err != nil {
		return err
	}
	info, err := doc.Where(doc.ID.Eq(req.ID)).Delete()
	if err != nil {
		ctx.L().Error("delete knowledge document fail", zap.Error(err), zap.Uint64("id", req.ID))
		return res.FailDefault
	}
	if info.RowsAffected == 0 {
		return res.FailMsg("文档不存在")
	}
	return nil
}

// @Summary 导入文件创建知识库文档
// @Tags KnowledgeDocument
// @Accept json
// @Produce json
// @Param req body ReqKnowledgeDocumentImportFile true "导入参数"
// @Success 200 {object} res.Response{data=models.KnowledgeDocument} "成功"
// @Router /api/knowledge/document/importFile [post]
func (h *KnowledgeDocumentHandler) ImportFile(ctx *handler.Ctx, req *ReqKnowledgeDocumentImportFile) (*models.KnowledgeDocument, error) {
	doc, err := h.Indexer.ImportFile(ctx.Context(), appsvc.ImportKnowledgeFileInput{
		CollectionID: req.CollectionID,
		FileAssetID:  req.FileAssetID,
		Title:        req.Title,
		ContentType:  req.ContentType,
		Remark:       req.Remark,
		OperatorID:   ctx.SessionInfo.Id,
	})
	if err != nil {
		return nil, err
	}
	if err = h.enqueueIndexDocument(ctx, doc.ID); err != nil {
		return nil, err
	}
	return doc, nil
}

const knowledgeDocumentIndexJobCodePrefix = "knowledge-document-index"

func KnowledgeDocumentIndexJobCode(documentID uint64) string {
	return fmt.Sprintf("%s-%d", knowledgeDocumentIndexJobCodePrefix, documentID)
}

func (h *KnowledgeDocumentHandler) enqueueIndexDocument(ctx *handler.Ctx, documentID uint64) error {
	if h.Temporal == nil || !h.Temporal.IsConnected() || strings.TrimSpace(h.TemporalTaskQueue) == "" {
		ctx.L().Warn("Temporal not connected, falling back to direct indexing",
			zap.Uint64("documentID", documentID),
			zap.String("taskQueue", h.TemporalTaskQueue))
		return h.runDirectIndexWithExecutionRecord(ctx.Context(), documentID)
	}

	workflowID := fmt.Sprintf("%s-dispatch-%d", KnowledgeDocumentIndexJobCode(documentID), time.Now().UnixNano())
	_, err := h.Temporal.ExecuteWorkflow(ctx.Context(), client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: h.TemporalTaskQueue,
	}, temporaljob.DispatchWorkflowName, temporaljob.DispatchInput{
		JobCode:          KnowledgeDocumentIndexJobCode(documentID),
		WorkflowType:     knowledgedocument.IndexWorkflowName,
		TaskQueue:        h.TemporalTaskQueue,
		WorkflowIDPrefix: KnowledgeDocumentIndexJobCode(documentID),
		Input: knowledgedocument.IndexInput{
			DocumentID: documentID,
		},
	})
	if err != nil {
		ctx.L().Error("dispatch knowledge document index workflow fail, falling back to direct indexing", zap.Error(err), zap.Uint64("documentID", documentID), zap.String("workflowID", workflowID))
		return h.runDirectIndexWithExecutionRecord(ctx.Context(), documentID)
	}
	return nil
}

func (h *KnowledgeDocumentHandler) runDirectIndexWithExecutionRecord(ctx context.Context, documentID uint64) error {
	jobCode := KnowledgeDocumentIndexJobCode(documentID)
	now := time.Now()
	execution := &models.JobExecution{
		JobCode:     jobCode,
		TriggerTime: now,
		StartTime:   &now,
		Status:      models.JobExecutionStatusRunning,
		RetryCount:  0,
	}
	if err := h.Q.JobExecution.WithContext(ctx).Create(execution); err != nil {
		return err
	}

	indexErr := h.Indexer.IndexDocument(ctx, documentID)

	if indexErr != nil {
		_, _ = h.Q.JobExecution.WithContext(ctx).
			Where(h.Q.JobExecution.ID.Eq(execution.ID)).
			UpdateSimple(
				h.Q.JobExecution.Status.Value(models.JobExecutionStatusFailed),
				h.Q.JobExecution.ErrorMessage.Value(sanitizeError(indexErr)),
				h.Q.JobExecution.EndTime.Value(time.Now()),
			)
		return indexErr
	}

	_, _ = h.Q.JobExecution.WithContext(ctx).
		Where(h.Q.JobExecution.ID.Eq(execution.ID)).
		UpdateSimple(
			h.Q.JobExecution.Status.Value(models.JobExecutionStatusSuccess),
			h.Q.JobExecution.EndTime.Value(time.Now()),
		)
	return nil
}

func sanitizeError(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	if len(msg) > 512 {
		return msg[:512]
	}
	return msg
}

func shouldReindexDocument(before *models.KnowledgeDocument, req *ReqKnowledgeDocumentUpdate) bool {
	if before == nil {
		return false
	}
	if req.CollectionID != nil && *req.CollectionID != before.CollectionID {
		return true
	}
	if req.DocumentID != nil && *req.DocumentID != before.DocumentID {
		return true
	}
	if req.Title != nil && *req.Title != before.Title {
		return true
	}
	if req.Content != nil && *req.Content != before.Content {
		return true
	}
	if req.ContentType != nil && *req.ContentType != before.ContentType {
		return true
	}
	if req.Source != nil && *req.Source != before.Source {
		return true
	}
	if req.Metadata != nil {
		current, _ := sonic.MarshalString(before.Metadata)
		next := strings.TrimSpace(*req.Metadata)
		if next == "" {
			next = "{}"
		}
		return current != next
	}
	return false
}

type noopKnowledgeDocumentIndexer struct{}

func (noopKnowledgeDocumentIndexer) ImportFile(ctx context.Context, input appsvc.ImportKnowledgeFileInput) (*models.KnowledgeDocument, error) {
	return &models.KnowledgeDocument{}, nil
}

func (noopKnowledgeDocumentIndexer) IndexDocument(ctx context.Context, documentID uint64) error {
	return nil
}

func (noopKnowledgeDocumentIndexer) DeleteDocumentVectors(ctx context.Context, documentID uint64) error {
	return nil
}

func (noopKnowledgeDocumentIndexer) EnsureCollection(ctx context.Context, collection *models.KnowledgeCollection) error {
	return nil
}

func (noopKnowledgeDocumentIndexer) DropCollection(ctx context.Context, collectionName string) error {
	return nil
}
