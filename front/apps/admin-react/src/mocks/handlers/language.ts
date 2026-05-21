import type { PagingRequest } from '@vp/core'
import type { LanguageEntry, LanguageType } from '~/api/business/sysLanguage'
import { http, HttpResponse } from 'msw'
import { fail, success, url } from '.'

let langTypeId = 3
let langEntryId = 6

let langTypes: LanguageType[] = [
  { id: 1, typeCode: 'zh-CN', typeName: '简体中文', isDefault: true, sortOrder: 1, isEnabled: true },
  { id: 2, typeCode: 'en-US', typeName: 'English', isDefault: false, sortOrder: 2, isEnabled: true },
]

let langEntries: LanguageEntry[] = [
  { id: 1, entryCode: 'welcome', entryValue: '欢迎', sysLanguageTypeId: 1, sortOrder: 1, isEnabled: true, remark: '' },
  { id: 2, entryCode: 'welcome', entryValue: 'Welcome', sysLanguageTypeId: 2, sortOrder: 1, isEnabled: true, remark: '' },
  { id: 3, entryCode: 'login', entryValue: '登录', sysLanguageTypeId: 1, sortOrder: 2, isEnabled: true, remark: '' },
  { id: 4, entryCode: 'login', entryValue: 'Login', sysLanguageTypeId: 2, sortOrder: 2, isEnabled: true, remark: '' },
  { id: 5, entryCode: 'logout', entryValue: '退出', sysLanguageTypeId: 1, sortOrder: 3, isEnabled: true, remark: '' },
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

export const languageHandlers = [
  http.post(url('/api/sys/language/type/list'), async ({ request }) => {
    const req = (await request.json()) as PagingRequest
    let items = [...langTypes]
    if (req.query) {
      try {
        const q = JSON.parse(req.query) as Record<string, string>
        if (q.typeCode) { items = items.filter(i => i.typeCode.includes(q.typeCode)) }
        if (q.typeName) { items = items.filter(i => i.typeName.includes(q.typeName)) }
      }
      catch { /* ignore */ }
    }
    return HttpResponse.json(success(paginate(items, req)))
  }),

  http.post(url('/api/sys/language/type/create'), async ({ request }) => {
    const body = await request.json() as Record<string, unknown>
    langTypes.push({
      id: langTypeId++,
      typeCode: body.typeCode as string,
      typeName: body.typeName as string,
      isDefault: (body.isDefault as boolean) ?? false,
      sortOrder: (body.sortOrder as number) ?? 0,
      isEnabled: (body.isEnabled as boolean) ?? true,
    })
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/language/type/update'), async ({ request }) => {
    const body = await request.json() as Record<string, unknown>
    const id = body.id as number
    const idx = langTypes.findIndex(i => i.id === id)
    if (idx === -1) { return HttpResponse.json(fail('语言类型不存在')) }
    langTypes[idx] = { ...langTypes[idx], ...body }
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/language/type/del'), async ({ request }) => {
    const body = await request.json() as { ids: number[] }
    langTypes = langTypes.filter(i => !body.ids.includes(i.id))
    langEntries = langEntries.filter(i => !body.ids.includes(i.sysLanguageTypeId))
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/language/entry/list'), async ({ request }) => {
    const req = (await request.json()) as PagingRequest & { sysLanguageTypeId?: number }
    let items = [...langEntries]
    if (req.sysLanguageTypeId) {
      items = items.filter(i => i.sysLanguageTypeId === req.sysLanguageTypeId)
    }
    if (req.query) {
      try {
        const q = JSON.parse(req.query) as Record<string, string>
        if (q.entryCode) { items = items.filter(i => i.entryCode.includes(q.entryCode)) }
        if (q.entryValue) { items = items.filter(i => i.entryValue.includes(q.entryValue)) }
      }
      catch { /* ignore */ }
    }
    return HttpResponse.json(success(paginate(items, req)))
  }),

  http.post(url('/api/sys/language/entry/create'), async ({ request }) => {
    const body = await request.json() as Record<string, unknown>
    const type = langTypes.find(t => t.id === (body.sysLanguageTypeId as number))
    langEntries.push({
      id: langEntryId++,
      entryCode: body.entryCode as string,
      entryValue: body.entryValue as string,
      sysLanguageTypeId: body.sysLanguageTypeId as number,
      sortOrder: (body.sortOrder as number) ?? 0,
      isEnabled: (body.isEnabled as boolean) ?? true,
      remark: (body.remark as string) ?? '',
      sysLanguageType: type ? { typeCode: type.typeCode, typeName: type.typeName } : undefined,
    })
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/language/entry/update'), async ({ request }) => {
    const body = await request.json() as Record<string, unknown>
    const id = (body.id as number) ?? 0
    const idx = id > 0 ? langEntries.findIndex(i => i.id === id) : -1
    if (idx !== -1) {
      const type = langTypes.find(t => t.id === ((body.sysLanguageTypeId as number) ?? langEntries[idx].sysLanguageTypeId))
      langEntries[idx] = {
        ...langEntries[idx],
        ...body,
        sysLanguageType: type ? { typeCode: type.typeCode, typeName: type.typeName } : langEntries[idx].sysLanguageType,
      } as LanguageEntry
    }
    else if (body.updates && Array.isArray(body.updates)) {
      for (const upd of body.updates as Array<Record<string, unknown>>) {
        const uid = upd.id as number
        const uidx = langEntries.findIndex(i => i.id === uid)
        if (uidx !== -1) {
          const type = langTypes.find(t => t.id === ((upd.sysLanguageTypeId as number) ?? langEntries[uidx].sysLanguageTypeId))
          langEntries[uidx] = {
            ...langEntries[uidx],
            ...upd,
            sysLanguageType: type ? { typeCode: type.typeCode, typeName: type.typeName } : langEntries[uidx].sysLanguageType,
          } as LanguageEntry
        }
      }
    }
    else {
      return HttpResponse.json(fail('语言条目不存在'))
    }
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/language/entry/del'), async ({ request }) => {
    const body = await request.json() as { ids: number[] }
    langEntries = langEntries.filter(i => !body.ids.includes(i.id))
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/language/entry/batch/create'), async ({ request }) => {
    const body = await request.json() as Record<string, unknown>
    const entryCode = body.entryCode as string
    const values = body.values as Record<string, string>
    const sortOrder = (body.sortOrder as number) ?? 0
    const isEnabled = (body.isEnabled as boolean) ?? true
    for (const [typeCode, entryValue] of Object.entries(values)) {
      const type = langTypes.find(t => t.typeCode === typeCode)
      if (type) {
        langEntries.push({
          id: langEntryId++,
          entryCode,
          entryValue,
          sysLanguageTypeId: type.id,
          sortOrder,
          isEnabled,
          remark: '',
          sysLanguageType: { typeCode: type.typeCode, typeName: type.typeName },
        })
      }
    }
    return HttpResponse.json(success(undefined))
  }),
]
