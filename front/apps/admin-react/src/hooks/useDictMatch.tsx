import { createDictMatchHooks } from '@vp/react-core'
import { DictApi } from '~/api/business/sysDict'
import { renderDictEntryLabel } from '~/components/dictEntryLabel'
import i18n from '~/i18n'

export type { DictMatchesResult, DictMatchResult } from '@vp/react-core'

const dictMatchHooks = createDictMatchHooks({
  i18n,
  entryMatch: req => DictApi.entryMatch(req),
  renderEntryLabel: renderDictEntryLabel,
})

export const useDictMatch = dictMatchHooks.useDictMatch
export const useDictMatches = dictMatchHooks.useDictMatches
