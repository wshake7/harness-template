import { createEnv } from '@t3-oss/env-core'
import { z } from 'zod'

export function createVpEnv(runtimeEnv: Record<string, string | number | boolean | undefined>) {
  return createEnv({
    server: {},
    clientPrefix: 'VITE_',
    client: {
      VITE_API_URL: z.url().optional(),
      VITE_PORT: z.coerce.number().optional(),
      VITE_MOCK: z.string().transform(val => val === 'true'),
    },
    runtimeEnv,
    emptyStringAsUndefined: true,
  })
}
