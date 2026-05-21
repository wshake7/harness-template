import type { ComponentType, CSSProperties } from 'react'
import * as AntIcons from '@ant-design/icons'

interface AntIconProps {
  className?: string
  style?: CSSProperties
}

const antIconNamePattern = /(?:Outlined|Filled|TwoTone)$/
export type AntIconStyle = 'Outlined' | 'Filled' | 'TwoTone'

const antIconMap = Object.fromEntries(
  Object.entries(AntIcons)
    .filter(([name]) => antIconNamePattern.test(name)),
) as Record<string, ComponentType<AntIconProps>>

export const antIconNames = Object.keys(antIconMap).sort((a, b) => a.localeCompare(b))

export const antIconNamesByStyle: Record<AntIconStyle, string[]> = {
  Outlined: antIconNames.filter(name => name.endsWith('Outlined')),
  Filled: antIconNames.filter(name => name.endsWith('Filled')),
  TwoTone: antIconNames.filter(name => name.endsWith('TwoTone')),
}

export function getAntIconStyle(name?: string): AntIconStyle | undefined {
  if (!name) {
    return undefined
  }

  if (name.endsWith('Filled')) {
    return 'Filled'
  }

  if (name.endsWith('TwoTone')) {
    return 'TwoTone'
  }

  if (name.endsWith('Outlined')) {
    return 'Outlined'
  }

  return undefined
}

export function AntIcon({ name, className, style }: AntIconProps & { name?: string }) {
  if (!name) {
    return null
  }

  const Icon = antIconMap[name]
  if (!Icon) {
    return null
  }

  return <Icon className={className} style={style} />
}

export function renderAntIcon(name?: string) {
  return name ? <AntIcon name={name} /> : undefined
}
