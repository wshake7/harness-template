import { createMockStore } from '@vp/react-core'

const mockStore = createMockStore()

export const useMockStore = mockStore.useMockStore
export const waitForHydration = mockStore.waitForHydration
