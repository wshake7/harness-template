import type { PagingRequest, PagingResult } from '@vp/core'
import API from '../index'

export interface SysUser {
  id: number
  username: string
  nickname: string
  languageCode: string
  isEnabled: boolean
  remark: string
  createdAt?: string
  updatedAt?: string
  createdBy?: number
  updatedBy?: number
  canWrite?: boolean
  canDelete?: boolean
}

export interface ReqSysUserCreate {
  username: string
  nickname?: string
  password: string
  languageCode: string
  isEnabled: boolean
  remark: string
}

export interface ReqSysUserUpdate extends Partial<Omit<ReqSysUserCreate, 'password'>> {
  id: number
}

export interface ReqSysUserBatchDelete {
  ids: number[]
}

function apiList(req: PagingRequest) {
  return API.Post<Res<PagingResult<SysUser>>>('/api/sys/user/list', req, {
    cacheFor: 0,
  })
}

async function apiCreate(req: ReqSysUserCreate) {
  await API.Post<Res>('/api/sys/user/create', req, {
    cacheFor: 0,
  }).send()
}

async function apiUpdate(req: ReqSysUserUpdate) {
  await API.Post<Res>('/api/sys/user/update', req, {
    cacheFor: 0,
  }).send()
}

async function apiDel(req: ReqSysUserBatchDelete) {
  await API.Post<Res>('/api/sys/user/del', req, {
    cacheFor: 0,
  }).send()
}

export const SysUserApi = {
  apiList,
  apiCreate,
  apiUpdate,
  apiDel,
}
