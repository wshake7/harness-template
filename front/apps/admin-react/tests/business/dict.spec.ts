import { expect, test } from '@playwright/test'
import { gotoAsAdmin } from '../helpers/auth'

test.describe('字典管理', () => {
  test.beforeEach(async ({ page }) => {
    await gotoAsAdmin(page, '/system/dict')
  })

  test('页面加载显示字典类型', async ({ page }) => {
    await expect(page.getByText('状态').first()).toBeVisible()
    await expect(page.getByText('性别').first()).toBeVisible()
  })
})
