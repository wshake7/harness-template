import { StartClient } from '@tanstack/react-start/client'
import { StrictMode } from 'react'
import { hydrateRoot } from 'react-dom/client'
import { worker } from './mocks/browser'

if (import.meta.env.VITE_MOCK === 'true') {
  worker.start({
    onUnhandledRequest: 'bypass',
  })
}

hydrateRoot(
  document,
  <StrictMode>
    <StartClient />
  </StrictMode>,
)
