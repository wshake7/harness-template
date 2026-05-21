package logic

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	"admin/internal/service"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"go-common/utils/slices_utils"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc"
	"orm-crud/gormc/mixin"
	paginationFilter "orm-crud/pagination/filter"

	"go.uber.org/zap"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type SysDictHandler struct {
	Q              *query.Query
	DataPermission service.DataPermissionService
}

func NewSysDictHandler(q *query.Query, dp service.DataPermissionService) *SysDictHandler {
	return &SysDictHandler{Q: q, DataPermission: dp}
}

// --- 字典类型 (DictType) ---

type ReqDictTypeCreate struct {
	TypeCode  string `json:"typeCode" change:"字典类型代码" binding:"required,max=128" binding_msg:"required=字典类型代码不能为空,max=字典类型代码最多128位"`
	TypeName  string `json:"typeName" change:"字典类型名称" binding:"required,max=255" binding_msg:"required=字典类型名称不能为空,max=字典类型名称最多255位"`
	IsEnabled bool   `json:"isEnabled" change:"启用状态"`
	SortOrder int32  `json:"sortOrder" change:"排序"`
	Remark    string `json:"remark" change:"描述" binding:"max=255" binding_msg:"max=描述最多255位"`
}

type ReqDictTypeUpdate struct {
	ID        uint64  `json:"id" binding:"required" binding_msg:"required=请求错误"`
	TypeCode  *string `json:"typeCode" change:"字典类型代码" binding:"omitempty,max=128" binding_msg:"max=字典类型代码最多128位"`
	TypeName  *string `json:"typeName" change:"字典类型名称" binding:"omitempty,max=255" binding_msg:"max=字典类型名称最多255位"`
	IsEnabled *bool   `json:"isEnabled" change:"启用状态"`
	SortOrder *int32  `json:"sortOrder" change:"排序"`
	Remark    *string `json:"remark" change:"描述" binding:"omitempty,max=255" binding_msg:"max=描述最多255位"`
}

type ReqDictTypeBatchDelete struct {
	IDs []uint64 `json:"ids" binding:"required,min=1" binding_msg:"required=请选择字典类型,min=至少选择一项"`
}

type RespDictType struct {
	models.SysDictType
	CanWrite  bool `json:"canWrite"`
	CanDelete bool `json:"canDelete"`
}

// @Summary 获取字典类型分页列表
// @Remark 分页查询字典类型信息
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[RespDictType]} "成功"
// @Router /api/dict/type/list [post]
func (h *SysDictHandler) TypeList(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[RespDictType], error) {
	permissionExprs, err := h.DataPermission.BuildFilterExprsForCtx(
		ctx,
		models.SysDictType{}.TableName(),
		models.DataPermissionActionRead,
		models.DataPermissionActionWrite,
		models.DataPermissionActionDelete,
	)
	if err != nil {
		ctx.L().Error("build dict type permission expressions failed", zap.Error(err))
		return nil, res.FailDefault
	}
	if err := h.DataPermission.ApplyPagePermissionExpr(req, permissionExprs[models.DataPermissionActionRead]); err != nil {
		ctx.L().Error("apply dict type read permission failed", zap.Error(err))
		return nil, res.FailDefault
	}

	pagination, err := h.Q.SysDictType.PageWithPaging(req)
	if err != nil {
		return nil, res.FailDefault
	}

	items := make([]*RespDictType, 0, len(pagination.Items))
	ids := make([]uint64, 0, len(pagination.Items))
	for _, item := range pagination.Items {
		ids = append(ids, item.ID)
	}

	writeIDSet, err := queryAllowedDictTypeIDSetByExpr(h.DataPermission, ids, permissionExprs[models.DataPermissionActionWrite])
	if err != nil {
		ctx.L().Error("apply dict type write permission failed", zap.Error(err), zap.Uint64s("ids", ids))
		return nil, res.FailDefault
	}
	deleteIDSet := writeIDSet
	if !reflect.DeepEqual(permissionExprs[models.DataPermissionActionWrite], permissionExprs[models.DataPermissionActionDelete]) {
		deleteIDSet, err = queryAllowedDictTypeIDSetByExpr(h.DataPermission, ids, permissionExprs[models.DataPermissionActionDelete])
		if err != nil {
			ctx.L().Error("apply dict type delete permission failed", zap.Error(err), zap.Uint64s("ids", ids))
			return nil, res.FailDefault
		}
	}

	for _, item := range pagination.Items {
		items = append(items, &RespDictType{
			SysDictType: *item,
			CanWrite:    writeIDSet[item.ID],
			CanDelete:   deleteIDSet[item.ID],
		})
	}
	return &gormc.PagingResult[RespDictType]{
		Items: items,
		Total: pagination.Total,
	}, nil
}

func queryAllowedDictTypeIDSetByExpr(dp service.DataPermissionService, ids []uint64, permissionExpr *v1.FilterExpr) (map[uint64]bool, error) {
	allowedIDSet := make(map[uint64]bool, len(ids))
	if len(ids) == 0 {
		return allowedIDSet, nil
	}

	permissionQuery, err := dp.BuildPermissionQueryFromExpr(permissionExpr)
	if err != nil {
		return nil, err
	}

	sysDictType := permissionQuery.SysDictType
	allowedTypes, err := sysDictType.
		Select(sysDictType.ID).
		Where(sysDictType.ID.In(ids...)).
		Find()
	if err != nil {
		return nil, err
	}
	for _, item := range allowedTypes {
		allowedIDSet[item.ID] = true
	}
	return allowedIDSet, nil
}

func queryAllowedDictTypeIDSet(ctx *handler.Ctx, ids []uint64, buildPermissionQuery func(*handler.Ctx, string) (*query.Query, error)) (map[uint64]bool, error) {
	allowedIDSet := make(map[uint64]bool, len(ids))
	if len(ids) == 0 {
		return allowedIDSet, nil
	}

	permissionQuery, err := buildPermissionQuery(ctx, models.SysDictType{}.TableName())
	if err != nil {
		return nil, err
	}

	sysDictType := permissionQuery.SysDictType
	allowedTypes, err := sysDictType.
		Select(sysDictType.ID).
		Where(sysDictType.ID.In(ids...)).
		Find()
	if err != nil {
		return nil, err
	}
	for _, item := range allowedTypes {
		allowedIDSet[item.ID] = true
	}
	return allowedIDSet, nil
}

func queryAllowedDictTypeIDs(ctx *handler.Ctx, buildPermissionQuery func(*handler.Ctx, string) (*query.Query, error)) ([]uint64, error) {
	permissionQuery, err := buildPermissionQuery(ctx, models.SysDictType{}.TableName())
	if err != nil {
		return nil, err
	}

	sysDictType := permissionQuery.SysDictType
	allowedTypes, err := sysDictType.
		Select(sysDictType.ID).
		Find()
	if err != nil {
		return nil, err
	}

	ids := make([]uint64, 0, len(allowedTypes))
	for _, item := range allowedTypes {
		ids = append(ids, item.ID)
	}
	return ids, nil
}

func applyDictEntryTypeIDPageFilter(req *v1.PagingRequest, typeIDs []uint64) error {
	values := make([]string, 0, len(typeIDs))
	for _, id := range typeIDs {
		values = append(values, strconv.FormatUint(id, 10))
	}

	typeFilterExpr := &v1.FilterExpr{
		Type: v1.ExprType_AND,
		Conditions: []*v1.FilterCondition{
			{
				Field:  "sysDictTypeId",
				Op:     v1.Operator_IN,
				Values: values,
			},
		},
	}

	currentExpr, err := paginationFilter.ConvertFilterByPagingRequest(req)
	if err != nil {
		return err
	}
	if currentExpr == nil {
		req.FilteringType = &v1.PagingRequest_FilterExpr{FilterExpr: typeFilterExpr}
		return nil
	}

	req.FilteringType = &v1.PagingRequest_FilterExpr{FilterExpr: &v1.FilterExpr{
		Type:   v1.ExprType_AND,
		Groups: []*v1.FilterExpr{currentExpr, typeFilterExpr},
	}}
	return nil
}

// @Summary 创建字典类型
// @Remark 创建新的字典类型
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictTypeCreate true "字典类型创建参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/dict/type/create [post]
func (h *SysDictHandler) TypeCreate(ctx *handler.Ctx, req *ReqDictTypeCreate) error {
	operationID := ctx.SessionInfo.Id

	err := h.Q.SysDictType.Create(&models.SysDictType{
		OperatorID: mixin.OperatorID{
			CreatedBy: mixin.CreatedBy{CreatedBy: operationID},
			UpdatedBy: mixin.UpdatedBy{UpdatedBy: operationID},
		},
		IsEnabled: mixin.IsEnabled{IsEnabled: req.IsEnabled},
		SortOrder: mixin.SortOrder{SortOrder: req.SortOrder},
		Remark:    mixin.Remark{Remark: req.Remark},
		TypeCode:  req.TypeCode,
		TypeName:  req.TypeName,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return res.FailMsg("类型编码已存在")
		}
		return res.FailDefault
	}
	return nil
}

// @Summary 更新字典类型
// @Remark 根据 ID 更新字典类型信息
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictTypeUpdate true "字典类型更新参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/dict/type/update [post]
func (h *SysDictHandler) TypeUpdate(ctx *handler.Ctx, req *ReqDictTypeUpdate) error {
	operationID := ctx.SessionInfo.Id
	permissionQuery, err := h.DataPermission.BuildWritePermissionQuery(ctx, models.SysDictType{}.TableName())
	if err != nil {
		ctx.L().Error("apply dict type write permission failed", zap.Error(err), zap.Uint64("id", req.ID))
		return res.FailDefault
	}

	sysDictType := permissionQuery.SysDictType

	exprs := []field.AssignExpr{sysDictType.UpdatedBy.Value(operationID)}
	query.ExprAppendSelf(&exprs, req.TypeCode, sysDictType.TypeCode.Value)
	query.ExprAppendSelf(&exprs, req.TypeName, sysDictType.TypeName.Value)
	query.ExprAppendSelf(&exprs, req.IsEnabled, sysDictType.IsEnabled.Value)
	query.ExprAppendSelf(&exprs, req.SortOrder, sysDictType.SortOrder.Value)
	query.ExprAppendSelf(&exprs, req.Remark, sysDictType.Remark.Value)

	info, err := sysDictType.Where(sysDictType.ID.Eq(req.ID)).UpdateSimple(exprs...)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return res.FailMsg("类型编码已存在")
		}
		return res.FailDefault
	}
	if info.RowsAffected == 0 {
		return res.FailMsg("无权限或数据不存在")
	}
	return nil
}

// @Summary 批量删除字典类型
// @Remark 根据 ID 列表批量删除字典类型及其关联的所有字典项
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictTypeBatchDelete true "批量删除参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/dict/type/del [post]
func (h *SysDictHandler) TypeDel(ctx *handler.Ctx, req *ReqDictTypeBatchDelete) error {
	ids := slices_utils.Distinct(req.IDs)
	permissionQuery, err := h.DataPermission.BuildDeletePermissionQuery(ctx, models.SysDictType{}.TableName())
	if err != nil {
		ctx.L().Error("apply dict type delete permission failed", zap.Error(err), zap.Uint64s("ids", ids))
		return res.FailDefault
	}

	sysDictType := permissionQuery.SysDictType
	allowedTypes, err := sysDictType.
		Select(sysDictType.ID).
		Where(sysDictType.ID.In(ids...)).
		Find()
	if err != nil {
		ctx.L().Error("校验字典类型删除权限失败", zap.Error(err), zap.Uint64s("ids", ids))
		return res.FailDefault
	}
	if len(allowedTypes) != len(ids) {
		return res.FailMsg("无权限或数据不存在")
	}

	err = query.Q.Transaction(func(tx *query.Query) error {
		sysDictEntry := tx.SysDictEntry
		_, err = sysDictEntry.
			Where(sysDictEntry.SysDictTypeId.In(ids...)).
			Delete()
		if err != nil {
			return err
		}
		sysDictTypeSub := tx.SysDictType
		_, err = sysDictTypeSub.
			Where(sysDictTypeSub.ID.In(ids...)).
			Delete()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		ctx.L().Error("批量删除字典类型失败", zap.Error(err), zap.Uint64s("ids", ids))
		return res.FailDefault
	}
	return nil
}

// --- 字典数据项 (DictEntry) ---

type ReqDictEntryCreate struct {
	LabelComponent string `json:"labelComponent" change:"标签组件" binding:"omitempty,max=255" binding_msg:"max=显示标签组件最多255位"`
	EntryLabel     string `json:"entryLabel" change:"显示标签" binding:"required,max=255" binding_msg:"required=显示标签不能为空,max=显示标签最多255位"`
	EntryValue     string `json:"entryValue" change:"数据值" binding:"required,max=255" binding_msg:"required=数据值不能为空,max=数据值最多255位"`
	LanguageCode   string `json:"languageCode" change:"语言" binding:"max=32" binding_msg:"max=语言代码最多32位"`
	SysDictTypeId  uint64 `json:"sysDictTypeId" change:"字典类型" binding:"required" binding_msg:"required=字典类型ID不能为空"`
	SortOrder      int32  `json:"sortOrder" change:"排序"`
	IsEnabled      bool   `json:"isEnabled" change:"启用状态"`
	Remark         string `json:"remark" change:"备注" binding:"max=255" binding_msg:"max=备注最多255位"`
}

type ReqDictEntryUpdate struct {
	ID             *uint64                  `json:"id"`
	LabelComponent *string                  `json:"labelComponent" binding:"omitempty,max=255" binding_msg:"max=显示标签组件最多255位"`
	EntryLabel     *string                  `json:"entryLabel" binding:"omitempty,max=255" binding_msg:"max=显示标签最多255位"`
	EntryValue     *string                  `json:"entryValue" binding:"omitempty,max=255" binding_msg:"max=数据值最多255位"`
	LanguageCode   *string                  `json:"languageCode" binding:"omitempty,max=32" binding_msg:"max=语言代码最多32位"`
	SysDictTypeId  *uint64                  `json:"sysDictTypeId"`
	SortOrder      *int32                   `json:"sortOrder"`
	IsEnabled      *bool                    `json:"isEnabled"`
	Remark         *string                  `json:"remark" binding:"omitempty,max=255" binding_msg:"max=备注最多255位"`
	Updates        []ReqDictEntryUpdateItem `json:"updates" change:"更新列表"`
}

type ReqDictEntryUpdateItem struct {
	ID             uint64  `json:"id" binding:"required" binding_msg:"required=请求错误"`
	LabelComponent *string `json:"labelComponent" change:"标签组件" binding:"omitempty,max=255" binding_msg:"max=显示标签组件最多255位"`
	EntryLabel     *string `json:"entryLabel" change:"显示标签" binding:"omitempty,max=255" binding_msg:"max=显示标签最多255位"`
	EntryValue     *string `json:"entryValue" change:"数据值" binding:"omitempty,max=255" binding_msg:"max=数据值最多255位"`
	LanguageCode   *string `json:"languageCode" change:"语言" binding:"omitempty,max=32" binding_msg:"max=语言代码最多32位"`
	SysDictTypeId  *uint64 `json:"sysDictTypeId" change:"字典类型"`
	SortOrder      *int32  `json:"sortOrder" change:"排序"`
	IsEnabled      *bool   `json:"isEnabled" change:"启用状态"`
	Remark         *string `json:"remark" change:"备注" binding:"omitempty,max=255" binding_msg:"max=备注最多255位"`
}

type ReqDictEntryBatchDelete struct {
	IDs []uint64 `json:"ids" binding:"required,min=1" binding_msg:"required=请选择字典项,min=至少选择一项"`
}

type ReqDictEntryBatchCopy struct {
	EntryIds     []uint64 `json:"entryIds" change:"源字典项" binding:"required,min=1" binding_msg:"required=请选择字典项,min=至少选择一项"`
	TargetTypeId uint64   `json:"targetTypeId" change:"目标字典类型" binding:"required" binding_msg:"required=目标字典类型不能为空"`
}

type ReqDictEntryMatch struct {
	Codes []string `json:"codes" binding:"omitempty,min=1,dive,required,max=128" binding_msg:"min=至少选择一个字典编码,required=字典类型编码不能为空,max=字典类型编码最多128位"`
}

type RespDictEntryByCode struct {
	ID             uint64 `json:"id"`
	LabelComponent string `json:"labelComponent"`
	EntryLabel     string `json:"entryLabel"`
	EntryValue     string `json:"entryValue"`
}

type RespDictEntryMatch map[string][]*RespDictEntryByCode

// @Summary 通过字典编码批量获取启用字典项
// @Remark 根据字典类型编码批量查询启用字典项；若字典项配置了语言条目编码，则按当前请求语言替换显示标签
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictEntryMatch true "字典类型编码"
// @Success 200 {object} res.Response{data=RespDictEntryMatch} "成功"
// @Router /api/sys/dict/entry/match [post]
func (h *SysDictHandler) EntryMatch(ctx *handler.Ctx, req *ReqDictEntryMatch) (*RespDictEntryMatch, error) {
	codes := normalizeDictEntryMatchCodes(req)
	if len(codes) == 0 {
		return nil, res.FailMsg("字典类型编码不能为空")
	}

	permissionQuery, err := h.DataPermission.BuildReadPermissionQuery(ctx, models.SysDictType{}.TableName())
	if err != nil {
		ctx.L().Error("apply dict type read permission failed", zap.Error(err), zap.Strings("codes", codes))
		return nil, res.FailDefault
	}

	items := make(RespDictEntryMatch, len(codes))
	for _, code := range codes {
		items[code] = []*RespDictEntryByCode{}
	}

	sysDictType := permissionQuery.SysDictType
	dictTypes, err := sysDictType.
		Select(sysDictType.ID, sysDictType.TypeCode).
		Where(sysDictType.TypeCode.In(codes...), sysDictType.IsEnabled.Is(true)).
		Find()
	if err != nil {
		ctx.L().Error("query dict types by code failed", zap.Error(err), zap.Strings("codes", codes))
		return nil, res.FailDefault
	}
	if len(dictTypes) == 0 {
		return &items, nil
	}

	typeIDs := make([]uint64, 0, len(dictTypes))
	codeByTypeID := make(map[uint64]string, len(dictTypes))
	for _, dictType := range dictTypes {
		typeIDs = append(typeIDs, dictType.ID)
		codeByTypeID[dictType.ID] = dictType.TypeCode
	}

	sysDictEntry := h.Q.SysDictEntry
	entries, err := sysDictEntry.
		Where(sysDictEntry.SysDictTypeId.In(typeIDs...), sysDictEntry.IsEnabled.Is(true)).
		Order(sysDictEntry.SortOrder.Asc(), sysDictEntry.ID.Asc()).
		Find()
	if err != nil {
		ctx.L().Error("query dict entries by codes failed", zap.Error(err), zap.Strings("codes", codes), zap.Uint64s("typeIds", typeIDs))
		return nil, res.FailDefault
	}
	if len(entries) == 0 {
		return &items, nil
	}

	translationMap, err := queryDictEntryTranslationMap(h.Q, ctx, entries)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		code, ok := codeByTypeID[entry.SysDictTypeId]
		if !ok {
			continue
		}
		entryLabel := entry.EntryLabel
		if translation, ok := translationMap[strings.TrimSpace(entry.LanguageCode)]; ok {
			entryLabel = translation
		}
		items[code] = append(items[code], &RespDictEntryByCode{
			ID:             entry.ID,
			LabelComponent: entry.LabelComponent,
			EntryLabel:     entryLabel,
			EntryValue:     entry.EntryValue,
		})
	}
	return &items, nil
}

func normalizeDictEntryMatchCodes(req *ReqDictEntryMatch) []string {
	codeSet := make(map[string]struct{}, len(req.Codes))
	codes := make([]string, 0, len(req.Codes))
	appendCode := func(code string) {
		code = strings.TrimSpace(code)
		if code == "" {
			return
		}
		if _, exists := codeSet[code]; exists {
			return
		}
		codeSet[code] = struct{}{}
		codes = append(codes, code)
	}

	for _, code := range req.Codes {
		appendCode(code)
	}
	return codes
}

func queryDictEntryTranslationMap(q *query.Query, ctx *handler.Ctx, entries []*models.SysDictEntry) (map[string]string, error) {
	languageCodes := make([]string, 0, len(entries))
	languageCodeSet := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		languageCode := strings.TrimSpace(entry.LanguageCode)
		if languageCode == "" {
			continue
		}
		if _, exists := languageCodeSet[languageCode]; exists {
			continue
		}
		languageCodeSet[languageCode] = struct{}{}
		languageCodes = append(languageCodes, languageCode)
	}
	language := strings.TrimSpace(ctx.Language)
	if len(languageCodes) == 0 || language == "" {
		return map[string]string{}, nil
	}

	sysLanguageType := q.SysLanguageType
	languageType, err := sysLanguageType.
		Select(sysLanguageType.ID).
		Where(sysLanguageType.TypeCode.Eq(language), sysLanguageType.IsEnabled.Is(true)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return map[string]string{}, nil
		}
		ctx.L().Error("query language type failed", zap.Error(err), zap.String("language", language))
		return nil, res.FailDefault
	}

	sysLanguageEntry := q.SysLanguageEntry
	languageEntries, err := sysLanguageEntry.
		Select(sysLanguageEntry.EntryCode, sysLanguageEntry.EntryValue).
		Where(
			sysLanguageEntry.SysLanguageTypeId.Eq(languageType.ID),
			sysLanguageEntry.EntryCode.In(languageCodes...),
			sysLanguageEntry.IsEnabled.Is(true),
		).
		Find()
	if err != nil {
		ctx.L().Error("query language entries failed", zap.Error(err), zap.String("language", language), zap.Strings("entryCodes", languageCodes))
		return nil, res.FailDefault
	}

	translationMap := make(map[string]string, len(languageEntries))
	for _, entry := range languageEntries {
		translationMap[entry.EntryCode] = entry.EntryValue
	}
	return translationMap, nil
}

// @Summary 获取字典数据项分页列表
// @Remark 分页查询字典数据项信息
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[models.SysDictEntry]} "成功"
// @Router /api/dict/entry/list [post]
func (h *SysDictHandler) EntryList(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[models.SysDictEntry], error) {
	readableTypeIDs, err := queryAllowedDictTypeIDs(ctx, h.DataPermission.BuildReadPermissionQuery)
	if err != nil {
		ctx.L().Error("apply dict type read permission failed", zap.Error(err))
		return nil, res.FailDefault
	}
	if len(readableTypeIDs) == 0 {
		return &gormc.PagingResult[models.SysDictEntry]{
			Items: []*models.SysDictEntry{},
			Total: 0,
		}, nil
	}
	if err := applyDictEntryTypeIDPageFilter(req, readableTypeIDs); err != nil {
		ctx.L().Error("apply dict entry type read permission failed", zap.Error(err), zap.Uint64s("typeIds", readableTypeIDs))
		return nil, res.FailDefault
	}

	pagination, err := h.Q.SysDictEntry.PageWithPaging(req)
	if err != nil {
		return nil, res.FailDefault
	}
	if len(pagination.Items) == 0 {
		return pagination, nil
	}

	typeIDSet := make(map[uint64]struct{}, len(pagination.Items))
	typeIDs := make([]uint64, 0, len(pagination.Items))
	for _, item := range pagination.Items {
		if _, exists := typeIDSet[item.SysDictTypeId]; exists {
			continue
		}
		typeIDSet[item.SysDictTypeId] = struct{}{}
		typeIDs = append(typeIDs, item.SysDictTypeId)
	}
	if len(typeIDs) == 0 {
		return pagination, nil
	}

	sysDictType := h.Q.SysDictType
	typeList, err := sysDictType.
		Select(sysDictType.ID, sysDictType.TypeCode, sysDictType.TypeName).
		Where(sysDictType.ID.In(typeIDs...)).
		Find()
	if err != nil {
		ctx.L().Error("查询字典类型失败", zap.Error(err), zap.Uint64s("typeIds", typeIDs))
		return nil, res.FailDefault
	}
	typeMap := make(map[uint64]*models.SysDictType, len(typeList))
	for _, item := range typeList {
		typeMap[item.ID] = item
	}
	for _, item := range pagination.Items {
		item.SysDictType = typeMap[item.SysDictTypeId]
	}
	return pagination, nil
}

// @Summary 创建字典数据项
// @Remark 创建新的字典数据项
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictEntryCreate true "字典数据项创建参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/dict/entry/create [post]
func (h *SysDictHandler) EntryCreate(ctx *handler.Ctx, req *ReqDictEntryCreate) error {
	operationID := ctx.SessionInfo.Id
	permissionQuery, err := h.DataPermission.BuildWritePermissionQuery(ctx, models.SysDictType{}.TableName())
	if err != nil {
		ctx.L().Error("apply dict type write permission failed", zap.Error(err), zap.Uint64("typeId", req.SysDictTypeId))
		return res.FailDefault
	}

	sysDictType := permissionQuery.SysDictType
	_, err = sysDictType.
		Select(sysDictType.ID).
		Where(sysDictType.ID.Eq(req.SysDictTypeId)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res.FailMsg("无权限或字典类型不存在")
		}
		return res.FailDefault
	}

	err = h.Q.SysDictEntry.Create(&models.SysDictEntry{
		OperatorID: mixin.OperatorID{
			CreatedBy: mixin.CreatedBy{CreatedBy: operationID},
			UpdatedBy: mixin.UpdatedBy{UpdatedBy: operationID},
		},
		SortOrder:      mixin.SortOrder{SortOrder: req.SortOrder},
		IsEnabled:      mixin.IsEnabled{IsEnabled: req.IsEnabled},
		Remark:         mixin.Remark{Remark: req.Remark},
		LabelComponent: req.LabelComponent,
		EntryLabel:     req.EntryLabel,
		EntryValue:     req.EntryValue,
		LanguageCode:   req.LanguageCode,
		SysDictTypeId:  req.SysDictTypeId,
	})
	if err != nil {
		return res.FailDefault
	}
	return nil
}

// @Summary 更新字典数据项
// @Remark 根据 ID 更新字典数据项信息
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictEntryUpdate true "字典数据项更新参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/dict/entry/update [post]
func (h *SysDictHandler) EntryUpdate(ctx *handler.Ctx, req *ReqDictEntryUpdate) error {
	operationID := ctx.SessionInfo.Id
	typeWriteQuery, err := h.DataPermission.BuildWritePermissionQuery(ctx, models.SysDictType{}.TableName())
	if err != nil {
		ctx.L().Error("apply dict type write permission failed", zap.Error(err))
		return res.FailDefault
	}

	if len(req.Updates) > 0 {
		for _, item := range req.Updates {
			if err := updateDictEntry(h.Q, operationID, typeWriteQuery, &item); err != nil {
				return err
			}
		}
		return nil
	}
	if req.ID == nil {
		return res.FailMsg("请求错误")
	}
	return updateDictEntry(h.Q, operationID, typeWriteQuery, &ReqDictEntryUpdateItem{
		ID:             *req.ID,
		LabelComponent: req.LabelComponent,
		EntryLabel:     req.EntryLabel,
		EntryValue:     req.EntryValue,
		LanguageCode:   req.LanguageCode,
		SysDictTypeId:  req.SysDictTypeId,
		SortOrder:      req.SortOrder,
		IsEnabled:      req.IsEnabled,
		Remark:         req.Remark,
	})
}

func updateDictEntry(q *query.Query, operationID uint64, typeWriteQuery *query.Query, req *ReqDictEntryUpdateItem) error {
	entry, err := q.SysDictEntry.
		Select(q.SysDictEntry.ID, q.SysDictEntry.SysDictTypeId).
		Where(q.SysDictEntry.ID.Eq(req.ID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res.FailMsg("无权限或数据不存在")
		}
		return res.FailDefault
	}

	if err := ensureDictTypeAllowed(typeWriteQuery, entry.SysDictTypeId); err != nil {
		return err
	}
	if req.SysDictTypeId != nil && *req.SysDictTypeId != entry.SysDictTypeId {
		if err := ensureDictTypeAllowed(typeWriteQuery, *req.SysDictTypeId); err != nil {
			return err
		}
	}

	sysDictEntry := q.SysDictEntry
	exprs := []field.AssignExpr{sysDictEntry.UpdatedBy.Value(operationID)}
	query.ExprAppendSelf(&exprs, req.LabelComponent, sysDictEntry.LabelComponent.Value)
	query.ExprAppendSelf(&exprs, req.EntryLabel, sysDictEntry.EntryLabel.Value)
	query.ExprAppendSelf(&exprs, req.EntryValue, sysDictEntry.EntryValue.Value)
	query.ExprAppendSelf(&exprs, req.LanguageCode, sysDictEntry.LanguageCode.Value)
	query.ExprAppendSelf(&exprs, req.SysDictTypeId, sysDictEntry.SysDictTypeId.Value)
	query.ExprAppendSelf(&exprs, req.SortOrder, sysDictEntry.SortOrder.Value)
	query.ExprAppendSelf(&exprs, req.IsEnabled, sysDictEntry.IsEnabled.Value)
	query.ExprAppendSelf(&exprs, req.Remark, sysDictEntry.Remark.Value)

	info, err := sysDictEntry.Where(sysDictEntry.ID.Eq(req.ID)).UpdateSimple(exprs...)
	if err != nil {
		return res.FailDefault
	}
	if info.RowsAffected == 0 {
		return res.FailMsg("无权限或数据不存在")
	}
	return nil
}

func ensureDictTypeAllowed(permissionQuery *query.Query, typeID uint64) error {
	sysDictType := permissionQuery.SysDictType
	_, err := sysDictType.
		Select(sysDictType.ID).
		Where(sysDictType.ID.Eq(typeID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res.FailMsg("无权限或字典类型不存在")
		}
		return res.FailDefault
	}
	return nil
}

// @Summary 批量删除字典数据项
// @Remark 根据 ID 列表批量删除字典数据项
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictEntryBatchDelete true "批量删除参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/dict/entry/del [post]
func (h *SysDictHandler) EntryDel(ctx *handler.Ctx, req *ReqDictEntryBatchDelete) error {
	ids := slices_utils.Distinct(req.IDs)
	sysDictEntry := h.Q.SysDictEntry
	entries, err := sysDictEntry.
		Select(sysDictEntry.ID, sysDictEntry.SysDictTypeId).
		Where(sysDictEntry.ID.In(ids...)).
		Find()
	if err != nil {
		ctx.L().Error("查询字典项删除范围失败", zap.Error(err), zap.Uint64s("ids", ids))
		return res.FailDefault
	}
	if len(entries) != len(ids) {
		return res.FailMsg("无权限或数据不存在")
	}

	typeIDSet := make(map[uint64]struct{}, len(entries))
	typeIDs := make([]uint64, 0, len(entries))
	for _, entry := range entries {
		if _, exists := typeIDSet[entry.SysDictTypeId]; exists {
			continue
		}
		typeIDSet[entry.SysDictTypeId] = struct{}{}
		typeIDs = append(typeIDs, entry.SysDictTypeId)
	}
	deleteIDSet, err := queryAllowedDictTypeIDSet(ctx, typeIDs, h.DataPermission.BuildDeletePermissionQuery)
	if err != nil {
		ctx.L().Error("apply dict type delete permission failed", zap.Error(err), zap.Uint64s("typeIds", typeIDs))
		return res.FailDefault
	}
	if len(deleteIDSet) != len(typeIDs) {
		return res.FailMsg("无权限或数据不存在")
	}

	info, err := sysDictEntry.Where(sysDictEntry.ID.In(ids...)).Delete()
	if err != nil {
		return res.FailDefault
	}
	if info.RowsAffected != int64(len(ids)) {
		return res.FailMsg("无权限或数据不存在")
	}
	return nil
}

// @Summary 批量复制字典数据项
// @Remark 将选中的字典数据项批量复制到指定字典类型下（不支持复制到同一类型）
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictEntryBatchCopy true "批量复制参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/dict/entry/batch/copy [post]
func (h *SysDictHandler) EntryBatchCopy(ctx *handler.Ctx, req *ReqDictEntryBatchCopy) error {
	writePermissionQuery, err := h.DataPermission.BuildWritePermissionQuery(ctx, models.SysDictType{}.TableName())
	if err != nil {
		ctx.L().Error("apply dict type write permission failed", zap.Error(err), zap.Uint64("targetTypeId", req.TargetTypeId))
		return res.FailDefault
	}

	if err := ensureDictTypeAllowed(writePermissionQuery, req.TargetTypeId); err != nil {
		return err
	}

	entryIDs := slices_utils.Distinct(req.EntryIds)
	sourceEntries, err := h.Q.SysDictEntry.
		Where(h.Q.SysDictEntry.ID.In(entryIDs...)).
		Find()
	if err != nil {
		ctx.L().Error("查询源字典项失败", zap.Error(err), zap.Uint64s("entryIds", entryIDs))
		return res.FailDefault
	}
	if len(sourceEntries) != len(entryIDs) {
		return res.FailMsg("无权限或源字典项不存在")
	}

	sourceTypeIDSet := make(map[uint64]struct{}, len(sourceEntries))
	sourceTypeIDs := make([]uint64, 0, len(sourceEntries))
	for _, entry := range sourceEntries {
		if _, exists := sourceTypeIDSet[entry.SysDictTypeId]; exists {
			continue
		}
		sourceTypeIDSet[entry.SysDictTypeId] = struct{}{}
		sourceTypeIDs = append(sourceTypeIDs, entry.SysDictTypeId)
	}
	readIDSet, err := queryAllowedDictTypeIDSet(ctx, sourceTypeIDs, h.DataPermission.BuildReadPermissionQuery)
	if err != nil {
		ctx.L().Error("apply source dict type read permission failed", zap.Error(err), zap.Uint64s("typeIds", sourceTypeIDs))
		return res.FailDefault
	}
	if len(readIDSet) != len(sourceTypeIDs) {
		return res.FailMsg("无权限或源字典项不存在")
	}

	var newEntries []*models.SysDictEntry
	for _, entry := range sourceEntries {
		newEntries = append(newEntries, &models.SysDictEntry{
			LabelComponent: entry.LabelComponent,
			EntryLabel:     entry.EntryLabel,
			EntryValue:     entry.EntryValue,
			LanguageCode:   entry.LanguageCode,
			SysDictTypeId:  req.TargetTypeId,
			SortOrder:      mixin.SortOrder{SortOrder: entry.SortOrder.SortOrder},
			IsEnabled:      mixin.IsEnabled{IsEnabled: entry.IsEnabled.IsEnabled},
			Remark:         mixin.Remark{Remark: entry.Remark.Remark},
		})
	}

	err = h.Q.SysDictEntry.Create(newEntries...)
	if err != nil {
		ctx.L().Error("批量复制字典项失败", zap.Error(err), zap.Uint64s("entryIds", entryIDs), zap.Uint64("targetTypeId", req.TargetTypeId))
		return res.FailDefault
	}

	return nil
}
