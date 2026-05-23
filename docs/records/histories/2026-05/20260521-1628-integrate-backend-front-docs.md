## [2026-05-21 16:28] | Task: integrate backend and front docs

### Execution Context

- **Agent ID**: `Codex`
- **Base Model**: `GPT-5`
- **Runtime**: `Codex desktop`

### User Query

> 用户已将之前代码里的 backend 和 front 加入仓库，要求按计划补齐所有相关文档与知识源。

### Changes Overview

**Scope:** documentation, service matrix, project context

**Key Actions:**

- **Architecture docs**: Replaced template architecture notes with the real Go backend and React frontend topology.
- **Service matrix**: Replaced the example service with real admin/backend/frontend/shared module entries.
- **Frontend docs**: Documented pnpm workspace, Vite/vite-plus scripts, env vars, mock mode, request/auth flow, and validation commands.
- **Backend docs**: Added Go workspace, admin service config, dependencies, API conventions, Swagger, and test commands.
- **Governance docs**: Updated README, contributing, CI/CD, reliability, security, quality score, debt tracker, and project context.

### Design Intent (Why)

The repository now contains real backend and frontend code, so keeping template placeholder documentation would force future humans and Agents to rediscover basic architecture, startup, dependency, and validation facts from source code. This update moves those facts into versioned docs and the service matrix.

### Files Modified

- `README.md`
- `CONTRIBUTING.md`
- `AGENTS.md`
- `.service-matrix/dependencies.yaml`
- `docs/ARCHITECTURE.md`
- `docs/BACKEND.md`
- `docs/FRONTEND.md`
- `docs/RELIABILITY.md`
- `docs/SECURITY.md`
- `docs/CICD.md`
- `docs/QUALITY_SCORE.md`
- `docs/exec-plans/tech-debt-tracker.md`
- `context/project/INDEX.md`
- `context/project/admin/INDEX.md`
