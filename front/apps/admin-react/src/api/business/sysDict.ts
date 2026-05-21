import type { DictMatchedEntriesByCode, PagingRequest, PagingResult, ReqDictEntryMatch } from '@vp/core'
import API from '../index'

export type { DictMatchedEntriesByCode, DictMatchedEntry, ReqDictEntryMatch } from '@vp/core'

export interface DictType {
  id: number
  typeCode: string
  typeName: string
  isEnabled: boolean
  sortOrder: number
  remark: string
  canWrite?: boolean
  canDelete?: boolean
}

export interface DictEntry {
  id: number
  labelComponent: string
  entryLabel: string
  entryValue: string
  languageCode: string
  sysDictTypeId: number
  sysDictType?: {
    typeCode: string
    typeName: string
  }
  sortOrder: number
  isEnabled: boolean
  remark: string
}

export interface ReqDictTypeCreate {
  typeCode: string
  typeName: string
  isEnabled: boolean
  sortOrder: number
  remark: string
}

export interface ReqDictTypeUpdate extends Partial<ReqDictTypeCreate> {
  id: number
}

export interface ReqDictTypeBatchDelete {
  ids: number[]
}

export interface ReqDictEntryCreate {
  labelComponent: string
  entryLabel: string
  entryValue: string
  languageCode: string
  sysDictTypeId: number
  sortOrder: number
  isEnabled: boolean
  remark: string
}

export interface ReqDictEntryUpdate extends Partial<ReqDictEntryCreate> {
  id: number
}

export interface ReqDictEntryBatchDelete {
  ids: number[]
}

export interface ReqDictEntryBatchCopy {
  entryIds: number[]
  targetTypeId: number
}

function typeList(req: PagingRequest) {
  return API.Post<Res<PagingResult<DictType>>>('/api/sys/dict/type/list', req, {
    cacheFor: 0,
  })
}

async function typeCreate(req: ReqDictTypeCreate) {
  await API.Post<Res>('/api/sys/dict/type/create', req, {
    cacheFor: 0,
  }).send()
}

async function typeUpdate(req: ReqDictTypeUpdate) {
  await API.Post<Res>('/api/sys/dict/type/update', req, {
    cacheFor: 0,
  }).send()
}

async function typeDel(req: ReqDictTypeBatchDelete) {
  await API.Post<Res>('/api/sys/dict/type/del', req, {
    cacheFor: 0,
  }).send()
}

function entryList(req: PagingRequest) {
  return API.Post<Res<PagingResult<DictEntry>>>('/api/sys/dict/entry/list', req, {
    cacheFor: 0,
  })
}

function entryMatch(req: ReqDictEntryMatch) {
  return API.Post<Res<DictMatchedEntriesByCode>>('/api/sys/dict/entry/match', req, {
    cacheFor: 0,
  })
}

async function entryCreate(req: ReqDictEntryCreate) {
  await API.Post<Res>('/api/sys/dict/entry/create', req, {
    cacheFor: 0,
  }).send()
}

async function entryUpdate(req: ReqDictEntryUpdate) {
  await API.Post<Res>('/api/sys/dict/entry/update', req, {
    cacheFor: 0,
  }).send()
}

async function entryDel(req: ReqDictEntryBatchDelete) {
  await API.Post<Res>('/api/sys/dict/entry/del', req, {
    cacheFor: 0,
  }).send()
}

async function entryBatchCopy(req: ReqDictEntryBatchCopy) {
  await API.Post<Res>('/api/sys/dict/entry/batch/copy', req, {
    cacheFor: 0,
  }).send()
}

export const DictApi = {
  typeList,
  typeCreate,
  typeUpdate,
  typeDel,
  entryList,
  entryMatch,
  entryCreate,
  entryUpdate,
  entryDel,
  entryBatchCopy,
}
