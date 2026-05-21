import { expect, test } from '@playwright/test'
import { gotoAsAdmin } from '../helpers/auth'

test.describe('登录日志', () => {
  test.beforeEach(async ({ page }) => {
    await gotoAsAdmin(page, '/logger/login/log')
  })

  test('页面加载显示登录日志列表', async ({ page }) => {
    await expect(page.getByText('admin').first()).toBeVisible()
  })

  test('打开登录日志详情', async ({ page }) => {
    await page.getByRole('row', { name: /^1 admin/ }).getByText('详情').click()

    const dialog = page.getByRole('dialog', { name: '登录日志详情' })
    await expect(dialog).toBeVisible()
    await expect(dialog).toContainText('登录账号')
  })
})
