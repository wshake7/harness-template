import type { Page } from '@playwright/test'

export async function clearAuthState(page: Page) {
  await page.goto('/login', { waitUntil: 'commit' })
  await page.evaluate(() => {
    localStorage.clear()
    sessionStorage.clear()
  })
  await page.waitForLoadState('load')
}

export async function loginAsAdmin(page: Page) {
  await clearAuthState(page)
  await page.getByPlaceholder('用户名: admin or user').fill('admin')
  await page.getByPlaceholder('密码: ant.design').fill('123456')
  await page.getByRole('button', { name: '登 录' }).click()
  await page.waitForURL('**/')
}

export async function gotoAsAdmin(page: Page, path: string) {
  await loginAsAdmin(page)
  await page.goto(path)
  await page.waitForLoadState('load')
}
