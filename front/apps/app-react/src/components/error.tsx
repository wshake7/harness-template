interface Props {
  onBack: () => void
}

export function ErrorComponent(props: Props) {
  const { onBack } = props
  return (
    <>
      error
      <button onClick={onBack}>Back Home</button>
    </>
  )
}
