## [2026-05-23 16:01] | Task: 重构项目文档结构

### 🤖 Execution Context

- **Agent ID**: `Codex`
- **Base Model**: `GPT-5`
- **Runtime**: `Codex desktop`

### 📥 User Query

> 按工作流分层彻底重构项目文档结构，保持 `AGENTS.md` 为薄入口，新增 `docs/README.md`，迁移并合并现有文档，同步脚本、上下文和 Harness 入口。

### 🛠 Changes Overview

**Scope:** `docs/`, `context/`, `.harness/`, `scripts/`, repository entry docs

**Key Actions:**

- **[重构文档目录]**: 将横铺文档迁移到 `docs/start/`、`docs/build/`、`docs/operate/`、`docs/govern/`、`docs/product/` 和 `docs/records/`。
- **[合并入口文档]**: 新增 `docs/README.md`，并将核心理念、产品判断和设计原则合并到 `docs/start/principles.md`。
- **[合并安全文档]**: 将应用安全和供应链安全统一到 `docs/operate/security.md`。
- **[同步工具链]**: 更新 `make new-plan`、`make new-history`、`make check-docs` 和 `make harness-sync` 对应的新路径。

### 🧠 Design Intent (Why)

原文档入口一次性暴露太多横铺文件，Agent 和人都需要在多个根级文档之间跳转。按工作流分层后，`AGENTS.md` 继续保持轻量，`docs/README.md` 承担完整地图，日常任务可以按阶段和职责找到对应知识源。

### 📁 Files Modified

- `AGENTS.md`
- `docs/README.md`
- `docs/start/principles.md`
- `docs/build/`
- `docs/operate/`
- `docs/govern/`
- `docs/product/specs/`
- `docs/records/`
- `context/`
- `.harness/`
- `scripts/`
