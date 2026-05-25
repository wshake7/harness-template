import { expect, test } from 'vite-plus/test'
import z from 'zod'
import {
  fieldZodValidator,
  getFirstIssueMessage,
  getZodIssues,
  globalZodValidator,
  mapErrorFromZodIssue,
  useZodForm,
} from './zod'

test('mapErrorFromZodIssue groups issues by field path', () => {
  const schema = z.object({
    name: z.string().min(1, '名称不能为空'),
    age: z.number().min(0, '年龄不能为负数'),
  })
  const result = schema.safeParse({ name: '', age: -1 })
  expect(result.success).toBe(false)
  if (!result.success) {
    const mapped = mapErrorFromZodIssue(result.error.issues)
    expect(mapped.name).toEqual(['名称不能为空'])
    expect(mapped.age).toEqual(['年龄不能为负数'])
  }
})

test('mapErrorFromZodIssue handles nested paths', () => {
  const schema = z.object({
    user: z.object({
      email: z.string().email('邮箱格式错误'),
    }),
  })
  const result = schema.safeParse({ user: { email: 'invalid' } })
  expect(result.success).toBe(false)
  if (!result.success) {
    const mapped = mapErrorFromZodIssue(result.error.issues)
    expect(mapped['user.email']).toEqual(['邮箱格式错误'])
  }
})

test('getZodIssues extracts issues from ZodError', () => {
  const schema = z.object({ name: z.string().min(1) })
  const result = schema.safeParse({ name: '' })
  expect(result.success).toBe(false)
  if (!result.success) {
    const issues = getZodIssues(result.error)
    expect(issues.length).toBeGreaterThan(0)
    expect(issues[0].path).toEqual(['name'])
  }
})

test('getZodIssues returns empty array for non-ZodError', () => {
  const issues = getZodIssues(new Error('random error'))
  expect(issues).toEqual([])
})

test('getFirstIssueMessage returns first issue message', () => {
  const result = z.object({
    name: z.string().min(1, '名称不能为空'),
    email: z.string().email('邮箱格式错误'),
  }).safeParse({
    name: '',
    email: 'invalid',
  })
  expect(result.success).toBe(false)
  if (result.success) {
    throw new Error('expected zod validation error')
  }
  const issues: z.core.$ZodIssue[] = result.error.issues
  expect(getFirstIssueMessage(issues)).toBe('名称不能为空')
})

test('getFirstIssueMessage returns default for empty array', () => {
  expect(getFirstIssueMessage([])).toBe('请检查表单信息')
})

test('fieldZodValidator and globalZodValidator are exported functions', () => {
  expect(typeof fieldZodValidator).toBe('function')
  expect(typeof globalZodValidator).toBe('function')
})

test('useZodForm is exported function', () => {
  expect(typeof useZodForm).toBe('function')
})
