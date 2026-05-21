import type { JointContent, MessageInstance } from 'antd/es/message/interface'
import { message } from 'antd'
import { useEffect } from 'react'

class GMessage {
  public messageApi: MessageInstance | null = null

  public setMessageApi(messageApi: any) {
    this.messageApi = messageApi
  }

  public error(content: JointContent, duration?: number | VoidFunction, onClose?: VoidFunction) {
    this.messageApi?.error(content, duration, onClose)
  }

  public success(content: JointContent, duration?: number | VoidFunction, onClose?: VoidFunction) {
    this.messageApi?.success(content, duration, onClose)
  }

  public info(content: JointContent, duration?: number | VoidFunction, onClose?: VoidFunction) {
    this.messageApi?.info(content, duration, onClose)
  }

  public warning(content: JointContent, duration?: number | VoidFunction, onClose?: VoidFunction) {
    this.messageApi?.warning(content, duration, onClose)
  }
}

const gMessage = new GMessage()
export { gMessage }

export default function GlobalMessage() {
  const [messageApi, contextHolder] = message.useMessage()

  useEffect(() => {
    gMessage.setMessageApi(messageApi)
    return () => {
      gMessage.setMessageApi(null)
    }
  }, [messageApi])

  return contextHolder
}
