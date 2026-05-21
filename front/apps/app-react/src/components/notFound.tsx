interface Props {
  onBack: () => void
}

export function NotFoundComponent(props: Props) {
  const { onBack } = props

  return (
    <>
      404
      <button onClick={onBack}>Back Home</button>
    </>
  )
}
