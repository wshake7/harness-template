import type { AccountState } from '@vp/core'
import { base64ToArrayBuffer } from '@vp/utils'
import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { immer } from 'zustand/middleware/immer'

export interface AccountActions {
  account: () => AccountState
  login: (token: string) => void
  logout: () => void
}

export interface CreateAccountStoreOptions {
  storageName?: string
  onLogout?: () => void
}

export function createAccountStore(options: CreateAccountStoreOptions = {}) {
  const useAccountStore = create<AccountState & AccountActions>()(
    persist(
      immer((set, get) => ({
        token: '',
        account() {
          const { token } = get()
          return { token }
        },
        login(token: string) {
          set((state) => {
            state.token = token
          })
        },
        logout() {
          set((state) => {
            state.token = ''
          })
          useAccountStore.persist.clearStorage()
          options.onLogout?.()
        },
      })),
      {
        name: options.storageName ?? 'account-store',
      },
    ),
  )

  return useAccountStore
}

export interface DeviceState {
  publicKey: string
}

export interface DeviceActions {
  setPublicKey: (key: string) => void
  getPublicCryptoKey: () => Promise<CryptoKey | undefined>
  clear: () => void
}

export interface CreateDeviceStoreOptions {
  storageName?: string
}

export function createDeviceStore(options: CreateDeviceStoreOptions = {}) {
  const useDeviceStore = create<DeviceState & DeviceActions>()(
    persist(
      immer((set, get) => ({
        publicKey: '',
        setPublicKey(key) {
          set((state) => {
            state.publicKey = key
          })
        },
        getPublicCryptoKey: async () => {
          const state = get()
          if (state.publicKey !== '') {
            const keyData = base64ToArrayBuffer(state.publicKey)
            const publicKey = await window.crypto.subtle.importKey(
              'spki',
              keyData,
              { name: 'RSA-OAEP', hash: 'SHA-256' },
              false,
              ['encrypt'],
            )
            return publicKey
          }
          return undefined
        },
        clear: () => {
          set((state) => {
            state.publicKey = ''
          })
          useDeviceStore.persist.clearStorage()
        },
      })),
      {
        name: options.storageName ?? 'device-store',
      },
    ),
  )

  return useDeviceStore
}

export interface MockStore {
  account: {
    username: string
    pwd: string
    token: string
  }
  updatePwd: (pwd: string) => void
  _hasHydrated: boolean
  setHasHydrated: (state: boolean) => void
}

export interface CreateMockStoreOptions {
  storageName?: string
  account?: MockStore['account']
}

export function createMockStore(options: CreateMockStoreOptions = {}) {
  const useMockStore = create<MockStore>()(
    persist(
      immer(set => ({
        account: options.account ?? {
          username: 'admin',
          pwd: '123456',
          token: 'mock_token',
        },
        updatePwd: (pwd: string) =>
          set((state) => {
            state.account.pwd = pwd
          }),
        _hasHydrated: false,
        setHasHydrated: (state: boolean) =>
          set((s) => {
            s._hasHydrated = state
          }),
      })),
      {
        name: options.storageName ?? 'mock-store',
        onRehydrateStorage: () => (state) => {
          state?.setHasHydrated(true)
        },
      },
    ),
  )

  const waitForHydration = (): Promise<void> => {
    return new Promise((resolve) => {
      if (useMockStore.getState()._hasHydrated) {
        resolve()
        return
      }
      const unsub = useMockStore.subscribe((state) => {
        if (state._hasHydrated) {
          unsub()
          resolve()
        }
      })
    })
  }

  return {
    useMockStore,
    waitForHydration,
  }
}
