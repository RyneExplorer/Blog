<script setup>
import { computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'

import PaginationBar from '@/components/article/PaginationBar.vue'
import { useArticleStore } from '@/stores/articleStore'

const route = useRoute()
const router = useRouter()
const articleStore = useArticleStore()
const {
  categories,
  filteredMyArticles,
  myArticlesLoading,
  myArticlePagination,
  myArticleFilters,
  myArticleStatusFilter,
} = storeToRefs(articleStore)

const selectedCategory = computed(() => route.query.category_id || '')

async function syncList() {
  await articleStore.loadCategories()
  await articleStore.loadMyArticles({
    page: Number(route.query.page || 1),
    category_id: selectedCategory.value,
  })
}

watch(() => route.fullPath, syncList)
onMounted(syncList)

function updateQuery(patch) {
  router.push({
    name: 'my-articles',
    query: {
      ...route.query,
      ...patch,
    },
  })
}

const statusOptions = [
  { label: '全部', value: '' },
  { label: '草稿', value: '0' },
  { label: '待审核', value: '1' },
  { label: '已发布', value: '2' },
  { label: '已驳回', value: '3' },
  { label: '已封禁', value: '4' },
]
</script>

<template>
  <section class="page-section page-panel">
    <div class="section-heading">
      <div>
        <h1>我的文章</h1>
        <p>管理草稿、待审核和已发布内容。</p>
      </div>
      <button type="button" class="primary-btn" @click="router.push({ name: 'editor-new' })">新建文章</button>
    </div>

    <div class="chip-row">
      <button
        type="button"
        class="chip-button"
        :class="{ 'is-active': myArticleFilters.category_id === '' }"
        @click="updateQuery({ category_id: '', page: 1 })"
      >
        全部分类
      </button>
      <button
        v-for="category in categories"
        :key="category.id"
        type="button"
        class="chip-button"
        :class="{ 'is-active': String(selectedCategory) === String(category.id) }"
        @click="updateQuery({ category_id: category.id, page: 1 })"
      >
        {{ category.name }}
      </button>
    </div>

    <div class="chip-row">
      <button
        v-for="item in statusOptions"
        :key="item.label"
        type="button"
        class="chip-button"
        :class="{ 'is-active': myArticleStatusFilter === item.value }"
        @click="myArticleStatusFilter = item.value"
      >
        {{ item.label }}
      </button>
    </div>

    <div v-if="myArticlesLoading" class="empty-box">文章加载中...</div>
    <div v-else-if="!filteredMyArticles.length" class="empty-box">暂无文章</div>

    <div v-else class="item-stack">
      <article v-for="article in filteredMyArticles" :key="article.id" class="item-card">
        <div>
          <h3>{{ article.title }}</h3>
          <p class="muted-text">{{ article.summary }}</p>
          <div class="chip-row item-meta">
            <span class="chip">{{ article.statusLabel }}</span>
            <span class="chip">{{ article.categoryName }}</span>
            <span class="chip">{{ article.updatedAt }}</span>
          </div>
        </div>
        <div class="item-actions">
          <button type="button" class="secondary-btn" @click="router.push({ name: 'editor-edit', params: { id: article.id } })">
            继续编辑
          </button>
          <button
            v-if="article.statusCode === 0 || article.statusCode === 3"
            type="button"
            class="primary-btn"
            @click="articleStore.submitArticle(article.id)"
          >
            提交审核
          </button>
          <button type="button" class="danger-btn" @click="articleStore.removeArticle(article.id)">删除</button>
        </div>
      </article>
    </div>

    <PaginationBar :pagination="myArticlePagination" @change="updateQuery({ page: $event })" />
  </section>
</template>

<style scoped>
.page-panel {
  padding: 24px;
}

.item-stack {
  display: grid;
  gap: 16px;
}

.item-card {
  display: flex;
  justify-content: space-between;
  gap: 18px;
  padding: 18px;
  border: 1px solid var(--color-border);
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.92);
}

.item-card h3 {
  font-size: 22px;
  font-weight: 700;
}

.item-meta {
  margin-top: 12px;
}

.item-actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
}

@media (max-width: 760px) {
  .item-card {
    flex-direction: column;
  }
}
</style>
