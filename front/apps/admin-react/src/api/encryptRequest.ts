import { createEncryptedRequestHelpers } from '@vp/request'
import { appNotifier } from '~/utils/notifier'

const encryptedRequestHelpers = createEncryptedRequestHelpers({
  getCachedPublicKey: () => useDeviceStore.getState().publicKey,
  setPublicKey: publicKey => useDeviceStore.getState().setPublicKey(publicKey),
  getPublicCryptoKey: () => useDeviceStore.getState().getPublicCryptoKey(),
  notifier: appNotifier,
  systemErrorMessage: '系统异常',
})

export const ensurePublicKey = encryptedRequestHelpers.ensurePublicKey
export const encryptRequest = encryptedRequestHelpers.encryptRequest
export const createEncryptedRequestConfig = encryptedRequestHelpers.createEncryptedRequestConfig
export const decryptText = encryptedRequestHelpers.decryptText
