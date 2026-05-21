import type { AppNotifier } from '@vp/core'
import { gMessage } from './message'

export const appNotifier: AppNotifier = {
  success: message => gMessage.success(message),
  error: message => gMessage.error(message),
  warning: message => gMessage.warning(message),
  info: message => gMessage.info(message),
}
