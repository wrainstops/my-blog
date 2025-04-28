<template>
  <div class="pl-30px pr-30px h-full overflow-auto" @scroll="handleScroll($event as UIEvent)">
    <div class="flex justify-between mt-30px">
      <div class="flex">
        <div class="flex items-center">
          <el-avatar :size="40" class="avatar">{{ getLastLetter(articleData?.authName) }}</el-avatar>
        </div>
        <div class="flex flex-col flex-1 p-10px gap-3px">
          <h3>{{ articleData?.authName || '' }}</h3>
          <p class="font-size-13px">{{ formatTime(articleData?.createdTime) }}</p>
        </div>
      </div>

      <div class="flex items-end pointer" @click="clickLike(articleData)">
        <el-icon v-if="articleData?.likeFlag" class="big_icon"><StarFilled /></el-icon>
        <el-icon v-else class="big_icon"><Star /></el-icon>
        <p class="mt-1px">{{ articleData?.likeNum }}</p>
      </div>
    </div>

    <el-divider />

    <h2 class="text-center">{{ articleData?.title }}</h2>

    <div v-html="articleData?.content" />

    <div class="h-50px" />

    <el-input
      v-model="comment"
      type="textarea"
      placeholder="发布你的评论"
      row="5"
      maxlength="200"
      show-word-limit
    />

    <div class="blank">
      <el-button @click="handleReplyMom">评论</el-button>
    </div>

    <template v-if="commentList.length > 0">
      <!-- reply list -->
      <div v-for="(item, index) in commentList" :key="index" class="mt-10px">
        <div>
          <!-- head -->
          <div class="flex justify-between">
            <div class="flex items-center">
              <el-avatar :size="32" class="avatar">{{ getLastLetter(item?.authName) }}</el-avatar>
              <div class="flex h-full pl-10px gap-10px">
                <p class="h-full flex items-center font-size-16px">{{ item?.authName || '' }}</p>
                <p class="h-full flex items-center font-size-13px">{{ formatTime(item?.createdTime) }}</p>
              </div>
            </div>
            <div class="flex items-end pointer">
              <el-icon class="mr-10px" @click="handleReplySon(item)"><EditPen /></el-icon>
              <el-icon v-if="item?.likeFlag" @click="clickLike(item)"><StarFilled /></el-icon>
              <el-icon v-else @click="clickLike(item)"><Star /></el-icon>
              <p>{{ item?.likeNum }}</p>
            </div>
          </div>

          <!-- content -->
          <div class="pl-42px" v-html="item.content" />

          <div class="text-end font-size-14px">
            <el-button v-if="!item._showReply && item.replyNum" link @click="showSonReplyList(item)">展开{{ item.replyNum }}条回复</el-button>
          </div>
          
          <template v-if="item._showReply">
            <!-- child reply list -->
            <div v-for="(eachReply, index) in item.reply" :key="index" class="mt-20px">
              <div class="pl-42px">
                <!-- child reply head -->
                <div class="flex justify-between">
                  <div class="flex items-center">
                    <el-avatar :size="28" class="avatar">{{ getLastLetter(eachReply?.authName) }}</el-avatar>
                    <div class="flex h-full pl-10px gap-10px">
                      <p class="h-full flex items-center font-size-16px">{{ eachReply?.authName || '' }}<span v-if="eachReply.toAuthName" class="pl-5px pr-5px">TO</span>{{ eachReply.toAuthName || '' }}</p>
                      <p class="h-full flex items-center font-size-13px">{{ formatTime(eachReply?.createdTime) }}</p>
                    </div>
                  </div>
                  <div class="flex items-end pointer">
                    <el-icon class="mr-10px" @click="handleReplySonSon(eachReply)"><EditPen /></el-icon>
                    <el-icon v-if="eachReply?.likeFlag" @click="clickLike(eachReply)"><StarFilled /></el-icon>
                    <el-icon v-else @click="clickLike(eachReply)"><Star /></el-icon>
                    <p>{{ eachReply?.likeNum }}</p>
                  </div>
                </div>

                <!-- child reply content -->
                <div class="pl-38px">{{ eachReply.content }}</div>
              </div>
            </div>
            <div class="flex justify-end items-center font-size-13px">
              <el-button link class="mr-10px" @click="hideSonReplyList(item)">收起回复</el-button>
              <el-button v-if="item.reply && item.reply.length < item.replyNum" link @click="getMoreSonReply(item)">加载更多</el-button>
              <span v-else>没有更多了...</span>
            </div>
          </template>
        </div>
      </div>
      <div class="text-align-center">
        <span v-if="isEnd">到底了...</span>
        <span v-else>下滑加载更多</span>
      </div>
    </template>
    <login-component v-if="dialogShow" v-model:isShow="dialogShow" />
    <template v-if="dialogProps.dialogShow">
      <component
        :is="dialogCom[dialogProps.dialogName]"
        v-model:isShow="dialogProps.dialogShow"
        :row-data="dialogProps.rowData"
        @confirm="dialogConfirm"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, unref, reactive, inject, shallowRef, defineComponent, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Star, StarFilled, EditPen } from '@element-plus/icons-vue'
import { useLogin } from '@/store/modules/loginMark'
import { useDivLoadMore } from '@/composables/loadMore'
import { formatTime, hasWorth, getLastLetter } from '@/utils/common'
import replyDialog from './replyDialog.vue'
import { ajaxArticle } from '@/api/article'

const id = inject('id')
const getNewList = inject('getNewList') as Function
const numId = ref()

const router = useRouter()
const store = useLogin()

const dialogShow = ref(false)

const articleData = ref() // 博客内容
const commentList = ref([] as any[]) // 评论列表
const comment = ref('') // 评论输入框
const replyListParams = reactive({
  page: 0,
  pageSize: 10
})

const isEnd = ref(false) // 是否到底标记
// 获取博客详情和评论列表
async function getArticleDetail() {
  if (numId.value === 0 || numId.value + '' === 'NaN') {
    router.push('/')
    return
  }
  try {
    const res: any = await ajaxArticle.getById(numId.value)
    articleData.value = res
  } catch {
    router.push('/')
  }
}
async function getReplyList() {
  if (numId.value === 0 || numId.value + '' === 'NaN') {
    router.push('/')
    return
  }
  replyListParams.page = 0
  const res: any = await ajaxArticle.queryReply({ ...replyListParams, parentId: numId.value })
  commentList.value = res.content
  replyListParams.page++
  isEnd.value = false
}
async function getMoreReplyList() {
  if (!isEnd.value) {
    const res: any = await ajaxArticle.queryReply({ ...replyListParams, parentId: numId.value })
    if (!(res.content && res.content.length)) {
      isEnd.value = true
      return
    }
    commentList.value.push(...res.content)
    replyListParams.page++
  }
}
function handleScroll(e: UIEvent) {
  useDivLoadMore(e, () => isEnd.value, getMoreReplyList)
}

watch(
  () => unref(id),
  (val) => {
    numId.value = Number(val)
    getArticleDetail()
    getReplyList()
  },
  { immediate: true }
)

async function handleReplyMom() {
  if (!store.loginMark) {
    dialogShow.value = !dialogShow.value
    return
  }
  if (!comment.value) {
    ElMessage.warning('无内容！')
    return
  }
  const params = {
    parentId: articleData.value.ID,
    content: comment.value,
    toAuthId: articleData.value.authId
  }
  await ElMessageBox.confirm(
    '确认评论此条博客吗?',
    '确认',
    { showClose: false }
  )
  await ajaxArticle.addReply(params)
  ElMessage.success('评论成功!')
  getReplyList()
  getNewList() // 更新左侧列表
  comment.value = ''
}
function handleReplySon(row: any) {
  if (!store.loginMark) {
    dialogShow.value = !dialogShow.value
    return
  }
  dialogToggle('replyDialog', { ...row, targetId: row.ID })
}
function handleReplySonSon(row: any) {
  if (!store.loginMark) {
    dialogShow.value = !dialogShow.value
    return
  }
  dialogToggle('replyDialog', { ...row, targetId: row.parentId })
}
function clickLike(e: any) {
  if (!store.loginMark) {
    dialogShow.value = !dialogShow.value
    return
  }
  const opeType = e?.likeFlag ? 'cancelLike' : 'addLike'
  ajaxArticle[opeType]({ articleId: e.ID})
    .then(() => {
      e.likeFlag = !e?.likeFlag
      if (e.likeFlag) {
        e.likeNum++
      } else {
        e.likeNum--
      }
      getNewList() // 更新左侧列表
    })
}

const dialogCom: ReturnType<typeof defineComponent>  = reactive({
  replyDialog: shallowRef(replyDialog),
})
const dialogProps = reactive({
  dialogName: '',
  dialogShow: false,
  rowData: {}
})

function dialogToggle(name: string, rowData = {}) {
  dialogProps.dialogName = name
  dialogProps.rowData = { ...rowData }
  dialogProps.dialogShow = !dialogProps.dialogShow
}
function dialogConfirm() {
  getReplyList()
  getNewList() // 更新左侧列表
}

const moreReplyParams = {
  page: 0,
  pageSize: 10
}
// 第一次点开子回复
async function showSonReplyList(row: any) {
  if (row.replyNum === 0) return
  row._showReply = true
  if (row.replyNum && !(row.reply && row.reply.length)) {
    const rep: any = await ajaxArticle.queryReply({ ...moreReplyParams, parentId: Number(row.ID) })
    row.reply = rep.content
    row.page = hasWorth(row.page) ? row.page + 1 : 0
  }
}
// 收起子回复
function hideSonReplyList(row: any) {
  row._showReply = false
}
// 加载更多子回复
async function getMoreSonReply(item: any) {
  const rep: any = await ajaxArticle.queryReply({ page: item.page, pageSize: 20, parentId: Number(item.ID) })
  item.reply.push(rep.content)
  item.page++
}
</script>

<style scoped lang="scss">
h3, p {
  margin: 0;
}
.el-icon {
  height: 1.2em; // 为了与p标签对齐
}
.big_icon {
  height: 2em;
  width: 2em;

  svg {
    height: 100%;
    width: 100%;
  }
}
.el-divider--horizontal {
  border-top: 1px solid #646464;
}
.blank {
  padding: 10px 0 30px;
  text-align: end;
}
.text-end {
  text-align: end;
}
</style>
