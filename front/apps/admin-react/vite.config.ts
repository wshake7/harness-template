import fs from 'node:fs'
import path from 'node:path'
import process from 'node:process'
import tailwindcss from '@tailwindcss/vite'
import { devtools } from '@tanstack/devtools-vite'
import { tanstackRouter } from '@tanstack/router-plugin/vite'
// import { DevTools } from "@vitejs/devtools";
import react from '@vitejs/plugin-react'
import AutoImport from 'unplugin-auto-import/vite'
import { loadEnv } from 'vite'
// import Inject from "vite-plugin-inspect";
import { VitePWA } from 'vite-plugin-pwa'
import { defineConfig } from 'vite-plus'

const workspacePackageAliases = {
  '@vp/build-config': path.resolve(import.meta.dirname, '../../packages/build-config/src'),
  '@vp/core': path.resolve(import.meta.dirname, '../../packages/core/src'),
  '@vp/react-core': path.resolve(import.meta.dirname, '../../packages/react-core/src'),
  '@vp/request': path.resolve(import.meta.dirname, '../../packages/request/src'),
  '@vp/utils': path.resolve(import.meta.dirname, '../../packages/utils/src'),
}

export default defineConfig(({ mode }: { mode: string }) => {
  const env = loadEnv(mode, process.cwd(), '')

  return {
    base: env.GITHUB_ACTIONS === 'true' ? '/template/' : '/',
    staged: {
      '*': '',
    },
    resolve: {
      alias: workspacePackageAliases,
      dedupe: ['react', 'react-dom', 'react-i18next', 'zustand'],
      tsconfigPaths: true,
    },
    preview: {
      strictPort: true,
    },
    server: {
      port: Number(env.VITE_PORT),
      watch: {
        // wsl下热更新必须开
        usePolling: true,
        interval: 500,
      },
      proxy: {
        '/api': {
          target: env.VITE_API_URL,
          changeOrigin: true,
          secure: false,
          // rewrite: path => path.replace(/^\/api/, ''),
        },
      },
    },
    define: {
      __VUE_OPTIONS_API__: true,
      __VUE_PROD_DEVTOOLS__: true,
      __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: true,
    },
    build: {
      rollupOptions: {
        output: {
          manualChunks(id: string) {
            // ── React 核心 ──────────────────────────────────────────
            if (id.includes('node_modules/react/') || id.includes('node_modules/react-dom/')) {
              return 'vendor-react'
            }
            if (id.includes('node_modules/scheduler/') || id.includes('node_modules/use-sync-external-store/')) {
              return 'vendor-react'
            }

            // ── Monaco（体积最大，单独隔离）────────────────────────
            if (
              id.includes('node_modules/modern-monaco')
            ) {
              return 'vendor-monaco'
            }

            // ── Shiki 拆分：核心 / 语言包 / 主题 ──────────────────
            if (id.includes('node_modules/@shikijs/langs') || id.includes('/shiki/langs/')) {
              return 'vendor-shiki-langs'
            }
            if (id.includes('node_modules/@shikijs/themes') || id.includes('/shiki/themes/')) {
              return 'vendor-shiki-themes'
            }
            if (id.includes('node_modules/shiki') || id.includes('node_modules/@shikijs')) {
              return 'vendor-shiki-core'
            }

            // ── TanStack 细分 ───────────────────────────────────────
            if (id.includes('node_modules/@tanstack/react-query')) {
              return 'vendor-tanstack-query'
            }
            if (
              id.includes('node_modules/@tanstack/react-router')
              || id.includes('node_modules/@tanstack/router')
            ) {
              return 'vendor-tanstack-router'
            }
            if (id.includes('node_modules/@tanstack')) {
              return 'vendor-tanstack-misc'
            }

            // ── Ant Design：icons 单独分包 ──────────────────────────
            if (id.includes('node_modules/@ant-design/icons')) {
              return 'vendor-antd-icons'
            }
            if (id.includes('node_modules/antd')) {
              return 'vendor-antd'
            }
            if (id.includes('node_modules/@ant-design/')) {
              return 'vendor-antd-design'
            }
            if (id.includes('node_modules/rc-')) {
              return 'vendor-rc' // rc-* 是 antd 底层依赖
            }

            // ── Zod ─────────────────────────────────────────────────
            if (id.includes('node_modules/zod')) {
              return 'vendor-zod'
            }

            // ── 常见工具库 ──────────────────────────────────────────
            if (id.includes('node_modules/i18next') || id.includes('node_modules/react-i18next')) {
              return 'vendor-i18n'
            }
            if (id.includes('node_modules/lodash') || id.includes('node_modules/lodash-es')) {
              return 'vendor-lodash'
            }
            if (id.includes('node_modules/dayjs')) {
              return 'vendor-dayjs'
            }
            if (id.includes('node_modules/immer')) {
              return 'vendor-immer'
            }
            if (id.includes('node_modules/ahooks') || id.includes('node_modules/@ahooksjs')) {
              return 'vendor-ahooks'
            }

            // ── 其余 node_modules：按包名首字母分桶，避免单桶过大 ──
            if (id.includes('node_modules')) {
              const match = id.match(/node_modules\/(@[^/]+\/[^/]+|[^/]+)/)
              if (match) {
                const first = match[1].replace('@', '')[0].toLowerCase()
                if (first <= 'h') { return 'vendor-a-h' }
                if (first <= 'p') { return 'vendor-i-p' }
                return 'vendor-q-z'
              }
              return 'vendor-misc'
            }
          },
        },
      },
    },
    plugins: [
      // DevTools(),
      // Inject(),
      devtools(),
      tailwindcss(),
      tanstackRouter(),
      react(),
      AutoImport({
        include: [/\.[jt]sx?$/, /\.md$/, /tsr-split/],
        imports: ['react', 'ahooks', 'react-i18next'],
        dts: 'src/auto-imports.d.ts',
        dtsMode: 'overwrite',
        dirsScanOptions: {
          fileFilter: file => !file.replaceAll('\\', '/').includes('src/components/business/logger/'),
        },
        dirs: [
          'src/stores/**',
          'src/domains/**',
          'src/utils/**',
          'src/components/**',
          'src/api/**',
          'src/config/**',
        ],
      }) as any,
      VitePWA({
        registerType: 'autoUpdate', // 当有新版本时自动更新
        workbox: {
          globIgnores: [
            '**/vendor-monaco-*.js',
            '**/vendor-monaco-*.css',
            '**/codicon-*.ttf',
          ],
        },
        manifest: {
          name: '我的 React PWA 应用',
          short_name: 'ReactPWA',
          description: '一个使用 Vite 和 React 构建的 PWA 应用',
          theme_color: '#ffffff',
        },
      }),
      {
        name: 'remove-msw-in-prod',
        apply: 'build',
        closeBundle() {
          if (mode === 'prod') {
            const file = path.resolve('dist/mockServiceWorker.js')
            if (fs.existsSync(file)) {
              fs.unlinkSync(file)
              console.log('✅ removed mockServiceWorker.js')
            }
          }
        },
      },
    ],
  }
})
