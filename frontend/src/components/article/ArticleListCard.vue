<script setup>
import { computed } from 'vue'

const props = defineProps({
  article: {
    type: Object,
    required: true,
  },
})

const authorFallback = computed(() => {
  const name = props.article.authorName || '作者'
  return name.trim().slice(0, 1)
})
</script>

<template>
  <article class="article-list-card">
    <RouterLink :to="{ name: 'article-detail', params: { id: article.id } }" class="article-list-card__link">
      <div class="article-list-card__cover">
        <img v-if="article.coverImage" :src="article.coverImage" :alt="article.title" />
        <div v-else class="article-list-card__fallback">
          <span>{{ article.title }}</span>
        </div>
      </div>

      <div class="article-list-card__content">
        <header class="article-list-card__header">
          <h2 class="article-list-card__title">{{ article.title }}</h2>
          <p class="article-list-card__summary">{{ article.summary }}</p>
        </header>

        <footer class="article-list-card__footer">
          <div class="article-list-card__author">
            <div class="article-list-card__avatar">
              <img v-if="article.authorAvatar" :src="article.authorAvatar" :alt="article.authorName" />
              <span v-else>{{ authorFallback }}</span>
            </div>
            <span class="article-list-card__author-name">{{ article.authorName || '匿名作者' }}</span>
          </div>

          <div class="article-list-card__meta" aria-label="文章元信息">
            <span class="article-list-card__meta-item">
              <svg viewBox="0 0 24 24" aria-hidden="true">
                <rect x="3.5" y="5" width="17" height="15.5" rx="2.5" />
                <path d="M7.5 3.5v3M16.5 3.5v3M3.5 9.5h17" />
              </svg>
              {{ article.createdDate || article.createdAt }}
            </span>

            <span class="article-list-card__meta-item">
              <svg viewBox="0 0 24 24" aria-hidden="true">
                <path d="M2.5 12s3.5-6 9.5-6 9.5 6 9.5 6-3.5 6-9.5 6-9.5-6-9.5-6Z" />
                <circle cx="12" cy="12" r="3" />
              </svg>
              {{ article.viewCount || 0 }}
            </span>

            <span class="article-list-card__meta-item">
              <svg viewBox="0 0 24 24" aria-hidden="true">
                <path d="M4 6.75A2.75 2.75 0 0 1 6.75 4h10.5A2.75 2.75 0 0 1 20 6.75v7.5A2.75 2.75 0 0 1 17.25 17H10l-4.25 3v-3H6.75A2.75 2.75 0 0 1 4 14.25Z" />
              </svg>
              {{ article.commentCount || 0 }}
            </span>

            <span class="article-list-card__meta-item">
              <svg viewBox="0 0 24 24" aria-hidden="true">
                <path d="m12 20.5-1.45-1.32C5.4 14.5 2 11.42 2 7.63 2 4.55 4.42 2.2 7.5 2.2c1.74 0 3.41.8 4.5 2.05 1.09-1.25 2.76-2.05 4.5-2.05 3.08 0 5.5 2.35 5.5 5.43 0 3.8-3.4 6.87-8.55 11.55L12 20.5Z" />
              </svg>
              {{ article.likeCount || 0 }}
            </span>

            <span class="article-list-card__meta-item">
              <svg viewBox="0 0 24 24" aria-hidden="true">
                <path d="m12 2.8 2.86 5.79 6.39.93-4.62 4.5 1.09 6.36L12 17.35 6.28 20.38l1.09-6.36-4.62-4.5 6.39-.93L12 2.8Z" />
              </svg>
              {{ article.favoriteCount || 0 }}
            </span>
          </div>
        </footer>
      </div>
    </RouterLink>
  </article>
</template>

<style scoped>
.article-list-card {
  border-radius: 22px;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.96), rgba(255, 251, 244, 0.9));
  box-shadow:
    0 1px 0 rgba(24, 65, 66, 0.05),
    0 10px 24px rgba(35, 47, 49, 0.08);
  transition:
    transform 0.2s ease,
    box-shadow 0.2s ease;
}

.article-list-card:hover {
  transform: translateY(-3px);
  box-shadow:
    0 1px 0 rgba(24, 65, 66, 0.08),
    0 16px 32px rgba(35, 47, 49, 0.12);
}

.article-list-card__link {
  display: grid;
  grid-template-columns: 168px minmax(0, 1fr);
  gap: 18px;
  align-items: center;
  padding: 16px 18px;
}

.article-list-card__cover {
  position: relative;
  width: 168px;
  overflow: hidden;
  border-radius: 18px;
  background: linear-gradient(135deg, rgba(93, 163, 158, 0.14), rgba(24, 65, 66, 0.14));
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.18);
}

.article-list-card__cover::before {
  content: '';
  display: block;
  padding-top: 56.25%;
}

.article-list-card__cover img,
.article-list-card__fallback {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
}

.article-list-card__cover img {
  object-fit: cover;
}

.article-list-card__fallback {
  display: grid;
  place-items: center;
  padding: 14px;
  text-align: center;
  color: rgba(24, 65, 66, 0.76);
  font-size: 0.86rem;
  font-weight: 700;
  line-height: 1.45;
}

.article-list-card__content {
  min-width: 0;
  display: grid;
  gap: 14px;
}

.article-list-card__header {
  display: grid;
  gap: 8px;
}

.article-list-card__title {
  color: #183f45;
  font-size: 1.05rem;
  font-weight: 800;
  line-height: 1.35;
  letter-spacing: -0.01em;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 1;
  overflow: hidden;
}

.article-list-card__summary {
  color: #72808d;
  font-size: 0.92rem;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
}

.article-list-card__footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px 22px;
  padding-top: 12px;
  border-top: 1px solid rgba(93, 109, 109, 0.1);
}

.article-list-card__author {
  min-width: 0;
  display: inline-flex;
  align-items: center;
  gap: 10px;
}

.article-list-card__avatar {
  width: 26px;
  height: 26px;
  flex-shrink: 0;
  display: grid;
  place-items: center;
  overflow: hidden;
  border-radius: 999px;
  background: linear-gradient(135deg, var(--color-primary-dark), var(--color-primary));
  color: #fff;
  font-size: 0.74rem;
  font-weight: 700;
  box-shadow: 0 6px 12px rgba(24, 65, 66, 0.16);
}

.article-list-card__avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.article-list-card__author-name {
  min-width: 0;
  color: #30404f;
  font-size: 0.95rem;
  font-weight: 700;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.article-list-card__meta {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 8px 14px;
  color: #8592a0;
  font-size: 0.82rem;
}

.article-list-card__meta-item {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  white-space: nowrap;
}

.article-list-card__meta-item svg {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
  fill: none;
  stroke: currentColor;
  stroke-width: 1.8;
  stroke-linecap: round;
  stroke-linejoin: round;
}

@media (max-width: 900px) {
  .article-list-card__link {
    grid-template-columns: 1fr;
    gap: 14px;
  }

  .article-list-card__cover {
    width: 100%;
  }
}

@media (max-width: 640px) {
  .article-list-card {
    border-radius: 18px;
  }

  .article-list-card__link {
    padding: 14px;
  }

  .article-list-card__cover {
    border-radius: 14px;
  }

  .article-list-card__title {
    font-size: 1rem;
  }

  .article-list-card__summary {
    font-size: 0.88rem;
  }

  .article-list-card__footer {
    flex-direction: column;
    align-items: flex-start;
  }

  .article-list-card__meta {
    justify-content: flex-start;
  }
}
</style>
