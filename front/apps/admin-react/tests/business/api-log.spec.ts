import { expect, test } from '@playwright/test'
import { gotoAsAdmin } from '../helpers/auth'

test.describe('API日志', () => {
  test.beforeEach(async ({ page }) => {
    await gotoAsAdmin(page, '/logger/api/log')
  })

  test('页面加载显示 API 日志列表', async ({ page }) => {
    await expect(page.getByRole('row', { name: /\/api\/sys\/user\/list/ }).first()).toBeVisible()
  })

  test('打开 API 日志详情', async ({ page }) => {
    await page.getByRole('row', { name: /req-1000/ }).getByText('详情').click()

    const dialog = page.getByRole('dialog', { name: 'API日志详情' })
    await expect(dialog).toBeVisible()
    await expect(dialog).toContainText('req-1000')
  })
})
