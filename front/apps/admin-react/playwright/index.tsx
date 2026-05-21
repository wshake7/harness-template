import type { ReactNode } from 'react'
import { ConfigProvider } from 'antd'

export default function CtWrapper({ children }: { children: ReactNode }) {
  return <ConfigProvider>{children}</ConfigProvider>
}
