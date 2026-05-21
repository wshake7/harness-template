import { http, HttpResponse } from 'msw'
import { success, url } from '.'

const MOCK_PUBLIC_KEY = 'MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQD2TJoOBk8KZE9z2CIYQsip86ANw7N7Wl15s4GmxBokdOv+pYdkHmyc24Jt2W4hqUZGetKsmcNjEu8ZtDj3dVrRjuopg+phzGmUrrW+oeGiacK35sLKLkK7ZZv9u9s4HQVOMq9M9v+1leCEJ4g17dGM2YppxTJjzrxNzPVg/XWJYwIDAQAB'

export const encryptHandlers = [
  http.get(url('/api/encrypt/public/key'), async () => {
    return HttpResponse.json(success({
      publicKey: MOCK_PUBLIC_KEY,
    }))
  }),
]
