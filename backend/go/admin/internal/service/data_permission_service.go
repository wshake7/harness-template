package service

import (
	"admin/internal/fiberc/handler"
	"admin/internal/services/orm/data_permission"
	"admin/internal/services/orm/query"
	v1 "orm-crud/api/gen/go/pagination/v1"
)

//go:generate mockgen -source=data_permission_service.go -destination=../mock/mock_data_permission_service.go -package=mock -typed

type DataPermissionService interface {
	BuildFilterExprsForCtx(ctx *handler.Ctx, resourceTable string, actions ...string) (map[string]*v1.FilterExpr, error)
	ApplyPagePermissionExpr(req *v1.PagingRequest, expr *v1.FilterExpr) error
	BuildReadPermissionQuery(ctx *handler.Ctx, resourceTable string) (*query.Query, error)
	BuildWritePermissionQuery(ctx *handler.Ctx, resourceTable string) (*query.Query, error)
	BuildDeletePermissionQuery(ctx *handler.Ctx, resourceTable string) (*query.Query, error)
	BuildPermissionQueryFromExpr(expr *v1.FilterExpr) (*query.Query, error)
}

type dataPermissionServiceImpl struct{}

func NewDataPermissionService() DataPermissionService {
	return &dataPermissionServiceImpl{}
}

func (s *dataPermissionServiceImpl) BuildFilterExprsForCtx(ctx *handler.Ctx, resourceTable string, actions ...string) (map[string]*v1.FilterExpr, error) {
	return datapermission.BuildPermissionFilterExprsForCtx(ctx, resourceTable, actions...)
}

func (s *dataPermissionServiceImpl) ApplyPagePermissionExpr(req *v1.PagingRequest, expr *v1.FilterExpr) error {
	return datapermission.ApplyPagePermissionExpr(req, expr)
}

func (s *dataPermissionServiceImpl) BuildReadPermissionQuery(ctx *handler.Ctx, resourceTable string) (*query.Query, error) {
	return datapermission.BuildReadPermissionQuery(ctx, resourceTable)
}

func (s *dataPermissionServiceImpl) BuildWritePermissionQuery(ctx *handler.Ctx, resourceTable string) (*query.Query, error) {
	return datapermission.BuildWritePermissionQuery(ctx, resourceTable)
}

func (s *dataPermissionServiceImpl) BuildDeletePermissionQuery(ctx *handler.Ctx, resourceTable string) (*query.Query, error) {
	return datapermission.BuildDeletePermissionQuery(ctx, resourceTable)
}

func (s *dataPermissionServiceImpl) BuildPermissionQueryFromExpr(expr *v1.FilterExpr) (*query.Query, error) {
	return datapermission.BuildPermissionQueryFromExpr(expr)
}
