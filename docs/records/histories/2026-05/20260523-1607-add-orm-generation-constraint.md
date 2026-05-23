## [2026-05-23 16:07] | Task: 添加 ORM model 生成约束

### 🤖 Execution Context

- **Agent ID**: `Codex`
- **Base Model**: `GPT-5`
- **Runtime**: `Codex desktop`

### 📥 User Query

> 添加一个后端生成 ORM 相关代码的约束，当 `models` 里的文件有修改时运行 Makefile 里的 `script-orm` 命令。

### 🛠 Changes Overview

**Scope:** `docs/build/backend.md`, `.harness/skills/`

**Key Actions:**

- **[后端文档]**: 在后端常用命令和变更要求中补充 `make script-orm`。
- **[Harness 检查]**: 在交付检查和代码审查 Skill 中补充 model 变更后的 ORM 生成代码检查项。

### 🧠 Design Intent (Why)

`backend/go/admin/internal/services/orm/models/` 是 ORM 查询代码生成的源头。如果 model 改动后没有重新运行生成命令，`internal/services/orm/query/` 可能与模型定义不一致，后续编译或运行时查询会出现滞后问题。

### 📁 Files Modified

- `.harness/skills/code-review.md`
- `.harness/skills/delivery-check.md`
- `docs/build/backend.md`
