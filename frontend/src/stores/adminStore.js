import { reactive, ref } from 'vue'
import { defineStore } from 'pinia'

import {
  approveReview,
  banReview,
  getAdminReviewDetail,
  getAdminReviewList,
  getAdminUserList,
  rejectReview,
  updateReviewCategory,
} from '@/api/admin'
import { useUiStore } from '@/stores/uiStore'
import {
  createPagination,
  mapAdminReviewDetail,
  mapArticleListItem,
  mapAdminUser,
} from '@/utils/formatters'

export const useAdminStore = defineStore('admin', () => {
  const uiStore = useUiStore()

  const reviewList = ref([])
  const reviewPagination = reactive(createPagination(10))
  const reviewFilters = reactive({
    status: '',
    username: '',
    category_id: '',
  })
  const reviewLoading = ref(false)
  const reviewDetail = ref(null)
  const reviewDetailLoading = ref(false)
  const reviewActionLoading = ref(false)

  const userList = ref([])
  const userPagination = reactive(createPagination(10))
  const userFilters = reactive({
    username: '',
    nickname: '',
    status: '',
  })
  const userLoading = ref(false)

  async function loadReviews(params = {}) {
    reviewLoading.value = true

    try {
      const query = {
        page: params.page ?? reviewPagination.page,
        pageSize: params.pageSize ?? reviewPagination.pageSize,
      }
      const status = params.status ?? reviewFilters.status
      const username = params.username ?? reviewFilters.username
      const categoryId = params.category_id ?? reviewFilters.category_id

      if (status !== '') {
        query.status = Number(status)
      }
      if (username) {
        query.username = username
      }
      if (categoryId) {
        query.category_id = Number(categoryId)
      }

      const page = await getAdminReviewList(query)
      reviewList.value = (page.list || []).map(mapArticleListItem)
      Object.assign(reviewPagination, page)
      reviewFilters.status = status
      reviewFilters.username = username
      reviewFilters.category_id = categoryId
    } catch (error) {
      reviewList.value = []
      uiStore.showToast(error.message, 'error')
    } finally {
      reviewLoading.value = false
    }
  }

  async function loadReviewDetail(id) {
    reviewDetailLoading.value = true

    try {
      const data = await getAdminReviewDetail(id)
      reviewDetail.value = mapAdminReviewDetail(data)
    } catch (error) {
      reviewDetail.value = null
      uiStore.showToast(error.message, 'error')
    } finally {
      reviewDetailLoading.value = false
    }
  }

  async function refreshReviewState(id) {
    await Promise.all([loadReviewDetail(id), loadReviews()])
  }

  async function handleApprove(id) {
    reviewActionLoading.value = true

    try {
      await approveReview(id)
      uiStore.showToast('审核已通过', 'success')
      await refreshReviewState(id)
    } catch (error) {
      uiStore.showToast(error.message, 'error')
    } finally {
      reviewActionLoading.value = false
    }
  }

  async function handleReject(id, reason) {
    const trimmedReason = `${reason || ''}`.trim()
    if (!trimmedReason) {
      uiStore.showToast('请输入驳回原因', 'error')
      return
    }

    reviewActionLoading.value = true

    try {
      await rejectReview(id, trimmedReason)
      uiStore.showToast('文章已驳回', 'success')
      await refreshReviewState(id)
    } catch (error) {
      uiStore.showToast(error.message, 'error')
    } finally {
      reviewActionLoading.value = false
    }
  }

  async function handleBan(id, reason) {
    const trimmedReason = `${reason || ''}`.trim()
    if (!trimmedReason) {
      uiStore.showToast('请输入封禁原因', 'error')
      return
    }

    reviewActionLoading.value = true

    try {
      await banReview(id, trimmedReason)
      uiStore.showToast('文章已封禁', 'success')
      await refreshReviewState(id)
    } catch (error) {
      uiStore.showToast(error.message, 'error')
    } finally {
      reviewActionLoading.value = false
    }
  }

  async function handleUpdateCategory(id, categoryIds) {
    reviewActionLoading.value = true

    try {
      await updateReviewCategory(id, categoryIds)
      uiStore.showToast('分类已更新', 'success')
      await refreshReviewState(id)
    } catch (error) {
      uiStore.showToast(error.message, 'error')
    } finally {
      reviewActionLoading.value = false
    }
  }

  async function loadUsers(params = {}) {
    userLoading.value = true

    try {
      const query = {
        page: params.page ?? userPagination.page,
        pageSize: params.pageSize ?? userPagination.pageSize,
      }
      const username = params.username ?? userFilters.username
      const nickname = params.nickname ?? userFilters.nickname
      const status = params.status ?? userFilters.status

      if (username) {
        query.username = username
      }
      if (nickname) {
        query.nickname = nickname
      }
      if (status !== '') {
        query.status = Number(status)
      }

      const page = await getAdminUserList(query)
      userList.value = (page.list || []).map(mapAdminUser)
      Object.assign(userPagination, page)
      userFilters.username = username
      userFilters.nickname = nickname
      userFilters.status = status
    } catch (error) {
      userList.value = []
      uiStore.showToast(error.message, 'error')
    } finally {
      userLoading.value = false
    }
  }

  return {
    reviewList,
    reviewPagination,
    reviewFilters,
    reviewLoading,
    reviewDetail,
    reviewDetailLoading,
    reviewActionLoading,
    userList,
    userPagination,
    userFilters,
    userLoading,
    loadReviews,
    loadReviewDetail,
    handleApprove,
    handleReject,
    handleBan,
    handleUpdateCategory,
    loadUsers,
  }
})
