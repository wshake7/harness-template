package logic

import (
	"errors"
	"slices"
	"strings"

	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	"admin/internal/service"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"go-common/utils/slices_utils"
	"go-common/utils/str"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc"
	"orm-crud/gormc/mixin"

	"gorm.io/gorm"
)

type SysRoleHandler struct {
	Q      *query.Query
	Casbin service.CasbinService
}

func NewSysRoleHandler(q *query.Query, casbin service.CasbinService) *SysRoleHandler {
	return &SysRoleHandler{Q: q, Casbin: casbin}
}

type RespSysRole struct {
	models.SysRole
	Children  []*RespSysRole `json:"children,omitempty"`
	CanWrite  bool           `json:"canWrite"`
	CanDelete bool           `json:"canDelete"`
}

type RespRolePermission struct {
	MenuIDs []uint64 `json:"menuIDs"`
	ApiIDs  []uint64 `json:"apiIDs"`
}

type ReqSysRoleCreate struct {
	Name      string  `json:"name" change:"角色名称" binding:"required,max=255" binding_msg:"required=角色名称不能为空,max=角色名称最多255位"`
	Code      string  `json:"code" change:"角色标识" binding:"required,max=128" binding_msg:"required=角色标识不能为空,max=角色标识最多128位"`
	ParentID  *uint64 `json:"parentID" change:"父级角色"`
	IsEnabled bool    `json:"isEnabled" change:"启用状态"`
	Remark    string  `json:"remark" change:"备注" binding:"max=255" binding_msg:"max=备注最多255位"`
}

type ReqSysRoleUpdate struct {
	ID        uint64  `json:"id" binding:"required" binding_msg:"required=请求错误"`
	Name      *string `json:"name" change:"角色名称" binding:"omitempty,max=255" binding_msg:"max=角色名称最多255位"`
	Code      *string `json:"code" change:"角色标识" binding:"omitempty,max=128" binding_msg:"max=角色标识最多128位"`
	ParentID  *uint64 `json:"parentID" change:"父级角色"`
	IsEnabled *bool   `json:"isEnabled" change:"启用状态"`
	Remark    *string `json:"remark" change:"备注" binding:"omitempty,max=255" binding_msg:"max=备注最多255位"`
}

type ReqSysRoleBatchDelete struct {
	IDs []uint64 `json:"ids" binding:"required,min=1" binding_msg:"required=请选择角色,min=至少选择一个角色"`
}

type ReqSysRolePermissionQuery struct {
	ID uint64 `params:"id" binding:"required" binding_msg:"required=请求错误"`
}

type ReqSysRolePermissionSave struct {
	ID      uint64   `json:"id" binding:"required" binding_msg:"required=请求错误"`
	MenuIDs []uint64 `json:"menuIDs" change:"菜单权限"`
	ApiIDs  []uint64 `json:"apiIDs" change:"API权限"`
}

// @Summary 获取角色列表
// @Remark 分页或不分页查询角色基础信息
// @Tags Role
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[RespSysRole]} "成功"
// @Router /api/sys/role/list [post]
func (h *SysRoleHandler) List(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[RespSysRole], error) {
	if req.GetOrderBy() == "" {
		orderBy := "id desc"
		req.OrderBy = &orderBy
	}
	pagination, err := h.Q.SysRole.PageWithPaging(req)
	if err != nil {
		return nil, res.FailDefault
	}

	items := make([]*RespSysRole, 0, len(pagination.Items))
	for _, item := range pagination.Items {
		if item == nil {
			continue
		}
		items = append(items, &RespSysRole{
			SysRole:   *item,
			CanWrite:  true,
			CanDelete: true,
		})
	}
	if req.GetNoPaging() {
		items = buildRoleRespTreeFromResp(items)
	}
	return &gormc.PagingResult[RespSysRole]{
		Items: items,
		Total: pagination.Total,
	}, nil
}

// @Summary 获取角色树
// @Remark 查询全部角色并按 parentID 组装树
// @Tags Role
// @Produce json
// @Success 200 {object} res.Response{data=[]RespSysRole} "成功"
// @Router /api/sys/role/tree [get]
func (h *SysRoleHandler) Tree(ctx *handler.Ctx) (*[]*RespSysRole, error) {
	sysRole := h.Q.SysRole
	items, err := sysRole.Order(sysRole.ID.Asc()).Find()
	if err != nil {
		return nil, res.FailDefault
	}
	tree := buildRoleRespTree(items)
	return &tree, nil
}

// @Summary 创建角色
// @Remark 创建新的后台角色
// @Tags Role
// @Accept json
// @Produce json
// @Param req body ReqSysRoleCreate true "角色创建参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/role/create [post]
func (h *SysRoleHandler) Create(ctx *handler.Ctx, req *ReqSysRoleCreate) error {
	req.normalize()
	if err := validateSysRoleValues(req.Name, req.Code); err != nil {
		return err
	}

	operationID := ctx.SessionInfo.Id
	return h.Q.Transaction(func(tx *query.Query) error {
		parentID, err := normalizeSysRoleParentID(tx, 0, req.ParentID)
		if err != nil {
			return err
		}
		err = tx.SysRole.Create(&models.SysRole{
			OperatorID: mixin.OperatorID{
				CreatedBy: mixin.CreatedBy{CreatedBy: operationID},
				UpdatedBy: mixin.UpdatedBy{UpdatedBy: operationID},
			},
			Remark:    mixin.Remark{Remark: req.Remark},
			IsEnabled: mixin.IsEnabled{IsEnabled: req.IsEnabled},
			Name:      req.Name,
			Code:      req.Code,
			ParentID:  parentID,
		})
		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return res.FailMsg("角色标识已存在")
			}
			return res.FailDefault
		}
		return nil
	})
}

// @Summary 更新角色
// @Remark 根据 ID 更新角色基础信息，启停状态也通过该接口更新
// @Tags Role
// @Accept json
// @Produce json
// @Param req body ReqSysRoleUpdate true "角色更新参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/role/update [post]
func (h *SysRoleHandler) Update(ctx *handler.Ctx, req *ReqSysRoleUpdate) error {
	req.normalize()
	var oldCode string
	var oldEnabled bool
	var newCode string
	var newEnabled bool
	err := h.Q.Transaction(func(tx *query.Query) error {
		sysRole := tx.SysRole
		current, err := sysRole.Where(sysRole.ID.Eq(req.ID)).First()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return res.FailMsg("角色不存在")
			}
			return res.FailDefault
		}

		name := current.Name
		if req.Name != nil {
			name = *req.Name
		}
		code := current.Code
		if req.Code != nil {
			code = *req.Code
		}
		if err := validateSysRoleValues(name, code); err != nil {
			return err
		}

		updates := map[string]any{
			"updated_by": ctx.SessionInfo.Id,
		}
		if req.Name != nil {
			updates["name"] = *req.Name
		}
		if req.Code != nil {
			updates["code"] = *req.Code
		}
		if req.IsEnabled != nil {
			updates["is_enabled"] = *req.IsEnabled
		}
		if req.Remark != nil {
			updates["remark"] = *req.Remark
		}
		if req.ParentID != nil {
			parentID, err := normalizeSysRoleParentID(tx, req.ID, req.ParentID)
			if err != nil {
				return err
			}
			updates["parent_id"] = parentID
		}

		info, err := sysRole.Where(sysRole.ID.Eq(req.ID)).Updates(updates)
		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return res.FailMsg("角色标识已存在")
			}
			return res.FailDefault
		}
		if info.RowsAffected == 0 {
			return res.FailMsg("角色不存在")
		}
		oldCode = current.Code
		oldEnabled = current.IsEnabled.IsEnabled
		newCode = code
		newEnabled = current.IsEnabled.IsEnabled
		if req.IsEnabled != nil {
			newEnabled = *req.IsEnabled
		}
		return nil
	})
	if err != nil {
		return err
	}
	if oldCode != newCode || oldEnabled != newEnabled {
		if err := h.Casbin.SyncRoleState(req.ID, oldCode, oldEnabled, newCode, newEnabled); err != nil {
			return res.FailDefault
		}
	}
	return nil
}

// @Summary 删除角色
// @Remark 根据 ID 列表批量软删除角色；存在子角色、用户绑定或授权绑定时拒绝删除
// @Tags Role
// @Accept json
// @Produce json
// @Param req body ReqSysRoleBatchDelete true "批量删除参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/role/del [post]
func (h *SysRoleHandler) Del(ctx *handler.Ctx, req *ReqSysRoleBatchDelete) error {
	ids := slices_utils.Distinct(req.IDs)
	if len(ids) == 0 {
		return res.FailMsg("请选择角色")
	}
	if err := ensureSysRoleCanDelete(h.Q, ids); err != nil {
		return err
	}

	info, err := h.Q.SysRole.Where(h.Q.SysRole.ID.In(ids...)).Delete()
	if err != nil {
		return res.FailDefault
	}
	if info.RowsAffected != int64(len(ids)) {
		return res.FailMsg("角色不存在")
	}
	return nil
}

// @Summary 获取角色授权
// @Remark 查询角色已绑定的菜单资源和 API 资源
// @Tags Role
// @Produce json
// @Param id path int true "角色 ID"
// @Success 200 {object} res.Response{data=RespRolePermission} "成功"
// @Router /api/sys/role/{id}/permissions [get]
func (h *SysRoleHandler) Permissions(ctx *handler.Ctx, req *ReqSysRolePermissionQuery) (*RespRolePermission, error) {
	if err := ensureSysRoleExists(query.Q, req.ID); err != nil {
		return nil, err
	}

	roleMenu := h.Q.SysRoleMenu
	menus, err := roleMenu.Select(roleMenu.MenuID).Where(roleMenu.RoleID.Eq(req.ID)).Find()
	if err != nil {
		return nil, res.FailDefault
	}
	roleAPI := h.Q.SysRoleApi
	apis, err := roleAPI.Select(roleAPI.ApiID).Where(roleAPI.RoleID.Eq(req.ID)).Find()
	if err != nil {
		return nil, res.FailDefault
	}

	menuIDs := make([]uint64, 0, len(menus))
	for _, item := range menus {
		menuIDs = append(menuIDs, item.MenuID)
	}
	apiIDs := make([]uint64, 0, len(apis))
	for _, item := range apis {
		apiIDs = append(apiIDs, item.ApiID)
	}
	return &RespRolePermission{
		MenuIDs: menuIDs,
		ApiIDs:  apiIDs,
	}, nil
}

// @Summary 保存角色授权
// @Remark 事务内替换角色绑定的菜单资源和 API 资源
// @Tags Role
// @Accept json
// @Produce json
// @Param req body ReqSysRolePermissionSave true "角色授权参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/role/permissions [post]
func (h *SysRoleHandler) SavePermissions(ctx *handler.Ctx, req *ReqSysRolePermissionSave) error {
	menuIDs := slices_utils.Distinct(req.MenuIDs)
	apiIDs := slices_utils.Distinct(req.ApiIDs)
	var roleCode string
	var roleEnabled bool
	var oldAPIIDs []uint64
	err := h.Q.Transaction(func(tx *query.Query) error {
		sysRole := tx.SysRole
		role, err := sysRole.Select(sysRole.Code, sysRole.IsEnabled).Where(sysRole.ID.Eq(req.ID)).First()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return res.FailMsg("角色不存在")
			}
			return res.FailDefault
		}
		if err := ensureSysRoleMenuIDsExist(tx, menuIDs); err != nil {
			return err
		}
		if err := ensureSysRoleAPIIDsExist(tx, apiIDs); err != nil {
			return err
		}

		roleMenu := tx.SysRoleMenu
		if _, err := roleMenu.Where(roleMenu.RoleID.Eq(req.ID)).Delete(); err != nil {
			return res.FailDefault
		}
		roleAPI := tx.SysRoleApi
		oldAPIs, err := roleAPI.Select(roleAPI.ApiID).Where(roleAPI.RoleID.Eq(req.ID)).Find()
		if err != nil {
			return res.FailDefault
		}
		oldAPIIDs = make([]uint64, 0, len(oldAPIs))
		for _, item := range oldAPIs {
			oldAPIIDs = append(oldAPIIDs, item.ApiID)
		}
		if _, err := roleAPI.Where(roleAPI.RoleID.Eq(req.ID)).Delete(); err != nil {
			return res.FailDefault
		}

		operationID := ctx.SessionInfo.Id
		roleMenus := make([]*models.SysRoleMenu, 0, len(menuIDs))
		for _, id := range menuIDs {
			roleMenus = append(roleMenus, &models.SysRoleMenu{
				OperatorID: mixin.OperatorID{
					CreatedBy: mixin.CreatedBy{CreatedBy: operationID},
					UpdatedBy: mixin.UpdatedBy{UpdatedBy: operationID},
				},
				RoleID: req.ID,
				MenuID: id,
			})
		}
		if len(roleMenus) > 0 {
			if err := roleMenu.CreateInBatches(roleMenus, 100); err != nil {
				return res.FailDefault
			}
		}

		roleAPIs := make([]*models.SysRoleApi, 0, len(apiIDs))
		for _, id := range apiIDs {
			roleAPIs = append(roleAPIs, &models.SysRoleApi{
				OperatorID: mixin.OperatorID{
					CreatedBy: mixin.CreatedBy{CreatedBy: operationID},
					UpdatedBy: mixin.UpdatedBy{UpdatedBy: operationID},
				},
				RoleID: req.ID,
				ApiID:  id,
			})
		}
		if len(roleAPIs) > 0 {
			if err := roleAPI.CreateInBatches(roleAPIs, 100); err != nil {
				return res.FailDefault
			}
		}
		roleCode = role.Code
		roleEnabled = role.IsEnabled.IsEnabled
		return nil
	})
	if err != nil {
		return err
	}
	if roleEnabled {
		if err := h.Casbin.SyncRoleAPIPermissions(roleCode, oldAPIIDs, apiIDs); err != nil {
			return res.FailDefault
		}
	}
	return nil
}

func (req *ReqSysRoleCreate) normalize() {
	req.Name = strings.TrimSpace(req.Name)
	req.Code = strings.TrimSpace(req.Code)
	req.Remark = strings.TrimSpace(req.Remark)
}

func (req *ReqSysRoleUpdate) normalize() {
	str.TrimStringPtr(req.Name, nil)
	str.TrimStringPtr(req.Code, nil)
	str.TrimStringPtr(req.Remark, nil)
}

func validateSysRoleValues(name, code string) error {
	if strings.TrimSpace(name) == "" {
		return res.FailMsg("角色名称不能为空")
	}
	if strings.TrimSpace(code) == "" {
		return res.FailMsg("角色标识不能为空")
	}
	return nil
}

func normalizeSysRoleParentID(tx *query.Query, currentID uint64, parentID *uint64) (*uint64, error) {
	if parentID == nil || *parentID == 0 {
		return nil, nil
	}
	if currentID > 0 {
		if *parentID == currentID {
			return nil, res.FailMsg("父级角色不能选择自身")
		}
		if err := ensureSysRoleParentNotDescendant(tx, currentID, *parentID); err != nil {
			return nil, err
		}
	}
	if err := ensureSysRoleExists(tx, *parentID); err != nil {
		var resp res.Response
		if errors.As(err, &resp) {
			return nil, res.FailMsg("父级角色不存在")
		}
		return nil, err
	}
	return parentID, nil
}

func ensureSysRoleParentNotDescendant(tx *query.Query, currentID, parentID uint64) error {
	sysRole := tx.SysRole
	nextID := parentID
	visited := map[uint64]struct{}{}
	for nextID > 0 {
		if nextID == currentID {
			return res.FailMsg("父级角色不能选择自身或子角色")
		}
		if _, ok := visited[nextID]; ok {
			return res.FailMsg("角色层级存在循环")
		}
		visited[nextID] = struct{}{}
		parent, err := sysRole.Select(sysRole.ID, sysRole.ParentID).Where(sysRole.ID.Eq(nextID)).First()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return res.FailMsg("父级角色不存在")
			}
			return res.FailDefault
		}
		if parent.ParentID == nil {
			return nil
		}
		nextID = *parent.ParentID
	}
	return nil
}

func ensureSysRoleExists(tx *query.Query, id uint64) error {
	sysRole := tx.SysRole
	_, err := sysRole.Select(sysRole.ID).Where(sysRole.ID.Eq(id)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res.FailMsg("角色不存在")
		}
		return res.FailDefault
	}
	return nil
}

func ensureSysRoleCanDelete(q *query.Query, ids []uint64) error {
	sysRole := q.SysRole
	children, err := sysRole.Select(sysRole.ID).Where(sysRole.ParentID.In(ids...)).Find()
	if err != nil {
		return res.FailDefault
	}
	if len(children) > 0 {
		return res.FailMsg("存在子角色，不能删除")
	}

	userRole := q.SysUserRole
	userRoles, err := userRole.Select(userRole.ID).Where(userRole.RoleID.In(ids...)).Limit(1).Find()
	if err != nil {
		return res.FailDefault
	}
	if len(userRoles) > 0 {
		return res.FailMsg("角色已绑定用户，不能删除")
	}

	roleMenu := q.SysRoleMenu
	roleMenus, err := roleMenu.Select(roleMenu.ID).Where(roleMenu.RoleID.In(ids...)).Limit(1).Find()
	if err != nil {
		return res.FailDefault
	}
	if len(roleMenus) > 0 {
		return res.FailMsg("角色已绑定菜单权限，不能删除")
	}

	roleAPI := q.SysRoleApi
	roleAPIs, err := roleAPI.Select(roleAPI.ID).Where(roleAPI.RoleID.In(ids...)).Limit(1).Find()
	if err != nil {
		return res.FailDefault
	}
	if len(roleAPIs) > 0 {
		return res.FailMsg("角色已绑定API权限，不能删除")
	}
	return nil
}

func ensureSysRoleMenuIDsExist(tx *query.Query, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}
	sysResourceMenu := tx.SysResourceMenu
	items, err := sysResourceMenu.Select(sysResourceMenu.ID).Where(sysResourceMenu.ID.In(ids...)).Find()
	if err != nil {
		return res.FailDefault
	}
	if len(items) != len(ids) {
		return res.FailMsg("菜单资源不存在")
	}
	return nil
}

func ensureSysRoleAPIIDsExist(tx *query.Query, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}
	sysResourceAPI := tx.SysResourceApi
	items, err := sysResourceAPI.Select(sysResourceAPI.ID).Where(sysResourceAPI.ID.In(ids...)).Find()
	if err != nil {
		return res.FailDefault
	}
	if len(items) != len(ids) {
		return res.FailMsg("API资源不存在")
	}
	return nil
}

func buildRoleRespTree(items []*models.SysRole) []*RespSysRole {
	nodes := make([]*RespSysRole, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		nodes = append(nodes, &RespSysRole{
			SysRole:   *item,
			CanWrite:  true,
			CanDelete: true,
		})
	}
	return buildRoleRespTreeFromResp(nodes)
}

func buildRoleRespTreeFromResp(items []*RespSysRole) []*RespSysRole {
	nodes := make(map[uint64]*RespSysRole, len(items))
	roots := make([]*RespSysRole, 0)
	for _, item := range items {
		if item == nil {
			continue
		}
		item.Children = nil
		nodes[item.ID] = item
	}
	for _, item := range items {
		if item == nil {
			continue
		}
		if item.ParentID != nil {
			if parent, ok := nodes[*item.ParentID]; ok {
				parent.Children = append(parent.Children, item)
				continue
			}
		}
		roots = append(roots, item)
	}
	sortRoleRespNodes(roots)
	return roots
}

func sortRoleRespNodes(nodes []*RespSysRole) {
	slices.SortFunc(nodes, func(a, b *RespSysRole) int {
		return int(a.ID - b.ID)
	})
	for _, node := range nodes {
		sortRoleRespNodes(node.Children)
	}
}
