## [2026-05-24 20:06] | Task: add-test-plan-requirement

### Execution Context

- **Agent ID**: `Codex`
- **Base Model**: `GPT-5`
- **Runtime**: `Codex desktop`

### User Query

> 沉淀一下文档，让后续编写开发计划的时候记得编写测试计划

### Changes Overview

**Scope:** Harness execution plan governance docs.

**Key Actions:**

- **Plan governance**: Updated `docs/govern/plans.md` so future execution plans that touch code, config, APIs, scripts, or user-visible behavior must include a dedicated test plan.
- **Plan template**: Added a `测试计划` section to `docs/records/exec-plans/templates/execution-plan.md`, including test files, failing-test-first steps, commands, and pass criteria.

### Design Intent (Why)

Development plans should preserve not only implementation steps but also test intent. Making test planning explicit in the governance doc and template helps future Agents write executable plans that include where tests live, what behavior they cover, which commands prove failure before implementation, and what counts as passing.

### Files Modified

- `docs/govern/plans.md`
- `docs/records/exec-plans/templates/execution-plan.md`
- `docs/records/histories/2026-05/20260524-2006-add-test-plan-requirement.md`
