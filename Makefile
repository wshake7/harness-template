PROJECT ?=
SLUG ?=
SCOPE ?=

# 定义伪目标（不生成同名文件）
.PHONY: init check-docs check-repo ci release-package new-history new-plan new-experience validate-matrix harness-sync test-admin test-admin-unit test-admin-ct test-admin-e2e test-admin-report

# 初始化新项目
init:
	@if [ -z "$(PROJECT)" ]; then echo "用法: make init PROJECT=项目名"; exit 1; fi
	./scripts/init-project.sh "$(PROJECT)"

# 检查文档完整性
check-docs:
	./scripts/check-docs.sh

# 检查仓库卫生（文档 + 代码规范）
check-repo:
	./scripts/check-docs.sh
	./scripts/check-repo-hygiene.sh

# 运行完整 CI 检查（含前端测试）
ci:
	./scripts/ci.sh

# 验证服务矩阵配置
validate-matrix:
	./scripts/validate-service-matrix.sh

# 同步 Harness 流程变更到文档
harness-sync:
	./scripts/harness-sync.sh --workspace

# 发布新包版本
release-package:
	./scripts/release-package.sh

# 创建变更历史记录
new-history:
	@if [ -z "$(SLUG)" ]; then echo "用法: make new-history SLUG=变更名"; exit 1; fi
	./scripts/new-history.sh "$(SLUG)"

# 创建执行计划
new-plan:
	@if [ -z "$(SLUG)" ]; then echo "用法: make new-plan SLUG=计划名"; exit 1; fi
	./scripts/new-exec-plan.sh "$(SLUG)"

# 创建经验沉淀文档
new-experience:
	@if [ -z "$(SCOPE)" ]; then echo "用法: make new-experience SCOPE=project/example SLUG=经验名"; exit 1; fi
	@if [ -z "$(SLUG)" ]; then echo "用法: make new-experience SCOPE=project/example SLUG=经验名"; exit 1; fi
	./scripts/new-experience.sh "$(SCOPE)" "$(SLUG)"

# 运行 admin-react 单元测试（Vitest，21 个测试）
test-admin-unit:
	@echo "→ Running admin-react unit tests (Vitest)..."
	@cd front/apps/admin-react && node_modules/.bin/vp test --run

# 运行 admin-react 组件测试（Playwright CT，10 个测试）
test-admin-ct:
	@echo "→ Running admin-react component tests (Playwright CT)..."
	@cd front/apps/admin-react && node_modules/.bin/playwright test -c playwright-ct.config.ts

# 运行 admin-react E2E 测试（Playwright，15 个测试）
test-admin-e2e:
	@echo "→ Running admin-react E2E tests (Playwright)..."
	@cd front/apps/admin-react && node_modules/.bin/playwright test -c playwright.config.ts --workers=1

# 运行 admin-react 全部测试（单元 + CT + E2E，共 46 个）
test-admin: test-admin-unit test-admin-ct test-admin-e2e
	@echo "→ All admin-react tests passed."

# 打开 Playwright HTML 测试报告
# 报告路径：front/apps/admin-react/playwright-report/index.html
test-admin-report:
	@echo "→ Opening Playwright test report..."
	@cd front/apps/admin-react && npx playwright show-report
