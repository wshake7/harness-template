import type { ReqChangePwd, ReqPwdLogin } from '@vp/core'
import Cookies from 'js-cookie'
import { XHeader } from '~/domains/http'
import { router } from '~/router'
import API from '../index'

async function loginPwd(req: ReqPwdLogin) {
  await API.Post<Res>('/api/account/login/pwd', req, {
    cacheFor: 0,
    meta: {
      authRole: 'login',
    },
  }).send()
}

function logout() {
  useAccountStore.getState().logout()
  useDeviceStore.getState().setPublicKey('')
  Cookies.remove(XHeader.Token)
  router.update({
    context: {
      account: {
        token: '',
      },
    },
  })
  API.Get<Res>('/api/account/logout', {
    cacheFor: 0,
    meta: {
      authRole: 'logout',
    },
  }).send()
  router.navigate({ to: '/login' })
}

async function changePwd(req: ReqChangePwd) {
  await API.Post<Res>('/api/account/changePwd', req, {
    cacheFor: 0,
  }).send()
}

export const AccountApi = {
  loginPwd,
  logout,
  changePwd,
}
