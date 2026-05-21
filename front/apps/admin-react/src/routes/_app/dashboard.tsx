import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_app/dashboard')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <>
      <p>---------------</p>
      <div>Test Action12 Hello "/_app/dashboard"!</div>
      <p>---------------</p>
    </>
  )
}
