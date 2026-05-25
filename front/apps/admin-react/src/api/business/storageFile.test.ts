import { beforeEach, expect, test, vi } from 'vite-plus/test'

const send = vi.fn()
const post = vi.fn(() => ({ send }))

vi.mock('../index', () => ({
  default: {
    Post: post,
  },
}))

beforeEach(() => {
  post.mockClear()
  send.mockReset()
  send.mockResolvedValue(undefined)
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
  const [url, formData, options] = post.mock.calls[0]
  expect(url).toBe('/api/storage/file/upload')
  expect(formData).toBeInstanceOf(FormData)
  expect(formData.get('file')).toBe(file)
  expect(formData.get('bizType')).toBe('attachment')
  expect(formData.get('bizID')).toBe('biz-1')
  expect(formData.get('metadata')).toBe('{"source":"test"}')
  expect(options.meta).toEqual({ skipEncrypt: true })
})

test('presigned and delete call expected endpoints', async () => {
  const { StorageFileApi } = await import('./storageFile')

  await StorageFileApi.presigned({ id: 1, disposition: 'attachment' })
  await StorageFileApi.del({ id: 1 })

  expect(post).toHaveBeenNthCalledWith(1, '/api/storage/file/presigned', { id: 1, disposition: 'attachment' }, { cacheFor: 0 })
  expect(post).toHaveBeenNthCalledWith(2, '/api/storage/file/del', { id: 1 }, { cacheFor: 0 })
})
