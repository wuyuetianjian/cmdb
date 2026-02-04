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
        <a-card title="SSO 设置" class="card" :bordered="false">
          <a-space>
            <span>启用 SSO</span>
            <a-switch :checked="ssoEnabled" :loading="ssoLoading" @change="handleSSOToggle" />
          </a-space>
        </a-card>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { clearAuthUser } from '../utils/auth'

const menuItems = reactive([
  { key: 'dashboard', label: '仪表盘' },
  { key: 'assets', label: '资产' },
  { key: 'settings', label: '设置' },
])

const router = useRouter()
const ssoEnabled = ref(false)
const ssoLoading = ref(false)

const handleLogout = async () => {
  clearAuthUser()
  await router.replace({ name: 'login' })
}

const loadSSOConfig = async () => {
  ssoLoading.value = true
  try {
    const response = await fetch('/sso/config')
    if (!response.ok) {
      throw new Error('load failed')
    }
    const data = (await response.json()) as { enabled: boolean }
    ssoEnabled.value = Boolean(data.enabled)
  } catch (error) {
    message.error('读取 SSO 配置失败')
  } finally {
    ssoLoading.value = false
  }
}

const handleSSOToggle = async (checked: boolean) => {
  ssoLoading.value = true
  try {
    const response = await fetch('/sso/config', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ enabled: checked, protocol: 'saml2' }),
    })
    if (!response.ok) {
      throw new Error('update failed')
    }
    ssoEnabled.value = checked
    message.success('SSO 配置已更新')
  } catch (error) {
    message.error('更新 SSO 配置失败')
  } finally {
    ssoLoading.value = false
  }
}

onMounted(() => {
  loadSSOConfig()
})
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

.card {
  margin-top: 16px;
}
</style>
