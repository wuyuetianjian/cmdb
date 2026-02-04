<template>
  <div class="login">
    <a-card title="登录" class="login-card">
      <p>需要认证后才能访问控制台。</p>
      <a-space direction="vertical" style="width: 100%">
        <a-form layout="vertical" @submit.prevent="handleLocalLogin">
          <a-form-item label="用户名">
            <a-input v-model:value="formState.username" placeholder="请输入用户名" />
          </a-form-item>
          <a-form-item label="密码">
            <a-input-password v-model:value="formState.password" placeholder="请输入密码" />
          </a-form-item>
          <a-button type="primary" block :loading="submitting" @click="handleLocalLogin">
            本地登录
          </a-button>
        </a-form>
        <a-button v-if="ssoEnabled" block @click="handleSSOLogin">使用 SSO 登录</a-button>
      </a-space>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { setAuthUser } from '../utils/auth'

const router = useRouter()
const route = useRoute()
const submitting = ref(false)
const ssoEnabled = import.meta.env.VITE_SSO_ENABLED === 'true'

const formState = reactive({
  username: '',
  password: '',
})

const handleSSOLogin = () => {
  window.location.href = '/auth/sso/login'
}

const handleLocalLogin = async () => {
  if (!formState.username || !formState.password) {
    message.error('请输入用户名和密码')
    return
  }

  submitting.value = true
  try {
    const response = await fetch('/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(formState),
    })
    if (!response.ok) {
      throw new Error('login failed')
    }
    setAuthUser(formState.username)
    await router.replace((route.query.redirect as string) || '/')
  } catch (error) {
    message.error('登录失败，请检查账号信息')
  } finally {
    submitting.value = false
  }
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
