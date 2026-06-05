# Go Admin Schema Baseline

> 基线冻结时间：2026-05-26
> 基于：`backend/go/admin/internal/services/orm/models/` 所有 Go model 文件
> 对应 Flyway 迁移：`V1__baseline_schema.sql`

## 表清单与功能说明

| 表名 | 功能 | 关键特征 |
|------|------|----------|
| `sys_user` | 系统用户 | 用户名唯一索引（active），密码 json:"-" |
| `sys_role` | 系统角色 | 角色标识唯一索引（active），树形结构（parent_id） |
| `sys_user_role` | 用户-角色关联 | user_id + role_id 唯一索引（active），软删除 |
| `sys_resource_menu` | 菜单资源 | 树形结构（parent_id, tree_path），多类型（CATALOG/MENU/BUTTON/EMBEDDED/LINK） |
| `sys_resource_api` | API 资源 | module + path + method 唯一索引（active） |
| `sys_resource_menu_api` | 菜单-API 关联 | menu_id + api_id 唯一索引（active） |
| `sys_role_menu` | 角色-菜单授权 | role_id + menu_id 唯一索引（active） |
| `sys_role_api` | 角色-API 授权 | role_id + api_id 唯一索引（active） |
| `sys_data_permission` | 数据权限 | 主体类型+ID+资源表+操作唯一索引，支持多种 scope |
| `sys_dict_type` | 字典类型 | type_code 唯一索引（active） |
| `sys_dict_entry` | 字典条目 | sys_dict_type_id + language_code + entry_value 唯一索引（active） |
| `sys_language_type` | 语言类型 | type_code 唯一索引（active），最多一个 is_default=true（active） |
| `sys_language_entry` | 语言条目 | entry_code + sys_language_type_id 唯一索引（active） |
| `sys_api_log` | API 操作日志 | 多维度索引，记录请求/响应/变更/耗时 |
| `sys_login_log` | 登录日志 | 记录登录成功/失败、IP、浏览器信息 |
| `sys_casbin_model` | Casbin 模型定义 | name 唯一索引（active） |
| `casbin_rule` | Casbin 策略规则 | 由 jCasbin 框架管理，不在 Go model 中定义 |
| `job_schedule` | 任务调度配置 | job_code 唯一索引，支持 ONCE/CRON/INTERVAL |
| `job_execution` | 任务执行记录 | workflow_id + run_id 唯一索引，状态机 RUNNING/SUCCESS/FAILED/CANCELED/TIMEOUT |
| `file_asset` | 文件资产 | object_key 唯一索引（active），状态机 pending_upload/active/deleted |
| `knowledge_collection` | 知识库集合 | collection_name 唯一索引（active），关联 Milvus Collection |
| `knowledge_document` | 知识库文档 | document_id 唯一索引（active），vector_status 状态机 pending/indexed/failed |
| `agent_session` | AI Agent 会话 | session_id 唯一索引（active），状态 active/archived |
| `agent_message` | AI Agent 消息 | 按 session_id 分组，角色 user/assistant/system/tool |

## 通用字段约定

所有实体表（除 `sys_api_log` 外）统一使用以下 mixin 字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | `BIGSERIAL PRIMARY KEY` | 自增主键 |
| `created_at` | `TIMESTAMPTZ` | 创建时间（Go 侧 Hook 自动设置） |
| `updated_at` | `TIMESTAMPTZ` | 更新时间（Go 侧 Hook 自动设置） |
| `created_by` | `BIGINT NOT NULL DEFAULT 0` | 创建者 ID |
| `updated_by` | `BIGINT NOT NULL DEFAULT 0` | 更新者 ID |
| `deleted_by` | `BIGINT NOT NULL DEFAULT 0` | 删除者 ID |
| `deleted_at` | `BIGINT NOT NULL DEFAULT 0` | 软删除标记（毫秒时间戳，0=未删除） |

部分表额外使用：

| 字段 | 类型 | 说明 |
|------|------|------|
| `remark` | `TEXT DEFAULT ''` | 备注 |
| `is_enabled` | `BOOLEAN DEFAULT TRUE` | 启用状态 |
| `sort_order` | `INT DEFAULT 0` | 排序序号 |
| `metadata` | `JSONB DEFAULT '{}'` | 扩展元数据 |

## 关键索引策略

- 所有软删除表使用 `WHERE deleted_at = 0` 条件唯一索引保证 active 记录唯一性
- 外键关系使用 `ON UPDATE CASCADE ON DELETE RESTRICT`
- 时间字段和状态字段均有独立索引和组合索引支持列表查询
- `sys_data_permission` 支持 subject_type + subject_id 组合索引
