import type { PagingRequest, PagingResult } from '@vp/core'
import API from '../index'

export interface KnowledgeCollection {
  id: number
  collectionName: string
  displayName: string
  description: string
  embeddingModel: string
  vectorDimension: number
  metricType: string
  indexType: string
  idMaxLength: number
  contentMaxLength: number
  documentCount: number
  status: string
  isEnabled: boolean
  remark: string
  createdAt?: string
  updatedAt?: string
  canWrite?: boolean
  canDelete?: boolean
}

export interface ReqKnowledgeCollectionCreate {
  collectionName: string
  displayName: string
  description?: string
  embeddingModel: string
  vectorDimension: number
  metricType?: string
  indexType?: string
  idMaxLength?: number
  contentMaxLength?: number
  isEnabled: boolean
  remark?: string
}

export interface ReqKnowledgeCollectionUpdate extends Partial<ReqKnowledgeCollectionCreate> {
  id: number
}

export interface ReqKnowledgeCollectionID {
  id: number
}

function list(req: PagingRequest) {
  return API.Post<Res<PagingResult<KnowledgeCollection>>>('/api/knowledge/collection/list', req, {
    cacheFor: 0,
  })
}

async function detail(req: ReqKnowledgeCollectionID) {
  return await API.Post<Res<KnowledgeCollection>>('/api/knowledge/collection/detail', req, {
    cacheFor: 0,
  }).send()
}

async function create(req: ReqKnowledgeCollectionCreate) {
  await API.Post<Res>('/api/knowledge/collection/create', req, {
    cacheFor: 0,
  }).send()
}

async function update(req: ReqKnowledgeCollectionUpdate) {
  await API.Post<Res>('/api/knowledge/collection/update', req, {
    cacheFor: 0,
  }).send()
}

async function del(req: ReqKnowledgeCollectionID) {
  await API.Post<Res>('/api/knowledge/collection/del', req, {
    cacheFor: 0,
  }).send()
}

export const KnowledgeCollectionApi = {
  list,
  detail,
  create,
  update,
  del,
}
