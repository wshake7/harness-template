import path from 'node:path'
import { createPlaywrightCtConfig } from '@vp/build-config'

const base = createPlaywrightCtConfig()

export default {
  ...base,
  testMatch: ['src/**/*.spec.tsx'],
  use: {
    ...base.use,
    channel: 'chrome',
  },
  ctViteConfig: {
    resolve: {
      alias: {
        '~': path.resolve(import.meta.dirname, './src'),
        '@': path.resolve(import.meta.dirname, '.'),
      },
    },
  },
}
