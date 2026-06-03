<script setup>
import { computed } from 'vue'

const props = defineProps({
  status: {
    type: Number,
    required: true,
  },
})

const text = computed(() => {
  return {
    0: '草稿',
    1: '待审核',
    2: '已发布',
    3: '已拒绝',
    4: '已封禁',
  }[props.status] || '未知状态'
})
</script>

<template>
  <span class="admin-status" :class="`admin-status--${status}`">{{ text }}</span>
</template>

<style scoped>
.admin-status {
  display: inline-flex;
  align-items: center;
  padding: 6px 12px;
  border-radius: 999px;
  font-size: 13px;
}

.admin-status--1 {
  background: rgba(255, 179, 0, 0.14);
  color: #8b6300;
}

.admin-status--2 {
  background: rgba(93, 163, 158, 0.12);
  color: var(--color-primary-dark);
}

.admin-status--3,
.admin-status--4 {
  background: var(--color-danger-soft);
  color: var(--color-danger);
}

.admin-status--0 {
  background: rgba(108, 124, 122, 0.12);
  color: #4b5a58;
}
</style>
