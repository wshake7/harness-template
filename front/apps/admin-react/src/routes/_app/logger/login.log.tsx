import type { ProColumns } from '@ant-design/pro-components'
import type { SysLoginLog } from '~/api/business/sysLoginLog'
import { ProTable } from '@ant-design/pro-components'
import { createFileRoute } from '@tanstack/react-router'
import { DEFAULT_PAGE_SIZE } from '@vp/core'
import { formatDateYYYYMMDDHHmmss } from '@vp/utils'
import { usePagination } from 'alova/client'
import {
  Descriptions,
  Input,
  Modal,
  Select,
  Tag,
} from 'antd'
import { useCallback, useState } from 'react'
import { LoginLogApi } from '~/api/business/sysLoginLog'

export const Route = createFileRoute('/_app/logger/login/log')({
  staleTime: 1000 * 60 * 2,
  component: LoginLogManagement,
})

function successTag(success: boolean) {
  if (success) {
    return <Tag color="success">成功</Tag>
  }
  return <Tag color="error">失败</Tag>
}

function DetailText({ children }: { children?: string | number | null }) {
  return (
    <span style={{ overflowWrap: 'anywhere', wordBreak: 'break-word' }}>
      {children === undefined || children === null || children === '' ? '-' : children}
    </span>
  )
}

function operatorName(record?: SysLoginLog | null) {
  if (record?.sysUser?.username) {
    return record.sysUser.nickname
      ? `${record.sysUser.username} (${record.sysUser.nickname})`
      : record.sysUser.username
  }
  return record?.sysUserID ? String(record.sysUserID) : '-'
}

function operationPlatform(record?: SysLoginLog | null) {
  const platform = [
    record?.osName,
    record?.browserName,
  ].filter(Boolean).join(' ')

  return platform || '-'
}

function LoginLogManagement() {
  const [detailOpen, setDetailOpen] = useState(false)
  const [detailData, setDetailData] = useState<SysLoginLog | null>(null)
  const [searchText, setSearchText] = useState('')
  const [successFilter, setSuccessFilter] = useState<boolean | undefined>()

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
            { username__icontains: keyword },
            { login_ip__icontains: keyword },
          ],
        })
      }
      if (successFilter !== undefined) {
        filters.push({ success: successFilter })
      }

      return LoginLogApi.list({
        page: nextPage,
        pageSize: nextPageSize,
        orderBy: 'id desc',
        query: filters.length > 0 ? JSON.stringify({ $and: filters }) : undefined,
      })
    },
    {
      initialData: {
        total: 0,
        items: [],
      },
      initialPage: 1,
      initialPageSize: DEFAULT_PAGE_SIZE,
      total: response => response.data?.total ?? 0,
      data: response => response.data?.items ?? [],
      watchingStates: [searchText, successFilter],
      debounce: [500, 0],
    },
  )

  const openDetail = useCallback(async (record: SysLoginLog) => {
    try {
      const res = await LoginLogApi.detail({ id: record.id })
      if (res.data) {
        setDetailData(res.data)
        setDetailOpen(true)
      }
    }
    catch {
      // ignore
    }
  }, [])

  const columns: ProColumns<SysLoginLog>[] = [
    {
      title: '序号',
      dataIndex: 'index',
      width: 60,
      render: (_, __, index) => (page - 1) * pageSize + index + 1,
    },
    {
      title: '登录账号',
      dataIndex: 'username',
      width: 140,
      ellipsis: true,
      render: (_, record) => record.username || '-',
    },
    {
      title: '关联用户',
      dataIndex: ['sysUser', 'username'],
      width: 160,
      ellipsis: true,
      render: (_, record) => operatorName(record),
    },
    {
      title: '登录IP',
      dataIndex: 'loginIP',
      width: 140,
      ellipsis: true,
      render: (_, record) => record.loginIP || '-',
    },
    {
      title: '登录平台',
      dataIndex: 'operationPlatform',
      width: 180,
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
      title: '状态码',
      dataIndex: 'statusCode',
      width: 90,
    },
    {
      title: '结果',
      dataIndex: 'success',
      width: 80,
      render: (_, record) => successTag(record.success),
    },
    {
      title: '登录时间',
      dataIndex: 'loginTime',
      width: 170,
      render: (_, record) => formatDateYYYYMMDDHHmmss(record.loginTime || record.createdAt),
    },
    {
      title: '操作',
      valueType: 'option',
      width: 80,
      fixed: 'right',
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
      <ProTable<SysLoginLog>
        rowKey="id"
        search={false}
        columns={columns}
        dataSource={data}
        loading={loading}
        headerTitle="登录日志"
        scroll={{ x: 1300 }}
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
        toolBarRender={() => [
          <Input.Search
            key="search"
            placeholder="搜索账号或登录IP"
            allowClear
            value={searchText}
            onChange={(event) => {
              setSearchText(event.target.value)
            }}
            onSearch={setSearchText}
            style={{ width: 240 }}
          />,
          <Select
            key="success"
            allowClear
            placeholder="登录结果"
            value={successFilter}
            style={{ width: 120 }}
            options={[
              { label: '成功', value: true },
              { label: '失败', value: false },
            ]}
            onChange={setSuccessFilter}
          />,
        ]}
      />
      <Modal
        title="登录日志详情"
        open={detailOpen}
        onCancel={() => {
          setDetailOpen(false)
        }}
        footer={null}
        width={920}
        style={{ top: 48 }}
        styles={{
          body: {
            maxHeight: 'calc(100vh - 160px)',
            overflowY: 'auto',
          },
        }}
      >
        {detailData && (
          <Descriptions column={2} bordered size="small">
            <Descriptions.Item label="ID">{detailData.id}</Descriptions.Item>
            <Descriptions.Item label="登录账号"><DetailText>{detailData.username}</DetailText></Descriptions.Item>
            <Descriptions.Item label="关联用户">{operatorName(detailData)}</Descriptions.Item>
            <Descriptions.Item label="用户ID"><DetailText>{detailData.sysUserID}</DetailText></Descriptions.Item>
            <Descriptions.Item label="登录IP"><DetailText>{detailData.loginIP}</DetailText></Descriptions.Item>
            <Descriptions.Item label="状态码"><DetailText>{detailData.statusCode}</DetailText></Descriptions.Item>
            <Descriptions.Item label="结果">{successTag(detailData.success)}</Descriptions.Item>
            <Descriptions.Item label="登录时间"><DetailText>{formatDateYYYYMMDDHHmmss(detailData.loginTime || detailData.createdAt)}</DetailText></Descriptions.Item>
            <Descriptions.Item label="失败原因" span={2}><DetailText>{detailData.reason}</DetailText></Descriptions.Item>
            <Descriptions.Item label="地理位置" span={2}><DetailText>{detailData.location}</DetailText></Descriptions.Item>
            <Descriptions.Item label="浏览器" span={2}>
              <DetailText>
                {[detailData.browserName, detailData.browserVersion].filter(Boolean).join(' ')}
              </DetailText>
            </Descriptions.Item>
            <Descriptions.Item label="操作系统" span={2}>
              <DetailText>
                {[detailData.osName, detailData.osVersion].filter(Boolean).join(' ')}
              </DetailText>
            </Descriptions.Item>
            <Descriptions.Item label="客户端" span={2}>
              <DetailText>
                {[detailData.clientName, detailData.clientID ? `(${detailData.clientID})` : ''].filter(Boolean).join(' ')}
              </DetailText>
            </Descriptions.Item>
            <Descriptions.Item label="User-Agent" span={2}><DetailText>{detailData.userAgent}</DetailText></Descriptions.Item>
            <Descriptions.Item label="创建时间"><DetailText>{formatDateYYYYMMDDHHmmss(detailData.createdAt)}</DetailText></Descriptions.Item>
            <Descriptions.Item label="登录MAC"><DetailText>{detailData.loginMAC}</DetailText></Descriptions.Item>
          </Descriptions>
        )}
      </Modal>
    </>
  )
}
