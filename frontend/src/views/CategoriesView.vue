<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'

import { useArticleStore } from '@/stores/articleStore'

const router = useRouter()
const articleStore = useArticleStore()
const { categories } = storeToRefs(articleStore)

onMounted(() => {
  articleStore.loadCategories()
})
</script>

<template>
  <section class="page-section category-view">
    <div class="section-heading">
      <div>
        <h1>分类导航</h1>
        <p>按分类进入公开文章列表。</p>
      </div>
    </div>

    <div class="category-grid">
      <button
        v-for="category in categories"
        :key="category.id"
        type="button"
        class="category-card"
        @click="router.push({ name: 'articles', query: { category_id: category.id } })"
      >
        <strong>{{ category.name }}</strong>
        <span>{{ category.slug || '未设置 slug' }}</span>
      </button>
    </div>
  </section>
</template>

<style scoped>
.category-view {
  padding: 24px;
}

.category-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
}

.category-card {
  min-height: 120px;
  border: 1px solid var(--color-border);
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.95);
  padding: 18px;
  text-align: left;
  cursor: pointer;
}

.category-card strong {
  display: block;
  font-size: 20px;
  font-weight: 700;
}

.category-card span {
  display: block;
  margin-top: 10px;
  color: var(--color-muted);
}
</style>
