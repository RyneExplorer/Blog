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
const { categories, publicArticles, publicPagination, publicFilters, publicLoading } = storeToRefs(articleStore)

const selectedCategory = computed(() => route.query.category_id || '')
const selectedSort = computed(() => route.query.sort || 'latest')

async function syncList() {
  await articleStore.loadCategories()
  await articleStore.loadPublicArticles({
    page: Number(route.query.page || 1),
    category_id: selectedCategory.value,
    sort: selectedSort.value,
  })
}

watch(() => route.fullPath, syncList)

onMounted(syncList)

function updateQuery(patch) {
  router.push({
    name: 'articles',
    query: {
      ...route.query,
      ...patch,
    },
  })
}
</script>

<template>
  <section class="list-view">
    <div class="page-section list-header">
      <div class="section-heading">
        <div>
          <h1>文章列表</h1>
          <p>按分类、热度和时间浏览所有公开文章。</p>
        </div>
      </div>

      <div class="chip-row">
        <button type="button" class="chip-button" :class="{ 'is-active': selectedCategory === '' }" @click="updateQuery({ category_id: '', page: 1 })">
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
        <button type="button" class="chip-button" :class="{ 'is-active': selectedSort === 'latest' }" @click="updateQuery({ sort: 'latest', page: 1 })">
          最新
        </button>
        <button type="button" class="chip-button" :class="{ 'is-active': selectedSort === 'hottest' }" @click="updateQuery({ sort: 'hottest', page: 1 })">
          最热
        </button>
      </div>
    </div>

    <div class="article-stack">
      <div v-if="publicLoading" class="empty-box">文章加载中...</div>
      <div v-else-if="!publicArticles.length" class="empty-box">当前筛选条件下暂无文章</div>
      <ArticleListCard v-for="article in publicArticles" :key="article.id" :article="article" />
    </div>

    <PaginationBar :pagination="publicPagination" @change="updateQuery({ page: $event })" />
    <p v-if="!publicLoading && publicArticles.length" class="bottom-tip">已经浏览到底了~</p>
  </section>
</template>

<style scoped>
.list-view,
.article-stack {
  display: grid;
  gap: 18px;
}

.list-header {
  padding: 24px;
}

.bottom-tip {
  text-align: center;
  color: var(--color-muted);
  padding: 8px 0 24px;
}
</style>
