import { createDictMatchHooks } from '@vp/react-core'
import { DictApi } from '~/api/business/sysDict'
import i18n from '~/i18n'

export type { DictMatchesResult, DictMatchResult } from '@vp/react-core'

const dictMatchHooks = createDictMatchHooks({
  i18n,
  entryMatch: req => DictApi.entryMatch(req),
})

export const useDictMatch = dictMatchHooks.useDictMatch
export const useDictMatches = dictMatchHooks.useDictMatches
