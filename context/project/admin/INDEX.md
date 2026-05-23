# 管理后台项目上下文

这个入口面向 Agent 执行任务时快速建立项目感。更详细的人类说明以 `docs/` 为准。

## 代码范围

- 后端服务：`backend/go/admin`
- 后端共享模块：`backend/go/go-common`、`backend/go/orm-crud`、`backend/go/sa-token`
- 管理后台前端：`front/apps/admin-react`
- 前端共享包：`front/packages/*`
- 其他前端模板：`front/apps/app-react`、`front/apps/app-react-ssr`

## 默认运行关系

- 后端 `admin` 默认监听 `0.0.0.0:3001`，API 前缀为 `/api`。
- 前端通过 `VITE_API_URL` 配置 Vite proxy 目标，通过 `VITE_MOCK=true` 启用 MSW。
- 后端默认依赖 Postgres、Redis 和 Temporal，本地配置在 `backend/go/admin/etc/config.yaml`。
- 前端请求层统一处理 Token、加密请求、公钥、响应解密和业务 code 检查。

## 任务前必读

- 架构或跨端联调：`docs/build/architecture.md`
- 后端服务、配置或 API：`docs/build/backend.md`
- 前端应用、环境变量或 UI 验证：`docs/build/frontend.md`
- 安全相关：`docs/operate/security.md`
- 稳定性或排障：`docs/operate/reliability.md`
- 服务路径归属：`.service-matrix/dependencies.yaml`

## 项目经验

- 前端表单校验：`context/project/admin/experience/zod-form-validation.md`
- 登录互踢兼容性：`context/project/admin/experience/auth-kickout-rueidis.md`

## Agent 注意事项

- 不要把 `backend/go/admin/internal` 的内部实现当作其他服务可直接依赖的公共 API。
- 前端跨应用能力优先放入 `front/packages/*`，避免在多个 app 中复制。
- 配置、路由、认证、加密、外部依赖或验证命令发生变化时，同步更新 `docs/` 和 history。
