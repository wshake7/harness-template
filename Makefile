PROJECT ?=
SLUG ?=
SCOPE ?=

.PHONY: init check-docs check-repo ci release-package new-history new-plan new-experience validate-matrix harness-sync

init:
	@if [ -z "$(PROJECT)" ]; then echo "用法: make init PROJECT=项目名"; exit 1; fi
	./scripts/init-project.sh "$(PROJECT)"

check-docs:
	./scripts/check-docs.sh

check-repo:
	./scripts/check-docs.sh
	./scripts/check-repo-hygiene.sh

ci:
	./scripts/ci.sh

validate-matrix:
	./scripts/validate-service-matrix.sh

harness-sync:
	./scripts/harness-sync.sh --workspace

release-package:
	./scripts/release-package.sh

new-history:
	@if [ -z "$(SLUG)" ]; then echo "用法: make new-history SLUG=变更名"; exit 1; fi
	./scripts/new-history.sh "$(SLUG)"

new-plan:
	@if [ -z "$(SLUG)" ]; then echo "用法: make new-plan SLUG=计划名"; exit 1; fi
	./scripts/new-exec-plan.sh "$(SLUG)"

new-experience:
	@if [ -z "$(SCOPE)" ]; then echo "用法: make new-experience SCOPE=project/example SLUG=经验名"; exit 1; fi
	@if [ -z "$(SLUG)" ]; then echo "用法: make new-experience SCOPE=project/example SLUG=经验名"; exit 1; fi
	./scripts/new-experience.sh "$(SCOPE)" "$(SLUG)"
