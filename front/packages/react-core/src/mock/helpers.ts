import type { ApiResponse } from '@vp/core'
import { HttpCode } from '@vp/core'

export function createMockUrl(port?: number) {
  const baseUrl = `http://localhost:${port}`
  return (path: string) => `${baseUrl}${path}`
}

export function fail(msg: string): ApiResponse {
  return {
    code: HttpCode.ERROR,
    msg,
  }
}

export function success<T>(data: T): ApiResponse<T> {
  return {
    code: HttpCode.SUCCESS,
    msg: 'success',
    data,
  }
}
