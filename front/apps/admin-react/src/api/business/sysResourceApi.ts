import type { PagingRequest, PagingResult } from '@vp/core'
import API from '../index'

export type ResourceApiMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE' | 'OPTIONS' | 'HEAD'

export interface ResourceApi {
  id: number
  module: string
  path: string
  method: ResourceApiMethod
  sortOrder: number
  isEnabled: boolean
  remark: string
  canWrite?: boolean
  canDelete?: boolean
}

export interface ReqResourceApiCreate {
  module: string
  path: string
  method: ResourceApiMethod
  sortOrder: number
  isEnabled: boolean
  remark: string
}

export interface ReqResourceApiUpdate extends Partial<ReqResourceApiCreate> {
  id: number
}

export interface ReqResourceApiBatchDelete {
  ids: number[]
}

function apiList(req: PagingRequest) {
  return API.Post<Res<PagingResult<ResourceApi>>>('/api/sys/resource/api/list', req, {
    cacheFor: 0,
  })
}

async function apiCreate(req: ReqResourceApiCreate) {
  await API.Post<Res>('/api/sys/resource/api/create', req, {
    cacheFor: 0,
  }).send()
}

async function apiUpdate(req: ReqResourceApiUpdate) {
  await API.Post<Res>('/api/sys/resource/api/update', req, {
    cacheFor: 0,
  }).send()
}

async function apiDel(req: ReqResourceApiBatchDelete) {
  await API.Post<Res>('/api/sys/resource/api/del', req, {
    cacheFor: 0,
  }).send()
}

export const ResourceApiApi = {
  apiList,
  apiCreate,
  apiUpdate,
  apiDel,
}
