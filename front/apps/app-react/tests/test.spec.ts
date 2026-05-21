import { expect, test } from '@playwright/test'

test('test', async ({ page }) => {
  await page.goto('http://localhost:3000/login')
  await page.getByRole('textbox', { name: '用户名: admin or user' }).fill('admin')
  await page.getByRole('textbox', { name: '密码: ant.design' }).fill('123456')
  await page.getByRole('button', { name: '登 录' }).click()

  // 登录成功后，等待跳转到 dashboard 页面
  await page.waitForURL('**/dashboard**')

  // 断言 dashboard 页面中是否含有 hello dashboard 几个字
  await expect(page.getByText('Hello "/_app/dashboard"!')).toBeVisible()
})
