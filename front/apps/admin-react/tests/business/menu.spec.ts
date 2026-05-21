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

test.describe('菜单管理', () => {
  test.beforeEach(async ({ page }) => {
    await login(page)
    await page.goto('/system/resource/menu')
  })

  test('页面加载显示菜单列表', async ({ page }) => {
    await expect(page.getByText('系统管理').first()).toBeVisible()
    await expect(page.getByText('字典管理').first()).toBeVisible()
  })

  test('打开创建菜单抽屉', async ({ page }) => {
    await page.getByRole('button', { name: '创建菜单' }).click()
    await expect(page.locator('.ant-drawer-title').getByText('创建菜单')).toBeVisible()
  })
})
