# Go Admin API Contract

> 基线冻结时间：2026-05-26
> 基于：`backend/go/admin/internal/router/` 路由注册代码

## 路由分组与中间件

```
/api (TimestampMiddleware, LanguageMiddleware)
├── /account (PublicMiddleware, EncryptMiddleware for login/changePwd)
│   ├── POST /login/pwd          → 密码登录
│   ├── POST /changePwd          → 修改密码 (Auth)
│   └── GET  /logout             → 登出 (Auth)
├── /encrypt
│   └── GET  /public/key         → 获取加密公钥
└── [Auth + Encrypt] 业务路由组
    ├── /events                  → GET SSE events (Auth, no encrypt body)
    ├── /storage/file/*          → 文件上传 (Auth, no encrypt body)
    ├── /sys/role/*              → 角色管理
    ├── /sys/user/*              → 用户管理
    ├── /sys/dict/*              → 字典管理
    ├── /sys/language/*          → 语言管理
    ├── /sys/api/log/*           → API 日志
    ├── /sys/login/log/*         → 登录日志
    ├── /sys/resource/menu/*     → 菜单资源
    ├── /sys/resource/api/*      → API 资源
    ├── /sys/job/schedule/*      → 任务调度
    ├── /sys/job/execution/*     → 任务执行
    └── /knowledge/*             → 知识库
```

## 完整接口清单

### 账号接口

| 方法 | 路径 | 认证 | 加密 | 说明 |
|------|------|------|------|------|
| POST | `/api/account/login/pwd` | Public | Encrypt | 密码登录，返回 token + publicKey，设置 Cookie |
| POST | `/api/account/changePwd` | Auth | Encrypt | 修改密码 |
| GET | `/api/account/logout` | Auth | No | 登出，清除 session |

### 加密接口

| 方法 | 路径 | 认证 | 加密 | 说明 |
|------|------|------|------|------|
| GET | `/api/encrypt/public/key` | Public | No | 获取本次会话 RSA 公钥 |

### SSE 事件

| 方法 | 路径 | 认证 | 加密 | 说明 |
|------|------|------|------|------|
| GET | `/api/events` | Auth | Encrypted SSE | 服务端推送事件（加密） |

### 存储文件接口

| 方法 | 路径 | 认证 | 加密 | 说明 |
|------|------|------|------|------|
| POST | `/api/storage/file/upload` | Auth | No (multipart) | 服务端代理上传 |
| POST | `/api/storage/file/prepareUpload` | Auth | No | 创建 pending_upload 资产，返回 presigned URL |
| POST | `/api/storage/file/completeUpload` | Auth | No | 校验对象存在，切换到 active |
| POST | `/api/storage/file/detail` | Auth | No | 获取文件详情 |
| POST | `/api/storage/file/presigned` | Auth | No | 获取下载 presigned URL |
| POST | `/api/storage/file/del` | Auth | No | 标记删除 |

### 角色管理

| 方法 | 路径 | 认证 | 加密 | 说明 |
|------|------|------|------|------|
| POST | `/api/sys/role/list` | Auth+Casbin | Encrypt | 分页列表 |
| GET | `/api/sys/role/tree` | Auth+Casbin | Encrypt | 角色树 |
| GET | `/api/sys/role/:id/permissions` | Auth+Casbin | Encrypt | 获取角色权限 |
| POST | `/api/sys/role/create` | Auth+Casbin | Encrypt | 创建角色 |
| POST | `/api/sys/role/update` | Auth+Casbin | Encrypt | 更新角色 |
| POST | `/api/sys/role/del` | Auth+Casbin | Encrypt | 删除角色 |
| POST | `/api/sys/role/permissions` | Auth+Casbin | Encrypt | 保存角色权限 |

### 用户管理

| 方法 | 路径 | 认证 | 加密 | 说明 |
|------|------|------|------|------|
| POST | `/api/sys/user/list` | Auth+Casbin | Encrypt | 分页列表 |
| POST | `/api/sys/user/create` | Auth+Casbin | Encrypt | 创建用户 |
| POST | `/api/sys/user/update` | Auth+Casbin | Encrypt | 更新用户 |
| POST | `/api/sys/user/del` | Auth+Casbin | Encrypt | 删除用户 |

### 字典管理

| 方法 | 路径 | 认证 | 加密 | 说明 |
|------|------|------|------|------|
| POST | `/api/sys/dict/type/list` | Auth+Casbin | Encrypt | 字典类型分页列表 |
| POST | `/api/sys/dict/type/create` | Auth+Casbin | Encrypt | 创建字典类型 |
| POST | `/api/sys/dict/type/update` | Auth+Casbin | Encrypt | 更新字典类型 |
| POST | `/api/sys/dict/type/del` | Auth+Casbin | Encrypt | 批量删除字典类型 |
| POST | `/api/sys/dict/entry/list` | Auth+Casbin | Encrypt | 字典项分页列表 |
| POST | `/api/sys/dict/entry/match` | Auth+Casbin | Encrypt | 按类型和 key 匹配条目 |
| POST | `/api/sys/dict/entry/create` | Auth+Casbin | Encrypt | 创建字典项 |
| POST | `/api/sys/dict/entry/update` | Auth+Casbin | Encrypt | 更新字典项 |
| POST | `/api/sys/dict/entry/del` | Auth+Casbin | Encrypt | 删除字典项 |
| POST | `/api/sys/dict/entry/batch/copy` | Auth+Casbin | Encrypt | 批量复制字典项 |

### 语言管理

| 方法 | 路径 | 认证 | 加密 | 说明 |
|------|------|------|------|------|
| POST | `/api/sys/language/type/list` | Auth+Casbin | Encrypt | 语言类型分页列表 |
| POST | `/api/sys/language/type/create` | Auth+Casbin | Encrypt | 创建语言类型 |
| POST | `/api/sys/language/type/update` | Auth+Casbin | Encrypt | 更新语言类型 |
| POST | `/api/sys/language/type/del` | Auth+Casbin | Encrypt | 删除语言类型 |
| POST | `/api/sys/language/entry/list` | Auth+Casbin | Encrypt | 语言项分页列表 |
| POST | `/api/sys/language/entry/create` | Auth+Casbin | Encrypt | 创建语言项 |
| POST | `/api/sys/language/entry/update` | Auth+Casbin | Encrypt | 更新语言项 |
| POST | `/api/sys/language/entry/del` | Auth+Casbin | Encrypt | 删除语言项 |
| POST | `/api/sys/language/entry/batch/create` | Auth+Casbin | Encrypt | 批量创建语言项 |

### API 日志

| 方法 | 路径 | 认证 | 加密 | 说明 |
|------|------|------|------|------|
| POST | `/api/sys/api/log/list` | Auth+Casbin | Encrypt | API 日志分页列表 |
| POST | `/api/sys/api/log/detail` | Auth+Casbin | Encrypt | API 日志详情 |

### 登录日志

| 方法 | 路径 | 认证 | 加密 | 说明 |
|------|------|------|------|------|
| POST | `/api/sys/login/log/list` | Auth+Casbin | Encrypt | 登录日志分页列表 |
| POST | `/api/sys/login/log/detail` | Auth+Casbin | Encrypt | 登录日志详情 |

### 菜单资源

| 方法 | 路径 | 认证 | 加密 | 说明 |
|------|------|------|------|------|
| POST | `/api/sys/resource/menu/list` | Auth+Casbin | Encrypt | 菜单分页列表 |
| GET | `/api/sys/resource/menu/tree` | Auth+Casbin | Encrypt | 菜单树 |
| POST | `/api/sys/resource/menu/create` | Auth+Casbin | Encrypt | 创建菜单 |
| POST | `/api/sys/resource/menu/update` | Auth+Casbin | Encrypt | 更新菜单 |
| POST | `/api/sys/resource/menu/del` | Auth+Casbin | Encrypt | 删除菜单 |

### API 资源管理

| 方法 | 路径 | 认证 | 加密 | 说明 |
|------|------|------|------|------|
| POST | `/api/sys/resource/api/list` | Auth+Casbin | Encrypt | API 资源分页列表 |
| POST | `/api/sys/resource/api/create` | Auth+Casbin | Encrypt | 创建 API 资源 |
| POST | `/api/sys/resource/api/update` | Auth+Casbin | Encrypt | 更新 API 资源 |
| POST | `/api/sys/resource/api/del` | Auth+Casbin | Encrypt | 删除 API 资源 |

### 任务调度

| 方法 | 路径 | 认证 | 加密 | 说明 |
|------|------|------|------|------|
| POST | `/api/sys/job/schedule/options` | Auth+Casbin | Encrypt | 获取任务选项 |
| POST | `/api/sys/job/schedule/list` | Auth+Casbin | Encrypt | 任务分页列表 |
| POST | `/api/sys/job/schedule/detail` | Auth+Casbin | Encrypt | 任务详情 |
| POST | `/api/sys/job/schedule/create` | Auth+Casbin | Encrypt | 创建任务 |
| POST | `/api/sys/job/schedule/update` | Auth+Casbin | Encrypt | 更新任务 |
| POST | `/api/sys/job/schedule/del` | Auth+Casbin | Encrypt | 删除任务 |
| POST | `/api/sys/job/schedule/switch` | Auth+Casbin | Encrypt | 启用/停用任务 |
| POST | `/api/sys/job/schedule/sync` | Auth+Casbin | Encrypt | 同步到 Temporal |
| POST | `/api/sys/job/schedule/trigger` | Auth+Casbin | Encrypt | 手动触发 |

### 任务执行

| 方法 | 路径 | 认证 | 加密 | 说明 |
|------|------|------|------|------|
| POST | `/api/sys/job/execution/list` | Auth+Casbin | Encrypt | 执行记录分页列表 |
| POST | `/api/sys/job/execution/detail` | Auth+Casbin | Encrypt | 执行记录详情 |
| POST | `/api/sys/job/execution/cancel` | Auth+Casbin | Encrypt | 取消执行 |
| POST | `/api/sys/job/execution/retry` | Auth+Casbin | Encrypt | 重试执行 |

### 知识库

| 方法 | 路径 | 认证 | 加密 | 说明 |
|------|------|------|------|------|
| POST | `/api/knowledge/collection/list` | Auth+Casbin | Encrypt | 集合分页列表 |
| POST | `/api/knowledge/collection/detail` | Auth+Casbin | Encrypt | 集合详情 |
| POST | `/api/knowledge/collection/create` | Auth+Casbin | Encrypt | 创建集合 |
| POST | `/api/knowledge/collection/update` | Auth+Casbin | Encrypt | 更新集合 |
| POST | `/api/knowledge/collection/del` | Auth+Casbin | Encrypt | 删除集合 |
| POST | `/api/knowledge/document/list` | Auth+Casbin | Encrypt | 文档分页列表 |
| POST | `/api/knowledge/document/listByCollection` | Auth+Casbin | Encrypt | 按集合查文档 |
| POST | `/api/knowledge/document/detail` | Auth+Casbin | Encrypt | 文档详情 |
| POST | `/api/knowledge/document/create` | Auth+Casbin | Encrypt | 创建文档 |
| POST | `/api/knowledge/document/importFile` | Auth+Casbin | Encrypt | 从 file_asset 导入文档 |
| POST | `/api/knowledge/document/update` | Auth+Casbin | Encrypt | 更新文档 |
| POST | `/api/knowledge/document/del` | Auth+Casbin | Encrypt | 删除文档 |

## 业务响应码约定

| Code | 含义 |
|------|------|
| `1` | 成功 |
| `2` | 通用失败 |
| 其他 | 特定业务状态（如请求超时/重放、认证失败、授权失败） |

## 加密协议约定

1. `GET /api/encrypt/public/key` 返回 base64 编码的 RSA SPKI 公钥，客户端生成 AES key 并用 RSA 公钥加密
2. 登录后请求头携带：`X-Request-Encrypted-Key`（RSA 加密的 AES key）、`X-Request-Signature`（签名）、`X-Request-Timestamp`、`X-Request-ID`
3. AAD 拼接字符串由 query string (sort keys) + body 组成
4. 响应体通过 `X-Response-Is-Encrypt` 头标识是否加密
5. multipart upload 路径 `/api/storage/file/upload` 需要认证但跳过 JSON body 加密
6. SSE `/api/events` 使用独立加密通道

## 分页请求格式

```json
{
  "page": 1,
  "pageSize": 10,
  "orderBy": "id desc",
  "filters": [
    {"field": "username", "operator": "contains", "value": "admin"}
  ]
}
```

## 分页响应格式

```json
{
  "code": 1,
  "msg": "success",
  "data": {
    "items": [...],
    "total": 100
  }
}
```
