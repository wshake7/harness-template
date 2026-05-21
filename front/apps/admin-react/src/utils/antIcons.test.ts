import { expect, test } from 'vite-plus/test'
import { antIconNames, antIconNamesByStyle, getAntIconStyle } from './antIcons'

test('getAntIconStyle returns Outlined for Outlined suffix', () => {
  expect(getAntIconStyle('HomeOutlined')).toBe('Outlined')
})

test('getAntIconStyle returns Filled for Filled suffix', () => {
  expect(getAntIconStyle('HomeFilled')).toBe('Filled')
})

test('getAntIconStyle returns TwoTone for TwoTone suffix', () => {
  expect(getAntIconStyle('HomeTwoTone')).toBe('TwoTone')
})

test('getAntIconStyle returns undefined for names without valid suffix', () => {
  expect(getAntIconStyle('Home')).toBeUndefined()
  expect(getAntIconStyle('')).toBeUndefined()
})

test('all icon names end with valid style suffix', () => {
  const validSuffixes = ['Outlined', 'Filled', 'TwoTone']
  for (const name of antIconNames) {
    const hasValidSuffix = validSuffixes.some(suffix => name.endsWith(suffix))
    expect(hasValidSuffix, `Icon "${name}" has invalid suffix`).toBe(true)
  }
})

test('antIconNamesByStyle categories are mutually exclusive', () => {
  const outlined = new Set(antIconNamesByStyle.Outlined)
  const filled = new Set(antIconNamesByStyle.Filled)
  const twoTone = new Set(antIconNamesByStyle.TwoTone)

  for (const name of filled) {
    expect(outlined.has(name), `"${name}" in both Filled and Outlined`).toBe(false)
  }
  for (const name of twoTone) {
    expect(outlined.has(name), `"${name}" in both TwoTone and Outlined`).toBe(false)
    expect(filled.has(name), `"${name}" in both TwoTone and Filled`).toBe(false)
  }
})

test('total icons equal sum of style categories', () => {
  const total
    = antIconNamesByStyle.Outlined.length
      + antIconNamesByStyle.Filled.length
      + antIconNamesByStyle.TwoTone.length
  expect(antIconNames.length).toBe(total)
})
