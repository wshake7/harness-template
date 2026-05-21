import { expect, test } from '@playwright/test'
import { gotoAsAdmin } from '../helpers/auth'

test.describe('菜单管理', () => {
  test.beforeEach(async ({ page }) => {
    await gotoAsAdmin(page, '/system/resource/menu')
  })

  test('页面加载显示菜单列表', async ({ page }) => {
    await expect(page.getByText('系统管理').first()).toBeVisible()
    await expect(page.getByText('字典管理').first()).toBeVisible()
  })

  test('打开创建菜单抽屉', async ({ page }) => {
    await page.getByRole('button', { name: '创建菜单' }).click()
    await expect(page.locator('.ant-drawer-title').getByText('创建菜单')).toBeVisible()
  })
})
