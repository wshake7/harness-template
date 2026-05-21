import { createAccountStore } from '@vp/react-core'

export const useAccountStore = createAccountStore({
  onLogout: () => {
    useMenuTabsStore.getState().removeAll()
    useDeviceStore.getState().clear()
    useResourceMenuStore.getState().clear()
  },
})
