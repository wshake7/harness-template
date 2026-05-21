package datapermission

import (
	"admin/internal/services/orm/models"
	"errors"

	"admin/internal/fiberc/handler"
	v1 "orm-crud/api/gen/go/pagination/v1"
	paginationFilter "orm-crud/pagination/filter"
)

// ApplyReadPagePermission merges the current request filter with the read-scope
// filter configured in sys_data_permission for the given resource table.
func ApplyReadPagePermission(req *v1.PagingRequest, resourceTable string, subjects []Subject) error {
	if req == nil {
		return errors.New("paging request is nil")
	}

	permissionExpr, err := buildPermissionFilterExpr(resourceTable, models.DataPermissionActionRead, subjects)
	if err != nil {
		return err
	}
	return ApplyPagePermissionExpr(req, permissionExpr)
}

// ApplyPagePermissionExpr merges the current request filter with a permission
// expression. A nil permission expression means full access.
func ApplyPagePermissionExpr(req *v1.PagingRequest, permissionExpr *v1.FilterExpr) error {
	if req == nil {
		return errors.New("paging request is nil")
	}
	if permissionExpr == nil {
		return nil
	}
	currentExpr, err := paginationFilter.ConvertFilterByPagingRequest(req)
	if err != nil {
		return err
	}
	if currentExpr == nil {
		req.FilteringType = &v1.PagingRequest_FilterExpr{FilterExpr: permissionExpr}
		return nil
	}

	req.FilteringType = &v1.PagingRequest_FilterExpr{FilterExpr: &v1.FilterExpr{
		Type:   v1.ExprType_AND,
		Groups: []*v1.FilterExpr{currentExpr, permissionExpr},
	}}
	return nil
}

// ApplyReadPagePermissionForCtx builds permission subjects from the current
// session and applies read data permission to a paging request.
func ApplyReadPagePermissionForCtx(ctx *handler.Ctx, req *v1.PagingRequest, resourceTable string) error {
	return ApplyReadPagePermission(req, resourceTable, BuildPermissionSubjectsFromCtx(ctx))
}
