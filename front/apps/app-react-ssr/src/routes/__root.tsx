import { createRootRouteWithContext, HeadContent, redirect, Scripts, useRouter } from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'
import { seo } from '~/lib/seo'
import appCss from '~/styles/index.css?url'

export const Route = createRootRouteWithContext<RouterContext>()({
  head: () => ({
    meta: [
      {
        charSet: 'utf-8',
      },
      {
        name: 'viewport',
        content: 'width=device-width, initial-scale=1',
      },
      ...seo({
        title: 'TanStack Start | Type-Safe, Client-First, Full-Stack React Framework',
        description: `TanStack Start is a type-safe, client-first, full-stack React framework. `,
      }),
    ],
    links: [{ rel: 'stylesheet', href: appCss }],
  }),
  beforeLoad: async ({ location, context }) => {
    const { account } = context
    if (account.token === '' && location.pathname !== '/login') {
      throw redirect({
        to: '/login',
      })
    }
  },
  notFoundComponent: NotFoundRouteComponent,
  errorComponent: ErrorRouteComponent,
  shellComponent: RootDocument,
})

function useBackHandler() {
  const router = useRouter()

  return () => {
    const token = useAccountStore.getState().token

    if (token !== '') {
      const isFromSameOrigin = document.referrer
        && new URL(document.referrer).origin === window.location.origin

      if (router.history.length > 1 && isFromSameOrigin) {
        router.history.back()
      }
      else {
        router.navigate({ to: '/' })
      }
    }
    else {
      router.navigate({ to: '/login' })
    }
  }
}

function NotFoundRouteComponent() {
  const onBack = useBackHandler()

  return <NotFoundComponent onBack={onBack} />
}

function ErrorRouteComponent() {
  const onBack = useBackHandler()

  return <ErrorComponent onBack={onBack} />
}

function RootDocument({ children }: { children: React.ReactNode }) {
  return (
    <html>
      <head>
        <HeadContent />
      </head>
      <body>
        {children}
        <TanStackRouterDevtools position="bottom-right" />
        <Scripts />
      </body>
    </html>
  )
}
