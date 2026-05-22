package logic

import (
	"errors"
	"strings"
	"unicode"

	"admin/internal/appsvc"
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"go-common/utils/slices_utils"
	"go-common/utils/str"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc"
	"orm-crud/gormc/mixin"

	"gorm.io/gen/field"
	"gorm.io/gorm"
)

const (
	HttpMethodGet     = "GET"
	HttpMethodPost    = "POST"
	HttpMethodPut     = "PUT"
	HttpMethodPatch   = "PATCH"
	HttpMethodDelete  = "DELETE"
	HttpMethodOptions = "OPTIONS"
	HttpMethodHead    = "HEAD"
)

var validResourceApiMethods = map[string]struct{}{
	HttpMethodGet:     {},
	HttpMethodPost:    {},
	HttpMethodPut:     {},
	HttpMethodPatch:   {},
	HttpMethodDelete:  {},
	HttpMethodOptions: {},
	HttpMethodHead:    {},
}

type SysResourceApiHandler struct {
	Q      *query.Query
	Casbin appsvc.CasbinService
}

func NewSysResourceApiHandler(q *query.Query, casbin appsvc.CasbinService) *SysResourceApiHandler {
	return &SysResourceApiHandler{Q: q, Casbin: casbin}
}

type RespSysResourceApi struct {
	models.SysResourceApi
	CanWrite  bool `json:"canWrite"`
	CanDelete bool `json:"canDelete"`
}

type ReqResourceApiCreate struct {
	Module    string `json:"module" change:"业务模块" binding:"max=128" binding_msg:"max=业务模块最多128位"`
	Path      string `json:"path" change:"接口路径" binding:"required,max=255" binding_msg:"required=接口路径不能为空,max=接口路径最多255位"`
	Method    string `json:"method" change:"请求方法" binding:"required,max=16" binding_msg:"required=请求方法不能为空,max=请求方法最多16位"`
	SortOrder int32  `json:"sortOrder" change:"排序"`
	IsEnabled bool   `json:"isEnabled" change:"启用状态"`
	Remark    string `json:"remark" change:"备注" binding:"max=255" binding_msg:"max=备注最多255位"`
}

type ReqResourceApiUpdate struct {
	ID        uint64  `json:"id" binding:"required" binding_msg:"required=请求错误"`
	Module    *string `json:"module" change:"业务模块" binding:"omitempty,max=128" binding_msg:"max=业务模块最多128位"`
	Path      *string `json:"path" change:"接口路径" binding:"omitempty,max=255" binding_msg:"max=接口路径最多255位"`
	Method    *string `json:"method" change:"请求方法" binding:"omitempty,max=16" binding_msg:"max=请求方法最多16位"`
	SortOrder *int32  `json:"sortOrder" change:"排序"`
	IsEnabled *bool   `json:"isEnabled" change:"启用状态"`
	Remark    *string `json:"remark" change:"备注" binding:"omitempty,max=255" binding_msg:"max=备注最多255位"`
}

type ReqResourceApiBatchDelete struct {
	IDs []uint64 `json:"ids" binding:"required,min=1" binding_msg:"required=请选择API资源,min=至少选择一项"`
}

// @Summary 获取API资源分页列表
// @Remark 分页查询 API 资源信息
// @Tags ResourceApi
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[RespSysResourceApi]} "成功"
// @Router /api/sys/resource/api/list [post]
func (h *SysResourceApiHandler) List(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[RespSysResourceApi], error) {
	if req.GetOrderBy() == "" {
		orderBy := "sort_order asc,id desc"
		req.OrderBy = &orderBy
	}
	pagination, err := h.Q.SysResourceApi.PageWithPaging(req)
	if err != nil {
		return nil, res.FailDefault
	}
	items := make([]*RespSysResourceApi, 0, len(pagination.Items))
	for _, item := range pagination.Items {
		items = append(items, &RespSysResourceApi{
			SysResourceApi: *item,
			CanWrite:       true,
			CanDelete:      true,
		})
	}
	return &gormc.PagingResult[RespSysResourceApi]{
		Items: items,
		Total: pagination.Total,
	}, nil
}

// @Summary 创建API资源
// @Remark 创建新的 API 资源记录
// @Tags ResourceApi
// @Accept json
// @Produce json
// @Param req body ReqResourceApiCreate true "API资源创建参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/resource/api/create [post]
func (h *SysResourceApiHandler) Create(ctx *handler.Ctx, req *ReqResourceApiCreate) error {
	req.normalize()
	if err := validateResourceApiValues(req.Method, req.Path); err != nil {
		return err
	}
	operationID := ctx.SessionInfo.Id
	sysResourceApi := h.Q.SysResourceApi
	err := sysResourceApi.Create(&models.SysResourceApi{
		OperatorID: mixin.OperatorID{
			CreatedBy: mixin.CreatedBy{CreatedBy: operationID},
			UpdatedBy: mixin.UpdatedBy{UpdatedBy: operationID},
		},
		SortOrder: mixin.SortOrder{SortOrder: req.SortOrder},
		Remark:    mixin.Remark{Remark: req.Remark},
		IsEnabled: mixin.IsEnabled{IsEnabled: req.IsEnabled},
		Module:    req.Module,
		Path:      req.Path,
		Method:    req.Method,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return res.FailMsg("API资源已存在")
		}
		return res.FailDefault
	}
	return nil
}

// @Summary 更新API资源
// @Remark 根据 ID 更新 API 资源信息，并同步权限策略
// @Tags ResourceApi
// @Accept json
// @Produce json
// @Param req body ReqResourceApiUpdate true "API资源更新参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/resource/api/update [post]
func (h *SysResourceApiHandler) Update(ctx *handler.Ctx, req *ReqResourceApiUpdate) error {
	req.normalize()

	sysResourceApi := h.Q.SysResourceApi
	current, err := sysResourceApi.Select(sysResourceApi.ID, sysResourceApi.Method, sysResourceApi.Path, sysResourceApi.IsEnabled).Where(sysResourceApi.ID.Eq(req.ID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res.FailMsg("API资源不存在")
		}
		return res.FailDefault
	}

	method := current.Method
	path := current.Path
	isEnabled := current.IsEnabled.IsEnabled
	if req.Method != nil {
		method = *req.Method
	}
	if req.Path != nil {
		path = *req.Path
	}
	if req.IsEnabled != nil {
		isEnabled = *req.IsEnabled
	}
	if err := validateResourceApiValues(method, path); err != nil {
		return err
	}

	exprs := []field.AssignExpr{sysResourceApi.UpdatedBy.Value(ctx.SessionInfo.Id)}
	query.ExprAppendSelf(&exprs, req.Module, sysResourceApi.Module.Value)
	query.ExprAppendSelf(&exprs, req.Path, sysResourceApi.Path.Value)
	query.ExprAppendSelf(&exprs, req.Method, sysResourceApi.Method.Value)
	query.ExprAppendSelf(&exprs, req.SortOrder, sysResourceApi.SortOrder.Value)
	query.ExprAppendSelf(&exprs, req.IsEnabled, sysResourceApi.IsEnabled.Value)
	query.ExprAppendSelf(&exprs, req.Remark, sysResourceApi.Remark.Value)

	info, err := sysResourceApi.Where(sysResourceApi.ID.Eq(req.ID)).UpdateSimple(exprs...)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return res.FailMsg("API资源已存在")
		}
		return res.FailDefault
	}
	if info.RowsAffected == 0 {
		return res.FailMsg("API资源不存在")
	}
	if current.IsEnabled.IsEnabled != isEnabled || current.Method != method || current.Path != path {
		if err := h.Casbin.SyncAPIResourcePolicies(
			current,
			&models.SysResourceApi{
				AutoIncrementID: current.AutoIncrementID,
				IsEnabled:       mixin.IsEnabled{IsEnabled: isEnabled},
				Path:            path,
				Method:          method,
			},
		); err != nil {
			return res.FailDefault
		}
	}
	return nil
}

// @Summary 删除API资源
// @Remark 根据 ID 列表批量删除 API 资源，并移除对应权限策略
// @Tags ResourceApi
// @Accept json
// @Produce json
// @Param req body ReqResourceApiBatchDelete true "批量删除参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/resource/api/del [post]
func (h *SysResourceApiHandler) Del(ctx *handler.Ctx, req *ReqResourceApiBatchDelete) error {
	ids := slices_utils.Distinct(req.IDs)
	if len(ids) == 0 {
		return res.FailMsg("请选择API资源")
	}
	sysResourceApi := h.Q.SysResourceApi
	apis, err := sysResourceApi.Select(sysResourceApi.ID, sysResourceApi.Path, sysResourceApi.Method, sysResourceApi.IsEnabled).Where(sysResourceApi.ID.In(ids...)).Find()
	if err != nil {
		return res.FailDefault
	}
	info, err := sysResourceApi.Where(sysResourceApi.ID.In(ids...)).Delete()
	if err != nil {
		return res.FailDefault
	}
	if info.RowsAffected != int64(len(ids)) {
		return res.FailMsg("API资源不存在")
	}
	for _, api := range apis {
		if err := h.Casbin.RemoveAPIResourcePolicies(api); err != nil {
			return res.FailDefault
		}
	}
	return nil
}

func (req *ReqResourceApiCreate) normalize() {
	req.Module = strings.TrimSpace(req.Module)
	req.Path = normalizeResourceApiPath(strings.TrimSpace(req.Path))
	req.Method = strings.ToUpper(strings.TrimSpace(req.Method))
	req.Remark = strings.TrimSpace(req.Remark)
}

func (req *ReqResourceApiUpdate) normalize() {
	str.TrimStringPtr(req.Module, nil)
	str.TrimStringPtr(req.Remark, nil)
	str.TrimStringPtr(req.Method, func(value string) string {
		return strings.ToUpper(value)
	})
	if req.Path != nil {
		str.TrimStringPtr(req.Path, normalizeResourceApiPath)
	}
}

func validateResourceApiValues(method, path string) error {
	if _, ok := validResourceApiMethods[method]; !ok {
		return res.FailMsg("请求方法无效")
	}
	if path == "" {
		return res.FailMsg("接口路径不能为空")
	}
	if !strings.HasPrefix(path, "/") {
		return res.FailMsg("接口路径必须以 / 开头")
	}
	return validateResourceApiPathTemplate(path)
}

func normalizeResourceApiPath(path string) string {
	if path == "" {
		return ""
	}
	segments := strings.Split(path, "/")
	for i, segment := range segments {
		if len(segment) >= 3 && strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") {
			name := strings.TrimSpace(segment[1 : len(segment)-1])
			if name != "" {
				segments[i] = ":" + name
			}
		}
	}
	return strings.Join(segments, "/")
}

func validateResourceApiPathTemplate(path string) error {
	for segment := range strings.SplitSeq(path, "/") {
		if segment == "" {
			continue
		}
		if strings.Contains(segment, "{") || strings.Contains(segment, "}") {
			return res.FailMsg("路径参数请使用 {name} 或 :name 格式")
		}
		if after, ok := strings.CutPrefix(segment, ":"); ok {
			if !isValidResourceApiParamName(after) {
				return res.FailMsg("路径参数名只能包含字母、数字、下划线，且不能以数字开头")
			}
		}
	}
	return nil
}

func isValidResourceApiParamName(name string) bool {
	if name == "" {
		return false
	}
	for i, r := range name {
		if i == 0 {
			if r != '_' && !unicode.IsLetter(r) {
				return false
			}
			continue
		}
		if r != '_' && !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
