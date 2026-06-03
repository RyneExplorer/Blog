import { computed, reactive, ref, watch } from 'vue'
import { defineStore } from 'pinia'

import {
  createArticle,
  deleteArticle,
  favoriteArticle,
  getArticleList,
  getFavoriteArticleList,
  getMyArticleDetail,
  getMyArticleList,
  getPublicArticleDetail,
  likeArticle,
  publishArticle,
  recordArticleView,
  unfavoriteArticle,
  unlikeArticle,
  updateArticle,
  uploadCoverImage,
} from '@/api/article'
import { getCategoryList } from '@/api/category'
import { useAuthStore } from '@/stores/authStore'
import { useUiStore } from '@/stores/uiStore'
import {
  createPagination,
  mapArticleDetail,
  mapArticleListItem,
  mapCategory,
} from '@/utils/formatters'

function createEditorForm() {
  return {
    id: 0,
    title: '',
    summary: '',
    cover_image: '',
    category_ids: [],
    content: '',
  }
}

function hasExplicitFlag(value) {
  return typeof value === 'boolean'
}

export const useArticleStore = defineStore('article', () => {
  const uiStore = useUiStore()
  const authStore = useAuthStore()

  const categories = ref([])
  const categoriesLoaded = ref(false)
  let categoriesRequest = null

  const publicArticles = ref([])
  const publicPagination = reactive(createPagination(8))
  const publicFilters = reactive({
    category_id: '',
    sort: 'latest',
  })
  const publicLoading = ref(false)

  const articleDetail = ref(null)
  const detailLoading = ref(false)

  const myArticles = ref([])
  const myArticlePagination = reactive(createPagination(10))
  const myArticleFilters = reactive({
    category_id: '',
    sort: 'latest',
  })
  const myArticlesLoading = ref(false)
  const myArticleStatusFilter = ref('')

  const favoriteArticles = ref([])
  const favoritePagination = reactive(createPagination(10))
  const favoriteFilters = reactive({
    category_id: '',
    sort: 'latest',
  })
  const favoriteLoading = ref(false)

  const editorForm = reactive(createEditorForm())
  const editorMode = ref('create')
  const submitIntent = ref('draft')
  const editorLoading = ref(false)
  const editorSaving = ref(false)

  const likedArticleIds = ref(new Set())
  const favoritedArticleIds = ref(new Set())
  const interactionUserId = computed(() => Number(authStore.currentUser?.id || 0))

  const filteredMyArticles = computed(() => {
    if (myArticleStatusFilter.value === '') {
      return myArticles.value
    }

    return myArticles.value.filter((article) => String(article.statusCode) === String(myArticleStatusFilter.value))
  })

  // restoreInteractionState 切换用户时清空本地交互状态
  function restoreInteractionState() {
    // 1. 文章点赞/收藏红色状态以后端 liked/favorited 字段为准。
    // 2. 清空旧用户状态，避免退出或切换账号后沿用上一个用户的 icon 状态。
    likedArticleIds.value = new Set()
    favoritedArticleIds.value = new Set()
  }

  watch(interactionUserId, restoreInteractionState, { immediate: true })

  function applyKnownInteractionState(article, options = {}) {
    if (!article?.id) {
      return
    }

    const targetId = Number(article.id)

    if (hasExplicitFlag(article.liked)) {
      const next = new Set(likedArticleIds.value)
      if (article.liked) {
        next.add(targetId)
      } else {
        next.delete(targetId)
      }
      likedArticleIds.value = next
    }

    const favorited = options.forceFavorited === true ? true : article.favorited
    if (hasExplicitFlag(favorited)) {
      const next = new Set(favoritedArticleIds.value)
      if (favorited) {
        next.add(targetId)
      } else {
        next.delete(targetId)
      }
      favoritedArticleIds.value = next
    }
  }

  function isAlreadyActiveError(error, action) {
    const message = String(error?.message || '')
    const status = Number(error?.status || 0)
    const code = Number(error?.code || 0)

    if (status === 409 || code === 409) {
      return true
    }

    if (action === 'like') {
      return /\u5df2\u70b9\u8d5e|\u91cd\u590d\u70b9\u8d5e|duplicate/i.test(message)
    }

    return /\u5df2\u6536\u85cf|\u91cd\u590d\u6536\u85cf|duplicate/i.test(message)
  }

  async function loadCategories() {
    if (categoriesLoaded.value) {
      return
    }

    if (!categoriesRequest) {
      // 1. 布局和页面会同时需要分类数据。
      // 2. 在首个请求完成前复用同一个 Promise，避免刷新时重复请求 /api/categories。
      // 3. 请求失败时清空 Promise，允许用户后续重新加载。
      categoriesRequest = getCategoryList()
        .then((data) => {
          categories.value = (data || []).map(mapCategory)
          categoriesLoaded.value = true
        })
        .catch((error) => {
          uiStore.showToast(error.message, 'error')
          throw error
        })
        .finally(() => {
          categoriesRequest = null
        })
    }

    return categoriesRequest
  }

  async function loadPublicArticles(params = {}) {
    publicLoading.value = true

    try {
      const query = {
        page: params.page ?? publicPagination.page,
        pageSize: params.pageSize ?? publicPagination.pageSize,
        sort: params.sort ?? publicFilters.sort,
      }

      const categoryId = params.category_id ?? publicFilters.category_id
      if (categoryId) {
        query.category_id = categoryId
      }

      const page = await getArticleList(query)
      publicArticles.value = (page.list || []).map(mapArticleListItem)
      publicArticles.value.forEach((article) => applyKnownInteractionState(article))
      Object.assign(publicPagination, page)
      publicFilters.sort = query.sort
      publicFilters.category_id = categoryId || ''
    } catch (error) {
      publicArticles.value = []
      uiStore.showToast(error.message, 'error')
    } finally {
      publicLoading.value = false
    }
  }

  async function loadArticleDetail(id) {
    detailLoading.value = true

    try {
      const data = await getPublicArticleDetail(id)
      articleDetail.value = mapArticleDetail(data)
      applyKnownInteractionState(articleDetail.value)
      await recordArticleView(id)
      articleDetail.value.viewCount += 1
    } catch (error) {
      articleDetail.value = null
      uiStore.showToast(error.message, 'error')
    } finally {
      detailLoading.value = false
    }
  }

  async function loadMyArticles(params = {}) {
    myArticlesLoading.value = true

    try {
      const query = {
        page: params.page ?? myArticlePagination.page,
        pageSize: params.pageSize ?? myArticlePagination.pageSize,
        sort: params.sort ?? myArticleFilters.sort,
      }
      const categoryId = params.category_id ?? myArticleFilters.category_id

      if (categoryId) {
        query.category_id = categoryId
      }

      const page = await getMyArticleList(query)
      myArticles.value = (page.list || []).map(mapArticleListItem)
      myArticles.value.forEach((article) => applyKnownInteractionState(article))
      Object.assign(myArticlePagination, page)
      myArticleFilters.sort = query.sort
      myArticleFilters.category_id = categoryId || ''
    } catch (error) {
      myArticles.value = []
      uiStore.showToast(error.message, 'error')
    } finally {
      myArticlesLoading.value = false
    }
  }

  async function loadFavoriteArticles(params = {}) {
    favoriteLoading.value = true

    try {
      const query = {
        page: params.page ?? favoritePagination.page,
        pageSize: params.pageSize ?? favoritePagination.pageSize,
        sort: params.sort ?? favoriteFilters.sort,
      }
      const categoryId = params.category_id ?? favoriteFilters.category_id

      if (categoryId) {
        query.category_id = categoryId
      }

      const page = await getFavoriteArticleList(query)
      favoriteArticles.value = (page.list || []).map(mapArticleListItem)
      favoriteArticles.value.forEach((article) => applyKnownInteractionState(article, { forceFavorited: true }))
      Object.assign(favoritePagination, page)
      favoriteFilters.sort = query.sort
      favoriteFilters.category_id = categoryId || ''
    } catch (error) {
      favoriteArticles.value = []
      uiStore.showToast(error.message, 'error')
    } finally {
      favoriteLoading.value = false
    }
  }

  function resetEditor() {
    Object.assign(editorForm, createEditorForm())
    editorMode.value = 'create'
    submitIntent.value = 'draft'
  }

  async function loadEditorArticle(id) {
    resetEditor()
    if (!id) {
      return
    }

    editorLoading.value = true

    try {
      const data = await getMyArticleDetail(id)
      editorMode.value = 'edit'
      editorForm.id = data?.id || 0
      editorForm.title = data?.title || ''
      editorForm.summary = data?.summary || ''
      editorForm.cover_image = data?.cover_image || ''
      editorForm.category_ids = data?.category_ids || []
      editorForm.content = data?.content || ''
    } catch (error) {
      uiStore.showToast(error.message, 'error')
    } finally {
      editorLoading.value = false
    }
  }

  async function uploadCover(file) {
    try {
      const data = await uploadCoverImage(file)
      editorForm.cover_image = data?.cover_image || data?.url || ''
      uiStore.showToast('\u5c01\u9762\u4e0a\u4f20\u6210\u529f', 'success')
      return editorForm.cover_image
    } catch (error) {
      uiStore.showToast(error.message, 'error')
      return ''
    }
  }

  async function saveEditorDraft() {
    editorSaving.value = true

    try {
      const payload = {
        title: editorForm.title,
        summary: editorForm.summary,
        cover_image: editorForm.cover_image,
        category_ids: editorForm.category_ids,
        content: editorForm.content,
      }

      let articleId = editorForm.id
      if (editorForm.id) {
        await updateArticle(editorForm.id, payload)
      } else {
        const data = await createArticle(payload)
        articleId = Number(data?.article_id || 0)
        editorForm.id = articleId
        editorMode.value = 'edit'
      }

      uiStore.showToast('\u8349\u7a3f\u5df2\u4fdd\u5b58', 'success')
      return articleId
    } catch (error) {
      uiStore.showToast(error.message, 'error')
      return 0
    } finally {
      editorSaving.value = false
    }
  }

  async function submitEditorReview() {
    if (!editorForm.title.trim() || !editorForm.content.trim()) {
      uiStore.showToast('\u63d0\u4ea4\u5ba1\u6838\u524d\u8bf7\u586b\u5199\u6807\u9898\u548c\u6b63\u6587', 'error')
      return false
    }

    const articleId = await saveEditorDraft()
    if (!articleId) {
      return false
    }

    try {
      await publishArticle(articleId)
      uiStore.showToast('\u5df2\u63d0\u4ea4\u5ba1\u6838', 'success')
      return true
    } catch (error) {
      uiStore.showToast(error.message, 'error')
      return false
    }
  }

  async function removeArticle(id) {
    try {
      await deleteArticle(id)
      uiStore.showToast('\u6587\u7ae0\u5df2\u5220\u9664', 'success')
      await loadMyArticles()
    } catch (error) {
      uiStore.showToast(error.message, 'error')
    }
  }

  async function submitArticle(id) {
    try {
      await publishArticle(id)
      uiStore.showToast('\u5df2\u63d0\u4ea4\u5ba1\u6838', 'success')
      await loadMyArticles()
    } catch (error) {
      uiStore.showToast(error.message, 'error')
    }
  }

  function updateArticleCollection(collection, id, updater) {
    const targetId = Number(id)
    collection.value = collection.value.map((article) => {
      if (Number(article.id) !== targetId) {
        return article
      }

      return updater(article)
    })
  }

  function syncArticleLikeState(id, active, options = {}) {
    const targetId = Number(id)
    const next = new Set(likedArticleIds.value)

    if (active) {
      next.add(targetId)
    } else {
      next.delete(targetId)
    }
    likedArticleIds.value = next

    if (options.adjustCount === false) {
      return
    }

    const update = (article) => ({
      ...article,
      likeCount: Math.max(0, Number(article.likeCount || 0) + (active ? 1 : -1)),
    })

    updateArticleCollection(publicArticles, targetId, update)
    updateArticleCollection(myArticles, targetId, update)
    updateArticleCollection(favoriteArticles, targetId, update)

    if (Number(articleDetail.value?.id) === targetId) {
      articleDetail.value = {
        ...articleDetail.value,
        likeCount: Math.max(0, Number(articleDetail.value.likeCount || 0) + (active ? 1 : -1)),
      }
    }
  }

  function syncArticleFavoriteState(id, active, options = {}) {
    const targetId = Number(id)
    const next = new Set(favoritedArticleIds.value)

    if (active) {
      next.add(targetId)
    } else {
      next.delete(targetId)
    }
    favoritedArticleIds.value = next

    if (options.adjustCount === false) {
      return
    }

    const update = (article) => ({
      ...article,
      favoriteCount: Math.max(0, Number(article.favoriteCount || 0) + (active ? 1 : -1)),
    })

    updateArticleCollection(publicArticles, targetId, update)
    updateArticleCollection(myArticles, targetId, update)
    updateArticleCollection(favoriteArticles, targetId, update)

    if (Number(articleDetail.value?.id) === targetId) {
      articleDetail.value = {
        ...articleDetail.value,
        favoriteCount: Math.max(0, Number(articleDetail.value.favoriteCount || 0) + (active ? 1 : -1)),
      }
    }
  }

  async function toggleLike(id) {
    if (!authStore.isLoggedIn) {
      authStore.openAuth('login')
      return
    }

    try {
      if (likedArticleIds.value.has(id)) {
        await unlikeArticle(id)
        syncArticleLikeState(id, false)
        uiStore.showToast('\u5df2\u53d6\u6d88\u70b9\u8d5e', 'success')
      } else {
        await likeArticle(id)
        syncArticleLikeState(id, true)
        uiStore.showToast('点赞成功', 'success')
      }
    } catch (error) {
      if (isAlreadyActiveError(error, 'like')) {
        syncArticleLikeState(id, true, { adjustCount: false })
      }
      uiStore.showToast(error.message, 'error')
    }
  }

  async function toggleFavorite(id) {
    if (!authStore.isLoggedIn) {
      authStore.openAuth('login')
      return
    }

    try {
      if (favoritedArticleIds.value.has(id)) {
        await unfavoriteArticle(id)
        syncArticleFavoriteState(id, false)
        uiStore.showToast('\u5df2\u53d6\u6d88\u6536\u85cf', 'success')
      } else {
        await favoriteArticle(id)
        syncArticleFavoriteState(id, true)
        uiStore.showToast('收藏成功', 'success')
      }
    } catch (error) {
      if (isAlreadyActiveError(error, 'favorite')) {
        syncArticleFavoriteState(id, true, { adjustCount: false })
      }
      uiStore.showToast(error.message, 'error')
    }
  }

  return {
    categories,
    categoriesLoaded,
    publicArticles,
    publicPagination,
    publicFilters,
    publicLoading,
    articleDetail,
    detailLoading,
    myArticles,
    myArticlePagination,
    myArticleFilters,
    myArticlesLoading,
    myArticleStatusFilter,
    filteredMyArticles,
    favoriteArticles,
    favoritePagination,
    favoriteFilters,
    favoriteLoading,
    editorForm,
    editorMode,
    submitIntent,
    editorLoading,
    editorSaving,
    likedArticleIds,
    favoritedArticleIds,
    loadCategories,
    loadPublicArticles,
    loadArticleDetail,
    loadMyArticles,
    loadFavoriteArticles,
    resetEditor,
    loadEditorArticle,
    uploadCover,
    saveEditorDraft,
    submitEditorReview,
    removeArticle,
    submitArticle,
    toggleLike,
    toggleFavorite,
  }
})
