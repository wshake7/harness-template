import type { ApiResponse, AppNotifier, ResLogin } from '@vp/core'
import type { EncryptableMethod } from '../encryption/helpers'
import { HttpCode, noopNotifier, XHeader } from '@vp/core'
import { createAlova } from 'alova'
import { createClientTokenAuthentication } from 'alova/client'
import adapterFetch from 'alova/fetch'
import reactHook from 'alova/react'
import Cookies from 'js-cookie'
import NProgress from 'nprogress'

export interface RequestMethodMeta {
  meta?: {
    authRole?: string
    aesKey?: CryptoKey
  }
  url: string
}

export interface CreateVpApiClientOptions {
  baseURL?: string
  getToken: () => string
  setToken: (token: string) => void
  setPublicKey: (publicKey: string) => void
  encryptRequest: (method: EncryptableMethod) => Promise<void>
  decryptText: (encryptedText: string, aesKey: CryptoKey) => Promise<string>
  checkResponseCode: (res: ApiResponse) => Promise<void> | void
  afterLogin?: (token: string) => void
  notifier?: AppNotifier
  httpErrorMessage?: string | ((response: Response) => string)
}

export function canRequestWithoutToken(method: RequestMethodMeta) {
  const authRole = method.meta?.authRole
  return authRole === 'login'
    || authRole === 'logout'
    || authRole === 'visitor'
    || method.url === '/api/encrypt/public/key'
}

export function createVpApiClient(options: CreateVpApiClientOptions) {
  const notifier = options.notifier ?? noopNotifier

  const { onAuthRequired, onResponseRefreshToken } = createClientTokenAuthentication<typeof reactHook>({
    visitorMeta: {
      visitor: true,
    },
    assignToken(method) {
      let token = Cookies.get(XHeader.Token)
      if (token === undefined || token === '') {
        token = options.getToken()
      }
      method.config.headers[XHeader.Token] = token
    },
    async login(response, method) {
      if (response.ok) {
        const encryptedText = await response.clone().text()
        const aesKey = (method.meta as { aesKey?: CryptoKey }).aesKey
        if (!aesKey) {
          return
        }
        const decrypted = await options.decryptText(encryptedText, aesKey)
        const json = JSON.parse(decrypted)
        const res = json as ApiResponse<ResLogin>
        const data = res.data
        if (res.code === HttpCode.SUCCESS && data) {
          options.setToken(data.token)
          options.setPublicKey(data.publicKey)
          Cookies.set(XHeader.Token, data.token, {
            path: '/',
            sameSite: 'Lax',
          })
          options.afterLogin?.(data.token)
        }
      }
    },
    logout() {
    },
  })

  return createAlova({
    baseURL: options.baseURL ?? '',
    statesHook: reactHook,
    cacheFor: null,
    requestAdapter: adapterFetch(),
    shareRequest: false,
    beforeRequest: onAuthRequired(async (method) => {
      const token = Cookies.get(XHeader.Token) || options.getToken()
      if (!token && !canRequestWithoutToken(method)) {
        throw new Error('request canceled: unauthenticated')
      }
      NProgress.start()
      await options.encryptRequest(method)
    }),
    responded: onResponseRefreshToken({
      onSuccess: async (response, method) => {
        if (!response.ok) {
          const message = typeof options.httpErrorMessage === 'function'
            ? options.httpErrorMessage(response)
            : options.httpErrorMessage
          if (message) {
            notifier.error(message)
          }
          throw new Error(`[${response.status}]${response.statusText}`)
        }
        const contentType = response.headers.get('Content-Type') || ''
        if (response.headers.get(XHeader.XResponseIsEncrypt) === 'true') {
          const encryptedText = await response.clone().text()
          const aesKey = (method.meta as { aesKey?: CryptoKey }).aesKey
          if (!aesKey) {
            throw new Error('missing AES key for encrypted response')
          }
          const decrypted = await options.decryptText(encryptedText, aesKey)
          response = new Response(decrypted, {
            status: response.status,
            statusText: response.statusText,
            headers: response.headers,
          })
        }
        if (contentType.includes('application/json')) {
          const json = await response.clone().json()
          const res = json as ApiResponse
          console.log(`${method.url}:response`, res)
          if (method.url !== '/api/account/logout') {
            await options.checkResponseCode(res)
          }
          return json
        }
        return response
      },
      onError: async (error) => {
        const message = typeof options.httpErrorMessage === 'function'
          ? '请求错误'
          : options.httpErrorMessage
        if (message) {
          notifier.error(message)
        }
        console.error('[API Error]', error)
        throw error
      },
      onComplete: async () => {
        NProgress.done()
      },
    }),
  })
}
