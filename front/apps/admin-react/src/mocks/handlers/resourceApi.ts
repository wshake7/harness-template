import type { PagingRequest } from '@vp/core'
import type { ResourceApi } from '~/api/business/sysResourceApi'
import { http, HttpResponse } from 'msw'
import { fail, success, url } from '.'

let apiId = 4

let apis: ResourceApi[] = [
  { id: 1, module: 'system', path: '/api/sys/user/list', method: 'POST', sortOrder: 1, isEnabled: true, remark: '用户列表', canWrite: true, canDelete: true },
  { id: 2, module: 'system', path: '/api/sys/role/list', method: 'POST', sortOrder: 2, isEnabled: true, remark: '角色列表', canWrite: true, canDelete: true },
  { id: 3, module: 'system', path: '/api/sys/resource/menu/list', method: 'POST', sortOrder: 3, isEnabled: true, remark: '菜单列表', canWrite: true, canDelete: true },
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

export const resourceApiHandlers = [
  http.post(url('/api/sys/resource/api/list'), async ({ request }) => {
    const req = (await request.json()) as PagingRequest
    let items = [...apis]
    if (req.query) {
      try {
        const q = JSON.parse(req.query) as Record<string, string>
        if (q.module) { items = items.filter(i => i.module.includes(q.module)) }
        if (q.path) { items = items.filter(i => i.path.includes(q.path)) }
      }
      catch { /* ignore */ }
    }
    return HttpResponse.json(success(paginate(items, req)))
  }),

  http.post(url('/api/sys/resource/api/create'), async ({ request }) => {
    const body = await request.json() as Record<string, unknown>
    apis.push({
      id: apiId++,
      module: body.module as string,
      path: body.path as string,
      method: body.method as ResourceApi['method'],
      sortOrder: (body.sortOrder as number) ?? 0,
      isEnabled: body.isEnabled as boolean,
      remark: (body.remark as string) ?? '',
      canWrite: true,
      canDelete: true,
    })
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/resource/api/update'), async ({ request }) => {
    const body = await request.json() as Record<string, unknown>
    const id = body.id as number
    const idx = apis.findIndex(i => i.id === id)
    if (idx === -1) { return HttpResponse.json(fail('API资源不存在')) }
    apis[idx] = { ...apis[idx], ...body }
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/resource/api/del'), async ({ request }) => {
    const body = await request.json() as { ids: number[] }
    apis = apis.filter(i => !body.ids.includes(i.id))
    return HttpResponse.json(success(undefined))
  }),
]
