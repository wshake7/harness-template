# admin-react 前端验收清单

## 一、测试覆盖率

### 1.1 单元测试（Vitest）

| 模块 | 文件 | 覆盖项 | 状态 |
|------|------|--------|------|
| Zod 工具 | `src/utils/zod.test.ts` | mapErrorFromZodIssue, getZodIssues, getFirstIssueMessage, field/global ZodValidator | 已完成 |
| 图标工具 | `src/utils/antIcons.test.ts` | getAntIconStyle, antIconNamesByStyle 分类正确性 | 已完成 |
| 标签页 Store | `src/stores/menuTabs.test.ts` | add, remove, removeAll 行为 | 已完成 |
| 资源菜单 Store | `src/stores/resourceMenu.test.ts` | setDynamicMenuTree, clear 行为 | 已完成 |

### 1.2 组件测试（Playwright CT）

| 组件 | 文件 | 覆盖项 | 状态 |
|------|------|--------|------|
| ErrorComponent | `src/components/error.spec.tsx` | 渲染、按钮存在、点击回调 | 已完成 |
| NotFoundComponent | `src/components/notFound.spec.tsx` | 渲染、按钮存在、点击回调 | 已完成 |
| JsonCodeBlock | `src/components/business/logger/jsonCodeBlock.spec.tsx` | 空值、有效 JSON、无效 JSON、XSS 防护 | 已有 |

### 1.3 E2E 流程测试（Playwright）

| 流程 | 文件 | 覆盖项 | 状态 |
|------|------|--------|------|
| 登录 | `tests/auth/login.spec.ts` | 成功登录、失败处理、标签切换、未登录重定向 | 已完成 |
| Dashboard | `tests/business/dashboard.spec.ts` | 欢迎信息、菜单显示、导航跳转 | 已完成 |
| 用户管理 | `tests/business/user.spec.ts` | 列表加载、创建抽屉、表单验证 | 已完成 |
| 字典管理 | `tests/business/dict.spec.ts` | 类型列表显示 | 已完成 |
| 角色管理 | `tests/business/role.spec.ts` | 角色列表、创建抽屉 | 已完成 |
| 菜单管理 | `tests/business/menu.spec.ts` | 菜单列表、创建抽屉 | 已完成 |

## 二、本地测试命令

```bash
# 单元测试
cd front/apps/admin-react
pnpm test

# 组件测试
pnpm ct

# E2E 测试
pnpm e2e

# 全量测试
pnpm test && pnpm ct && pnpm e2e
```

## 三、自动化门禁

### 3.1 pre-commit
- [x] ESLint 自动修复

### 3.2 CI
- [x] `scripts/ci.sh` 包含前端单元测试、组件测试和 E2E 测试

### 3.3 PR Check
- 所有单元测试通过
- 所有组件测试通过
- 所有 E2E 测试通过

## 四、已知限制

1. **Monaco Editor**：DictEntryPanel 中的 MonacoCodeEditorBlock 在 CT 测试中难以覆盖，暂不处理
2. **AntIconPicker**：依赖 ResizeObserver 和大量 DOM 操作，CT 测试较复杂，后续迭代
3. **文件上传/下载**：当前业务不涉及，暂不覆盖
4. **移动端适配**：当前只测试 Desktop Chrome
