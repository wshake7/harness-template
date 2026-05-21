import type { PagingRequest } from '@vp/core'
import type { ResourceMenu, ResourceMenuNode } from '~/api/business/sysResourceMenu'
import { http, HttpResponse } from 'msw'
import { fail, success, url } from '.'

let menuId = 16

let menus: ResourceMenu[] = [
  { id: 1, parentID: undefined, treePath: '1', menuType: 'CATALOG', path: '/system', redirect: '', alias: 'system', name: '系统管理', component: '', metadata: { icon: 'Settings', order: 1 }, apiIDs: [], sortOrder: 1, isEnabled: true, remark: '' },
  { id: 2, parentID: 1, treePath: '1,2', menuType: 'MENU', path: '/system/dict', redirect: '', alias: 'dict', name: '字典管理', component: 'system/dict/index', metadata: { icon: 'Book', order: 1 }, apiIDs: [1], sortOrder: 1, isEnabled: true, remark: '' },
  { id: 3, parentID: 1, treePath: '1,3', menuType: 'MENU', path: '/system/language', redirect: '', alias: 'language', name: '语言管理', component: 'system/language/index', metadata: { icon: 'Global', order: 2 }, apiIDs: [1], sortOrder: 2, isEnabled: true, remark: '' },
  { id: 4, parentID: 1, treePath: '1,4', menuType: 'MENU', path: '/system/resource/menu', redirect: '', alias: 'resource-menu', name: '菜单管理', component: 'system/resource/menu/index', metadata: { icon: 'Menu', order: 3 }, apiIDs: [1], sortOrder: 3, isEnabled: true, remark: '' },
  { id: 5, parentID: 1, treePath: '1,5', menuType: 'MENU', path: '/system/resource/api', redirect: '', alias: 'resource-api', name: 'API管理', component: 'system/resource/api/index', metadata: { icon: 'Api', order: 4 }, apiIDs: [1], sortOrder: 4, isEnabled: true, remark: '' },
  { id: 6, parentID: undefined, treePath: '6', menuType: 'CATALOG', path: '/account', redirect: '', alias: 'account', name: '权限管理', component: '', metadata: { icon: 'Safety', order: 2 }, apiIDs: [], sortOrder: 2, isEnabled: true, remark: '' },
  { id: 7, parentID: 6, treePath: '6,7', menuType: 'MENU', path: '/account/user', redirect: '', alias: 'user', name: '用户管理', component: 'account/user/index', metadata: { icon: 'User', order: 1 }, apiIDs: [1], sortOrder: 1, isEnabled: true, remark: '' },
  { id: 8, parentID: 6, treePath: '6,8', menuType: 'MENU', path: '/account/role', redirect: '', alias: 'role', name: '角色管理', component: 'account/role/index', metadata: { icon: 'Shield', order: 2 }, apiIDs: [1], sortOrder: 2, isEnabled: true, remark: '' },
  { id: 9, parentID: undefined, treePath: '9', menuType: 'CATALOG', path: '/job', redirect: '', alias: 'job', name: '任务调度', component: '', metadata: { icon: 'ClockCircle', order: 3 }, apiIDs: [], sortOrder: 3, isEnabled: true, remark: '' },
  { id: 10, parentID: 9, treePath: '9,10', menuType: 'MENU', path: '/job/schedule', redirect: '', alias: 'schedule', name: '调度管理', component: 'job/schedule/index', metadata: { icon: 'Schedule', order: 1 }, apiIDs: [1], sortOrder: 1, isEnabled: true, remark: '' },
  { id: 11, parentID: 9, treePath: '9,11', menuType: 'MENU', path: '/job/execution', redirect: '', alias: 'execution', name: '执行记录', component: 'job/execution/index', metadata: { icon: 'History', order: 2 }, apiIDs: [1], sortOrder: 2, isEnabled: true, remark: '' },
  { id: 12, parentID: undefined, treePath: '12', menuType: 'CATALOG', path: '/logger', redirect: '', alias: 'logger', name: '系统日志', component: '', metadata: { icon: 'FileText', order: 4 }, apiIDs: [], sortOrder: 4, isEnabled: true, remark: '' },
  { id: 13, parentID: 12, treePath: '12,13', menuType: 'MENU', path: '/logger/api/log', redirect: '', alias: 'api-log', name: 'API日志', component: 'logger/api/log/index', metadata: { icon: 'Code', order: 1 }, apiIDs: [1], sortOrder: 1, isEnabled: true, remark: '' },
  { id: 14, parentID: 12, treePath: '12,14', menuType: 'MENU', path: '/logger/login/log', redirect: '', alias: 'login-log', name: '登录日志', component: 'logger/login/log/index', metadata: { icon: 'Login', order: 2 }, apiIDs: [1], sortOrder: 2, isEnabled: true, remark: '' },
  { id: 15, parentID: undefined, treePath: '15', menuType: 'MENU', path: '/dashboard', redirect: '', alias: 'dashboard', name: '数据看板', component: 'dashboard/index', metadata: { icon: 'Dashboard', order: 0 }, apiIDs: [1], sortOrder: 0, isEnabled: true, remark: '' },
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

function toMenuNodes(items: ResourceMenu[]): ResourceMenuNode[] {
  const map = new Map<number, ResourceMenuNode>()
  const roots: ResourceMenuNode[] = []
  for (const item of items) {
    map.set(item.id, {
      id: item.id,
      parentID: item.parentID,
      menuType: item.menuType,
      path: item.path,
      redirect: item.redirect,
      name: item.name,
      component: item.component,
      icon: item.metadata.icon ?? '',
      order: item.metadata.order ?? 0,
      sortOrder: item.sortOrder,
      hidden: item.metadata.hidden ?? false,
      authorities: item.metadata.authorities ?? [],
      isUrl: item.menuType === 'LINK',
      children: [],
    })
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

export const resourceMenuHandlers = [
  http.post(url('/api/sys/resource/menu/list'), async ({ request }) => {
    const req = (await request.json()) as PagingRequest
    let items = [...menus]
    if (req.query) {
      try {
        const q = JSON.parse(req.query) as Record<string, string>
        if (q.name) { items = items.filter(i => i.name.includes(q.name)) }
        if (q.path) { items = items.filter(i => i.path.includes(q.path)) }
      }
      catch { /* ignore */ }
    }
    return HttpResponse.json(success(paginate(items, req)))
  }),

  http.get(url('/api/sys/resource/menu/tree'), async () => {
    return HttpResponse.json(success(toMenuNodes([...menus])))
  }),

  http.post(url('/api/sys/resource/menu/create'), async ({ request }) => {
    const body = await request.json() as Record<string, unknown>
    const parentID = body.parentID as number | undefined
    let treePath = String(menuId)
    if (parentID) {
      const parent = menus.find(m => m.id === parentID)
      if (parent) { treePath = `${parent.treePath},${menuId}` }
    }
    menus.push({
      id: menuId++,
      parentID,
      treePath,
      menuType: body.menuType as ResourceMenu['menuType'],
      path: body.path as string,
      redirect: (body.redirect as string) ?? '',
      alias: (body.alias as string) ?? '',
      name: body.name as string,
      component: (body.component as string) ?? '',
      metadata: (body.metadata as ResourceMenu['metadata']) ?? {},
      apiIDs: (body.apiIDs as number[]) ?? [],
      sortOrder: (body.sortOrder as number) ?? 0,
      isEnabled: body.isEnabled as boolean,
      remark: (body.remark as string) ?? '',
      canWrite: true,
      canDelete: true,
    })
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/resource/menu/update'), async ({ request }) => {
    const body = await request.json() as Record<string, unknown>
    const id = body.id as number
    const idx = menus.findIndex(i => i.id === id)
    if (idx === -1) { return HttpResponse.json(fail('菜单不存在')) }
    menus[idx] = { ...menus[idx], ...body }
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/resource/menu/del'), async ({ request }) => {
    const body = await request.json() as { ids: number[] }
    menus = menus.filter(i => !body.ids.includes(i.id))
    return HttpResponse.json(success(undefined))
  }),
]
