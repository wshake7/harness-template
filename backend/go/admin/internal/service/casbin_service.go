package service

import (
	"admin/internal/services/casbin"
	"admin/internal/services/orm/models"
)

//go:generate mockgen -source=casbin_service.go -destination=../mock/mock_casbin_service.go -package=mock -typed

type CasbinService interface {
	SyncRoleState(roleID uint64, oldCode string, oldEnabled bool, newCode string, newEnabled bool) error
	SyncRoleAPIPermissions(roleCode string, oldAPIIDs, newAPIIDs []uint64) error
	SyncAPIResourcePolicies(oldAPI, newAPI *models.SysResourceApi) error
	RemoveAPIResourcePolicies(api *models.SysResourceApi) error
}

type casbinServiceImpl struct{}

func NewCasbinService() CasbinService {
	return &casbinServiceImpl{}
}

func (s *casbinServiceImpl) SyncRoleState(roleID uint64, oldCode string, oldEnabled bool, newCode string, newEnabled bool) error {
	return casbin.SyncRoleState(roleID, oldCode, oldEnabled, newCode, newEnabled)
}

func (s *casbinServiceImpl) SyncRoleAPIPermissions(roleCode string, oldAPIIDs, newAPIIDs []uint64) error {
	return casbin.SyncRoleAPIPermissions(roleCode, oldAPIIDs, newAPIIDs)
}

func (s *casbinServiceImpl) SyncAPIResourcePolicies(oldAPI, newAPI *models.SysResourceApi) error {
	return casbin.SyncAPIResourcePolicies(oldAPI, newAPI)
}

func (s *casbinServiceImpl) RemoveAPIResourcePolicies(api *models.SysResourceApi) error {
	return casbin.RemoveAPIResourcePolicies(api)
}
