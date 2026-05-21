import type { FormInstance } from 'antd'
import type { FormListFieldData } from 'antd/lib/form'
import type { ReactNode } from 'react'
import { ModalForm, ProFormText } from '@ant-design/pro-components'
import { Space } from 'antd'
import z from 'zod'
import { useZodForm } from '~/utils/zod'

interface Props {
  children: ReactNode
  onSubmit: (values?: ChangePwdFormValues, error?: FormListFieldData) => Promise<boolean> | boolean | void
  username: string
  form?: FormInstance<ChangePwdFormValues>
}

const ChangePwdSchema = z
  .object({
    oldPwd: z.string('请输入旧密码').min(6, '旧密码长度不能小于6位'),
    newPwd: z.string('请输入新密码').min(6, '新密码长度不能小于6位'),
    confirmPwd: z.string('请输入确认新密码').optional(),
  })
  .superRefine((data, ctx) => {
    if (data.newPwd === data.oldPwd) {
      ctx.addIssue({
        code: 'custom',
        message: '不能与旧密码一致',
        path: ['newPwd'],
      })
    }
    if (data.confirmPwd !== data.newPwd) {
      ctx.addIssue({
        code: 'custom',
        message: '必须与新密码一致',
        path: ['confirmPwd'],
      })
    }
  })

export type ChangePwdFormValues = z.infer<typeof ChangePwdSchema>

export default (props: Props) => {
  const { form, rules, onFinish } = useZodForm({
    form: props.form,
    schema: ChangePwdSchema,
    onSubmit(values, error) {
      return props.onSubmit(values, error)
    },
  })
  return (
    <Space>
      <ModalForm
        form={form}
        onFinish={onFinish}
        title="修改密码"
        width={400}
        onOpenChange={(open) => {
          if (!open) {
            form.resetFields()
          }
        }}
        trigger={<div>{props.children}</div>}
        submitTimeout={2000}
      >
        <ProFormText
          width="md"
          name="username"
          label="用户名"
          fieldProps={{
            autoComplete: 'username',
            defaultValue: props.username,
          }}
          disabled
          placeholder="请输入用户名"
        />

        <ProFormText.Password
          required
          width="md"
          name="oldPwd"
          label="旧密码"
          fieldProps={{
            autoComplete: 'current-password',
          }}
          rules={rules}
          placeholder="请输入旧密码"
        />

        <ProFormText.Password
          required
          width="md"
          name="newPwd"
          label="新密码"
          fieldProps={{
            autoComplete: 'new-password',
          }}
          rules={rules}
          placeholder="请输入新密码"
        />

        <ProFormText.Password
          required
          width="md"
          name="confirmPwd"
          label="确认新密码"
          fieldProps={{
            autoComplete: 'new-password',
          }}
          rules={rules}
          placeholder="请确认新密码"
        />
      </ModalForm>
    </Space>
  )
}
