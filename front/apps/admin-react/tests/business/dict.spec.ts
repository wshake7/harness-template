import { expect, test } from '@playwright/test'

async function login(page: any) {
  await page.goto('/login', { waitUntil: 'commit' })
  await page.evaluate(() => {
    localStorage.clear()
    sessionStorage.clear()
  })
  await page.waitForLoadState('load')
  await page.getByPlaceholder('用户名: admin or user').fill('admin')
  await page.getByPlaceholder('密码: ant.design').fill('123456')
  await page.getByRole('button', { name: '登 录' }).click()
  await page.waitForURL('**/')
}

test.describe('字典管理', () => {
  test.beforeEach(async ({ page }) => {
    await login(page)
    await page.goto('/system/dict')
  })

  test('页面加载显示字典类型', async ({ page }) => {
    await expect(page.getByText('状态').first()).toBeVisible()
    await expect(page.getByText('性别').first()).toBeVisible()
  })
})
