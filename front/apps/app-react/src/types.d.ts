import type {
  AccountState as VpAccountState,
  ApiResponse,
  CodeType as VpCodeType,
  RouterContext as VpRouterContext,
} from '@vp/core'

export {}

declare global {
  type AccountState = VpAccountState
  type RouterContext = VpRouterContext
  type CodeType = VpCodeType
  type Res<T = unknown> = ApiResponse<T>
}
