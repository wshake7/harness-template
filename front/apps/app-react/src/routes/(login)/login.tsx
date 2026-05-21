import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/(login)/login')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/(login)/login"!</div>
}
