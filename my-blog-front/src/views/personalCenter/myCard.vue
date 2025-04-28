<template>
  <div class="pl-5px pr-5px">
    <el-card shadow="hover" class="border-rd-20px">
      <template #header>
        <div class="flex justify-between">
          <p class="font-size-13px m-0">{{ formatTime(row?.createdTime) }}</p>
          <el-icon class="pointer" @click="handleDelete"><Delete /></el-icon>
        </div>
      </template>

      <div class="pointer" @click="handleClickCard">
        <h3 class="font-bold">{{ row?.title || '' }}</h3>
        <p v-html="row.content" class="content3line m-0"></p>
      </div>

      <template #footer>
        <div class="flex gap-10px">
          <div class="flex items-center pointer">
            <el-icon v-if="row?.replyNum"><ChatDotSquare /></el-icon>
            <el-icon v-else><ChatSquare /></el-icon>
            <p class="m-0">{{ row?.replyNum }}</p>
          </div>
          <div class="flex items-center pointer">
            <el-icon v-if="row?.likeFlag"><StarFilled /></el-icon>
            <el-icon v-else><Star /></el-icon>
            <p class="m-0">{{ row?.likeNum }}</p>
          </div>
        </div>
      </template>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete, Star, StarFilled, ChatSquare, ChatDotSquare } from '@element-plus/icons-vue'
import { formatTime } from '@/utils/common'
import { ajaxArticle } from '@/api/article'

const props = defineProps({
  row: {
    type: Object,
    default: () => {}
  }
})
const emit = defineEmits(['confirm'])

const router = useRouter()
function handleClickCard() {
  router.push(`/details/${ props.row.ID }`)
}
async function handleDelete() {
  await ElMessageBox.confirm(
    '确认删除此条博客吗?',
    '确认',
    { showClose: false }
  )
  await ajaxArticle.deleteArticle(props.row.ID)
  ElMessage.success('删除成功!')
  emit('confirm')
}

const randomLine = Math.floor(Math.random() * 3) + 1 // 向下取整，生成随机数 1~3
</script>

<style scoped lang="scss">
.content3line {
  display: -webkit-box;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  line-clamp: v-bind(randomLine);
  -webkit-line-clamp: v-bind(randomLine);
}
:deep(.el-card__footer) {
  padding: 10px var(--el-card-padding);
}
</style>
