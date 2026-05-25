import { beforeEach, expect, test } from 'vite-plus/test'
import { useResourceMenuStore } from './resourceMenu'

beforeEach(() => {
  useResourceMenuStore.setState({ dynamicMenuTree: [] })
  useResourceMenuStore.persist.clearStorage()
})

test('setDynamicMenuTree updates tree', () => {
  const tree = [
    {
      id: 1,
      name: 'Test',
      path: '/test',
      menuType: 'MENU' as const,
      redirect: '',
      component: '',
      icon: '',
      order: 0,
      sortOrder: 0,
      hidden: false,
      authorities: [],
      isUrl: false,
      children: [],
    },
  ]
  useResourceMenuStore.getState().setDynamicMenuTree(tree)
  expect(useResourceMenuStore.getState().dynamicMenuTree).toHaveLength(1)
  expect(useResourceMenuStore.getState().dynamicMenuTree[0].name).toBe('Test')
})

test('clear resets tree to empty', () => {
  const tree = [
    {
      id: 1,
      name: 'Test',
      path: '/test',
      menuType: 'MENU' as const,
      redirect: '',
      component: '',
      icon: '',
      order: 0,
      sortOrder: 0,
      hidden: false,
      authorities: [],
      isUrl: false,
      children: [],
    },
  ]
  useResourceMenuStore.getState().setDynamicMenuTree(tree)
  useResourceMenuStore.getState().clear()
  expect(useResourceMenuStore.getState().dynamicMenuTree).toHaveLength(0)
})
