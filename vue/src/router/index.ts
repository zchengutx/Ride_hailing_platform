import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: {
      title: '首页'
    }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: {
      title: '登录'
    }
  },
  {
    path: '/login-callback',
    name: 'LoginCallback',
    component: () => import('@/views/LoginCallback.vue'),
    meta: {
      title: '登录处理中'
    }
  },
  {
    path: '/address-input',
    name: 'AddressInput',
    component: () => import('@/views/AddressInput.vue'),
    meta: {
      title: '选择地点'
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title} - 打车平台` : '打车平台'
  next()
})

export default router 