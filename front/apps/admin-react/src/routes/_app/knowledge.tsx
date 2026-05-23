import { createFileRoute, Outlet } from '@tanstack/react-router'

export const Route = createFileRoute('/_app/knowledge')({
  component: RouteComponent,
})

function RouteComponent() {
  return <Outlet />
}
