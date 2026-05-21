import { expect, test } from '@playwright/test'
import { gotoAsAdmin } from '../helpers/auth'

test.describe('API资源管理', () => {
  test.beforeEach(async ({ page }) => {
    await gotoAsAdmin(page, '/system/resource/api')
  })

  test('页面加载显示 API 资源列表', async ({ page }) => {
    await expect(page.getByText('/api/sys/user/list')).toBeVisible()
  })

  test('打开创建 API 资源抽屉并校验路径必填', async ({ page }) => {
    await page.getByRole('button', { name: '创建API资源' }).click()
    await expect(page.locator('.ant-drawer-title').getByText('创建API资源')).toBeVisible()

    await page.getByRole('button', { name: '保 存' }).click()

    await expect(page.getByText('请输入接口路径')).toBeVisible()
  })
})
