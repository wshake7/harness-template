import { expect, test } from '@playwright/experimental-ct-react'
import { ErrorComponent } from './error'

test('ErrorComponent', async ({ mount }) => {
  const component = await mount(<ErrorComponent onBack={() => { }} />)
  await expect(component).toContainText('There are some problems with your operation.')
})
