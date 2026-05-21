import { createFileRoute, Outlet } from '@tanstack/react-router'

export const Route = createFileRoute('/_app/logger')({
  component: RouteComponent,
})

function RouteComponent() {
  return <Outlet />
}
