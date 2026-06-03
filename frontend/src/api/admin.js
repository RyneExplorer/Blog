import http, { normalizePage } from './http'

export async function getAdminReviewList(params) {
  const payload = await http.get('/api/super/articles', { params })
  return normalizePage(payload, params?.pageSize)
}

export async function getAdminReviewDetail(id) {
  const payload = await http.get(`/api/super/articles/${id}`)
  return payload?.data || null
}

export async function approveReview(id) {
  const payload = await http.post(`/api/super/articles/${id}/approve`)
  return payload?.data || null
}

export async function rejectReview(id, reason) {
  const payload = await http.post(`/api/super/articles/${id}/reject`, { reason })
  return payload?.data || null
}

export async function banReview(id, reason) {
  const payload = await http.post(`/api/super/articles/${id}/ban`, { reason })
  return payload?.data || null
}

export async function updateReviewCategory(id, categoryIds) {
  const payload = await http.put(`/api/super/articles/${id}/category`, {
    category_ids: categoryIds,
  })

  return payload?.data || null
}

export async function getAdminUserList(params) {
  const payload = await http.get('/api/super/userlist', { params })
  return normalizePage(payload, params?.pageSize)
}
