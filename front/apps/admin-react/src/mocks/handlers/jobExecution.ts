import type { PagingRequest } from '@vp/core'
import type { JobExecution } from '~/api/business/jobExecution'
import { http, HttpResponse } from 'msw'
import { fail, success, url } from '.'

const executions: JobExecution[] = Array.from({ length: 15 }, (_, i) => ({
  id: i + 1,
  jobCode: ['cleanup', 'backup', 'sync'][i % 3],
  temporalWorkflowID: `workflow-${1000 + i}`,
  temporalRunID: i % 4 === 0 ? null : `run-${2000 + i}`,
  triggerTime: new Date(Date.now() - i * 1800000).toISOString(),
  startTime: i % 5 === 0 ? null : new Date(Date.now() - i * 1800000 + 1000).toISOString(),
  endTime: i % 5 === 0 ? null : new Date(Date.now() - i * 1800000 + 5000).toISOString(),
  status: ['RUNNING', 'SUCCESS', 'FAILED', 'CANCELED', 'TIMEOUT'][i % 5] as JobExecution['status'],
  inputJSON: {},
  resultJSON: i % 5 === 1 ? { count: 100 } : undefined,
  errorMessage: i % 5 === 2 ? '执行超时' : undefined,
  retryCount: 0,
  createdAt: new Date(Date.now() - i * 1800000).toISOString(),
  updatedAt: new Date(Date.now() - i * 1800000).toISOString(),
  canWrite: true,
  canDelete: true,
}))

function paginate<T>(items: T[], req: PagingRequest): { items: T[], total: number } {
  if (req.noPaging) {
    return { items, total: items.length }
  }
  const page = req.page ?? 1
  const pageSize = req.pageSize ?? 20
  const start = (page - 1) * pageSize
  return { items: items.slice(start, start + pageSize), total: items.length }
}

export const jobExecutionHandlers = [
  http.post(url('/api/sys/job/execution/list'), async ({ request }) => {
    const req = (await request.json()) as PagingRequest
    let items = [...executions]
    if (req.query) {
      try {
        const q = JSON.parse(req.query) as Record<string, string>
        if (q.jobCode) { items = items.filter(i => i.jobCode.includes(q.jobCode)) }
        if (q.status) { items = items.filter(i => i.status === q.status) }
      }
      catch { /* ignore */ }
    }
    return HttpResponse.json(success(paginate(items, req)))
  }),

  http.post(url('/api/sys/job/execution/detail'), async ({ request }) => {
    const body = await request.json() as { id: number }
    const item = executions.find(i => i.id === body.id)
    if (!item) { return HttpResponse.json(fail('任务执行不存在')) }
    return HttpResponse.json(success(item))
  }),

  http.post(url('/api/sys/job/execution/cancel'), async ({ request }) => {
    const body = await request.json() as { id: number }
    const idx = executions.findIndex(i => i.id === body.id)
    if (idx === -1) { return HttpResponse.json(fail('任务执行不存在')) }
    executions[idx].status = 'CANCELED'
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/job/execution/retry'), async ({ request }) => {
    const body = await request.json() as { id: number }
    const idx = executions.findIndex(i => i.id === body.id)
    if (idx === -1) { return HttpResponse.json(fail('任务执行不存在')) }
    executions[idx].retryCount += 1
    executions[idx].status = 'RUNNING'
    return HttpResponse.json(success(undefined))
  }),
]
