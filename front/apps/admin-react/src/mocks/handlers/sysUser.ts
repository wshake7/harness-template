import type { PagingRequest } from '@vp/core'
import { http, HttpResponse } from 'msw'
import { fail, success, url } from '.'

let userId = 3

let users = [
  { id: 1, username: 'admin', nickname: '管理员', languageCode: 'zh-CN', isEnabled: true, remark: '系统管理员', createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z', createdBy: 1, updatedBy: 1, canWrite: true, canDelete: false },
  { id: 2, username: 'test', nickname: '测试用户', languageCode: 'zh-CN', isEnabled: true, remark: '测试账号', createdAt: '2024-01-02T00:00:00Z', updatedAt: '2024-01-02T00:00:00Z', createdBy: 1, updatedBy: 1, canWrite: true, canDelete: true },
]

function paginate<T>(items: T[], req: PagingRequest): { items: T[], total: number } {
  if (req.noPaging) {
    return { items, total: items.length }
  }
  const page = req.page ?? 1
  const pageSize = req.pageSize ?? 20
  const start = (page - 1) * pageSize
  return { items: items.slice(start, start + pageSize), total: items.length }
}

export const sysUserHandlers = [
  http.post(url('/api/sys/user/list'), async ({ request }) => {
    const req = (await request.json()) as PagingRequest
    let items = [...users]
    if (req.query) {
      try {
        const q = JSON.parse(req.query) as Record<string, string>
        if (q.username) { items = items.filter(i => i.username.includes(q.username)) }
        if (q.nickname) { items = items.filter(i => i.nickname.includes(q.nickname)) }
      }
      catch { /* ignore */ }
    }
    return HttpResponse.json(success(paginate(items, req)))
  }),

  http.post(url('/api/sys/user/create'), async ({ request }) => {
    const body = await request.json() as Record<string, unknown>
    const now = new Date().toISOString()
    users.push({
      id: userId++,
      username: body.username as string,
      nickname: (body.nickname as string) ?? '',
      languageCode: body.languageCode as string,
      isEnabled: body.isEnabled as boolean,
      remark: (body.remark as string) ?? '',
      createdAt: now,
      updatedAt: now,
      createdBy: 1,
      updatedBy: 1,
      canWrite: true,
      canDelete: true,
    })
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/user/update'), async ({ request }) => {
    const body = await request.json() as Record<string, unknown>
    const id = body.id as number
    const idx = users.findIndex(i => i.id === id)
    if (idx === -1) { return HttpResponse.json(fail('用户不存在')) }
    users[idx] = { ...users[idx], ...body, updatedAt: new Date().toISOString() }
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/user/del'), async ({ request }) => {
    const body = await request.json() as { ids: number[] }
    users = users.filter(i => !body.ids.includes(i.id))
    return HttpResponse.json(success(undefined))
  }),
]
