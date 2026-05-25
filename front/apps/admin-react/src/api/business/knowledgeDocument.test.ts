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

test('importFile calls expected endpoint', async () => {
  const { KnowledgeDocumentApi } = await import('./knowledgeDocument')

  await KnowledgeDocumentApi.importFile({
    collectionID: 1,
    fileAssetID: 10,
    title: 'runbook',
    contentType: 'markdown',
  })

  expect(post).toHaveBeenCalledWith('/api/knowledge/document/importFile', {
    collectionID: 1,
    fileAssetID: 10,
    title: 'runbook',
    contentType: 'markdown',
  }, { cacheFor: 0 })
})
