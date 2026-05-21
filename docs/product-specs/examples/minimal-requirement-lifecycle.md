# 最小需求生命周期示例

这份示例展示一个需求如何从定义走到交付，供 Agent 和工程师对齐口径。

## 背景

模板需要能证明 Harness 基础闭环可运行，而不只是文档存在。

## 目标

- 新项目 clone 后可以运行 `make ci`。
- Agent 能从 `AGENTS.md` 找到流程、上下文和门禁入口。
- 服务矩阵、Skill、Command 和经验目录都有最小结构。

## 非目标

- 不绑定具体业务技术栈。
- 不实现多运行时渲染脚本。
- 不内置复杂 Agent 调度系统。

## 验收标准

- `make ci` 通过。
- `docs/HARNESS_PROCESS.md` 存在并定义阶段和门禁。
- `context/`、`.harness/`、`.service-matrix/` 的入口文件存在。
- 交付时有 history 记录。
