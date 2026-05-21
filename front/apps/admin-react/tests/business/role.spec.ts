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

test.describe('角色管理', () => {
  test.beforeEach(async ({ page }) => {
    await login(page)
    await page.goto('/account/role')
  })

  test('页面加载显示角色列表', async ({ page }) => {
    await expect(page.getByText('超级管理员').first()).toBeVisible()
    await expect(page.getByText('管理员').first()).toBeVisible()
    await expect(page.getByText('普通用户').first()).toBeVisible()
  })

  test('打开创建角色抽屉', async ({ page }) => {
    await page.getByRole('button', { name: '创建角色' }).click()
    await expect(page.locator('.ant-drawer-title').getByText('创建角色')).toBeVisible()
    await expect(page.getByPlaceholder('请输入角色名称')).toBeVisible()
  })
})
