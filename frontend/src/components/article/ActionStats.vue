<script setup>
import { computed } from 'vue'

const props = defineProps({
  views: {
    type: Number,
    default: 0,
  },
  comments: {
    type: Number,
    default: 0,
  },
  favorites: {
    type: Number,
    default: 0,
  },
  likes: {
    type: Number,
    default: 0,
  },
  favorited: {
    type: Boolean,
    default: false,
  },
  liked: {
    type: Boolean,
    default: false,
  },
  clickable: {
    type: Boolean,
    default: true,
  },
})

const emit = defineEmits(['views-click', 'comments-click', 'favorites-click', 'likes-click'])

function formatCount(value) {
  const numeric = Number(value || 0)
  if (numeric >= 1000) {
    return `${(numeric / 1000).toFixed(numeric >= 10000 ? 0 : 1).replace(/\.0$/, '')}k`
  }
  return String(numeric)
}

const items = computed(() => [
  {
    key: 'views',
    label: formatCount(props.views),
    subLabel: '',
    onClick: () => emit('views-click'),
  },
  {
    key: 'comments',
    label: formatCount(props.comments),
    subLabel: '条评论',
    onClick: () => emit('comments-click'),
  },
  {
    key: 'favorites',
    label: formatCount(props.favorites),
    subLabel: '',
    active: props.favorited,
    onClick: () => emit('favorites-click'),
  },
  {
    key: 'likes',
    label: formatCount(props.likes),
    subLabel: '',
    active: props.liked,
    onClick: () => emit('likes-click'),
  },
])
</script>

<template>
  <div class="action-stats">
    <button
      v-for="item in items"
      :key="item.key"
      type="button"
      class="action-stats__item"
      :class="{
        'is-clickable': clickable,
        'is-active': item.active,
      }"
      :disabled="!clickable"
      @click="item.onClick"
    >
      <span class="action-stats__icon">
        <svg
          v-if="item.key === 'views'"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.9"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M2 12s3.6-6 10-6 10 6 10 6-3.6 6-10 6-10-6-10-6Z" />
          <circle cx="12" cy="12" r="3" />
        </svg>
        <svg
          v-else-if="item.key === 'comments'"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.9"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M21 11.5a8.5 8.5 0 0 1-8.5 8.5c-1.42 0-2.76-.34-3.94-.95L3 21l1.83-4.96A8.47 8.47 0 0 1 4 11.5 8.5 8.5 0 1 1 21 11.5Z" />
        </svg>
        <svg
          v-else-if="item.key === 'favorites'"
          viewBox="0 0 24 24"
          :fill="item.active ? 'currentColor' : 'none'"
          stroke="currentColor"
          stroke-width="1.9"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="m12 2.8 2.86 5.79 6.39.93-4.62 4.5 1.09 6.36L12 17.35 6.28 20.38l1.09-6.36-4.62-4.5 6.39-.93L12 2.8Z" />
        </svg>
        <svg
          v-else
          viewBox="0 0 24 24"
          :fill="item.active ? 'currentColor' : 'none'"
          stroke="currentColor"
          stroke-width="1.9"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="m12 20.5-1.45-1.32C5.4 14.5 2 11.42 2 7.63 2 4.55 4.42 2.2 7.5 2.2c1.74 0 3.41.8 4.5 2.05 1.09-1.25 2.76-2.05 4.5-2.05 3.08 0 5.5 2.35 5.5 5.43 0 3.8-3.4 6.87-8.55 11.55L12 20.5Z" />
        </svg>
      </span>

      <span class="action-stats__text">
        <span>{{ item.label }}</span>
        <span v-if="item.subLabel" class="action-stats__sub">{{ item.subLabel }}</span>
      </span>
    </button>
  </div>
</template>

<style scoped>
.action-stats {
  display: flex;
  align-items: center;
  gap: 18px;
  flex-wrap: wrap;
}

.action-stats__item {
  display: inline-flex;
  align-items: center;
  justify-content: flex-start;
  gap: 6px;
  min-height: 28px;
  padding: 2px 0;
  border: none;
  background: transparent;
  color: rgb(132, 145, 167);
  line-height: 1.1;
  transition:
    transform 0.2s ease,
    color 0.2s ease,
    background-color 0.2s ease,
    opacity 0.2s ease;
}

.action-stats__item.is-clickable {
  cursor: pointer;
}

.action-stats__item:hover {
  transform: translateY(-2px);
  color: rgb(98, 115, 144);
}

.action-stats__item.is-active {
  color: #e85757;
}

.action-stats__item.is-active:hover {
  color: #dd4545;
}

.action-stats__item.is-active .action-stats__icon svg {
  transform: scale(1.03);
}

.action-stats__icon {
  width: 16px;
  height: 16px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.action-stats__icon svg {
  width: 16px;
  height: 16px;
  transition:
    fill 0.2s ease,
    stroke 0.2s ease,
    transform 0.2s ease;
}

.action-stats__text {
  display: inline-flex;
  align-items: baseline;
  gap: 4px;
  font-size: 13px;
  white-space: nowrap;
}

.action-stats__sub {
  font-size: 13px;
}

@media (max-width: 640px) {
  .action-stats {
    gap: 14px;
  }
}
</style>
