import { createAlova } from 'alova'
import { useSSE } from 'alova/client'
import adapterFetch from 'alova/fetch'
import reactHook from 'alova/react'
import { useEffect, useRef } from 'react'
import { createEncryptedRequestConfig, decryptText } from '~/api/encryptRequest'
import { XHeader } from '~/domains/http'

const SSE_API = createAlova({
  baseURL: '',
  statesHook: reactHook,
  requestAdapter: adapterFetch(),
  cacheFor: null,
  shareRequest: false,
  beforeRequest(method) {
    if (!useAccountStore.getState().token) {
      throw new Error('event stream canceled: unauthenticated')
    }
    method.config.headers[XHeader.Token] = useAccountStore.getState().token
  },
})

export function useEventStream() {
  const token = useAccountStore(state => state.token)
  const connectedRef = useRef(false)
  const tokenRef = useRef(token)
  const aesKeyRef = useRef<CryptoKey | undefined>(undefined)
  const headersRef = useRef<Record<string, string>>({})
  tokenRef.current = token

  const eventStream = useSSE(
    () => SSE_API.Get('/api/events', {
      headers: headersRef.current,
    }),
    {
      immediate: false,
      withCredentials: true,
      reconnectionTime: 10000,
      interceptByGlobalResponded: false,
      responseType: 'json',
    },
  )

  useEffect(() => {
    const offCount = eventStream.on<{ payload: string }>('count', async (event) => {
      if (!aesKeyRef.current) {
        console.error('[SSE decrypt error]', 'missing aes key')
        return
      }
      const decrypted = await decryptText(event.data.payload, aesKeyRef.current)
      const data = JSON.parse(decrypted) as { count: number }
      console.log('[SSE count]', data.count)
    }) as unknown

    eventStream.onError((event) => {
      console.error('[SSE error]', event.error)
      if (tokenRef.current === '') {
        eventStream.close()
        connectedRef.current = false
      }
    })

    return () => {
      if (typeof offCount === 'function') {
        offCount()
      }
    }
  }, [eventStream])

  useEffect(() => {
    if (token === '') {
      eventStream.close()
      connectedRef.current = false
      aesKeyRef.current = undefined
      headersRef.current = {}
      return
    }
    if (connectedRef.current) {
      return
    }

    const connectingToken = token
    connectedRef.current = true
    createEncryptedRequestConfig({
      headers: {
        [XHeader.Token]: connectingToken,
      },
    })
      .then((config) => {
        if (!config.aesKey || tokenRef.current !== connectingToken) {
          connectedRef.current = false
          return
        }
        aesKeyRef.current = config.aesKey
        headersRef.current = config.headers
        return eventStream.send()
      })
      .catch((error) => {
        connectedRef.current = false
        console.error('[SSE connect error]', error)
      })

    return () => {
      eventStream.close()
      connectedRef.current = false
      aesKeyRef.current = undefined
      headersRef.current = {}
    }
  }, [eventStream, token])
}
