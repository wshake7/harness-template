import { RouterProvider } from '@tanstack/react-router'
import { StrictMode } from 'react'
import ReactDOM from 'react-dom/client'
import { router } from './router'
import '~/styles/index.css'
import './i18n'

async function prepare() {
  if (import.meta.env.VITE_MOCK === 'true') {
    const { worker } = await import('./mocks/browser')
    await waitForHydration()
    return worker.start({ onUnhandledRequest: 'bypass' })
  }
}

const Root = () => {
  // todo useEventStream()

  return (
    <RouterProvider router={router} />
  )
}

prepare().then(() => {
  // Render the app
  const rootElement = document.getElementById('root')!
  if (!rootElement.innerHTML) {
    const root = ReactDOM.createRoot(rootElement)
    root.render(
      <StrictMode>
        <Root />
      </StrictMode>,
    )
  }
})
