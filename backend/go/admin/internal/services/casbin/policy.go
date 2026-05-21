package casbin

import (
	"slices"

	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
)

func AddRoleAPIPolicies(roleCode string, apis []*models.SysResourceApi) error {
	if E == nil {
		return nil
	}
	rules := buildRoleAPIPolicyRules(roleCode, apis)
	if len(rules) == 0 {
		return nil
	}
	_, err := E.AddPoliciesEx(rules)
	return err
}

func RemoveRoleAPIPolicies(roleCode string, apis []*models.SysResourceApi) error {
	if E == nil {
		return nil
	}
	for _, api := range apis {
		if api == nil {
			continue
		}
		if _, err := E.RemoveFilteredPolicy(0, roleSubject(roleCode), api.Path, api.Method); err != nil {
			return err
		}
	}
	return nil
}

func SyncRoleAPIPermissions(roleCode string, oldAPIIDs, newAPIIDs []uint64) error {
	if E == nil {
		return nil
	}
	removeIDs, addIDs := diffUint64(oldAPIIDs, newAPIIDs)
	if len(removeIDs) > 0 {
		apis, err := findResourceAPIsByIDs(removeIDs, false)
		if err != nil {
			return err
		}
		if err := RemoveRoleAPIPolicies(roleCode, apis); err != nil {
			return err
		}
	}
	if len(addIDs) > 0 {
		apis, err := findResourceAPIsByIDs(addIDs, true)
		if err != nil {
			return err
		}
		if err := AddRoleAPIPolicies(roleCode, apis); err != nil {
			return err
		}
	}
	return nil
}

func SyncRoleState(roleID uint64, oldCode string, oldEnabled bool, newCode string, newEnabled bool) error {
	if E == nil {
		return nil
	}
	apiIDs, err := findRoleAPIIDs(roleID)
	if err != nil {
		return err
	}
	apis, err := findResourceAPIsByIDs(apiIDs, true)
	if err != nil {
		return err
	}
	if oldEnabled {
		if err := RemoveRoleAPIPolicies(oldCode, apis); err != nil {
			return err
		}
	}
	if newEnabled {
		if err := AddRoleAPIPolicies(newCode, apis); err != nil {
			return err
		}
	}
	return nil
}

func RemoveAPIResourcePolicies(api *models.SysResourceApi) error {
	if E == nil || api == nil {
		return nil
	}
	_, err := E.RemoveFilteredPolicy(1, api.Path, api.Method)
	return err
}

func AddAPIResourcePolicies(api *models.SysResourceApi) error {
	if E == nil || api == nil || !api.IsEnabled.IsEnabled {
		return nil
	}
	roles, err := findEnabledRolesByAPIID(api.ID)
	if err != nil {
		return err
	}
	rules := make([][]string, 0, len(roles))
	for _, role := range roles {
		if role == nil {
			continue
		}
		rules = append(rules, buildRoleAPIPolicyRule(role.Code, api))
	}
	if len(rules) == 0 {
		return nil
	}
	_, err = E.AddPoliciesEx(rules)
	return err
}

func SyncAPIResourcePolicies(oldAPI, newAPI *models.SysResourceApi) error {
	if E == nil {
		return nil
	}
	if oldAPI != nil && oldAPI.IsEnabled.IsEnabled {
		if err := RemoveAPIResourcePolicies(oldAPI); err != nil {
			return err
		}
	}
	if newAPI != nil && newAPI.IsEnabled.IsEnabled {
		if err := AddAPIResourcePolicies(newAPI); err != nil {
			return err
		}
	}
	return nil
}

func buildRoleAPIPolicyRules(roleCode string, apis []*models.SysResourceApi) [][]string {
	rules := make([][]string, 0, len(apis))
	for _, api := range apis {
		if api == nil || !api.IsEnabled.IsEnabled {
			continue
		}
		rules = append(rules, buildRoleAPIPolicyRule(roleCode, api))
	}
	return rules
}

func buildRoleAPIPolicyRule(roleCode string, api *models.SysResourceApi) []string {
	return []string{roleSubject(roleCode), api.Path, api.Method}
}

func roleSubject(roleCode string) string {
	return "role:" + roleCode
}

func diffUint64(oldIDs, newIDs []uint64) ([]uint64, []uint64) {
	oldIDs = slices.Clone(oldIDs)
	newIDs = slices.Clone(newIDs)
	slices.Sort(oldIDs)
	slices.Sort(newIDs)
	oldIDs = slices.Compact(oldIDs)
	newIDs = slices.Compact(newIDs)

	removeIDs := make([]uint64, 0)
	addIDs := make([]uint64, 0)
	for _, id := range oldIDs {
		if _, ok := slices.BinarySearch(newIDs, id); !ok {
			removeIDs = append(removeIDs, id)
		}
	}
	for _, id := range newIDs {
		if _, ok := slices.BinarySearch(oldIDs, id); !ok {
			addIDs = append(addIDs, id)
		}
	}
	return removeIDs, addIDs
}

func findResourceAPIsByIDs(ids []uint64, onlyEnabled bool) ([]*models.SysResourceApi, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	sysResourceAPI := query.SysResourceApi
	dao := sysResourceAPI.
		Select(sysResourceAPI.ID, sysResourceAPI.Path, sysResourceAPI.Method, sysResourceAPI.IsEnabled).
		Where(sysResourceAPI.ID.In(ids...))
	if onlyEnabled {
		dao = dao.Where(sysResourceAPI.IsEnabled.Is(true))
	}
	return dao.Find()
}

func findRoleAPIIDs(roleID uint64) ([]uint64, error) {
	roleAPI := query.SysRoleApi
	items, err := roleAPI.Select(roleAPI.ApiID).Where(roleAPI.RoleID.Eq(roleID)).Find()
	if err != nil {
		return nil, err
	}
	ids := make([]uint64, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ApiID)
	}
	return ids, nil
}

func findEnabledRolesByAPIID(apiID uint64) ([]*models.SysRole, error) {
	roleAPI := query.SysRoleApi
	roleAPIs, err := roleAPI.Select(roleAPI.RoleID).Where(roleAPI.ApiID.Eq(apiID)).Find()
	if err != nil {
		return nil, err
	}
	if len(roleAPIs) == 0 {
		return nil, nil
	}
	roleIDs := make([]uint64, 0, len(roleAPIs))
	for _, item := range roleAPIs {
		roleIDs = append(roleIDs, item.RoleID)
	}
	sysRole := query.SysRole
	return sysRole.
		Select(sysRole.ID, sysRole.Code, sysRole.IsEnabled).
		Where(sysRole.ID.In(roleIDs...), sysRole.IsEnabled.Is(true)).
		Find()
}
