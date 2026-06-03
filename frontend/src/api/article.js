import http, { normalizePage } from './http'

export async function getArticleList(params) {
  const payload = await http.get('/api/articles', { params })
  return normalizePage(payload, params?.pageSize)
}

export async function getHotArticleList(params = {}) {
  return getArticleList({
    page: 1,
    pageSize: 5,
    sort: 'hottest',
    ...params,
  })
}

export async function getMyArticleList(params) {
  const payload = await http.get('/api/articles/mine', { params })
  return normalizePage(payload, params?.pageSize)
}

export async function getFavoriteArticleList(params) {
  const payload = await http.get('/api/articles/favorites', { params })
  return normalizePage(payload, params?.pageSize)
}

export async function getPublicArticleDetail(id) {
  const payload = await http.get(`/api/articles/${id}`)
  return payload?.data || null
}

export async function getMyArticleDetail(id) {
  const payload = await http.get(`/api/articles/mine/${id}`)
  return payload?.data || null
}

export async function createArticle(data) {
  const payload = await http.post('/api/articles', data)
  return payload?.data || null
}

export async function updateArticle(id, data) {
  const payload = await http.put(`/api/articles/${id}`, data)
  return payload?.data || null
}

export async function uploadCoverImage(file) {
  const formData = new FormData()
  formData.append('file', file)
  const payload = await http.post('/api/articles/cover_image', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })

  return payload?.data || null
}

export async function uploadContentImage(file) {
  const formData = new FormData()
  formData.append('file', file)
  const payload = await http.post('/api/articles/content_image', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })

  return payload?.data || null
}

export async function publishArticle(id) {
  const payload = await http.post(`/api/articles/${id}/publish`)
  return payload?.data || null
}

export async function deleteArticle(id) {
  const payload = await http.delete(`/api/articles/${id}`)
  return payload?.data || null
}

export async function recordArticleView(id) {
  const payload = await http.post(`/api/articles/${id}/view`)
  return payload?.data || null
}

export async function likeArticle(id) {
  const payload = await http.post(`/api/articles/${id}/like`)
  return payload?.data || null
}

export async function unlikeArticle(id) {
  const payload = await http.post(`/api/articles/${id}/unlike`)
  return payload?.data || null
}

export async function favoriteArticle(id) {
  const payload = await http.post(`/api/articles/${id}/favorite`)
  return payload?.data || null
}

export async function unfavoriteArticle(id) {
  const payload = await http.post(`/api/articles/${id}/unfavorite`)
  return payload?.data || null
}
