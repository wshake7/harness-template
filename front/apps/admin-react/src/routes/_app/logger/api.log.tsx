import type { ProColumns } from '@ant-design/pro-components'
import type { SysApiLog } from '~/api/business/sysApiLog'
import { ProTable } from '@ant-design/pro-components'
import { createFileRoute } from '@tanstack/react-router'
import { formatDateYYYYMMDDHHmmss } from '@vp/utils'
import { usePagination } from 'alova/client'
import {
  Descriptions,
  Modal,
  Tag,
} from 'antd'
import { useCallback, useState } from 'react'
import { JsonCodeBlock } from '~/components/business/logger/jsonCodeBlock'

export const Route = createFileRoute('/_app/logger/api/log')({
  staleTime: 1000 * 60 * 2,
  component: ApiLogManagement,
})

function successTag(success: boolean) {
  if (success) {
    return <Tag color="success">成功</Tag>
  }
  return <Tag color="error">失败</Tag>
}

function methodTag(method: string) {
  const colorMap: Record<string, string> = {
    GET: 'cyan',
    POST: 'blue',
    PUT: 'orange',
    DELETE: 'red',
    PATCH: 'green',
  }
  return <Tag color={colorMap[method] || 'default'}>{method}</Tag>
}

function costTimeDisplay(costTime: number) {
  if (costTime < 1000) {
    return `${costTime}ms`
  }
  return `${(costTime / 1000).toFixed(2)}s`
}

function DetailText({ children }: { children?: string }) {
  return (
    <span style={{ overflowWrap: 'anywhere', wordBreak: 'break-word' }}>
      {children || '-'}
    </span>
  )
}

function operatorName(record?: SysApiLog | null) {
  return record?.sysUser?.username || (record?.sysUserID ? String(record.sysUserID) : '-')
}

function osName(record?: SysApiLog | null) {
  return record?.oSName || record?.osName || ''
}

function osVersion(record?: SysApiLog | null) {
  return record?.oSVersion || record?.osVersion || ''
}

function operationPlatform(record?: SysApiLog | null) {
  const platform = [
    osName(record),
    record?.browserName,
  ].filter(Boolean).join(' ')

  return platform || '-'
}

function ApiLogManagement() {
  const [detailOpen, setDetailOpen] = useState(false)
  const [detailData, setDetailData] = useState<SysApiLog | null>(null)

  const {
    data,
    total,
    page,
    pageSize,
    loading,
    update,
    send,
  } = usePagination(
    (nextPage, nextPageSize) =>
      ApiLogApi.list({
        page: nextPage,
        pageSize: nextPageSize,
        orderBy: 'id desc',
      }),
    {
      initialData: {
        total: 0,
        items: [],
      },
      initialPage: 1,
      initialPageSize: 10,
      total: response => response.data?.total ?? 0,
      data: response => response.data?.items ?? [],
    },
  )

  const openDetail = useCallback(async (record: SysApiLog) => {
    try {
      const res = await ApiLogApi.detail({ id: record.id })
      if (res.data) {
        setDetailData(res.data)
        setDetailOpen(true)
      }
    }
    catch {
      // ignore
    }
  }, [])

  const columns: ProColumns<SysApiLog>[] = [
    {
      title: '序号',
      dataIndex: 'index',
      width: 60,
      render: (_, __, index) => (page - 1) * pageSize + index + 1,
    },
    {
      title: '模块',
      dataIndex: 'module',
      width: 100,
    },
    {
      title: '方法',
      dataIndex: 'method',
      width: 80,
      render: (_, record) => methodTag(record.method),
    },
    {
      title: '请求路径',
      dataIndex: 'path',
      ellipsis: true,
    },
    {
      title: '请求ID',
      dataIndex: 'requestID',
      width: 220,
      ellipsis: true,
      render: (_, record) => record.requestID || '-',
    },
    {
      title: '操作者',
      dataIndex: ['sysUser', 'username'],
      width: 120,
      render: (_, record) => operatorName(record),
    },
    {
      title: '客户端IP',
      dataIndex: 'clientIP',
      width: 130,
    },
    {
      title: '操作平台',
      dataIndex: 'operationPlatform',
      width: 160,
      ellipsis: true,
      render: (_, record) => operationPlatform(record),
    },
    {
      title: '地理位置',
      dataIndex: 'location',
      width: 240,
      ellipsis: true,
      render: (_, record) => record.location || '-',
    },
    {
      title: '变更项',
      dataIndex: 'formatChange',
      width: 260,
      ellipsis: true,
      render: (_, record) => <DetailText>{record.formatChange}</DetailText>,
    },
    {
      title: '状态码',
      dataIndex: 'statusCode',
      width: 80,
    },
    {
      title: '结果',
      dataIndex: 'success',
      width: 70,
      render: (_, record) => successTag(record.success),
    },
    {
      title: '耗时',
      dataIndex: 'costTime',
      width: 80,
      render: (_, record) => costTimeDisplay(record.costTime),
    },
    {
      title: '操作时间',
      dataIndex: 'createdAt',
      width: 240,
      render: (_, record) => formatDateYYYYMMDDHHmmss(record.createdAt),
    },
    {
      title: '操作',
      valueType: 'option',
      width: 80,
      render: (_, record) => [
        <a
          key="detail"
          onClick={() => {
            openDetail(record)
          }}
        >
          详情
        </a>,
      ],
    },
  ]

  return (
    <>
      <ProTable<SysApiLog>
        rowKey="id"
        search={false}
        columns={columns}
        dataSource={data}
        loading={loading}
        pagination={{
          showSizeChanger: true,
          current: page,
          pageSize,
          total,
          onChange: (nextPage, nextPageSize) => {
            update({
              page: nextPage,
              pageSize: nextPageSize,
            })
          },
        }}
        options={{
          reload: () => send(),
        }}
      />
      <Modal
        title="API日志详情"
        open={detailOpen}
        onCancel={() => {
          setDetailOpen(false)
        }}
        footer={null}
        width={1080}
        style={{ top: 48 }}
        styles={{
          body: {
            maxHeight: 'calc(100vh - 160px)',
            overflowY: 'auto',
          },
        }}
      >
        {detailData && (
          <div>
            <style>
              {`
                .api-log-detail .ant-descriptions-view {
                  overflow: hidden;
                }
                .api-log-detail .ant-descriptions-item-label {
                  width: 96px;
                  min-width: 96px;
                  max-width: 96px;
                  white-space: nowrap;
                }
                .api-log-detail .ant-descriptions-item-content {
                  min-width: 0;
                  max-width: 0;
                }
              `}
            </style>
            <Descriptions className="api-log-detail" column={2} bordered size="small">
              <Descriptions.Item label="ID">{detailData.id}</Descriptions.Item>
              <Descriptions.Item label="请求ID"><DetailText>{detailData.requestID}</DetailText></Descriptions.Item>
              <Descriptions.Item label="模块">{detailData.module}</Descriptions.Item>
              <Descriptions.Item label="方法">{detailData.method}</Descriptions.Item>
              <Descriptions.Item label="请求路径" span={2}><DetailText>{detailData.path}</DetailText></Descriptions.Item>
              <Descriptions.Item label="请求URI" span={2}><DetailText>{detailData.requestURI}</DetailText></Descriptions.Item>
              <Descriptions.Item label="操作者">{operatorName(detailData)}</Descriptions.Item>
              <Descriptions.Item label="客户端IP">{detailData.clientIP}</Descriptions.Item>
              <Descriptions.Item label="状态码">{detailData.statusCode}</Descriptions.Item>
              <Descriptions.Item label="结果">{successTag(detailData.success)}</Descriptions.Item>
              <Descriptions.Item label="耗时">{costTimeDisplay(detailData.costTime)}</Descriptions.Item>
              <Descriptions.Item label="操作时间">{formatDateYYYYMMDDHHmmss(detailData.createdAt)}</Descriptions.Item>
              <Descriptions.Item label="失败原因" span={2}><DetailText>{detailData.reason}</DetailText></Descriptions.Item>
              <Descriptions.Item label="变更项" span={2}><DetailText>{detailData.formatChange}</DetailText></Descriptions.Item>
              <Descriptions.Item label="地理位置" span={2}><DetailText>{detailData.location}</DetailText></Descriptions.Item>
              <Descriptions.Item label="来源" span={2}><DetailText>{detailData.referer}</DetailText></Descriptions.Item>
              <Descriptions.Item label="浏览器" span={2}>
                {detailData.browserName}
                {' '}
                {detailData.browserVersion}
              </Descriptions.Item>
              <Descriptions.Item label="操作系统" span={2}>
                {osName(detailData)}
                {' '}
                {osVersion(detailData)}
              </Descriptions.Item>
              <Descriptions.Item label="客户端" span={2}>
                {detailData.clientName}
                {' '}
                (
                {detailData.clientID}
                )
              </Descriptions.Item>
              <Descriptions.Item label="User-Agent" span={2}><DetailText>{detailData.userAgent}</DetailText></Descriptions.Item>
              <Descriptions.Item label="请求体" span={2}>
                <JsonCodeBlock value={detailData.requestBody} />
              </Descriptions.Item>
              <Descriptions.Item label="请求头" span={2}>
                <JsonCodeBlock value={detailData.requestHeader} />
              </Descriptions.Item>
              <Descriptions.Item label="响应信息" span={2}>
                <JsonCodeBlock value={detailData.response} />
              </Descriptions.Item>
            </Descriptions>
          </div>
        )}
      </Modal>
    </>
  )
}
