import type { PagingRequest, PagingResult } from '@vp/core'
import API from '../index'

export interface SysLoginLog {
  id: number
  username: string
  loginIP: string
  loginMAC: string
  loginTime?: string
  userAgent: string
  browserName: string
  browserVersion: string
  clientID: string
  clientName: string
  osName?: string
  osVersion?: string
  sysUserID?: number | null
  sysUser?: {
    id: number
    username: string
    nickname: string
  } | null
  statusCode: number
  success: boolean
  reason: string
  location: string
  createdAt: string
}

export interface ReqLoginLogDetail {
  id: number
}

function list(req: PagingRequest) {
  return API.Post<Res<PagingResult<SysLoginLog>>>('/api/sys/login/log/list', req, {
    cacheFor: 0,
  })
}

async function detail(req: ReqLoginLogDetail) {
  return await API.Post<Res<SysLoginLog>>('/api/sys/login/log/detail', req, {
    cacheFor: 0,
  }).send()
}

export const LoginLogApi = {
  list,
  detail,
}
