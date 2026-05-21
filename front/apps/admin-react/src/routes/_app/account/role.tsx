import type { ProColumns } from '@ant-design/pro-components'
import type { Key } from 'react'
import type { ResourceApi } from '~/api/business/sysResourceApi'
import type { ResourceMenu } from '~/api/business/sysResourceMenu'
import type { SysRole } from '~/api/business/sysRole'
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
  Table,
  Tag,
  Tree,
  TreeSelect,
} from 'antd'
import { useMemo, useState } from 'react'
import z from 'zod'
import { ResourceApiApi } from '~/api/business/sysResourceApi'
import { ResourceMenuApi } from '~/api/business/sysResourceMenu'
import { RoleApi } from '~/api/business/sysRole'
import { useDictMatch } from '~/hooks/useDictMatch'
import { gMessage } from '~/utils/message'

export const Route = createFileRoute('/_app/account/role')({
  staleTime: 1000 * 60 * 2,
  component: RoleManagement,
})

const RoleFormSchema = z.object({
  name: z.string(),
  code: z.string(),
  parentID: z.number().optional(),
  isEnabled: z.boolean(),
  remark: z.string().optional(),
})

type RoleFormValues = z.infer<typeof RoleFormSchema>

const defaultFormValues = RoleFormSchema.parse({
  ...RoleFormSchema.partial().parse({}),
  name: '',
  code: '',
  parentID: undefined,
  isEnabled: true,
  remark: '',
})

const enabledStatusValue = (isEnabled: boolean) => isEnabled ? '1' : '0'
const fallbackEnabledStatusLabel = (isEnabled: boolean) => isEnabled ? '启用' : '停用'

const RoleSubmitSchema = RoleFormSchema.superRefine((values, ctx) => {
  if (!values.name.trim()) {
    ctx.addIssue({
      code: 'custom',
      path: ['name'],
      message: '请输入角色名称',
    })
  }
  if (!values.code.trim()) {
    ctx.addIssue({
      code: 'custom',
      path: ['code'],
      message: '请输入角色标识',
    })
  }
})

function toFormValues(record: SysRole): RoleFormValues {
  return RoleFormSchema.parse({
    ...defaultFormValues,
    ...RoleFormSchema.partial().parse({
      name: record.name ?? '',
      code: record.code ?? '',
      parentID: record.parentID ?? undefined,
      isEnabled: Boolean(record.isEnabled),
      remark: record.remark ?? '',
    }),
  })
}

function roleTreeToSelectData(items: SysRole[], editingID?: number): any[] {
  return items.map(item => ({
    title: `${item.name}（${item.code}）`,
    value: item.id,
    disabled: item.id === editingID,
    children: item.children ? roleTreeToSelectData(item.children, editingID) : undefined,
  }))
}

function menuTreeToTreeData(items: ResourceMenu[]): any[] {
  return items.map(item => ({
    title: item.name || item.path,
    key: item.id,
    children: item.children ? menuTreeToTreeData(item.children) : undefined,
  }))
}

function flattenResourceMenus(items: ResourceMenu[]) {
  const result: ResourceMenu[] = []
  const walk = (nodes: ResourceMenu[]) => {
    for (const node of nodes) {
      result.push(node)
      if (node.children?.length) {
        walk(node.children)
      }
    }
  }
  walk(items)
  return result
}

function buildMenuApiMap(items: ResourceMenu[]) {
  const map = new Map<number, number[]>()
  for (const item of flattenResourceMenus(items)) {
    map.set(item.id, item.apiIDs ?? [])
  }
  return map
}

function collectRoleRowKeys(items: SysRole[]) {
  const keys: Key[] = []
  const walk = (nodes: SysRole[]) => {
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

function RoleManagement() {
  const [expandedRowKeys, setExpandedRowKeys] = useState<Key[]>([])
  const [drawerOpen, setDrawerOpen] = useState(false)
  const [permissionDrawerOpen, setPermissionDrawerOpen] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [permissionSubmitting, setPermissionSubmitting] = useState(false)
  const [editing, setEditing] = useState<SysRole>()
  const [permissionRole, setPermissionRole] = useState<SysRole>()
  const [searchText, setSearchText] = useState('')
  const [enabledFilter, setEnabledFilter] = useState<boolean | undefined>()
  const [roleTree, setRoleTree] = useState<SysRole[]>([])
  const [menuTree, setMenuTree] = useState<ResourceMenu[]>([])
  const [apiItems, setApiItems] = useState<ResourceApi[]>([])
  const [apiSearchText, setApiSearchText] = useState('')
  const [selectedMenuKeys, setSelectedMenuKeys] = useState<Key[]>([])
  const [selectedApiKeys, setSelectedApiKeys] = useState<Key[]>([])
  const enabledStatus = useDictMatch(DictCode.SYS_IS_ENABLED_DICT_CODE)

  const {
    data,
    loading,
    send,
  } = usePagination(
    (nextPage, nextPageSize) => {
      const filters: Record<string, unknown>[] = []
      const keyword = searchText.trim()
      if (keyword) {
        filters.push({
          $or: [
            { name__icontains: keyword },
            { code__icontains: keyword },
          ],
        })
      }
      if (enabledFilter !== undefined) {
        filters.push({ isEnabled: enabledFilter })
      }

      return RoleApi.apiList({
        page: nextPage,
        pageSize: nextPageSize,
        noPaging: true,
        orderBy: 'id asc',
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

  const roleOptions = useMemo(() => roleTreeToSelectData(roleTree, editing?.id), [editing?.id, roleTree])
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
  const menuTreeData = useMemo(() => menuTreeToTreeData(menuTree), [menuTree])
  const menuTotal = useMemo(() => flattenResourceMenus(menuTree).length, [menuTree])
  const menuApiMap = useMemo(() => buildMenuApiMap(menuTree), [menuTree])
  const filteredApiItems = useMemo(() => {
    const keyword = apiSearchText.trim().toLowerCase()
    if (!keyword) {
      return apiItems
    }
    return apiItems.filter(item =>
      [
        item.module,
        item.method,
        item.path,
      ].some(value => String(value ?? '').toLowerCase().includes(keyword)),
    )
  }, [apiItems, apiSearchText])

  const loadRoleTree = async () => {
    const response = await RoleApi.apiTree().send()
    setRoleTree(response.data ?? [])
  }
  const { form, rules, onFinish } = useZodForm<RoleFormValues>({
    schema: RoleSubmitSchema,
    async onSubmit(values) {
      if (!values) {
        gMessage.error('请填写完整信息')
        return
      }

      setSubmitting(true)
      try {
        const payload = {
          name: values.name.trim(),
          code: values.code.trim(),
          parentID: values.parentID,
          isEnabled: values.isEnabled ?? editing?.isEnabled ?? true,
          remark: values.remark?.trim() ?? '',
        }
        if (editing) {
          await RoleApi.apiUpdate({
            id: editing.id,
            ...payload,
            parentID: payload.parentID ?? 0,
          })
        }
        else {
          await RoleApi.apiCreate(payload)
        }
        gMessage.success('保存成功')
        setDrawerOpen(false)
        await Promise.all([send(), loadRoleTree()])
      }
      catch {
        gMessage.error('保存失败')
      }
      finally {
        setSubmitting(false)
      }
    },
  })

  const saveRole = async () => {
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

  const openCreateForm = async (parentID?: number) => {
    setEditing(undefined)
    form.resetFields()
    form.setFieldsValue({
      ...defaultFormValues,
      parentID,
    })
    await loadRoleTree()
    setDrawerOpen(true)
  }

  const openEditForm = async (record: SysRole) => {
    setEditing(record)
    form.resetFields()
    form.setFieldsValue(toFormValues(record))
    await loadRoleTree()
    setDrawerOpen(true)
  }

  const toggleEnabled = async (record: SysRole) => {
    try {
      await RoleApi.apiUpdate({
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

  const openPermissionDrawer = async (record: SysRole) => {
    setPermissionRole(record)
    setApiSearchText('')
    setSelectedMenuKeys([])
    setSelectedApiKeys([])
    try {
      const [permissionRes, menuRes, apiRes] = await Promise.all([
        RoleApi.apiPermissions(record.id).send(),
        ResourceMenuApi.menuList({
          noPaging: true,
          orderBy: 'sort_order asc,id asc',
        }).send(),
        ResourceApiApi.apiList({
          noPaging: true,
          orderBy: 'sort_order asc,id desc',
        }).send(),
      ])
      setSelectedMenuKeys(permissionRes.data?.menuIDs ?? [])
      setSelectedApiKeys(permissionRes.data?.apiIDs ?? [])
      setMenuTree(menuRes.data?.items ?? [])
      setApiItems(apiRes.data?.items ?? [])
      setPermissionDrawerOpen(true)
    }
    catch {
      gMessage.error('加载授权信息失败')
    }
  }

  const savePermissions = async () => {
    if (!permissionRole) {
      return
    }
    setPermissionSubmitting(true)
    try {
      await RoleApi.apiSavePermissions({
        id: permissionRole.id,
        menuIDs: selectedMenuKeys.map(Number),
        apiIDs: selectedApiKeys.map(Number),
      })
      gMessage.success('授权保存成功')
      setPermissionDrawerOpen(false)
      await send()
    }
    catch {
      gMessage.error('授权保存失败')
    }
    finally {
      setPermissionSubmitting(false)
    }
  }

  const columns: ProColumns<SysRole>[] = [
    {
      title: '角色名称',
      dataIndex: 'name',
      width: 160,
      ellipsis: true,
    },
    {
      title: '角色标识',
      dataIndex: 'code',
      width: 140,
      ellipsis: true,
    },
    {
      title: '状态',
      dataIndex: 'isEnabled',
      width: 90,
      render: (_, record) => (
        enabledStatus.renderLabel(enabledStatusValue(record.isEnabled), fallbackEnabledStatusLabel(record.isEnabled))
      ),
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
      width: 260,
      fixed: 'right',
      render: (_, record) => (
        <Space size="small">
          <Button type="link" size="small" disabled={record.canWrite === false} onClick={() => openCreateForm(record.id)}>
            新增
          </Button>
          <Button type="link" size="small" disabled={record.canWrite === false} onClick={() => toggleEnabled(record)}>
            {enabledStatus.getLabel(enabledStatusValue(!record.isEnabled), fallbackEnabledStatusLabel(!record.isEnabled))}
          </Button>
          <Button type="link" size="small" disabled={record.canWrite === false} onClick={() => openEditForm(record)}>
            编辑
          </Button>
          <Button type="link" size="small" disabled={record.canWrite === false} onClick={() => openPermissionDrawer(record)}>
            授权
          </Button>
          <Popconfirm
            title="确认删除该角色？"
            onConfirm={async () => {
              try {
                await RoleApi.apiDel({ ids: [record.id] })
                gMessage.success('删除成功')
                await Promise.all([send(), loadRoleTree()])
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

  const apiColumns = [
    {
      title: '模块',
      dataIndex: 'module',
      width: 140,
      render: (value: string) => value || '-',
    },
    {
      title: '方法',
      dataIndex: 'method',
      width: 90,
      render: (value: string) => <Tag>{value}</Tag>,
    },
    {
      title: '路径',
      dataIndex: 'path',
      ellipsis: true,
    },
    {
      title: '备注',
      dataIndex: 'remark',
      ellipsis: true,
      render: (value: string) => value || '-',
    },
  ]

  return (
    <>
      <ProTable<SysRole>
        rowKey="id"
        columns={columns}
        dataSource={data}
        loading={loading}
        headerTitle="角色管理"
        search={false}
        scroll={{ x: 1100 }}
        pagination={false}
        expandable={{
          expandedRowKeys,
          onExpandedRowsChange: keys => setExpandedRowKeys([...keys]),
        }}
        options={{
          reload: () => send(),
        }}
        toolBarRender={() => [
          <Button key="create" type="primary" onClick={() => openCreateForm()}>
            创建角色
          </Button>,
          <Button key="expand" onClick={() => setExpandedRowKeys(collectRoleRowKeys(data))}>
            展开全部
          </Button>,
          <Button key="collapse" onClick={() => setExpandedRowKeys([])}>
            折叠全部
          </Button>,
          <Input.Search
            key="search"
            placeholder="搜索角色名称或标识"
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
        title={editing ? '编辑角色' : '创建角色'}
        open={drawerOpen}
        size={640}
        destroyOnHidden
        forceRender
        onClose={() => setDrawerOpen(false)}
        extra={(
          <Space>
            <Button onClick={() => setDrawerOpen(false)}>取消</Button>
            <Button type="primary" loading={submitting} onClick={saveRole}>
              保存
            </Button>
          </Space>
        )}
      >
        <Form<RoleFormValues>
          form={form}
          layout="horizontal"
          labelCol={{ span: 5 }}
          wrapperCol={{ span: 19 }}
          onFinish={onFinish}
        >
          <ProFormText
            name="name"
            label="角色名称"
            placeholder="请输入角色名称"
            rules={rules}
            fieldProps={{ maxLength: 255, showCount: true }}
          />
          <ProFormText
            name="code"
            label="角色标识"
            placeholder="例如：admin"
            rules={rules}
            fieldProps={{ maxLength: 128, showCount: true }}
          />
          <Form.Item name="parentID" label="父级角色" rules={rules}>
            <TreeSelect
              allowClear
              showSearch
              treeDefaultExpandAll
              placeholder="请选择父级角色"
              treeData={roleOptions}
              filterTreeNode={(input, node) => String(node.title ?? '').toLowerCase().includes(input.toLowerCase())}
            />
          </Form.Item>
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

      <Drawer
        title={permissionRole ? `角色授权：${permissionRole.name}` : '角色授权'}
        open={permissionDrawerOpen}
        size={960}
        destroyOnHidden
        onClose={() => setPermissionDrawerOpen(false)}
        extra={(
          <Space>
            <Button onClick={() => setPermissionDrawerOpen(false)}>取消</Button>
            <Button type="primary" loading={permissionSubmitting} onClick={savePermissions}>
              保存授权
            </Button>
          </Space>
        )}
      >
        <Space orientation="vertical" size={16} style={{ width: '100%' }}>
          <section>
            <div style={{ marginBottom: 8, fontWeight: 600 }}>
              菜单权限（
              {selectedMenuKeys.length}
              /
              {menuTotal}
              ）
            </div>
            <Tree
              checkable
              defaultExpandAll
              checkedKeys={selectedMenuKeys}
              treeData={menuTreeData}
              onCheck={(keys) => {
                const nextKeys = Array.isArray(keys) ? keys : keys.checked
                const previous = new Set(selectedMenuKeys.map(Number))
                const nextApiKeys = new Set(selectedApiKeys.map(Number))
                for (const key of nextKeys) {
                  const menuID = Number(key)
                  if (previous.has(menuID)) {
                    continue
                  }
                  for (const apiID of menuApiMap.get(menuID) ?? []) {
                    nextApiKeys.add(apiID)
                  }
                }
                setSelectedMenuKeys(nextKeys)
                setSelectedApiKeys([...nextApiKeys])
              }}
            />
          </section>
          <section>
            <div style={{ marginBottom: 8, fontWeight: 600 }}>
              API 权限（
              {selectedApiKeys.length}
              /
              {apiItems.length}
              ）
            </div>
            <Input.Search
              allowClear
              placeholder="搜索模块、方法或路径"
              value={apiSearchText}
              onChange={(event) => {
                setApiSearchText(event.target.value)
              }}
              onSearch={setApiSearchText}
              style={{ marginBottom: 12, width: 320 }}
            />
            <Table<ResourceApi>
              rowKey="id"
              size="small"
              columns={apiColumns}
              dataSource={filteredApiItems}
              pagination={{ pageSize: 8, showSizeChanger: false }}
              rowSelection={{
                preserveSelectedRowKeys: true,
                selectedRowKeys: selectedApiKeys,
                onChange: setSelectedApiKeys,
              }}
            />
          </section>
        </Space>
      </Drawer>
    </>
  )
}
