import type { ProColumns } from '@ant-design/pro-components'
import type { ResourceApi, ResourceApiMethod } from '~/api/business/sysResourceApi'
import {
  ProFormDigit,
  ProFormSelect,
  ProFormText,
  ProFormTextArea,
  ProTable,
} from '@ant-design/pro-components'
import { createFileRoute } from '@tanstack/react-router'
import { DEFAULT_PAGE_SIZE } from '@vp/core'
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
import { ResourceApiApi } from '~/api/business/sysResourceApi'
import { useDictMatch } from '~/hooks/useDictMatch'
import { gMessage } from '~/utils/message'

export const Route = createFileRoute('/_app/system/resource/api')({
  staleTime: 1000 * 60 * 2,
  component: ResourceApiManagement,
})

const resourceApiMethodValues: ResourceApiMethod[] = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'OPTIONS', 'HEAD']

const ResourceApiFormSchema = z.object({
  module: z.string().optional(),
  path: z.string(),
  method: z.enum(['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'OPTIONS', 'HEAD']),
  sortOrder: z.number().default(0),
  isEnabled: z.boolean(),
  remark: z.string().optional(),
})

type ResourceApiFormValues = z.infer<typeof ResourceApiFormSchema>

const defaultFormValues = ResourceApiFormSchema.parse({
  ...ResourceApiFormSchema.partial().parse({}),
  module: '',
  path: '',
  method: 'GET',
  sortOrder: 0,
  isEnabled: true,
  remark: '',
})

const enabledStatusValue = (isEnabled: boolean) => isEnabled ? '1' : '0'
const fallbackEnabledStatusLabel = (isEnabled: boolean) => isEnabled ? '启用' : '停用'

const ResourceApiSubmitSchema = ResourceApiFormSchema.superRefine((values, ctx) => {
  const path = values.path.trim()
  if (!path) {
    ctx.addIssue({
      code: 'custom',
      path: ['path'],
      message: '请输入接口路径',
    })
    return
  }
  if (!path.startsWith('/')) {
    ctx.addIssue({
      code: 'custom',
      path: ['path'],
      message: '接口路径必须以 / 开头',
    })
    return
  }
  const error = validatePathTemplate(path)
  if (error) {
    ctx.addIssue({
      code: 'custom',
      path: ['path'],
      message: error,
    })
  }
})

function validatePathTemplate(path: string) {
  for (const segment of path.split('/')) {
    if (!segment) {
      continue
    }
    const braceMatch = segment.match(/^\{(.+)\}$/)
    if (braceMatch) {
      if (!isValidParamName(braceMatch[1].trim())) {
        return '路径参数名只能包含字母、数字、下划线，且不能以数字开头'
      }
      continue
    }
    if (segment.includes('{') || segment.includes('}')) {
      return '路径参数请使用 {name} 或 :name 格式'
    }
    if (segment.startsWith(':') && !isValidParamName(segment.slice(1))) {
      return '路径参数名只能包含字母、数字、下划线，且不能以数字开头'
    }
  }
  return ''
}

function isValidParamName(name: string) {
  return /^[_A-Z]\w*$/i.test(name)
}

function normalizePathTemplate(path: string) {
  return path
    .trim()
    .split('/')
    .map((segment) => {
      const braceMatch = segment.match(/^\{(.+)\}$/)
      return braceMatch ? `:${braceMatch[1].trim()}` : segment
    })
    .join('/')
}

function toFormValues(record: ResourceApi): ResourceApiFormValues {
  return ResourceApiFormSchema.parse({
    ...defaultFormValues,
    ...ResourceApiFormSchema.partial().parse({
      module: record.module ?? '',
      path: record.path ?? '',
      method: record.method,
      sortOrder: record.sortOrder ?? 0,
      isEnabled: Boolean(record.isEnabled),
      remark: record.remark ?? '',
    }),
  })
}

function toPayload(values: ResourceApiFormValues) {
  return {
    module: values.module?.trim() ?? '',
    path: normalizePathTemplate(values.path),
    method: values.method,
    sortOrder: values.sortOrder ?? 0,
    isEnabled: values.isEnabled ?? true,
    remark: values.remark?.trim() ?? '',
  }
}

function ResourceApiManagement() {
  const [drawerOpen, setDrawerOpen] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [editing, setEditing] = useState<ResourceApi>()
  const [searchText, setSearchText] = useState('')
  const [methodFilter, setMethodFilter] = useState<ResourceApiMethod | undefined>()
  const [enabledFilter, setEnabledFilter] = useState<boolean | undefined>()
  const enabledStatus = useDictMatch(DictCode.SYS_IS_ENABLED_DICT_CODE)
  const resourceApiType = useDictMatch(DictCode.SYS_RESOURCE_API_TYPE_DICT_CODE)

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
            { path__icontains: keyword },
            { module__icontains: keyword },
          ],
        })
      }
      if (methodFilter) {
        filters.push({ method: methodFilter })
      }
      if (enabledFilter !== undefined) {
        filters.push({ isEnabled: enabledFilter })
      }

      return ResourceApiApi.apiList({
        page: nextPage,
        pageSize: nextPageSize,
        orderBy: 'sort_order asc,id desc',
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
      watchingStates: [searchText, methodFilter, enabledFilter],
      debounce: [500, 0, 0],
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

  const methodOptions = useMemo(() => {
    const options = resourceApiType.entries
      .filter(entry => resourceApiMethodValues.includes(entry.entryValue as ResourceApiMethod))
      .map(entry => ({
        label: resourceApiType.getLabel(entry.entryValue, entry.entryLabel),
        value: entry.entryValue as ResourceApiMethod,
      }))

    return options.length > 0
      ? options
      : resourceApiMethodValues.map(value => ({
          label: value,
          value,
        }))
  }, [resourceApiType])

  const closeDrawer = () => {
    setDrawerOpen(false)
  }

  const { form, rules, onFinish } = useZodForm<ResourceApiFormValues>({
    schema: ResourceApiSubmitSchema,
    async onSubmit(values) {
      if (!values) {
        gMessage.error('请填写完整信息')
        return
      }

      setSubmitting(true)
      try {
        const payload = toPayload(values)
        if (editing) {
          await ResourceApiApi.apiUpdate({
            id: editing.id,
            ...payload,
          })
        }
        else {
          await ResourceApiApi.apiCreate(payload)
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

  const openCreateForm = () => {
    setEditing(undefined)
    form.resetFields()
    form.setFieldsValue(defaultFormValues)
    setDrawerOpen(true)
  }

  const openEditForm = (record: ResourceApi) => {
    setEditing(record)
    form.resetFields()
    form.setFieldsValue(toFormValues(record))
    setDrawerOpen(true)
  }

  const toggleEnabled = async (record: ResourceApi) => {
    try {
      await ResourceApiApi.apiUpdate({
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

  const columns: ProColumns<ResourceApi>[] = [
    {
      title: '模块',
      dataIndex: 'module',
      width: 160,
      ellipsis: true,
      render: (_, record) => record.module || '-',
    },
    {
      title: '方法',
      dataIndex: 'method',
      width: 100,
      render: (_, record) => resourceApiType.renderLabel(record.method, record.method),
    },
    {
      title: '路径模板',
      dataIndex: 'path',
      ellipsis: true,
    },
    {
      title: '状态',
      dataIndex: 'isEnabled',
      width: 90,
      render: (_, record) =>
        enabledStatus.renderLabel(enabledStatusValue(record.isEnabled), fallbackEnabledStatusLabel(record.isEnabled)),
    },
    {
      title: '排序',
      dataIndex: 'sortOrder',
      width: 80,
    },
    {
      title: '备注',
      dataIndex: 'remark',
      ellipsis: true,
    },
    {
      title: '操作',
      valueType: 'option',
      width: 190,
      render: (_, record) => (
        <Space>
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
            title="确认删除该API资源？"
            onConfirm={async () => {
              try {
                await ResourceApiApi.apiDel({ ids: [record.id] })
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
      <ProTable<ResourceApi>
        rowKey="id"
        columns={columns}
        dataSource={data}
        loading={loading}
        headerTitle="API资源管理"
        search={false}
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
            创建API资源
          </Button>,
          <Input.Search
            key="search"
            placeholder="搜索路径或模块"
            allowClear
            value={searchText}
            onChange={(event) => {
              setSearchText(event.target.value)
            }}
            onSearch={setSearchText}
            style={{ width: 240 }}
          />,
          <Select
            key="method"
            allowClear
            placeholder="方法"
            value={methodFilter}
            style={{ width: 130 }}
            options={methodOptions}
            onChange={setMethodFilter}
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
        title={editing ? '编辑API资源' : '创建API资源'}
        open={drawerOpen}
        size={640}
        destroyOnHidden
        forceRender
        onClose={closeDrawer}
        extra={(
          <Space>
            <Button onClick={closeDrawer}>取消</Button>
            <Button type="primary" loading={submitting} onClick={() => form.submit()}>
              保存
            </Button>
          </Space>
        )}
      >
        <Form<ResourceApiFormValues>
          form={form}
          layout="horizontal"
          labelCol={{ span: 5 }}
          wrapperCol={{ span: 19 }}
          onFinish={onFinish}
        >
          <ProFormText
            name="module"
            label="模块"
            placeholder="例如：用户中心"
            rules={rules}
          />
          <ProFormSelect
            name="method"
            label="请求方法"
            placeholder="请选择"
            options={methodOptions}
            rules={rules}
          />
          <ProFormText
            name="path"
            label="路径模板"
            placeholder="/api/users/:id"
            tooltip="也可以输入 /api/users/{id}，保存后会统一为 /api/users/:id"
            rules={rules}
          />
          <ProFormDigit
            name="sortOrder"
            label="排序"
            min={0}
            rules={rules}
            fieldProps={{ precision: 0 }}
          />
          <Form.Item name="isEnabled" label="状态" valuePropName="checked" rules={rules}>
            <Switch
              checkedChildren={enabledStatus.getLabel(enabledStatusValue(true), fallbackEnabledStatusLabel(true))}
              unCheckedChildren={enabledStatus.getLabel(enabledStatusValue(false), fallbackEnabledStatusLabel(false))}
            />
          </Form.Item>
          <ProFormTextArea name="remark" label="备注" rules={rules} fieldProps={{ maxLength: 255, showCount: true }} />
        </Form>
      </Drawer>
    </>
  )
}
