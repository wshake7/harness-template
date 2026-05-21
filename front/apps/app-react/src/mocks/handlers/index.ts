import { createMockUrl, fail, success } from '@vp/react-core'
import { gEnv } from '~/env'

export const url = createMockUrl(gEnv.VITE_PORT)
export { fail, success }
