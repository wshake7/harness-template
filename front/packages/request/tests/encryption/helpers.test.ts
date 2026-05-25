import { expect, test, vi } from 'vite-plus/test'

const XHeader = {
  XRequestTimestamp: 'X-Request-Timestamp',
  XRequestID: 'X-Request-Id',
  XRequestEncryptedKey: 'X-Request-Encrypted-Key',
  XRequestSignature: 'X-Request-Signature',
} as const

vi.mock('@vp/core', () => ({
  XHeader,
  noopNotifier: {
    error: vi.fn(),
    success: vi.fn(),
    warning: vi.fn(),
    info: vi.fn(),
  },
}))

vi.mock('@vp/utils', () => ({
  aesDecrypt: vi.fn(),
  aesEncrypt: vi.fn(),
  generateAesKey: vi.fn(),
  rsaEncrypt: vi.fn(),
  uriSort: vi.fn(),
}))

test('keeps FormData body when skipEncrypt is true', async () => {
  const { createEncryptedRequestHelpers } = await import('../../src/encryption/helpers')
  const helpers = createEncryptedRequestHelpers({
    getCachedPublicKey: () => '',
    setPublicKey: () => {},
    getPublicCryptoKey: vi.fn(),
    fetchPublicKey: vi.fn(),
  })

  const formData = new FormData()
  formData.append('file', new Blob(['hello'], { type: 'text/plain' }), 'hello.txt')
  const method = {
    url: '/api/storage/file/upload',
    data: formData,
    config: {},
    meta: {
      skipEncrypt: true,
    },
  }

  await helpers.encryptRequest(method)

  expect(method.data).toBe(formData)
  expect(method.config.headers?.[XHeader.XRequestTimestamp]).toBeDefined()
  expect(method.config.headers?.[XHeader.XRequestID]).toBeDefined()
  expect(method.config.headers?.[XHeader.XRequestEncryptedKey]).toBeUndefined()
})
