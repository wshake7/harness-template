# Admin Appsvc Rename Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rename `backend/go/admin/internal/service` to `backend/go/admin/internal/appsvc` so application-facing service ports are clearly separated from infrastructure lifecycle services in `internal/services`.

**Architecture:** Keep the current code structure and behavior intact, only clarify the package boundary name. `internal/appsvc` will continue to provide router/logic-facing interfaces plus thin default implementations, while `internal/services` remains the startup and infrastructure lifecycle layer.

**Tech Stack:** Go, Go Fiber, gomock, repo docs/history conventions

---

### Task 1: Rename the application service package

**Files:**
- Modify: `backend/go/admin/internal/service/*`
- Modify: `backend/go/admin/internal/router/**/*.go`
- Modify: `backend/go/admin/internal/mock/*.go`

- [ ] **Step 1: Move `internal/service` to `internal/appsvc`**

Run: `mv backend/go/admin/internal/service backend/go/admin/internal/appsvc`
Expected: package files now live under `internal/appsvc`

- [ ] **Step 2: Update package declarations and imports**

Update `package service` to `package appsvc` in moved files, and replace imports from `admin/internal/service` to `admin/internal/appsvc` throughout the admin module.

- [ ] **Step 3: Verify mock generation paths still point to `internal/mock`**

Check the moved `//go:generate mockgen` directives and keep destinations unchanged so existing mocks remain in `backend/go/admin/internal/mock`.

### Task 2: Sync repository documentation

**Files:**
- Modify: `docs/BACKEND.md`
- Modify: `docs/ARCHITECTURE.md`

- [ ] **Step 1: Update backend boundary wording**

Document that `internal/appsvc` is the application-facing port/thin-adapter layer and `internal/services` is the infrastructure lifecycle layer.

- [ ] **Step 2: Keep wording aligned across docs**

Make sure architecture and backend docs use the same naming and boundary description.

### Task 3: Record the change and verify it

**Files:**
- Create: `docs/histories/2026-05/20260522-<time>-rename-admin-appsvc.md`

- [ ] **Step 1: Add a history entry**

Record the naming cleanup, affected packages, and why the rename improves boundary clarity.

- [ ] **Step 2: Run focused verification**

Run: `cd backend/go/admin && go test ./...`
Expected: tests pass and imports resolve after the rename

- [ ] **Step 3: Review working tree**

Run: `git status --short`
Expected: only the intended rename, doc, and history changes appear besides pre-existing unrelated worktree changes
