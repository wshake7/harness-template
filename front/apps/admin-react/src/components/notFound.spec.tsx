import { expect, test } from '@playwright/experimental-ct-react'
import { NotFoundComponent } from './notFound'

test('NotFoundComponent renders 404 status', async ({ mount }) => {
  const component = await mount(<NotFoundComponent onBack={() => {}} />)
  await expect(component).toContainText('404')
  await expect(component).toContainText('Sorry, the page you visited does not exist.')
})

test('NotFoundComponent has Back Home button', async ({ mount }) => {
  const component = await mount(<NotFoundComponent onBack={() => {}} />)
  await expect(component.getByRole('button', { name: 'Back Home' })).toBeVisible()
})

test('NotFoundComponent calls onBack when button clicked', async ({ mount }) => {
  let clicked = false
  const component = await mount(<NotFoundComponent onBack={() => { clicked = true }} />)
  await component.getByRole('button', { name: 'Back Home' }).click()
  expect(clicked).toBe(true)
})
