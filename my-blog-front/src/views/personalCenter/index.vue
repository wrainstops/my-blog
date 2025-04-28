<template>
  <div class="main_Container">
    <div class="relative h-40vh mt-10px ml-5px mr-5px">
      <div class="absolute h-full w-full">
        <el-image :src="isDark ? darkBg : bg" fit="cover" class="h-full w-full" />
      </div>
      <div class="absolute bottom-80px left-200px">
        <div class="font-size-22px font-bold">{{ userName }}</div>
        <div>
          <span>发布文章数量:{{ stats?.articleNum || 0 }}</span>
          <span class="pl-18px">点赞数量:{{ stats?.likeNum || 0 }}</span>
          <span class="pl-18px">被回复数量:{{ stats?.beRepliedNum || 0 }}</span>
          <span class="pl-18px">被点赞数量:{{ stats?.beLikedNum || 0 }}</span>
        </div>
      </div>
      <div class="absolute bottom-16px left-24px">
        <el-avatar :size="150" src="https://img0.baidu.com/it/u=4044907637,4062909661&fm=253&fmt=auto&app=138&f=JPEG?w=1422&h=800" />
      </div>
    </div>

    <div class="content">
      <div v-for="(item, index) in articleList" :key="index" class="item">
        <my-card :row="item" @confirm="getNewList" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { isDark } from '@/composables/dark'
import { useWindowLoadMore } from '@/composables/loadMore'
import { getLocalUserInfo } from '@/utils/auth'
import myCard from './myCard.vue'
import darkBg from '@/assets/image/personal_bg_dark.png'
import bg from '@/assets/image/personal_bg.png'
import { ajaxLog } from '@/api/auth'
import { ajaxArticle } from '@/api/article'

const userName = getLocalUserInfo()?.name || ''
const stats: any = ref()
async function getStats() {
  const res: any = await ajaxLog.getStats()
  stats.value = res
}

const articleList: any = ref([])
const queryParameter = reactive({
  page: 0,
  pageSize: 20,
  key: '',
  descCreatedTime: true
})
const isEnd = ref(false) // 是否到底标记

async function getNewList() {
  queryParameter.page = 0
  articleList.value = []
  const res: any = await ajaxArticle.query(queryParameter)
  articleList.value.push(...res.content)
  queryParameter.page++
  isEnd.value = false
}
async function getMoreList() {
  if (!isEnd.value) {
    const res: any = await ajaxArticle.query(queryParameter)
    if (!(res.content && res.content.length)) {
      isEnd.value = true
      return Promise.reject(false) // 不return Promise会浏览器崩溃
    }
    articleList.value.push(...res.content)
    queryParameter.page++
  }
}
function handleScroll() {
  useWindowLoadMore(() => isEnd.value, getMoreList)
}

onMounted(async () => {
  getStats()
  await getNewList()
  const target = document.querySelector('.main_Container') as Element
  while (target.clientHeight <= window.innerHeight) {
    await getMoreList()
  }
  window.addEventListener('scroll', handleScroll)
})
onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped>
.content {
  column-count: 4;
  column-gap: 10px;
  padding: 0;
  margin-top: 10px;
}
.item {
  break-inside: avoid;
  margin-bottom: 10px;
}
</style>
