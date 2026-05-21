import { expect, test } from '@playwright/test'
import { gotoAsAdmin } from '../helpers/auth'

test.describe('用户管理', () => {
  test.beforeEach(async ({ page }) => {
    await gotoAsAdmin(page, '/account/user')
  })

  test('页面加载显示用户列表', async ({ page }) => {
    await expect(page.getByText('admin')).toBeVisible()
  })

  test('打开创建用户抽屉', async ({ page }) => {
    await page.getByRole('button', { name: '创建用户' }).click()
    await expect(page.locator('.ant-drawer-title').getByText('创建用户')).toBeVisible()
    await expect(page.getByPlaceholder('请输入用户名')).toBeVisible()
  })

  test('关闭创建用户抽屉', async ({ page }) => {
    await page.getByRole('button', { name: '创建用户' }).click()
    await expect(page.locator('.ant-drawer-title').getByText('创建用户')).toBeVisible()
    await page.getByRole('button', { name: '取 消' }).click()
    await expect(page.locator('.ant-drawer-title').getByText('创建用户')).not.toBeVisible()
  })
})
