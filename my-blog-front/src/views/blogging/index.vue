<template>
  <div class="main_Container">
    <div class="pt-30px max-w-600px ma">
      <el-input
        v-model="title"
        placeholder="博客主题"
        clearable
        maxlength="50"
        show-word-limit
      />
      <rich-txt ref="richTxtRef" class="mt-20px" />
      <div class="mt-10px text-end">
        <el-button @click="submit">发布</el-button>
      </div>
    </div>
    <login-component v-if="dialogShow" v-model:isShow="dialogShow" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import richTxt from '@/components/richTxt.vue'
import { useLogin } from '@/store/modules/loginMark'
import { ajaxArticle } from '@/api/article'

const router = useRouter()
const store = useLogin()
const richTxtRef = ref()
const dialogShow = ref(false)
const title = ref('')
async function submit() {
  if (!store.loginMark) {
    dialogShow.value = !dialogShow.value
    return
  }
  const content = richTxtRef.value.tinymce!.activeEditor!.getContent()
  if(!content) {
    ElMessage.warning('无内容！')
    return
  }
  await ElMessageBox.confirm(
    '确认发布此条博客吗?',
    '确认',
    { showClose: false }
  )
  await ajaxArticle.add({ title: title.value, content: content })
  ElMessage.success('发布成功')
  router.push('/')
}
</script>
