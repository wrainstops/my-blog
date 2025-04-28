/// <reference types="vite/client" />

// 加下两部分代码，可导入vue文件 --- import wave from '@/components/wave.vue'
declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}
declare module 'aplayer'
declare interface Window {
  PIXI?: any,
  existLoading: boolean
}
