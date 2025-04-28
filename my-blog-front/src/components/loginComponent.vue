<template>
  <el-dialog
    v-model="dialogData.dialogShow"
    :append-to-body="true"
    width="500"
  >
    <el-tabs class="tabs" v-model="activeTab">
      <div class="whole-page">
          <el-tab-pane label="登录" name="login">
            <el-form :model="form1" class="login-form">
              <el-form-item label="昵称">
                <el-input v-model="form1.name" clearable />
              </el-form-item>

              <el-form-item label="密码">
                <el-input v-model="form1.password" type="password" show-password clearable />
              </el-form-item>

              <el-form-item>
                <el-button @click="handleLogin">登录</el-button>
              </el-form-item>
            </el-form>
          </el-tab-pane>
          <el-tab-pane label="注册" name="register">
            <el-form :model="form2" class="login-form">
              <el-form-item label="昵称">
                <el-input v-model="form2.name" clearable />
              </el-form-item>

              <el-form-item label="密码">
                <el-input v-model="form2.password" type="password" show-password clearable />
              </el-form-item>

              <el-form-item>
                <el-button @click="handleRegister">注册</el-button>
              </el-form-item>
            </el-form>
          </el-tab-pane>
      </div>
    </el-tabs>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useLogin } from '@/store/modules/loginMark'
import { setToken, setLocalUserInfo } from '@/utils/auth'
import { ajaxLog } from '@/api/auth'

const emit = defineEmits(['update:isShow'])
const store = useLogin()

const dialogData = reactive({
  dialogShow: true
})
const activeTab = ref('login')
const form1 = reactive({
  name: '',
  password: ''
})
const form2 = reactive({
  name: '',
  password: ''
})
async function handleLogin() {
  const ret: any = await ajaxLog.login(form1)
  setToken(ret)
  await getUserInfo()
  ElMessage.success('登录成功√')
  store.hasLogin()
  location.reload()
}

async function handleRegister() {
  await ajaxLog.register(form2)
  ElMessage.success('注册成功√')
  form1.name = form2.name
  form1.password = form2.password
  await handleLogin()
}

async function getUserInfo() {
  const userInfo: any = await ajaxLog.getUserInfo()
  setLocalUserInfo(userInfo)
}

watch(
  () => dialogData.dialogShow,
  val => {
    emit('update:isShow', val)
  }
)
</script>

<style scoped lang="scss">
.whole-page {
  width: 100%;
  display: grid;
  justify-content: center;
  align-items: center;
}
.tab {
  width: 100%;
}
.login-form {
  .el-input {
    width: 200px;
  }

  :last-child {
    justify-content: center;
  }
}
:deep(.el-tabs__active-bar) {
  height: 0;
}
:deep(.el-tabs__item.is-active) {
  color: var(--el-menu-active-color);
}
:deep(.el-tabs__item:hover) {
  color: var(--el-menu-hover-text-color);
}
</style>
