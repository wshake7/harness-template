import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export type ThemeType = 'default' | 'shadcn' | 'cartoon'

interface ThemeStore {
  themeType: ThemeType
  setThemeType: (themeType: ThemeType) => void
  resolveTheme: <T>(themes: Record<ThemeType, T>) => T
}

export const useThemeStore = create<ThemeStore>()(
  persist(
    (set, get) => ({
      themeType: 'shadcn',
      setThemeType: themeType => set({ themeType }),
      resolveTheme: themes => themes[get().themeType],
    }),
    {
      name: 'theme-store',
    },
  ),
)

export const useAntTheme = () => {
  const { themeType, setThemeType, resolveTheme } = useThemeStore()

  const theme = resolveTheme({
    cartoon: useCartoonTheme(),
    shadcn: useShadcnTheme(),
    default: useDefaultTheme(),
  })

  return {
    theme,
    themeType,
    setThemeType,
  }
}
