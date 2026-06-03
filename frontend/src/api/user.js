import http from './http'

export async function getProfile() {
  const payload = await http.get('/api/user/profile')
  return payload?.data || null
}

export async function updateProfile(data) {
  const payload = await http.put('/api/user/profile', data)
  return payload?.data || null
}

export async function uploadAvatar(file) {
  const formData = new FormData()
  formData.append('file', file)
  const payload = await http.post('/api/user/avatar', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })

  return payload?.data || null
}

export async function changePassword(data) {
  const payload = await http.post('/api/user/password', data)
  return payload?.data || null
}
