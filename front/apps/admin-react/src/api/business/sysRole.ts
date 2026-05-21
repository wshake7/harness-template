import type { PagingRequest, PagingResult } from '@vp/core'
import API from '../index'

export interface SysRole {
  id: number
  name: string
  code: string
  parentID?: number | null
  isEnabled: boolean
  remark: string
  createdAt?: string
  updatedAt?: string
  createdBy?: number
  updatedBy?: number
  children?: SysRole[]
  canWrite?: boolean
  canDelete?: boolean
}

export interface ReqSysRoleCreate {
  name: string
  code: string
  parentID?: number
  isEnabled: boolean
  remark: string
}

export interface ReqSysRoleUpdate extends Partial<ReqSysRoleCreate> {
  id: number
}

export interface ReqSysRoleBatchDelete {
  ids: number[]
}

export interface RolePermission {
  menuIDs: number[]
  apiIDs: number[]
}

export interface ReqRolePermissionSave extends RolePermission {
  id: number
}

function apiList(req: PagingRequest) {
  return API.Post<Res<PagingResult<SysRole>>>('/api/sys/role/list', req, {
    cacheFor: 0,
  })
}

function apiTree() {
  return API.Get<Res<SysRole[]>>('/api/sys/role/tree', {
    cacheFor: 0,
  })
}

function apiPermissions(id: number) {
  return API.Get<Res<RolePermission>>(`/api/sys/role/${id}/permissions`, {
    cacheFor: 0,
  })
}

async function apiCreate(req: ReqSysRoleCreate) {
  await API.Post<Res>('/api/sys/role/create', req, {
    cacheFor: 0,
  }).send()
}

async function apiUpdate(req: ReqSysRoleUpdate) {
  await API.Post<Res>('/api/sys/role/update', req, {
    cacheFor: 0,
  }).send()
}

async function apiDel(req: ReqSysRoleBatchDelete) {
  await API.Post<Res>('/api/sys/role/del', req, {
    cacheFor: 0,
  }).send()
}

async function apiSavePermissions(req: ReqRolePermissionSave) {
  await API.Post<Res>('/api/sys/role/permissions', req, {
    cacheFor: 0,
  }).send()
}

export const RoleApi = {
  apiList,
  apiTree,
  apiPermissions,
  apiCreate,
  apiUpdate,
  apiDel,
  apiSavePermissions,
}
