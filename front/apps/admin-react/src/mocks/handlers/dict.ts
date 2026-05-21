import type { PagingRequest } from '@vp/core'
import type { DictEntry, DictType } from '~/api/business/sysDict'
import { http, HttpResponse } from 'msw'
import { fail, success, url } from '.'

let dictTypeId = 3
let dictEntryId = 6

let dictTypes: DictType[] = [
  { id: 1, typeCode: 'status', typeName: '状态', isEnabled: true, sortOrder: 1, remark: '通用状态字典', canWrite: true, canDelete: true },
  { id: 2, typeCode: 'gender', typeName: '性别', isEnabled: true, sortOrder: 2, remark: '性别字典', canWrite: true, canDelete: true },
]

let dictEntries: DictEntry[] = [
  { id: 1, labelComponent: 'Tag', entryLabel: '启用', entryValue: 'enabled', languageCode: 'zh-CN', sysDictTypeId: 1, sysDictType: { typeCode: 'status', typeName: '状态' }, sortOrder: 1, isEnabled: true, remark: '' },
  { id: 2, labelComponent: 'Tag', entryLabel: '禁用', entryValue: 'disabled', languageCode: 'zh-CN', sysDictTypeId: 1, sysDictType: { typeCode: 'status', typeName: '状态' }, sortOrder: 2, isEnabled: true, remark: '' },
  { id: 3, labelComponent: 'Tag', entryLabel: '男', entryValue: 'male', languageCode: 'zh-CN', sysDictTypeId: 2, sysDictType: { typeCode: 'gender', typeName: '性别' }, sortOrder: 1, isEnabled: true, remark: '' },
  { id: 4, labelComponent: 'Tag', entryLabel: '女', entryValue: 'female', languageCode: 'zh-CN', sysDictTypeId: 2, sysDictType: { typeCode: 'gender', typeName: '性别' }, sortOrder: 2, isEnabled: true, remark: '' },
  { id: 5, labelComponent: 'Tag', entryLabel: '未知', entryValue: 'unknown', languageCode: 'zh-CN', sysDictTypeId: 2, sysDictType: { typeCode: 'gender', typeName: '性别' }, sortOrder: 3, isEnabled: true, remark: '' },
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

export const dictHandlers = [
  http.post(url('/api/sys/dict/type/list'), async ({ request }) => {
    const req = (await request.json()) as PagingRequest
    let items = [...dictTypes]
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

  http.post(url('/api/sys/dict/type/create'), async ({ request }) => {
    const body = await request.json() as Record<string, unknown>
    const item: DictType = {
      id: dictTypeId++,
      typeCode: body.typeCode as string,
      typeName: body.typeName as string,
      isEnabled: body.isEnabled as boolean,
      sortOrder: (body.sortOrder as number) ?? 0,
      remark: (body.remark as string) ?? '',
      canWrite: true,
      canDelete: true,
    }
    dictTypes.push(item)
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/dict/type/update'), async ({ request }) => {
    const body = await request.json() as Record<string, unknown>
    const id = body.id as number
    const idx = dictTypes.findIndex(i => i.id === id)
    if (idx === -1) { return HttpResponse.json(fail('字典类型不存在')) }
    dictTypes[idx] = { ...dictTypes[idx], ...body } as DictType
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/dict/type/del'), async ({ request }) => {
    const body = await request.json() as { ids: number[] }
    dictTypes = dictTypes.filter(i => !body.ids.includes(i.id))
    dictEntries = dictEntries.filter(i => !body.ids.includes(i.sysDictTypeId))
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/dict/entry/list'), async ({ request }) => {
    const req = (await request.json()) as PagingRequest & { sysDictTypeId?: number }
    let items = [...dictEntries]
    if (req.sysDictTypeId) {
      items = items.filter(i => i.sysDictTypeId === req.sysDictTypeId)
    }
    if (req.query) {
      try {
        const q = JSON.parse(req.query) as Record<string, string>
        if (q.entryLabel) { items = items.filter(i => i.entryLabel.includes(q.entryLabel)) }
        if (q.entryValue) { items = items.filter(i => i.entryValue.includes(q.entryValue)) }
      }
      catch { /* ignore */ }
    }
    return HttpResponse.json(success(paginate(items, req)))
  }),

  http.post(url('/api/sys/dict/entry/match'), async ({ request }) => {
    const body = await request.json() as { codes: string[] }
    const result: Record<string, Array<{ id: number, labelComponent: string, entryLabel: string, entryValue: string }>> = {}
    for (const code of body.codes) {
      const type = dictTypes.find(t => t.typeCode === code)
      if (type) {
        result[code] = dictEntries
          .filter(e => e.sysDictTypeId === type.id)
          .map(e => ({ id: e.id, labelComponent: e.labelComponent, entryLabel: e.entryLabel, entryValue: e.entryValue }))
      }
    }
    return HttpResponse.json(success(result))
  }),

  http.post(url('/api/sys/dict/entry/create'), async ({ request }) => {
    const body = await request.json() as Record<string, unknown>
    const type = dictTypes.find(t => t.id === (body.sysDictTypeId as number))
    const item: DictEntry = {
      id: dictEntryId++,
      labelComponent: body.labelComponent as string,
      entryLabel: body.entryLabel as string,
      entryValue: body.entryValue as string,
      languageCode: body.languageCode as string,
      sysDictTypeId: body.sysDictTypeId as number,
      sysDictType: type ? { typeCode: type.typeCode, typeName: type.typeName } : undefined,
      sortOrder: (body.sortOrder as number) ?? 0,
      isEnabled: body.isEnabled as boolean,
      remark: (body.remark as string) ?? '',
    }
    dictEntries.push(item)
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/dict/entry/update'), async ({ request }) => {
    const body = await request.json() as Record<string, unknown>
    const id = body.id as number
    const idx = dictEntries.findIndex(i => i.id === id)
    if (idx === -1) { return HttpResponse.json(fail('字典条目不存在')) }
    const type = dictTypes.find(t => t.id === ((body.sysDictTypeId as number) ?? dictEntries[idx].sysDictTypeId))
    dictEntries[idx] = {
      ...dictEntries[idx],
      ...body,
      sysDictType: type ? { typeCode: type.typeCode, typeName: type.typeName } : dictEntries[idx].sysDictType,
    } as DictEntry
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/dict/entry/del'), async ({ request }) => {
    const body = await request.json() as { ids: number[] }
    dictEntries = dictEntries.filter(i => !body.ids.includes(i.id))
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/dict/entry/batch/copy'), async ({ request }) => {
    const body = await request.json() as { entryIds: number[], targetTypeId: number }
    const type = dictTypes.find(t => t.id === body.targetTypeId)
    for (const entryId of body.entryIds) {
      const entry = dictEntries.find(e => e.id === entryId)
      if (entry) {
        dictEntries.push({
          ...entry,
          id: dictEntryId++,
          sysDictTypeId: body.targetTypeId,
          sysDictType: type ? { typeCode: type.typeCode, typeName: type.typeName } : undefined,
        })
      }
    }
    return HttpResponse.json(success(undefined))
  }),
]
