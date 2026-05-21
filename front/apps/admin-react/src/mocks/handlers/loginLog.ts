import type { PagingRequest } from '@vp/core'
import type { SysLoginLog } from '~/api/business/sysLoginLog'
import { http, HttpResponse } from 'msw'
import { success, url } from '.'

const loginLogs: SysLoginLog[] = Array.from({ length: 20 }, (_, i) => ({
  id: i + 1,
  username: ['admin', 'test', 'user'][i % 3],
  loginIP: '127.0.0.1',
  loginMAC: '',
  loginTime: new Date(Date.now() - i * 7200000).toISOString(),
  userAgent: 'Mozilla/5.0',
  browserName: 'Chrome',
  browserVersion: '120.0',
  clientID: 'web',
  clientName: 'Web Client',
  osName: 'macOS',
  osVersion: '14.0',
  sysUserID: i % 3 === 0 ? 1 : null,
  sysUser: i % 3 === 0 ? { id: 1, username: 'admin', nickname: '管理员' } : null,
  statusCode: i % 5 === 0 ? 401 : 200,
  success: i % 5 !== 0,
  reason: i % 5 === 0 ? '密码错误' : '',
  location: '',
  createdAt: new Date(Date.now() - i * 7200000).toISOString(),
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

export const loginLogHandlers = [
  http.post(url('/api/sys/login/log/list'), async ({ request }) => {
    const req = (await request.json()) as PagingRequest
    let items = [...loginLogs]
    if (req.query) {
      try {
        const q = JSON.parse(req.query) as Record<string, string>
        if (q.username) { items = items.filter(i => i.username.includes(q.username)) }
      }
      catch { /* ignore */ }
    }
    return HttpResponse.json(success(paginate(items, req)))
  }),

  http.post(url('/api/sys/login/log/detail'), async ({ request }) => {
    const body = await request.json() as { id: number }
    const log = loginLogs.find(i => i.id === body.id)
    if (!log) { return HttpResponse.json(success(undefined)) }
    return HttpResponse.json(success(log))
  }),
]
