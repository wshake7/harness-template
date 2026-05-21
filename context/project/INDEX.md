# 项目级上下文

项目级上下文记录具体项目、模块和服务的知识。当前仓库已经接入前后端代码，Agent 开始项目相关任务时应优先读取这里，再按任务选择更细文档。

## 当前项目入口

- `context/project/admin/INDEX.md`：管理后台项目前后端上下文入口。
- `.service-matrix/dependencies.yaml`：模块、服务和路径归属的单一真相源。
- `docs/ARCHITECTURE.md`：整体拓扑、依赖边界和前后端数据流。
- `docs/BACKEND.md`：Go 后端启动、配置、依赖、API 和验证方式。
- `docs/FRONTEND.md`：React/pnpm workspace、环境变量、联调、请求和验证方式。

## 维护规则

- 如果新增真实项目或服务，先更新 `.service-matrix/dependencies.yaml`，再新增对应 `context/project/<name>/INDEX.md`。
- 如果某类经验会反复影响 Agent 执行，沉淀到 `context/project/<name>/experience/`，不要只写在 history。
- 模板示例仍保留在 `context/project/example/`，仅作为目录结构参考。
