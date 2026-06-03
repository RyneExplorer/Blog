<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'

import ArticleCard from '@/components/article/ArticleCard.vue'
import PaginationBar from '@/components/article/PaginationBar.vue'
import { useArticleStore } from '@/stores/articleStore'

const router = useRouter()
const articleStore = useArticleStore()
const { categories, publicArticles, publicPagination, publicFilters, publicLoading } = storeToRefs(articleStore)

onMounted(async () => {
  await articleStore.loadCategories()
  await articleStore.loadPublicArticles({ page: 1, pageSize: 6 })
})
</script>

<template>
  <section class="home-view">
    <div class="page-section hero-card">
      <div>
        <p class="hero-kicker">公开博客首页</p>
        <h1>内容、创作和审核都在一个前台里完成</h1>
        <p class="hero-copy">
          这里展示最近公开发布的文章。你可以继续进入文章详情、评论互动，或者登录后开始写作和管理自己的内容。
        </p>
      </div>
      <div class="hero-actions">
        <button type="button" class="primary-btn" @click="router.push({ name: 'articles' })">浏览全部文章</button>
        <button type="button" class="secondary-btn" @click="router.push({ name: 'editor-new' })">开始写作</button>
      </div>
    </div>

    <div class="chip-row filters-row">
      <button
        type="button"
        class="chip-button"
        :class="{ 'is-active': publicFilters.category_id === '' }"
        @click="articleStore.loadPublicArticles({ page: 1, category_id: '' })"
      >
        全部分类
      </button>
      <button
        v-for="category in categories"
        :key="category.id"
        type="button"
        class="chip-button"
        :class="{ 'is-active': String(publicFilters.category_id) === String(category.id) }"
        @click="articleStore.loadPublicArticles({ page: 1, category_id: category.id })"
      >
        {{ category.name }}
      </button>
      <button
        type="button"
        class="chip-button"
        :class="{ 'is-active': publicFilters.sort === 'latest' }"
        @click="articleStore.loadPublicArticles({ page: 1, sort: 'latest' })"
      >
        最新
      </button>
      <button
        type="button"
        class="chip-button"
        :class="{ 'is-active': publicFilters.sort === 'hottest' }"
        @click="articleStore.loadPublicArticles({ page: 1, sort: 'hottest' })"
      >
        最热
      </button>
    </div>

    <div class="article-stack">
      <div v-if="publicLoading" class="empty-box">文章加载中...</div>
      <div v-else-if="!publicArticles.length" class="empty-box">暂无公开文章</div>
      <ArticleCard v-for="article in publicArticles" :key="article.id" :article="article" />
    </div>

    <PaginationBar :pagination="publicPagination" @change="articleStore.loadPublicArticles({ page: $event })" />
  </section>
</template>

<style scoped>
.home-view {
  display: grid;
  gap: 24px;
}

.hero-card {
  padding: 28px;
  display: flex;
  justify-content: space-between;
  gap: 18px;
}

.hero-kicker {
  color: var(--color-primary);
  letter-spacing: 0.2em;
  text-transform: uppercase;
  font-size: 12px;
}

.hero-card h1 {
  margin-top: 12px;
  font-size: 40px;
  font-weight: 800;
  line-height: 1.15;
}

.hero-copy {
  margin-top: 14px;
  max-width: 720px;
  color: var(--color-muted);
}

.hero-actions {
  display: flex;
  align-items: flex-end;
  gap: 12px;
}

.filters-row,
.article-stack {
  display: grid;
  gap: 18px;
}

@media (max-width: 760px) {
  .hero-card {
    flex-direction: column;
  }

  .hero-card h1 {
    font-size: 30px;
  }

  .hero-actions {
    align-items: stretch;
    flex-direction: column;
  }
}
</style>
