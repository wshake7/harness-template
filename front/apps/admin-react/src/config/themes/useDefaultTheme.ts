// ========== App.tsx ==========
import type { ConfigProviderProps } from 'antd'
import { theme } from 'antd'

const useDefaultTheme = () => {
  return useMemo<ConfigProviderProps>(
    () => ({
      theme: { algorithm: theme.defaultAlgorithm },
    }),
    [],
  )
}
export default useDefaultTheme
