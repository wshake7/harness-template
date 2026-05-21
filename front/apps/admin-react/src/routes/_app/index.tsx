import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_app/')({
  beforeLoad: () => { },
  component: () => null,
})
