import type { PagingRequest } from '@vp/core'
import type { JobSchedule, JobScheduleOptions } from '~/api/business/jobSchedule'
import { http, HttpResponse } from 'msw'
import { fail, success, url } from '.'

let scheduleId = 4

let schedules: JobSchedule[] = [
  { id: 1, jobCode: 'cleanup', jobName: '数据清理', workflowType: 'cleanup', taskQueue: 'default', scheduleType: 'CRON', cronExpr: '0 0 * * *', status: 'ENABLED', description: '每日清理过期数据', createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z', canWrite: true, canDelete: true },
  { id: 2, jobCode: 'backup', jobName: '数据备份', workflowType: 'backup', taskQueue: 'default', scheduleType: 'ONCE', startTime: '2024-12-31T00:00:00Z', status: 'ENABLED', description: '一次性全量备份', createdAt: '2024-01-02T00:00:00Z', updatedAt: '2024-01-02T00:00:00Z', canWrite: true, canDelete: true },
  { id: 3, jobCode: 'sync', jobName: '数据同步', workflowType: 'sync', taskQueue: 'high', scheduleType: 'INTERVAL', intervalSeconds: 3600, status: 'DISABLED', description: '每小时同步一次', createdAt: '2024-01-03T00:00:00Z', updatedAt: '2024-01-03T00:00:00Z', canWrite: true, canDelete: true },
]

const scheduleOptions: JobScheduleOptions = {
  workflowTypes: [
    { label: '清理', value: 'cleanup' },
    { label: '备份', value: 'backup' },
    { label: '同步', value: 'sync' },
  ],
  taskQueues: [
    { label: '默认队列', value: 'default' },
    { label: '高优先级', value: 'high' },
  ],
  defaultTaskQueue: 'default',
}

function paginate<T>(items: T[], req: PagingRequest): { items: T[], total: number } {
  if (req.noPaging) {
    return { items, total: items.length }
  }
  const page = req.page ?? 1
  const pageSize = req.pageSize ?? 20
  const start = (page - 1) * pageSize
  return { items: items.slice(start, start + pageSize), total: items.length }
}

export const jobScheduleHandlers = [
  http.post(url('/api/sys/job/schedule/list'), async ({ request }) => {
    const req = (await request.json()) as PagingRequest
    let items = [...schedules]
    if (req.query) {
      try {
        const q = JSON.parse(req.query) as Record<string, string>
        if (q.jobCode) { items = items.filter(i => i.jobCode.includes(q.jobCode)) }
        if (q.jobName) { items = items.filter(i => i.jobName.includes(q.jobName)) }
      }
      catch { /* ignore */ }
    }
    return HttpResponse.json(success(paginate(items, req)))
  }),

  http.post(url('/api/sys/job/schedule/options'), async () => {
    return HttpResponse.json(success(scheduleOptions))
  }),

  http.post(url('/api/sys/job/schedule/detail'), async ({ request }) => {
    const body = await request.json() as { id: number }
    const item = schedules.find(i => i.id === body.id)
    if (!item) { return HttpResponse.json(fail('任务调度不存在')) }
    return HttpResponse.json(success(item))
  }),

  http.post(url('/api/sys/job/schedule/create'), async ({ request }) => {
    const body = await request.json() as Record<string, unknown>
    const now = new Date().toISOString()
    schedules.push({
      id: scheduleId++,
      jobCode: body.jobCode as string,
      jobName: body.jobName as string,
      workflowType: body.workflowType as string,
      taskQueue: body.taskQueue as string,
      scheduleType: body.scheduleType as JobSchedule['scheduleType'],
      cronExpr: body.cronExpr as string | undefined,
      intervalSeconds: (body.intervalSeconds as number | null) ?? null,
      startTime: (body.startTime as string | null) ?? null,
      endTime: (body.endTime as string | null) ?? null,
      inputJSON: body.inputJSON,
      status: body.status as JobSchedule['status'],
      description: (body.description as string) ?? '',
      createdAt: now,
      updatedAt: now,
      canWrite: true,
      canDelete: true,
    })
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/job/schedule/update'), async ({ request }) => {
    const body = await request.json() as Record<string, unknown>
    const id = body.id as number
    const idx = schedules.findIndex(i => i.id === id)
    if (idx === -1) { return HttpResponse.json(fail('任务调度不存在')) }
    schedules[idx] = { ...schedules[idx], ...body, updatedAt: new Date().toISOString() }
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/job/schedule/del'), async ({ request }) => {
    const body = await request.json() as { id: number }
    schedules = schedules.filter(i => i.id !== body.id)
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/job/schedule/switch'), async ({ request }) => {
    const body = await request.json() as { id: number, enabled: boolean }
    const idx = schedules.findIndex(i => i.id === body.id)
    if (idx === -1) { return HttpResponse.json(fail('任务调度不存在')) }
    schedules[idx].status = body.enabled ? 'ENABLED' : 'DISABLED'
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/job/schedule/sync'), async ({ request }) => {
    const body = await request.json() as { id: number }
    const idx = schedules.findIndex(i => i.id === body.id)
    if (idx === -1) { return HttpResponse.json(fail('任务调度不存在')) }
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/sys/job/schedule/trigger'), async ({ request }) => {
    const body = await request.json() as { id: number }
    const idx = schedules.findIndex(i => i.id === body.id)
    if (idx === -1) { return HttpResponse.json(fail('任务调度不存在')) }
    return HttpResponse.json(success(undefined))
  }),
]
