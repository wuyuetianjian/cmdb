<template>
  <a-layout style="min-height: 100vh">
    <a-layout-sider collapsible>
      <div class="logo">CMDB</div>
      <a-menu theme="dark" mode="inline" :items="menuItems" />
    </a-layout-sider>
    <a-layout>
      <a-layout-header class="header">资产管理控制台</a-layout-header>
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

const menuItems = reactive([
  { key: 'dashboard', label: '仪表盘' },
  { key: 'assets', label: '资产' },
  { key: 'settings', label: '设置' },
])
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
}

.content {
  margin: 24px;
}
</style>
