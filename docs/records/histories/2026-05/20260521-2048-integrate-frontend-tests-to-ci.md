## [2026-05-21 20:48] | Task: 前端测试体系接入 CI

### Execution Context

- **Agent ID**: `Claude`
- **Base Model**: `Opus 4.7`
- **Runtime**: `Claude Code VSCode`

### User Query

> make harness-sync 同步一下相关文档

### Changes Overview

**Scope:** 前端 admin-react 测试体系、CI 流程、Makefile

**Key Actions:**

- **Makefile 新增测试目标**: 添加 `test-admin-unit`、`test-admin-ct`、`test-admin-e2e`、`test-admin`、`test-admin-report` 五个目标，覆盖 Vitest 单元测试（21 个）、Playwright 组件测试（10 个）和 E2E 测试（15 个）。
- **CI 接入前端测试**: `scripts/ci.sh` 中新增前端测试执行步骤，使 `make ci` 自动运行 admin-react 的全部 46 个测试。
- **同步 CI 文档**: 更新 `docs/CICD.md`，反映前端测试已正式进入 CI 流程，并调整推荐接入顺序。

### Design Intent

让前端验证从"手动运行"变为"CI 自动执行"，确保每次提交都能通过单元、组件和 E2E 三层测试。同时通过 Makefile 提供本地快速运行单个测试层的能力，保持开发体验。

### Files Modified

- `Makefile`
- `scripts/ci.sh`
- `docs/CICD.md`
- `front/apps/admin-react/package.json`
- `front/apps/admin-react/playwright-ct.config.ts`
- `front/apps/admin-react/playwright.config.ts`
- `front/apps/admin-react/vite.config.ts`
- `front/apps/admin-react/src/components/error.spec.tsx`
- `front/apps/admin-react/src/components/notFound.spec.tsx`
- `front/apps/admin-react/src/stores/menuTabs.test.ts`
- `front/apps/admin-react/src/stores/resourceMenu.test.ts`
- `front/apps/admin-react/src/utils/antIcons.test.ts`
- `front/apps/admin-react/src/utils/zod.test.ts`
- `front/apps/admin-react/src/utils/zod.ts`
- `front/apps/admin-react/tests/auth/login.spec.ts`
- `front/apps/admin-react/tests/business/dashboard.spec.ts`
- `front/apps/admin-react/tests/business/dict.spec.ts`
- `front/apps/admin-react/tests/business/menu.spec.ts`
- `front/apps/admin-react/tests/business/role.spec.ts`
- `front/apps/admin-react/tests/business/user.spec.ts`
- `pnpm-workspace.yaml`
- `vite.config.ts`
