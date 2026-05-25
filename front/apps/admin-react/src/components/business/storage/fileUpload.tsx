import type { UploadFile, UploadProps } from 'antd'
import type { StorageFileAsset, StorageFilePresignedReq } from '~/api/business/storageFile'
import { UploadOutlined } from '@ant-design/icons'
import { Button, Upload } from 'antd'
import { StorageFileApi } from '~/api/business/storageFile'
import { gMessage } from '~/utils/message'

export interface StorageFileUploadProps {
  value?: StorageFileAsset[]
  onChange?: (value: StorageFileAsset[]) => void
  bizType?: string
  bizID?: string
  maxCount?: number
  accept?: string
  disabled?: boolean
  storageApi?: typeof StorageFileApi
  openWindow?: (url: string) => void
}

function toUploadFile(asset: StorageFileAsset): UploadFile<StorageFileAsset> {
  return {
    uid: String(asset.id),
    name: asset.originalName,
    status: 'done',
    size: asset.size,
    type: asset.contentType,
    response: asset,
  }
}

function appendAsset(list: StorageFileAsset[], asset: StorageFileAsset, maxCount?: number) {
  const next = [...list, asset]
  if (!maxCount || next.length <= maxCount) {
    return next
  }
  return next.slice(next.length - maxCount)
}

export function StorageFileUpload(props: StorageFileUploadProps) {
  const {
    value = [],
    onChange,
    bizType,
    bizID,
    maxCount,
    accept,
    disabled,
    storageApi = StorageFileApi,
    openWindow = url => window.open(url, '_blank', 'noopener,noreferrer'),
  } = props

  const uploadProps: UploadProps<StorageFileAsset> = {
    accept,
    disabled,
    multiple: false,
    maxCount,
    fileList: value.map(toUploadFile),
    customRequest: async ({ file, onError, onSuccess }) => {
      try {
        const res = await storageApi.upload(file as File, {
          bizType,
          bizID,
        })
        const asset = res.data
        if (!asset) {
          throw new Error('上传返回为空')
        }
        onChange?.(appendAsset(value, asset, maxCount))
        onSuccess?.(asset)
      }
      catch (error) {
        const message = error instanceof Error ? error.message : '上传失败'
        gMessage.error(message)
        onError?.(error as Error)
      }
    },
    onRemove: async (file) => {
      const target = file.response ?? value.find(item => String(item.id) === file.uid)
      if (!target) {
        return false
      }
      onChange?.(value.filter(item => item.id !== target.id))
      return true
    },
    onPreview: async (file) => {
      const target = file.response ?? value.find(item => String(item.id) === file.uid)
      if (!target) {
        return
      }
      try {
        const req: StorageFilePresignedReq = {
          id: target.id,
          disposition: 'attachment',
        }
        const res = await storageApi.presigned(req)
        if (res.data?.url) {
          openWindow(res.data.url)
        }
      }
      catch (error) {
        const message = error instanceof Error ? error.message : '获取下载链接失败'
        gMessage.error(message)
      }
    },
  }

  return (
    <Upload {...uploadProps}>
      <Button icon={<UploadOutlined />} disabled={disabled}>
        上传文件
      </Button>
    </Upload>
  )
}

export { appendAsset }
