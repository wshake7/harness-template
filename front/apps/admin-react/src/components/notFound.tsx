import { Button, Result } from 'antd'

interface Props {
  onBack: () => void
}

export function NotFoundComponent(props: Props) {
  const { onBack } = props

  return (
    <Result
      status="404"
      title="404"
      subTitle="Sorry, the page you visited does not exist."
      extra={(
        <Button
          type="primary"
          onClick={onBack}
        >
          Back Home
        </Button>
      )}
    />
  )
}
