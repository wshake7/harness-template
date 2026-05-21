import type { ProColumns } from '@ant-design/pro-components'
import type { Key } from 'react'
import type { ResourceApi } from '~/api/business/sysResourceApi'
import type { ResourceMenu, ResourceMenuType } from '~/api/business/sysResourceMenu'
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
  Segmented,
  Select,
  Space,
  Switch,
  Tag,
} from 'antd'
import { useEffect, useMemo, useState } from 'react'
import z from 'zod'
import { ResourceApiApi } from '~/api/business/sysResourceApi'
import { ResourceMenuApi } from '~/api/business/sysResourceMenu'
import AntIconPicker from '~/components/common/antIconPicker'
import { useDictMatch } from '~/hooks/useDictMatch'
import { AntIcon } from '~/utils/antIcons'
import { gMessage } from '~/utils/message'
import { useZodForm } from '~/utils/zod'

export const Route = createFileRoute('/_app/system/resource/menu')({
  staleTime: 1000 * 60 * 2,
  component: ResourceMenuManagement,
})

const resourceMenuTypeLabels: Record<ResourceMenuType, string> = {
  CATALOG: '目录',
  MENU: '菜单',
  BUTTON: '按钮',
  EMBEDDED: '内嵌',
  LINK: '外链',
}

const resourceMenuTypeValues = Object.keys(resourceMenuTypeLabels) as ResourceMenuType[]

const parentMenuTypes: Record<ResourceMenuType, ResourceMenuType[]> = {
  CATALOG: ['CATALOG'],
  MENU: ['CATALOG'],
  BUTTON: ['MENU'],
  EMBEDDED: ['MENU'],
  LINK: ['MENU'],
}

const ResourceMenuFormSchema = z.object({
  parentID: z.number().optional(),
  menuType: z.enum(['CATALOG', 'MENU', 'BUTTON', 'EMBEDDED', 'LINK']),
  path: z.string().optional(),
  redirect: z.string().optional(),
  alias: z.string().optional(),
  name: z.string(),
  component: z.string().optional(),
  icon: z.string().optional(),
  sortOrder: z.number().default(0),
  hidden: z.boolean(),
  isEnabled: z.boolean(),
  authorities: z.string().optional(),
  apiIDs: z.array(z.number()).default([]),
  remark: z.string().optional(),
})

type ResourceMenuFormValues = z.infer<typeof ResourceMenuFormSchema>

const defaultFormValues = ResourceMenuFormSchema.parse({
  ...ResourceMenuFormSchema.partial().parse({}),
  parentID: 0,
  menuType: 'CATALOG',
  path: '',
  redirect: '',
  alias: '',
  name: '',
  component: '',
  icon: '',
  sortOrder: 0,
  hidden: false,
  isEnabled: true,
  authorities: '',
  apiIDs: [],
  remark: '',
})

const enabledStatusValue = (isEnabled: boolean) => isEnabled ? '1' : '0'
const fallbackEnabledStatusLabel = (isEnabled: boolean) => isEnabled ? '启用' : '停用'
const fallbackResourceMenuTypeLabel = (menuType: ResourceMenuType) => resourceMenuTypeLabels[menuType]

const ResourceMenuSubmitSchema = ResourceMenuFormSchema.superRefine((values, ctx) => {
  if (!values.name.trim()) {
    ctx.addIssue({
      code: 'custom',
      path: ['name'],
      message: '请输入菜单标题',
    })
  }
  if (isPathType(values.menuType) && !values.path?.trim()) {
    ctx.addIssue({
      code: 'custom',
      path: ['path'],
      message: '请输入路由地址',
    })
  }
  if (values.menuType === 'MENU' && !values.component?.trim()) {
    ctx.addIssue({
      code: 'custom',
      path: ['component'],
      message: '请输入组件路径',
    })
  }
})

function isIconType(type: ResourceMenuType) {
  return type === 'CATALOG' || type === 'MENU'
}

function isAuthorityType(type: ResourceMenuType) {
  return type === 'BUTTON' || type === 'EMBEDDED' || type === 'LINK'
}

function canAssociateApi(type: ResourceMenuType) {
  return type === 'MENU' || type === 'BUTTON'
}

function isPathType(type: ResourceMenuType) {
  return type === 'CATALOG' || type === 'MENU' || type === 'EMBEDDED' || type === 'LINK'
}

function canHaveChildren(type: ResourceMenuType) {
  return type === 'CATALOG' || type === 'MENU'
}

function canUseRootParent(type: ResourceMenuType) {
  return type === 'CATALOG'
}

function findMenuByID(items: ResourceMenu[], id: number): ResourceMenu | undefined {
  for (const item of items) {
    if (item.id === id) {
      return item
    }
    const child = findMenuByID(item.children ?? [], id)
    if (child) {
      return child
    }
  }
  return undefined
}

function isValidParentForType(type: ResourceMenuType, parentID?: number, items: ResourceMenu[] = []) {
  if (!parentID) {
    return canUseRootParent(type)
  }
  const parent = findMenuByID(items, parentID)
  return parent ? parentMenuTypes[type].includes(parent.menuType) : false
}

function defaultChildType(parent?: ResourceMenu): ResourceMenuType {
  if (!parent) {
    return 'CATALOG'
  }
  return parent.menuType === 'MENU' ? 'BUTTON' : 'MENU'
}

function flattenMenuOptions(
  items: ResourceMenu[],
  allowedTypes: ResourceMenuType[],
  editingID?: number,
  level = 0,
): { label: string, value: number }[] {
  const options: { label: string, value: number }[] = []
  for (const item of items) {
    if (item.id !== editingID && allowedTypes.includes(item.menuType)) {
      options.push({
        label: `${'  '.repeat(level)}${item.name || item.path || item.id}`,
        value: item.id,
      })
    }
    options.push(...flattenMenuOptions(item.children ?? [], allowedTypes, editingID, level + 1))
  }
  return options
}

function collectRowKeys(items: ResourceMenu[]) {
  const keys: Key[] = []
  const walk = (nodes: ResourceMenu[]) => {
    for (const node of nodes) {
      keys.push(node.id)
      if (node.children?.length) {
        walk(node.children)
      }
    }
  }
  walk(items)
  return keys
}

function authoritiesToText(value: unknown) {
  if (Array.isArray(value)) {
    return value.map(String).filter(Boolean).join(', ')
  }
  return ''
}

function splitAuthorities(value: string) {
  return value
    .split(',')
    .map(item => item.trim())
    .filter(Boolean)
}

function toFormValues(record: ResourceMenu): ResourceMenuFormValues {
  return ResourceMenuFormSchema.parse({
    ...defaultFormValues,
    ...ResourceMenuFormSchema.partial().parse({
      parentID: record.parentID ?? 0,
      menuType: record.menuType,
      path: record.path ?? '',
      redirect: record.redirect ?? '',
      alias: record.alias ?? '',
      name: record.name ?? '',
      component: record.component ?? '',
      icon: String(record.metadata?.icon ?? ''),
      sortOrder: record.sortOrder ?? Number(record.metadata?.order ?? 0),
      hidden: Boolean(record.metadata?.hidden),
      isEnabled: Boolean(record.isEnabled),
      authorities: authoritiesToText(record.metadata?.authorities),
      apiIDs: record.apiIDs ?? [],
      remark: record.remark ?? '',
    }),
  })
}

function toPayload(values: ResourceMenuFormValues) {
  const metadata: Record<string, unknown> = {}
  if (isIconType(values.menuType)) {
    metadata.icon = values.icon ?? ''
    metadata.order = values.sortOrder ?? 0
    metadata.hidden = values.hidden ?? false
  }
  if (isAuthorityType(values.menuType)) {
    metadata.authorities = splitAuthorities(values.authorities ?? '')
  }

  return {
    parentID: values.parentID ?? 0,
    menuType: values.menuType,
    path: isPathType(values.menuType) ? values.path ?? '' : '',
    redirect: values.redirect ?? '',
    alias: values.alias ?? '',
    name: values.name ?? '',
    component: values.menuType === 'MENU' ? values.component ?? '' : '',
    metadata,
    apiIDs: canAssociateApi(values.menuType) ? values.apiIDs ?? [] : [],
    sortOrder: values.sortOrder ?? 0,
    isEnabled: values.isEnabled ?? true,
    remark: values.remark ?? '',
  }
}

function ResourceMenuManagement() {
  const [parentRows, setParentRows] = useState<ResourceMenu[]>([])
  const [expandedRowKeys, setExpandedRowKeys] = useState<Key[]>([])
  const [drawerOpen, setDrawerOpen] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [editing, setEditing] = useState<ResourceMenu>()
  const [apiItems, setApiItems] = useState<ResourceApi[]>([])
  const [searchText, setSearchText] = useState('')
  const [enabledFilter, setEnabledFilter] = useState<boolean | undefined>()
  const enabledStatus = useDictMatch(DictCode.SYS_IS_ENABLED_DICT_CODE)
  const resourceMenuType = useDictMatch(DictCode.SYS_RESOURCE_MENU_TYPE_DICT_CODE)

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
      if (searchText.trim()) {
        filters.push({ name__icontains: searchText.trim() })
      }
      if (enabledFilter !== undefined) {
        filters.push({ isEnabled: enabledFilter })
      }

      return ResourceMenuApi.menuList({
        page: nextPage,
        pageSize: nextPageSize,
        orderBy: 'sort_order asc,id asc',
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

  const menuTypeOptions = useMemo(() => {
    const options = resourceMenuType.entries
      .filter(entry => resourceMenuTypeValues.includes(entry.entryValue as ResourceMenuType))
      .map(entry => ({
        label: resourceMenuType.getLabel(entry.entryValue, entry.entryLabel),
        value: entry.entryValue as ResourceMenuType,
      }))

    return options.length > 0
      ? options
      : resourceMenuTypeValues.map(value => ({
          label: fallbackResourceMenuTypeLabel(value),
          value,
        }))
  }, [resourceMenuType])

  const apiOptions = useMemo(() => apiItems.map(item => ({
    label: `${item.method} ${item.path}${item.module ? ` (${item.module})` : ''}`,
    value: item.id,
  })), [apiItems])

  const refreshApiItems = async () => {
    const res = await ResourceApiApi.apiList({
      noPaging: true,
      orderBy: 'module asc,sort_order asc,id asc',
    }).send()
    setApiItems(res.data?.items ?? [])
  }

  const refreshParentRows = async () => {
    const res = await ResourceMenuApi.menuList({
      noPaging: true,
      orderBy: 'sort_order asc,id asc',
    }).send()
    setParentRows(res.data?.items ?? [])
  }

  const refreshData = async () => {
    await Promise.all([
      send(),
      refreshParentRows(),
    ])
  }

  const closeDrawer = () => {
    setDrawerOpen(false)
  }

  const { form, rules, onFinish } = useZodForm<ResourceMenuFormValues>({
    schema: ResourceMenuSubmitSchema,
    async onSubmit(values) {
      if (!values) {
        gMessage.error('请填写完整信息')
        return
      }

      setSubmitting(true)
      try {
        const payload = toPayload(values)
        if (editing) {
          await ResourceMenuApi.menuUpdate({
            id: editing.id,
            ...payload,
          })
        }
        else {
          await ResourceMenuApi.menuCreate(payload)
        }
        gMessage.success('保存成功')
        closeDrawer()
        await refreshData()
      }
      catch {
        gMessage.error('保存失败')
      }
      finally {
        setSubmitting(false)
      }
    },
  })

  const menuType = Form.useWatch('menuType', form) ?? defaultFormValues.menuType

  const parentOptions = useMemo(() => [
    ...(canUseRootParent(menuType) ? [{ label: '根节点', value: 0 }] : []),
    ...flattenMenuOptions(parentRows, parentMenuTypes[menuType], editing?.id),
  ], [editing?.id, menuType, parentRows])

  const openCreateForm = (parentID = 0) => {
    const parent = parentID ? findMenuByID(parentRows, parentID) : undefined
    const nextFormValues = {
      ...defaultFormValues,
      parentID,
      menuType: defaultChildType(parent),
    }
    setEditing(undefined)
    form.resetFields()
    form.setFieldsValue(nextFormValues)
    refreshApiItems().catch(() => undefined)
    setDrawerOpen(true)
  }

  const openEditForm = (record: ResourceMenu) => {
    const nextFormValues = toFormValues(record)
    setEditing(record)
    form.resetFields()
    form.setFieldsValue(nextFormValues)
    refreshApiItems().catch(() => undefined)
    setDrawerOpen(true)
  }

  const reload = () => {
    refreshData()
  }

  useEffect(() => {
    if (!drawerOpen) {
      return
    }
    const parentID = form.getFieldValue('parentID')
    if (!isValidParentForType(menuType, parentID, parentRows)) {
      form.setFieldValue('parentID', canUseRootParent(menuType) ? 0 : undefined)
    }
    form.validateFields(['parentID']).catch(() => undefined)
  }, [drawerOpen, form, menuType, parentRows])

  useEffect(() => {
    refreshParentRows().catch(() => undefined)
  }, [])

  const columns: ProColumns<ResourceMenu>[] = [
    {
      title: '标题',
      dataIndex: 'name',
      width: 220,
      render: (_, record) => (
        <Space size={6}>
          <AntIcon name={String(record.metadata?.icon ?? '')} />
          <span>{record.name}</span>
          {record.metadata?.hidden ? <Tag>隐藏</Tag> : null}
        </Space>
      ),
    },
    {
      title: '类型',
      dataIndex: 'menuType',
      width: 90,
      search: false,
      render: (_, record) =>
        resourceMenuType.renderLabel(record.menuType, fallbackResourceMenuTypeLabel(record.menuType)),
    },
    {
      title: '权限标识',
      dataIndex: ['metadata', 'authorities'],
      search: false,
      width: 280,
      render: (_, record) => (
        <Space wrap size={4}>
          {(record.metadata?.authorities ?? []).map(item => (
            <Tag key={item}>{item}</Tag>
          ))}
        </Space>
      ),
    },
    {
      title: '路由地址',
      dataIndex: 'path',
      search: false,
      ellipsis: true,
    },
    {
      title: '组件路径',
      dataIndex: 'component',
      search: false,
      ellipsis: true,
      width: 220,
    },
    {
      title: '状态',
      dataIndex: 'isEnabled',
      width: 90,
      valueType: 'select',
      render: (_, record) =>
        enabledStatus.renderLabel(enabledStatusValue(record.isEnabled), fallbackEnabledStatusLabel(record.isEnabled)),
    },
    {
      title: '排序',
      dataIndex: 'sortOrder',
      search: false,
      width: 80,
    },
    {
      title: '操作',
      valueType: 'option',
      width: 180,
      render: (_, record) => (
        <Space>
          <Button type="link" size="small" disabled={!canHaveChildren(record.menuType)} onClick={() => openCreateForm(record.id)}>
            新增
          </Button>
          <Button type="link" size="small" disabled={record.canWrite === false} onClick={() => openEditForm(record)}>
            编辑
          </Button>

          <Popconfirm
            title="确认删除该菜单？"
            onConfirm={async () => {
              try {
                await ResourceMenuApi.menuDel({ ids: [record.id] })
                gMessage.success('删除成功')
                reload()
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
      <ProTable<ResourceMenu>
        rowKey="id"
        columns={columns}
        dataSource={data}
        loading={loading}
        headerTitle="菜单管理"
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
        expandable={{
          expandedRowKeys,
          onExpandedRowsChange: keys => setExpandedRowKeys([...keys]),
        }}
        options={{
          reload,
        }}
        toolBarRender={() => [
          <Button key="create" type="primary" onClick={() => openCreateForm()}>
            创建菜单
          </Button>,
          <Button key="expand" onClick={() => setExpandedRowKeys(collectRowKeys(data))}>
            展开全部
          </Button>,
          <Button key="collapse" onClick={() => setExpandedRowKeys([])}>
            折叠全部
          </Button>,
          <Input.Search
            key="search"
            placeholder="搜索菜单标题"
            allowClear
            value={searchText}
            onChange={(event) => {
              setSearchText(event.target.value)
            }}
            onSearch={(value) => {
              setSearchText(value)
            }}
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
        title={editing ? '编辑菜单' : '创建菜单'}
        open={drawerOpen}
        size={760}
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
        <Form<ResourceMenuFormValues>
          form={form}
          layout="horizontal"
          labelCol={{ span: 5 }}
          wrapperCol={{ span: 19 }}
          onFinish={onFinish}
        >
          <Form.Item
            name="menuType"
            label="类型"
            rules={rules}
          >
            <Segmented options={menuTypeOptions} />
          </Form.Item>
          <ProFormText
            name="name"
            label="菜单标题"
            placeholder="请输入"
            rules={rules}
          />
          <ProFormSelect
            name="parentID"
            label="上级菜单"
            placeholder="请选择"
            options={parentOptions}
            rules={[
              ...rules,
              {
                validator: async (_, value) => {
                  if (!isValidParentForType(menuType, value, parentRows)) {
                    throw new Error(
                      menuType === 'CATALOG'
                        ? '上级菜单只能选择目录或根节点'
                        : menuType === 'MENU'
                          ? '上级菜单只能选择目录'
                          : '上级菜单只能选择菜单',
                    )
                  }
                },
              },
            ]}
            fieldProps={{ showSearch: true, optionFilterProp: 'label' }}
          />
          <ProFormDigit
            name="sortOrder"
            label="排序"
            placeholder="请输入"
            rules={rules}
            fieldProps={{ precision: 0 }}
          />
          {isIconType(menuType)
            ? (
                <Form.Item name="icon" label="图标" rules={rules}>
                  <AntIconPicker />
                </Form.Item>
              )
            : null}
          {isPathType(menuType)
            ? (
                <ProFormText
                  name="path"
                  label="路由地址"
                  placeholder="请输入"
                  rules={rules}
                />
              )
            : null}
          {menuType === 'MENU'
            ? (
                <ProFormText
                  name="component"
                  label="组件路径"
                  placeholder="BasicLayout"
                  rules={rules}
                />
              )
            : null}
          {isAuthorityType(menuType)
            ? (
                <ProFormText
                  name="authorities"
                  label="权限标识"
                  placeholder="多个权限用逗号隔开，如：sys:user:list, sys:user:create"
                  rules={rules}
                />
              )
            : null}
          {canAssociateApi(menuType)
            ? (
                <ProFormSelect
                  name="apiIDs"
                  label="关联 API"
                  placeholder="请选择关联 API"
                  options={apiOptions}
                  rules={rules}
                  fieldProps={{
                    mode: 'multiple',
                    showSearch: true,
                    optionFilterProp: 'label',
                  }}
                />
              )
            : null}
          {menuType === 'LINK'
            ? (
                <ProFormText name="redirect" label="重定向" placeholder="请输入" rules={rules} />
              )
            : null}
          <Form.Item name="isEnabled" label="状态" valuePropName="checked" rules={rules}>
            <Switch
              checkedChildren={enabledStatus.getLabel(enabledStatusValue(true), fallbackEnabledStatusLabel(true))}
              unCheckedChildren={enabledStatus.getLabel(enabledStatusValue(false), fallbackEnabledStatusLabel(false))}
            />
          </Form.Item>
          {isIconType(menuType)
            ? (
                <Form.Item name="hidden" valuePropName="checked" wrapperCol={{ offset: 5, span: 19 }} rules={rules}>
                  <Switch checkedChildren="隐藏菜单" unCheckedChildren="显示菜单" />
                </Form.Item>
              )
            : null}
          <ProFormText name="alias" label="别名" placeholder="请输入" rules={rules} />
          <ProFormTextArea name="remark" label="备注" rules={rules} fieldProps={{ maxLength: 255, showCount: true }} />
        </Form>
      </Drawer>
    </>
  )
}
