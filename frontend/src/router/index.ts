import { createRouter, createWebHistory } from 'vue-router'
import BasicLayout from '../layouts/BasicLayout.vue'

const routes = [
  {
    path: '/',
    name: 'home',
    component: BasicLayout,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
