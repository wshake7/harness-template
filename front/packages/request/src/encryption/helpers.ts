import type { ApiResponse, AppNotifier } from '@vp/core'
import { noopNotifier, XHeader } from '@vp/core'
import { aesDecrypt, aesEncrypt, generateAesKey, rsaEncrypt, uriSort } from '@vp/utils'

export interface EncryptableMethod {
  url: string
  data?: any
  config: {
    headers?: Record<string, any>
    params?: Record<string, any> | string | URLSearchParams
  }
  meta?: Record<string, any>
}

export interface EncryptedRequestConfig {
  headers: Record<string, any>
  data?: any
  aesKey?: CryptoKey
  nonce: string
}

export interface CreateEncryptedRequestHelpersOptions {
  getCachedPublicKey: () => string
  setPublicKey: (publicKey: string) => void
  getPublicCryptoKey: () => Promise<CryptoKey | undefined>
  fetchPublicKey?: () => Promise<string>
  notifier?: AppNotifier
  systemErrorMessage?: string
}

function normalizeParams(params: Record<string, any> | string | URLSearchParams | undefined): Record<string, any> {
  if (!params) { return {} }

  if (typeof params === 'string') {
    return Object.fromEntries(new URLSearchParams(params))
  }

  if (params instanceof URLSearchParams) {
    return Object.fromEntries(params.entries())
  }

  if (typeof params === 'object') {
    return params as Record<string, any>
  }

  return {}
}

async function defaultFetchPublicKey() {
  const response = await fetch('/api/encrypt/public/key')
  if (!response.ok) {
    return ''
  }
  const res = await response.json() as ApiResponse<{ publicKey: string }>
  return res.data?.publicKey || ''
}

export function createEncryptedRequestHelpers(options: CreateEncryptedRequestHelpersOptions) {
  const notifier = options.notifier ?? noopNotifier
  const systemErrorMessage = options.systemErrorMessage ?? '系统异常'

  async function ensurePublicKey() {
    let publicKey = options.getCachedPublicKey()
    if (publicKey !== '') {
      return publicKey
    }

    publicKey = await (options.fetchPublicKey ?? defaultFetchPublicKey)()
    if (publicKey !== '') {
      options.setPublicKey(publicKey)
    }
    return publicKey
  }

  async function encryptRequest(method: EncryptableMethod) {
    if (method.url === '/api/encrypt/public/key') {
      const timestamp = Date.now()
      const nonce = Math.random().toString(36).substring(2, 18)
      method.config.headers = method.config.headers || {}
      method.config.headers[XHeader.XRequestTimestamp] = timestamp
      method.config.headers[XHeader.XRequestID] = nonce
      return
    }

    const encryptedConfig = await createEncryptedRequestConfig({
      data: method.data,
      headers: method.config.headers,
      params: method.config.params,
    })
    method.config.headers = encryptedConfig.headers
    method.data = encryptedConfig.data
    method.meta = {
      ...method.meta,
      aesKey: encryptedConfig.aesKey,
      nonce: encryptedConfig.nonce,
    }
  }

  async function createEncryptedRequestConfig(optionsConfig: {
    data?: any
    headers?: Record<string, any>
    params?: Record<string, any> | string | URLSearchParams
  } = {}): Promise<EncryptedRequestConfig> {
    const headers = {
      ...(optionsConfig.headers || {}),
    }
    const timestamp = Date.now()
    const nonce = Math.random().toString(36).substring(2, 18)
    headers[XHeader.XRequestTimestamp] = timestamp
    headers[XHeader.XRequestID] = nonce

    const publicKey = await ensurePublicKey()
    if (publicKey === '') {
      notifier.error(systemErrorMessage)
      return { headers, nonce }
    }

    const publicCryptoKey = await options.getPublicCryptoKey()
    if (!publicCryptoKey) {
      notifier.error(systemErrorMessage)
      return { headers, nonce }
    }

    const { key, keyBase64 } = await generateAesKey()
    headers[XHeader.XRequestEncryptedKey] = await rsaEncrypt(keyBase64, publicCryptoKey)

    const sort = uriSort({
      [XHeader.XRequestTimestamp]: timestamp,
      [XHeader.XRequestID]: nonce,
      ...normalizeParams(optionsConfig.params),
    })
    const aesData = await aesEncrypt(key, sort, optionsConfig.data)
    headers[XHeader.XRequestSignature] = aesData.TagIv

    return {
      headers,
      data: aesData.Ciphertext !== '' ? aesData.Ciphertext : optionsConfig.data,
      aesKey: key,
      nonce,
    }
  }

  async function decryptText(encryptedText: string, aesKey: CryptoKey) {
    return aesDecrypt(encryptedText, aesKey, '')
  }

  return {
    ensurePublicKey,
    encryptRequest,
    createEncryptedRequestConfig,
    decryptText,
  }
}
