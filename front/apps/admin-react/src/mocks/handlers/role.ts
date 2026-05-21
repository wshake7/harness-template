import type { PagingRequest } from '@vp/core'
import type { SysRole } from '~/api/business/sysRole'
import { http, HttpResponse } from 'msw'
import { fail, success, url } from '.'

let roleId = 4

let roles: SysRole[] = [
  { id: 1, name: '超级管理员', code: 'super_admin', parentID: null, isEnabled: true, remark: '最高权限', createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z', createdBy: 1, updatedBy: 1, canWrite: true, canDelete: false },
  { id: 2, name: '管理员', code: 'admin', parentID: 1, isEnabled: true, remark: '普通管理员', createdAt: '2024-01-02T00:00:00Z', updatedAt: '2024-01-02T00:00:00Z', createdBy: 1, updatedBy: 1, canWrite: true, canDelete: true },
  { id: 3, name: '普通用户', code: 'user', parentID: null, isEnabled: true, remark: '普通用户角色', createdAt: '2024-01-03T00:00:00Z', updatedAt: '2024-01-03T00:00:00Z', createdBy: 1, updatedBy: 1, canWrite: true, canDelete: true },
]

const rolePermissions: Record<number, { menuIDs: number[], apiIDs: number[] }> = {
  1: { menuIDs: [1, 2, 3, 4, 5], apiIDs: [1, 2, 3] },
  2: { menuIDs: [1, 2], apiIDs: [1] },
  3: { menuIDs: [1], apiIDs: [] },
}

function paginate<T>(items: T[], req: PagingRequest): { items: T[], total: number } {
  if (req.noPaging) {
    return { items, total: items.length }
  }
  const page = req.page ?? 1
  const pageSize = req.pageSize ?? 20
  const start = (page - 1) * pageSize
  return { items: items.slice(start, start + pageSize), total: items.length }
}

function buildTree(items: SysRole[]): SysRole[] {
  const map = new Map<number, SysRole>()
  const roots: SysRole[] = []
  for (const item of items) {
    map.set(item.id, { ...item, children: [] })
  }
  for (const item of map.values()) {
    if (item.parentID && map.has(item.parentID)) {
      const parent = map.get(item.parentID)!
      parent.children = parent.children ?? []
      parent.children.push(item)
    }
    else {
      roots.push(item)
    }
  }
  return roots
}

export const roleHandlers = [
  http.post(url('/api/sys/role/list'), async ({ request }) => {
    const req = (await request.json()) as PagingRequest
    let items = [...roles]
    if (req.query) {
      try {
        const q = JSON.parse(req.query) as Record<string, string>
        if (q.name) { items = items.filter(i => i.name.includes(q.name)) }
        if (q.code) { items = items.filter(i => i.code.includes(q.code)) }
      }
      catch { /* ignore */ }
    }
    return HttpResponse.json(success(paginate(items, req)))
  }),

  http.get(url('/api/sys/role/tree'), async () => {
    return HttpResponse.json(success(buildTree([...roles])))
  }),

  http.get(url('/api/sys/role/:id/permissions'), async ({ params }) => {
    const id = Number(params.id)
    return HttpResponse.json(success(rolePermissions[id] ?? { menuIDs: [], apiIDs: [] }))
  }),

  http.post(url('/api/sys/role/create'), async ({ request }) => {
    const body = await request.json() as Record<string, unknown>
    const now = new Date().toISOString()
    roles.push({
      id: roleId++,
      name: body.name as string,
      code: body.code as string,
      parentID: (body.parentID as number | null | undefined) ?? null,
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

  http.post(url('/api/sys/role/update'), async ({ request }) => {
    const body = await request.json() as Record<string, unknown>
    const id = body.id as number
    const idx = roles.findIndex(i => i.id === id)
    if (idx === -1) { return HttpResponse.json(fail('角色不存在')) }
    roles[idx] = { ...roles[idx], ...body, updatedAt: new Date().toISOString() }
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/role/del'), async ({ request }) => {
    const body = await request.json() as { ids: number[] }
    roles = roles.filter(i => !body.ids.includes(i.id))
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/role/permissions'), async ({ request }) => {
    const body = await request.json() as { id: number, menuIDs: number[], apiIDs: number[] }
    rolePermissions[body.id] = { menuIDs: body.menuIDs, apiIDs: body.apiIDs }
    return HttpResponse.json(success(undefined))
  }),
]
