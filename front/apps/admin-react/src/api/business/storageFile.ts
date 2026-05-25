import API from '../index'

export interface StorageFileAsset {
  id: number
  engine: string
  bucket: string
  objectKey: string
  originalName: string
  contentType: string
  extension?: string
  size: number
  sha256: string
  bizType?: string
  bizID?: string
  metadata?: Record<string, unknown>
  status?: string
  remark?: string
  createdAt?: string
  updatedAt?: string
}

export interface StorageUploadFields {
  bizType?: string
  bizID?: string
  metadata?: string
  remark?: string
}

export interface StorageFileIDReq {
  id: number
}

export interface StorageFilePresignedReq extends StorageFileIDReq {
  expiresSeconds?: number
  disposition?: 'inline' | 'attachment'
}

export interface StorageFilePresignedRes {
  url: string
  expiresAt: string
}

async function upload(file: File, fields: StorageUploadFields = {}) {
  const formData = new FormData()
  formData.append('file', file)
  Object.entries(fields).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') {
      formData.append(key, String(value))
    }
  })

  return await API.Post<Res<StorageFileAsset>>('/api/storage/file/upload', formData, {
    cacheFor: 0,
    meta: { skipEncrypt: true },
  }).send()
}

async function detail(req: StorageFileIDReq) {
  return await API.Post<Res<StorageFileAsset>>('/api/storage/file/detail', req, {
    cacheFor: 0,
  }).send()
}

async function presigned(req: StorageFilePresignedReq) {
  return await API.Post<Res<StorageFilePresignedRes>>('/api/storage/file/presigned', req, {
    cacheFor: 0,
  }).send()
}

async function del(req: StorageFileIDReq) {
  return await API.Post<Res>('/api/storage/file/del', req, {
    cacheFor: 0,
  }).send()
}

export const StorageFileApi = {
  upload,
  detail,
  presigned,
  del,
}
