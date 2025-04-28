<template>
  <div class="l_container">
    <div class="header">
      <div class="header_Nav" :class="darkMode ? 'header_Nav_Dark' : 'header_Nav_Light'">
        <el-scrollbar>
          <div class="header_Nav_container">
            <div>
              <el-image class="h-50px w-auto" fit="contain" :src="darkMode ? LogoDark : Logo"></el-image>
            </div>
            <div>
              <el-menu
                :default-active="activeIndex"
                class="El_menu"
                mode="horizontal"
                @select="handleSelectMenu"
              >
                <el-menu-item index="/">首页</el-menu-item>
                <el-menu-item index="/blogging">发博</el-menu-item>
              </el-menu>
            </div>

            <div class="person">
              <el-switch
                ref="switchRef"
                v-model="darkMode"
                :before-change="beforeChange"
                :active-action-icon="Moon"
                :inactive-action-icon="Sunny"
              />

              <el-divider direction="vertical" style="border-left: 2px dashed var(--main-color-opposite);" />

              <el-avatar v-if="!store.loginMark" :size="42" class="pointer" @click="getLoginComponent">登录</el-avatar>
              <el-dropdown v-else class="pointer">
                <el-avatar :size="42">{{ getLastLetter(name) }}</el-avatar>
                
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="goPersonalCenter">个人中心</el-dropdown-item>
                    <el-dropdown-item @click="logout">退出</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
              <div>{{ name }}</div>
            </div>
          </div>
        </el-scrollbar>
      </div>
      <div class="header_Nav_holder" />
    </div>
    <div class="main">
      <router-view />
    </div>

    <plug v-if="showPlug" />
    <login-component v-if="dialogShow" v-model:isShow="dialogShow" />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import type { SwitchInstance } from 'element-plus'
import { Sunny, Moon } from '@element-plus/icons-vue'
import { isDark, toggleDark } from '@/composables/dark'
import { useLogin } from '@/store/modules/loginMark'
import { getLocalUserInfo, removeToken, removeLocalUserInfo } from '@/utils/auth'
import { getLastLetter } from '@/utils/common'
import plug from '@/components/plug.vue'
import loginComponent from '@/components/loginComponent.vue'
import Logo from '@/assets/image/logo.png'
import LogoDark from '@/assets/image/logo-dark.png'

const myRouter = useRouter()
const activeIndex = ref('/')
const handleSelectMenu = (key: string) => {
  myRouter.push(key)
}

const darkMode = ref(isDark.value)
const switchRef = ref<SwitchInstance>()
const name = ref(getLocalUserInfo()?.name || '')
const store = useLogin()
const showPlug = ref(true)

watch(
  () => darkMode.value,
  () => {
    toggleDark()
  }
)
const beforeChange = () => {
  return new Promise<boolean>((resolve) => {
    const isAppearanceTransition =
            // @ts-expect-error
            document.startViewTransition &&
            !window.matchMedia('(prefers-reduced-motion: reduce)').matches
    if (!isAppearanceTransition) {
      resolve(true)
      return
    }

    const switchElement = switchRef.value?.$el
    const rect = switchElement.getBoundingClientRect()
    const x = rect.left + rect.width / 2
    const y = rect.top + rect.height / 2

    const endRadius = Math.hypot(
      Math.max(x, innerWidth - x),
      Math.max(y, innerHeight - y)
    )
    const transition = document.startViewTransition(async () => {
      resolve(true)
      await nextTick()
    })
    transition.ready.then(() => {
      const clipPath = [
        `circle(0px at ${ x }px ${ y }px)`,
        `circle(${ endRadius }px at ${ x }px ${ y }px)`,
      ]
      document.documentElement.animate(
        {
          clipPath: darkMode.value ? [...clipPath].reverse() : clipPath,
        },
        {
          duration: 400,
          easing: 'ease-in',
          pseudoElement: darkMode.value
            ? '::view-transition-old(root)'
            : '::view-transition-new(root)',
        }
      )
    })
  })
}
const dialogShow = ref(false)
function getLoginComponent() {
  dialogShow.value = !dialogShow.value
}

function goPersonalCenter() {
  myRouter.push('/personalCenter')
}

function logout() {
  removeToken()
  removeLocalUserInfo()
  name.value = ''
  myRouter.push('/')
  store.hasLogout()
}

function handleResizeInnerWidth() {
  window.innerWidth > 1280 ? (showPlug.value = true) : (showPlug.value = false)
}
onMounted(() => {
  window.addEventListener('resize', handleResizeInnerWidth)
})
onUnmounted(() => {
  window.removeEventListener('resize', handleResizeInnerWidth)
})
</script>

<style scoped lang="scss">
.l_container {
  min-height: 100%;
  height: 100%;
}
.header {

  .header_Nav_Dark {
    background: url('@/assets/image/dash_t_dark.jpg') no-repeat;
  }

  .header_Nav_Light {
    background: url('@/assets/image/dash_t.jpg') no-repeat;
  }

  .header_Nav {
    height: 60px;
    width: 100%;
    background-size: cover;
    position: fixed;
    left: 0;
    top: 0;
    z-index: 999;
    border-bottom: 1px solid var(--main-color);

    .header_Nav_container {
      height: 100%;
      width: 90%;
      min-width: 600px;
      margin: auto;
      display: flex;
      justify-content: space-between;
      align-items: center;

      .person {
        display: flex;
        align-items: center;
        gap: 10px;
      }
    }
  }

  .header_Nav_holder {
    height: 60px;
  }
}
.El_menu {
  min-width: 400px;
  justify-content: center;
  border-bottom: none;
}
:deep(.el-menu-item.is-active) {
  border-bottom: none;
}
.main {
  width: 100%;
}
</style>
