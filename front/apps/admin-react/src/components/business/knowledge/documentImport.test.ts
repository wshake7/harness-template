import * as React from 'react'
import { act } from 'react'
import { createRoot } from 'react-dom/client'
import { afterEach, beforeEach, expect, test, vi } from 'vite-plus/test'

const h = React.createElement
const uploadDirect = vi.fn()
const importFile = vi.fn()
const detail = vi.fn(() => Promise.resolve({ data: { displayName: 'Biz 1', collectionName: 'biz-1' } }))
const sendMock = vi.fn(() => Promise.resolve())
const success = vi.fn()
const mainForm = {
  resetFields: vi.fn(),
  setFieldsValue: vi.fn(),
  submit: vi.fn(),
}
const importForm = {
  resetFields: vi.fn(),
  setFieldsValue: vi.fn(),
  getFieldValue: vi.fn(() => ''),
  validateFields: vi.fn(() => Promise.resolve({
    title: 'Runbook',
    contentType: 'markdown',
    remark: 'imported',
  })),
}

function flush() {
  return new Promise(resolve => setTimeout(resolve, 0))
}

vi.mock('@tanstack/react-router', () => ({
  createFileRoute: () => () => ({}),
  useRouter: () => ({
    state: {
      location: {
        search: {
          collectionId: 1,
        },
      },
    },
  }),
}))

vi.mock('alova/client', () => ({
  usePagination: () => ({
    data: [],
    total: 0,
    page: 1,
    pageSize: 20,
    loading: false,
    update: vi.fn(),
    send: sendMock,
  }),
}))

vi.mock('~/api/business/knowledgeCollection', () => ({
  KnowledgeCollectionApi: {
    detail,
  },
}))

vi.mock('~/api/business/knowledgeDocument', () => ({
  KnowledgeDocumentApi: {
    listByCollection: vi.fn(() => Promise.resolve({ data: { items: [], total: 0 } })),
    create: vi.fn(),
    update: vi.fn(),
    del: vi.fn(),
    importFile,
  },
}))

vi.mock('~/api/business/storageFile', () => ({
  StorageFileApi: {
    uploadDirect,
  },
}))

vi.mock('~/hooks/useDictMatch', () => ({
  useDictMatch: () => ({
    entries: [],
    getLabel: (_value: string, fallback: string) => fallback,
    renderLabel: (_value: string, fallback: React.ReactNode) => fallback,
  }),
}))

vi.mock('~/utils/message', () => ({
  gMessage: {
    success,
    error: vi.fn(),
  },
}))

vi.mock('~/utils/zod', () => ({
  useZodForm: ({ onSubmit }: any) => ({
    rules: [],
    onFinish: onSubmit,
  }),
}))

vi.mock('@ant-design/pro-components', () => ({
  ProTable: ({ toolBarRender }: any) => h('div', null, toolBarRender?.().filter(Boolean)),
  ProFormDigit: () => null,
  ProFormSelect: () => null,
  ProFormText: () => null,
  ProFormTextArea: () => null,
}))

vi.mock('antd', () => ({
  Button: ({ children, onClick, loading, ...props }: any) => h('button', { type: 'button', onClick, ...(loading ? { 'data-loading': 'true' } : {}), ...props }, children),
  Card: ({ children }: any) => h('div', null, children),
  Drawer: ({ children, open, title, extra }: any) => open ? h('div', { 'data-testid': 'drawer' }, h('div', null, title), extra, children) : null,
  Form: Object.assign(
    ({ children }: any) => h('div', null, children),
    {
      useForm: vi.fn(() => {
        const index = (globalThis as any).__mockFormUseCount__ || 0
        ;(globalThis as any).__mockFormUseCount__ = index + 1
        return [index === 0 ? mainForm : importForm]
      }),
      Item: ({ children }: any) => h('div', null, children),
    },
  ),
  Input: Object.assign(
    ({ value, onChange, ...props }: any) => h('input', { value, onChange, ...props }),
    {
      Search: ({ value, onChange }: any) => h('input', { value, onChange }),
      TextArea: ({ value, onChange, ...props }: any) => h('textarea', { value, onChange, ...props }),
    },
  ),
  Popconfirm: ({ children }: any) => h('div', null, children),
  Space: ({ children }: any) => h('div', null, children),
  Tag: ({ children }: any) => h('span', null, children),
  Upload: ({ beforeUpload, children }: any) => h('div', null,
    h('button', {
      type: 'button',
      'data-testid': 'upload-file',
      onClick: () => beforeUpload?.(new File(['# Title'], 'runbook.md', { type: 'text/markdown' })),
    }, 'mock-upload'),
    children,
  ),
}))

let container: HTMLDivElement

beforeEach(() => {
  container = document.createElement('div')
  document.body.appendChild(container)
  uploadDirect.mockReset()
  importFile.mockReset()
  detail.mockClear()
  sendMock.mockClear()
  success.mockClear()
  ;(globalThis as any).__mockFormUseCount__ = 0
  ;(globalThis as any).IS_REACT_ACT_ENVIRONMENT = true
  vi.stubGlobal('DictCode', {
    KNOWLEDGE_CONTENT_TYPE_DICT_CODE: 'knowledge:content_type',
    KNOWLEDGE_VECTOR_STATUS_DICT_CODE: 'knowledge:vector_status',
  })
  vi.stubGlobal('notifyError', vi.fn())
})

afterEach(() => {
  container.remove()
})

test('imports uploaded file from the document page', async () => {
  uploadDirect.mockResolvedValue({
    data: {
      id: 10,
      originalName: 'runbook.md',
    },
  })
  importFile.mockResolvedValue({
    data: {
      id: 20,
      documentID: 'file-10',
    },
  })

  const { KnowledgeDocumentManagement } = await import('~/routes/_app/knowledge/document')
  const root = createRoot(container)

  await act(async () => {
    root.render(h(KnowledgeDocumentManagement))
  })

  const importButton = Array.from(container.querySelectorAll('button')).find(node => node.textContent === '导入文件')
  expect(importButton).toBeTruthy()

  await act(async () => {
    importButton?.dispatchEvent(new MouseEvent('click', { bubbles: true }))
    await flush()
  })

  await act(async () => {
    container.querySelector('[data-testid="upload-file"]')?.dispatchEvent(new MouseEvent('click', { bubbles: true }))
    await flush()
  })

  const submitButton = Array.from(container.querySelectorAll('button')).find(node => node.textContent === '开始导入')
  await act(async () => {
    submitButton?.dispatchEvent(new MouseEvent('click', { bubbles: true }))
    await flush()
  })

  expect(uploadDirect).toHaveBeenCalledTimes(1)
  expect(importFile).toHaveBeenCalledWith({
    collectionID: 1,
    fileAssetID: 10,
    title: 'Runbook',
    contentType: 'markdown',
    remark: 'imported',
  })
  expect(sendMock).toHaveBeenCalled()
  expect(success).toHaveBeenCalledWith('导入成功')
})
