import { LockOutlined, MobileOutlined, UserOutlined } from '@ant-design/icons'
import {
  LoginForm,
  ProConfigProvider,
  ProFormCaptcha,
  ProFormCheckbox,
  ProFormText,
} from '@ant-design/pro-components'
import { createFileRoute } from '@tanstack/react-router'
import { App, Form, Tabs, theme } from 'antd'
import { useState } from 'react'

export const Route = createFileRoute('/(login)/login')({
  component: RouteComponent,
})

type LoginType = 'phone' | 'account'

const PHONE_REGEX = /^1\d{10}$/

const PasswordStatus = ({ value = '' }: { value?: string }) => {
  const { token } = theme.useToken()
  const getStatus = () => {
    if (value && value.length > 12) {
      return 'ok'
    }
    if (value && value.length > 6) {
      return 'pass'
    }
    return 'poor'
  }
  const status = getStatus()
  if (status === 'pass') {
    return <div style={{ color: token.colorWarning }}>强度：中</div>
  }
  if (status === 'ok') {
    return <div style={{ color: token.colorSuccess }}>强度：强</div>
  }
  return <div style={{ color: token.colorError }}>强度：弱</div>
}

function RouteComponent() {
  const { token } = theme.useToken()
  const [loginType, setLoginType] = useState<LoginType>('account')
  const [form] = Form.useForm()
  const { message } = App.useApp()

  return (
    <div className="flex justify-center items-start min-h-screen pt-16 sm:pt-24 md:pt-32 ">
      <ProConfigProvider hashed={false}>
        <div style={{ backgroundColor: token.colorBgContainer }}>
          <LoginForm
            form={form}
            name="login"
            title="Wshake"
            onFinish={async (values: { username: string, pwd: string }) => {
              if (loginType === 'account') {
                await AccountApi.loginPwd({ username: values.username, pwd: values.pwd })
              }
            }}
          >
            <Tabs
              centered
              activeKey={loginType}
              onChange={activeKey => setLoginType(activeKey as LoginType)}
              items={[
                {
                  key: 'account',
                  label: '账号密码登录',
                },
                {
                  key: 'phone',
                  label: '手机号登录',
                },
              ]}
            />
            {loginType === 'account' && (
              <>
                <ProFormText
                  name="username"
                  fieldProps={{
                    size: 'large',
                    prefix: <UserOutlined className="prefixIcon" />,
                    autoComplete: 'username',
                  }}
                  placeholder="用户名: admin or user"
                  rules={[
                    {
                      required: true,
                      message: '请输入用户名!',
                    },
                  ]}
                />
                <ProFormText.Password
                  name="pwd"
                  fieldProps={{
                    type: 'password',
                    size: 'large',
                    prefix: <LockOutlined className="prefixIcon" />,
                    autoComplete: 'current-password',
                  }}
                  placeholder="密码: ant.design"
                  rules={[
                    {
                      required: true,
                      message: '请输入密码！',
                    },
                  ]}
                />
                <div style={{ marginTop: -20, marginBottom: 24 }}>
                  <Form.Item noStyle shouldUpdate={(prev, curr) => prev.pwd !== curr.pwd}>
                    {() => <PasswordStatus value={form.getFieldValue('pwd')} />}
                  </Form.Item>
                </div>
              </>
            )}
            {loginType === 'phone' && (
              <>
                <ProFormText
                  fieldProps={{
                    size: 'large',
                    prefix: <MobileOutlined className="prefixIcon" />,
                  }}
                  name="mobile"
                  placeholder="手机号"
                  rules={[
                    {
                      required: true,
                      message: '请输入手机号！',
                    },
                    {
                      pattern: PHONE_REGEX,
                      message: '手机号格式错误！',
                    },
                  ]}
                />
                <ProFormCaptcha
                  fieldProps={{
                    size: 'large',
                    prefix: <LockOutlined className="prefixIcon" />,
                  }}
                  captchaProps={{
                    size: 'large',
                  }}
                  placeholder="请输入验证码"
                  captchaTextRender={(timing, count) => {
                    if (timing) {
                      return `${count} ${'获取验证码'}`
                    }
                    return '获取验证码'
                  }}
                  name="captcha"
                  rules={[
                    {
                      required: true,
                      message: '请输入验证码！',
                    },
                  ]}
                  onGetCaptcha={async () => {
                    message.success('获取验证码成功!')
                  }}
                />
              </>
            )}
            <div
              style={{
                marginBlockEnd: 24,
              }}
            >
              <ProFormCheckbox noStyle>记住密码</ProFormCheckbox>
              <a
                style={{
                  float: 'right',
                }}
              >
                忘记密码
              </a>
            </div>
          </LoginForm>
        </div>
      </ProConfigProvider>
    </div>
  )
}
