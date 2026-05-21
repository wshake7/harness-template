import type { FormInstance, FormListFieldData } from 'antd'
import type { Rule } from 'antd/es/form'
import type { NamePath } from 'antd/es/form/interface'
import type z from 'zod'
import { Form } from 'antd'
import { useCallback, useMemo } from 'react'
import { appNotifier } from '~/utils/notifier'

interface UseFormProps<T extends object> {
  form?: FormInstance<T>
  schema?: z.ZodObject
  onSubmit: (data?: T, error?: FormListFieldData) => Promise<boolean> | Promise<void> | boolean | void
  zodValidator?: <T extends Record<string, any>>(props: ZodValidatorProps<T>) => Rule
}

interface UseFormTypes<T> {
  form: FormInstance<T>
  onFinish?: (formData: T) => Promise<boolean | void>
  rules: Rule[]
}

const mapErrorFromZodIssue = (issues: z.core.$ZodIssue[]) =>
  issues.reduce((obj: Record<string, string[]>, issue) => {
    const fieldName = issue.path.join('.')
    if (!obj[fieldName]) {
      obj[fieldName] = [issue.message]
    }
    else {
      obj[fieldName] = [...obj[fieldName], issue.message]
    }
    return obj
  }, {})

const getZodIssues = (error: unknown) => {
  const issues = (error as { issues?: unknown })?.issues
  return Array.isArray(issues) ? issues as z.core.$ZodIssue[] : []
}

const getIssueFieldName = (issue: z.core.$ZodIssue) => issue.path.join('.') || '表单'

const getIssueMessage = (issue: z.core.$ZodIssue) => {
  if (issue.message && !issue.message.startsWith('Invalid input')) {
    return issue.message
  }
  return `${getIssueFieldName(issue)} 校验失败`
}

const getFirstIssueMessage = (issues: z.core.$ZodIssue[]) => {
  const issue = issues[0]
  return issue ? getIssueMessage(issue) : '请检查表单信息'
}

const setZodFormErrors = <T extends Record<string, any>>(
  form: FormInstance<T>,
  issues: z.core.$ZodIssue[],
) => {
  const errorMap = mapErrorFromZodIssue(issues)
  const values = form.getFieldsValue()
  const fieldNames = new Set([
    ...Object.keys(values),
    ...Object.keys(errorMap),
  ])

  form.setFields(
    Array.from(fieldNames, key => ({
      name: key as NamePath,
      errors: errorMap[key] ?? [],
    })),
  )
}

interface ZodValidatorProps<T extends Record<string, any>> {
  form: FormInstance<T>
  schema: z.ZodObject
}

export const fieldZodValidator = <T extends object>(props: ZodValidatorProps<T>): Rule => ({
  async validator(rule: any) {
    const values = props.form.getFieldsValue()
    await props.schema.parseAsync(values).catch((e: z.ZodError) => {
      const issues = getZodIssues(e)
      if (!issues.length) {
        throw e
      }
      const errorMap = mapErrorFromZodIssue(issues)
      const currentFieldError = errorMap[rule.field]
      if (currentFieldError?.length) {
        throw new Error(currentFieldError[0])
      }
    })
  },
})

export const globalZodValidator = <T extends Record<string, any>>(
  props: ZodValidatorProps<T>,
): Rule => ({
  async validator(rule: any) {
    const values = props.form.getFieldsValue()
    const result = await props.schema.safeParseAsync(values)
    const issues = result.success ? [] : result.error.issues
    const errorMap = mapErrorFromZodIssue(issues)
    setZodFormErrors(props.form, issues)
    const currentFieldError = errorMap[rule.field]
    if (currentFieldError?.length) {
      throw new Error(currentFieldError[0])
    }
  },
})

export const defaultZodValidator = globalZodValidator

export function useZodForm<T extends Record<string, any>>({
  form: propsForm,
  schema,
  onSubmit,
  zodValidator = defaultZodValidator,
}: UseFormProps<T>): UseFormTypes<T> {
  const [form] = Form.useForm<T>(propsForm)
  const rules = useMemo(() => (schema ? [zodValidator({ form, schema })] : []), [schema, form, zodValidator])

  const onFinish = useCallback(
    async (formData: T) => {
      if (!schema) {
        return onSubmit(formData)
      }

      try {
        await schema.parseAsync(formData)
        return await onSubmit(formData)
      }
      catch (e: unknown) {
        const issues = getZodIssues(e)
        if (!issues.length) {
          throw e
        }
        setZodFormErrors(form, issues)
        appNotifier.error(getFirstIssueMessage(issues))
        return false // 校验失败 → 弹窗不关闭
      }
    },
    [onSubmit, schema, form],
  )

  return {
    form,
    onFinish,
    rules,
  }
}
