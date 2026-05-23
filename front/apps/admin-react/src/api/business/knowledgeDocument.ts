import type { PagingRequest, PagingResult } from '@vp/core'
import API from '../index'

export interface KnowledgeDocument {
  id: number
  collectionID: number
  documentID: string
  title: string
  content: string
  contentType: string
  source: string
  chunkIndex: number
  totalChunks: number
  vectorStatus: string
  vectorID: string
  metadata: Record<string, unknown>
  indexingError: string
  lastIndexedAt: number
  isEnabled: boolean
  remark: string
  createdAt?: string
  updatedAt?: string
  canWrite?: boolean
  canDelete?: boolean
  collection?: {
    collectionName: string
    displayName: string
  }
}

export interface ReqKnowledgeDocumentCreate {
  collectionID: number
  documentID: string
  title: string
  content: string
  contentType?: string
  source?: string
  chunkIndex?: number
  totalChunks?: number
  metadata?: string
  isEnabled: boolean
  remark?: string
}

export interface ReqKnowledgeDocumentUpdate extends Partial<ReqKnowledgeDocumentCreate> {
  id: number
}

export interface ReqKnowledgeDocumentID {
  id: number
}

export interface ReqKnowledgeDocumentListByCollection {
  collectionID: number
  page?: number
  pageSize?: number
  orderBy?: string
  query?: string
}

function list(req: PagingRequest) {
  return API.Post<Res<PagingResult<KnowledgeDocument>>>('/api/knowledge/document/list', req, {
    cacheFor: 0,
  })
}

function listByCollection(req: ReqKnowledgeDocumentListByCollection) {
  return API.Post<Res<PagingResult<KnowledgeDocument>>>('/api/knowledge/document/listByCollection', req, {
    cacheFor: 0,
  })
}

async function detail(req: ReqKnowledgeDocumentID) {
  return await API.Post<Res<KnowledgeDocument>>('/api/knowledge/document/detail', req, {
    cacheFor: 0,
  }).send()
}

async function create(req: ReqKnowledgeDocumentCreate) {
  await API.Post<Res>('/api/knowledge/document/create', req, {
    cacheFor: 0,
  }).send()
}

async function update(req: ReqKnowledgeDocumentUpdate) {
  await API.Post<Res>('/api/knowledge/document/update', req, {
    cacheFor: 0,
  }).send()
}

async function del(req: ReqKnowledgeDocumentID) {
  await API.Post<Res>('/api/knowledge/document/del', req, {
    cacheFor: 0,
  }).send()
}

export const KnowledgeDocumentApi = {
  list,
  listByCollection,
  detail,
  create,
  update,
  del,
}
