# 前端协作说明

前端代码位于 `front/`，使用 pnpm workspace 管理三个 React 应用和一组 `@vp/*` 共享包。根目录的 `package.json` 和 `pnpm-workspace.yaml` 是前端工具链入口。

## 工作区结构

- `front/apps/admin-react`：管理后台 SPA，包含系统、账号、日志、任务等路由。
- `front/apps/app-react`：客户端 React 应用模板。
- `front/apps/app-react-ssr`：SSR React 应用模板，接入 TanStack Start/Nitro。
- `front/packages/build-config`：共享构建配置。
- `front/packages/core`：共享类型、HTTP 常量、字典、分页、账号、通知和加密类型。
- `front/packages/react-core`：React 侧环境变量、i18n、mock 和 store 工厂。
- `front/packages/request`：alova 请求客户端、认证、加密请求和响应处理。
- `front/packages/utils`：日期、样式 class、WebCrypto 等工具函数。

## 技术栈

- Node.js 要求：`>=22.12.0`。
- 包管理器：`pnpm@10.33.0`。
- 构建工具：`vite-plus` 与 Vite，应用脚本使用 `vp dev`、`vp build`、`vp preview`。
- UI 与路由：React、TanStack Router、Ant Design、lucide-react。
- 状态、国际化和请求：zustand、i18next/react-i18next、alova。
- Mock 与测试：MSW、Playwright component/e2e 配置、Vitest/vite-plus test。

## 环境变量

环境变量由 `front/packages/react-core/src/env/createVpEnv.ts` 定义，客户端变量必须以 `VITE_` 开头。

| 变量 | 说明 |
| --- | --- |
| `VITE_API_URL` | Vite dev server 的 `/api` proxy 目标，例如 `http://127.0.0.1:3001`。 |
| `VITE_PORT` | 当前前端应用本地开发端口。 |
| `VITE_MOCK` | 字符串 `true` 时启用 MSW mock；其他值视为关闭。 |

## 常用命令

在仓库根目录运行：

```bash
pnpm install
pnpm lint
pnpm ready
```

在具体应用目录运行：

```bash
pnpm dev
pnpm build
pnpm preview
pnpm e2e:test
pnpm e2e:test-ui
```

`pnpm ready` 当前会执行 `eslint . && vp run test -r && vp run build -r`，适合作为前端完整验证入口。

## 本地联调

1. 启动后端 `admin` 服务，默认监听 `0.0.0.0:3001`。
2. 在目标前端应用中设置 `VITE_API_URL=http://127.0.0.1:3001` 和合适的 `VITE_PORT`。
3. 运行 `pnpm dev`。
4. 前端请求 `/api/*`，由 Vite proxy 转发到后端。
5. 需要脱离后端验证页面时，设置 `VITE_MOCK=true`，应用会按各自 `src/mocks` 入口启用 MSW。

## 请求与认证

- 请求客户端在 `front/packages/request/src/client/alova.ts` 中统一创建。
- Token 从 Cookie 或账号 store 读取，并写入 `@vp/core` 定义的 Token header。
- 登录、登出、访客接口和 `/api/encrypt/public/key` 可以无 Token 请求。
- 加密请求和响应解密由各应用的 `src/api/encryptRequest.ts` 与共享请求包协作完成。
- 业务响应 code 检查由应用传入的 `HttpCodeCheck` 负责，HTTP 非 2xx 会按统一错误路径提示。
- 文件上传分两类：传统 `multipart/form-data` 上传，以及 `prepareUpload -> presigned PUT -> completeUpload` 的对象存储直传。直传的 prepare/complete 请求仍走登录态，真正的文件字节通过浏览器 `fetch` PUT 到预签名 URL。
- 管理后台知识库文档页 `front/apps/admin-react/src/routes/_app/knowledge/document.tsx` 现支持“导入文件”：先调用 `StorageFileApi.uploadDirect` 上传文件，再调用 `KnowledgeDocumentApi.importFile` 触发后端建文档和向量化。

## 变更要求

- 新增跨应用能力优先放入 `front/packages/*`，避免复制到多个 app。
- 新增环境变量必须同步更新本文件，并确认 `createVpEnv` 校验规则。
- 修改路由、请求协议、mock 或构建配置时，同步更新 `docs/develop/architecture.md`、`docs/operate/reliability.md` 或 history。
- UI 验证优先使用应用内 Playwright 配置；涉及真实联调时记录后端依赖和启动顺序。
