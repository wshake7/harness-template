package logic

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"errors"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc"
	"orm-crud/gormc/mixin"

	"gorm.io/gen/field"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SysLanguageHandler struct {
	Q *query.Query
}

func NewSysLanguageHandler(q *query.Query) *SysLanguageHandler {
	return &SysLanguageHandler{Q: q}
}

// --- 语言类型 (LanguageType) ---

type ReqLangTypeCreate struct {
	TypeCode  string `json:"typeCode" change:"语言编码" binding:"required,max=128" binding_msg:"required=语言编码不能为空,max=语言编码最多128位"`
	TypeName  string `json:"typeName" change:"语言名称" binding:"required,max=255" binding_msg:"required=语言名称不能为空,max=语言名称最多255位"`
	IsDefault bool   `json:"isDefault" change:"默认语言"`
	IsEnabled bool   `json:"isEnabled" change:"启用状态"`
	SortOrder int32  `json:"sortOrder" change:"排序"`
}

type ReqLangTypeUpdate struct {
	ID        uint64  `json:"id" binding:"required" binding_msg:"required=请求错误"`
	TypeCode  *string `json:"typeCode" change:"语言编码" binding:"omitempty,max=128" binding_msg:"max=语言编码最多128位"`
	TypeName  *string `json:"typeName" change:"语言名称" binding:"omitempty,max=255" binding_msg:"max=语言名称最多255位"`
	IsDefault *bool   `json:"isDefault" change:"默认语言"`
	IsEnabled *bool   `json:"isEnabled" change:"启用状态"`
	SortOrder *int32  `json:"sortOrder" change:"排序"`
}

type ReqLangTypeBatchDelete struct {
	IDs []uint64 `json:"ids" binding:"required,min=1" binding_msg:"required=请选择语言类型,min=至少选择一项"`
}

// @Summary 获取语言类型分页列表
// @Remark 分页查询语言类型信息
// @Tags Language
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[models.SysLanguageType]} "成功"
// @Router /api/sys/language/type/list [post]
func (h *SysLanguageHandler) TypeList(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[models.SysLanguageType], error) {
	pagination, err := h.Q.SysLanguageType.PageWithPaging(req)
	if err != nil {
		return nil, res.FailDefault
	}
	return pagination, nil
}

// @Summary 创建语言类型
// @Remark 创建新的语言类型，可指定默认语言
// @Tags Language
// @Accept json
// @Produce json
// @Param req body ReqLangTypeCreate true "语言类型创建参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/language/type/create [post]
func (h *SysLanguageHandler) TypeCreate(ctx *handler.Ctx, req *ReqLangTypeCreate) error {
	if req.IsDefault {
		sysLanguageType := h.Q.SysLanguageType
		if _, err := sysLanguageType.Where(sysLanguageType.IsDefault.Is(true)).Update(sysLanguageType.IsDefault, false); err != nil {
			ctx.L().Error("重置默认语言失败", zap.Error(err))
			return res.FailDefault
		}
	}
	operationID := ctx.SessionInfo.Id
	sysLanguageType := h.Q.SysLanguageType
	err := sysLanguageType.Create(&models.SysLanguageType{
		OperatorID: mixin.OperatorID{
			CreatedBy: mixin.CreatedBy{CreatedBy: operationID},
			UpdatedBy: mixin.UpdatedBy{UpdatedBy: operationID},
		},
		IsEnabled: mixin.IsEnabled{IsEnabled: req.IsEnabled},
		SortOrder: mixin.SortOrder{SortOrder: req.SortOrder},
		TypeCode:  req.TypeCode,
		TypeName:  req.TypeName,
		IsDefault: req.IsDefault,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return res.FailMsg("语言编码已存在")
		}
		return res.FailDefault
	}
	return nil
}

// @Summary 更新语言类型
// @Remark 根据 ID 更新语言类型信息
// @Tags Language
// @Accept json
// @Produce json
// @Param req body ReqLangTypeUpdate true "语言类型更新参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/language/type/update [post]
func (h *SysLanguageHandler) TypeUpdate(ctx *handler.Ctx, req *ReqLangTypeUpdate) error {
	sysLanguageType := h.Q.SysLanguageType
	if req.IsEnabled != nil && !*req.IsEnabled {
		target, err := sysLanguageType.Select(sysLanguageType.IsDefault).Where(sysLanguageType.ID.Eq(req.ID)).First()
		if err == nil && target.IsDefault {
			return res.FailMsg("默认语言不能停用")
		}
	}
	if req.IsDefault != nil && *req.IsDefault {
		if _, err := sysLanguageType.Where(sysLanguageType.ID.Neq(req.ID), sysLanguageType.IsDefault.Is(true)).Update(sysLanguageType.IsDefault, false); err != nil {
			ctx.L().Error("重置默认语言失败", zap.Error(err))
			return res.FailDefault
		}
	}
	operationID := ctx.SessionInfo.Id
	exprs := []field.AssignExpr{sysLanguageType.UpdatedBy.Value(operationID)}
	query.ExprAppendSelf(&exprs, req.TypeCode, sysLanguageType.TypeCode.Value)
	query.ExprAppendSelf(&exprs, req.TypeName, sysLanguageType.TypeName.Value)
	query.ExprAppendSelf(&exprs, req.IsDefault, sysLanguageType.IsDefault.Value)
	query.ExprAppendSelf(&exprs, req.IsEnabled, sysLanguageType.IsEnabled.Value)
	query.ExprAppendSelf(&exprs, req.SortOrder, sysLanguageType.SortOrder.Value)

	_, err := sysLanguageType.Where(sysLanguageType.ID.Eq(req.ID)).UpdateSimple(exprs...)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return res.FailMsg("语言编码已存在")
		}
		return res.FailDefault
	}
	return nil
}

// @Summary 批量删除语言类型
// @Remark 根据 ID 列表批量删除语言类型及其关联条目
// @Tags Language
// @Accept json
// @Produce json
// @Param req body ReqLangTypeBatchDelete true "批量删除参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/language/type/del [post]
func (h *SysLanguageHandler) TypeDel(ctx *handler.Ctx, req *ReqLangTypeBatchDelete) error {
	sysLanguageType := h.Q.SysLanguageType
	defaultLang, err := sysLanguageType.Select(sysLanguageType.ID).Where(sysLanguageType.ID.In(req.IDs...), sysLanguageType.IsDefault.Is(true)).First()
	if err == nil && defaultLang != nil {
		return res.FailMsg("默认语言不能删除")
	}
	err = h.Q.Transaction(func(tx *query.Query) error {
		sysLanguageEntry := tx.SysLanguageEntry
		if _, err = sysLanguageEntry.Where(sysLanguageEntry.SysLanguageTypeId.In(req.IDs...)).Delete(); err != nil {
			return err
		}
		sysLanguageTypeSub := tx.SysLanguageType
		if _, err = sysLanguageTypeSub.Where(sysLanguageTypeSub.ID.In(req.IDs...)).Delete(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		ctx.L().Error("批量删除语言类型失败", zap.Error(err), zap.Uint64s("ids", req.IDs))
		return res.FailDefault
	}
	return nil
}

// --- 语言条目 (LanguageEntry) ---

type ReqLangEntryCreate struct {
	EntryCode         string `json:"entryCode" change:"条目编码" binding:"required,max=128" binding_msg:"required=条目编码不能为空,max=条目编码最多128位"`
	EntryValue        string `json:"entryValue" change:"语言值" binding:"required,max=255" binding_msg:"required=语言值不能为空,max=语言值最多255位"`
	SysLanguageTypeId uint64 `json:"sysLanguageTypeId" change:"语言类型" binding:"required" binding_msg:"required=语言类型ID不能为空"`
	SortOrder         int32  `json:"sortOrder" change:"排序"`
	IsEnabled         bool   `json:"isEnabled" change:"启用状态"`
	Remark            string `json:"remark" change:"备注" binding:"max=255" binding_msg:"max=备注最多255位"`
}

type ReqLangEntryUpdate struct {
	ID                *uint64                  `json:"id"`
	EntryCode         *string                  `json:"entryCode" binding:"omitempty,max=128" binding_msg:"max=条目编码最多128位"`
	EntryValue        *string                  `json:"entryValue" binding:"omitempty,max=255" binding_msg:"max=语言值最多255位"`
	SysLanguageTypeId *uint64                  `json:"sysLanguageTypeId"`
	SortOrder         *int32                   `json:"sortOrder"`
	IsEnabled         *bool                    `json:"isEnabled"`
	Remark            *string                  `json:"remark" binding:"omitempty,max=255" binding_msg:"max=备注最多255位"`
	Updates           []ReqLangEntryUpdateItem `json:"updates" change:"更新列表"`
}

type ReqLangEntryUpdateItem struct {
	ID                uint64  `json:"id" binding:"required" binding_msg:"required=请求错误"`
	EntryCode         *string `json:"entryCode" change:"条目编码" binding:"omitempty,max=128" binding_msg:"max=条目编码最多128位"`
	EntryValue        *string `json:"entryValue" change:"语言值" binding:"omitempty,max=255" binding_msg:"max=语言值最多255位"`
	SysLanguageTypeId *uint64 `json:"sysLanguageTypeId" change:"语言类型"`
	SortOrder         *int32  `json:"sortOrder" change:"排序"`
	IsEnabled         *bool   `json:"isEnabled" change:"启用状态"`
	Remark            *string `json:"remark" change:"备注" binding:"omitempty,max=255" binding_msg:"max=备注最多255位"`
}

type ReqLangEntryBatchDelete struct {
	IDs []uint64 `json:"ids" binding:"required,min=1" binding_msg:"required=请选择语言条目,min=至少选择一项"`
}

type ReqLangEntryBatchCreate struct {
	EntryCode string            `json:"entryCode" change:"条目编码" binding:"required,max=128" binding_msg:"required=条目编码不能为空,max=条目编码最多128位"`
	Values    map[string]string `json:"values" change:"语言值"`
	SortOrder int32             `json:"sortOrder" change:"排序"`
	IsEnabled bool              `json:"isEnabled" change:"启用状态"`
}

// @Summary 获取语言条目分页列表
// @Remark 分页查询语言条目信息
// @Tags Language
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[models.SysLanguageEntry]} "成功"
// @Router /api/sys/language/entry/list [post]
func (h *SysLanguageHandler) EntryList(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[models.SysLanguageEntry], error) {
	pagination, err := h.Q.SysLanguageEntry.PageWithPaging(req)
	if err != nil {
		return nil, res.FailDefault
	}
	return pagination, nil
}

// @Summary 创建语言条目
// @Remark 在指定语言类型下创建语言条目
// @Tags Language
// @Accept json
// @Produce json
// @Param req body ReqLangEntryCreate true "语言条目创建参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/language/entry/create [post]
func (h *SysLanguageHandler) EntryCreate(ctx *handler.Ctx, req *ReqLangEntryCreate) error {
	operationID := ctx.SessionInfo.Id
	sysLanguageType := h.Q.SysLanguageType
	_, err := sysLanguageType.Select(sysLanguageType.ID).Where(sysLanguageType.ID.Eq(req.SysLanguageTypeId)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res.FailMsg("语言类型不存在")
		}
		return res.FailDefault
	}
	err = h.Q.SysLanguageEntry.Create(&models.SysLanguageEntry{
		OperatorID: mixin.OperatorID{
			CreatedBy: mixin.CreatedBy{CreatedBy: operationID},
			UpdatedBy: mixin.UpdatedBy{UpdatedBy: operationID},
		},
		SortOrder:         mixin.SortOrder{SortOrder: req.SortOrder},
		IsEnabled:         mixin.IsEnabled{IsEnabled: req.IsEnabled},
		Remark:            mixin.Remark{Remark: req.Remark},
		EntryCode:         req.EntryCode,
		EntryValue:        req.EntryValue,
		SysLanguageTypeId: req.SysLanguageTypeId,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return res.FailMsg("该语言下已存在相同编码的条目")
		}
		return res.FailDefault
	}
	return nil
}

// @Summary 更新语言条目
// @Remark 更新单个或批量更新语言条目
// @Tags Language
// @Accept json
// @Produce json
// @Param req body ReqLangEntryUpdate true "语言条目更新参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/language/entry/update [post]
func (h *SysLanguageHandler) EntryUpdate(ctx *handler.Ctx, req *ReqLangEntryUpdate) error {
	operationID := ctx.SessionInfo.Id
	if len(req.Updates) > 0 {
		for _, item := range req.Updates {
			if err := h.updateLanguageEntry(operationID, &item); err != nil {
				return err
			}
		}
		return nil
	}
	if req.ID == nil {
		return res.FailMsg("请求错误")
	}
	return h.updateLanguageEntry(operationID, &ReqLangEntryUpdateItem{
		ID:                *req.ID,
		EntryCode:         req.EntryCode,
		EntryValue:        req.EntryValue,
		SysLanguageTypeId: req.SysLanguageTypeId,
		SortOrder:         req.SortOrder,
		IsEnabled:         req.IsEnabled,
		Remark:            req.Remark,
	})
}

func (h *SysLanguageHandler) updateLanguageEntry(operationID uint64, req *ReqLangEntryUpdateItem) error {
	if req.SysLanguageTypeId != nil {
		sysLanguageType := h.Q.SysLanguageType
		_, err := sysLanguageType.Select(sysLanguageType.ID).Where(sysLanguageType.ID.Eq(*req.SysLanguageTypeId)).First()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return res.FailMsg("语言类型不存在")
			}
			return res.FailDefault
		}
	}
	sysLanguageEntry := h.Q.SysLanguageEntry
	exprs := []field.AssignExpr{sysLanguageEntry.UpdatedBy.Value(operationID)}
	query.ExprAppendSelf(&exprs, req.EntryCode, sysLanguageEntry.EntryCode.Value)
	query.ExprAppendSelf(&exprs, req.EntryValue, sysLanguageEntry.EntryValue.Value)
	query.ExprAppendSelf(&exprs, req.SysLanguageTypeId, sysLanguageEntry.SysLanguageTypeId.Value)
	query.ExprAppendSelf(&exprs, req.SortOrder, sysLanguageEntry.SortOrder.Value)
	query.ExprAppendSelf(&exprs, req.IsEnabled, sysLanguageEntry.IsEnabled.Value)
	query.ExprAppendSelf(&exprs, req.Remark, sysLanguageEntry.Remark.Value)
	_, err := sysLanguageEntry.Where(sysLanguageEntry.ID.Eq(req.ID)).UpdateSimple(exprs...)
	if err != nil {
		return res.FailDefault
	}
	return nil
}

// @Summary 批量删除语言条目
// @Remark 根据 ID 列表批量删除语言条目
// @Tags Language
// @Accept json
// @Produce json
// @Param req body ReqLangEntryBatchDelete true "批量删除参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/language/entry/del [post]
func (h *SysLanguageHandler) EntryDel(ctx *handler.Ctx, req *ReqLangEntryBatchDelete) error {
	sysLanguageEntry := h.Q.SysLanguageEntry
	_, err := sysLanguageEntry.Where(sysLanguageEntry.ID.In(req.IDs...)).Delete()
	if err != nil {
		return res.FailDefault
	}
	return nil
}

// @Summary 批量创建语言条目
// @Remark 按全部语言类型批量创建同编码的语言条目
// @Tags Language
// @Accept json
// @Produce json
// @Param req body ReqLangEntryBatchCreate true "批量创建参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/language/entry/batch/create [post]
func (h *SysLanguageHandler) EntryBatchCreate(ctx *handler.Ctx, req *ReqLangEntryBatchCreate) error {
	operationID := ctx.SessionInfo.Id

	var typeCodes []string
	for tc := range req.Values {
		typeCodes = append(typeCodes, tc)
	}

	sysLanguageType := h.Q.SysLanguageType
	typeList, err := sysLanguageType.Select(sysLanguageType.ID, sysLanguageType.TypeCode).Where(sysLanguageType.TypeCode.In(typeCodes...)).Find()
	if err != nil {
		ctx.L().Error("查询语言类型失败", zap.Error(err))
		return res.FailDefault
	}

	typeCodeToID := make(map[string]uint64, len(typeList))
	var typeIDs []uint64
	for _, t := range typeList {
		typeCodeToID[t.TypeCode] = t.ID
		typeIDs = append(typeIDs, t.ID)
	}

	existingEntries, err := h.Q.SysLanguageEntry.Where(h.Q.SysLanguageEntry.EntryCode.Eq(req.EntryCode), h.Q.SysLanguageEntry.SysLanguageTypeId.In(typeIDs...)).Find()
	if err != nil {
		ctx.L().Error("查询已有语言条目失败", zap.Error(err))
		return res.FailDefault
	}

	existingMap := make(map[uint64]*models.SysLanguageEntry, len(existingEntries))
	for _, e := range existingEntries {
		existingMap[e.SysLanguageTypeId] = e
	}

	var createEntries []*models.SysLanguageEntry
	var updateEntries []*models.SysLanguageEntry
	for typeCode, entryValue := range req.Values {
		typeID, ok := typeCodeToID[typeCode]
		if !ok {
			continue
		}
		if existing, ok := existingMap[typeID]; ok {
			existing.EntryValue = entryValue
			existing.UpdatedBy = mixin.UpdatedBy{UpdatedBy: operationID}
			updateEntries = append(updateEntries, existing)
		} else {
			createEntries = append(createEntries, &models.SysLanguageEntry{
				OperatorID: mixin.OperatorID{
					CreatedBy: mixin.CreatedBy{CreatedBy: operationID},
					UpdatedBy: mixin.UpdatedBy{UpdatedBy: operationID},
				},
				SortOrder:         mixin.SortOrder{SortOrder: req.SortOrder},
				IsEnabled:         mixin.IsEnabled{IsEnabled: req.IsEnabled},
				EntryCode:         req.EntryCode,
				EntryValue:        entryValue,
				SysLanguageTypeId: typeID,
			})
		}
	}

	if len(createEntries) > 0 {
		if err = h.Q.SysLanguageEntry.Create(createEntries...); err != nil {
			ctx.L().Error("批量创建语言条目失败", zap.Error(err))
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return res.FailMsg("语言条目有重复")
			}
			return res.FailDefault
		}
	}

	if len(updateEntries) > 0 {
		if err = h.Q.SysLanguageEntry.Save(updateEntries...); err != nil {
			ctx.L().Error("批量更新语言条目失败", zap.Error(err))
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return res.FailMsg("语言条目有重复")
			}
			return res.FailDefault
		}
	}

	return nil
}
