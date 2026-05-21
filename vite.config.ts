import { defineConfig } from 'vite-plus'

export default defineConfig({
  staged: {
    '*': 'vp run lint:fix',
  },
})
