<template>
  <a-layout style="min-height: 100vh">
    <a-layout-sider collapsible>
      <div class="logo">CMDB</div>
      <a-menu theme="dark" mode="inline" :items="menuItems" />
    </a-layout-sider>
    <a-layout>
      <a-layout-header class="header">
        <span>资产管理控制台</span>
        <a-button size="small" @click="handleLogout">退出</a-button>
      </a-layout-header>
      <a-layout-content class="content">
        <a-card title="欢迎" :bordered="false">
          使用 Kratos + Ant Design Pro Vue 的骨架已就绪。
        </a-card>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import { clearAuthUser } from '../utils/auth'

const menuItems = reactive([
  { key: 'dashboard', label: '仪表盘' },
  { key: 'assets', label: '资产' },
  { key: 'settings', label: '设置' },
])

const router = useRouter()

const handleLogout = async () => {
  clearAuthUser()
  await router.replace({ name: 'login' })
}
</script>

<style scoped>
.logo {
  height: 32px;
  margin: 16px;
  color: #fff;
  font-weight: 600;
  text-align: center;
  line-height: 32px;
}

.header {
  background: #fff;
  padding: 0 24px;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.content {
  margin: 24px;
}
</style>
