import type { TabPaneProps } from 'antd'
import { create } from 'zustand'
import { immer } from 'zustand/middleware/immer'
import { router } from '~/router'

export interface MenuTabsItemProps extends TabPaneProps {
  key: string
  label: React.ReactNode
}

export type TargetKey = React.MouseEvent | React.KeyboardEvent | string

interface MenuTabsStore {
  items: MenuTabsItemProps[]
  add: (item: MenuTabsItemProps) => void
  remove: (targetKey: TargetKey) => void
  removeAll: () => void
}

export const useMenuTabsStore = create<MenuTabsStore>()(
  immer((set, get) => ({
    items: [],

    add(item) {
      set((state) => {
        const exists = state.items.some(i => i.key === item.key)
        if (!exists) {
          state.items.push(item)
        }
      })
    },

    remove(targetKey) {
      const items = get().items
      const index = items.findIndex(item => item.key === targetKey)
      if (index === -1) { return }
      const newItems = items.filter(item => item.key !== targetKey)
      let nextKey: string | undefined
      const activeKey = router.state.location.pathname
      if (targetKey === activeKey) {
        if (newItems.length > 0) {
          if (index < newItems.length) {
            nextKey = newItems[index].key
          }
          else {
            nextKey = newItems[index - 1].key
          }
        }
        else {
          nextKey = '/'
        }
      }
      set((state) => {
        state.items = newItems
      })
      if (nextKey) {
        router.navigate({ to: nextKey })
      }
    },

    removeAll() {
      set((state) => {
        state.items = []
      })
    },
  })),
)
