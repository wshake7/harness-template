import type { ProColumns } from '@ant-design/pro-components'
import type { Dayjs } from 'dayjs'
import type { JobSchedule, JobScheduleOptions, JobScheduleStatus, JobScheduleType } from '~/api/business/jobSchedule'
import {
  ProFormDigit,
  ProFormSelect,
  ProFormText,
  ProFormTextArea,
  ProTable,
} from '@ant-design/pro-components'
import { createFileRoute } from '@tanstack/react-router'
import { DEFAULT_PAGE_SIZE } from '@vp/core'
import { formatDateYYYYMMDDHHmmss } from '@vp/utils'
import { usePagination } from 'alova/client'
import { AutoComplete, Button, DatePicker, Drawer, Form, Input, Popconfirm, Select, Space, Tag } from 'antd'
import dayjs from 'dayjs'
import { useCallback, useEffect, useMemo, useState } from 'react'
import z from 'zod'
import { JobScheduleApi } from '~/api/business/jobSchedule'
import { useDictMatch } from '~/hooks/useDictMatch'
import { gMessage } from '~/utils/message'

export const Route = createFileRoute('/_app/job/schedule')({
  staleTime: 1000 * 60 * 2,
  component: JobScheduleManagement,
})

const scheduleTypeOptions: { label: string, value: JobScheduleType }[] = [
  { label: '单次', value: 'ONCE' },
  { label: 'Cron', value: 'CRON' },
  { label: '固定间隔', value: 'INTERVAL' },
]

const emptyJobScheduleOptions: JobScheduleOptions = {
  workflowTypes: [],
  taskQueues: [],
  defaultTaskQueue: 'admin',
}

const JobScheduleFormSchema = z.object({
  jobCode: z.string(),
  jobName: z.string(),
  workflowType: z.string(),
  taskQueue: z.string(),
  scheduleType: z.enum(['ONCE', 'CRON', 'INTERVAL']),
  cronExpr: z.string().optional(),
  intervalSeconds: z.number().optional().nullable(),
  startTime: z.custom<Dayjs>().optional().nullable(),
  endTime: z.custom<Dayjs>().optional().nullable(),
  inputJSON: z.string().optional(),
  status: z.enum(['ENABLED', 'DISABLED']),
  description: z.string().optional(),
})

type JobScheduleFormValues = z.infer<typeof JobScheduleFormSchema>

const defaultFormValues: JobScheduleFormValues = {
  jobCode: '',
  jobName: '',
  workflowType: '',
  taskQueue: '',
  scheduleType: 'CRON',
  cronExpr: '',
  intervalSeconds: undefined,
  startTime: undefined,
  endTime: undefined,
  inputJSON: '',
  status: 'ENABLED',
  description: '',
}

const JobScheduleSubmitSchema = JobScheduleFormSchema.superRefine((values, ctx) => {
  for (const [field, label] of [
    ['jobCode', '任务编码'],
    ['jobName', '任务名称'],
    ['workflowType', 'Workflow 类型'],
    ['taskQueue', 'Task Queue'],
  ] as const) {
    if (!values[field]?.trim()) {
      ctx.addIssue({ code: 'custom', path: [field], message: `请输入${label}` })
    }
  }
  if (values.scheduleType === 'CRON' && !values.cronExpr?.trim()) {
    ctx.addIssue({ code: 'custom', path: ['cronExpr'], message: '请输入 cron 表达式' })
  }
  if (values.scheduleType === 'INTERVAL' && (!values.intervalSeconds || values.intervalSeconds <= 0)) {
    ctx.addIssue({ code: 'custom', path: ['intervalSeconds'], message: '间隔秒数必须大于 0' })
  }
  if (values.scheduleType === 'ONCE' && !values.startTime) {
    ctx.addIssue({ code: 'custom', path: ['startTime'], message: '请选择开始时间' })
  }
  if (values.endTime && values.startTime && !values.endTime.isAfter(values.startTime)) {
    ctx.addIssue({ code: 'custom', path: ['endTime'], message: '结束时间必须晚于开始时间' })
  }
  const input = values.inputJSON?.trim()
  if (input) {
    try {
      JSON.parse(input)
    }
    catch {
      ctx.addIssue({ code: 'custom', path: ['inputJSON'], message: '请输入合法 JSON' })
    }
  }
})

const enabledStatusValue = (status: JobScheduleStatus) => status === 'ENABLED' ? '1' : '0'
const scheduleStatusFromEnabledValue = (value: string): Extract<JobScheduleStatus, 'ENABLED' | 'DISABLED'> =>
  value === '1' ? 'ENABLED' : 'DISABLED'
const fallbackEnabledStatusLabel = (status: JobScheduleStatus) => status === 'ENABLED' ? '启用' : '停用'

function scheduleTypeTag(type: JobScheduleType) {
  const map: Record<JobScheduleType, { color: string, label: string }> = {
    ONCE: { color: 'purple', label: '单次' },
    CRON: { color: 'blue', label: 'Cron' },
    INTERVAL: { color: 'green', label: '固定间隔' },
  }
  const item = map[type]
  return <Tag color={item.color}>{item.label}</Tag>
}

function stringifyJSON(value: unknown) {
  if (value === undefined || value === null || value === '') {
    return ''
  }
  if (typeof value === 'string') {
    try {
      return JSON.stringify(JSON.parse(value), null, 2)
    }
    catch {
      return value
    }
  }
  return JSON.stringify(value, null, 2)
}

function toDayjs(value?: string | null) {
  if (!value) {
    return undefined
  }
  const date = dayjs(value)
  return date.isValid() ? date : undefined
}

function toFormValues(record: JobSchedule): JobScheduleFormValues {
  return {
    ...defaultFormValues,
    jobCode: record.jobCode ?? '',
    jobName: record.jobName ?? '',
    workflowType: record.workflowType ?? '',
    taskQueue: record.taskQueue ?? '',
    scheduleType: record.scheduleType,
    cronExpr: record.cronExpr ?? '',
    intervalSeconds: record.intervalSeconds ?? undefined,
    startTime: toDayjs(record.startTime),
    endTime: toDayjs(record.endTime),
    inputJSON: stringifyJSON(record.inputJSON),
    status: record.status === 'DELETED' ? 'DISABLED' : record.status,
    description: record.description ?? '',
  }
}

function toPayload(values: JobScheduleFormValues) {
  const input = values.inputJSON?.trim()
  return {
    jobCode: values.jobCode.trim(),
    jobName: values.jobName.trim(),
    workflowType: values.workflowType.trim(),
    taskQueue: values.taskQueue.trim(),
    scheduleType: values.scheduleType,
    cronExpr: values.scheduleType === 'CRON' ? values.cronExpr?.trim() ?? '' : '',
    intervalSeconds: values.scheduleType === 'INTERVAL' ? values.intervalSeconds ?? null : null,
    startTime: values.startTime ? values.startTime.format('YYYY-MM-DD HH:mm:ss') : null,
    endTime: values.endTime ? values.endTime.format('YYYY-MM-DD HH:mm:ss') : null,
    inputJSON: input ? JSON.stringify(JSON.parse(input), null, 2) : '',
    status: values.status,
    description: values.description?.trim() ?? '',
  }
}

function filterAutoCompleteOption(input: string, option?: { label?: unknown, value?: unknown }) {
  const keyword = input.toLowerCase()
  return String(option?.label ?? '').toLowerCase().includes(keyword)
    || String(option?.value ?? '').toLowerCase().includes(keyword)
}

function JobScheduleManagement() {
  const [drawerOpen, setDrawerOpen] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [editing, setEditing] = useState<JobSchedule>()
  const [searchText, setSearchText] = useState('')
  const [scheduleTypeFilter, setScheduleTypeFilter] = useState<JobScheduleType | undefined>()
  const [statusFilter, setStatusFilter] = useState<JobScheduleStatus | undefined>()
  const [jobScheduleOptions, setJobScheduleOptions] = useState<JobScheduleOptions>(emptyJobScheduleOptions)
  const [form] = Form.useForm<JobScheduleFormValues>()
  const watchedScheduleType = Form.useWatch('scheduleType', form) ?? 'CRON'
  const enabledStatus = useDictMatch(DictCode.SYS_IS_ENABLED_DICT_CODE)

  useEffect(() => {
    let ignore = false
    JobScheduleApi.options()
      .then((response) => {
        if (!ignore && response.data) {
          setJobScheduleOptions(response.data)
        }
      })
      .catch(() => {
        gMessage.error('加载任务选项失败')
      })
    return () => {
      ignore = true
    }
  }, [])

  useEffect(() => {
    if (!drawerOpen || editing || !jobScheduleOptions.defaultTaskQueue || form.getFieldValue('taskQueue')) {
      return
    }
    form.setFieldValue('taskQueue', jobScheduleOptions.defaultTaskQueue)
  }, [drawerOpen, editing, form, jobScheduleOptions.defaultTaskQueue])

  const {
    data,
    total,
    page,
    pageSize,
    loading,
    update,
    send,
  } = usePagination(
    (nextPage, nextPageSize) => {
      const filters: Record<string, unknown>[] = []
      const keyword = searchText.trim()
      if (keyword) {
        filters.push({
          $or: [
            { jobCode__icontains: keyword },
            { jobName__icontains: keyword },
            { workflowType__icontains: keyword },
            { taskQueue__icontains: keyword },
          ],
        })
      }
      if (scheduleTypeFilter) {
        filters.push({ scheduleType: scheduleTypeFilter })
      }
      if (statusFilter) {
        filters.push({ status: statusFilter })
      }
      return JobScheduleApi.list({
        page: nextPage,
        pageSize: nextPageSize,
        orderBy: 'id desc',
        query: filters.length > 0 ? JSON.stringify({ $and: filters }) : undefined,
      })
    },
    {
      initialData: { total: 0, items: [] },
      initialPage: 1,
      initialPageSize: DEFAULT_PAGE_SIZE,
      total: response => response.data?.total ?? 0,
      data: response => response.data?.items ?? [],
      watchingStates: [searchText, scheduleTypeFilter, statusFilter],
      debounce: [500, 0, 0],
    },
  )

  const statusOptions = useMemo(() => {
    const options = enabledStatus.entries
      .filter(entry => entry.entryValue === '1' || entry.entryValue === '0')
      .map(entry => ({
        label: enabledStatus.getLabel(entry.entryValue, entry.entryLabel),
        value: scheduleStatusFromEnabledValue(entry.entryValue),
      }))

    return options.length > 0
      ? options
      : [
          { label: fallbackEnabledStatusLabel('ENABLED'), value: 'ENABLED' as const },
          { label: fallbackEnabledStatusLabel('DISABLED'), value: 'DISABLED' as const },
        ]
  }, [enabledStatus])

  const { rules, onFinish } = useZodForm<JobScheduleFormValues>({
    form,
    schema: JobScheduleSubmitSchema,
    async onSubmit(values) {
      if (!values) {
        return
      }
      setSubmitting(true)
      try {
        const payload = toPayload(values)
        if (editing) {
          await JobScheduleApi.update({ id: editing.id, ...payload })
        }
        else {
          await JobScheduleApi.create(payload)
        }
        gMessage.success('保存成功')
        setDrawerOpen(false)
        await send()
      }
      catch {
        gMessage.error('保存失败')
      }
      finally {
        setSubmitting(false)
      }
    },
  })

  const openCreateForm = () => {
    setEditing(undefined)
    form.resetFields()
    form.setFieldsValue({ ...defaultFormValues, taskQueue: jobScheduleOptions.defaultTaskQueue })
    setDrawerOpen(true)
  }

  const openEditForm = useCallback((record: JobSchedule) => {
    setEditing(record)
    form.resetFields()
    form.setFieldsValue(toFormValues(record))
    setDrawerOpen(true)
  }, [form])

  const action = useCallback(async (fn: () => Promise<void>, success: string, fail: string) => {
    try {
      await fn()
      gMessage.success(success)
      await send()
    }
    catch {
      gMessage.error(fail)
    }
  }, [send])

  const columns: ProColumns<JobSchedule>[] = useMemo(() => [
    { title: '任务编码', dataIndex: 'jobCode', width: 160, ellipsis: true },
    { title: '任务名称', dataIndex: 'jobName', width: 180, ellipsis: true },
    { title: 'Workflow 类型', dataIndex: 'workflowType', width: 220, ellipsis: true },
    { title: 'Task Queue', dataIndex: 'taskQueue', width: 160, ellipsis: true },
    {
      title: '调度类型',
      dataIndex: 'scheduleType',
      width: 100,
      render: (_, record) => scheduleTypeTag(record.scheduleType),
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 90,
      render: (_, record) => record.status === 'DELETED'
        ? <Tag color="error">已删除</Tag>
        : enabledStatus.renderLabel(
            enabledStatusValue(record.status),
            <Tag color={record.status === 'ENABLED' ? 'success' : 'default'}>{fallbackEnabledStatusLabel(record.status)}</Tag>,
          ),
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      width: 160,
      render: (_, record) => formatDateYYYYMMDDHHmmss(record.createdAt),
    },
    {
      title: '操作',
      valueType: 'option',
      width: 260,
      fixed: 'right',
      render: (_, record) => {
        const nextStatus: Extract<JobScheduleStatus, 'ENABLED' | 'DISABLED'> = record.status === 'ENABLED' ? 'DISABLED' : 'ENABLED'
        const nextStatusLabel = enabledStatus.getLabel(enabledStatusValue(nextStatus), fallbackEnabledStatusLabel(nextStatus))

        return (
          <Space>
            <Button
              type="link"
              size="small"
              disabled={record.canWrite === false}
              onClick={() => action(
                () => JobScheduleApi.switchStatus({ id: record.id, enabled: record.status !== 'ENABLED' }),
                `${nextStatusLabel}成功`,
                `${nextStatusLabel}失败`,
              )}
            >
              {nextStatusLabel}
            </Button>
            <Button type="link" size="small" disabled={record.canWrite === false} onClick={() => openEditForm(record)}>
              编辑
            </Button>
            <Button
              type="link"
              size="small"
              disabled={record.canWrite === false}
              onClick={() => action(() => JobScheduleApi.sync({ id: record.id }), '同步成功', '同步失败')}
            >
              同步
            </Button>
            <Popconfirm
              title="确认立即触发该任务？"
              onConfirm={() => action(() => JobScheduleApi.trigger({ id: record.id }), '触发成功', '触发失败')}
            >
              <Button type="link" size="small" disabled={record.canWrite === false}>
                触发
              </Button>
            </Popconfirm>
            <Popconfirm
              title="确认删除该任务配置？"
              onConfirm={() => action(() => JobScheduleApi.del({ id: record.id }), '删除成功', '删除失败')}
            >
              <Button type="link" size="small" danger disabled={record.canDelete === false}>
                删除
              </Button>
            </Popconfirm>
          </Space>
        )
      },
    },
  ], [action, enabledStatus, openEditForm])

  return (
    <>
      <ProTable<JobSchedule>
        rowKey="id"
        headerTitle="任务配置"
        columns={columns}
        dataSource={data}
        loading={loading}
        search={false}
        scroll={{ x: 1500 }}
        pagination={{
          showSizeChanger: true,
          current: page,
          pageSize,
          total,
          onChange: (nextPage, nextPageSize) => update({ page: nextPage, pageSize: nextPageSize }),
        }}
        options={{ reload: () => send() }}
        toolBarRender={() => [
          <Button key="create" type="primary" onClick={openCreateForm}>创建任务</Button>,
          <Input.Search
            key="search"
            allowClear
            placeholder="搜索编码、名称、Workflow、队列"
            value={searchText}
            onChange={event => setSearchText(event.target.value)}
            onSearch={setSearchText}
            style={{ width: 280 }}
          />,
          <Select
            key="scheduleType"
            allowClear
            placeholder="调度类型"
            value={scheduleTypeFilter}
            options={scheduleTypeOptions}
            onChange={setScheduleTypeFilter}
            style={{ width: 130 }}
          />,
          <Select
            key="status"
            allowClear
            placeholder="状态"
            value={statusFilter}
            options={statusOptions}
            onChange={setStatusFilter}
            style={{ width: 120 }}
          />,
        ]}
      />

      <Drawer
        title={editing ? '编辑任务配置' : '创建任务配置'}
        width={560}
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        extra={(
          <Space>
            <Button onClick={() => setDrawerOpen(false)}>取消</Button>
            <Button type="primary" loading={submitting} onClick={() => form.submit()}>保存</Button>
          </Space>
        )}
      >
        <Form form={form} layout="vertical" onFinish={onFinish}>
          <ProFormText name="jobCode" label="任务编码" disabled={Boolean(editing)} fieldProps={{ maxLength: 128 }} rules={rules} />
          <ProFormText name="jobName" label="任务名称" fieldProps={{ maxLength: 255 }} rules={rules} />
          <Form.Item name="workflowType" label="Workflow 类型" rules={rules}>
            <AutoComplete options={jobScheduleOptions.workflowTypes} filterOption={filterAutoCompleteOption}>
              <Input maxLength={255} />
            </AutoComplete>
          </Form.Item>
          <Form.Item name="taskQueue" label="Task Queue" rules={rules}>
            <AutoComplete options={jobScheduleOptions.taskQueues} filterOption={filterAutoCompleteOption}>
              <Input maxLength={255} />
            </AutoComplete>
          </Form.Item>
          <ProFormSelect name="scheduleType" label="调度类型" options={scheduleTypeOptions} rules={rules} />
          {watchedScheduleType === 'CRON' && (
            <ProFormText name="cronExpr" label="Cron 表达式" fieldProps={{ maxLength: 128 }} rules={rules} />
          )}
          {watchedScheduleType === 'INTERVAL' && (
            <ProFormDigit name="intervalSeconds" label="间隔秒数" min={1} precision={0} rules={rules} />
          )}
          {watchedScheduleType === 'ONCE' && (
            <Form.Item name="startTime" label="开始时间" rules={rules}>
              <DatePicker showTime style={{ width: '100%' }} />
            </Form.Item>
          )}
          <Form.Item name="endTime" label="结束时间" rules={rules}>
            <DatePicker showTime style={{ width: '100%' }} />
          </Form.Item>
          <ProFormSelect name="status" label="状态" options={statusOptions} rules={rules} />
          <ProFormTextArea name="inputJSON" label="Workflow 输入 JSON" fieldProps={{ rows: 8 }} rules={rules} />
          <ProFormTextArea name="description" label="描述" fieldProps={{ rows: 3, maxLength: 512 }} />
        </Form>
      </Drawer>
    </>
  )
}
