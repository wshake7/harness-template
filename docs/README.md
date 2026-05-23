# 文档总入口

`docs/` 是这个仓库面向人、Agent 和评审的正式知识源。`AGENTS.md` 只负责把每轮任务路由到这里；具体规则、流程、架构和记录都在本目录维护。

## 目录职责

- `start/`：最先阅读的原则和判断标准。
- `build/`：架构、后端、前端和开发协作说明。
- `operate/`：稳定性、CI/CD、安全和供应链治理。
- `govern/`：协作方式、Harness 流程、计划、history 和质量评分。
- `product/specs/`：需求、目标、非目标和验收标准。
- `records/`：execution plan、history、release note、外部参考和生成文档。

## 常用阅读路径

- 开始任意任务：读本文件、`start/principles.md`、`govern/collaboration.md`、`build/architecture.md`、`govern/harness-process.md`。
- 后端任务：再读 `build/backend.md`、`operate/reliability.md`，涉及认证或外部凭据时读 `operate/security.md`。
- 前端任务：再读 `build/frontend.md`、`operate/reliability.md`，涉及请求协议或 Token 时读 `operate/security.md`。
- 流程或脚本任务：再读 `govern/harness-process.md`、`operate/cicd.md`、`govern/histories.md`。
- 高风险或跨模块任务：先读 `govern/plans.md`，在 `records/exec-plans/active/` 建立或更新计划。
- 交付前：读 `govern/histories.md`、`govern/quality-score.md`，确认是否需要 history、release note 或质量评分更新。

## 当前入口

- `start/principles.md`：Agent-first、产品取舍和设计原则。
- `build/architecture.md`：仓库拓扑、依赖方向和边界。
- `build/backend.md`：Go 后端启动、配置、外部依赖、API 和验证。
- `build/frontend.md`：React/pnpm workspace、环境变量、联调和验证。
- `operate/reliability.md`：启动、外部依赖、可观测性和排障。
- `operate/cicd.md`：仓库级 CI/CD、前端测试和 release 制品。
- `operate/security.md`：认证、权限、加密、密钥、依赖、SBOM 和 provenance。
- `govern/collaboration.md`：协作、文档同步、Git、评审和验证。
- `govern/harness-process.md`：Harness 五阶段流程和轻量门禁。
- `govern/plans.md`：execution plan 使用规则。
- `govern/histories.md`：history 记录规则。
- `govern/quality-score.md`：当前质量水位和技术债关系。

## 记录目录

- `records/exec-plans/active/`：进行中的 execution plan。
- `records/exec-plans/completed/`：已完成的 execution plan，正文保留当时路径和上下文。
- `records/histories/`：已完成变更的审计记录，正文保留当时路径和上下文。
- `records/releases/`：面向用户的发布记录。
- `records/references/`：沉淀到仓库里的外部参考资料。
- `records/generated/`：自动生成文档的入口说明。

## 维护规则

- 新增稳定规则时，优先放进最贴近职责的分层文档，不要塞回 `AGENTS.md`。
- 行为、脚本、目录、接口或配置发生变化时，同步更新相关文档和必要的 `context/` 入口。
- 已完成记录属于审计材料，不为追求路径统一而批量改写旧正文。
