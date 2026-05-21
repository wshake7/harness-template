package logic

import (
	"errors"
	"strings"

	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"go-common/utils/passwd"
	"go-common/utils/slices_utils"
	"go-common/utils/str"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc"
	"orm-crud/gormc/mixin"

	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type SysUserHandler struct {
	Q *query.Query
}

func NewSysUserHandler(q *query.Query) *SysUserHandler {
	return &SysUserHandler{Q: q}
}

type RespSysUser struct {
	models.SysUser
	CanWrite  bool `json:"canWrite"`
	CanDelete bool `json:"canDelete"`
}

type ReqSysUserCreate struct {
	Username     string `json:"username" change:"用户名" binding:"notblank,max=64" binding_msg:"notblank=用户名不能为空,max=用户名最多64位"`
	Nickname     string `json:"nickname" change:"昵称" binding:"max=64" binding_msg:"max=昵称最多64位"`
	Password     string `json:"password" binding:"notblank,min=6,max=255" binding_msg:"notblank=密码不能为空,min=密码至少6位,max=密码最多255位"`
	LanguageCode string `json:"languageCode" change:"语言" binding:"max=32" binding_msg:"max=语言代码最多32位"`
	IsEnabled    bool   `json:"isEnabled" change:"启用状态"`
	Remark       string `json:"remark" change:"备注" binding:"max=255" binding_msg:"max=备注最多255位"`
}

type ReqSysUserUpdate struct {
	ID           uint64  `json:"id" binding:"required" binding_msg:"required=请求错误"`
	Username     *string `json:"username" change:"用户名" binding:"omitempty,notblank,max=64" binding_msg:"notblank=用户名不能为空,max=用户名最多64位"`
	Nickname     *string `json:"nickname" change:"昵称" binding:"omitempty,max=64" binding_msg:"max=昵称最多64位"`
	LanguageCode *string `json:"languageCode" change:"语言" binding:"omitempty,max=32" binding_msg:"max=语言代码最多32位"`
	IsEnabled    *bool   `json:"isEnabled" change:"启用状态"`
	Remark       *string `json:"remark" change:"备注" binding:"omitempty,max=255" binding_msg:"max=备注最多255位"`
}

type ReqSysUserBatchDelete struct {
	IDs []uint64 `json:"ids" binding:"required,min=1" binding_msg:"required=请选择用户,min=至少选择一个用户"`
}

// @Summary 获取用户分页列表
// @Remark 分页查询用户信息
// @Tags User
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[RespSysUser]} "成功"
// @Router /api/sys/user/list [post]
func (h *SysUserHandler) List(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[RespSysUser], error) {
	if req.GetOrderBy() == "" {
		orderBy := "id desc"
		req.OrderBy = &orderBy
	}
	pagination, err := h.Q.SysUser.PageWithPaging(req)
	if err != nil {
		return nil, res.FailDefault
	}

	items := make([]*RespSysUser, 0, len(pagination.Items))
	for _, item := range pagination.Items {
		if item == nil {
			continue
		}
		items = append(items, &RespSysUser{
			SysUser:   *item,
			CanWrite:  true,
			CanDelete: true,
		})
	}
	return &gormc.PagingResult[RespSysUser]{
		Items: items,
		Total: pagination.Total,
	}, nil
}

// @Summary 创建用户
// @Remark 创建新的后台用户
// @Tags User
// @Accept json
// @Produce json
// @Param req body ReqSysUserCreate true "用户创建参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/user/create [post]
func (h *SysUserHandler) Create(ctx *handler.Ctx, req *ReqSysUserCreate) error {
	req.normalize()

	encodedPwd, err := passwd.Encode(req.Password)
	if err != nil {
		return res.FailDefault
	}

	operationID := ctx.SessionInfo.Id
	sysUser := h.Q.SysUser
	err = sysUser.Create(&models.SysUser{
		OperatorID: mixin.OperatorID{
			CreatedBy: mixin.CreatedBy{CreatedBy: operationID},
			UpdatedBy: mixin.UpdatedBy{UpdatedBy: operationID},
		},
		Remark:       mixin.Remark{Remark: req.Remark},
		IsEnabled:    mixin.IsEnabled{IsEnabled: req.IsEnabled},
		Username:     req.Username,
		Nickname:     req.Nickname,
		Password:     encodedPwd,
		LanguageCode: req.LanguageCode,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return res.FailMsg("用户名已存在")
		}
		return res.FailDefault
	}
	return nil
}

// @Summary 更新用户
// @Remark 根据 ID 更新用户基础信息
// @Tags User
// @Accept json
// @Produce json
// @Param req body ReqSysUserUpdate true "用户更新参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/user/update [post]
func (h *SysUserHandler) Update(ctx *handler.Ctx, req *ReqSysUserUpdate) error {
	req.normalize()

	sysUser := h.Q.SysUser
	_, err := sysUser.Select(sysUser.ID).Where(sysUser.ID.Eq(req.ID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res.FailMsg("用户不存在")
		}
		return res.FailDefault
	}

	exprs := []field.AssignExpr{sysUser.UpdatedBy.Value(ctx.SessionInfo.Id)}
	query.ExprAppendSelf(&exprs, req.Username, sysUser.Username.Value)
	query.ExprAppendSelf(&exprs, req.Nickname, sysUser.Nickname.Value)
	query.ExprAppendSelf(&exprs, req.LanguageCode, sysUser.LanguageCode.Value)
	query.ExprAppendSelf(&exprs, req.IsEnabled, sysUser.IsEnabled.Value)
	query.ExprAppendSelf(&exprs, req.Remark, sysUser.Remark.Value)

	info, err := sysUser.Where(sysUser.ID.Eq(req.ID)).UpdateSimple(exprs...)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return res.FailMsg("用户名已存在")
		}
		return res.FailDefault
	}
	if info.RowsAffected == 0 {
		return res.FailMsg("用户不存在")
	}
	return nil
}

// @Summary 删除用户
// @Remark 根据 ID 列表批量删除用户
// @Tags User
// @Accept json
// @Produce json
// @Param req body ReqSysUserBatchDelete true "批量删除参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/user/del [post]
func (h *SysUserHandler) Del(ctx *handler.Ctx, req *ReqSysUserBatchDelete) error {
	ids := slices_utils.Distinct(req.IDs)
	if len(ids) == 0 {
		return res.FailMsg("请选择用户")
	}
	sysUser := h.Q.SysUser
	info, err := sysUser.Where(sysUser.ID.In(ids...)).Delete()
	if err != nil {
		return res.FailDefault
	}
	if info.RowsAffected != int64(len(ids)) {
		return res.FailMsg("用户不存在")
	}
	return nil
}

func (req *ReqSysUserCreate) normalize() {
	req.Username = strings.TrimSpace(req.Username)
	req.Nickname = strings.TrimSpace(req.Nickname)
	req.Password = strings.TrimSpace(req.Password)
	req.LanguageCode = strings.TrimSpace(req.LanguageCode)
	req.Remark = strings.TrimSpace(req.Remark)
}

func (req *ReqSysUserUpdate) normalize() {
	str.TrimStringPtr(req.Username, nil)
	str.TrimStringPtr(req.Nickname, nil)
	str.TrimStringPtr(req.LanguageCode, nil)
	str.TrimStringPtr(req.Remark, nil)
}
