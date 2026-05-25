import { beforeEach, expect, test, vi } from 'vite-plus/test'

const send = vi.fn()
const post = vi.fn(() => ({ send }))
const fetchMock = vi.fn()

function getPostCall(index: number): [string, FormData | Record<string, unknown>, { cacheFor?: number, meta?: Record<string, unknown> }] {
  return post.mock.calls[index] as unknown as [string, FormData | Record<string, unknown>, { cacheFor?: number, meta?: Record<string, unknown> }]
}

vi.mock('../index', () => ({
  default: {
    Post: post,
  },
}))

beforeEach(() => {
  post.mockClear()
  send.mockReset()
  send.mockResolvedValue(undefined)
  fetchMock.mockReset()
  vi.stubGlobal('fetch', fetchMock)
})

test('upload sends FormData with skipEncrypt meta', async () => {
  const { StorageFileApi } = await import('./storageFile')
  const file = new File(['hello'], 'hello.txt', { type: 'text/plain' })

  await StorageFileApi.upload(file, {
    bizType: 'attachment',
    bizID: 'biz-1',
    metadata: '{"source":"test"}',
  })

  expect(post).toHaveBeenCalledTimes(1)
  const [url, formData, options] = getPostCall(0)
  expect(url).toBe('/api/storage/file/upload')
  expect(formData).toBeInstanceOf(FormData)
  const multipart = formData as FormData
  expect(multipart.get('file')).toBe(file)
  expect(multipart.get('bizType')).toBe('attachment')
  expect(multipart.get('bizID')).toBe('biz-1')
  expect(multipart.get('metadata')).toBe('{"source":"test"}')
  expect(options.meta).toEqual({ skipEncrypt: true })
})

test('presigned and delete call expected endpoints', async () => {
  const { StorageFileApi } = await import('./storageFile')

  await StorageFileApi.presigned({ id: 1, disposition: 'attachment' })
  await StorageFileApi.del({ id: 1 })

  expect(post).toHaveBeenNthCalledWith(1, '/api/storage/file/presigned', { id: 1, disposition: 'attachment' }, { cacheFor: 0, meta: { skipEncrypt: true } })
  expect(post).toHaveBeenNthCalledWith(2, '/api/storage/file/del', { id: 1 }, { cacheFor: 0, meta: { skipEncrypt: true } })
})

test('prepare and complete upload call expected endpoints', async () => {
  const { StorageFileApi } = await import('./storageFile')

  await StorageFileApi.prepareUpload({
    originalName: 'runbook.md',
    contentType: 'text/markdown',
    size: 12,
    bizType: 'knowledge-document',
    bizID: 'collection-1',
  })
  await StorageFileApi.completeUpload({ id: 10 })

  expect(post).toHaveBeenNthCalledWith(1, '/api/storage/file/prepareUpload', {
    originalName: 'runbook.md',
    contentType: 'text/markdown',
    size: 12,
    bizType: 'knowledge-document',
    bizID: 'collection-1',
  }, { cacheFor: 0, meta: { skipEncrypt: true } })
  expect(post).toHaveBeenNthCalledWith(2, '/api/storage/file/completeUpload', { id: 10 }, { cacheFor: 0, meta: { skipEncrypt: true } })
})

test('uploadDirect performs prepare, put, and complete', async () => {
  const { StorageFileApi } = await import('./storageFile')
  const file = new File(['hello'], 'runbook.md', { type: 'text/markdown' })

  send
    .mockResolvedValueOnce({
      data: {
        asset: { id: 10 },
        uploadURL: 'https://mock-minio.local/upload',
        method: 'PUT',
        headers: { 'Content-Type': 'text/markdown' },
      },
    })
    .mockResolvedValueOnce({
      data: {
        id: 10,
        status: 'active',
      },
    })
  fetchMock.mockResolvedValue({ ok: true })

  const result = await StorageFileApi.uploadDirect(file, {
    bizType: 'knowledge-document',
    bizID: 'collection-1',
  })

  expect(fetchMock).toHaveBeenCalledWith('https://mock-minio.local/upload', {
    method: 'PUT',
    headers: { 'Content-Type': 'text/markdown' },
    body: file,
  })
  const asset = result.data
  expect(asset).toBeDefined()
  expect(asset?.id).toBe(10)
})
