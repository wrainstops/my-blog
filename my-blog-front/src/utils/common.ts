import moment from 'moment'

export function formatTime(time: any, formatType: string = 'YYYY-MM-DD HH:mm:ss') {
  return time ? moment(time).format(formatType) : ''
}

export function hasWorth(data: any) {
  return Boolean(data || data === 0)
}

// 获取字符串最后一个字符
export function getLastLetter(str: string) {
  return str ? str.slice(-1) : ''
}
