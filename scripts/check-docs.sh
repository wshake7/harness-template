#!/usr/bin/env bash

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

required_files=(
  "AGENTS.md"
  "README.md"
  "CONTRIBUTING.md"
  ".service-matrix/dependencies.yaml"
  ".markdownlint-cli2.yaml"
  ".harness/commands/harness-health.md"
  ".harness/commands/requirement-gate-check.md"
  ".harness/commands/requirement-new.md"
  ".harness/agents/README.md"
  ".harness/skills/code-review.md"
  ".harness/skills/delivery-check.md"
  ".harness/skills/design-review.md"
  ".harness/skills/knowledge-maintenance.md"
  ".harness/skills/requirement-review.md"
  "context/team/INDEX.md"
  "context/harness-framework/INDEX.md"
  "context/project/INDEX.md"
  "context/project/example/experience/README.md"
  "docs/README.md"
  "docs/start/principles.md"
  "docs/develop/architecture.md"
  "docs/develop/backend.md"
  "docs/develop/frontend.md"
  "docs/govern/collaboration.md"
  "docs/govern/harness-process.md"
  "docs/govern/quality-score.md"
  "docs/govern/histories.md"
  "docs/govern/plans.md"
  "docs/operate/security.md"
  "docs/operate/reliability.md"
  "docs/operate/cicd.md"
  "docs/product/specs/index.md"
  "docs/product/specs/examples/minimal-requirement-lifecycle.md"
  "docs/records/exec-plans/README.md"
  "docs/records/exec-plans/templates/execution-plan.md"
  "docs/records/exec-plans/tech-debt-tracker.md"
  "docs/records/histories/template.md"
  "docs/records/releases/README.md"
  "docs/records/releases/feature-release-notes.md"
  "docs/records/references/README.md"
  "docs/records/generated/README.md"
  "scripts/new-experience.sh"
  "scripts/harness-sync.sh"
  "scripts/validate-service-matrix.sh"
  "scripts/check-init-project.sh"
  ".vite-hooks/pre-commit"
)

missing=0

for path in "${required_files[@]}"; do
  if [[ ! -f "${repo_root}/${path}" ]]; then
    echo "缺少必要文件: ${path}"
    missing=1
  fi
done

for dir in .harness/agents context/project/example/experience docs/records/exec-plans/active docs/records/exec-plans/completed docs/records/histories; do
  if [[ ! -d "${repo_root}/${dir}" ]]; then
    echo "缺少必要目录: ${dir}"
    missing=1
  fi
done

if ! grep -q "docs/" "${repo_root}/AGENTS.md"; then
  echo "AGENTS.md 应明确指向 docs/，说明它是仓库知识的正式来源"
  missing=1
fi

if [[ "${missing}" -ne 0 ]]; then
  exit 1
fi

echo "文档骨架检查通过"
