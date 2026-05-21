import { beforeEach, expect, test, vi } from 'vite-plus/test'
import { useMenuTabsStore } from './menuTabs'

vi.mock('~/router', () => ({
  router: {
    state: { location: { pathname: '/dashboard' } },
    navigate: vi.fn(),
  },
}))

beforeEach(() => {
  useMenuTabsStore.setState({ items: [] })
})

test('add inserts new tab', () => {
  useMenuTabsStore.getState().add({ key: '/dashboard', label: 'Dashboard' })
  expect(useMenuTabsStore.getState().items).toHaveLength(1)
  expect(useMenuTabsStore.getState().items[0].key).toBe('/dashboard')
})

test('add ignores duplicate key', () => {
  const store = useMenuTabsStore.getState()
  store.add({ key: '/dashboard', label: 'Dashboard' })
  store.add({ key: '/dashboard', label: 'Dashboard Duplicate' })
  expect(useMenuTabsStore.getState().items).toHaveLength(1)
})

test('remove removes target tab', () => {
  const store = useMenuTabsStore.getState()
  store.add({ key: '/a', label: 'A' })
  store.add({ key: '/b', label: 'B' })
  store.remove('/a')
  expect(useMenuTabsStore.getState().items).toHaveLength(1)
  expect(useMenuTabsStore.getState().items[0].key).toBe('/b')
})

test('removeAll clears all tabs', () => {
  const store = useMenuTabsStore.getState()
  store.add({ key: '/a', label: 'A' })
  store.add({ key: '/b', label: 'B' })
  store.removeAll()
  expect(useMenuTabsStore.getState().items).toHaveLength(0)
})
