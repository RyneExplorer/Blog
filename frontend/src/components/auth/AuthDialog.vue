<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'

import { useAuthStore } from '@/stores/authStore'

const router = useRouter()
const authStore = useAuthStore()
const {
  authVisible,
  authTab,
  authSubmitting,
  authMessage,
  authMessageType,
  loginForm,
  registerForm,
  captchaImage,
} = storeToRefs(authStore)

const isLogin = computed(() => authTab.value === 'login')
const isRegister = computed(() => authTab.value === 'register')

async function handleLogin() {
  const ok = await authStore.loginAction()
  if (ok) {
    authStore.consumeRedirect(router)
  }
}
</script>

<template>
  <transition name="auth-fade">
    <div v-if="authVisible" class="auth-overlay" @click.self="authStore.closeAuth">
      <div class="auth-card">
        <button type="button" class="auth-close" @click="authStore.closeAuth">×</button>

        <div class="auth-tabs">
          <button type="button" class="auth-tab" :class="{ 'auth-tab--active': isLogin }" @click="authStore.openAuth('login')">
            登录
          </button>
          <button type="button" class="auth-tab" :class="{ 'auth-tab--active': isRegister }" @click="authStore.openAuth('register')">
            注册
          </button>
        </div>

        <div v-if="isLogin" class="auth-content">
          <h2 class="auth-title">登录博客系统</h2>
          <p class="auth-subtitle">继续阅读、评论、收藏和创作文章</p>
          <p v-if="authMessage" class="auth-message" :class="`auth-message--${authMessageType}`">{{ authMessage }}</p>

          <form class="auth-form" @submit.prevent="handleLogin">
            <div class="input-group">
              <span class="input-icon">👤</span>
              <input v-model="loginForm.username" type="text" placeholder="用户名" :disabled="authSubmitting" />
            </div>
            <div class="input-group">
              <span class="input-icon">🔒</span>
              <input v-model="loginForm.password" type="password" placeholder="密码" :disabled="authSubmitting" />
            </div>
            <div class="input-group input-group--captcha">
              <span class="input-icon">🛡️</span>
              <input v-model="loginForm.captcha" type="text" placeholder="图形验证码" :disabled="authSubmitting" />
              <img :src="captchaImage" alt="验证码" class="captcha-img" @click="authStore.loadCaptcha" />
            </div>
            <button type="submit" class="auth-submit" :disabled="authSubmitting">
              {{ authSubmitting ? '登录中...' : '登录' }}
            </button>
          </form>
        </div>

        <div v-else-if="isRegister" class="auth-content">
          <h2 class="auth-title">注册账号</h2>
          <p class="auth-subtitle">用户名只允许字母数字，邮箱需为 qq.com</p>
          <p v-if="authMessage" class="auth-message" :class="`auth-message--${authMessageType}`">{{ authMessage }}</p>

          <form class="auth-form" @submit.prevent="authStore.registerAction">
            <div class="input-group">
              <span class="input-icon">👤</span>
              <input v-model="registerForm.username" type="text" placeholder="用户名" :disabled="authSubmitting" />
            </div>
            <div class="input-group input-group--inline">
              <span class="input-icon">✉️</span>
              <input v-model="registerForm.email" type="email" placeholder="QQ 邮箱" :disabled="authSubmitting" />
              <button type="button" class="mini-button" @click="authStore.requestEmailCode(registerForm.email)">发送验证码</button>
            </div>
            <div class="input-group">
              <span class="input-icon">🔐</span>
              <input v-model="registerForm.captcha" type="text" placeholder="邮箱验证码" :disabled="authSubmitting" />
            </div>
            <div class="input-group">
              <span class="input-icon">🔒</span>
              <input v-model="registerForm.password" type="password" placeholder="密码" :disabled="authSubmitting" />
            </div>
            <div class="input-group">
              <span class="input-icon">🔒</span>
              <input v-model="registerForm.confirm_password" type="password" placeholder="确认密码" :disabled="authSubmitting" />
            </div>
            <button type="submit" class="auth-submit" :disabled="authSubmitting">
              {{ authSubmitting ? '注册中...' : '注册' }}
            </button>
          </form>
        </div>

      </div>
    </div>
  </transition>
</template>

<style scoped>
.auth-overlay {
  position: fixed;
  inset: 0;
  z-index: 100;
  display: grid;
  place-items: center;
  background: rgba(43, 48, 47, 0.42);
  backdrop-filter: blur(8px);
}

.auth-card {
  position: relative;
  width: 100%;
  max-width: 520px;
  padding: 32px 32px 36px;
  border-radius: 24px;
  background: #f8f6f2;
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.12);
}

.auth-close {
  position: absolute;
  top: 16px;
  right: 16px;
  width: 32px;
  height: 32px;
  border: 0;
  border-radius: 50%;
  background: transparent;
  color: #6e746f;
  font-size: 24px;
  cursor: pointer;
  display: grid;
  place-items: center;
}

.auth-tabs {
  display: flex;
  gap: 10px;
  margin-bottom: 24px;
}

.auth-tab {
  padding: 10px 14px;
  border: 0;
  border-radius: 12px;
  background: #ece8e0;
  color: #6e746f;
  cursor: pointer;
}

.auth-tab--active {
  background: #5da39e;
  color: #fff;
}

.auth-content {
  text-align: center;
}

.auth-title {
  margin-bottom: 12px;
  color: #2b302f;
  font-size: 30px;
  font-weight: 700;
}

.auth-subtitle {
  margin-bottom: 20px;
  color: #8c928d;
  font-size: 16px;
}

.auth-message {
  margin-bottom: 18px;
  padding: 12px 14px;
  border-radius: 12px;
  font-size: 14px;
  line-height: 1.6;
  text-align: left;
}

.auth-message--error {
  background: rgba(201, 94, 71, 0.1);
  color: #b24a35;
}

.auth-message--success {
  background: rgba(93, 163, 158, 0.12);
  color: #2f7d77;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.input-group {
  position: relative;
  display: flex;
  align-items: center;
  min-height: 60px;
  padding: 0 20px;
  border: 1.5px solid #e8e4db;
  border-radius: 16px;
  background: #fff;
}

.input-group:focus-within {
  border-color: #5da39e;
}

.input-group--captcha,
.input-group--inline {
  padding-right: 10px;
}

.input-icon {
  margin-right: 14px;
  font-size: 20px;
  opacity: 0.6;
}

.input-group input {
  flex: 1;
  border: 0;
  background: transparent;
  color: #2b302f;
  font-size: 16px;
  outline: none;
}

.captcha-img {
  height: 40px;
  border-radius: 8px;
  cursor: pointer;
}

.mini-button {
  padding: 8px 12px;
  border: 0;
  border-radius: 10px;
  background: rgba(93, 163, 158, 0.12);
  color: #2f7d77;
  cursor: pointer;
}

.auth-submit {
  display: grid;
  place-items: center;
  height: 60px;
  margin-top: 8px;
  border: 0;
  border-radius: 16px;
  background: #5da39e;
  color: #fff;
  font-size: 18px;
  font-weight: 500;
  cursor: pointer;
}

.auth-submit:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.auth-footer {
  margin-top: 24px;
  color: #6e746f;
  font-size: 14px;
}

.auth-footer button {
  border: 0;
  background: transparent;
  color: #5da39e;
  font-weight: 600;
  cursor: pointer;
}

.auth-fade-enter-active,
.auth-fade-leave-active {
  transition: opacity 0.3s ease;
}

.auth-fade-enter-from,
.auth-fade-leave-to {
  opacity: 0;
}

@media (max-width: 640px) {
  .auth-card {
    width: calc(100vw - 24px);
    padding: 24px 18px 28px;
  }

  .auth-tabs {
    flex-wrap: wrap;
  }
}
</style>
