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

test.describe('Dashboard 页面', () => {
  test.beforeEach(async ({ page }) => {
    await login(page)
    await page.goto('/dashboard')
  })

  test('页面加载显示欢迎信息', async ({ page }) => {
    await expect(page.getByText('Hello "/_app/dashboard"!')).toBeVisible()
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
