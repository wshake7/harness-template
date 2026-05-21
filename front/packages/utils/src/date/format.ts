import dayjs from 'dayjs'

export function formatDateYYYYMMDDHHmmss(value?: string) {
  if (!value) {
    return '-'
  }
  const date = dayjs(value)
  if (!date.isValid()) {
    return value
  }
  return date.format('YYYY-MM-DD HH:mm:ss')
}
