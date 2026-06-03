<script setup>
const props = defineProps({
  pagination: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['change'])

function go(page) {
  if (page < 1 || page > (props.pagination.pages || 1) || page === props.pagination.page) {
    return
  }
  emit('change', page)
}
</script>

<template>
  <div v-if="pagination.pages > 1" class="pagination-bar">
    <button type="button" class="secondary-btn" @click="go(pagination.page - 1)">上一页</button>
    <span>第 {{ pagination.page }} / {{ pagination.pages }} 页，共 {{ pagination.total }} 条</span>
    <button type="button" class="secondary-btn" @click="go(pagination.page + 1)">下一页</button>
  </div>
</template>

<style scoped>
.pagination-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  margin-top: 20px;
}

@media (max-width: 640px) {
  .pagination-bar {
    flex-direction: column;
  }
}
</style>
