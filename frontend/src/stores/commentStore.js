import { ref } from 'vue'
import { defineStore } from 'pinia'

import {
  createComment,
  createReply,
  deleteComment,
  getArticleComments,
  likeComment,
  unlikeComment,
} from '@/api/comment'
import { useAuthStore } from '@/stores/authStore'
import { useUiStore } from '@/stores/uiStore'
import { mapCommentNode } from '@/utils/formatters'

// collectLikedCommentIds 收集后端返回的已点赞评论 ID
function collectLikedCommentIds(nodes, result = new Set()) {
  nodes.forEach((node) => {
    // 1. 遍历一级评论和所有回复，保证嵌套评论也能恢复红色点赞状态。
    // 2. 只信任后端返回的 liked=true，避免浏览器本地缓存和数据库状态不一致。
    if (node?.comment?.liked === true) {
      result.add(Number(node.comment.id))
    }
    collectLikedCommentIds(node?.replies || [], result)
  })
  return result
}

export const useCommentStore = defineStore('comment', () => {
  const authStore = useAuthStore()
  const uiStore = useUiStore()

  const comments = ref([])
  const loading = ref(false)
  const rootComposer = ref('')
  const replyComposer = ref('')
  const replyTarget = ref(null)
  const replySheetVisible = ref(false)
  const likedCommentIds = ref(new Set())

  async function loadComments(articleId) {
    loading.value = true

    try {
      const page = await getArticleComments(articleId, {
        page: 1,
        pageSize: 100,
      })
      comments.value = (page.list || []).map(mapCommentNode)
      likedCommentIds.value = collectLikedCommentIds(comments.value)
    } catch (error) {
      comments.value = []
      uiStore.showToast(error.message, 'error')
    } finally {
      loading.value = false
    }
  }

  function startReply(node) {
    replyTarget.value = node
    replyComposer.value = ''
    replySheetVisible.value = true
  }

  function cancelReply() {
    replyTarget.value = null
    replyComposer.value = ''
    replySheetVisible.value = false
  }

  function updateCommentTree(nodes, id, updater) {
    return nodes.map((node) => {
      if (Number(node.comment.id) === Number(id)) {
        return {
          ...node,
          comment: updater(node.comment),
        }
      }

      if (!node.replies?.length) {
        return node
      }

      return {
        ...node,
        replies: updateCommentTree(node.replies, id, updater),
      }
    })
  }

  async function submitRootComment(articleId) {
    if (!authStore.isLoggedIn) {
      authStore.openAuth('login')
      return false
    }

    if (!rootComposer.value.trim()) {
      uiStore.showToast('评论内容不能为空', 'error')
      return false
    }

    try {
      await createComment({
        article_id: articleId,
        content: rootComposer.value,
      })

      rootComposer.value = ''
      uiStore.showToast('评论已发布', 'success')
      await loadComments(articleId)
      return true
    } catch (error) {
      uiStore.showToast(error.message, 'error')
      return false
    }
  }

  async function submitReply(articleId) {
    if (!authStore.isLoggedIn) {
      authStore.openAuth('login')
      return false
    }

    if (!replyTarget.value) {
      uiStore.showToast('未选择回复目标', 'error')
      return false
    }

    if (!replyComposer.value.trim()) {
      uiStore.showToast('回复内容不能为空', 'error')
      return false
    }

    try {
      await createReply({
        article_id: articleId,
        content: replyComposer.value,
        parent_id: replyTarget.value.comment.id,
        root_id: replyTarget.value.comment.rootId || replyTarget.value.comment.id,
      })

      replyComposer.value = ''
      replySheetVisible.value = false
      replyTarget.value = null
      uiStore.showToast('回复已发布', 'success')
      await loadComments(articleId)
      return true
    } catch (error) {
      uiStore.showToast(error.message, 'error')
      return false
    }
  }

  async function toggleLike(id) {
    if (!authStore.isLoggedIn) {
      authStore.openAuth('login')
      return
    }

    try {
      const next = new Set(likedCommentIds.value)

      if (likedCommentIds.value.has(id)) {
        await unlikeComment(id)
        next.delete(id)
        comments.value = updateCommentTree(comments.value, id, (comment) => ({
          ...comment,
          likeCount: Math.max(0, Number(comment.likeCount || 0) - 1),
          liked: false,
        }))
        uiStore.showToast('已取消点赞', 'success')
      } else {
        await likeComment(id)
        next.add(id)
        comments.value = updateCommentTree(comments.value, id, (comment) => ({
          ...comment,
          likeCount: Number(comment.likeCount || 0) + 1,
          liked: true,
        }))
        uiStore.showToast('点赞成功', 'success')
      }
      likedCommentIds.value = next
    } catch (error) {
      uiStore.showToast(error.message, 'error')
    }
  }

  async function removeComment(id, articleId) {
    const confirmed = await uiStore.openConfirm({
      title: '确认删除这条评论？',
      confirmText: '确认删除',
      cancelText: '取消',
      intent: 'danger',
    })

    if (!confirmed) {
      return false
    }

    try {
      await deleteComment(id)
      uiStore.showToast('评论已删除', 'success')
      await loadComments(articleId)
      return true
    } catch (error) {
      uiStore.showToast(error.message, 'error')
      return false
    }
  }

  return {
    comments,
    loading,
    rootComposer,
    replyComposer,
    replyTarget,
    replySheetVisible,
    likedCommentIds,
    loadComments,
    startReply,
    cancelReply,
    submitRootComment,
    submitReply,
    toggleLike,
    removeComment,
  }
})
