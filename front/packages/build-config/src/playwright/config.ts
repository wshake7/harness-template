import type { PlaywrightTestConfig } from '@playwright/test'
import process from 'node:process'
import { devices as ctDevices, defineConfig as defineCtConfig } from '@playwright/experimental-ct-react'
import { devices } from '@playwright/test'

export function createPlaywrightConfig(): PlaywrightTestConfig {
  return {
    testDir: './tests',
    timeout: 30 * 1000,
    expect: {
      timeout: 5000,
    },
    fullyParallel: true,
    forbidOnly: !!process.env.CI,
    retries: process.env.CI ? 2 : 0,
    workers: process.env.CI ? 1 : undefined,
    reporter: 'html',
    use: {
      actionTimeout: 0,
      trace: 'on-first-retry',
    },
    projects: [
      {
        name: 'chromium',
        use: {
          ...devices['Desktop Chrome'],
        },
      },
    ],
  }
}

export function createPlaywrightCtConfig() {
  return defineCtConfig({
    testDir: './',
    snapshotDir: './__snapshots__',
    timeout: 10 * 1000,
    fullyParallel: true,
    forbidOnly: !!process.env.CI,
    retries: process.env.CI ? 2 : 0,
    workers: process.env.CI ? 1 : undefined,
    reporter: 'html',
    use: {
      trace: 'on-first-retry',
      ctPort: 3100,
    },
    projects: [
      {
        name: 'chromium',
        use: { ...ctDevices['Desktop Chrome'] },
      },
    ],
  })
}
