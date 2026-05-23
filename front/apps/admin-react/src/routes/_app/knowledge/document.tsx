import type { ProColumns } from '@ant-design/pro-components'
import type { KnowledgeDocument } from '~/api/business/knowledgeDocument'
import { ProFormDigit, ProFormSelect, ProFormText, ProFormTextArea, ProTable } from '@ant-design/pro-components'
import { createFileRoute, useRouter } from '@tanstack/react-router'
import { DEFAULT_PAGE_SIZE } from '@vp/core'
import { formatDateYYYYMMDDHHmmss } from '@vp/utils'
import { usePagination } from 'alova/client'
import { Button, Card, Drawer, Form, Input, Popconfirm, Space, Tag } from 'antd'
import { useCallback, useEffect, useMemo, useState } from 'react'
import z from 'zod'
import { KnowledgeCollectionApi } from '~/api/business/knowledgeCollection'
import { KnowledgeDocumentApi } from '~/api/business/knowledgeDocument'
import { useDictMatch } from '~/hooks/useDictMatch'
import { useZodForm } from '~/utils/zod'
import { gMessage } from '~/utils/message'

const searchSchema = z.object({
  collectionId: z.number().optional().catch(0),
})

export const Route = createFileRoute('/_app/knowledge/document')({
  validateSearch: searchSchema,
  staleTime: 1000 * 60 * 2,
  component: KnowledgeDocumentManagement,
})

const KnowledgeDocumentFormSchema = z.object({
  collectionID: z.number().optional().nullable(),
  documentID: z.string(),
  title: z.string(),
  content: z.string(),
  contentType: z.string().optional(),
  source: z.string().optional(),
  chunkIndex: z.number().optional().nullable(),
  totalChunks: z.number().optional().nullable(),
  metadata: z.string().optional(),
  isEnabled: z.number(),
  remark: z.string().optional(),
})

type KnowledgeDocumentFormValues = z.infer<typeof KnowledgeDocumentFormSchema>

const defaultFormValues: KnowledgeDocumentFormValues = {
  collectionID: undefined,
  documentID: '',
  title: '',
  content: '',
  contentType: 'text',
  source: '',
  chunkIndex: 0,
  totalChunks: 1,
  metadata: '',
  isEnabled: 1,
  remark: '',
}

const KnowledgeDocumentSubmitSchema = KnowledgeDocumentFormSchema.superRefine((values, ctx) => {
  for (const [field, label] of [
    ['documentID', '文档ID'],
    ['title', '标题'],
    ['content', '内容'],
  ] as const) {
    if (!values[field]?.toString().trim()) {
      ctx.addIssue({ code: 'custom', path: [field], message: `${label}不能为空` })
    }
  }
  if (values.metadata?.trim()) {
    try {
      JSON.parse(values.metadata.trim())
    }
    catch {
      ctx.addIssue({ code: 'custom', path: ['metadata'], message: '元数据必须是合法 JSON' })
    }
  }
})

function KnowledgeDocumentManagement() {
  const router = useRouter()
  const rawSearch = router.state.location.search as Record<string, unknown>
  const parsed = searchSchema.safeParse(rawSearch)
  const collectionId = Number(parsed.success ? parsed.data.collectionId : rawSearch.collectionId) || 0
  const [searchText, setSearchText] = useState('')
  const [collectionName, setCollectionName] = useState('')
  const [drawerOpen, setDrawerOpen] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [editing, setEditing] = useState<KnowledgeDocument>()
  const [form] = Form.useForm<KnowledgeDocumentFormValues>()
  const contentTypeDict = useDictMatch(DictCode.KNOWLEDGE_CONTENT_TYPE_DICT_CODE)
  const vectorStatusDict = useDictMatch(DictCode.KNOWLEDGE_VECTOR_STATUS_DICT_CODE)

  const contentTypeOptions = useMemo(() =>
    contentTypeDict.entries.map(entry => ({
      label: contentTypeDict.getLabel(entry.entryValue, entry.entryLabel),
      value: entry.entryValue,
    })),
  [contentTypeDict])

  const vectorStatusColor: Record<string, string> = {
    pending: 'default',
    indexed: 'success',
    failed: 'error',
  }

  // 加载 Collection 信息
  useEffect(() => {
    if (!collectionId) return
    KnowledgeCollectionApi.detail({ id: collectionId })
      .then((res) => {
        if (res.data) {
          setCollectionName(res.data.displayName || res.data.collectionName)
        }
      })
      .catch(() => {})
  }, [collectionId])

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
            { title__icontains: keyword },
            { documentID__icontains: keyword },
          ],
        })
      }
      return KnowledgeDocumentApi.listByCollection({
        collectionID: collectionId,
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
      watchingStates: [searchText, collectionId],
      debounce: [500, 0],
    },
  )

  const { rules, onFinish } = useZodForm<KnowledgeDocumentFormValues>({
    form,
    schema: KnowledgeDocumentSubmitSchema,
    async onSubmit(values) {
      if (!values) {
        return
      }
      setSubmitting(true)
      try {
        const payload = {
          collectionID: collectionId || values.collectionID || 0,
          documentID: values.documentID.trim(),
          title: values.title.trim(),
          content: values.content.trim(),
          contentType: values.contentType?.trim() || 'text',
          source: values.source?.trim() || '',
          chunkIndex: values.chunkIndex ?? 0,
          totalChunks: values.totalChunks ?? 1,
          metadata: values.metadata?.trim() || '',
          isEnabled: Boolean(values.isEnabled),
          remark: values.remark?.trim() ?? '',
        }
        if (editing) {
          await KnowledgeDocumentApi.update({ id: editing.id, ...payload })
        }
        else {
          await KnowledgeDocumentApi.create(payload)
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
    form.setFieldsValue({
      ...defaultFormValues,
      collectionID: collectionId || undefined,
    })
    setDrawerOpen(true)
  }

  const openEditForm = useCallback((record: KnowledgeDocument) => {
    setEditing(record)
    form.resetFields()
    form.setFieldsValue({
      collectionID: record.collectionID,
      documentID: record.documentID,
      title: record.title,
      content: record.content,
      contentType: record.contentType,
      source: record.source,
      chunkIndex: record.chunkIndex,
      totalChunks: record.totalChunks,
      metadata: record.metadata ? JSON.stringify(record.metadata) : '',
      isEnabled: record.isEnabled ? 1 : 0,
      remark: record.remark,
    })
    setDrawerOpen(true)
  }, [form])

  const action = useMemo(() => async (fn: () => Promise<void>, success: string, fail: string) => {
    try {
      await fn()
      gMessage.success(success)
      await send()
    }
    catch (e) {
      notifyError(e, fail)
    }
  }, [send])

  const columns: ProColumns<KnowledgeDocument>[] = useMemo(() => [
    { title: '文档ID', dataIndex: 'documentID', width: 200, ellipsis: true },
    { title: '标题', dataIndex: 'title', width: 200, ellipsis: true },
    {
      title: '内容类型',
      dataIndex: 'contentType',
      width: 100,
      render: (_, record) => contentTypeDict.renderLabel(
        record.contentType,
        <Tag>{record.contentType}</Tag>,
      ),
    },
    {
      title: '向量状态',
      dataIndex: 'vectorStatus',
      width: 100,
      render: (_, record) => vectorStatusDict.renderLabel(
        record.vectorStatus,
        <Tag color={vectorStatusColor[record.vectorStatus] || 'default'}>{record.vectorStatus}</Tag>,
      ),
    },
    { title: '分块', dataIndex: 'chunkIndex', width: 80, render: (_, record) => `${record.chunkIndex + 1}/${record.totalChunks}` },
    { title: '创建时间', dataIndex: 'createdAt', width: 160, render: (_, record) => formatDateYYYYMMDDHHmmss(record.createdAt) },
    {
      title: '操作',
      valueType: 'option',
      width: 200,
      fixed: 'right',
      render: (_, record) => (
        <Space>
          <Button type="link" size="small" disabled={record.canWrite === false} onClick={() => openEditForm(record)}>编辑</Button>
          <Popconfirm title="确认删除该文档？" onConfirm={() => action(() => KnowledgeDocumentApi.del({ id: record.id }), '删除成功', '删除失败')}>
            <Button type="link" size="small" danger disabled={record.canDelete === false}>删除</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ], [action, contentTypeDict, openEditForm, vectorStatusDict])

  return (
    <>
      {/* 顶部 Collection 信息卡片 */}
      <Card style={{ marginBottom: 16 }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <span style={{ fontSize: 16, fontWeight: 600 }}>所属集合：</span>
          <Tag color="blue">{collectionName || '未选择'}</Tag>
          {!collectionId && <span style={{ color: '#999' }}>请从集合列表选择一个集合查看文档</span>}
        </div>
      </Card>

      <ProTable<KnowledgeDocument>
        rowKey="id"
        headerTitle="文档管理"
        columns={columns}
        dataSource={data}
        loading={loading}
        search={false}
        scroll={{ x: 1100 }}
        pagination={{
          showSizeChanger: true,
          current: page,
          pageSize,
          total,
          onChange: (nextPage, nextPageSize) => update({ page: nextPage, pageSize: nextPageSize }),
        }}
        options={{ reload: () => send() }}
        toolBarRender={() => [
          <Button key="create" type="primary" onClick={openCreateForm}>创建文档</Button>,
          <Input.Search
            key="search"
            allowClear
            placeholder="搜索标题、文档ID"
            value={searchText}
            onChange={e => setSearchText(e.target.value)}
            onSearch={setSearchText}
            style={{ width: 280 }}
          />,
        ]}
      />

      <Drawer
        title={editing ? '编辑文档' : '创建文档'}
        size={600}
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
          <ProFormText name="documentID" label="文档ID" fieldProps={{ maxLength: 255 }} rules={rules} />
          <ProFormText name="title" label="标题" fieldProps={{ maxLength: 512 }} rules={rules} />
          <ProFormTextArea name="content" label="内容" fieldProps={{ rows: 6 }} rules={rules} />
          <ProFormSelect name="contentType" label="内容类型" options={contentTypeOptions} />
          <ProFormText name="source" label="来源" fieldProps={{ maxLength: 512 }} />
          <ProFormDigit name="chunkIndex" label="分块索引" min={0} precision={0} />
          <ProFormDigit name="totalChunks" label="总分块数" min={1} precision={0} />
          <ProFormTextArea name="metadata" label="元数据 (JSON)" fieldProps={{ rows: 3 }} />
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
