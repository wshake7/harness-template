import handler, { createServerEntry } from '@tanstack/react-start/server-entry'
// import { server } from './mocks/node'

// server.listen({ onUnhandledRequest: 'bypass' })

export default createServerEntry({
  async fetch(request) {
    if (import.meta.env.VITE_MOCK === 'true') {
      const { server } = await import('./mocks/node')
      server.listen({ onUnhandledRequest: 'bypass' })
    }
    return handler.fetch(request)
  },
})
