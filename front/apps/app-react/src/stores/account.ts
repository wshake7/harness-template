import { createAccountStore } from '@vp/react-core'

export const useAccountStore = createAccountStore({
  onLogout: () => {
    useDeviceStore.getState().clear()
  },
})
