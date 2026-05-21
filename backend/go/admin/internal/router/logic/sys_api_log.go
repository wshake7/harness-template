package logic

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc"

	"go.uber.org/zap"
)

type SysApiLogHandler struct {
	Q *query.Query
}

func NewSysApiLogHandler(q *query.Query) *SysApiLogHandler {
	return &SysApiLogHandler{Q: q}
}

type ReqLogDetail struct {
	ID uint64 `json:"id" binding:"required" binding_msg:"required=请求错误"`
}

// @Summary 获取API日志分页列表
// @Remark 分页查询 API 调用日志
// @Tags Log
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[models.SysApiLog]} "成功"
// @Router /api/sys/api/log/list [post]
func (h *SysApiLogHandler) List(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[models.SysApiLog], error) {
	pagination, err := h.Q.SysApiLog.PageWithPaging(req)
	if err != nil {
		return nil, res.FailDefault
	}
	if err = h.reloadUsers(pagination.Items); err != nil {
		ctx.L().Error("查询API日志用户信息失败", zap.Error(err))
		return nil, res.FailDefault
	}
	return pagination, nil
}

func (h *SysApiLogHandler) reloadUsers(items []*models.SysApiLog) error {
	if len(items) == 0 {
		return nil
	}
	userIDSet := make(map[uint64]struct{})
	for _, item := range items {
		if item == nil || item.SysUserID == nil {
			continue
		}
		userIDSet[*item.SysUserID] = struct{}{}
	}
	if len(userIDSet) == 0 {
		return nil
	}
	userIDs := make([]uint64, 0, len(userIDSet))
	for id := range userIDSet {
		userIDs = append(userIDs, id)
	}
	sysUser := h.Q.SysUser
	users, err := sysUser.Select(sysUser.ID, sysUser.Username).Where(sysUser.ID.In(userIDs...)).Find()
	if err != nil {
		return err
	}
	userMap := make(map[uint64]*models.SysUser, len(users))
	for _, u := range users {
		if u != nil {
			userMap[u.ID] = u
		}
	}
	for _, item := range items {
		if item == nil || item.SysUserID == nil {
			continue
		}
		item.SysUser = userMap[*item.SysUserID]
	}
	return nil
}

// @Summary 获取API日志详情
// @Remark 根据 ID 查询 API 日志详情
// @Tags Log
// @Accept json
// @Produce json
// @Param req body ReqLogDetail true "日志详情参数"
// @Success 200 {object} res.Response{data=models.SysApiLog} "成功"
// @Router /api/sys/api/log/detail [post]
func (h *SysApiLogHandler) Detail(ctx *handler.Ctx, req *ReqLogDetail) (*models.SysApiLog, error) {
	sysApiLog := h.Q.SysApiLog
	sysUser := h.Q.SysUser
	logEntry, err := sysApiLog.
		Preload(sysApiLog.SysUser.Select(sysUser.ID, sysUser.Username)).
		Where(sysApiLog.ID.Eq(req.ID)).
		First()
	if err != nil {
		ctx.L().Error("查询API日志失败", zap.Error(err), zap.Uint64("id", req.ID))
		return nil, res.FailDefault
	}
	return logEntry, nil
}
