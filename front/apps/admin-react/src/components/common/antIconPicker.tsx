import type { AntIconStyle } from '~/utils/antIcons'
import { Button, Empty, Input, Segmented, Space, Tooltip } from 'antd'
import { useEffect, useMemo, useRef, useState } from 'react'
import { AntIcon, antIconNamesByStyle, getAntIconStyle } from '~/utils/antIcons'

const iconStyleOptions: { title: string, value: AntIconStyle, icon: string }[] = [
  { title: '线框风格', value: 'Outlined', icon: 'BorderOutlined' },
  { title: '实底风格', value: 'Filled', icon: 'AppstoreFilled' },
  { title: '双色风格', value: 'TwoTone', icon: 'PieChartTwoTone' },
]

const iconGridHeight = 220
const iconGridPadding = 8
const iconGridGap = 6
const iconButtonSize = 32
const iconGridRowHeight = iconButtonSize + iconGridGap
const iconGridColumnWidth = iconButtonSize + iconGridGap
const iconGridOverscanRows = 3

interface AntIconPickerProps {
  value?: string
  onChange?: (value?: string) => void
}

interface IconGridProps extends AntIconPickerProps {
  iconNames: string[]
}

function IconGrid({ iconNames, value, onChange }: IconGridProps) {
  const containerRef = useRef<HTMLDivElement>(null)
  const [scrollTop, setScrollTop] = useState(0)
  const [columnCount, setColumnCount] = useState(8)

  useEffect(() => {
    const container = containerRef.current
    if (!container) {
      return
    }

    const updateColumnCount = () => {
      const contentWidth = container.clientWidth - iconGridPadding * 2
      requestAnimationFrame(() => {
        setColumnCount(Math.max(1, Math.floor((contentWidth + iconGridGap) / iconGridColumnWidth)))
      })
    }

    updateColumnCount()
    const resizeObserver = new ResizeObserver(updateColumnCount)
    resizeObserver.observe(container)

    return () => {
      resizeObserver.disconnect()
    }
  }, [])

  useEffect(() => {
    const container = containerRef.current
    if (!container) {
      return
    }

    container.scrollTop = 0
    requestAnimationFrame(() => {
      setScrollTop(0)
    })
  }, [iconNames])

  const totalRows = Math.ceil(iconNames.length / columnCount)
  const startRow = Math.max(0, Math.floor(scrollTop / iconGridRowHeight) - iconGridOverscanRows)
  const visibleRowCount = Math.ceil(iconGridHeight / iconGridRowHeight) + iconGridOverscanRows * 2
  const endRow = Math.min(totalRows, startRow + visibleRowCount)
  const visibleIconNames = iconNames.slice(startRow * columnCount, endRow * columnCount)
  const topSpacerHeight = startRow * iconGridRowHeight
  const bottomSpacerHeight = Math.max(0, (totalRows - endRow) * iconGridRowHeight)

  return (
    <div
      ref={containerRef}
      style={{
        border: '1px solid rgba(5, 5, 5, 0.06)',
        borderRadius: 6,
        height: iconGridHeight,
        overflowY: 'auto',
        padding: iconGridPadding,
      }}
      onScroll={event => setScrollTop(event.currentTarget.scrollTop)}
    >
      {iconNames.length === 0
        ? <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} />
        : (
            <>
              <div style={{ height: topSpacerHeight }} />
              <div
                style={{
                  display: 'grid',
                  gap: iconGridGap,
                  gridTemplateColumns: `repeat(${columnCount}, ${iconButtonSize}px)`,
                }}
              >
                {visibleIconNames.map(name => (
                  <Tooltip key={name} title={name}>
                    <Button
                      aria-label={name}
                      icon={<AntIcon name={name} />}
                      style={{ height: iconButtonSize, width: iconButtonSize }}
                      type={value === name ? 'primary' : 'text'}
                      onClick={() => onChange?.(name)}
                    />
                  </Tooltip>
                ))}
              </div>
              <div style={{ height: bottomSpacerHeight }} />
            </>
          )}
    </div>
  )
}

export default function AntIconPicker({ value, onChange }: AntIconPickerProps) {
  const [iconStyle, setIconStyle] = useState<AntIconStyle>(() => getAntIconStyle(value) ?? 'Outlined')
  const [keyword, setKeyword] = useState('')
  const selectedStyle = getAntIconStyle(value) ?? iconStyle
  const iconNames = useMemo(() => {
    const names = antIconNamesByStyle[selectedStyle]
    const normalizedKeyword = keyword.trim().toLowerCase()
    if (!normalizedKeyword) {
      return names
    }

    return names.filter(name => name.toLowerCase().includes(normalizedKeyword))
  }, [keyword, selectedStyle])

  return (
    <Space orientation="vertical" size={8} style={{ width: '100%' }}>
      <Space size={8} style={{ width: '100%' }}>
        <Segmented<AntIconStyle>
          value={selectedStyle}
          options={iconStyleOptions.map(option => ({
            value: option.value,
            label: (
              <Tooltip title={option.title}>
                <span>
                  <AntIcon name={option.icon} />
                </span>
              </Tooltip>
            ),
          }))}
          onChange={(nextStyle) => {
            setIconStyle(nextStyle)
            if (value && getAntIconStyle(value) !== nextStyle) {
              onChange?.(undefined)
            }
          }}
        />
        <Input
          allowClear
          placeholder="搜索图标"
          value={keyword}
          onChange={event => setKeyword(event.target.value)}
        />
      </Space>
      <IconGrid iconNames={iconNames} value={value} onChange={onChange} />
    </Space>
  )
}
