<script setup>
import { computed, ref } from 'vue'
import { storeToRefs } from 'pinia'

import { useAuthStore } from '@/stores/authStore'
import { useCommentStore } from '@/stores/commentStore'

defineOptions({
  name: 'CommentItem',
})

const props = defineProps({
  node: {
    type: Object,
    required: true,
  },
  articleId: {
    type: Number,
    required: true,
  },
  depth: {
    type: Number,
    default: 0,
  },
})

const authStore = useAuthStore()
const commentStore = useCommentStore()
const { currentUser } = storeToRefs(authStore)
const showAllReplies = ref(false)

const canDelete = computed(() => Number(currentUser.value?.id || 0) === props.node.author.id)
const liked = computed(() => commentStore.likedCommentIds.has(props.node.comment.id))
const visibleReplies = computed(() => {
  if (props.depth > 0 || showAllReplies.value) {
    return props.node.replies
  }

  return props.node.replies.slice(0, 2)
})
const hiddenReplyCount = computed(() => Math.max(0, props.node.replies.length - visibleReplies.value.length))
</script>

<template>
  <article class="comment-item" :class="{ 'is-nested': depth > 0 }">
    <div class="comment-item__row">
      <div class="comment-author__avatar">
        <img v-if="node.author.avatar" :src="node.author.avatar" :alt="node.author.nickname" />
        <span v-else>{{ node.author.nickname.slice(0, 1) }}</span>
      </div>

      <div class="comment-item__main">
        <div class="comment-item__bubble">
          <div class="comment-item__content-line">
            <strong>{{ node.author.nickname }}</strong>
            <div class="comment-item__content" v-html="node.comment.content"></div>
          </div>
        </div>

        <div class="comment-item__meta">
          <span>{{ node.comment.createdAt }}</span>
          <button type="button" class="comment-action" @click="commentStore.startReply(node)">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round">
              <path d="M9 17 4 12l5-5" />
              <path d="M20 18v-2a4 4 0 0 0-4-4H4" />
            </svg>
            <span>回复</span>
          </button>
          <button type="button" class="comment-action" :class="{ 'is-active': liked }" @click="commentStore.toggleLike(node.comment.id)">
            <svg
              viewBox="0 0 24 24"
              :fill="liked ? 'currentColor' : 'none'"
              stroke="currentColor"
              stroke-width="1.9"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="m12 20.5-1.45-1.32C5.4 14.5 2 11.42 2 7.63 2 4.55 4.42 2.2 7.5 2.2c1.74 0 3.41.8 4.5 2.05 1.09-1.25 2.76-2.05 4.5-2.05 3.08 0 5.5 2.35 5.5 5.43 0 3.8-3.4 6.87-8.55 11.55L12 20.5Z" />
            </svg>
            <span>{{ node.comment.likeCount }}</span>
          </button>
          <button
            v-if="canDelete"
            type="button"
            class="comment-action comment-action--danger"
            @click="commentStore.removeComment(node.comment.id, articleId)"
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round">
              <path d="M3 6h18" />
              <path d="M8 6V4.8c0-.66.54-1.2 1.2-1.2h5.6c.66 0 1.2.54 1.2 1.2V6" />
              <path d="m18 6-1 13.2c-.05.71-.64 1.25-1.35 1.25H8.35c-.71 0-1.3-.54-1.35-1.25L6 6" />
              <path d="M10 10.2v6.6" />
              <path d="M14 10.2v6.6" />
            </svg>
            <span>删除</span>
          </button>
        </div>

        <div v-if="node.replies.length" class="comment-item__replies">
          <CommentItem
            v-for="reply in visibleReplies"
            :key="reply.comment.id"
            :node="reply"
            :article-id="articleId"
            :depth="depth + 1"
          />

          <button
            v-if="hiddenReplyCount > 0 && !showAllReplies"
            type="button"
            class="comment-item__more"
            @click="showAllReplies = true"
          >
            展开更多回复（{{ hiddenReplyCount }}）
          </button>
        </div>
      </div>
    </div>
  </article>
</template>

<style scoped>
.comment-item {
  padding: 16px 0;
  border-bottom: 1px solid rgba(93, 109, 109, 0.08);
}

.comment-item.is-nested {
  padding: 12px 0 0 24px;
  position: relative;
  border-bottom: none;
}

.comment-item.is-nested::before {
  content: '';
  position: absolute;
  left: 8px;
  top: 4px;
  bottom: 6px;
  width: 1px;
  background: rgba(93, 109, 109, 0.14);
}

.comment-item__row {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}

.comment-author__avatar {
  width: 36px;
  height: 36px;
  border-radius: 999px;
  overflow: hidden;
  background: var(--color-primary-soft);
  color: var(--color-primary-dark);
  display: grid;
  place-items: center;
  flex-shrink: 0;
  font-size: 14px;
  font-weight: 700;
}

.comment-author__avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.comment-item__main {
  min-width: 0;
  flex: 1;
}

.comment-item__bubble {
  min-width: 0;
}

.comment-item__content-line {
  color: #22313f;
  line-height: 1.6;
  word-break: break-word;
}

.comment-item__content-line strong {
  margin-right: 8px;
  font-size: 14px;
  font-weight: 700;
  color: #1f2d3d;
}

.comment-item__content {
  display: inline;
  color: #33465a;
  font-size: 14px;
}

:deep(.comment-item__content p) {
  display: inline;
  margin: 0;
}

.comment-item__meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px;
  margin-top: 6px;
  color: #8d98a8;
  font-size: 12px;
  line-height: 1;
}

.comment-action {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  border: none;
  background: transparent;
  padding: 0;
  color: #8d98a8;
  font-size: 12px;
  line-height: 1;
  cursor: pointer;
  transition:
    transform 0.18s ease,
    color 0.18s ease,
    opacity 0.18s ease;
}

.comment-action:hover {
  transform: translateY(-1px);
  color: #697789;
}

.comment-action svg {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
  transition:
    fill 0.18s ease,
    stroke 0.18s ease,
    transform 0.18s ease;
}

.comment-action.is-active {
  color: #e85757;
}

.comment-action.is-active:hover {
  color: #dd4545;
}

.comment-action.is-active svg {
  transform: scale(1.03);
}

.comment-action--danger {
  color: #d45a5a;
}

.comment-action--danger:hover {
  color: #bc4343;
}

.comment-item__replies {
  margin-top: 10px;
  display: grid;
  gap: 8px;
}

.comment-item__more {
  justify-self: flex-start;
  border: none;
  background: transparent;
  padding: 0;
  color: var(--color-primary-dark);
  font-size: 12px;
  cursor: pointer;
  transition: color 0.18s ease;
}

.comment-item__more:hover {
  color: var(--color-primary);
}

@media (max-width: 640px) {
  .comment-item.is-nested {
    padding-left: 16px;
  }

  .comment-item.is-nested::before {
    left: 4px;
  }

  .comment-item__meta {
    gap: 10px;
  }
}
</style>
