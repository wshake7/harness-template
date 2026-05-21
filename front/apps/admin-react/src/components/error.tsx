import { Button, Result } from 'antd'

interface Props {
  onBack: () => void
}

export function ErrorComponent(props: Props) {
  const { onBack } = props
  return (
    <Result
      status="warning"
      title="There are some problems with your operation."
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
