import { http, HttpResponse } from 'msw'
import { fail, success, url } from '.'

let nextID = 1
const files = new Map<number, any>()

export const storageFileHandlers = [
  http.post(url('/api/storage/file/upload'), async ({ request }) => {
    const formData = await request.formData()
    const file = formData.get('file')
    if (!(file instanceof File) || file.size <= 0) {
      return HttpResponse.json(fail('上传文件不能为空'))
    }

    const id = nextID++
    const asset = {
      id,
      engine: 'minio',
      bucket: 'admin-files',
      objectKey: `uploads/mock/${id}-${file.name}`,
      originalName: file.name,
      contentType: file.type || 'application/octet-stream',
      size: file.size,
      sha256: `mock-sha-${id}`,
      bizType: formData.get('bizType')?.toString() || '',
      bizID: formData.get('bizID')?.toString() || '',
      createdAt: new Date().toISOString(),
    }
    files.set(id, asset)
    return HttpResponse.json(success(asset))
  }),

  http.post(url('/api/storage/file/detail'), async ({ request }) => {
    const { id } = await request.json() as { id: number }
    const asset = files.get(id)
    if (!asset) {
      return HttpResponse.json(fail('文件不存在'))
    }
    return HttpResponse.json(success(asset))
  }),

  http.post(url('/api/storage/file/presigned'), async ({ request }) => {
    const { id } = await request.json() as { id: number }
    const asset = files.get(id)
    if (!asset) {
      return HttpResponse.json(fail('文件不存在'))
    }
    return HttpResponse.json(success({
      url: `https://mock-minio.local/${asset.objectKey}`,
      expiresAt: new Date(Date.now() + 3600 * 1000).toISOString(),
    }))
  }),

  http.post(url('/api/storage/file/del'), async ({ request }) => {
    const { id } = await request.json() as { id: number }
    files.delete(id)
    return HttpResponse.json(success(undefined))
  }),
]
