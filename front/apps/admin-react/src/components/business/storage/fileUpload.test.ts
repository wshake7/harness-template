import type { Root } from 'react-dom/client'
import type { StorageFileAsset } from '~/api/business/storageFile'
import * as React from 'react'
import { act, useState } from 'react'
import { createRoot } from 'react-dom/client'
import { afterEach, beforeEach, expect, test, vi } from 'vite-plus/test'

const messageError = vi.fn()
const h = React.createElement

vi.mock('~/utils/message', () => ({
  gMessage: {
    error: messageError,
  },
}))

vi.mock('~/api/business/storageFile', () => ({
  StorageFileApi: {
    upload: vi.fn(),
    presigned: vi.fn(),
    detail: vi.fn(),
    del: vi.fn(),
  },
}))

vi.mock('@ant-design/icons', () => ({
  UploadOutlined: () => null,
}))

vi.mock('antd', () => ({
  Button: ({ children, disabled, ...props }: any) => h('button', { type: 'button', disabled, ...props }, children),
  Upload: ({ children, customRequest, fileList = [], onPreview, onRemove, disabled }: any) => h(
    'div',
    null,
    h('button', {
      'type': 'button',
      'data-testid': 'upload-trigger',
      'disabled': disabled,
      'onClick': () => {
        if (disabled) {
          return
        }
        customRequest?.({
          file: new File(['hello'], 'hello.txt', { type: 'text/plain' }),
          onSuccess: () => {},
          onError: () => {},
        })
      },
    }, 'mock-upload'),
    h('button', {
      'type': 'button',
      'data-testid': 'preview-trigger',
      'onClick': () => {
        const current = fileList[0]
        if (current) {
          onPreview?.(current)
        }
      },
    }, 'mock-preview'),
    h('button', {
      'type': 'button',
      'data-testid': 'remove-trigger',
      'onClick': () => {
        const current = fileList[0]
        if (current) {
          onRemove?.(current)
        }
      },
    }, 'mock-remove'),
    children,
  ),
}))

function flush() {
  return new Promise(resolve => setTimeout(resolve, 0))
}

let container: HTMLDivElement
let root: Root

beforeEach(() => {
  container = document.createElement('div')
  document.body.appendChild(container)
  root = createRoot(container)
  messageError.mockClear()
})

afterEach(() => {
  root.unmount()
  container.remove()
})

test('uploads a file and emits StorageFileAsset list', async () => {
  const upload = vi.fn().mockResolvedValue({
    data: {
      id: 1,
      engine: 'minio',
      bucket: 'admin-files',
      objectKey: 'uploads/mock/1-hello.txt',
      originalName: 'hello.txt',
      contentType: 'text/plain',
      size: 5,
      sha256: 'sha',
    },
  })
  const storageApi = {
    upload,
    presigned: vi.fn(),
    detail: vi.fn(),
    del: vi.fn(),
  }
  const { StorageFileUpload } = await import('./fileUpload')

  function Harness() {
    const [value, setValue] = useState<StorageFileAsset[]>([])
    return h(
      'div',
      null,
      h(StorageFileUpload, { value, onChange: setValue, storageApi: storageApi as any, bizType: 'attachment' }),
      h('div', { 'data-testid': 'names' }, value.map(item => item.originalName).join(',')),
    )
  }

  await act(async () => {
    root.render(h(Harness))
  })
  await act(async () => {
    container.querySelector('[data-testid="upload-trigger"]')?.dispatchEvent(new MouseEvent('click', { bubbles: true }))
    await flush()
  })

  expect(upload).toHaveBeenCalledTimes(1)
  expect(container.querySelector('[data-testid="names"]')?.textContent).toBe('hello.txt')
})

test('removes uploaded file from controlled value', async () => {
  const storageApi = {
    upload: vi.fn(),
    presigned: vi.fn(),
    detail: vi.fn(),
    del: vi.fn(),
  }
  const { StorageFileUpload } = await import('./fileUpload')

  function Harness() {
    const [value, setValue] = useState<StorageFileAsset[]>([{
      id: 1,
      engine: 'minio',
      bucket: 'admin-files',
      objectKey: 'uploads/mock/1-hello.txt',
      originalName: 'hello.txt',
      contentType: 'text/plain',
      size: 5,
      sha256: 'sha',
    }])
    return h(
      'div',
      null,
      h(StorageFileUpload, { value, onChange: setValue, storageApi: storageApi as any }),
      h('div', { 'data-testid': 'count' }, String(value.length)),
    )
  }

  await act(async () => {
    root.render(h(Harness))
  })
  await act(async () => {
    container.querySelector('[data-testid="remove-trigger"]')?.dispatchEvent(new MouseEvent('click', { bubbles: true }))
    await flush()
  })

  expect(container.querySelector('[data-testid="count"]')?.textContent).toBe('0')
})

test('does not upload when disabled', async () => {
  const upload = vi.fn()
  const { StorageFileUpload } = await import('./fileUpload')

  await act(async () => {
    root.render(h(StorageFileUpload, {
      value: [],
      storageApi: { upload, presigned: vi.fn(), detail: vi.fn(), del: vi.fn() } as any,
      disabled: true,
    }))
  })
  await act(async () => {
    container.querySelector('[data-testid="upload-trigger"]')?.dispatchEvent(new MouseEvent('click', { bubbles: true }))
    await flush()
  })

  expect(upload).not.toHaveBeenCalled()
})

test('opens presigned url for preview or download', async () => {
  const openWindow = vi.fn()
  const presigned = vi.fn().mockResolvedValue({
    data: {
      url: 'https://mock-minio.local/file',
      expiresAt: new Date().toISOString(),
    },
  })
  const { StorageFileUpload } = await import('./fileUpload')

  await act(async () => {
    root.render(h(StorageFileUpload, {
      value: [{
        id: 1,
        engine: 'minio',
        bucket: 'admin-files',
        objectKey: 'uploads/mock/1-hello.txt',
        originalName: 'hello.txt',
        contentType: 'text/plain',
        size: 5,
        sha256: 'sha',
      }],
      storageApi: { upload: vi.fn(), presigned, detail: vi.fn(), del: vi.fn() } as any,
      openWindow,
    }))
  })
  await act(async () => {
    container.querySelector('[data-testid="preview-trigger"]')?.dispatchEvent(new MouseEvent('click', { bubbles: true }))
    await flush()
  })

  expect(presigned).toHaveBeenCalledWith({ id: 1, disposition: 'attachment' })
  expect(openWindow).toHaveBeenCalledWith('https://mock-minio.local/file')
})
