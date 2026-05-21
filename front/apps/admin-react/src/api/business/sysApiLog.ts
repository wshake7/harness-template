import type { PagingRequest, PagingResult } from '@vp/core'
import API from '../index'

export interface SysApiLog {
  id: number
  requestID: string
  method: string
  module: string
  path: string
  referer: string
  beforeChange: string
  afterChange: string
  formatChange: string
  requestURI: string
  requestBody: string
  requestHeader: string
  response: string
  costTime: number
  sysUserID?: number | null
  sysUser?: {
    id: number
    username: string
    nickname: string
  } | null
  clientIP: string
  statusCode: number
  reason: string
  success: boolean
  location: string
  userAgent: string
  browserName: string
  browserVersion: string
  clientID: string
  clientName: string
  osName?: string
  osVersion?: string
  oSName?: string
  oSVersion?: string
  createdAt: string
}

export interface ReqLogDetail {
  id: number
}

function list(req: PagingRequest) {
  return API.Post<Res<PagingResult<SysApiLog>>>('/api/sys/api/log/list', req, {
    cacheFor: 0,
  })
}

async function detail(req: ReqLogDetail) {
  return await API.Post<Res<SysApiLog>>('/api/sys/api/log/detail', req, {
    cacheFor: 0,
  }).send()
}

export const ApiLogApi = {
  list,
  detail,
}
