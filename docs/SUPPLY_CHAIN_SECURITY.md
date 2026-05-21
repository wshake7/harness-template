# 供应链安全

这份文档定义模板默认采用的供应链安全做法。

## 默认控制项

- 在 Pull Request 上做依赖变更审查。
- 在 PR、定时任务和手动触发时，用 OSV 对仓库中的依赖声明和 lockfile 做漏洞扫描。
- 为 release 产物生成 SBOM。
- 为 release 产物生成 build provenance attestation。
- 所有 GitHub Actions 都固定到不可变的 commit SHA，而不是漂移的版本标签。
- pnpm 保持 `trustPolicy: no-downgrade`；如果遇到已确认的传递依赖信任降级误报，必须用精确版本写入 `trustPolicyExclude`，不要关闭全局策略。

## 当前对应关系

- `actions/dependency-review-action`：阻止 PR 引入高风险依赖变更。
- `google/osv-scanner-action`：根据仓库里的依赖文件扫描已知漏洞。
- `anchore/sbom-action`：生成 SPDX 格式的 SBOM。
- `actions/attest-build-provenance`：为 release artifact 生成签名 provenance。
- `scripts/check-action-pinning.sh`：如果 workflow 里出现浮动 tag 而不是 SHA，直接让 CI 失败。
- `pnpm-workspace.yaml`：记录 pnpm trust policy 和精确版本例外。

## pnpm trust policy 例外

`trustPolicyExclude` 只能使用精确包名和版本，新增例外时要记录原因，不能扩大到包级或关闭 `trustPolicy: no-downgrade`。

| 包 | 原因 |
| --- | --- |
| `semver@6.3.1` | 由 `@babel/core` 经 `eslint-plugin-react-hooks` 传递引入；仓库保留全局 no-downgrade，只放行该精确旧版本。 |
| `@trickfilm400/rollup-plugin-off-main-thread@3.0.0-pre1` | `workbox-build@7.4.1` 传递引入；该 prerelease 版本缺少 provenance attestation，触发 `ERR_PNPM_TRUST_DOWNGRADE`，当前只放行报错精确版本以恢复安装。 |

## 限制和前提

- Dependency Review 在 public repo 可以直接使用；private repo 通常需要 GitHub Advanced Security 或对应的代码安全能力。
- OSV 和 SBOM 的效果依赖仓库里存在可识别的依赖清单或 lockfile。
- 只有当 `scripts/release-package.sh` 真的代表项目的构建产物时，provenance 才真正有意义。
- OpenSSF Scorecard 默认不启用，因为新模板仓库还没有真实分支保护、release 历史和 SAST 姿态可以评分；等仓库规则配置完成后再按需加回。

## 项目落地后建议继续做的事

- 锁定并提交项目真实依赖的 lockfile。
- 定期复查 `trustPolicyExclude`，依赖上游恢复 provenance 或替换传递依赖后及时删除例外。
- 让构建过程尽量可重复、可验证。
- 如果条件允许，在部署链路里增加对 provenance 的校验。
- 把 attestation 校验继续下沉到部署平台或准入层。
