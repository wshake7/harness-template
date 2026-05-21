package datapermission

import (
	"maps"
	"slices"

	"admin/internal/fiberc/handler"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	v1 "orm-crud/api/gen/go/pagination/v1"
	paginationFilter "orm-crud/pagination/filter"

	"github.com/bytedance/sonic"
	"gorm.io/gen/field"
)

type permissionAction string

type permissionRowFilter struct {
	query      map[string]any
	fullAccess bool
}

// composeActionFilterQuery collects all allow branches for one action.
// A nil return means at least one permission grants full access.
func composeActionFilterQuery(resourceTable string, action permissionAction, subjects []Subject) (map[string]any, error) {
	permissions, err := loadMatchingPermissionRows(resourceTable, subjects)
	if err != nil {
		return nil, err
	}

	return composeActionFilterQueryFromRows(permissions, action), nil
}

func composeActionFilterQueryFromRows(permissions []*models.SysDataPermission, action permissionAction) map[string]any {
	allowBranches := make([]map[string]any, 0, len(permissions))
	for _, permission := range permissions {
		if !permissionIncludesAction(permission, action) {
			continue
		}

		filter, ok := buildRowFilterFromPermission(permission)
		if !ok {
			continue
		}
		if filter.fullAccess {
			return nil
		}
		allowBranches = append(allowBranches, filter.query)
	}

	return mergeAllowFilters(allowBranches)
}

// BuildPermissionFilterExprs loads matching permission rows once, then builds a
// filter expression for each requested action.
func BuildPermissionFilterExprs(resourceTable string, subjects []Subject, actions ...string) (map[string]*v1.FilterExpr, error) {
	permissions, err := loadMatchingPermissionRows(resourceTable, subjects)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*v1.FilterExpr, len(actions))
	for _, action := range actions {
		queryMap := composeActionFilterQueryFromRows(permissions, permissionAction(action))
		if queryMap == nil {
			result[action] = nil
			continue
		}

		expr, err := convertQueryMapToFilterExpr(queryMap)
		if err != nil {
			return nil, err
		}
		result[action] = expr
	}
	return result, nil
}

// BuildPermissionFilterExprsForCtx builds permission subjects from the current
// session and returns filter expressions for all requested actions.
func BuildPermissionFilterExprsForCtx(ctx *handler.Ctx, resourceTable string, actions ...string) (map[string]*v1.FilterExpr, error) {
	return BuildPermissionFilterExprs(resourceTable, BuildPermissionSubjectsFromCtx(ctx), actions...)
}

// loadMatchingPermissionRows loads enabled permission rows for a table and any
// of the current permission subjects.
func loadMatchingPermissionRows(resourceTable string, subjects []Subject) ([]*models.SysDataPermission, error) {
	subjectPredicates := buildSubjectWhereConditions(subjects)
	if len(subjectPredicates) == 0 {
		return nil, nil
	}

	sysDataPermission := query.SysDataPermission
	return sysDataPermission.
		Where(
			sysDataPermission.IsEnabled.Is(true),
			sysDataPermission.ResourceTable.Eq(resourceTable),
			field.Or(subjectPredicates...),
		).
		Find()
}

// buildSubjectWhereConditions turns permission subjects into OR-able GORM Gen
// predicates: (subject_type, subject_id) IN current subjects.
func buildSubjectWhereConditions(subjects []Subject) []field.Expr {
	predicates := make([]field.Expr, 0, len(subjects))
	for _, subject := range subjects {
		predicates = append(predicates, field.And(
			query.SysDataPermission.SubjectType.Eq(subject.Type),
			query.SysDataPermission.SubjectID.Eq(subject.ID),
		))
	}
	return predicates
}

// permissionIncludesAction checks whether the JSON action list contains the
// requested action. The "all" action grants read/write/delete together.
func permissionIncludesAction(permission *models.SysDataPermission, action permissionAction) bool {
	if len(permission.Action) == 0 {
		return false
	}

	var actions []permissionAction
	if err := sonic.Unmarshal(permission.Action, &actions); err != nil {
		return false
	}
	return slices.Contains(actions, models.DataPermissionActionAll) || slices.Contains(actions, action)
}

// buildRowFilterFromPermission converts one permission row into a filter branch.
// Scope fields are appended to any custom conditions already stored on the row.
func buildRowFilterFromPermission(permission *models.SysDataPermission) (permissionRowFilter, bool) {
	condition := map[string]any{}
	maps.Copy(condition, permission.Conditions)

	switch permission.ScopeType {
	case "all":
		return permissionRowFilter{query: condition, fullAccess: len(condition) == 0}, true
	case "custom":
		return permissionRowFilter{query: condition, fullAccess: len(condition) == 0}, true
	case "include":
		values, ok := parseJSONScopeValues(permission)
		if !ok {
			return permissionRowFilter{}, false
		}
		condition[permission.ScopeField+"__in"] = values
		return permissionRowFilter{query: condition}, true
	case "exclude":
		values, ok := parseJSONScopeValues(permission)
		if !ok {
			return permissionRowFilter{}, false
		}
		condition[permission.ScopeField+"__not_in"] = values
		return permissionRowFilter{query: condition}, true
	default:
		return permissionRowFilter{}, false
	}
}

// parseJSONScopeValues decodes scope_values and rejects empty or invalid values,
// because include/exclude scopes need at least one concrete value.
func parseJSONScopeValues(permission *models.SysDataPermission) ([]any, bool) {
	if len(permission.ScopeValues) == 0 {
		return nil, false
	}

	var values []any
	if err := sonic.Unmarshal(permission.ScopeValues, &values); err != nil {
		return nil, false
	}
	if len(values) == 0 {
		return nil, false
	}
	return values, true
}

// mergeAllowFilters combines allowed branches with OR. No matching branch is a
// deliberate deny-all filter instead of falling through to unrestricted access.
func mergeAllowFilters(allowBranches []map[string]any) map[string]any {
	if len(allowBranches) == 0 {
		return map[string]any{"id": 0}
	}
	if len(allowBranches) == 1 {
		return allowBranches[0]
	}
	return map[string]any{"$or": allowBranches}
}

// convertQueryMapToFilterExpr reuses the pagination filter parser so permission
// conditions follow the same field/operator syntax as request filters.
func convertQueryMapToFilterExpr(queryMap map[string]any) (*v1.FilterExpr, error) {
	queryBytes, err := sonic.Marshal(queryMap)
	if err != nil {
		return nil, err
	}
	return paginationFilter.ConvertFilterByPagingRequest(&v1.PagingRequest{
		FilteringType: &v1.PagingRequest_Query{Query: string(queryBytes)},
	})
}
