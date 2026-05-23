# 知识库前后端实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现知识库 Collection 和 Document 的完整前后端 CRUD 功能，前端展示 Collection 列表，点击后跳转新 Tab 展示对应 Document 列表，并在 Document 页面顶部显示所属 Collection 信息。

**Architecture:** 后端采用 Handler + GORM Query 模式提供 REST API；前端使用 TanStack Router 路由、alova 请求、ProTable 表格、Drawer 表单，保持与现有系统（任务调度、字典管理）一致的交互风格。

**Tech Stack:** Go Fiber v3 + GORM + Gen Query / React + TanStack Router + Ant Design Pro + alova

---

## 文件结构

### 后端（Go）

| 文件 | 职责 |
|------|------|
| `backend/go/admin/internal/services/orm/models/knowledge_collection.go` | Collection 模型（已创建） |
| `backend/go/admin/internal/services/orm/models/knowledge_document.go` | Document 模型（已创建） |
| `backend/go/admin/internal/router/logic/knowledge_collection.go` | Collection Handler：列表、详情、创建、更新、删除 |
| `backend/go/admin/internal/router/logic/knowledge_document.go` | Document Handler：列表、详情、创建、更新、删除、按 Collection 筛选 |
| `backend/go/admin/internal/router/auth_router/knowledge.go` | 路由注册：/api/knowledge/collection/*、/api/knowledge/document/* |
| `backend/go/admin/internal/router/auth_router/auth_router.go` | 新增 knowledge 路由注册调用 |

### 前端（React）

| 文件 | 职责 |
|------|------|
| `front/apps/admin-react/src/api/business/knowledgeCollection.ts` | Collection API 封装 |
| `front/apps/admin-react/src/api/business/knowledgeDocument.ts` | Document API 封装 |
| `front/apps/admin-react/src/routes/_app/knowledge.tsx` | Knowledge 父路由布局 |
| `front/apps/admin-react/src/routes/_app/knowledge/collection.tsx` | Collection 列表页（ProTable + Drawer 表单） |
| `front/apps/admin-react/src/routes/_app/knowledge/document.tsx` | Document 列表页（顶部显示 Collection 信息 + ProTable） |

---

## 任务清单

### Task 1: 后端 Collection Handler

**Files:**
- Create: `backend/go/admin/internal/router/logic/knowledge_collection.go`

- [ ] **Step 1: 实现 Collection Handler 结构体和请求/响应类型**

```go
package logic

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"go-common/utils/str"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc"
	"orm-crud/gormc/mixin"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type KnowledgeCollectionHandler struct {
	Q *query.Query
}

func NewKnowledgeCollectionHandler(q *query.Query) *KnowledgeCollectionHandler {
	return &KnowledgeCollectionHandler{Q: q}
}

type RespKnowledgeCollection struct {
	models.KnowledgeCollection
	CanWrite  bool `json:"canWrite"`
	CanDelete bool `json:"canDelete"`
}

type ReqKnowledgeCollectionCreate struct {
	CollectionName   string `json:"collectionName" change:"集合名称" binding:"required,max=128" binding_msg:"required=集合名称不能为空,max=集合名称最多128位"`
	DisplayName      string `json:"displayName" change:"显示名称" binding:"required,max=255" binding_msg:"required=显示名称不能为空,max=显示名称最多255位"`
	Description      string `json:"description" change:"描述" binding:"max=512" binding_msg:"max=描述最多512位"`
	EmbeddingModel   string `json:"embeddingModel" change:"Embedding模型" binding:"required,max=128" binding_msg:"required=Embedding模型不能为空,max=Embedding模型最多128位"`
	VectorDimension  int    `json:"vectorDimension" change:"向量维度" binding:"required,min=1" binding_msg:"required=向量维度不能为空,min=向量维度必须大于0"`
	MetricType       string `json:"metricType" change:"度量类型" binding:"max=32" binding_msg:"max=度量类型最多32位"`
	IndexType        string `json:"indexType" change:"索引类型" binding:"max=32" binding_msg:"max=索引类型最多32位"`
	IDMaxLength      int    `json:"idMaxLength" change:"ID最大长度"`
	ContentMaxLength int    `json:"contentMaxLength" change:"Content最大长度"`
	IsEnabled        bool   `json:"isEnabled" change:"启用状态"`
	Remark           string `json:"remark" change:"备注" binding:"max=255" binding_msg:"max=备注最多255位"`
}

type ReqKnowledgeCollectionUpdate struct {
	ID               uint64  `json:"id" binding:"required" binding_msg:"required=请求错误"`
	CollectionName   *string `json:"collectionName" change:"集合名称" binding:"omitempty,max=128" binding_msg:"max=集合名称最多128位"`
	DisplayName      *string `json:"displayName" change:"显示名称" binding:"omitempty,max=255" binding_msg:"max=显示名称最多255位"`
	Description      *string `json:"description" change:"描述" binding:"omitempty,max=512" binding_msg:"max=描述最多512位"`
	EmbeddingModel   *string `json:"embeddingModel" change:"Embedding模型" binding:"omitempty,max=128" binding_msg:"max=Embedding模型最多128位"`
	VectorDimension  *int    `json:"vectorDimension" change:"向量维度" binding:"omitempty,min=1" binding_msg:"min=向量维度必须大于0"`
	MetricType       *string `json:"metricType" change:"度量类型" binding:"omitempty,max=32" binding_msg:"max=度量类型最多32位"`
	IndexType        *string `json:"indexType" change:"索引类型" binding:"omitempty,max=32" binding_msg:"max=索引类型最多32位"`
	IDMaxLength      *int    `json:"idMaxLength" change:"ID最大长度"`
	ContentMaxLength *int    `json:"contentMaxLength" change:"Content最大长度"`
	IsEnabled        *bool   `json:"isEnabled" change:"启用状态"`
	Remark           *string `json:"remark" change:"备注" binding:"omitempty,max=255" binding_msg:"max=备注最多255位"`
}

type ReqKnowledgeCollectionID struct {
	ID uint64 `json:"id" binding:"required" binding_msg:"required=请求错误"`
}

```

- [ ] **Step 2: 实现 List 方法**

```go
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
```

- [ ] **Step 3: 实现 Detail 方法**

```go
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
```

- [ ] **Step 4: 实现 Create 方法**

```go
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
	idMaxLength := req.IDMaxLength
	if idMaxLength == 0 {
		idMaxLength = 255
	}
	contentMaxLength := req.ContentMaxLength
	if contentMaxLength == 0 {
		contentMaxLength = 8192
	}

	err := h.Q.KnowledgeCollection.Create(&models.KnowledgeCollection{
		OperatorID: mixin.OperatorID{
			CreatedBy: mixin.CreatedBy{CreatedBy: operationID},
			UpdatedBy: mixin.UpdatedBy{UpdatedBy: operationID},
		},
		IsEnabled:        mixin.IsEnabled{IsEnabled: req.IsEnabled},
		Remark:           mixin.Remark{Remark: req.Remark},
		CollectionName:   req.CollectionName,
		DisplayName:      req.DisplayName,
		Description:      req.Description,
		EmbeddingModel:   req.EmbeddingModel,
		VectorDimension:  req.VectorDimension,
		MetricType:       metricType,
		IndexType:        indexType,
		IDMaxLength:      idMaxLength,
		ContentMaxLength: contentMaxLength,
		Status:           "active",
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return res.FailMsg("集合名称已存在")
		}
		ctx.L().Error("create knowledge collection fail", zap.Error(err))
		return res.FailDefault
	}
	return nil
}
```

- [ ] **Step 5: 实现 Update 方法**

```go
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
	query.ExprAppendSelf(&exprs, req.Description, collection.Description.Value)
	query.ExprAppendSelf(&exprs, req.EmbeddingModel, collection.EmbeddingModel.Value)
	query.ExprAppendSelf(&exprs, req.VectorDimension, collection.VectorDimension.Value)
	query.ExprAppendSelf(&exprs, req.MetricType, collection.MetricType.Value)
	query.ExprAppendSelf(&exprs, req.IndexType, collection.IndexType.Value)
	query.ExprAppendSelf(&exprs, req.IDMaxLength, collection.IDMaxLength.Value)
	query.ExprAppendSelf(&exprs, req.ContentMaxLength, collection.ContentMaxLength.Value)
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
```

- [ ] **Step 6: 实现 Del 方法**

```go
// @Summary 删除知识库集合
// @Tags KnowledgeCollection
// @Accept json
// @Produce json
// @Param req body ReqKnowledgeCollectionID true "集合ID"
// @Success 200 {object} res.Response "成功"
// @Router /api/knowledge/collection/del [post]
func (h *KnowledgeCollectionHandler) Del(ctx *handler.Ctx, req *ReqKnowledgeCollectionID) error {
	collection := h.Q.KnowledgeCollection
	info, err := collection.Where(collection.ID.Eq(req.ID)).Delete()
	if err != nil {
		ctx.L().Error("delete knowledge collection fail", zap.Error(err), zap.Uint64("id", req.ID))
		return res.FailDefault
	}
	if info.RowsAffected == 0 {
		return res.FailMsg("集合不存在")
	}
	return nil
}
```

---

### Task 2: 后端 Document Handler

**Files:**
- Create: `backend/go/admin/internal/router/logic/knowledge_document.go`

- [ ] **Step 1: 实现 Document Handler 结构体和请求/响应类型**

```go
package logic

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc"
	"orm-crud/gormc/mixin"

	"github.com/bytedance/sonic"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type KnowledgeDocumentHandler struct {
	Q *query.Query
}

func NewKnowledgeDocumentHandler(q *query.Query) *KnowledgeDocumentHandler {
	return &KnowledgeDocumentHandler{Q: q}
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

type ReqKnowledgeDocumentListByCollection struct {
	CollectionID uint64 `json:"collectionID" binding:"required" binding_msg:"required=请选择集合"`
	v1.PagingRequest
}
```

- [ ] **Step 2: 实现 List 方法（支持按 Collection 筛选）**

```go
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
```

- [ ] **Step 3: 实现 ListByCollection 方法**

```go
// @Summary 按集合获取知识库文档分页列表
// @Tags KnowledgeDocument
// @Accept json
// @Produce json
// @Param req body ReqKnowledgeDocumentListByCollection true "查询参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[RespKnowledgeDocument]} "成功"
// @Router /api/knowledge/document/listByCollection [post]
func (h *KnowledgeDocumentHandler) ListByCollection(ctx *handler.Ctx, req *ReqKnowledgeDocumentListByCollection) (*gormc.PagingResult[RespKnowledgeDocument], error) {
	if req.PagingRequest.GetOrderBy() == "" {
		orderBy := "id desc"
		req.PagingRequest.OrderBy = &orderBy
	}

	doc := h.Q.KnowledgeDocument
	pagination, err := doc.Where(doc.CollectionID.Eq(req.CollectionID)).PageWithPaging(&req.PagingRequest)
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
```

- [ ] **Step 4: 实现 Detail 方法**

```go
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
```

- [ ] **Step 5: 实现 Create 方法**

```go
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

	err := h.Q.KnowledgeDocument.Create(&models.KnowledgeDocument{
		OperatorID: mixin.OperatorID{
			CreatedBy: mixin.CreatedBy{CreatedBy: operationID},
			UpdatedBy: mixin.UpdatedBy{UpdatedBy: operationID},
		},
		IsEnabled:   mixin.IsEnabled{IsEnabled: req.IsEnabled},
		Remark:      mixin.Remark{Remark: req.Remark},
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
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return res.FailMsg("文档ID已存在")
		}
		ctx.L().Error("create knowledge document fail", zap.Error(err))
		return res.FailDefault
	}
	return nil
}
```

- [ ] **Step 6: 实现 Update 方法**

```go
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
	return nil
}
```

- [ ] **Step 7: 实现 Del 方法**

```go
// @Summary 删除知识库文档
// @Tags KnowledgeDocument
// @Accept json
// @Produce json
// @Param req body ReqKnowledgeDocumentID true "文档ID"
// @Success 200 {object} res.Response "成功"
// @Router /api/knowledge/document/del [post]
func (h *KnowledgeDocumentHandler) Del(ctx *handler.Ctx, req *ReqKnowledgeDocumentID) error {
	doc := h.Q.KnowledgeDocument
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
```

---

### Task 3: 后端路由注册

**Files:**
- Create: `backend/go/admin/internal/router/auth_router/knowledge.go`
- Modify: `backend/go/admin/internal/router/auth_router/auth_router.go`

- [ ] **Step 1: 创建 knowledge 路由注册文件**

```go
package auth_router

import (
	"admin/internal/fiberc/handler"
	"admin/internal/router/logic"
	"admin/internal/services/orm/query"

	"github.com/gofiber/fiber/v3"
)

func registerKnowledgeRouters(router fiber.Router) {
	collectionHandler := logic.NewKnowledgeCollectionHandler(query.Q)
	documentHandler := logic.NewKnowledgeDocumentHandler(query.Q)

	collectionGroup := router.Group("/knowledge/collection")
	collectionGroup.Post("/list", handler.CtxHandlerFunc(collectionHandler.List))
	collectionGroup.Post("/detail", handler.CtxHandlerFunc(collectionHandler.Detail))
	collectionGroup.Post("/create", handler.CtxHandlerNilFunc(collectionHandler.Create))
	collectionGroup.Post("/update", handler.CtxHandlerNilFunc(collectionHandler.Update))
	collectionGroup.Post("/del", handler.CtxHandlerNilFunc(collectionHandler.Del))

	documentGroup := router.Group("/knowledge/document")
	documentGroup.Post("/list", handler.CtxHandlerFunc(documentHandler.List))
	documentGroup.Post("/listByCollection", handler.CtxHandlerFunc(documentHandler.ListByCollection))
	documentGroup.Post("/detail", handler.CtxHandlerFunc(documentHandler.Detail))
	documentGroup.Post("/create", handler.CtxHandlerNilFunc(documentHandler.Create))
	documentGroup.Post("/update", handler.CtxHandlerNilFunc(documentHandler.Update))
	documentGroup.Post("/del", handler.CtxHandlerNilFunc(documentHandler.Del))
}
```

- [ ] **Step 2: 在 auth_router.go 中注册 knowledge 路由**

在 `registerJobExecutionRouters` 调用后添加：

```go
registerKnowledgeRouters(group)
```

---

### Task 4: 前端 API 封装

**Files:**
- Create: `front/apps/admin-react/src/api/business/knowledgeCollection.ts`
- Create: `front/apps/admin-react/src/api/business/knowledgeDocument.ts`

- [ ] **Step 1: 创建 Collection API**

```typescript
import type { PagingRequest, PagingResult } from '@vp/core'
import API from '../index'

export interface KnowledgeCollection {
  id: number
  collectionName: string
  displayName: string
  description: string
  embeddingModel: string
  vectorDimension: number
  metricType: string
  indexType: string
  idMaxLength: number
  contentMaxLength: number
  documentCount: number
  status: string
  isEnabled: boolean
  remark: string
  createdAt?: string
  updatedAt?: string
  canWrite?: boolean
  canDelete?: boolean
}

export interface ReqKnowledgeCollectionCreate {
  collectionName: string
  displayName: string
  description?: string
  embeddingModel: string
  vectorDimension: number
  metricType?: string
  indexType?: string
  idMaxLength?: number
  contentMaxLength?: number
  isEnabled: boolean
  remark?: string
}

export interface ReqKnowledgeCollectionUpdate extends Partial<ReqKnowledgeCollectionCreate> {
  id: number
}

export interface ReqKnowledgeCollectionID {
  id: number
}

function list(req: PagingRequest) {
  return API.Post<Res<PagingResult<KnowledgeCollection>>>('/api/knowledge/collection/list', req, {
    cacheFor: 0,
  })
}

async function detail(req: ReqKnowledgeCollectionID) {
  return await API.Post<Res<KnowledgeCollection>>('/api/knowledge/collection/detail', req, {
    cacheFor: 0,
  }).send()
}

async function create(req: ReqKnowledgeCollectionCreate) {
  await API.Post<Res>('/api/knowledge/collection/create', req, {
    cacheFor: 0,
  }).send()
}

async function update(req: ReqKnowledgeCollectionUpdate) {
  await API.Post<Res>('/api/knowledge/collection/update', req, {
    cacheFor: 0,
  }).send()
}

async function del(req: ReqKnowledgeCollectionID) {
  await API.Post<Res>('/api/knowledge/collection/del', req, {
    cacheFor: 0,
  }).send()
}

export const KnowledgeCollectionApi = {
  list,
  detail,
  create,
  update,
  del,
}
```

- [ ] **Step 2: 创建 Document API**

```typescript
import type { PagingRequest, PagingResult } from '@vp/core'
import API from '../index'

export interface KnowledgeDocument {
  id: number
  collectionID: number
  documentID: string
  title: string
  content: string
  contentType: string
  source: string
  chunkIndex: number
  totalChunks: number
  vectorStatus: string
  vectorID: string
  metadata: Record<string, unknown>
  indexingError: string
  lastIndexedAt: number
  isEnabled: boolean
  remark: string
  createdAt?: string
  updatedAt?: string
  canWrite?: boolean
  canDelete?: boolean
  collection?: {
    collectionName: string
    displayName: string
  }
}

export interface ReqKnowledgeDocumentCreate {
  collectionID: number
  documentID: string
  title: string
  content: string
  contentType?: string
  source?: string
  chunkIndex?: number
  totalChunks?: number
  metadata?: string
  isEnabled: boolean
  remark?: string
}

export interface ReqKnowledgeDocumentUpdate extends Partial<ReqKnowledgeDocumentCreate> {
  id: number
}

export interface ReqKnowledgeDocumentID {
  id: number
}

export interface ReqKnowledgeDocumentListByCollection {
  collectionID: number
  page?: number
  pageSize?: number
  orderBy?: string
  query?: string
}

function list(req: PagingRequest) {
  return API.Post<Res<PagingResult<KnowledgeDocument>>>('/api/knowledge/document/list', req, {
    cacheFor: 0,
  })
}

function listByCollection(req: ReqKnowledgeDocumentListByCollection) {
  return API.Post<Res<PagingResult<KnowledgeDocument>>>('/api/knowledge/document/listByCollection', req, {
    cacheFor: 0,
  })
}

async function detail(req: ReqKnowledgeDocumentID) {
  return await API.Post<Res<KnowledgeDocument>>('/api/knowledge/document/detail', req, {
    cacheFor: 0,
  }).send()
}

async function create(req: ReqKnowledgeDocumentCreate) {
  await API.Post<Res>('/api/knowledge/document/create', req, {
    cacheFor: 0,
  }).send()
}

async function update(req: ReqKnowledgeDocumentUpdate) {
  await API.Post<Res>('/api/knowledge/document/update', req, {
    cacheFor: 0,
  }).send()
}

async function del(req: ReqKnowledgeDocumentID) {
  await API.Post<Res>('/api/knowledge/document/del', req, {
    cacheFor: 0,
  }).send()
}

export const KnowledgeDocumentApi = {
  list,
  listByCollection,
  detail,
  create,
  update,
  del,
}
```

---

### Task 5: 前端路由

**Files:**
- Create: `front/apps/admin-react/src/routes/_app/knowledge.tsx`
- Create: `front/apps/admin-react/src/routes/_app/knowledge/collection.tsx`
- Create: `front/apps/admin-react/src/routes/_app/knowledge/document.tsx`

- [ ] **Step 1: 创建 Knowledge 父路由**

```typescript
import { createFileRoute, Outlet } from '@tanstack/react-router'

export const Route = createFileRoute('/_app/knowledge')({
  component: RouteComponent,
})

function RouteComponent() {
  return <Outlet />
}
```

- [ ] **Step 2: 创建 Collection 列表页**

参考 `job/schedule.tsx` 实现 ProTable + Drawer 表单，包含：
- 表格列：ID、集合名称、显示名称、Embedding模型、向量维度、文档数量、状态、创建时间、操作
- 搜索框：按名称搜索
- 操作按钮：创建、编辑、启用/停用、删除
- Drawer 表单：创建/编辑 Collection

关键代码结构：

```typescript
import type { ProColumns } from '@ant-design/pro-components'
import type { KnowledgeCollection } from '~/api/business/knowledgeCollection'
import { ProFormText, ProFormDigit, ProFormSelect, ProFormTextArea, ProTable } from '@ant-design/pro-components'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { usePagination } from 'alova/client'
import { Button, Drawer, Form, Input, Popconfirm, Space, Tag } from 'antd'
import { useCallback, useMemo, useState } from 'react'
import { KnowledgeCollectionApi } from '~/api/business/knowledgeCollection'
import { gMessage } from '~/utils/message'

export const Route = createFileRoute('/_app/knowledge/collection')({
  staleTime: 1000 * 60 * 2,
  component: KnowledgeCollectionManagement,
})

function KnowledgeCollectionManagement() {
  const [drawerOpen, setDrawerOpen] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [editing, setEditing] = useState<KnowledgeCollection>()
  const [searchText, setSearchText] = useState('')
  const [form] = Form.useForm()
  const navigate = useNavigate()

  const { data, total, page, pageSize, loading, update, send } = usePagination(
    (nextPage, nextPageSize) => {
      const filters: Record<string, unknown>[] = []
      const keyword = searchText.trim()
      if (keyword) {
        filters.push({
          $or: [
            { collectionName__icontains: keyword },
            { displayName__icontains: keyword },
          ],
        })
      }
      return KnowledgeCollectionApi.list({
        page: nextPage,
        pageSize: nextPageSize,
        orderBy: 'id desc',
        query: filters.length > 0 ? JSON.stringify({ $and: filters }) : undefined,
      })
    },
    {
      initialData: { total: 0, items: [] },
      initialPage: 1,
      initialPageSize: 20,
      total: response => response.data?.total ?? 0,
      data: response => response.data?.items ?? [],
      watchingStates: [searchText],
      debounce: [500, 0],
    },
  )

  const handleViewDocuments = useCallback((record: KnowledgeCollection) => {
    navigate({ to: '/knowledge/document', search: { collectionId: record.id } })
  }, [navigate])

  // 启用/停用通过 update 接口实现
  // () => KnowledgeCollectionApi.update({ id: record.id, isEnabled: nextEnabled })

  // ... 其余实现参考 job/schedule.tsx
}
```

- [ ] **Step 3: 创建 Document 列表页**

Document 页面顶部显示所属 Collection 信息，使用 `useSearch` 获取 `collectionId` 参数。

```typescript
import type { ProColumns } from '@ant-design/pro-components'
import type { KnowledgeDocument } from '~/api/business/knowledgeDocument'
import { ProTable } from '@ant-design/pro-components'
import { createFileRoute, useSearch } from '@tanstack/react-router'
import { usePagination } from 'alova/client'
import { Button, Card, Input, Popconfirm, Space, Tag } from 'antd'
import { useMemo, useState } from 'react'
import { KnowledgeDocumentApi } from '~/api/business/knowledgeDocument'
import { KnowledgeCollectionApi } from '~/api/business/knowledgeCollection'
import { gMessage } from '~/utils/message'

export const Route = createFileRoute('/_app/knowledge/document')({
  staleTime: 1000 * 60 * 2,
  component: KnowledgeDocumentManagement,
})

function KnowledgeDocumentManagement() {
  const search = useSearch({ from: '/_app/knowledge/document' })
  const collectionId = Number(search.collectionId) || 0
  const [searchText, setSearchText] = useState('')
  const [collectionName, setCollectionName] = useState('')

  // 加载 Collection 信息
  useEffect(() => {
    if (!collectionId) return
    KnowledgeCollectionApi.detail({ id: collectionId })
      .then((res) => {
        if (res.data) {
          setCollectionName(res.data.displayName || res.data.collectionName)
        }
      })
      .catch(() => {})
  }, [collectionId])

  const { data, total, page, pageSize, loading, update, send } = usePagination(
    (nextPage, nextPageSize) => {
      const filters: Record<string, unknown>[] = []
      const keyword = searchText.trim()
      if (keyword) {
        filters.push({
          $or: [
            { title__icontains: keyword },
            { documentID__icontains: keyword },
          ],
        })
      }
      return KnowledgeDocumentApi.listByCollection({
        collectionID: collectionId,
        page: nextPage,
        pageSize: nextPageSize,
        orderBy: 'id desc',
        query: filters.length > 0 ? JSON.stringify({ $and: filters }) : undefined,
      })
    },
    {
      initialData: { total: 0, items: [] },
      initialPage: 1,
      initialPageSize: 20,
      total: response => response.data?.total ?? 0,
      data: response => response.data?.items ?? [],
      watchingStates: [searchText, collectionId],
      debounce: [500, 0],
    },
  )

  const columns: ProColumns<KnowledgeDocument>[] = useMemo(() => [
    { title: '文档ID', dataIndex: 'documentID', width: 200, ellipsis: true },
    { title: '标题', dataIndex: 'title', width: 200, ellipsis: true },
    { title: '内容类型', dataIndex: 'contentType', width: 100 },
    {
      title: '向量状态',
      dataIndex: 'vectorStatus',
      width: 100,
      render: (_, record) => {
        const statusMap: Record<string, { color: string, label: string }> = {
          pending: { color: 'default', label: '待处理' },
          indexed: { color: 'success', label: '已索引' },
          failed: { color: 'error', label: '失败' },
        }
        const item = statusMap[record.vectorStatus] || { color: 'default', label: record.vectorStatus }
        return <Tag color={item.color}>{item.label}</Tag>
      },
    },
    { title: '分块', dataIndex: 'chunkIndex', width: 80, render: (_, record) => `${record.chunkIndex + 1}/${record.totalChunks}` },
    { title: '创建时间', dataIndex: 'createdAt', width: 160 },
    {
      title: '操作',
      valueType: 'option',
      width: 200,
      fixed: 'right',
      render: (_, record) => (
        <Space>
          <Button type="link" size="small">编辑</Button>
          <Popconfirm title="确认删除该文档？" onConfirm={() => handleDelete(record.id)}>
            <Button type="link" size="small" danger>删除</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ], [])

  return (
    <>
      {/* 顶部 Collection 信息卡片 */}
      <Card style={{ marginBottom: 16 }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <span style={{ fontSize: 16, fontWeight: 600 }}>所属集合：</span>
          <Tag color="blue">{collectionName || '未选择'}</Tag>
          {!collectionId && <span style={{ color: '#999' }}>请从集合列表选择一个集合查看文档</span>}
        </div>
      </Card>

      <ProTable<KnowledgeDocument>
        rowKey="id"
        headerTitle="文档管理"
        columns={columns}
        dataSource={data}
        loading={loading}
        search={false}
        pagination={{
          showSizeChanger: true,
          current: page,
          pageSize,
          total,
          onChange: (nextPage, nextPageSize) => update({ page: nextPage, pageSize: nextPageSize }),
        }}
        toolBarRender={() => [
          <Button key="create" type="primary">创建文档</Button>,
          <Input.Search
            key="search"
            allowClear
            placeholder="搜索标题、文档ID"
            value={searchText}
            onChange={e => setSearchText(e.target.value)}
            style={{ width: 280 }}
          />,
        ]}
      />
    </>
  )
}
```

---

### Task 6: 生成路由并验证

- [ ] **Step 1: 重新生成 TanStack Router 路由树**

```bash
cd front/apps/admin-react
pnpm dev
# 或手动运行路由生成
npx tsr generate
```

- [ ] **Step 2: 验证后端编译**

```bash
cd backend/go/admin
go build ./...
```

- [ ] **Step 3: 验证前端编译**

```bash
cd front/apps/admin-react
pnpm build
```

---

## 验证方式

- **命令：**
  - `cd backend/go/admin && go build ./...` - 后端编译通过
  - `cd front/apps/admin-react && pnpm build` - 前端编译通过
  - `cd backend/go/admin && go test ./internal/router/logic/ -run "TestKnowledge" -v` - 后端 Handler 单元测试
  - 启动后端后访问 Swagger 文档确认 API 已注册

- **手工检查：**
  - 访问 `/knowledge/collection` 页面，确认 Collection 列表正常显示
  - 点击 Collection 行内的"查看文档"按钮，确认跳转到 `/knowledge/document?collectionId=xxx`
  - 确认 Document 页面顶部显示正确的 Collection 名称
  - 测试创建、编辑、删除、启用/停用功能

- **观测检查：**
  - 后端日志确认 API 调用正常
  - 数据库确认数据写入正确
  - 浏览器 Network 面板确认请求响应正常

---

## 测试文件

### 后端测试

| 文件 | 测试内容 |
|------|----------|
| `backend/go/admin/internal/router/logic/knowledge_collection_test.go` | Collection Handler 单元测试：列表、详情、创建、更新、删除、状态切换（Update）、默认值、不存在场景 |
| `backend/go/admin/internal/router/logic/knowledge_document_test.go` | Document Handler 单元测试：列表、按集合筛选、详情、创建、元数据解析、更新、删除、不存在场景 |
| `backend/go/admin/internal/router/logic/sqlite_test.go` | 新增 `mustMigrateKnowledge` 辅助函数 |

### 前端 E2E 测试

| 文件 | 测试内容 |
|------|----------|
| `front/apps/admin-react/tests/business/knowledge.spec.ts` | 集合页面加载、创建抽屉校验、搜索框、文档页面集合提示、集合到文档跳转 |

### 初始化 SQL

| 文件 | 内容 |
|------|------|
| `backend/go/admin/cmd/scripts/init.sql` | 知识库菜单（sys_resource_menu）、API（sys_resource_api）、菜单-API关联、Casbin权限规则、角色-菜单/角色-API关联 |

---

## 进度记录

- [x] Task 1: 后端 Collection Handler
- [x] Task 2: 后端 Document Handler
- [x] Task 3: 后端路由注册
- [x] Task 4: 前端 API 封装
- [x] Task 5: 前端路由和页面
- [x] Task 6: 生成路由并验证

## 完成说明

- 后端 `go build ./...` 编译通过
- 前端知识库相关代码编译通过（`pnpm build` 中知识库文件零错误）
- 后端 20 个知识库 Handler 单元测试全部通过
- 修复内容：
  - `knowledge_document.go:ListByCollection`：将 `Where(...).PageWithPaging()` 链式调用改为通过 query filter 合并 `collectionID` 条件后调用 `PageWithPaging`
  - `routeTree.gen.ts`：手动添加 knowledge 路由树注册
  - `collection.tsx`：`ProFormSelect` 的 `isEnabled` options 值从 boolean 改为 number（1/0），同步修改表单初始值和回填值
  - `document.tsx`：添加 `validateSearch` zod schema 声明 `collectionId` 搜索参数类型
- 移除内容：
  - 独立 Switch 接口：Collection 状态切换统一通过 Update 接口实现，减少 API 面
  - 前端 `switchStatus` API 和路由注册同步移除
- 新增内容：
  - `collection.tsx` / `document.tsx`：引入 zod schema 表单校验（`useZodForm` + `superRefine`），实现必填项提示、数值范围校验、JSON 格式校验
  - `document.tsx`：补充 Drawer 表单，实现创建/编辑文档功能
  - `context/project/admin/experience/zod-form-validation.md`：沉淀前端表单校验实现模式
  - 字典化枚举字段：`MetricType`、`IndexType`、`ContentType`、`VectorStatus` 改用 `sys_dict` 管理，前端通过 `useDictMatch` 取选项和渲染标签
  - 新增字典类型 `knowledge:metric_type`（COSINE/IP/L2）、`knowledge:index_type`（auto/hnsw/ivf_flat）、`knowledge:content_type`（text/markdown/pdf/html）、`knowledge:vector_status`（pending/indexed/failed）
  - `domains/dict.ts` 注册 4 个新字典 code
- 删除内容：
  - `KnowledgeCollection.Status` 字段（与 `IsEnabled` 重复，移除后通过 ORM 重新生成 query 文件）

---

## 决策记录

- 2026-05-23: Document 列表使用独立路由 `/knowledge/document?collectionId=xxx` 而非内嵌面板，符合用户"跳转新 tab"的需求，同时保持 URL 可分享。
- 2026-05-23: 后端不引入数据权限（DataPermission）简化实现，后续如需可按 sys_dict 模式补充。
- 2026-05-23: Collection 和 Document 的 CanWrite/CanDelete 暂时固定为 true，后续接入权限系统时动态计算。
- 2026-05-23: 移除独立 Switch 接口，Collection 启用/停用统一通过 Update 接口（仅传 `id` + `isEnabled`）实现，减少 API 面并保持一致性。
- 2026-05-23: 删除 `KnowledgeCollection.Status` 字段。`Status` 与 `IsEnabled` 语义重复（active/archived ≈ true/false），保留 `IsEnabled` 即可，避免双状态字段维护成本。
- 2026-05-23: `MetricType`、`IndexType`、`ContentType`、`VectorStatus` 用 `sys_dict` 管理而非硬编码枚举。优势：选项可以在管理后台动态维护、复用现有翻译/渲染基础设施、与 `is_enabled`、`menu_type` 等系统字典保持一致风格。
