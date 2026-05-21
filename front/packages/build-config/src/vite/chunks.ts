export interface CreateManualChunksOptions {
  antd?: boolean
  monaco?: boolean
  shiki?: boolean
}

export function createManualChunks(options: CreateManualChunksOptions = {}) {
  return function manualChunks(id: string) {
    if (id.includes('node_modules/react/') || id.includes('node_modules/react-dom/')) {
      return 'vendor-react'
    }
    if (id.includes('node_modules/scheduler/') || id.includes('node_modules/use-sync-external-store/')) {
      return 'vendor-react'
    }

    if (options.monaco && id.includes('node_modules/modern-monaco')) {
      return 'vendor-monaco'
    }

    if (options.shiki) {
      if (id.includes('node_modules/@shikijs/langs') || id.includes('/shiki/langs/')) {
        return 'vendor-shiki-langs'
      }
      if (id.includes('node_modules/@shikijs/themes') || id.includes('/shiki/themes/')) {
        return 'vendor-shiki-themes'
      }
      if (id.includes('node_modules/shiki') || id.includes('node_modules/@shikijs')) {
        return 'vendor-shiki-core'
      }
    }

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

    if (options.antd) {
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
        return 'vendor-rc'
      }
    }

    if (id.includes('node_modules/zod')) {
      return 'vendor-zod'
    }

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
  }
}
