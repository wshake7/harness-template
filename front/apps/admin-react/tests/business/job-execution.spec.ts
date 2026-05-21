import { expect, test } from '@playwright/test'
import { gotoAsAdmin } from '../helpers/auth'

test.describe('执行记录', () => {
  test.beforeEach(async ({ page }) => {
    await gotoAsAdmin(page, '/job/execution')
  })

  test('页面加载显示执行记录列表', async ({ page }) => {
    await expect(page.getByText('workflow-1000')).toBeVisible()
  })

  test('打开执行记录详情', async ({ page }) => {
    await page.getByRole('button', { name: '详情' }).first().click()

    const dialog = page.getByRole('dialog', { name: '执行记录详情' })
    await expect(dialog).toBeVisible()
    await expect(dialog).toContainText('workflow-1000')
    await expect(dialog).toContainText('cleanup')
  })
})
