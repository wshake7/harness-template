import { expect, test } from '@playwright/test'
import { gotoAsAdmin } from '../helpers/auth'

test.describe('Dashboard 页面', () => {
  test.beforeEach(async ({ page }) => {
    await gotoAsAdmin(page, '/dashboard')
  })

  test('页面加载显示欢迎信息', async ({ page }) => {
    await expect(page.getByText('Test Action12 Hello "/_app/dashboard"!')).toBeVisible()
  })

  test('侧边栏显示菜单项', async ({ page }) => {
    const sider = page.getByTestId('pro-layout-sider')
    await expect(sider.getByText('系统管理')).toBeVisible()
    await expect(sider.getByText('权限管理')).toBeVisible()
    await expect(sider.getByText('数据看板').first()).toBeVisible()
  })

  test('点击菜单导航到用户管理', async ({ page }) => {
    await page.getByTestId('pro-layout-sider').getByText('权限管理').click()
    await page.getByTestId('pro-layout-sider').getByText('用户管理').click()
    await page.waitForURL('**/account/user**')
    await expect(page).toHaveURL(/\/account\/user/)
  })
})
