import http from './http'

export async function getCaptcha() {
  const payload = await http.get('/api/auth/captcha')
  return payload?.data || null
}

export async function login(data) {
  const payload = await http.post('/api/auth/login', data)
  return payload?.data || null
}

export async function refreshToken(token) {
  const payload = await http.post('/api/auth/refresh', { token }, {
    skipAuthRefresh: true,
    skipAuthFailureHandler: true,
  })
  return payload?.data || null
}

export async function logout() {
  const payload = await http.post('/api/auth/logout', null, {
    skipAuthRefresh: true,
    skipAuthFailureHandler: true,
  })
  return payload?.data || null
}

export async function register(data) {
  const payload = await http.post('/api/auth/register', data)
  return payload?.data || null
}

export async function sendEmailCode(email) {
  const payload = await http.post('/api/auth/email/code', { email })
  return payload?.data || null
}

export async function resetPassword(data) {
  const payload = await http.post('/api/auth/password/reset', data)
  return payload?.data || null
}
