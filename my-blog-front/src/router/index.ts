import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useLogin } from '@/store/modules/loginMark'
import { loading } from '@/utils/loading'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/',
    component: () => import('@/layout/index.vue'),
    children: [
      {
        path: '/dashboard',
        name: 'dashboard',
        component: () => import('@/views/dashboard/index.vue')
      },
      {
        path: '/blogging',
        name: 'blogging',
        component: () => import('@/views/blogging/index.vue')
      },
      {
        path: '/details/:id',
        name: 'details',
        component: () => import('@/views/details/index.vue'),
        props: true
      },
      {
        path: '/personalCenter',
        name: 'personalCenter',
        component: () => import('@/views/personalCenter/index.vue')
      }
    ]
  },
  {
    path: '/404',
    name: '404',
    component: () => import('@/views/error/404.vue')
  },
  {
    path: '/:pathMatch(.*)*',
    component: () => import('@/views/error/404.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

const routerWhiteList = ['dashboard', 'blogging', '404']

router.beforeEach((to, _from, next) => {
  if (!window.existLoading) {
    loading.show()
    window.existLoading = true
  }

  /**
   * !!! 必须在内部调用useLogin()
   * 否则报错：
   * "getActivePinia()" was called but there was no active Pinia.
   * Are you trying to use a store before calling "app.use(pinia)"?
   */
  const store = useLogin()

  if (!routerWhiteList.includes(to.name as string)) {
    if (!store.loginMark) { // 未登录
      next('/404')
    } else { // 已登录
      next()
    }
  } else {
    next()
  }
})
router.afterEach(() => {
  if (window.existLoading) {
    loading.hide()
    window.existLoading = false
  }
})

export default router
