export type NotifyMessage = string

export interface AppNotifier {
  success: (message: NotifyMessage) => void
  error: (message: NotifyMessage) => void
  warning: (message: NotifyMessage) => void
  info: (message: NotifyMessage) => void
}

export const noopNotifier: AppNotifier = {
  success: () => {},
  error: () => {},
  warning: () => {},
  info: () => {},
}

export function createConsoleNotifier(prefix = '[notify]'): AppNotifier {
  return {
    success: message => console.info(prefix, message),
    error: message => console.error(prefix, message),
    warning: message => console.warn(prefix, message),
    info: message => console.info(prefix, message),
  }
}
