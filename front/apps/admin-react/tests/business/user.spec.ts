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

test.describe('用户管理', () => {
  test.beforeEach(async ({ page }) => {
    await login(page)
    await page.goto('/account/user')
  })

  test('页面加载显示用户列表', async ({ page }) => {
    await expect(page.getByText('admin')).toBeVisible()
  })

  test('打开创建用户抽屉', async ({ page }) => {
    await page.getByRole('button', { name: '创建用户' }).click()
    await expect(page.locator('.ant-drawer-title').getByText('创建用户')).toBeVisible()
    await expect(page.getByPlaceholder('请输入用户名')).toBeVisible()
  })

  test('关闭创建用户抽屉', async ({ page }) => {
    await page.getByRole('button', { name: '创建用户' }).click()
    await expect(page.locator('.ant-drawer-title').getByText('创建用户')).toBeVisible()
    await page.getByRole('button', { name: '取 消' }).click()
    await expect(page.locator('.ant-drawer-title').getByText('创建用户')).not.toBeVisible()
  })
})
