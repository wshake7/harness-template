import type { DictMatchedEntriesByCode, ReqDictEntryMatch } from '@vp/core'
import API from '../index'

export type { DictMatchedEntriesByCode, DictMatchedEntry, ReqDictEntryMatch } from '@vp/core'

function entryMatch(req: ReqDictEntryMatch) {
  return API.Post<Res<DictMatchedEntriesByCode>>('/api/sys/dict/entry/match', req, {
    cacheFor: 0,
  })
}

export const DictApi = {
  entryMatch,
}
