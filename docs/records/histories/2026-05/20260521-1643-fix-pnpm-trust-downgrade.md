## [2026-05-21 16:43] | Task: fix pnpm trust downgrade

### Execution Context

- **Agent ID**: `Codex`
- **Base Model**: `GPT-5`
- **Runtime**: `Codex desktop`

### User Query

> `vp i` failed with `ERR_PNPM_TRUST_DOWNGRADE` for `@trickfilm400/rollup-plugin-off-main-thread@3.0.0-pre1`, pulled by `workbox-build@7.4.1`; user asked to fix it.

### Changes Overview

**Scope:** supply-chain configuration, lockfile, documentation

**Key Actions:**

- **[Trust policy]**: Added an exact-version `trustPolicyExclude` entry for `@trickfilm400/rollup-plugin-off-main-thread@3.0.0-pre1`.
- **[Install verification]**: Re-ran `vp i`; install completed successfully.
- **[Documentation]**: Recorded the current pnpm trust policy exceptions and rationale in supply-chain security docs.

### Design Intent (Why)

The repository intentionally uses `trustPolicy: no-downgrade`, so the fix keeps that global protection enabled and scopes the exception to the exact transitive version that blocked installation.

### Files Modified

- `pnpm-workspace.yaml`
- `pnpm-lock.yaml`
- `docs/SUPPLY_CHAIN_SECURITY.md`
