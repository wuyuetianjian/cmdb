<template>
  <div class="login">
    <a-card title="登录" class="login-card">
      <p>需要认证后才能访问控制台。</p>
      <a-space direction="vertical" style="width: 100%">
        <a-button type="primary" block @click="handleSSOLogin">使用 SSO 登录</a-button>
        <a-button block @click="handleDevLogin">本地模拟登录</a-button>
      </a-space>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { setAuthUser } from '../utils/auth'

const router = useRouter()
const route = useRoute()

const handleSSOLogin = () => {
  window.location.href = '/auth/sso/login'
}

const handleDevLogin = async () => {
  setAuthUser('local-user')
  await router.replace((route.query.redirect as string) || '/')
}
</script>

<style scoped>
.login {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f2f5;
}

.login-card {
  width: 360px;
}
</style>
