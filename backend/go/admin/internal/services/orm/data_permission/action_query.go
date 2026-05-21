package datapermission

import (
	"admin/internal/fiberc/handler"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	v1 "orm-crud/api/gen/go/pagination/v1"
	gormcFilter "orm-crud/gormc/filter"

	"gorm.io/gorm"
)

// BuildWritePermissionQuery returns a query facade whose DB is constrained by
// write-scope data permissions for the current session.
func BuildWritePermissionQuery(ctx *handler.Ctx, resourceTable string) (*query.Query, error) {
	scopes, err := BuildWritePermissionScopes(resourceTable, BuildPermissionSubjectsFromCtx(ctx))
	if err != nil {
		return nil, err
	}
	return query.WithDBScopes(scopes...), nil
}

// BuildReadPermissionQuery returns a query facade whose DB is constrained by
// read-scope data permissions for the current session.
func BuildReadPermissionQuery(ctx *handler.Ctx, resourceTable string) (*query.Query, error) {
	scopes, err := BuildReadPermissionScopes(resourceTable, BuildPermissionSubjectsFromCtx(ctx))
	if err != nil {
		return nil, err
	}
	return query.WithDBScopes(scopes...), nil
}

// BuildPermissionQueryFromExpr returns a query facade constrained by the given
// permission expression. A nil expression means full access.
func BuildPermissionQueryFromExpr(permissionExpr *v1.FilterExpr) (*query.Query, error) {
	scopes, err := BuildPermissionScopesFromExpr(permissionExpr)
	if err != nil {
		return nil, err
	}
	return query.WithDBScopes(scopes...), nil
}

// BuildDeletePermissionQuery returns a query facade whose DB is constrained by
// delete-scope data permissions for the current session.
func BuildDeletePermissionQuery(ctx *handler.Ctx, resourceTable string) (*query.Query, error) {
	scopes, err := BuildDeletePermissionScopes(resourceTable, BuildPermissionSubjectsFromCtx(ctx))
	if err != nil {
		return nil, err
	}
	return query.WithDBScopes(scopes...), nil
}

// BuildReadPermissionScopes builds GORM scopes from read permissions.
func BuildReadPermissionScopes(resourceTable string, subjects []Subject) ([]func(*gorm.DB) *gorm.DB, error) {
	return buildPermissionScopesForAction(resourceTable, models.DataPermissionActionRead, subjects)
}

// BuildPermissionScopesFromExpr builds GORM scopes from a precomputed
// permission expression. A nil expression means full access.
func BuildPermissionScopesFromExpr(permissionExpr *v1.FilterExpr) ([]func(*gorm.DB) *gorm.DB, error) {
	if permissionExpr == nil {
		return nil, nil
	}
	return gormcFilter.NewStructuredFilter().BuildSelectors(permissionExpr)
}

// BuildWritePermissionScopes builds GORM scopes from write permissions.
func BuildWritePermissionScopes(resourceTable string, subjects []Subject) ([]func(*gorm.DB) *gorm.DB, error) {
	return buildPermissionScopesForAction(resourceTable, models.DataPermissionActionWrite, subjects)
}

// BuildDeletePermissionScopes builds GORM scopes from delete permissions.
func BuildDeletePermissionScopes(resourceTable string, subjects []Subject) ([]func(*gorm.DB) *gorm.DB, error) {
	return buildPermissionScopesForAction(resourceTable, models.DataPermissionActionDelete, subjects)
}

func buildPermissionScopesForAction(resourceTable string, action permissionAction, subjects []Subject) ([]func(*gorm.DB) *gorm.DB, error) {
	permissionExpr, err := buildPermissionFilterExpr(resourceTable, action, subjects)
	if err != nil {
		return nil, err
	}
	if permissionExpr == nil {
		return nil, nil
	}
	return BuildPermissionScopesFromExpr(permissionExpr)
}

// buildPermissionFilterExpr converts matched permission rows into a structured
// filter expression shared by pagination filters and GORM scopes.
func buildPermissionFilterExpr(resourceTable string, action permissionAction, subjects []Subject) (*v1.FilterExpr, error) {
	queryMap, err := composeActionFilterQuery(resourceTable, action, subjects)
	if err != nil {
		return nil, err
	}
	if queryMap == nil {
		return nil, nil
	}
	return convertQueryMapToFilterExpr(queryMap)
}
