import type { ReactNode } from 'react'
import { Tag } from 'antd'

export const ENTRY_LABEL_PLACEHOLDER = '$' + '{EntryLabel}'

interface ParsedTagTemplate {
  content: string
  color?: string
  bordered?: boolean
}

const blockedHtmlTags = new Set(['script', 'style', 'iframe', 'object', 'embed', 'link', 'meta', 'base'])

function parseTagTemplate(template?: string): ParsedTagTemplate | null {
  const source = template?.trim()
  if (!source) {
    return { content: ENTRY_LABEL_PLACEHOLDER }
  }

  const openEnd = source.indexOf('>')
  if (openEnd <= 1) {
    return null
  }
  const openTag = source.slice(1, openEnd).trim()
  const attrStart = openTag.search(/\s/)
  const tag = attrStart === -1 ? openTag : openTag.slice(0, attrStart)
  if (tag !== 'Tag') {
    return null
  }
  const closeTag = `</${tag}>`
  if (!source.endsWith(closeTag)) {
    return null
  }
  const attrs = attrStart === -1 ? '' : openTag.slice(attrStart)
  const content = source.slice(openEnd + 1, -closeTag.length)

  const parsed: ParsedTagTemplate = {
    content: content.trim(),
  }
  let consumed = ''
  const attrRegex = /\s+([a-zA-Z]+)=(?:"([^"]*)"|'([^']*)'|\{(true|false)\})/g
  for (let attrMatch = attrRegex.exec(attrs); attrMatch; attrMatch = attrRegex.exec(attrs)) {
    consumed += attrMatch[0]
    const [, name, doubleQuotedValue, singleQuotedValue, boolValue] = attrMatch
    const stringValue = doubleQuotedValue ?? singleQuotedValue
    if (name === 'color' && stringValue && /^[\w#-]+$/.test(stringValue)) {
      parsed.color = stringValue
      continue
    }
    if (name === 'bordered' && boolValue) {
      parsed.bordered = boolValue === 'true'
      continue
    }
    return null
  }

  if (consumed.trim() !== attrs.trim()) {
    return null
  }
  return parsed
}

function sanitizeLabelHtml(labelComponent: string, entryLabel: string) {
  if (typeof document === 'undefined') {
    return ''
  }
  const template = document.createElement('template')
  template.innerHTML = labelComponent.replaceAll(ENTRY_LABEL_PLACEHOLDER, entryLabel)

  const walk = document.createTreeWalker(template.content, NodeFilter.SHOW_ELEMENT)
  const blockedNodes: Element[] = []
  while (walk.nextNode()) {
    const element = walk.currentNode as Element
    const tagName = element.tagName.toLowerCase()
    if (blockedHtmlTags.has(tagName)) {
      blockedNodes.push(element)
      continue
    }
    for (const attr of Array.from(element.attributes)) {
      const attrName = attr.name.toLowerCase()
      const attrValue = attr.value.trim().toLowerCase()
      if (attrName.startsWith('on')) {
        element.removeAttribute(attr.name)
        continue
      }
      if ((attrName === 'href' || attrName === 'src') && /^(?:javascript|data):/.test(attrValue)) {
        element.removeAttribute(attr.name)
      }
    }
  }
  for (const node of blockedNodes) {
    node.remove()
  }
  return template.innerHTML
}

export function renderDictEntryLabel(labelComponent: string | undefined, entryLabel: string): ReactNode {
  if (!labelComponent?.trim()) {
    return entryLabel
  }
  const parsed = parseTagTemplate(labelComponent)
  if (!parsed) {
    const html = sanitizeLabelHtml(labelComponent, entryLabel)
    if (!html) {
      return entryLabel
    }
    return (
      // eslint-disable-next-line react-dom/no-dangerously-set-innerhtml
      <span dangerouslySetInnerHTML={{ __html: html }} />
    )
  }
  const children = parsed.content.includes(ENTRY_LABEL_PLACEHOLDER)
    ? parsed.content.replaceAll(ENTRY_LABEL_PLACEHOLDER, entryLabel)
    : parsed.content

  return (
    <Tag color={parsed.color} bordered={parsed.bordered}>
      {children}
    </Tag>
  )
}
