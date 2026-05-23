import type { ProColumns } from '@ant-design/pro-components'
import type { KnowledgeCollection } from '~/api/business/knowledgeCollection'
import {
  ProFormDigit,
  ProFormSelect,
  ProFormText,
  ProFormTextArea,
  ProTable,
} from '@ant-design/pro-components'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { DEFAULT_PAGE_SIZE } from '@vp/core'
import { formatDateYYYYMMDDHHmmss } from '@vp/utils'
import { usePagination } from 'alova/client'
import { Button, Drawer, Form, Input, Popconfirm, Space, Tag } from 'antd'
import { useCallback, useMemo, useState } from 'react'
import z from 'zod'
import { KnowledgeCollectionApi } from '~/api/business/knowledgeCollection'
import { useDictMatch } from '~/hooks/useDictMatch'
import { useZodForm } from '~/utils/zod'
import { gMessage } from '~/utils/message'

export const Route = createFileRoute('/_app/knowledge/collection')({
  staleTime: 1000 * 60 * 2,
  component: KnowledgeCollectionManagement,
})

const KnowledgeCollectionFormSchema = z.object({
  collectionName: z.string(),
  displayName: z.string(),
  description: z.string().optional(),
  embeddingModel: z.string(),
  vectorDimension: z.number().optional().nullable(),
  metricType: z.string().optional(),
  indexType: z.string().optional(),
  idMaxLength: z.number().optional().nullable(),
  contentMaxLength: z.number().optional().nullable(),
  isEnabled: z.number(),
  remark: z.string().optional(),
})

type KnowledgeCollectionFormValues = z.infer<typeof KnowledgeCollectionFormSchema>

const defaultFormValues: KnowledgeCollectionFormValues = {
  collectionName: '',
  displayName: '',
  description: '',
  embeddingModel: '',
  vectorDimension: undefined,
  metricType: '',
  indexType: '',
  idMaxLength: undefined,
  contentMaxLength: undefined,
  isEnabled: 1,
  remark: '',
}

const KnowledgeCollectionSubmitSchema = KnowledgeCollectionFormSchema.superRefine((values, ctx) => {
  for (const [field, label] of [
    ['collectionName', '集合名称'],
    ['displayName', '显示名称'],
    ['embeddingModel', 'Embedding模型'],
  ] as const) {
    if (!values[field]?.toString().trim()) {
      ctx.addIssue({ code: 'custom', path: [field], message: `${label}不能为空` })
    }
  }
  if (!values.vectorDimension || Number(values.vectorDimension) <= 0) {
    ctx.addIssue({ code: 'custom', path: ['vectorDimension'], message: '向量维度必须大于0' })
  }
})

function KnowledgeCollectionManagement() {
  const [drawerOpen, setDrawerOpen] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [editing, setEditing] = useState<KnowledgeCollection>()
  const [searchText, setSearchText] = useState('')
  const [form] = Form.useForm<KnowledgeCollectionFormValues>()
  const navigate = useNavigate()
  const enabledStatus = useDictMatch(DictCode.SYS_IS_ENABLED_DICT_CODE)
  const metricTypeDict = useDictMatch(DictCode.KNOWLEDGE_METRIC_TYPE_DICT_CODE)
  const indexTypeDict = useDictMatch(DictCode.KNOWLEDGE_INDEX_TYPE_DICT_CODE)

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
            { collectionName__icontains: keyword },
            { displayName__icontains: keyword },
          ],
        })
      }
      return KnowledgeCollectionApi.list({
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
      watchingStates: [searchText],
      debounce: [500, 0],
    },
  )

  const { rules, onFinish } = useZodForm<KnowledgeCollectionFormValues>({
    form,
    schema: KnowledgeCollectionSubmitSchema,
    async onSubmit(values) {
      if (!values) {
        return
      }
      setSubmitting(true)
      try {
        const payload = {
          collectionName: values.collectionName.trim(),
          displayName: values.displayName.trim(),
          description: values.description?.trim() ?? '',
          embeddingModel: values.embeddingModel.trim(),
          vectorDimension: Number(values.vectorDimension) || 0,
          metricType: values.metricType?.trim() ?? '',
          indexType: values.indexType?.trim() ?? '',
          idMaxLength: values.idMaxLength ? Number(values.idMaxLength) : undefined,
          contentMaxLength: values.contentMaxLength ? Number(values.contentMaxLength) : undefined,
          isEnabled: Boolean(values.isEnabled),
          remark: values.remark?.trim() ?? '',
        }
        if (editing) {
          await KnowledgeCollectionApi.update({ id: editing.id, ...payload })
        }
        else {
          await KnowledgeCollectionApi.create(payload)
        }
        gMessage.success('保存成功')
        setDrawerOpen(false)
        await send()
      }
      catch (e) {
        notifyError(e, '保存失败')
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

  const openEditForm = useCallback((record: KnowledgeCollection) => {
    setEditing(record)
    form.resetFields()
    form.setFieldsValue({
      collectionName: record.collectionName,
      displayName: record.displayName,
      description: record.description,
      embeddingModel: record.embeddingModel,
      vectorDimension: record.vectorDimension,
      metricType: record.metricType,
      indexType: record.indexType,
      idMaxLength: record.idMaxLength,
      contentMaxLength: record.contentMaxLength,
      isEnabled: record.isEnabled ? 1 : 0,
      remark: record.remark,
    })
    setDrawerOpen(true)
  }, [form])

  const action = useCallback(async (fn: () => Promise<void>, success: string, fail: string) => {
    try {
      await fn()
      gMessage.success(success)
      await send()
    }
    catch (e) {
      notifyError(e, fail)
    }
  }, [send])

  const handleViewDocuments = useCallback((record: KnowledgeCollection) => {
    const label = record.displayName || record.collectionName
    useMenuTabsStore.getState().add({
      key: '/knowledge/document',
      label: `文档 - ${label}`,
    })
    navigate({ to: '/knowledge/document', search: { collectionId: record.id } })
  }, [navigate])

  const metricTypeOptions = useMemo(() =>
    metricTypeDict.entries.map(entry => ({
      label: metricTypeDict.getLabel(entry.entryValue, entry.entryLabel),
      value: entry.entryValue,
    })),
  [metricTypeDict])

  const indexTypeOptions = useMemo(() =>
    indexTypeDict.entries.map(entry => ({
      label: indexTypeDict.getLabel(entry.entryValue, entry.entryLabel),
      value: entry.entryValue,
    })),
  [indexTypeDict])

  const columns: ProColumns<KnowledgeCollection>[] = useMemo(() => [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '集合名称', dataIndex: 'collectionName', width: 160, ellipsis: true },
    { title: '显示名称', dataIndex: 'displayName', width: 180, ellipsis: true },
    { title: 'Embedding模型', dataIndex: 'embeddingModel', width: 200, ellipsis: true },
    { title: '向量维度', dataIndex: 'vectorDimension', width: 100 },
    { title: '文档数量', dataIndex: 'documentCount', width: 100 },
    {
      title: '度量类型',
      dataIndex: 'metricType',
      width: 100,
      render: (_, record) => metricTypeDict.renderLabel(
        record.metricType,
        <Tag>{record.metricType}</Tag>,
      ),
    },
    {
      title: '索引类型',
      dataIndex: 'indexType',
      width: 100,
      render: (_, record) => indexTypeDict.renderLabel(
        record.indexType,
        <Tag>{record.indexType}</Tag>,
      ),
    },
    {
      title: '状态',
      dataIndex: 'isEnabled',
      width: 90,
      render: (_, record) => enabledStatus.renderLabel(
        record.isEnabled ? '1' : '0',
        <Tag color={record.isEnabled ? 'success' : 'default'}>{record.isEnabled ? '启用' : '停用'}</Tag>,
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
        const nextEnabled = !record.isEnabled
        const nextLabel = nextEnabled ? '启用' : '停用'

        return (
          <Space>
            <Button type="link" size="small" onClick={() => handleViewDocuments(record)}>
              查看文档
            </Button>
            <Button
              type="link"
              size="small"
              disabled={record.canWrite === false}
              onClick={() => action(
                () => KnowledgeCollectionApi.update({ id: record.id, isEnabled: nextEnabled }),
                `${nextLabel}成功`,
                `${nextLabel}失败`,
              )}
            >
              {nextLabel}
            </Button>
            <Button type="link" size="small" disabled={record.canWrite === false} onClick={() => openEditForm(record)}>
              编辑
            </Button>
            <Popconfirm
              title="确认删除该集合？"
              onConfirm={() => action(() => KnowledgeCollectionApi.del({ id: record.id }), '删除成功', '删除失败')}
            >
              <Button type="link" size="small" danger disabled={record.canDelete === false}>
                删除
              </Button>
            </Popconfirm>
          </Space>
        )
      },
    },
  ], [action, enabledStatus, handleViewDocuments, openEditForm])

  return (
    <>
      <ProTable<KnowledgeCollection>
        rowKey="id"
        headerTitle="知识库集合"
        columns={columns}
        dataSource={data}
        loading={loading}
        search={false}
        scroll={{ x: 1300 }}
        pagination={{
          showSizeChanger: true,
          current: page,
          pageSize,
          total,
          onChange: (nextPage, nextPageSize) => update({ page: nextPage, pageSize: nextPageSize }),
        }}
        options={{ reload: () => send() }}
        toolBarRender={() => [
          <Button key="create" type="primary" onClick={openCreateForm}>创建集合</Button>,
          <Input.Search
            key="search"
            allowClear
            placeholder="搜索集合名称、显示名称"
            value={searchText}
            onChange={event => setSearchText(event.target.value)}
            onSearch={setSearchText}
            style={{ width: 280 }}
          />,
        ]}
      />

      <Drawer
        title={editing ? '编辑集合' : '创建集合'}
        size={560}
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
          <ProFormText name="collectionName" label="集合名称" fieldProps={{ maxLength: 128 }} rules={rules} />
          <ProFormText name="displayName" label="显示名称" fieldProps={{ maxLength: 255 }} rules={rules} />
          <ProFormTextArea name="description" label="描述" fieldProps={{ rows: 3, maxLength: 512 }} />
          <ProFormText name="embeddingModel" label="Embedding模型" fieldProps={{ maxLength: 128 }} rules={rules} />
          <ProFormDigit name="vectorDimension" label="向量维度" min={1} precision={0} rules={rules} />
          <ProFormSelect name="metricType" label="度量类型" options={metricTypeOptions} />
          <ProFormSelect name="indexType" label="索引类型" options={indexTypeOptions} />
          <ProFormDigit name="idMaxLength" label="ID最大长度" min={1} precision={0} />
          <ProFormDigit name="contentMaxLength" label="Content最大长度" min={1} precision={0} />
          <ProFormSelect
            name="isEnabled"
            label="启用状态"
            options={[
              { label: '启用', value: 1 },
              { label: '停用', value: 0 },
            ]}
            rules={rules}
          />
          <ProFormTextArea name="remark" label="备注" fieldProps={{ rows: 2, maxLength: 255 }} />
        </Form>
      </Drawer>
    </>
  )
}
