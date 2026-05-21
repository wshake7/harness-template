export interface DictMatchedEntry {
  id: number
  labelComponent: string
  entryLabel: string
  entryValue: string
}

export type DictMatchedEntriesByCode = Record<string, DictMatchedEntry[]>

export interface ReqDictEntryMatch {
  codes: string[]
}
