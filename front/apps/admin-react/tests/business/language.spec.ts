import { expect, test } from '@playwright/test'
import { gotoAsAdmin } from '../helpers/auth'

test.describe('语言管理', () => {
  test.beforeEach(async ({ page }) => {
    await gotoAsAdmin(page, '/system/language')
  })

  test('页面加载显示语言类型', async ({ page }) => {
    await expect(page.getByText('简体中文')).toBeVisible()
    await expect(page.getByText('English')).toBeVisible()
  })

  test('打开新增语言类型弹窗', async ({ page }) => {
    await page.getByRole('button', { name: '新增语言' }).click()

    await expect(page.getByText('新增语言类型')).toBeVisible()
    await expect(page.getByPlaceholder('如: zh-CN')).toBeVisible()
  })
})
