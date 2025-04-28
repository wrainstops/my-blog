<template>
  <div id="aplayer" />
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import APlayer from "aplayer"
import "aplayer/dist/APlayer.min.css"

// 歌曲列表
const audio = ref([
  {
    name: "光るなら", // 歌曲名称
    artist: "Goose house", // 歌曲作者
    // 歌曲url地址
    url: "/music/lieInApr.mp3",
    // 歌曲封面
    cover: "/music/lieInApr.jpg",
    // 歌词
    lrc: "/music/lieInApr.lrc",
    // 主题
    theme: "#fff",
  }
])
// 其他配置信息
const info = ref({
  fixed: true, // 开启吸底模式
  listFolded: true, // 折叠歌曲列表
  autoplay: false, // 开启自动播放
  preload: "auto", // 自动预加载歌曲
  loop: "all", // 播放循环模式、all全部循环 one单曲循环 none只播放一次
  order: "list", //  播放模式，list列表播放, random随机播放
  lrcType: 3, //使用lrc歌词
  volume: 0.1, // 播放音量
})
const initAudio = () => {
  // 创建一个音乐播放器实例，并挂载到DOM上，同时进行相关配置
  new APlayer({
    container: document.getElementById("aplayer"),
    audio: audio.value, // 音乐信息
    ...info.value, // 其他配置信息
  })
}
onMounted(() => {
  // 初始化播放器
  initAudio()
})
</script>

<style scoped lang="scss">
:deep(.aplayer-info) {
  border-top: 0 !important;
}
</style>
