import { expect, test } from 'vite-plus/test'
import { cn, formatDateYYYYMMDDHHmmss, uriSort } from '../src/index.ts'

test('cn merges class names', () => {
  expect(cn('px-2', false, 'px-4')).toBe('px-4')
})

test('formatDateYYYYMMDDHHmmss formats valid dates', () => {
  expect(formatDateYYYYMMDDHHmmss('2026-05-17T01:02:03Z')).toContain('2026-05-17')
})

test('uriSort sorts and filters query values', () => {
  expect(uriSort({ b: 2, a: 1, c: '' })).toBe('a=1&b=2')
})
