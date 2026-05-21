import { createRootRouteWithContext, Outlet, redirect } from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'
import { router } from '~/router'

function onBack() {
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

export const Route = createRootRouteWithContext<RouterContext>()({
  beforeLoad: async ({ location, context }) => {
    const { account } = context
    if (account.token === '' && location.pathname !== '/login') {
      throw redirect({
        to: '/login',
      })
    }
  },
  notFoundComponent: () => <NotFoundComponent onBack={() => onBack()} />,
  errorComponent: () => <ErrorComponent onBack={() => onBack()} />,
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <>
      <Outlet />
      {/* <TanStackDevtools
        plugins={[
          {
            name: 'TanStack Router',
            render: <TanStackRouterDevtoolsPanel />,
          },
        ]}
      /> */}
      <TanStackRouterDevtools position="bottom-right" />
    </>
  )
}
