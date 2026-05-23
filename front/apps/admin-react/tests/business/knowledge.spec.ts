import { expect, test } from '@playwright/test'
import { gotoAsAdmin } from '../helpers/auth'

test.describe('知识库管理', () => {
  test.describe('集合管理', () => {
    test.beforeEach(async ({ page }) => {
      await gotoAsAdmin(page, '/knowledge/collection')
    })

    test('页面加载显示知识库集合列表', async ({ page }) => {
      await expect(page.getByText('知识库集合')).toBeVisible()
      await expect(page.getByRole('button', { name: '创建集合' })).toBeVisible()
    })

    test('打开创建集合抽屉并校验必填项', async ({ page }) => {
      await page.getByRole('button', { name: '创建集合' }).click()
      await expect(page.locator('.ant-drawer-title').getByText('创建集合')).toBeVisible()

      await page.getByRole('button', { name: '保 存' }).click()

      await expect(page.getByText('集合名称不能为空')).toBeVisible()
      await expect(page.getByText('显示名称不能为空')).toBeVisible()
      await expect(page.getByText('Embedding模型不能为空')).toBeVisible()
      await expect(page.getByText('向量维度不能为空')).toBeVisible()
    })

    test('搜索框可输入并触发查询', async ({ page }) => {
      const searchBox = page.locator('input[placeholder*="搜索集合名称"]').first()
      await expect(searchBox).toBeVisible()
      await searchBox.fill('test')
      await expect(searchBox).toHaveValue('test')
    })
  })

  test.describe('文档管理', () => {
    test('文档页面显示集合信息提示', async ({ page }) => {
      await gotoAsAdmin(page, '/knowledge/document')

      await expect(page.getByText('所属集合：')).toBeVisible()
      await expect(page.getByText('未选择')).toBeVisible()
      await expect(page.getByText('请从集合列表选择一个集合查看文档')).toBeVisible()
    })

    test('从集合页面跳转到文档页面', async ({ page }) => {
      await gotoAsAdmin(page, '/knowledge/collection')

      // 创建测试集合
      await page.getByRole('button', { name: '创建集合' }).click()
      await page.locator('.ant-drawer-body input#collectionName').fill('test-jump-collection')
      await page.locator('.ant-drawer-body input#displayName').fill('Test Jump Collection')
      await page.locator('.ant-drawer-body input#embeddingModel').fill('text-embedding-3-small')
      await page.locator('.ant-drawer-body input#vectorDimension').fill('1536')
      await page.getByRole('button', { name: '保 存' }).click()

      // 等待抽屉关闭
      await expect(page.locator('.ant-drawer:visible')).not.toBeVisible()

      // 点击查看文档按钮（新创建的行）
      const row = page.getByRole('row', { name: /test-jump-collection/ }).first()
      await expect(row).toBeVisible()
      await row.getByRole('button', { name: '查看文档' }).click()

      // 验证跳转到文档页面
      await page.waitForURL('**/knowledge/document?collectionId=*')
      await expect(page.getByText('Test Jump Collection')).toBeVisible()
    })
  })
})
