import type { ReqChangePwd, ReqPwdLogin } from '@vp/core'

import { http, HttpResponse } from 'msw'
import { useMockStore } from '~/stores/mock'
import { fail, success, url } from '.'

export const accountHandlers = [
  http.post(url('/api/account/login/pwd'), async ({ request }) => {
    const { username, pwd } = (await request.json()) as ReqPwdLogin
    const { account } = useMockStore.getState()
    if (username !== account.username || pwd !== account.pwd) {
      return HttpResponse.json(fail('用户名或密码错误'))
    }
    return HttpResponse.json(success({
      token: account.token,
    }))
  }),

  http.get(url('/api/account/logout'), async () => {
    return HttpResponse.json(success(undefined))
  }),

  http.post(url('/api/account/changePwd'), async ({ request }) => {
    const { oldPwd, newPwd } = (await request.json()) as ReqChangePwd
    const { account, updatePwd } = useMockStore.getState()

    if (oldPwd !== account.pwd) {
      return HttpResponse.json(fail('原密码错误'))
    }

    updatePwd(newPwd)

    return HttpResponse.json(success(undefined))
  }),
]
