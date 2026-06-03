import axios from 'axios'

const TOKEN_KEY = 'blog_token'
const AUTH_ERROR_CODES = new Set([401, 1005, 1006])
const DEFAULT_AUTH_ERROR_MESSAGE = '登录状态已失效，请重新登录'
const NETWORK_ERROR_MESSAGE = '网络异常，请稍后重试'
const SERVER_UNAVAILABLE_MESSAGE = '后端服务暂时不可用，请确认接口服务已经启动'
const BACKEND_UNREACHABLE_MESSAGE = '无法连接到后端服务，请确认接口服务已经启动'

const tokenListeners = new Set()

let authFailureHandler = null
let refreshRequest = null
let authFailureRequest = null

export function getToken() {
  return window.localStorage.getItem(TOKEN_KEY) || ''
}

function notifyTokenChange(token) {
  tokenListeners.forEach((listener) => {
    listener(token)
  })
}

export function setToken(token) {
  window.localStorage.setItem(TOKEN_KEY, token)
  notifyTokenChange(token)
}

export function clearToken() {
  window.localStorage.removeItem(TOKEN_KEY)
  notifyTokenChange('')
}

export function subscribeTokenChange(listener) {
  tokenListeners.add(listener)

  return () => {
    tokenListeners.delete(listener)
  }
}

export function setAuthFailureHandler(handler) {
  authFailureHandler = handler
}

function normalizeBaseURL() {
  return import.meta.env.VITE_API_BASE_URL || ''
}

function getPayloadMessage(payload) {
  return payload?.msg || payload?.message || ''
}

function isAuthErrorCode(code) {
  return AUTH_ERROR_CODES.has(Number(code))
}

function createApiError(message, options = {}) {
  const error = new Error(message || NETWORK_ERROR_MESSAGE)
  error.code = options.code
  error.status = options.status
  error.isAuthError = Boolean(options.isAuthError)
  error.payload = options.payload
  return error
}

function normalizeError(error) {
  if (error instanceof Error && Object.prototype.hasOwnProperty.call(error, 'status')) {
    return error
  }

  const status = error?.response?.status
  const payload = error?.response?.data
  const code = payload?.code
  const payloadMessage = getPayloadMessage(payload)

  let message = payloadMessage || error?.message || NETWORK_ERROR_MESSAGE

  if (status === 502 || status === 503 || status === 504) {
    message = SERVER_UNAVAILABLE_MESSAGE
  } else if (!error?.response) {
    message = BACKEND_UNREACHABLE_MESSAGE
  } else if (status === 401 || isAuthErrorCode(code)) {
    message = payloadMessage || DEFAULT_AUTH_ERROR_MESSAGE
  }

  return createApiError(message, {
    code,
    status,
    isAuthError: status === 401 || isAuthErrorCode(code),
    payload,
  })
}

export function resolveAssetUrl(url) {
  if (!url) {
    return ''
  }

  if (/^https?:\/\//.test(url)) {
    return url
  }

  if (url.startsWith('/')) {
    const base = normalizeBaseURL()
    if (base) {
      return `${base.replace(/\/$/, '')}${url}`
    }

    return url
  }

  return url
}

export function normalizePage(payload, fallbackPageSize = 10) {
  return payload?.data || {
    list: [],
    total: 0,
    page: 1,
    pageSize: fallbackPageSize,
    pages: 0,
  }
}

const http = axios.create({
  baseURL: normalizeBaseURL(),
  timeout: 12000,
  headers: {
    'X-Requested-With': 'XMLHttpRequest',
  },
})

const authClient = axios.create({
  baseURL: normalizeBaseURL(),
  timeout: 12000,
  headers: {
    'X-Requested-With': 'XMLHttpRequest',
  },
})

http.interceptors.request.use((config) => {
  const token = getToken()

  if (token) {
    config.headers = config.headers || {}
    config.headers.Authorization = `Bearer ${token}`
  }

  return config
})

async function refreshAccessToken() {
  const currentToken = getToken()

  if (!currentToken) {
    throw createApiError(DEFAULT_AUTH_ERROR_MESSAGE, {
      code: 401,
      status: 401,
      isAuthError: true,
    })
  }

  if (!refreshRequest) {
    refreshRequest = authClient
      .post('/api/auth/refresh', { token: currentToken })
      .then((response) => {
        const payload = response.data

        if (typeof payload?.code === 'number' && payload.code !== 0) {
          throw createApiError(getPayloadMessage(payload) || DEFAULT_AUTH_ERROR_MESSAGE, {
            code: payload.code,
            status: response.status,
            isAuthError: isAuthErrorCode(payload.code),
            payload,
          })
        }

        const nextToken = payload?.data?.token || ''

        if (!nextToken) {
          throw createApiError(DEFAULT_AUTH_ERROR_MESSAGE, {
            code: 401,
            status: response.status,
            isAuthError: true,
            payload,
          })
        }

        setToken(nextToken)
        return nextToken
      })
      .catch((error) => {
        throw normalizeError(error)
      })
      .finally(() => {
        refreshRequest = null
      })
  }

  return refreshRequest
}

function shouldRefreshRequest(config = {}, error) {
  if (config.skipAuthRefresh || config.__isRetryRequest || !error?.isAuthError) {
    return false
  }

  if (!getToken()) {
    return false
  }

  return `${config.url || ''}`.includes('/api/auth/refresh') === false
}

async function handleAuthFailure(error, config = {}) {
  if (config.skipAuthFailureHandler) {
    return
  }

  if (!authFailureRequest) {
    authFailureRequest = Promise.resolve()
      .then(async () => {
        if (getToken()) {
          clearToken()
        }

        if (typeof authFailureHandler === 'function') {
          await authFailureHandler(error)
        }
      })
      .finally(() => {
        authFailureRequest = null
      })
  }

  return authFailureRequest
}

async function retryWithFreshToken(config) {
  const nextToken = await refreshAccessToken()

  return http.request({
    ...config,
    __isRetryRequest: true,
    headers: {
      ...(config.headers || {}),
      Authorization: `Bearer ${nextToken}`,
    },
  })
}

async function handleRejectedResponse(error, config = {}) {
  if (shouldRefreshRequest(config, error)) {
    try {
      return await retryWithFreshToken(config)
    } catch (refreshError) {
      const normalizedRefreshError = normalizeError(refreshError)
      await handleAuthFailure(normalizedRefreshError, config)
      return Promise.reject(normalizedRefreshError)
    }
  }

  if (error?.isAuthError) {
    await handleAuthFailure(error, config)
  }

  return Promise.reject(error)
}

http.interceptors.response.use(
  async (response) => {
    const payload = response.data

    if (typeof payload?.code === 'number' && payload.code !== 0) {
      const error = createApiError(getPayloadMessage(payload) || '请求失败', {
        code: payload.code,
        status: response.status,
        isAuthError: isAuthErrorCode(payload.code),
        payload,
      })

      return handleRejectedResponse(error, response.config)
    }

    return payload
  },
  async (error) => {
    const normalizedError = normalizeError(error)
    return handleRejectedResponse(normalizedError, error?.config)
  },
)

export default http
