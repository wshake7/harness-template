import type { ProColumns } from '@ant-design/pro-components'
import type { LanguageEntry, LanguageType } from '~/api/business/sysLanguage'
import { ModalForm, ProFormDigit, ProFormSwitch, ProFormText, ProTable } from '@ant-design/pro-components'
import { DEFAULT_PAGE_SIZE } from '@vp/core'
import { usePagination } from 'alova/client'
import {
  Button,
  Drawer,
  Form,
  Input,
  InputNumber,
  Popconfirm,
  Space,
  Tag,
} from 'antd'

import { useCallback, useMemo, useState } from 'react'
import z from 'zod'
import { LangApi } from '~/api/business/sysLanguage'
import { useDictMatch } from '~/hooks/useDictMatch'
import { gMessage } from '~/utils/message'

type EntryEditorMode = 'create' | 'edit'
const enabledStatusValue = (isEnabled: boolean) => isEnabled ? '1' : '0'
const fallbackEnabledStatusLabel = (isEnabled: boolean) => isEnabled ? '启用' : '停用'

function EntryEditorDrawer({
  mode,
  title,
  open,
  form,
  loading,
  page,
  pageSize,
  total,
  typeData,
  onClose,
  onSubmit,
  onTypePageChange,
}: {
  mode: EntryEditorMode
  title: string
  open: boolean
  form: ReturnType<typeof Form.useForm>[0]
  loading: boolean
  page: number
  pageSize: number
  total: number
  typeData: LanguageType[]
  onClose: () => void
  onSubmit: () => void
  onTypePageChange: (nextPage: number, nextPageSize: number) => void
}) {
  return (
    <Drawer
      title={title}
      open={open}
      onClose={onClose}
      width={680}
      placement="right"
      extra={(
        <Space>
          <Button onClick={onClose}>取消</Button>
          <Button type="primary" onClick={onSubmit}>保存</Button>
        </Space>
      )}
    >
      <Form form={form} layout="vertical">
        {mode === 'create'
          ? (
              <>
                <Form.Item
                  name="entryCode"
                  label="条目编码 (Code)"
                  rules={[{ required: true, message: '请输入条目编码' }]}
                >
                  <Input placeholder="如: common.submit" />
                </Form.Item>
                <Form.Item name="sortOrder" label="排序" initialValue={0}>
                  <InputNumber style={{ width: '100%' }} precision={0} />
                </Form.Item>
                <Form.Item name="remark" label="描述">
                  <Input placeholder="请输入描述" />
                </Form.Item>
              </>
            )
          : (
              <>
                <Form.Item name="currentRemark" label="描述">
                  <Input placeholder="请输入描述" />
                </Form.Item>
                <Form.Item name="currentSortOrder" label="排序">
                  <InputNumber style={{ width: '100%' }} precision={0} />
                </Form.Item>
              </>
            )}
        <div className="text-gray-500 mb-2 text-sm">语言值</div>
        <ProTable<LanguageType>
          rowKey="typeCode"
          search={false}
          options={false}
          loading={loading}
          pagination={{
            current: page,
            pageSize,
            total,
            showSizeChanger: true,
            pageSizeOptions: [8, 12, 20, 50],
            onChange: (nextPage, nextPageSize) => {
              onTypePageChange(nextPage, nextPageSize)
            },
          }}
          toolBarRender={false}
          dataSource={typeData}
          columns={[
            {
              title: '语言',
              dataIndex: 'typeName',
              width: 140,
            },
            {
              title: 'Code',
              dataIndex: 'typeCode',
              width: 120,
            },
            {
              title: '值',
              dataIndex: 'value',
              render: (_, record) => (
                <Form.Item
                  noStyle
                  name={['valuesByType', record.typeCode]}
                >
                  <Input
                    placeholder={mode === 'create' && record.isDefault ? '默认语言，建议填写' : `请输入${record.typeName}翻译`}
                    disabled={!record.isEnabled}
                  />
                </Form.Item>
              ),
            },
          ]}
        />
      </Form>
    </Drawer>
  )
}

const LangTypeSchema = z.object({
  typeCode: z.string('请输入语言编码').min(1, '请输入语言编码'),
  typeName: z.string('请输入语言名称').min(1, '请输入语言名称'),
  isDefault: z.boolean().default(false),
  isEnabled: z.boolean().default(true),
  sortOrder: z.number().default(0),
})

const langTypeDefaults = LangTypeSchema.partial().parse({})

export type LangTypeFormValues = z.infer<typeof LangTypeSchema>

export function LanguageTypePanel({
  selectedType,
  onSelectType,
  onDeleteSelectedType,
  onUpdateSelectedType,
}: {
  selectedType: LanguageType | undefined
  onSelectType: (record: LanguageType) => void
  onDeleteSelectedType: () => void
  onUpdateSelectedType: (record: LanguageType) => void
}) {
  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<LanguageType>()
  const [searchText, setSearchText] = useState('')
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
      const params: Record<string, unknown> = {
        page: nextPage,
        pageSize: nextPageSize,
        orderBy: 'sort_order asc,id desc',
      }
      if (searchText.trim()) {
        params.query = JSON.stringify({
          $or: [
            { typeCode__icontains: searchText.trim() },
            { typeName__icontains: searchText.trim() },
          ],
        })
      }
      return LangApi.typeList(params)
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
      watchingStates: [searchText],
      debounce: [500],
    },
  )

  const { form, rules, onFinish } = useZodForm<LangTypeFormValues>({
    schema: LangTypeSchema,
    async onSubmit(values) {
      if (!values) {
        gMessage.error('请填写完整信息')
        return
      }

      const payload = {
        typeCode: values.typeCode,
        typeName: values.typeName,
        isDefault: values.isDefault,
        isEnabled: values.isEnabled,
        sortOrder: values.sortOrder,
      }

      if (editing) {
        await LangApi.typeUpdate({
          id: editing.id,
          ...payload,
        })
        if (selectedType?.id === editing.id) {
          onUpdateSelectedType({
            ...selectedType,
            ...payload,
          })
        }
      }
      else {
        await LangApi.typeCreate(payload)
      }

      gMessage.success('保存成功')
      setEditing(undefined)
      form.resetFields()
      setFormOpen(false)
      await send()
    },
  })

  const openCreate = () => {
    setEditing(undefined)
    form.resetFields()
    form.setFieldsValue(langTypeDefaults as LangTypeFormValues)
    setFormOpen(true)
  }

  const openEdit = useCallback((record: LanguageType) => {
    setEditing(record)
    form.setFieldsValue({
      typeCode: record.typeCode,
      typeName: record.typeName,
      isDefault: record.isDefault,
      isEnabled: record.isEnabled,
      sortOrder: record.sortOrder,
    })
    setFormOpen(true)
  }, [form])

  const columns = useMemo<ProColumns<LanguageType>[]>(() => [
    {
      title: '序号',
      dataIndex: 'index',
      width: 60,
      render: (_, __, index) => (page - 1) * pageSize + index + 1,
    },
    {
      title: '类型编码',
      dataIndex: 'typeCode',
      width: 100,
    },
    {
      title: '语言名称',
      dataIndex: 'typeName',
      width: 140,
    },
    {
      title: '默认',
      dataIndex: 'isDefault',
      width: 70,
      render: (_, record) => (record.isDefault ? <Tag color="processing">默认</Tag> : '-'),
    },
    {
      title: '状态',
      dataIndex: 'isEnabled',
      width: 80,
      render: (_, record) => enabledStatus.renderLabel(enabledStatusValue(record.isEnabled), fallbackEnabledStatusLabel(record.isEnabled)),
    },
    {
      title: '排序',
      dataIndex: 'sortOrder',
      width: 70,
    },
    {
      title: '操作',
      valueType: 'option',
      fixed: 'right',
      width: 240,
      render: (_: unknown, record: LanguageType) => [
        !record.isDefault
          ? (
              <a
                key="setDefault"
                onClick={async (event) => {
                  event.stopPropagation()
                  await LangApi.typeUpdate({ id: record.id, isDefault: true })
                  gMessage.success('设置成功')
                  if (selectedType?.id === record.id) {
                    onUpdateSelectedType({ ...selectedType, isDefault: true })
                  }
                  await send()
                }}
              >
                设为默认
              </a>
            )
          : null,
        record.isDefault
          ? null
          : (
              <a
                key="switch"
                onClick={async (event) => {
                  event.stopPropagation()
                  await LangApi.typeUpdate({ id: record.id, isEnabled: !record.isEnabled })
                  gMessage.success('操作成功')
                  if (selectedType?.id === record.id) {
                    onUpdateSelectedType({ ...selectedType, isEnabled: !record.isEnabled })
                  }
                  await send()
                }}
              >
                {enabledStatus.getLabel(enabledStatusValue(!record.isEnabled), fallbackEnabledStatusLabel(!record.isEnabled))}
              </a>
            ),
        <a
          key="edit"
          onClick={(event) => {
            event.stopPropagation()
            openEdit(record)
          }}
        >
          编辑
        </a>,
        record.isDefault
          ? null
          : (
              <Popconfirm
                key="del"
                title="确认删除该语言类型吗？"
                onConfirm={async (event) => {
                  event?.stopPropagation()
                  await LangApi.typeDel({ ids: [record.id] })
                  gMessage.success('删除成功')
                  if (selectedType?.id === record.id) {
                    onDeleteSelectedType()
                  }
                  await send()
                }}
              >
                <a
                  onClick={(event) => {
                    event.stopPropagation()
                  }}
                >
                  删除
                </a>
              </Popconfirm>
            ),
      ],
    },
  ], [enabledStatus, openEdit, selectedType, onDeleteSelectedType, onUpdateSelectedType, send, page, pageSize])

  return (
    <>
      <ProTable<LanguageType>
        rowKey="id"
        search={false}
        columns={columns}
        dataSource={data}
        loading={loading}
        headerTitle={(
          <Space>
            语言类型
          </Space>
        )}
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
          <Button key="add" type="primary" onClick={openCreate}>
            新增语言
          </Button>,
          <Input.Search
            key="search"
            placeholder="搜索语言编码、名称"
            allowClear
            value={searchText}
            onChange={(e) => {
              setSearchText(e.target.value)
            }}
            onSearch={(value) => {
              setSearchText(value)
            }}
            style={{ width: 260 }}
          />,
        ]}
        rowClassName={(record) => {
          if (record.id === selectedType?.id) {
            return 'ant-table-row-selected'
          }
          return ''
        }}
        onRow={record => ({
          onClick: () => {
            onSelectType(record)
          },
        })}
      />
      <ModalForm
        title={editing ? '编辑语言类型' : '新增语言类型'}
        open={formOpen}
        onOpenChange={(open) => {
          if (!open) {
            setFormOpen(false)
            setEditing(undefined)
          }
        }}
        width={500}
        form={form}
        onFinish={onFinish}
        submitTimeout={2000}
      >
        <ProFormText required name="typeCode" label="语言编码" rules={rules} placeholder="如: zh-CN" />
        <ProFormText required name="typeName" label="语言名称" rules={rules} placeholder="如: 简体中文" />
        <ProFormDigit name="sortOrder" label="排序" fieldProps={{ precision: 0 }} />
        <ProFormSwitch name="isDefault" label="是否默认" />
        <ProFormSwitch
          name="isEnabled"
          label="启用状态"
          fieldProps={{
            checkedChildren: enabledStatus.getLabel(enabledStatusValue(true), fallbackEnabledStatusLabel(true)),
            unCheckedChildren: enabledStatus.getLabel(enabledStatusValue(false), fallbackEnabledStatusLabel(false)),
          }}
        />
      </ModalForm>
    </>
  )
}

export function LanguageEntryPanel({
  selectedType,
  refreshKey,
  onClearType,
}: {
  selectedType: LanguageType | undefined
  refreshKey: number
  onClearType: () => void
}) {
  const [searchText, setSearchText] = useState('')
  const [editDrawerOpen, setEditDrawerOpen] = useState(false)
  const [editDrawerData, setEditDrawerData] = useState<LanguageEntry[]>([])
  const [editingEntryId, setEditingEntryId] = useState<number>()
  const [drawerTypeCodeMap, setDrawerTypeCodeMap] = useState<Record<number, string>>({})
  const [editDrawerForm] = Form.useForm()
  const [batchCreateOpen, setBatchCreateOpen] = useState(false)
  const [batchCreateForm] = Form.useForm()
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
      const conditions: Record<string, unknown>[] = []
      if (selectedType) {
        conditions.push({ sysLanguageTypeId: String(selectedType.id) })
      }
      if (searchText.trim()) {
        conditions.push({
          $or: [
            { entryCode__icontains: searchText.trim() },
            { entryValue__icontains: searchText.trim() },
          ],
        })
      }
      return LangApi.entryList({
        page: nextPage,
        pageSize: nextPageSize,
        orderBy: 'sort_order asc,id desc',
        query: conditions.length > 0 ? JSON.stringify({ $and: conditions }) : undefined,
      })
    },
    {
      initialData: {
        total: 0,
        items: [],
      },
      initialPage: 1,
      initialPageSize: DEFAULT_PAGE_SIZE,
      watchingStates: [selectedType?.id, searchText, refreshKey],
      data: response => response.data?.items ?? [],
      total: response => response.data?.total ?? 0,
    },
  )

  const {
    data: drawerTypeData,
    total: drawerTypeTotal,
    page: drawerTypePage,
    pageSize: drawerTypePageSize,
    loading: drawerTypeLoading,
    update: updateDrawerTypes,
    send: sendDrawerTypes,
  } = usePagination(
    (nextPage, nextPageSize) => {
      return LangApi.typeList({
        page: nextPage,
        pageSize: nextPageSize,
        orderBy: 'sort_order asc,id desc',
      })
    },
    {
      initialData: {
        total: 0,
        items: [],
      },
      initialPage: 1,
      initialPageSize: 8,
      data: response => response.data?.items ?? [],
      total: response => response.data?.total ?? 0,
    },
  )

  const closeEditDrawer = useCallback(() => {
    setEditDrawerOpen(false)
    setEditingEntryId(undefined)
    setEditDrawerData([])
    setDrawerTypeCodeMap({})
    editDrawerForm.resetFields()
  }, [editDrawerForm])

  const closeBatchCreate = useCallback(() => {
    setBatchCreateOpen(false)
    batchCreateForm.resetFields()
  }, [batchCreateForm])

  const openEditDrawer = useCallback(async (record: LanguageEntry) => {
    editDrawerForm.resetFields()
    updateDrawerTypes({
      page: 1,
      pageSize: drawerTypePageSize,
    })
    await sendDrawerTypes()
    const allTypesRes = await LangApi.typeList({
      noPaging: true,
      pageSize: 1000,
      orderBy: 'sort_order asc,id desc',
    }).send()
    const allTypes = allTypesRes.data?.items ?? []
    const typeCodeMap: Record<number, string> = {}
    for (const type of allTypes) {
      typeCodeMap[type.id] = type.typeCode
    }
    setDrawerTypeCodeMap(typeCodeMap)
    const res = await LangApi.entryList({ noPaging: true, pageSize: 1000, query: JSON.stringify({ entryCode: record.entryCode }) }).send()
    const entries = res.data?.items ?? []
    setEditDrawerData(entries)
    setEditingEntryId(record.id)
    const valueFields: Record<string, string> = {}
    for (const entry of entries) {
      const typeCode = entry.sysLanguageType?.typeCode ?? typeCodeMap[entry.sysLanguageTypeId]
      if (!typeCode) {
        continue
      }
      valueFields[typeCode] = entry.entryValue
    }
    editDrawerForm.setFieldsValue({
      valuesByType: valueFields,
      currentSortOrder: record.sortOrder,
      currentRemark: record.remark,
    })
    setEditDrawerOpen(true)
  }, [drawerTypePageSize, updateDrawerTypes, sendDrawerTypes, editDrawerForm])

  const handleEditDrawerSave = useCallback(async () => {
    try {
      const values = await editDrawerForm.validateFields()
      const valuesByType: Record<string, string> = values.valuesByType ?? {}
      const updateMap = new Map<number, {
        id: number
        entryValue?: string
        sortOrder?: number
        remark?: string
      }>()
      for (const entry of editDrawerData) {
        const typeCode = entry.sysLanguageType?.typeCode ?? drawerTypeCodeMap[entry.sysLanguageTypeId]
        if (!typeCode) {
          continue
        }
        const newValue = valuesByType[typeCode] ?? ''
        if (newValue !== entry.entryValue) {
          updateMap.set(entry.id, {
            ...(updateMap.get(entry.id) ?? { id: entry.id }),
            entryValue: newValue,
          })
        }
      }
      const currentEntry = editDrawerData.find(entry => entry.id === editingEntryId)
      if (currentEntry) {
        const currentUpdatePayload = updateMap.get(currentEntry.id) ?? { id: currentEntry.id }
        if (values.currentSortOrder !== currentEntry.sortOrder) {
          currentUpdatePayload.sortOrder = values.currentSortOrder
        }
        if ((values.currentRemark ?? '') !== (currentEntry.remark ?? '')) {
          currentUpdatePayload.remark = values.currentRemark ?? ''
        }
        if (currentUpdatePayload.sortOrder !== undefined || currentUpdatePayload.remark !== undefined) {
          updateMap.set(currentEntry.id, currentUpdatePayload)
        }
      }
      const updates = Array.from(updateMap.values())
      if (updates.length > 0) {
        await LangApi.entryUpdate({
          updates,
        })
      }
      gMessage.success('保存成功')
      closeEditDrawer()
      await send()
    }
    catch {
      // validation error
    }
  }, [editDrawerForm, editDrawerData, editingEntryId, send, drawerTypeCodeMap, closeEditDrawer])

  const openBatchCreate = useCallback(async () => {
    updateDrawerTypes({
      page: 1,
      pageSize: drawerTypePageSize,
    })
    await sendDrawerTypes()
    batchCreateForm.resetFields()
    batchCreateForm.setFieldsValue({
      entryCode: '',
      sortOrder: 0,
      remark: '',
      valuesByType: {},
    })
    setBatchCreateOpen(true)
  }, [drawerTypePageSize, updateDrawerTypes, sendDrawerTypes, batchCreateForm])

  const handleBatchCreateSave = useCallback(async () => {
    try {
      const values = await batchCreateForm.validateFields()
      const allTypesRes = await LangApi.typeList({
        noPaging: true,
        pageSize: 1000,
        orderBy: 'sort_order asc,id desc',
      }).send()
      const allTypes = allTypesRes.data?.items ?? []
      const valuesByType: Record<string, string> = values.valuesByType ?? {}
      const valuesMap: Record<string, string> = {}
      for (const type of allTypes) {
        valuesMap[type.typeCode] = valuesByType[type.typeCode] ?? ''
      }
      await LangApi.entryBatchCreate({
        entryCode: values.entryCode,
        values: valuesMap,
        sortOrder: values.sortOrder ?? 0,
      })
      if ((values.remark ?? '').trim()) {
        const createEntriesRes = await LangApi.entryList({
          noPaging: true,
          pageSize: 1000,
          query: JSON.stringify({ entryCode: values.entryCode }),
        }).send()
        const createdEntries = createEntriesRes.data?.items ?? []
        await Promise.all(createdEntries.map(entry => LangApi.entryUpdate({
          id: entry.id,
          remark: values.remark.trim(),
        })))
      }
      gMessage.success('创建成功')
      closeBatchCreate()
      await send()
    }
    catch {
      // validation error
    }
  }, [batchCreateForm, send, closeBatchCreate])

  const handleToggleEntryEnabled = useCallback(async (record: LanguageEntry) => {
    const targetEnabled = !record.isEnabled
    const res = await LangApi.entryList({
      noPaging: true,
      pageSize: 1000,
      query: JSON.stringify({ entryCode: record.entryCode }),
    }).send()
    const relatedEntries = res.data?.items ?? []
    const updates = relatedEntries
      .filter(item => item.isEnabled !== targetEnabled)
      .map(item => ({
        id: item.id,
        isEnabled: targetEnabled,
      }))
    if (updates.length > 0) {
      await LangApi.entryUpdate({ updates })
    }
    gMessage.success(`${enabledStatus.getLabel(enabledStatusValue(targetEnabled), fallbackEnabledStatusLabel(targetEnabled))}成功`)
    await send()
  }, [enabledStatus, send])

  const columns = useMemo<ProColumns<LanguageEntry>[]>(() => [
    {
      title: '序号',
      dataIndex: 'index',
      width: 55,
      editable: false,
      render: (_: unknown, _record: LanguageEntry, index: number) => (page - 1) * pageSize + index + 1,
    },
    {
      title: 'Code',
      dataIndex: 'entryCode',
      width: 160,
      ellipsis: true,
      editable: false,
    },
    {
      title: '值',
      dataIndex: 'entryValue',
      fieldProps: {
        placeholder: '请输入值',
      },
      render: (_: unknown, record: LanguageEntry) =>
        record.entryValue || <span className="text-gray-300">(空)</span>,
    },
    {
      title: '排序',
      dataIndex: 'sortOrder',
      width: 65,
      valueType: 'digit',
      editable: false,
    },
    {
      title: '状态',
      dataIndex: 'isEnabled',
      width: 75,
      valueType: 'switch',
      editable: false,
      render: (_: unknown, record: LanguageEntry) => enabledStatus.renderLabel(enabledStatusValue(record.isEnabled), fallbackEnabledStatusLabel(record.isEnabled)),
    },
    {
      title: '备注',
      dataIndex: 'remark',
      ellipsis: true,
      editable: false,
    },
    {
      title: '操作',
      valueType: 'option',
      width: 180,
      fixed: 'right',
      render: (_, record: LanguageEntry, _index) => [
        <a
          key="switch"
          onClick={async (e: React.MouseEvent) => {
            e.stopPropagation()
            await handleToggleEntryEnabled(record)
          }}
        >
          {enabledStatus.getLabel(enabledStatusValue(!record.isEnabled), fallbackEnabledStatusLabel(!record.isEnabled))}
        </a>,
        <a
          key="modalEdit"
          onClick={(e: React.MouseEvent) => {
            e.stopPropagation()
            openEditDrawer(record)
          }}
        >
          编辑
        </a>,
        <Popconfirm
          key="del"
          title="确认删除该语言条目吗？"
          onConfirm={async () => {
            await LangApi.entryDel({ ids: [record.id] })
            gMessage.success('删除成功')
            await send()
          }}
        >
          <a onClick={e => e.stopPropagation()}>删除</a>
        </Popconfirm>,
      ],
    },
  ], [enabledStatus, page, pageSize, send, openEditDrawer, handleToggleEntryEnabled])

  return (
    <>
      <ProTable<LanguageEntry>
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
        headerTitle={(
          <Space>
            {selectedType
              ? (
                  <>
                    {`语言条目 - ${selectedType.typeName}`}
                    <Tag color="processing">{selectedType.typeCode}</Tag>
                    <Tag color="error" style={{ cursor: 'pointer' }} onClick={onClearType}>
                      清除筛选
                    </Tag>
                  </>
                )
              : '语言条目'}
          </Space>
        )}
        toolBarRender={() => [
          <Button key="batchAdd" type="primary" onClick={openBatchCreate}>
            新增
          </Button>,
          <Input.Search
            key="search"
            placeholder="搜索 Code、值"
            allowClear
            value={searchText}
            onChange={(e) => {
              setSearchText(e.target.value)
            }}
            onSearch={(value) => {
              setSearchText(value)
            }}
            style={{ width: 260 }}
          />,
        ]}
      />

      <EntryEditorDrawer
        mode="edit"
        title={`编辑条目 - ${editDrawerData[0]?.entryCode ?? ''}`}
        open={editDrawerOpen}
        form={editDrawerForm}
        loading={drawerTypeLoading}
        page={drawerTypePage ?? 1}
        pageSize={drawerTypePageSize ?? 8}
        total={drawerTypeTotal ?? 0}
        typeData={drawerTypeData}
        onClose={closeEditDrawer}
        onSubmit={handleEditDrawerSave}
        onTypePageChange={(nextPage, nextPageSize) => {
          updateDrawerTypes({
            page: nextPage,
            pageSize: nextPageSize,
          })
        }}
      />

      <EntryEditorDrawer
        mode="create"
        title="新增条目"
        open={batchCreateOpen}
        form={batchCreateForm}
        loading={drawerTypeLoading}
        page={drawerTypePage ?? 1}
        pageSize={drawerTypePageSize ?? 8}
        total={drawerTypeTotal ?? 0}
        typeData={drawerTypeData}
        onClose={closeBatchCreate}
        onSubmit={handleBatchCreateSave}
        onTypePageChange={(nextPage, nextPageSize) => {
          updateDrawerTypes({
            page: nextPage,
            pageSize: nextPageSize,
          })
        }}
      />

    </>
  )
}
