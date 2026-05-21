import type { PagingRequest } from '@vp/core'
import type { SysApiLog } from '~/api/business/sysApiLog'
import { http, HttpResponse } from 'msw'
import { success, url } from '.'

const apiLogs: SysApiLog[] = Array.from({ length: 25 }, (_, i) => ({
  id: i + 1,
  requestID: `req-${1000 + i}`,
  method: ['GET', 'POST', 'PUT', 'DELETE'][i % 4],
  module: ['system', 'account', 'dict'][i % 3],
  path: ['/api/sys/user/list', '/api/sys/role/list', '/api/account/login/pwd'][i % 3],
  referer: '',
  beforeChange: '',
  afterChange: '',
  formatChange: '',
  requestURI: '',
  requestBody: '',
  requestHeader: '',
  response: '{"code":1}',
  costTime: [50, 120, 80, 200][i % 4],
  sysUserID: i % 3 === 0 ? 1 : null,
  sysUser: i % 3 === 0 ? { id: 1, username: 'admin', nickname: '管理员' } : null,
  clientIP: '127.0.0.1',
  statusCode: 200,
  reason: '',
  success: true,
  location: '',
  userAgent: 'Mozilla/5.0',
  browserName: 'Chrome',
  browserVersion: '120.0',
  clientID: 'web',
  clientName: 'Web Client',
  osName: 'macOS',
  osVersion: '14.0',
  oSName: 'macOS',
  oSVersion: '14.0',
  createdAt: new Date(Date.now() - i * 3600000).toISOString(),
}))

function paginate<T>(items: T[], req: PagingRequest): { items: T[], total: number } {
  if (req.noPaging) {
    return { items, total: items.length }
  }
  const page = req.page ?? 1
  const pageSize = req.pageSize ?? 20
  const start = (page - 1) * pageSize
  return { items: items.slice(start, start + pageSize), total: items.length }
}

export const apiLogHandlers = [
  http.post(url('/api/sys/api/log/list'), async ({ request }) => {
    const req = (await request.json()) as PagingRequest
    let items = [...apiLogs]
    if (req.query) {
      try {
        const q = JSON.parse(req.query) as Record<string, string>
        if (q.path) { items = items.filter(i => i.path.includes(q.path)) }
        if (q.method) { items = items.filter(i => i.method === q.method) }
      }
      catch { /* ignore */ }
    }
    return HttpResponse.json(success(paginate(items, req)))
  }),

  http.post(url('/api/sys/api/log/detail'), async ({ request }) => {
    const body = await request.json() as { id: number }
    const log = apiLogs.find(i => i.id === body.id)
    if (!log) { return HttpResponse.json(success(undefined)) }
    return HttpResponse.json(success(log))
  }),
]
