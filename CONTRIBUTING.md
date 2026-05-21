# 参与协作

这个仓库是为 Agent-first 前后端协作准备的。人和 Agent 都应把稳定知识落到版本化文件里，而不是只保存在聊天记录、口头同步或工单评论里。

## 基本协作方式

- 从 `AGENTS.md` 开始，再按任务类型读取对应文档。
- 涉及架构、路径或服务归属时，先看 `.service-matrix/dependencies.yaml`。
- 涉及后端时，读 `docs/BACKEND.md`；涉及前端时，读 `docs/FRONTEND.md`。
- 行为变更要同步更新代码、文档、测试、服务矩阵和 history。
- 跨端、跨模块、高风险或会分多轮推进的任务，先在 `docs/exec-plans/active/` 下建 execution plan。

## 本地验证

仓库基础门禁：

```bash
make ci
```

前端验证：

```bash
pnpm lint
pnpm ready
cd front/apps/admin-react && pnpm e2e:test
```

后端验证：

```bash
cd backend/go
go test ./...
```

如果某个项目级命令因为本地外部依赖缺失失败，要在 PR 或 history 中说明失败命令、原因和下一步。

## 发起 Pull Request 之前

- 运行 `make ci`。
- 代码或流程变更要补齐 `docs/histories/`。
- 用户可感知变更要补 release note。
- 服务、路径、环境变量或启动方式变化时，同步更新 README、架构、前端/后端、稳定性或安全文档。
- 确认示例、脚本、说明文档和当前实现一致。

## Review 默认要求

- 优先拆成范围清晰的小 PR。
- 明确写出风险点、迁移影响和后续待办。
- 引用具体文档、plan、spec 或 history，不依赖评审者自己猜上下文。
- 对认证、权限、加密、外部依赖和生成代码改动保持额外审查。
