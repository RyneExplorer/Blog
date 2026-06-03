import { resolveAssetUrl } from '@/api/http'

export function createPagination(pageSize = 10) {
  return {
    list: [],
    total: 0,
    page: 1,
    pageSize,
    pages: 0,
  }
}

export function formatTime(value) {
  if (!value) {
    return '刚刚'
  }

  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }

  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function formatDateOnly(value) {
  if (!value) {
    return ''
  }

  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }

  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}.${month}.${day}`
}

export function formatArticleStatus(status) {
  const map = {
    0: '草稿',
    1: '待审核',
    2: '已发布',
    3: '已拒绝',
    4: '已封禁',
  }

  return map[status] || '未知状态'
}

export function mapCategory(item) {
  return {
    id: item?.id || 0,
    name: item?.name || '未分类',
    slug: item?.slug || '',
  }
}

export function mapArticleListItem(item) {
  const article = item?.article || {}
  const author = item?.author || {}
  const category = item?.category || {}
  const liked = article?.liked ?? article?.is_liked ?? item?.liked ?? item?.is_liked
  const favorited = article?.favorited ?? article?.is_favorited ?? item?.favorited ?? item?.is_favorited

  return {
    id: article.id || 0,
    title: article.title || '未命名文章',
    summary: article.summary || '暂无摘要。',
    coverImage: resolveAssetUrl(article.cover_image || ''),
    statusCode: article.status ?? -1,
    statusLabel: formatArticleStatus(article.status),
    createdAt: formatTime(article.created_at),
    createdDate: formatDateOnly(article.created_at),
    updatedAt: formatTime(article.updated_at),
    viewCount: article.view_count || 0,
    likeCount: article.like_count || 0,
    favoriteCount: article.favorite_count || 0,
    commentCount: article.comment_count || 0,
    authorId: author.id || 0,
    authorName: author.nickname || author.username || '匿名作者',
    authorAvatar: resolveAssetUrl(author.avatar || ''),
    authorBio: author.bio || '',
    categoryId: category.id || 0,
    categoryName: category.name || '未分类',
    categorySlug: category.slug || '',
    liked: typeof liked === 'boolean' ? liked : null,
    favorited: typeof favorited === 'boolean' ? favorited : null,
  }
}

export function mapArticleDetail(item) {
  const liked = item?.liked ?? item?.is_liked
  const favorited = item?.favorited ?? item?.is_favorited

  return {
    id: item?.id || 0,
    title: item?.title || '未命名文章',
    summary: item?.summary || '',
    content: item?.content || '',
    coverImage: resolveAssetUrl(item?.cover_image || ''),
    statusCode: item?.status ?? -1,
    statusLabel: formatArticleStatus(item?.status),
    viewCount: item?.view_count || 0,
    likeCount: item?.like_count || 0,
    favoriteCount: item?.favorite_count || 0,
    commentCount: item?.comment_count || 0,
    categoryIds: item?.category_ids || [],
    categoryName: item?.category_name || '',
    authorName: item?.nickname || item?.username || '',
    authorAvatar: resolveAssetUrl(item?.avatar || ''),
    authorBio: item?.bio || '',
    createdAt: formatTime(item?.created_at),
    updatedAt: formatTime(item?.updated_at),
    liked: typeof liked === 'boolean' ? liked : null,
    favorited: typeof favorited === 'boolean' ? favorited : null,
  }
}

export function mapCommentNode(node) {
  const comment = node?.comment || {}
  const author = node?.author || {}
  const liked = comment?.liked ?? comment?.is_liked ?? node?.liked ?? node?.is_liked

  return {
    comment: {
      id: comment.id || 0,
      parentId: comment.parent_id || 0,
      rootId: comment.root_id || 0,
      content: comment.content || '',
      likeCount: comment.like_count || 0,
      replyCount: comment.reply_count || 0,
      liked: typeof liked === 'boolean' ? liked : null,
      createdAt: formatTime(comment.created_at),
      updatedAt: formatTime(comment.updated_at),
    },
    author: {
      id: author.id || 0,
      nickname: author.nickname || '匿名用户',
      avatar: resolveAssetUrl(author.avatar || ''),
    },
    replies: (node?.replies || []).map(mapCommentNode),
  }
}

export function mapAdminReviewDetail(item) {
  return {
    article: mapArticleDetail(item?.article || {}),
    author: {
      id: item?.author?.id || 0,
      nickname: item?.author?.nickname || '未知作者',
      avatar: resolveAssetUrl(item?.author?.avatar || ''),
      bio: item?.author?.bio || '',
    },
    categories: (item?.categories || []).map(mapCategory),
    rejectReason: item?.reject_reason || '',
  }
}

export function mapAdminUser(item) {
  return {
    id: item?.id || 0,
    username: item?.username || '',
    email: item?.email || '',
    avatar: resolveAssetUrl(item?.avatar || ''),
    role: item?.role ?? 1,
    nickname: item?.nickname || '',
    status: item?.status ?? 1,
    createdAt: formatTime(item?.created_at),
    updatedAt: formatTime(item?.updated_at),
  }
}
