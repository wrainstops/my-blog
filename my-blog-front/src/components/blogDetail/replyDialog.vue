@ -0,0 +1,83 @@
<template>
  <el-dialog
    v-model="dialogData.dialogShow"
    :title="dialogData.title"
    :append-to-body="true"
    width="500"
  >
    <div>
      <el-input
        v-model="comment"
        type="textarea"
        placeholder="发布你的评论"
        row="5"
        maxlength="500"
        show-word-limit
      />
      <div class="blank">
        <el-button @click="handleReplyMom">评论</el-button>
      </div>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { ajaxArticle } from '@/api/article'

const props = defineProps({
  rowData: {
    type: Object,
    default: () => ({})
  }
})
const emit = defineEmits(['confirm', 'cancel', 'update:isShow'])

const dialogData = reactive({
  dialogShow: true,
  title: `TO ${ props.rowData.authName }`
})
const comment = ref('')
async function handleReplyMom() {
  const params = {
    parentId: props.rowData.targetId,
    content: comment.value,
    toAuthId: props.rowData.authId
  }
  await ajaxArticle.addReply(params)
  ElMessage.success(`回复${ props.rowData.authName }成功√`)
  toggleDialog()
  confirm()
}

watch(
  () => dialogData.dialogShow,
  val => {
    emit('update:isShow', val)
  }
)
function toggleDialog () {
  dialogData.dialogShow = !dialogData.dialogShow
}
function confirm() {
  emit('confirm',{
    fn: () => {},
    comName: 'replyDialog'
  })
}
</script>

<style scoped lang="scss">
.blank {
  padding: 10px 0 30px;
  text-align: end;
}
</style>
