import { expect, test } from '@playwright/test'
import { gotoAsAdmin } from '../helpers/auth'

test.describe('任务配置', () => {
  test.beforeEach(async ({ page }) => {
    await gotoAsAdmin(page, '/job/schedule')
  })

  test('页面加载显示任务配置列表', async ({ page }) => {
    await expect(page.getByRole('row', { name: /^cleanup / })).toBeVisible()
    await expect(page.getByRole('row', { name: /^backup / })).toBeVisible()
    await expect(page.getByRole('row', { name: /^sync / })).toBeVisible()
  })

  test('打开创建任务配置抽屉并校验必填项', async ({ page }) => {
    await page.getByRole('button', { name: '创建任务' }).click()
    await expect(page.locator('.ant-drawer-title').getByText('创建任务配置')).toBeVisible()

    await page.getByRole('button', { name: '保 存' }).click()

    await expect(page.getByText('请输入任务编码')).toBeVisible()
    await expect(page.getByText('请输入任务名称')).toBeVisible()
    await expect(page.getByText('请输入Workflow 类型')).toBeVisible()
  })
})
