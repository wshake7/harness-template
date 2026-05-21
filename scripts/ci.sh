#!/usr/bin/env bash

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

"${repo_root}/scripts/check-docs.sh"
"${repo_root}/scripts/check-repo-hygiene.sh"
"${repo_root}/scripts/check-action-pinning.sh"
"${repo_root}/scripts/validate-service-matrix.sh"
"${repo_root}/scripts/check-init-project.sh"

while IFS= read -r file; do
  bash -n "$file"
done < <(find "${repo_root}/scripts" -type f -name '*.sh' | sort)

# 前端测试
echo "Running frontend tests..."
cd "${repo_root}/front/apps/admin-react"

echo "  -> Unit tests"
node_modules/.bin/vp test --run

echo "  -> Component tests"
node_modules/.bin/playwright test -c playwright-ct.config.ts

echo "  -> E2E tests"
node_modules/.bin/playwright test -c playwright.config.ts

echo "基础 CI 检查通过"
