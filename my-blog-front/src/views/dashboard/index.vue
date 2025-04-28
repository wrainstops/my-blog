<template>
  <div>
    <div style="height: 80px;">
      <el-image class="h-full w-full" :src="isDark ? darkBgT : bgT" fit="cover"></el-image>
    </div>

    <div class="w-full sticky pos-top-61px bg-[var(--bg-color-2)] z-999">
      <div class="flex items-center gap-20px ma h-40px max-w-1122px">
        <el-button :icon="queryParameter.descCreatedTime ? SortDown : SortUp" @click="handleChangeSort">发布时间</el-button>
        <el-input
          v-model="queryParameter.key"
          placeholder="搜索博客"
          clearable
          class="w-400px!"
          @keyup.enter="getNewList"
        >
          <template #append>
            <el-button :icon="Search" @click="getNewList" />
          </template>
        </el-input>
      </div>
    </div>

    <div class="main_Container">
      <div class="flex justify-between">
        <div class="w-28%">
          <div class="mt-10px" v-for="item in articleList" :key="item.authId">
            <left-card
              :row="item"
              class="pointer"
              :class="selectId === item.ID ? 'selectcard' : ''"
              @click="handleClickCard(item.ID)"
            />
          </div>
          <div class="text-align-center">
            <span v-if="isEnd">到底了...</span>
            <span v-else>下滑加载更多</span>
          </div>
        </div>
        <div v-if="articleList.length" class="w-68% h-[calc(100vh-111px)] mt-10px overflow-auto sticky pos-top-111px">
          <right-detail />
        </div>
      </div>

      <el-tooltip
        content="返回顶部"
        effect="customized"
        placement="left"
      >
        <el-backtop :right="100" :bottom="200" />
      </el-tooltip>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, provide, onMounted, onUnmounted } from 'vue'
import { SortDown, SortUp, Search } from '@element-plus/icons-vue'
import { isDark } from '@/composables/dark'
import { useWindowLoadMore } from '@/composables/loadMore'
import leftCard from './leftCard.vue'
import rightDetail from './rightDetail.vue'
import darkBgT from '@/assets/image/dash_t_dark.jpg'
import bgT from '@/assets/image/dash_t.jpg'
import { ajaxArticle } from '@/api/article'

const articleList: any = ref([])
const queryParameter = reactive({
  page: 0,
  pageSize: 10,
  key: '',
  descCreatedTime: true
})

const isEnd = ref(false) // 是否到底标记

function handleChangeSort() {
  queryParameter.descCreatedTime = !queryParameter.descCreatedTime
  getNewList()
}

async function getNewList() {
  queryParameter.page = 0
  articleList.value = []
  const res: any = await ajaxArticle.query(queryParameter)
  articleList.value.push(...res.content)
  queryParameter.page++
  isEnd.value = false
  if (!selectId.value) {
    selectId.value = res.content[0]?.ID
  }
}
async function getMoreList() {
  if (!isEnd.value) {
    const res: any = await ajaxArticle.query(queryParameter)
    if (!(res.content && res.content.length)) {
      isEnd.value = true
      return
    }
    articleList.value.push(...res.content)
    queryParameter.page++
  }
}
provide('getNewList', getNewList)

const selectId = ref(0)
function handleClickCard(val: number) {
  selectId.value = val
}
provide('id', selectId)

function handleScroll() {
  useWindowLoadMore(() => isEnd.value, getMoreList)
}
onMounted(() => {
  getNewList()
  window.addEventListener('scroll', handleScroll)
})
onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped lang="scss">
.selectcard {
  box-shadow: 0px 0px 16px rgba(84, 84, 150, 0.805);
}
</style>
