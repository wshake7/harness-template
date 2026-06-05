# Frontend API Call Samples

> 基线冻结时间：2026-05-26
> 基于：`front/apps/admin-react` 请求层与 `@vp/request` 加密协议
> 用途：Java Admin 实现时的 contract test 参考样例

## 1. 登录成功样例

**Request:**
```http
POST /api/account/login/pwd
Content-Type: application/json
X-Request-Timestamp: 1711411200000
X-Request-ID: uuid-v4

{
  "username": "admin",
  "password": "encrypted_password_here"
}
```

**Response (HTTP 200):**
```json
{
  "code": 1,
  "msg": "success",
  "data": {
    "token": "sa-token-value-here",
    "publicKey": "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A..."
  }
}
```

**Cookie:** `token=sa-token-value-here; SameSite=Lax; HttpOnly`

---

## 2. 登录失败样例（账号不存在）

**Response (HTTP 200):**
```json
{
  "code": <non-1 code>,
  "msg": "账号或密码错误"
}
```

---

## 3. 加密 POST 请求样例

**前提：** 前端已通过 `/api/encrypt/public/key` 获取 base64 编码的 RSA SPKI 公钥

**流程：**
1. 前端生成随机 AES key (256-bit)
2. 用 RSA 公钥加密 AES key → `X-Request-Encrypted-Key` 请求头
3. 用 AES key 加密 JSON body → ciphertext
4. 计算签名：`SHA256(sortQueryKeys(query) + bodyPlaintext + timestamp + requestId + AESKey)` → `X-Request-Signature`
5. 发送加密请求

**Request:**
```http
POST /api/sys/user/list
Content-Type: application/json
X-Request-Encrypted-Key: <RSA-encrypted-aes-key-base64>
X-Request-Signature: <sha256-hex-signature>
X-Request-Timestamp: 1711411200000
X-Request-ID: uuid-v4
Cookie: token=sa-token-value-here

<RSA-encrypted-json-body>
```

**Response (HTTP 200):**
```http
X-Response-Is-Encrypt: true

{
  "code": 1,
  "msg": "success",
  "data": "<AES-encrypted-response-json>"
}
```

**非加密响应（X-Response-Is-Encrypt 缺失或为 false）：**
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

---

## 4. 权限拒绝样例

**Request (已登录但无权限):**
```http
POST /api/sys/role/list
Cookie: token=valid-token
```

**Response (HTTP 200):**
```json
{
  "code": <authorization-failure-code>,
  "msg": "无权限访问"
}
```

---

## 5. 未登录访问样例

**Request (无 Cookie):**
```http
POST /api/sys/role/list
```

**Response (HTTP 200):**
```json
{
  "code": <auth-failure-code>,
  "msg": "未登录或会话已过期"
}
```

---

## 6. 时间戳过期样例

**Request (timestamp 超过允许窗口):**
```http
POST /api/sys/user/list
X-Request-Timestamp: 1711324800000
Cookie: token=valid-token
```

**Response (HTTP 200):**
```json
{
  "code": <request-timeout-code>,
  "msg": "请求已过期"
}
```

---

## 7. 分页列表查询样例

**Request:**
```http
POST /api/sys/user/list
Content-Type: application/json

{
  "page": 1,
  "pageSize": 10,
  "orderBy": "id desc",
  "filters": [
    {"field": "username", "operator": "contains", "value": "admin"}
  ]
}
```

**Response (HTTP 200):**
```json
{
  "code": 1,
  "msg": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "username": "admin",
        "nickname": "管理员",
        "languageCode": "zh-CN",
        "isEnabled": true,
        "createdAt": "2026-05-26T00:00:00Z",
        "updatedAt": "2026-05-26T00:00:00Z"
      }
    ],
    "total": 1
  }
}
```

---

## 8. 文件 Prepare/Complete 直传样例

**Prepare Upload:**
```http
POST /api/storage/file/prepareUpload
Content-Type: application/json
Cookie: token=valid-token

{
  "originalName": "document.md",
  "contentType": "text/markdown",
  "size": 1024
}
```

**Response:**
```json
{
  "code": 1,
  "msg": "success",
  "data": {
    "fileAssetId": 123,
    "objectKey": "uploads/uuid-doc.md",
    "presignedUrl": "http://minio:9000/bucket/uploads/uuid-doc.md?signature=...",
    "status": "pending_upload"
  }
}
```

**Complete Upload:**
```http
POST /api/storage/file/completeUpload
Content-Type: application/json
Cookie: token=valid-token

{
  "fileAssetId": 123
}
```

**Response:**
```json
{
  "code": 1,
  "msg": "success",
  "data": {
    "fileAssetId": 123,
    "objectKey": "uploads/uuid-doc.md",
    "status": "active",
    "size": 1024
  }
}
```

---

## 9. 知识文档索引失败样例

**ImportFile (成功创建文档但索引失败):**
```http
POST /api/knowledge/document/importFile
Content-Type: application/json
Cookie: token=valid-token

{
  "collectionId": 1,
  "fileAssetId": 123
}
```

**Response (索引异步失败，文档状态反映):**
```json
{
  "code": 1,
  "msg": "success",
  "data": {
    "id": 456,
    "documentId": "doc-uuid",
    "title": "document.md",
    "vectorStatus": "failed",
    "indexingError": "embedding API key not configured"
  }
}
```

---

## 加密协议 Golden Test 固定参数

以下固定参数用于跨语言加密协议验证（Java 实现 `CryptoService` 时使用）：

| 参数 | 固定值 |
|------|--------|
| RSA Private Key | 从测试资源 `golden/crypto-request.json` 加载 |
| RSA Public Key | 从私钥推导 |
| AES Key (plain) | 256-bit random, hex-encoded |
| Timestamp | `1711411200000` |
| Request ID | `test-request-id-0001` |
| Query String | `page=1&pageSize=10` |
| Request Body (plain) | `{"username":"admin"}` |
| AAD String | `page=1&pageSize=10{"username":"admin"}` (query keys sorted + body) |
