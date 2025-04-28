import { debounce } from 'lodash-es'

/**
 * useDivLoadMore 滚动到元素底部调用方法
 * @param event 滚动的dom元素 {UIEvent}
 * @param checkEnd 判断是否需要调用func的方法 {() => boolean}
 * @param func 滚动到底部时调用的方法 {() => void}
 */
export const useDivLoadMore = debounce((event: UIEvent, checkEnd: () => boolean, func: () => void) => {
  const container: HTMLDivElement = event.target as HTMLDivElement
  const scrollTop = container.scrollTop // 已滚动高度
  const scrollHeight = container.scrollHeight // 滚动内容高度
  const clientHeight = container.clientHeight // 容器高度

  if (scrollTop + clientHeight >= scrollHeight - 5) { // 预留5的缓冲高度
    if (!checkEnd()) {
      func()
    }
  }
}, 300)


/**
 * useWindowLoadMore 滚动到窗口底部调用方法
 * @param checkEnd 判断是否需要调用func的方法 {() => boolean}
 * @param func 滚动到底部时调用的方法 {() => void}
 */
export const useWindowLoadMore = debounce((checkEnd: () => boolean, func: () => void) => {
  const scrollTop = window.scrollY
  const scrollHeight = window.innerHeight
  const clientHeight = document.documentElement.scrollHeight

  if (scrollTop + scrollHeight >= clientHeight - 5) { // 预留5的缓冲高度
    if (!checkEnd()) {
      func()
    }
  }
}, 300)
