import { beforeEach, expect, test, vi } from 'vite-plus/test'

const importKey = vi.fn()
const storage = new Map<string, string>()

vi.mock('@vp/core', () => ({}))

vi.mock('@vp/utils', () => ({
  base64ToArrayBuffer: vi.fn(() => new Uint8Array([1, 2, 3]).buffer),
}))

beforeEach(() => {
  storage.clear()
  const localStorageMock = {
    getItem: (key: string) => storage.get(key) ?? null,
    setItem: (key: string, value: string) => {
      storage.set(key, value)
    },
    removeItem: (key: string) => {
      storage.delete(key)
    },
    clear: () => {
      storage.clear()
    },
  }
  Object.defineProperty(globalThis, 'localStorage', {
    value: localStorageMock,
    configurable: true,
  })
  importKey.mockReset()
  const cryptoMock = {
    subtle: {
      importKey,
    },
  }
  Object.defineProperty(globalThis, 'crypto', {
    value: cryptoMock,
    configurable: true,
  })
  Object.defineProperty(globalThis, 'window', {
    value: {
      localStorage: localStorageMock,
      crypto: cryptoMock,
    },
    configurable: true,
  })
})

test('setPublicKey drops invalid public key format', async () => {
  const { createDeviceStore } = await import('../../src/stores/factories')
  const useDeviceStore = createDeviceStore({ storageName: 'device-store-invalid-format' })

  useDeviceStore.getState().setPublicKey(`-----BEGIN PUBLIC KEY-----
bad-key
-----END PUBLIC KEY-----`)

  expect(useDeviceStore.getState().publicKey).toBe('')
})

test('getPublicCryptoKey clears invalid cached public key', async () => {
  const { createDeviceStore } = await import('../../src/stores/factories')
  const useDeviceStore = createDeviceStore({ storageName: 'device-store-legacy-invalid-key' })

  useDeviceStore.setState({
    publicKey: `-----BEGIN PUBLIC KEY-----
bad-key
-----END PUBLIC KEY-----`,
  })

  await expect(useDeviceStore.getState().getPublicCryptoKey()).resolves.toBeUndefined()
  expect(useDeviceStore.getState().publicKey).toBe('')
  expect(importKey).not.toHaveBeenCalled()
})
