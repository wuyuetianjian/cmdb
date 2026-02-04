import { createRouter, createWebHistory } from 'vue-router'
import BasicLayout from '../layouts/BasicLayout.vue'
import LoginView from '../views/LoginView.vue'
import { isAuthenticated } from '../utils/auth'

const routes = [
  {
    path: '/login',
    name: 'login',
    component: LoginView,
  },
  {
    path: '/',
    name: 'home',
    component: BasicLayout,
    meta: { requiresAuth: true },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  if (to.meta.requiresAuth && !isAuthenticated()) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  return true
})

export default router
