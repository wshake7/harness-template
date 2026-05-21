import fs from 'node:fs'
import path from 'node:path'

export function createRemoveMswPlugin(mode: string) {
  return {
    name: 'remove-msw-in-prod',
    apply: 'build',
    closeBundle() {
      if (mode === 'prod') {
        const file = path.resolve('dist/mockServiceWorker.js')
        if (fs.existsSync(file)) {
          fs.unlinkSync(file)
          console.log('removed mockServiceWorker.js')
        }
      }
    },
  } as const
}
