import type { ProColumns } from '@ant-design/pro-components'
import type { Dayjs } from 'dayjs'
import type { JobExecution, JobExecutionStatus } from '~/api/business/jobExecution'
import { ProTable } from '@ant-design/pro-components'
import { createFileRoute } from '@tanstack/react-router'
import { DEFAULT_PAGE_SIZE } from '@vp/core'
import { formatDateYYYYMMDDHHmmss } from '@vp/utils'
import { usePagination } from 'alova/client'
import { Button, DatePicker, Descriptions, Input, Modal, Popconfirm, Select, Space, Tag } from 'antd'
import { useCallback, useMemo, useState } from 'react'
import { JobExecutionApi } from '~/api/business/jobExecution'
import { JsonCodeBlock } from '~/components/business/logger/jsonCodeBlock'
import { gMessage } from '~/utils/message'

export const Route = createFileRoute('/_app/job/execution')({
  staleTime: 1000 * 60 * 2,
  component: JobExecutionManagement,
})

const statusOptions: { label: string, value: JobExecutionStatus }[] = [
  { label: '运行中', value: 'RUNNING' },
  { label: '成功', value: 'SUCCESS' },
  { label: '失败', value: 'FAILED' },
  { label: '已取消', value: 'CANCELED' },
  { label: '超时', value: 'TIMEOUT' },
]

function statusTag(status: JobExecutionStatus) {
  const map: Record<JobExecutionStatus, { color: string, label: string }> = {
    RUNNING: { color: 'processing', label: '运行中' },
    SUCCESS: { color: 'success', label: '成功' },
    FAILED: { color: 'error', label: '失败' },
    CANCELED: { color: 'default', label: '已取消' },
    TIMEOUT: { color: 'warning', label: '超时' },
  }
  const item = map[status]
  return <Tag color={item.color}>{item.label}</Tag>
}

function isRetryable(status: JobExecutionStatus) {
  return status === 'FAILED' || status === 'CANCELED' || status === 'TIMEOUT'
}

function stringifyJSON(value: unknown) {
  if (value === undefined || value === null || value === '') {
    return ''
  }
  if (typeof value === 'string') {
    return value
  }
  return JSON.stringify(value, null, 2)
}

function DetailText({ children }: { children?: string | null }) {
  return (
    <span style={{ overflowWrap: 'anywhere', wordBreak: 'break-word' }}>
      {children || '-'}
    </span>
  )
}

function JobExecutionManagement() {
  const [detailOpen, setDetailOpen] = useState(false)
  const [detailData, setDetailData] = useState<JobExecution | null>(null)
  const [searchText, setSearchText] = useState('')
  const [statusFilter, setStatusFilter] = useState<JobExecutionStatus | undefined>()
  const [triggerRange, setTriggerRange] = useState<[Dayjs | null, Dayjs | null] | null>(null)

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
            { temporalWorkflowID__icontains: keyword },
            { temporalRunID__icontains: keyword },
          ],
        })
      }
      if (statusFilter) {
        filters.push({ status: statusFilter })
      }
      if (triggerRange?.[0]) {
        filters.push({ triggerTime__gte: triggerRange[0].format('YYYY-MM-DD HH:mm:ss') })
      }
      if (triggerRange?.[1]) {
        filters.push({ triggerTime__lte: triggerRange[1].format('YYYY-MM-DD HH:mm:ss') })
      }
      return JobExecutionApi.list({
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
      watchingStates: [searchText, statusFilter, triggerRange],
      debounce: [500, 0, 0],
    },
  )

  const openDetail = useCallback(async (record: JobExecution) => {
    try {
      const res = await JobExecutionApi.detail({ id: record.id })
      if (res.data) {
        setDetailData(res.data)
        setDetailOpen(true)
      }
    }
    catch {
      gMessage.error('获取详情失败')
    }
  }, [])

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

  const columns: ProColumns<JobExecution>[] = useMemo(() => [
    { title: '任务编码', dataIndex: 'jobCode', width: 160, ellipsis: true },
    { title: 'Workflow ID', dataIndex: 'temporalWorkflowID', width: 260, ellipsis: true },
    { title: 'Run ID', dataIndex: 'temporalRunID', width: 260, ellipsis: true, render: (_, record) => record.temporalRunID || '-' },
    {
      title: '触发时间',
      dataIndex: 'triggerTime',
      width: 160,
      render: (_, record) => formatDateYYYYMMDDHHmmss(record.triggerTime),
    },
    {
      title: '开始时间',
      dataIndex: 'startTime',
      width: 160,
      render: (_, record) => formatDateYYYYMMDDHHmmss(record.startTime ?? undefined),
    },
    {
      title: '结束时间',
      dataIndex: 'endTime',
      width: 160,
      render: (_, record) => formatDateYYYYMMDDHHmmss(record.endTime ?? undefined),
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (_, record) => statusTag(record.status),
    },
    { title: '重试次数', dataIndex: 'retryCount', width: 90 },
    {
      title: '操作',
      valueType: 'option',
      width: 180,
      fixed: 'right',
      render: (_, record) => (
        <Space>
          <Button type="link" size="small" onClick={() => openDetail(record)}>
            详情
          </Button>
          <Popconfirm
            title="确认取消该 Workflow？"
            disabled={record.status !== 'RUNNING'}
            onConfirm={() => action(() => JobExecutionApi.cancel({ id: record.id }), '取消成功', '取消失败')}
          >
            <Button type="link" size="small" disabled={record.status !== 'RUNNING'}>
              取消
            </Button>
          </Popconfirm>
          <Popconfirm
            title="确认重试该执行记录？"
            disabled={!isRetryable(record.status)}
            onConfirm={() => action(() => JobExecutionApi.retry({ id: record.id }), '重试已发起', '重试失败')}
          >
            <Button type="link" size="small" disabled={!isRetryable(record.status)}>
              重试
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ], [action, openDetail])

  return (
    <>
      <ProTable<JobExecution>
        rowKey="id"
        headerTitle="执行记录"
        columns={columns}
        dataSource={data}
        loading={loading}
        search={false}
        scroll={{ x: 1550 }}
        pagination={{
          showSizeChanger: true,
          current: page,
          pageSize,
          total,
          onChange: (nextPage, nextPageSize) => update({ page: nextPage, pageSize: nextPageSize }),
        }}
        options={{ reload: () => send() }}
        toolBarRender={() => [
          <Input.Search
            key="search"
            allowClear
            placeholder="搜索任务编码、Workflow ID、Run ID"
            value={searchText}
            onChange={event => setSearchText(event.target.value)}
            onSearch={setSearchText}
            style={{ width: 310 }}
          />,
          <Select
            key="status"
            allowClear
            placeholder="状态"
            value={statusFilter}
            options={statusOptions}
            onChange={setStatusFilter}
            style={{ width: 130 }}
          />,
          <DatePicker.RangePicker
            key="triggerRange"
            showTime
            value={triggerRange}
            onChange={setTriggerRange}
            style={{ width: 360 }}
          />,
        ]}
      />

      <Modal
        title="执行记录详情"
        width={820}
        open={detailOpen}
        footer={null}
        onCancel={() => setDetailOpen(false)}
      >
        <Descriptions column={1} bordered size="small">
          <Descriptions.Item label="任务编码">{detailData?.jobCode || '-'}</Descriptions.Item>
          <Descriptions.Item label="Workflow ID">
            <DetailText>{detailData?.temporalWorkflowID}</DetailText>
          </Descriptions.Item>
          <Descriptions.Item label="Run ID">
            <DetailText>{detailData?.temporalRunID}</DetailText>
          </Descriptions.Item>
          <Descriptions.Item label="状态">{detailData ? statusTag(detailData.status) : '-'}</Descriptions.Item>
          <Descriptions.Item label="触发时间">{formatDateYYYYMMDDHHmmss(detailData?.triggerTime)}</Descriptions.Item>
          <Descriptions.Item label="开始时间">{formatDateYYYYMMDDHHmmss(detailData?.startTime ?? undefined)}</Descriptions.Item>
          <Descriptions.Item label="结束时间">{formatDateYYYYMMDDHHmmss(detailData?.endTime ?? undefined)}</Descriptions.Item>
          <Descriptions.Item label="输入参数">
            <JsonCodeBlock value={stringifyJSON(detailData?.inputJSON)} />
          </Descriptions.Item>
          <Descriptions.Item label="执行结果">
            <JsonCodeBlock value={stringifyJSON(detailData?.resultJSON)} />
          </Descriptions.Item>
          <Descriptions.Item label="错误信息">
            <DetailText>{detailData?.errorMessage}</DetailText>
          </Descriptions.Item>
        </Descriptions>
      </Modal>
    </>
  )
}
