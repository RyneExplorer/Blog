<script setup>
import { computed } from 'vue'

import ArticleMetaBar from '@/components/article/ArticleMetaBar.vue'
import { useArticleStore } from '@/stores/articleStore'

const props = defineProps({
  article: {
    type: Object,
    required: true,
  },
})

const articleStore = useArticleStore()
const liked = computed(() => articleStore.likedArticleIds.has(props.article.id))
const favorited = computed(() => articleStore.favoritedArticleIds.has(props.article.id))
</script>

<template>
  <RouterLink :to="{ name: 'article-detail', params: { id: article.id } }" class="article-card page-section">
    <div class="article-card__cover">
      <img v-if="article.coverImage" :src="article.coverImage" :alt="article.title" />
      <div v-else class="article-card__fallback">{{ article.title }}</div>
    </div>

    <div class="article-card__content">
      <h3>{{ article.title }}</h3>

      <p class="article-card__summary">{{ article.summary }}</p>
      <ArticleMetaBar
        :article="article"
        :liked="liked"
        :favorited="favorited"
        @like-click="articleStore.toggleLike(article.id)"
        @favorite-click="articleStore.toggleFavorite(article.id)"
      />
    </div>
  </RouterLink>
</template>

<style scoped>
.article-card {
  display: grid;
  grid-template-columns: 296px minmax(0, 1fr);
  align-items: start;
  gap: 24px;
  padding: 0;
  color: inherit;
  background: #fff;
  border: none;
  border-radius: 0;
  box-shadow: none;
  transition:
    opacity 0.18s ease,
    transform 0.18s ease;
  will-change: transform;
}

.article-card:hover {
  transform: translateY(-2px);
  opacity: 0.96;
}

.article-card__content {
  min-width: 0;
  padding: 6px 0;
}

.article-card__content h3 {
  font-size: 22px;
  font-weight: 700;
  line-height: 1.35;
  color: #1f2d3d;
  transition: color 0.18s ease;
}

.article-card:hover .article-card__content h3 {
  color: #134b97;
}

.article-card__summary {
  margin: 8px 0 14px;
  color: #4c5d73;
  font-size: 16px;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-card__cover {
  width: 296px;
  height: 166px;
  border-radius: 4px;
  overflow: hidden;
  background: #f2f4f7;
}

.article-card__cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.article-card__fallback {
  height: 100%;
  padding: 18px;
  display: grid;
  place-items: center;
  text-align: center;
  color: #6e7d90;
  font-weight: 700;
}

@media (max-width: 760px) {
  .article-card {
    grid-template-columns: 1fr;
    gap: 14px;
  }

  .article-card__cover {
    width: 100%;
    height: 200px;
  }

  .article-card__content h3 {
    font-size: 20px;
  }

  .article-card__summary {
    font-size: 15px;
  }
}
</style>
