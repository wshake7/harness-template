import { expect, test } from '@playwright/experimental-ct-react'
import { ErrorComponent } from './error'

test('ErrorComponent renders warning message', async ({ mount }) => {
  const component = await mount(<ErrorComponent onBack={() => { }} />)
  await expect(component).toContainText('There are some problems with your operation.')
})

test('ErrorComponent has Back Home button', async ({ mount }) => {
  const component = await mount(<ErrorComponent onBack={() => {}} />)
  await expect(component.getByRole('button', { name: 'Back Home' })).toBeVisible()
})

test('ErrorComponent calls onBack when button clicked', async ({ mount }) => {
  let clicked = false
  const component = await mount(<ErrorComponent onBack={() => { clicked = true }} />)
  await component.getByRole('button', { name: 'Back Home' }).click()
  expect(clicked).toBe(true)
})
