import { expect, test } from '@playwright/experimental-ct-react'
import { JsonCodeBlock } from './jsonCodeBlock'

test('JsonCodeBlock renders placeholder for empty content', async ({ mount }) => {
  const component = await mount(<JsonCodeBlock value="" />)

  await expect(component).toContainText('-')
})

test('JsonCodeBlock formats valid JSON content', async ({ mount }) => {
  const component = await mount(<JsonCodeBlock value={'{"name":"root","enabled":true}'} />)

  await expect(component).toContainText('"name": "root"')
  await expect(component).toContainText('"enabled": true')
})

test('JsonCodeBlock falls back to plain text for invalid JSON', async ({ mount }) => {
  const component = await mount(<JsonCodeBlock value="raw payload" />)

  await expect(component).toContainText('raw payload')
})

test('JsonCodeBlock treats HTML-like JSON values as inert text', async ({ mount, page }) => {
  const component = await mount(
    <JsonCodeBlock value={JSON.stringify({ body: '<script>window.__jsonCodeBlockXss = true</script>' })} />,
  )

  await expect(component).toContainText('<script>window.__jsonCodeBlockXss = true</script>')
  await expect(page.evaluate(() => Reflect.get(window, '__jsonCodeBlockXss'))).resolves.toBeUndefined()
})
