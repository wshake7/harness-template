export const HttpCode = {
  UN_KNOW: 0,
  SUCCESS: 1,
  ERROR: 2,
  FailRequestExpired: 3,
  FailRequestNonce: 4,
  FailRequestKey: 5,
  FailLogin: 100,
  FailAuth: 200,
} as const

export type CodeType = (typeof HttpCode)[keyof typeof HttpCode]

export const HttpCodeSet = new Set(Object.values(HttpCode))

export const XHeader = {
  XRequestTimestamp: 'X-Request-Timestamp',
  XRequestID: 'X-Request-ID',
  XRequestEncryptedKey: 'X-Request-Encrypted-Key',
  XRequestSignature: 'X-Request-Signature',
  XResponseIsEncrypt: 'X-Response-Is-Encrypt',
  Token: 'Token',
} as const

export interface ApiResponse<T = unknown> {
  readonly code: CodeType
  readonly msg: string
  readonly data?: T
}
