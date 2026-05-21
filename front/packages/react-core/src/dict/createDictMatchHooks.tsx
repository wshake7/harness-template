import type { ApiResponse, DictMatchedEntriesByCode, DictMatchedEntry, ReqDictEntryMatch } from '@vp/core'
import type { i18n as I18nInstance } from 'i18next'
import type { ReactNode } from 'react'
import { useEffect, useMemo, useSyncExternalStore } from 'react'

export interface DictMatchResult {
  entries: DictMatchedEntry[]
  getEntry: (value: string | number | boolean) => DictMatchedEntry | undefined
  getLabel: (value: string | number | boolean, fallback?: string) => string
  renderLabel: (value: string | number | boolean, fallback?: ReactNode) => ReactNode
}

export type DictMatchesResult = Record<string, DictMatchResult>

export interface DictEntryMatchMethod {
  send: () => Promise<ApiResponse<DictMatchedEntriesByCode>>
}

export interface CreateDictMatchHooksOptions {
  i18n: I18nInstance
  entryMatch: (req: ReqDictEntryMatch) => DictEntryMatchMethod
  renderEntryLabel?: (labelComponent: string | undefined, entryLabel: string) => ReactNode
}

function normalizeDictCodes(codes: string[]) {
  return Array.from(new Set(codes.map(code => code.trim()).filter(Boolean))).sort()
}

export function createDictMatchHooks(options: CreateDictMatchHooksOptions) {
  const dictMatchCache = new Map<string, DictMatchedEntry[]>()
  const loadingDictCodes = new Set<string>()
  const pendingDictCodes = new Set<string>()
  const dictMatchListeners = new Set<() => void>()
  let flushScheduled = false
  let dictMatchSnapshot: DictMatchedEntriesByCode = {}

  function subscribeDictMatches(listener: () => void) {
    dictMatchListeners.add(listener)
    return () => {
      dictMatchListeners.delete(listener)
    }
  }

  function emitDictMatchesChange() {
    dictMatchSnapshot = Object.fromEntries(dictMatchCache)
    dictMatchListeners.forEach(listener => listener())
  }

  function getDictMatchSnapshot() {
    return dictMatchSnapshot
  }

  function subscribeI18nLanguage(listener: () => void) {
    options.i18n.on('languageChanged', listener)
    options.i18n.on('loaded', listener)

    return () => {
      options.i18n.off('languageChanged', listener)
      options.i18n.off('loaded', listener)
    }
  }

  function getI18nLanguageSnapshot() {
    return options.i18n.resolvedLanguage ?? options.i18n.language
  }

  function flushPendingDictMatches() {
    flushScheduled = false
    const requestCodes = Array.from(pendingDictCodes)
      .filter(code => !dictMatchCache.has(code) && !loadingDictCodes.has(code))
    pendingDictCodes.clear()

    if (requestCodes.length === 0) {
      return
    }

    requestCodes.forEach(code => loadingDictCodes.add(code))
    options.entryMatch({ codes: requestCodes })
      .send()
      .then((res) => {
        const entriesByCode = res.data ?? {}
        requestCodes.forEach((code) => {
          dictMatchCache.set(code, entriesByCode[code] ?? [])
        })
      })
      .catch(() => {
        requestCodes.forEach((code) => {
          dictMatchCache.set(code, [])
        })
      })
      .finally(() => {
        requestCodes.forEach(code => loadingDictCodes.delete(code))
        emitDictMatchesChange()
      })
  }

  function scheduleDictMatches(codes: string[]) {
    const missingCodes = normalizeDictCodes(codes)
      .filter(code => !dictMatchCache.has(code) && !loadingDictCodes.has(code))

    if (missingCodes.length === 0) {
      return
    }

    missingCodes.forEach(code => pendingDictCodes.add(code))
    if (!flushScheduled) {
      flushScheduled = true
      queueMicrotask(flushPendingDictMatches)
    }
  }

  function createDictMatchResult(
    entries: DictMatchedEntry[],
    resolveEntryLabel: (entryLabel: string) => string,
  ): DictMatchResult {
    const entryByValue = new Map(entries.map(item => [item.entryValue, item]))

    const getEntry = (value: string | number | boolean) => {
      return entryByValue.get(String(value))
    }

    const getLabel = (value: string | number | boolean, fallback = '') => {
      const entry = getEntry(value)
      return entry ? resolveEntryLabel(entry.entryLabel) : fallback
    }

    const renderLabel = (value: string | number | boolean, fallback: ReactNode = '未知状态') => {
      const entry = getEntry(value)
      if (!entry) {
        return fallback
      }
      return options.renderEntryLabel?.(entry.labelComponent, resolveEntryLabel(entry.entryLabel))
        ?? resolveEntryLabel(entry.entryLabel)
    }

    return {
      entries,
      getEntry,
      getLabel,
      renderLabel,
    }
  }

  function useDictMatches(codes: string[]) {
    const language = useSyncExternalStore(
      subscribeI18nLanguage,
      getI18nLanguageSnapshot,
      getI18nLanguageSnapshot,
    )
    const codesKey = useMemo(() => normalizeDictCodes(codes).join('\u0000'), [codes])
    const cachedEntriesByCode = useSyncExternalStore(
      subscribeDictMatches,
      getDictMatchSnapshot,
      getDictMatchSnapshot,
    )

    useEffect(() => {
      const requestCodes = codesKey ? codesKey.split('\u0000') : []
      scheduleDictMatches(requestCodes)
    }, [codesKey])

    return useMemo(() => {
      const requestCodes = codesKey ? codesKey.split('\u0000') : []
      const resolveEntryLabel = (entryLabel: string) => {
        const translated = options.i18n.t(entryLabel, { defaultValue: entryLabel, lng: language })
        return typeof translated === 'string' ? translated : entryLabel
      }

      return Object.fromEntries(
        requestCodes.map(code => [
          code,
          createDictMatchResult(cachedEntriesByCode[code] ?? [], resolveEntryLabel),
        ]),
      ) as DictMatchesResult
    }, [cachedEntriesByCode, codesKey, language])
  }

  function useDictMatch(code: string) {
    const codes = useMemo(() => [code], [code])
    const matches = useDictMatches(codes)

    return useMemo(() => {
      return matches[code] ?? createDictMatchResult([], entryLabel => entryLabel)
    }, [code, matches])
  }

  return {
    useDictMatch,
    useDictMatches,
  }
}
