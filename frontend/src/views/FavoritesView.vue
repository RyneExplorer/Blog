<script setup>
import { computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'

import ArticleListCard from '@/components/article/ArticleListCard.vue'
import PaginationBar from '@/components/article/PaginationBar.vue'
import { useArticleStore } from '@/stores/articleStore'

const route = useRoute()
const router = useRouter()
const articleStore = useArticleStore()
const { categories, favoriteArticles, favoritePagination, favoriteFilters, favoriteLoading } = storeToRefs(articleStore)

const selectedCategory = computed(() => route.query.category_id || '')

async function syncList() {
  await articleStore.loadCategories()
  await articleStore.loadFavoriteArticles({
    page: Number(route.query.page || 1),
    category_id: selectedCategory.value,
  })
}

watch(() => route.fullPath, syncList)
onMounted(syncList)

function updateQuery(patch) {
  router.push({
    name: 'favorites',
    query: {
      ...route.query,
      ...patch,
    },
  })
}
</script>

<template>
  <section class="favorite-view">
    <div class="page-section favorite-header">
      <div class="section-heading">
        <div>
          <h1>我的收藏</h1>
          <p>这里展示你收藏过的公开文章。</p>
        </div>
      </div>

      <div class="chip-row">
        <button
          type="button"
          class="chip-button"
          :class="{ 'is-active': favoriteFilters.category_id === '' }"
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
    </div>

    <div v-if="favoriteLoading" class="empty-box">收藏列表加载中...</div>
    <div v-else-if="!favoriteArticles.length" class="empty-box">你还没有收藏任何文章</div>
    <div v-else class="favorite-stack">
      <ArticleListCard v-for="article in favoriteArticles" :key="article.id" :article="article" />
    </div>

    <PaginationBar :pagination="favoritePagination" @change="updateQuery({ page: $event })" />
  </section>
</template>

<style scoped>
.favorite-view,
.favorite-stack {
  display: grid;
  gap: 18px;
}

.favorite-header {
  padding: 24px;
}
</style>
