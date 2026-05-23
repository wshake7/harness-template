## [2026-05-23 16:12] | Task: 添加 ORM model 生成约束

### 🤖 Execution Context

- **Agent ID**: `Codex`
- **Base Model**: `GPT-5`
- **Runtime**: `Codex desktop`

### 📥 User Query

> 把后端 ORM model 变更后需要运行 `make script-orm` 的约束同步到该分支。

### 🛠 Changes Overview

**Scope:** `.harness/skills/`, `docs/`

**Key Actions:**

- **[协作文档]**: 在仓库协作验证规则中补充 ORM model 变更后的生成命令要求。
- **[Harness 检查]**: 在 code review 和 delivery check 中补充 `make script-orm` 检查项。

### 🧠 Design Intent (Why)

ORM model 是查询代码生成的源头。把生成命令写入稳定协作入口，可以降低 model 变更后漏提交 `internal/services/orm/query/` 产物的风险。

### 📁 Files Modified

- `.harness/skills/code-review.md`
- `.harness/skills/delivery-check.md`
- `docs/REPO_COLLAB_GUIDE.md`
