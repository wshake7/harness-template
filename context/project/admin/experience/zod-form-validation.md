# 前端表单校验：useZodForm + ProComponents

管理后台的 Drawer 表单统一使用 `useZodForm`（`~/utils/zod`）做校验，不依赖 ProComponents 自身的 rules，保持校验逻辑可复用、可测试。

## 实现模式

### 1. 定义两层 Schema

```typescript
import z from 'zod'

// 第一层：描述表单字段类型（与 ProForm 字段一一对应）
const MyFormSchema = z.object({
  name: z.string(),
  age: z.number().optional().nullable(),
  status: z.number(),  // ProFormSelect value 用 number
  description: z.string().optional(),
})

type MyFormValues = z.infer<typeof MyFormSchema>

// 第二层：在 superRefine 中做业务校验（必填、格式、联动等）
const MySubmitSchema = MyFormSchema.superRefine((values, ctx) => {
  // 必填校验
  for (const [field, label] of [
    ['name', '名称'],
    ['status', '状态'],
  ] as const) {
    const val = values[field]
    if (val === undefined || val === null || val === '') {
      ctx.addIssue({ code: 'custom', path: [field], message: `${label}不能为空` })
    }
  }

  // 数值范围校验
  if (values.age !== undefined && values.age !== null && values.age <= 0) {
    ctx.addIssue({ code: 'custom', path: ['age'], message: '年龄必须大于 0' })
  }

  // JSON 格式校验
  if (values.jsonField?.trim()) {
    try { JSON.parse(values.jsonField.trim()) }
    catch {
      ctx.addIssue({ code: 'custom', path: ['jsonField'], message: '请输入合法 JSON' })
    }
  }
})
```

### 2. 默认值与表单类型

```typescript
const defaultFormValues: MyFormValues = {
  name: '',
  age: undefined,
  status: 1,      // ProFormSelect 用 number
  description: '',
}

// Form.useForm 带泛型，保证类型安全
const [form] = Form.useForm<MyFormValues>()
```

### 3. useZodForm 绑定

```typescript
const { rules, onFinish } = useZodForm<MyFormValues>({
  form,
  schema: MySubmitSchema,
  async onSubmit(values) {
    if (!values) return
    // 提交逻辑
  },
})
```

### 4. ProForm 组件绑定

所有需要校验的字段统一传 `rules={rules}`，由 `useZodForm` 内部通过 schema 生成校验规则：

```tsx
<ProFormText name="name" label="名称" rules={rules} />
<ProFormDigit name="age" label="年龄" rules={rules} />
<ProFormSelect name="status" label="状态" options={[...]} rules={rules} />
```

### 5. 打开表单的两种方式

**创建模式**：
```typescript
const openCreateForm = () => {
  setEditing(undefined)
  form.resetFields()
  form.setFieldsValue(defaultFormValues)
  setDrawerOpen(true)
}
```

**编辑模式**：
```typescript
const openEditForm = useCallback((record: MyModel) => {
  setEditing(record)
  form.resetFields()
  form.setFieldsValue({
    name: record.name,
    status: record.isEnabled ? 1 : 0,  // bool -> number 转换
    // ...
  })
  setDrawerOpen(true)
}, [form])
```

## 关键注意点

1. **ProFormSelect 的 value 类型**：ProFormSelect 的 `options` 中 `value` 只接受 `string | number | null`，boolean 不行。用 `1/0` 表示启用/停用，提交时 `Boolean(values.isEnabled)` 转回 boolean。

2. **编辑回填时的类型转换**：后端返回 `isEnabled: boolean`，但表单字段是 `number`，回填时要手动转换：`record.isEnabled ? 1 : 0`。

3. **schema 中的 optional 与 nullable**：`z.number().optional().nullable()` 兼容 ProFormDigit 的清空（值为 `null`），避免类型报错。

4. **superRefine 校验失败后**：`useZodForm` 内部会自动调用 `form.setFields` 把错误映射到对应字段，无需手动处理。

5. **提交前检查**：`onSubmit` 第一个参数可能是 `undefined`（校验失败时不会进入），所以开头判断 `if (!values) return`。

## 枚举字段优先用字典

后端表中 `MetricType`、`IndexType`、`ContentType`、`VectorStatus` 这类**枚举字符串**字段，前端**不要硬编码**选项，统一通过 `sys_dict` 管理：

1. 后端只存原始 string（如 `"COSINE"`、`"pending"`），不引入新的枚举类型。
2. SQL 初始化时往 `sys_dict_type` + `sys_dict_entry` 注入字典数据。命名规则：`<模块>:<字段名>`，例如 `knowledge:metric_type`、`knowledge:vector_status`。
3. `domains/dict.ts` 注册 `DictCode` 常量。
4. 页面用 `useDictMatch(DictCode.XXX)` 取字典：
   - `entries` → ProFormSelect 的 options 数据源
   - `renderLabel(value, fallback)` → 列表列的渲染
   - `getLabel(value, fallbackLabel)` → 表单选项的 label

```tsx
const metricTypeDict = useDictMatch(DictCode.KNOWLEDGE_METRIC_TYPE_DICT_CODE)

const metricTypeOptions = useMemo(() =>
  metricTypeDict.entries.map(entry => ({
    label: metricTypeDict.getLabel(entry.entryValue, entry.entryLabel),
    value: entry.entryValue,
  })),
[metricTypeDict])

// 列表渲染
{
  title: '度量类型',
  dataIndex: 'metricType',
  render: (_, record) => metricTypeDict.renderLabel(record.metricType, <Tag>{record.metricType}</Tag>),
}

// 表单字段
<ProFormSelect name="metricType" label="度量类型" options={metricTypeOptions} />
```

好处：管理员可在管理后台动态调整选项；翻译能力复用 `sys_language_entry`；不需要每次新增枚举都改前端代码。

**避免 IsEnabled + Status 双状态字段**：模型设计阶段就只保留一个开关字段，避免后续同步两份状态的复杂度。

## API 错误提示去重

后端返回业务错误时，`HttpCodeCheck`（`~/domains/http.ts`）已经调 `appNotifier.error(msg)` 弹了一次 toast；业务侧 `catch {}` 里再 `gMessage.error('保存失败')` 会导致两个提示同时出现。

解决方案：`HttpCodeCheck` 抛出的 Error 带 `notified` 标记，业务代码统一使用 `notifyError(e, fallback)` 代替 `gMessage.error`：

```ts
// domains/http.ts — 给已通知的错误打标记
import { markNotified, notifiedError } from '~/utils/notifier'

// 在 HttpCodeCheck 中：
appNotifier.error(msg)
throw notifiedError(msg)   // Error 上挂 notified=true

// utils/notifier.ts — 工具函数
export function wasNotified(err: unknown): boolean { ... }
export function notifyError(err: unknown, fallbackMessage: string): void {
  if (wasNotified(err)) return          // 底层已经弹过了，跳过
  appNotifier.error(fallbackMessage)    // 否则用降级文案提示一次
}

// 页面 catch 块（notifyError 已 auto-import）：
catch (e) {
  notifyError(e, '保存失败')
}
```

关键规则：
- 业务校验失败（如 `gMessage.error('请填写完整信息')`）保持原样，不需要 `notifyError`。
- 涉及 API 调用的 `catch {}` 统一改为 `catch (e) { notifyError(e, 'xxx失败') }`。
