# CI/CD 说明

这个仓库保留 Harness 模板自带的基础 CI/CD 骨架，并已接入前后端项目文档。当前 CI 覆盖仓库治理、供应链安全和前端测试验证。

## 默认包含的内容

- `ci.yml`：仓库级检查，覆盖 docs、repo hygiene、Markdown、shell 脚本和服务矩阵校验。
- `.vite-hooks/pre-commit`：提交前运行 `scripts/harness-sync.sh --staged` 和 `vp staged`，提前暴露文档同步风险。
- `supply-chain-security.yml`：在 PR 上做依赖变更检查，并在 PR、定时任务和手动触发时运行 OSV 扫描。
- `release.yml`：手动触发的 release 流水线，用来打包仓库级制品、生成 provenance，并创建 GitHub Release。

## 当前仓库级门禁

`make ci` 当前执行：

```bash
scripts/check-docs.sh
scripts/check-repo-hygiene.sh
scripts/check-action-pinning.sh
scripts/validate-service-matrix.sh
scripts/check-init-project.sh
bash -n scripts/*.sh

# 前端测试
cd front/apps/admin-react
node_modules/.bin/vp test --run
node_modules/.bin/playwright test -c playwright-ct.config.ts
node_modules/.bin/playwright test -c playwright.config.ts
```

也就是说，`make ci` 能验证文档骨架、服务矩阵、基础仓库卫生、脚本语法，并自动运行前端 admin-react 的单元测试、组件测试和 E2E 测试。

## 项目级验证建议

前端：

```bash
pnpm lint
pnpm ready
cd front/apps/admin-react && pnpm e2e:test
```

后端：

```bash
cd backend/go
go test ./...
```

跨端联调：

```bash
cd backend/go/admin && make run
cd front/apps/admin-react && VITE_PORT=5173 VITE_API_URL=http://127.0.0.1:3001 VITE_MOCK=false pnpm dev
```

## 推荐接入顺序

1. 保留 `ci.yml`，作为唯一默认常驻的仓库基础门禁。
2. 确认前端 `pnpm ready` 在当前锁文件和 Node 版本下稳定。
3. 确认后端 `cd backend/go && go test ./...` 不依赖未文档化的本地服务。
4. ~~将稳定后的前后端验证命令接入 `scripts/ci.sh`。~~ 前端测试已接入；后端测试稳定后接入。
5. 用真实前端/后端构建产物替换 `scripts/release-package.sh` 的占位打包逻辑。
6. 技术栈和环境稳定后，再补具体部署 job。

## 默认 release 产物

当前 release 流水线会产出：

- `release-manifest.json`
- `repo-metadata.tgz`
- `sbom.spdx.json`
- 对 release artifact 生成的 GitHub artifact attestation

后续如果前端或后端进入真实发布流程，应继续保留 SBOM 和 provenance，不要另起一套绕过现有供应链能力的临时发布链路。
