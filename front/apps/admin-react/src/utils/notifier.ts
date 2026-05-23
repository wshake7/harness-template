import type { AppNotifier } from '@vp/core'
import { gMessage } from './message'

export const appNotifier: AppNotifier = {
  success: message => gMessage.success(message),
  error: message => gMessage.error(message),
  warning: message => gMessage.warning(message),
  info: message => gMessage.info(message),
}

// 在 Error 上挂一个 notified 标记，表示底层（HttpCodeCheck/请求层）已经弹过 toast。
// 业务层 catch 时通过 wasNotified 判断，避免重复提示。
type NotifiedError = Error & { notified?: true }

export function markNotified<E extends Error>(err: E): E & { notified: true } {
  (err as NotifiedError).notified = true
  return err as E & { notified: true }
}

export function notifiedError(message: string): Error {
  return markNotified(new Error(message))
}

export function wasNotified(err: unknown): boolean {
  return Boolean(err && typeof err === 'object' && (err as NotifiedError).notified)
}

// 业务 catch 的统一收口：底层已经提示过就不再提示，否则按 fallback 文案提示一次。
export function notifyError(err: unknown, fallbackMessage: string): void {
  if (wasNotified(err)) {
    return
  }
  appNotifier.error(fallbackMessage)
}
