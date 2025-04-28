import { defineStore } from 'pinia'
import { getToken } from '@/utils/auth'

export const useLogin = defineStore('login', {
  persist: true, // 开启持久化
  state: () => ({ loginMark: Boolean(getToken()) }),
  actions: {
    hasLogin() {
      this.loginMark = true
    },
    hasLogout() {
      this.loginMark = false
    }
  }
})