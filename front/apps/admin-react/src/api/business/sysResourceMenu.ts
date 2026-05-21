import type { PagingRequest, PagingResult } from '@vp/core'
import API from '../index'

export type ResourceMenuType = 'CATALOG' | 'MENU' | 'BUTTON' | 'EMBEDDED' | 'LINK'

export interface ResourceMenuMetadata {
  icon?: string
  order?: number
  hidden?: boolean
  authorities?: string[]
  [key: string]: unknown
}

export interface ResourceMenu {
  id: number
  parentID?: number
  treePath?: string
  menuType: ResourceMenuType
  path: string
  redirect: string
  alias: string
  name: string
  component: string
  metadata: ResourceMenuMetadata
  apiIDs?: number[]
  sortOrder: number
  isEnabled: boolean
  remark: string
  children?: ResourceMenu[]
  canWrite?: boolean
  canDelete?: boolean
}

export interface ResourceMenuNode {
  id: number
  parentID?: number
  menuType: ResourceMenuType
  path: string
  redirect: string
  name: string
  component: string
  icon: string
  order: number
  sortOrder?: number
  hidden: boolean
  authorities: string[]
  isUrl: boolean
  children?: ResourceMenuNode[]
}

export interface ReqResourceMenuCreate {
  parentID?: number
  menuType: ResourceMenuType
  path: string
  redirect: string
  alias: string
  name: string
  component: string
  metadata: ResourceMenuMetadata
  apiIDs: number[]
  sortOrder: number
  isEnabled: boolean
  remark: string
}

export interface ReqResourceMenuUpdate extends Partial<ReqResourceMenuCreate> {
  id: number
}

export interface ReqResourceMenuBatchDelete {
  ids: number[]
}

function menuList(req: PagingRequest) {
  return API.Post<Res<PagingResult<ResourceMenu>>>('/api/sys/resource/menu/list', req, {
    cacheFor: 0,
  })
}

function menuTree() {
  return API.Get<Res<ResourceMenuNode[]>>('/api/sys/resource/menu/tree', {
    cacheFor: 0,
  })
}

async function menuCreate(req: ReqResourceMenuCreate) {
  await API.Post<Res>('/api/sys/resource/menu/create', req, {
    cacheFor: 0,
  }).send()
}

async function menuUpdate(req: ReqResourceMenuUpdate) {
  await API.Post<Res>('/api/sys/resource/menu/update', req, {
    cacheFor: 0,
  }).send()
}

async function menuDel(req: ReqResourceMenuBatchDelete) {
  await API.Post<Res>('/api/sys/resource/menu/del', req, {
    cacheFor: 0,
  }).send()
}

export const ResourceMenuApi = {
  menuList,
  menuTree,
  menuCreate,
  menuUpdate,
  menuDel,
}
