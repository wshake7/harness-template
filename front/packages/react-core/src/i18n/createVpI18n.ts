import type { Resource } from 'i18next'
import i18next from 'i18next'
import LanguageDetector from 'i18next-browser-languagedetector'
import { initReactI18next } from 'react-i18next'

export interface CreateVpI18nOptions {
  resources: Resource
  fallbackLng?: string
  preload?: string[]
}

export function createVpI18n(options: CreateVpI18nOptions) {
  const instance = i18next.createInstance()
  void instance
    .use(initReactI18next)
    .use(LanguageDetector)
    .init({
      resources: options.resources,
      fallbackLng: options.fallbackLng ?? 'zh',
      preload: options.preload ?? ['en', 'zh'],
      interpolation: {
        escapeValue: false,
      },
    })

  return instance
}
