# harness-template-cn

这个仓库是一个面向 Agent 协作开发的基础模板。

`AGENTS.md` 只做最薄导航，不承载完整规则。仓库内稳定知识以 `docs/README.md` 和各分层文档为准。

如果一次代码或流程变更会让某份文档过期，就在同一轮任务里顺手把它改掉。

## 每轮开始先读

- `docs/README.md`：文档总入口、目录职责和常用阅读路径。
- `docs/start/principles.md`：Agent-first、产品取舍和设计原则。
- `docs/govern/collaboration.md`：仓库级协作、提交、文档同步与测试约定。
- `docs/develop/architecture.md`：仓库整体结构和预期边界。
- `docs/govern/harness-process.md`：Harness 五阶段流程、门禁和交付闭环。
- `context/harness-framework/INDEX.md`：Agent 执行时的框架级上下文入口。

## 代码改完前要读

- `docs/govern/histories.md`：什么时候记 history、怎么命名、怎么脱敏。
- `docs/govern/quality-score.md`：当前质量分层和主要短板。

## 按任务需要选读

- `docs/govern/plans.md`：什么时候要写 execution plan，怎么维护。
- `.service-matrix/dependencies.yaml`：项目、模块、服务和路径占位符的单一真相源。
- `context/project/admin/INDEX.md`：管理后台项目前后端上下文入口。
- `.harness/commands/` 和 `.harness/skills/`：稳定命令入口和可复用工作流。
- `docs/operate/reliability.md`：运行稳定性、观测性和上线前的基本要求。
- `docs/operate/security.md`：认证、数据处理、外部集成和供应链安全默认约束。
- `docs/operate/cicd.md`：仓库的 CI/CD 骨架以及后续如何接入真实项目。
- `docs/develop/frontend.md`：如果仓库包含前端界面，这里记录对应规范。
- `docs/develop/backend.md`：如果任务涉及 Go 后端、API、配置或外部依赖，这里记录对应规范。
- `CONTRIBUTING.md`：提 PR 前后的默认检查项和协作要求。
- `docs/records/releases/README.md`：如何维护面向用户的发布记录。
- `docs/records/references/README.md`：沉淀到仓库里的外部参考资料。

## 工作规则

- 优先选择小而清晰、对仓库和 Agent 都友好的抽象。
- prompt、规则、架构约束尽量都版本化落在仓库里。
- 复杂任务不要只靠聊天上下文，应该落 execution plan。
- 完成的代码变更要记到 `docs/records/histories/`。
- 遇到可复用教训时，优先沉淀到 `context/` 或项目经验目录，不要只写进本轮对话。
