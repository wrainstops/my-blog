<template>
  <div class="dashLeftCard p-3 border-rd blog-border bg-[var(--bg-color-2)]">
    <div class="flex">
      <div class="flex items-center">
        <el-avatar :size="40" class="avatar">{{ getLastLetter(row?.authName) }}</el-avatar>
      </div>
      <div class="flex flex-col flex-1 p-10px gap-3px">
        <p>{{ row?.authName || '' }}</p>
        <p class="font-size-13px">{{ formatTime(row?.createdTime) }}</p>
        <p class="font-bold">{{ row?.title || '' }}</p>
      </div>
    </div>
    <el-divider />
    <div class="flex justify-around mt-6px">
      <div />
      <div class="flex items-center pointer">
        <el-icon v-if="row?.replyNum"><ChatLineSquare /></el-icon>
        <el-icon v-else><ChatSquare /></el-icon>
        <p>{{ row?.replyNum }}</p>
      </div>
      <div class="flex items-center pointer" @click="clickLike">
        <el-icon v-if="row?.likeFlag"><StarFilled /></el-icon>
        <el-icon v-else><Star /></el-icon>
        <p>{{ row?.likeNum }}</p>
      </div>
    </div>
    <login-component v-if="dialogShow" v-model:isShow="dialogShow" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Star, StarFilled, ChatSquare, ChatLineSquare } from '@element-plus/icons-vue'
import { useLogin } from '@/store/modules/loginMark'
import { formatTime, getLastLetter } from '@/utils/common'
import { ajaxArticle } from '@/api/article'

const props = defineProps({
  row: {
    type: Object,
    default: () => {}
  }
})
const store = useLogin()
const dialogShow = ref(false)

function clickLike() {
  if (!store.loginMark) {
    dialogShow.value = !dialogShow.value
    return
  }
  const opeType = props.row?.likeFlag ? 'cancelLike' : 'addLike'
  ajaxArticle[opeType]({ articleId: props.row.ID})
    .then(() => {
      props.row.likeFlag = !props.row?.likeFlag
      if (props.row.likeFlag) {
        props.row.likeNum++
      } else {
        props.row.likeNum--
      }
    })
}
</script>

<style scoped lang="scss">
p {
  margin: 0;
  text-align: start;
}
.el-divider--horizontal {
  margin: 3px 0;
  border-top: 1px dashed #646464;
}
</style>
