import process from 'node:process'
import { createPlaywrightConfig } from '@vp/build-config'

const base = createPlaywrightConfig()

export default {
  ...base,
  workers: 1,
  use: {
    ...base.use,
    baseURL: 'http://localhost:3000',
    channel: 'chrome',
  },
  webServer: {
    command: 'VITE_MOCK=true node_modules/.bin/vp dev --mode dev',
    url: 'http://localhost:3000',
    reuseExistingServer: !process.env.CI,
    timeout: 120_000,
  },
}
