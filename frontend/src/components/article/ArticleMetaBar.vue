<script setup>
defineProps({
  article: {
    type: Object,
    required: true,
  },
  liked: {
    type: Boolean,
    default: false,
  },
  favorited: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['like-click', 'favorite-click'])
</script>

<template>
  <div class="meta-bar">
    <span class="meta-bar__tag">原创</span>
    <span>发布博客 {{ article.createdDate || article.createdAt }}</span>
    <span>{{ article.viewCount }} 阅读</span>
    <button type="button" class="meta-bar__action" :class="{ 'is-active': liked }" @click.stop.prevent="emit('like-click')">
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
      <span>{{ article.likeCount }} 点赞</span>
    </button>
    <span>{{ article.commentCount }} 评论</span>
    <button type="button" class="meta-bar__action" :class="{ 'is-active': favorited }" @click.stop.prevent="emit('favorite-click')">
      <svg
        viewBox="0 0 24 24"
        :fill="favorited ? 'currentColor' : 'none'"
        stroke="currentColor"
        stroke-width="1.9"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <path d="m12 2.8 2.86 5.79 6.39.93-4.62 4.5 1.09 6.36L12 17.35 6.28 20.38l1.09-6.36-4.62-4.5 6.39-.93L12 2.8Z" />
      </svg>
      <span>{{ article.favoriteCount }} 收藏</span>
    </button>
  </div>
</template>

<style scoped>
.meta-bar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 0;
  color: #7b8798;
  font-size: 13px;
  line-height: 1.2;
}

.meta-bar > span {
  position: relative;
  padding-right: 12px;
  margin-right: 12px;
}

.meta-bar > span:not(:last-child)::after,
.meta-bar__action::after {
  content: '·';
  position: absolute;
  right: 0;
  color: #b2bcc9;
}

.meta-bar__tag {
  color: #ff6b5f;
  background: rgba(255, 107, 95, 0.08);
  border-radius: 4px;
  padding: 3px 6px;
  margin-right: 12px;
  font-size: 12px;
}

.meta-bar__tag::after {
  top: 50%;
  transform: translateY(-50%);
}

.meta-bar__action {
  position: relative;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 0 12px 0 0;
  margin-right: 12px;
  border: none;
  background: transparent;
  color: #7b8798;
  font-size: 13px;
  line-height: 1.2;
  cursor: pointer;
  transition:
    transform 0.18s ease,
    color 0.18s ease;
}

.meta-bar__action:hover {
  transform: translateY(-1px);
  color: #5e6c81;
}

.meta-bar__action.is-active {
  color: #e85757;
}

.meta-bar__action.is-active:hover {
  color: #dd4545;
}

.meta-bar__action svg {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
  transition:
    fill 0.18s ease,
    stroke 0.18s ease,
    transform 0.18s ease;
}

.meta-bar__action.is-active svg {
  transform: scale(1.03);
}

@media (max-width: 760px) {
  .meta-bar {
    gap: 8px 0;
  }
}
</style>
