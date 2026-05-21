import type { ProColumns } from '@ant-design/pro-components'
import type { SysUser } from '~/api/business/sysUser'
import {
  ProFormText,
  ProFormTextArea,
  ProTable,
} from '@ant-design/pro-components'
import { createFileRoute } from '@tanstack/react-router'
import { DEFAULT_PAGE_SIZE } from '@vp/core'
import { formatDateYYYYMMDDHHmmss } from '@vp/utils'
import { usePagination } from 'alova/client'
import {
  Button,
  Drawer,
  Form,
  Input,
  Popconfirm,
  Select,
  Space,
  Switch,
} from 'antd'
import { useMemo, useState } from 'react'
import z from 'zod'
import { SysUserApi } from '~/api/business/sysUser'
import { useDictMatch } from '~/hooks/useDictMatch'
import { gMessage } from '~/utils/message'

export const Route = createFileRoute('/_app/account/user')({
  staleTime: 1000 * 60 * 2,
  component: UserManagement,
})

const UserFormSchema = z.object({
  username: z.string(),
  nickname: z.string().optional(),
  password: z.string().optional(),
  languageCode: z.string().optional(),
  isEnabled: z.boolean(),
  remark: z.string().optional(),
})

type UserFormValues = z.infer<typeof UserFormSchema>

const defaultFormValues = UserFormSchema.parse({
  ...UserFormSchema.partial().parse({}),
  username: '',
  nickname: '',
  password: '',
  languageCode: '',
  isEnabled: true,
  remark: '',
})

const enabledStatusValue = (isEnabled: boolean) => isEnabled ? '1' : '0'
const fallbackEnabledStatusLabel = (isEnabled: boolean) => isEnabled ? '启用' : '停用'

const UserSubmitSchema = UserFormSchema.superRefine((values, ctx) => {
  if (!values.username.trim()) {
    ctx.addIssue({
      code: 'custom',
      path: ['username'],
      message: '请输入用户名',
    })
  }
})

function toFormValues(record: SysUser): UserFormValues {
  return UserFormSchema.parse({
    ...defaultFormValues,
    ...UserFormSchema.partial().parse({
      username: record.username ?? '',
      nickname: record.nickname ?? '',
      password: '',
      languageCode: record.languageCode ?? '',
      isEnabled: Boolean(record.isEnabled),
      remark: record.remark ?? '',
    }),
  })
}

function UserManagement() {
  const [drawerOpen, setDrawerOpen] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [editing, setEditing] = useState<SysUser>()
  const [searchText, setSearchText] = useState('')
  const [enabledFilter, setEnabledFilter] = useState<boolean | undefined>()
  const enabledStatus = useDictMatch(DictCode.SYS_IS_ENABLED_DICT_CODE)

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
            { nickname__icontains: keyword },
          ],
        })
      }
      if (enabledFilter !== undefined) {
        filters.push({ isEnabled: enabledFilter })
      }

      return SysUserApi.apiList({
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
      watchingStates: [searchText, enabledFilter],
      debounce: [500, 0],
    },
  )

  const enabledStatusOptions = useMemo(() => {
    const options = enabledStatus.entries
      .filter(entry => entry.entryValue === '1' || entry.entryValue === '0')
      .map(entry => ({
        label: enabledStatus.getLabel(entry.entryValue, entry.entryLabel),
        value: entry.entryValue === '1',
      }))

    return options.length > 0
      ? options
      : [
          { label: fallbackEnabledStatusLabel(true), value: true },
          { label: fallbackEnabledStatusLabel(false), value: false },
        ]
  }, [enabledStatus])

  const closeDrawer = () => {
    setDrawerOpen(false)
  }

  const { form, rules, onFinish } = useZodForm<UserFormValues>({
    schema: UserSubmitSchema,
    async onSubmit(values) {
      if (!values) {
        gMessage.error('请填写完整信息')
        return
      }

      const username = values.username.trim()
      const nickname = values.nickname?.trim()
      const languageCode = values.languageCode?.trim() ?? ''
      const remark = values.remark?.trim() ?? ''
      const password = values.password?.trim() ?? ''

      if (!editing && !password) {
        form.setFields([
          {
            name: 'password',
            errors: ['请输入初始密码'],
          },
        ])
        return
      }

      setSubmitting(true)
      try {
        if (editing) {
          await SysUserApi.apiUpdate({
            id: editing.id,
            username,
            nickname,
            languageCode,
            isEnabled: values.isEnabled ?? editing.isEnabled,
            remark,
          })
        }
        else {
          await SysUserApi.apiCreate({
            username,
            nickname,
            password,
            languageCode,
            isEnabled: values.isEnabled ?? true,
            remark,
          })
        }
        gMessage.success('保存成功')
        closeDrawer()
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

  const saveUser = async () => {
    try {
      const values = await form.validateFields()
      await onFinish?.({
        ...defaultFormValues,
        ...values,
        isEnabled: values.isEnabled ?? editing?.isEnabled ?? true,
      })
    }
    catch {
      gMessage.error('请检查表单信息')
    }
  }

  const openCreateForm = () => {
    setEditing(undefined)
    form.resetFields()
    form.setFieldsValue(defaultFormValues)
    setDrawerOpen(true)
  }

  const openEditForm = (record: SysUser) => {
    setEditing(record)
    form.resetFields()
    form.setFieldsValue(toFormValues(record))
    setDrawerOpen(true)
  }

  const toggleEnabled = async (record: SysUser) => {
    try {
      await SysUserApi.apiUpdate({
        id: record.id,
        isEnabled: !record.isEnabled,
      })
      gMessage.success(`${enabledStatus.getLabel(enabledStatusValue(!record.isEnabled), fallbackEnabledStatusLabel(!record.isEnabled))}成功`)
      await send()
    }
    catch {
      gMessage.error(`${enabledStatus.getLabel(enabledStatusValue(!record.isEnabled), fallbackEnabledStatusLabel(!record.isEnabled))}失败`)
    }
  }

  const columns: ProColumns<SysUser>[] = [
    {
      title: '用户名',
      dataIndex: 'username',
      width: 140,
      ellipsis: true,
    },
    {
      title: '昵称',
      dataIndex: 'nickname',
      width: 120,
      ellipsis: true,
    },
    {
      title: '语言代码',
      dataIndex: 'languageCode',
      width: 110,
      ellipsis: true,
      render: (_, record) => record.languageCode || '-',
    },
    {
      title: '状态',
      dataIndex: 'isEnabled',
      width: 84,
      render: (_, record) =>
        enabledStatus.renderLabel(enabledStatusValue(record.isEnabled), fallbackEnabledStatusLabel(record.isEnabled)),
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      width: 150,
      ellipsis: true,
      render: (_, record) => formatDateYYYYMMDDHHmmss(record.createdAt),
    },
    {
      title: '备注',
      dataIndex: 'remark',
      ellipsis: true,
      render: (_, record) => record.remark || '-',
    },
    {
      title: '操作',
      valueType: 'option',
      width: 170,
      fixed: 'right',
      render: (_, record) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            disabled={record.canWrite === false}
            onClick={() => toggleEnabled(record)}
          >
            {enabledStatus.getLabel(enabledStatusValue(!record.isEnabled), fallbackEnabledStatusLabel(!record.isEnabled))}
          </Button>
          <Button type="link" size="small" disabled={record.canWrite === false} onClick={() => openEditForm(record)}>
            编辑
          </Button>
          <Popconfirm
            title="确认删除该用户？"
            onConfirm={async () => {
              try {
                await SysUserApi.apiDel({ ids: [record.id] })
                gMessage.success('删除成功')
                await send()
              }
              catch {
                gMessage.error('删除失败')
              }
            }}
          >
            <Button type="link" size="small" disabled={record.canDelete === false}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <>
      <ProTable<SysUser>
        rowKey="id"
        columns={columns}
        dataSource={data}
        loading={loading}
        headerTitle="用户管理"
        search={false}
        scroll={{ x: 1200 }}
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
          <Button key="create" type="primary" onClick={openCreateForm}>
            创建用户
          </Button>,
          <Input.Search
            key="search"
            placeholder="搜索用户名或昵称"
            allowClear
            value={searchText}
            onChange={(event) => {
              setSearchText(event.target.value)
            }}
            onSearch={setSearchText}
            style={{ width: 240 }}
          />,
          <Select
            key="enabled"
            allowClear
            placeholder="状态"
            value={enabledFilter}
            style={{ width: 120 }}
            options={enabledStatusOptions}
            onChange={setEnabledFilter}
          />,
        ]}
      />

      <Drawer
        title={editing ? '编辑用户' : '创建用户'}
        open={drawerOpen}
        size={640}
        destroyOnHidden
        forceRender
        onClose={closeDrawer}
        extra={(
          <Space>
            <Button onClick={closeDrawer}>取消</Button>
            <Button type="primary" loading={submitting} onClick={saveUser}>
              保存
            </Button>
          </Space>
        )}
      >
        <Form<UserFormValues>
          form={form}
          layout="horizontal"
          labelCol={{ span: 5 }}
          wrapperCol={{ span: 19 }}
          onFinish={onFinish}
        >
          <ProFormText
            name="username"
            label="用户名"
            placeholder="请输入用户名"
            rules={rules}
            fieldProps={{ maxLength: 64, showCount: true }}
          />
          <ProFormText
            name="nickname"
            label="昵称"
            placeholder="请输入昵称"
            rules={rules}
            fieldProps={{ maxLength: 64, showCount: true }}
          />
          {!editing
            ? (
                <ProFormText
                  name="password"
                  label="初始密码"
                  placeholder="请输入初始密码"
                  rules={rules}
                  fieldProps={{
                    type: 'password',
                    autoComplete: 'new-password',
                    maxLength: 255,
                  }}
                />
              )
            : null}
          <ProFormText
            name="languageCode"
            label="语言代码"
            placeholder="例如：zh-CN"
            rules={rules}
            fieldProps={{ maxLength: 32, showCount: true }}
          />
          {!editing
            ? (
                <Form.Item name="isEnabled" label="状态" valuePropName="checked" rules={rules}>
                  <Switch
                    checkedChildren={enabledStatus.getLabel(enabledStatusValue(true), fallbackEnabledStatusLabel(true))}
                    unCheckedChildren={enabledStatus.getLabel(enabledStatusValue(false), fallbackEnabledStatusLabel(false))}
                  />
                </Form.Item>
              )
            : null}
          <ProFormTextArea
            name="remark"
            label="备注"
            rules={rules}
            fieldProps={{ maxLength: 255, showCount: true }}
          />
        </Form>
      </Drawer>
    </>
  )
}
