import { expect, test } from '@playwright/test'

test.describe('登录流程', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login', { waitUntil: 'commit' })
    await page.evaluate(() => {
      localStorage.clear()
      sessionStorage.clear()
    })
    await page.waitForLoadState('load')
  })

  test('账号密码登录成功', async ({ page }) => {
    await page.getByPlaceholder('用户名: admin or user').fill('admin')
    await page.getByPlaceholder('密码: ant.design').fill('123456')
    await page.getByRole('button', { name: '登 录' }).click()
    await page.waitForURL('**/')
    await expect(page).toHaveURL(/\/$/)
  })

  test('账号密码登录失败停留在登录页', async ({ page }) => {
    await page.getByPlaceholder('用户名: admin or user').fill('wrong')
    await page.getByPlaceholder('密码: ant.design').fill('wrong')
    await page.getByRole('button', { name: '登 录' }).click()
    await expect(page).toHaveURL(/.*login.*/)
  })

  test('切换手机号登录标签', async ({ page }) => {
    await page.getByRole('tab', { name: '手机号登录' }).click()
    await expect(page.getByPlaceholder('手机号')).toBeVisible()
    await expect(page.getByText('获取验证码')).toBeVisible()
  })

  test('未登录访问受保护页面重定向到登录', async ({ page }) => {
    await page.goto('/dashboard')
    await page.waitForURL('**/login**')
    await expect(page).toHaveURL(/.*login.*/)
  })
})
