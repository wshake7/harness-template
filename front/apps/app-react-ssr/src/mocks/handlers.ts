import type { RequestHandler } from 'msw'

const modules = import.meta.glob('./handlers/**/*.ts', { eager: true }) as Record<string, any>
export const handlers: RequestHandler[] = Object.values(modules).flatMap((m) => {
  return Object.values(m).flatMap(exported => (Array.isArray(exported) ? exported : []))
})
