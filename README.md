# harness-template-cn

面向 Agent 协作开发的前后端一体仓库。它在 Harness 治理层之上接入了 Go 后端、React 前端、服务矩阵、上下文入口、文档同步和基础 CI。

## 快速开始

先运行仓库级检查：

```bash
make ci
```

后端本地启动：

```bash
cd backend/go/admin
make run
```

前端本地启动：

```bash
pnpm install
cd front/apps/admin-react
VITE_PORT=5173 VITE_API_URL=http://127.0.0.1:3001 VITE_MOCK=false pnpm dev
```

需要只看前端页面、不联调后端时，把 `VITE_MOCK` 设为 `true`。

## 当前项目结构

- `backend/go/admin`：Go Fiber 管理后台 API 服务，默认端口 `3001`，API 前缀 `/api`。
- `backend/go/go-common`：后端共享工具库。
- `backend/go/orm-crud`：ORM CRUD、分页 proto 和 GORM 查询辅助能力。
- `front/apps/admin-react`：管理后台 React SPA。
- `front/apps/app-react`：客户端 React 应用模板。
- `front/apps/app-react-ssr`：SSR React 应用模板。
- `front/packages/*`：`@vp/*` 前端共享包。
- `docs/`：面向人和评审的说明、指南、质量记录和发布记录。
- `context/`：面向 Agent 执行的稳定上下文入口。
- `.service-matrix/`：项目、模块、服务和路径占位符的单一真相源。

## 常用命令

仓库级：

```bash
make check-docs
make validate-matrix
make ci
make harness-sync
make new-plan SLUG=example-plan
make new-history SLUG=example-change
```

前端：

```bash
pnpm lint
pnpm ready
cd front/apps/admin-react && pnpm e2e:test
```

后端：

```bash
cd backend/go
go test ./...
cd admin && make swagger
```

## 文档入口

- `docs/README.md`：完整文档地图、目录职责和常用阅读路径。
- `docs/start/principles.md`：Agent-first、产品取舍和设计原则。
- `docs/develop/architecture.md`：整体拓扑、依赖边界和前后端数据流。
- `docs/develop/backend.md`：Go 后端启动、配置、依赖、API 和验证方式。
- `docs/develop/frontend.md`：React/pnpm workspace、环境变量、联调、请求和验证方式。
- `docs/operate/reliability.md`：启动、外部依赖、可观测性和排障入口。
- `docs/operate/security.md`：认证、权限、加密、配置和依赖安全约束。
- `docs/govern/harness-process.md`：Harness 五阶段流程和轻量门禁。
- `context/project/admin/INDEX.md`：Agent 执行管理后台任务时的项目级入口。

## Harness 闭环

每轮协作从 `AGENTS.md` 开始。长期稳定的仓库知识在 `docs/` 和 `context/` 中维护。

默认五阶段流程：

1. 初始化：建立项目、上下文和矩阵骨架。
2. 需求定义：写清目标、非目标和验收标准。
3. 设计：记录方案、影响面、风险和验证方式。
4. 实现：按任务修改代码、文档和测试。
5. 交付：运行 `make ci`，补 history 或 release note，沉淀可复用经验。

流程与门禁的单一真相源是 `docs/govern/harness-process.md`。
