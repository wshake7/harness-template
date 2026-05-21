import type { DictMatchedEntriesByCode, ReqDictEntryMatch } from '@vp/core'
import { http, HttpResponse } from 'msw'
import { success, url } from '.'

const dictTypes = [
  { id: 1, typeCode: 'status', typeName: '状态' },
  { id: 2, typeCode: 'gender', typeName: '性别' },
]

const dictEntries = [
  { id: 1, labelComponent: 'Tag', entryLabel: '启用', entryValue: 'enabled' },
  { id: 2, labelComponent: 'Tag', entryLabel: '禁用', entryValue: 'disabled' },
  { id: 3, labelComponent: 'Tag', entryLabel: '男', entryValue: 'male' },
  { id: 4, labelComponent: 'Tag', entryLabel: '女', entryValue: 'female' },
  { id: 5, labelComponent: 'Tag', entryLabel: '未知', entryValue: 'unknown' },
]

export const dictHandlers = [
  http.post(url('/api/sys/dict/entry/match'), async ({ request }) => {
    const body = (await request.json()) as ReqDictEntryMatch
    const result: DictMatchedEntriesByCode = {}
    for (const code of body.codes) {
      const type = dictTypes.find(t => t.typeCode === code)
      if (type) {
        result[code] = dictEntries
          .filter(e => e.id <= (type.id === 1 ? 2 : 5) && e.id >= (type.id === 1 ? 1 : 3))
          .map(e => ({ id: e.id, labelComponent: e.labelComponent, entryLabel: e.entryLabel, entryValue: e.entryValue }))
      }
    }
    return HttpResponse.json(success(result))
  }),
]
