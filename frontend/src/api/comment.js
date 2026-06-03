import http, { normalizePage } from './http'

export async function getArticleComments(articleId, params) {
  const payload = await http.get(`/api/articles/${articleId}/comments`, { params })
  return normalizePage(payload, params?.pageSize)
}

export async function createComment(data) {
  const payload = await http.post('/api/comments', data)
  return payload?.data || null
}

export async function createReply(data) {
  const payload = await http.post('/api/comments/replies', data)
  return payload?.data || null
}

export async function likeComment(id) {
  const payload = await http.post(`/api/comments/${id}/like`)
  return payload?.data || null
}

export async function unlikeComment(id) {
  const payload = await http.post(`/api/comments/${id}/unlike`)
  return payload?.data || null
}

export async function deleteComment(id) {
  const payload = await http.delete(`/api/comments/${id}`)
  return payload?.data || null
}
