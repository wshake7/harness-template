import type { PagingRequest, PagingResult } from '@vp/core'
import API from '../index'

export type JobExecutionStatus = 'RUNNING' | 'SUCCESS' | 'FAILED' | 'CANCELED' | 'TIMEOUT'

export interface JobExecution {
  id: number
  jobCode: string
  temporalWorkflowID: string
  temporalRunID?: string | null
  triggerTime: string
  startTime?: string | null
  endTime?: string | null
  status: JobExecutionStatus
  inputJSON?: unknown
  resultJSON?: unknown
  errorMessage?: string
  retryCount: number
  createdAt?: string
  updatedAt?: string
  canWrite?: boolean
  canDelete?: boolean
}

export interface ReqJobExecutionID {
  id: number
}

function list(req: PagingRequest) {
  return API.Post<Res<PagingResult<JobExecution>>>('/api/sys/job/execution/list', req, {
    cacheFor: 0,
  })
}

async function detail(req: ReqJobExecutionID) {
  return await API.Post<Res<JobExecution>>('/api/sys/job/execution/detail', req, {
    cacheFor: 0,
  }).send()
}

async function cancel(req: ReqJobExecutionID) {
  await API.Post<Res>('/api/sys/job/execution/cancel', req, {
    cacheFor: 0,
  }).send()
}

async function retry(req: ReqJobExecutionID) {
  await API.Post<Res>('/api/sys/job/execution/retry', req, {
    cacheFor: 0,
  }).send()
}

export const JobExecutionApi = {
  list,
  detail,
  cancel,
  retry,
}
