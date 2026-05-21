import type { PagingRequest, PagingResult } from '@vp/core'
import API from '../index'

export interface LanguageType {
  id: number
  typeCode: string
  typeName: string
  isDefault: boolean
  sortOrder: number
  isEnabled: boolean
}

export interface LanguageEntry {
  id: number
  entryCode: string
  entryValue: string
  sysLanguageTypeId: number
  sortOrder: number
  isEnabled: boolean
  remark: string
  sysLanguageType?: {
    typeCode: string
    typeName: string
  }
}

export interface ReqLangTypeCreate {
  typeCode: string
  typeName: string
  isDefault?: boolean
  sortOrder?: number
  isEnabled?: boolean
}

export interface ReqLangTypeUpdate extends Partial<ReqLangTypeCreate> {
  id: number
}

export interface ReqLangTypeDel {
  ids: number[]
}

export interface ReqLangEntryCreate {
  entryCode: string
  entryValue: string
  sysLanguageTypeId: number
  sortOrder?: number
  isEnabled?: boolean
  remark?: string
}

export interface ReqLangEntryUpdate extends Partial<ReqLangEntryCreate> {
  id?: number
  updates?: Array<{
    id: number
  } & Partial<ReqLangEntryCreate>>
}

export interface ReqLangEntryDel {
  ids: number[]
}

export interface ReqLangEntryBatchCreate {
  entryCode: string
  values: Record<string, string>
  sortOrder?: number
  isEnabled?: boolean
}

function typeList(req: PagingRequest) {
  return API.Post<Res<PagingResult<LanguageType>>>('/api/sys/language/type/list', req, {
    cacheFor: 0,
  })
}

async function typeCreate(req: ReqLangTypeCreate) {
  await API.Post<Res>('/api/sys/language/type/create', req, {
    cacheFor: 0,
  }).send()
}

async function typeUpdate(req: ReqLangTypeUpdate) {
  await API.Post<Res>('/api/sys/language/type/update', req, {
    cacheFor: 0,
  }).send()
}

async function typeDel(req: ReqLangTypeDel) {
  await API.Post<Res>('/api/sys/language/type/del', req, {
    cacheFor: 0,
  }).send()
}

function entryList(req: PagingRequest) {
  return API.Post<Res<PagingResult<LanguageEntry>>>('/api/sys/language/entry/list', req, {
    cacheFor: 0,
  })
}

async function entryCreate(req: ReqLangEntryCreate) {
  await API.Post<Res>('/api/sys/language/entry/create', req, {
    cacheFor: 0,
  }).send()
}

async function entryUpdate(req: ReqLangEntryUpdate) {
  await API.Post<Res>('/api/sys/language/entry/update', req, {
    cacheFor: 0,
  }).send()
}

async function entryDel(req: ReqLangEntryDel) {
  await API.Post<Res>('/api/sys/language/entry/del', req, {
    cacheFor: 0,
  }).send()
}

async function entryBatchCreate(req: ReqLangEntryBatchCreate) {
  await API.Post<Res>('/api/sys/language/entry/batch/create', req, {
    cacheFor: 0,
  }).send()
}

export const LangApi = {
  typeList,
  typeCreate,
  typeUpdate,
  typeDel,
  entryList,
  entryCreate,
  entryUpdate,
  entryDel,
  entryBatchCreate,
}
