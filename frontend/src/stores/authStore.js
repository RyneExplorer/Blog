import { computed, reactive, ref } from 'vue'
import { defineStore } from 'pinia'

import { getCaptcha, login, logout as requestLogout, register, sendEmailCode } from '@/api/auth'
import { getProfile } from '@/api/user'
import { clearToken, getToken, resolveAssetUrl, setAuthFailureHandler, setToken, subscribeTokenChange } from '@/api/http'
import { useUiStore } from '@/stores/uiStore'

function parseJwtRole(token) {
  if (!token) {
    return 1
  }

  try {
    const encoded = token.split('.')[1]
    const normalized = encoded.replace(/-/g, '+').replace(/_/g, '/')
    const payload = JSON.parse(window.atob(normalized))
    return Number(payload.role ?? 1)
  } catch (error) {
    return 1
  }
}

function validateUsername(username) {
  return /^[A-Za-z0-9]+$/.test(username)
}

function validateQQEmail(email) {
  return /^[^@\s]+@qq\.com$/i.test(email)
}

export const useAuthStore = defineStore('auth', () => {
  const uiStore = useUiStore()

  const token = ref(getToken())
  const currentUser = ref(null)
  const initialized = ref(false)
  let initializeRequest = null

  const authVisible = ref(false)
  const authTab = ref('login')
  const authSubmitting = ref(false)
  const authMessage = ref('')
  const authMessageType = ref('info')
  const redirectAfterLogin = ref('')

  const captchaImage = ref('')

  const loginForm = reactive({
    username: '',
    password: '',
    captcha_id: '',
    captcha: '',
  })
  const registerForm = reactive({
    username: '',
    email: '',
    captcha: '',
    password: '',
    confirm_password: '',
  })

  const isLoggedIn = computed(() => Boolean(token.value))
  const isAdmin = computed(() => Number(currentUser.value?.role ?? parseJwtRole(token.value)) === 0)
  const avatarText = computed(() => {
    const source = currentUser.value?.nickname || currentUser.value?.username || '访'
    return source.slice(0, 1).toUpperCase()
  })

  function applyUser(user) {
    if (!user) {
      currentUser.value = null
      return
    }

    const role = parseJwtRole(token.value)
    currentUser.value = {
      ...user,
      avatar: resolveAssetUrl(user?.avatar || ''),
      role,
      isAdmin: role === 0,
    }
  }

  function clearAuthMessage() {
    authMessage.value = ''
    authMessageType.value = 'info'
  }

  function resetLocalAuthState() {
    clearToken()
    currentUser.value = null
    redirectAfterLogin.value = ''
  }

  function rememberCurrentRoute() {
    const target = `${window.location.pathname}${window.location.search}${window.location.hash}`
    if (target && target !== '/') {
      redirectAfterLogin.value = target
    }
  }

  function openAuth(tab = 'login', redirect = '') {
    authTab.value = tab
    if (redirect) {
      redirectAfterLogin.value = redirect
    }
    clearAuthMessage()
    authVisible.value = true
    if (tab === 'login') {
      loadCaptcha()
    }
  }

  function closeAuth() {
    authVisible.value = false
  }

  async function loadCaptcha() {
    try {
      const data = await getCaptcha()
      captchaImage.value = data?.captcha || ''
      loginForm.captcha_id = data?.captcha_id || ''
      loginForm.captcha = ''
    } catch (error) {
      uiStore.showToast(error.message, 'error')
    }
  }

  async function fetchCurrentUser() {
    if (!token.value) {
      currentUser.value = null
      return null
    }

    const data = await getProfile()
    applyUser(data)
    return currentUser.value
  }

  async function initializeAuth() {
    if (initialized.value) {
      return
    }

    if (!initializeRequest) {
      // 1. 路由守卫和布局挂载可能同时初始化登录态。
      // 2. 复用同一个 Promise，避免刷新页面时重复请求 /api/user/profile。
      // 3. 初始化结束后清空引用，后续登录态变化仍可重新触发初始化流程。
      initializeRequest = (async () => {
        if (token.value) {
          try {
            await fetchCurrentUser()
          } catch (error) {
            resetLocalAuthState()
          }
        }

        initialized.value = true
      })().finally(() => {
        initializeRequest = null
      })
    }

    return initializeRequest
  }

  async function ensureInitialized() {
    if (!initialized.value) {
      await initializeAuth()
    }
  }

  async function loginAction() {
    if (!validateUsername(loginForm.username)) {
      authMessage.value = '用户名只能包含字母和数字'
      authMessageType.value = 'error'
      return false
    }

    authSubmitting.value = true

    try {
      const data = await login(loginForm)
      const nextToken = data?.token || ''
      token.value = nextToken
      setToken(nextToken)
      applyUser(data?.user || null)
      clearAuthMessage()
      authVisible.value = false
      uiStore.showToast('登录成功', 'success')
      await loadCaptcha()
      return true
    } catch (error) {
      authMessage.value = error.message
      authMessageType.value = 'error'
      await loadCaptcha()
      return false
    } finally {
      authSubmitting.value = false
    }
  }

  async function registerAction() {
    if (!validateUsername(registerForm.username)) {
      authMessage.value = '用户名只能包含字母和数字'
      authMessageType.value = 'error'
      return false
    }

    if (!validateQQEmail(registerForm.email)) {
      authMessage.value = '注册邮箱必须是 qq.com'
      authMessageType.value = 'error'
      return false
    }

    if (registerForm.password.length < 6) {
      authMessage.value = '密码至少 6 位'
      authMessageType.value = 'error'
      return false
    }

    if (registerForm.password !== registerForm.confirm_password) {
      authMessage.value = '两次输入的密码不一致'
      authMessageType.value = 'error'
      return false
    }

    authSubmitting.value = true

    try {
      await register(registerForm)
      authTab.value = 'login'
      authMessage.value = '注册成功，请登录'
      authMessageType.value = 'success'
      registerForm.username = ''
      registerForm.email = ''
      registerForm.captcha = ''
      registerForm.password = ''
      registerForm.confirm_password = ''
      await loadCaptcha()
      return true
    } catch (error) {
      authMessage.value = error.message
      authMessageType.value = 'error'
      return false
    } finally {
      authSubmitting.value = false
    }
  }

  async function requestEmailCode(email) {
    if (!validateQQEmail(email)) {
      authMessage.value = '请填写 qq.com 邮箱'
      authMessageType.value = 'error'
      return false
    }

    try {
      await sendEmailCode(email)
      authMessage.value = '验证码已发送'
      authMessageType.value = 'success'
      return true
    } catch (error) {
      authMessage.value = error.message
      authMessageType.value = 'error'
      return false
    }
  }

  function consumeRedirect(router) {
    const target = redirectAfterLogin.value
    redirectAfterLogin.value = ''
    if (target) {
      router.push(target)
    } else {
      router.push({ name: 'articles' })
    }
  }

  async function handleSessionExpired(error) {
    resetLocalAuthState()
    rememberCurrentRoute()
    authTab.value = 'login'
    authMessage.value = error?.message || '登录状态已失效，请重新登录'
    authMessageType.value = 'error'
    authVisible.value = true
    await loadCaptcha()
    uiStore.showToast(authMessage.value, 'error')
  }

  async function logout() {
    const hadToken = Boolean(token.value)
    let logoutError = null

    try {
      if (hadToken) {
        await requestLogout()
      }
    } catch (error) {
      logoutError = error
    }

    resetLocalAuthState()
    authVisible.value = false
    clearAuthMessage()

    if (!hadToken) {
      uiStore.showToast('当前未登录', 'info')
      return true
    }

    if (!logoutError) {
      uiStore.showToast('已退出登录', 'success')
      return true
    }

    if (logoutError.isAuthError) {
      uiStore.showToast('登录状态已失效，已清理本地登录信息', 'info')
      return true
    }

    uiStore.showToast(`${logoutError.message}，已清理本地登录信息`, 'error')
    return false
  }

  subscribeTokenChange((nextToken) => {
    token.value = nextToken || ''

    if (!nextToken) {
      currentUser.value = null
    }
  })

  setAuthFailureHandler(handleSessionExpired)

  return {
    token,
    currentUser,
    initialized,
    authVisible,
    authTab,
    authSubmitting,
    authMessage,
    authMessageType,
    redirectAfterLogin,
    captchaImage,
    loginForm,
    registerForm,
    isLoggedIn,
    isAdmin,
    avatarText,
    applyUser,
    clearAuthMessage,
    openAuth,
    closeAuth,
    loadCaptcha,
    fetchCurrentUser,
    initializeAuth,
    ensureInitialized,
    loginAction,
    registerAction,
    requestEmailCode,
    consumeRedirect,
    handleSessionExpired,
    logout,
  }
})
