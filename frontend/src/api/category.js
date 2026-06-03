import http from './http'

export async function getCategoryList() {
  const payload = await http.get('/api/categories')

  return payload?.data || []
}
