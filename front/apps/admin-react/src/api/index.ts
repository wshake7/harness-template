import { createVpApiClient } from '@vp/request'
import { decryptText, encryptRequest } from '~/api/encryptRequest'
import { gEnv } from '~/env'
import { router } from '~/router'
import { appNotifier } from '~/utils/notifier'

const API = createVpApiClient({
  baseURL: gEnv.VITE_MOCK ? '' : '',
  getToken: () => useAccountStore.getState().token,
  setToken: token => useAccountStore.getState().login(token),
  setPublicKey: publicKey => useDeviceStore.getState().setPublicKey(publicKey),
  encryptRequest,
  decryptText,
  checkResponseCode: HttpCodeCheck,
  notifier: appNotifier,
  httpErrorMessage: '请求错误',
  afterLogin: (token) => {
    router.update({
      context: {
        account: {
          token,
        },
      },
    })
    router.navigate({ to: '/' })
  },
})

export default API
