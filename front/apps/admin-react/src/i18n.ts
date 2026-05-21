import { createVpI18n } from '@vp/react-core'
import zh from '@/locales/zh.json'

const i18n = createVpI18n({
  resources: {
    zh: { translation: zh },
  },
})

export default i18n
