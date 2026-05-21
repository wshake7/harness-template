import type { PagingRequest, PagingResult } from '@vp/core'
import API from '../index'

export type JobScheduleType = 'ONCE' | 'CRON' | 'INTERVAL'
export type JobScheduleStatus = 'ENABLED' | 'DISABLED' | 'DELETED'

export interface JobScheduleOption {
  label: string
  value: string
}

export interface JobScheduleOptions {
  workflowTypes: JobScheduleOption[]
  taskQueues: JobScheduleOption[]
  defaultTaskQueue: string
}

export interface JobSchedule {
  id: number
  jobCode: string
  jobName: string
  workflowType: string
  taskQueue: string
  scheduleType: JobScheduleType
  cronExpr?: string
  intervalSeconds?: number | null
  startTime?: string | null
  endTime?: string | null
  inputJSON?: unknown
  status: JobScheduleStatus
  temporalScheduleID?: string
  temporalWorkflowIDPrefix?: string
  description?: string
  createdAt?: string
  updatedAt?: string
  canWrite?: boolean
  canDelete?: boolean
}

export interface ReqJobScheduleCreate {
  jobCode: string
  jobName: string
  workflowType: string
  taskQueue: string
  scheduleType: JobScheduleType
  cronExpr?: string
  intervalSeconds?: number | null
  startTime?: string | null
  endTime?: string | null
  inputJSON?: string
  status: JobScheduleStatus
  description?: string
}

export interface ReqJobScheduleUpdate extends Partial<ReqJobScheduleCreate> {
  id: number
}

export interface ReqJobScheduleID {
  id: number
}

export interface ReqJobScheduleSwitch extends ReqJobScheduleID {
  enabled: boolean
}

function list(req: PagingRequest) {
  return API.Post<Res<PagingResult<JobSchedule>>>('/api/sys/job/schedule/list', req, {
    cacheFor: 0,
  })
}

async function options() {
  return await API.Post<Res<JobScheduleOptions>>('/api/sys/job/schedule/options', {}, {
    cacheFor: 0,
  }).send()
}

async function detail(req: ReqJobScheduleID) {
  return await API.Post<Res<JobSchedule>>('/api/sys/job/schedule/detail', req, {
    cacheFor: 0,
  }).send()
}

async function create(req: ReqJobScheduleCreate) {
  await API.Post<Res>('/api/sys/job/schedule/create', req, {
    cacheFor: 0,
  }).send()
}

async function update(req: ReqJobScheduleUpdate) {
  await API.Post<Res>('/api/sys/job/schedule/update', req, {
    cacheFor: 0,
  }).send()
}

async function del(req: ReqJobScheduleID) {
  await API.Post<Res>('/api/sys/job/schedule/del', req, {
    cacheFor: 0,
  }).send()
}

async function switchStatus(req: ReqJobScheduleSwitch) {
  await API.Post<Res>('/api/sys/job/schedule/switch', req, {
    cacheFor: 0,
  }).send()
}

async function sync(req: ReqJobScheduleID) {
  await API.Post<Res>('/api/sys/job/schedule/sync', req, {
    cacheFor: 0,
  }).send()
}

async function trigger(req: ReqJobScheduleID) {
  await API.Post<Res>('/api/sys/job/schedule/trigger', req, {
    cacheFor: 0,
  }).send()
}

export const JobScheduleApi = {
  list,
  options,
  detail,
  create,
  update,
  del,
  switchStatus,
  sync,
  trigger,
}
