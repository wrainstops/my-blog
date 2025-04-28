import { $get, $post } from '@/utils/request'

export const ajaxLog = {
  register: (data: any) => $post('/auth/register', data),
  login: (data: any) => $post('/auth/login', data),
  logout: () => $post('/auth/logout'),
  getUserInfo: () => $get('/auth/info'),
  getStats: () => $get('/auth/getStats')
}

export const logout = () => {
  ajaxLog.logout().finally(() => {
    localStorage.clear()
    location.reload()
  })
}
