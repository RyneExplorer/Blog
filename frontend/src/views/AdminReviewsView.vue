<script setup>
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'

import PaginationBar from '@/components/article/PaginationBar.vue'
import { useAdminStore } from '@/stores/adminStore'
import { useArticleStore } from '@/stores/articleStore'

const router = useRouter()
const adminStore = useAdminStore()
const articleStore = useArticleStore()
const { reviewList, reviewPagination, reviewFilters, reviewLoading, reviewActionLoading } = storeToRefs(adminStore)
const { categories } = storeToRefs(articleStore)

const reviewStatusOptions = computed(() => [
  { label: '全部状态', value: '' },
  { label: '草稿', value: '0' },
  { label: '待审核', value: '1' },
  { label: '已发布', value: '2' },
  { label: '已驳回', value: '3' },
  { label: '已封禁', value: '4' },
])

onMounted(async () => {
  await articleStore.loadCategories()
  await adminStore.loadReviews()
})

function openDetail(id) {
  router.push({ name: 'admin-review-detail', params: { id } })
}

function handleCardKeydown(event, id) {
  if (event.key === 'Enter' || event.key === ' ') {
    event.preventDefault()
    openDetail(id)
  }
}

function quickApprove(event, id) {
  event.stopPropagation()
  adminStore.handleApprove(id)
}
</script>

<template>
  <section class="page-section admin-review-view">
    <div class="section-heading">
      <div>
        <h1>审核中心</h1>
        <p>管理员可以在这里查看并处理待审核文章。</p>
      </div>
    </div>

    <div class="filter-grid">
      <input v-model="reviewFilters.username" class="field" placeholder="作者用户名" />

      <select v-model="reviewFilters.status" class="select-field">
        <option v-for="option in reviewStatusOptions" :key="option.value" :value="option.value">
          {{ option.label }}
        </option>
      </select>

      <select v-model="reviewFilters.category_id" class="select-field">
        <option value="">全部分类</option>
        <option v-for="category in categories" :key="category.id" :value="category.id">
          {{ category.name }}
        </option>
      </select>

      <button type="button" class="primary-btn" @click="adminStore.loadReviews({ page: 1 })">筛选</button>
    </div>

    <div v-if="reviewLoading" class="empty-box">审核列表加载中...</div>
    <div v-else-if="!reviewList.length" class="empty-box">暂无审核文章</div>

    <div v-else class="review-stack">
      <article
        v-for="item in reviewList"
        :key="item.id"
        class="review-card"
        tabindex="0"
        @click="openDetail(item.id)"
        @keydown="handleCardKeydown($event, item.id)"
      >
        <div class="review-card__body">
          <div class="review-card__cover">
            <img v-if="item.coverImage" :src="item.coverImage" :alt="item.title" />
            <div v-else class="review-card__fallback">
              <span>{{ item.title }}</span>
            </div>
          </div>

          <div class="review-card__content">
            <header class="review-card__header">
              <h2 class="review-card__title">{{ item.title }}</h2>
              <p class="review-card__summary">{{ item.summary }}</p>
            </header>

            <footer class="review-card__footer">
              <div class="review-card__author">
                <div class="review-card__avatar">
                  <img v-if="item.authorAvatar" :src="item.authorAvatar" :alt="item.authorName" />
                  <span v-else>{{ (item.authorName || '匿').slice(0, 1) }}</span>
                </div>

                <div class="review-card__author-copy">
                  <span class="review-card__author-name">{{ item.authorName || '匿名作者' }}</span>
                  <span class="review-card__author-bio">{{ item.authorBio || '这个作者还没有填写个人简介。' }}</span>
                </div>
              </div>

              <div class="review-card__meta" aria-label="审核文章元信息">
                <span class="review-card__meta-item">
                  <svg viewBox="0 0 24 24" aria-hidden="true">
                    <path d="M12 2.8 14.86 8.59l6.39.93-4.62 4.5 1.09 6.36L12 17.35 6.28 20.38l1.09-6.36-4.62-4.5 6.39-.93L12 2.8Z" />
                  </svg>
                  待审核
                </span>

                <span class="review-card__meta-item">
                  <svg viewBox="0 0 24 24" aria-hidden="true">
                    <path d="M12 12c2.76 0 5-2.46 5-5.5S14.76 1 12 1 7 3.46 7 6.5 9.24 12 12 12Z" />
                    <path d="M4 21c0-4.14 3.58-7.5 8-7.5s8 3.36 8 7.5" />
                  </svg>
                  {{ item.authorName || '匿名作者' }}
                </span>

                <span class="review-card__meta-item">
                  <svg viewBox="0 0 24 24" aria-hidden="true">
                    <path d="M4 6.5h16M7.5 3.5v3M16.5 3.5v3" />
                    <rect x="3.5" y="5.5" width="17" height="15" rx="2.5" />
                  </svg>
                  {{ item.updatedAt }}
                </span>

                <span class="review-card__meta-item">
                  <svg viewBox="0 0 24 24" aria-hidden="true">
                    <path d="M5 7.5h14M5 12h10M5 16.5h8" />
                    <path d="M4.5 4.5h15v15h-15z" opacity=".001" />
                  </svg>
                  {{ item.categoryName || '未分类' }}
                </span>
              </div>
            </footer>
          </div>
        </div>

        <div class="review-card__actions">
          <button
            type="button"
            class="primary-btn review-card__approve"
            :disabled="reviewActionLoading"
            @click="quickApprove($event, item.id)"
          >
            {{ reviewActionLoading ? '处理中...' : '审核通过' }}
          </button>
        </div>
      </article>
    </div>

    <PaginationBar :pagination="reviewPagination" @change="adminStore.loadReviews({ page: $event })" />
  </section>
</template>

<style scoped>
.admin-review-view {
  padding: 24px;
}

.filter-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
  margin-bottom: 20px;
}

.review-stack {
  display: grid;
  gap: 16px;
}

.review-card {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 18px;
  align-items: center;
  border-radius: 22px;
  padding: 16px 18px;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.96), rgba(255, 251, 244, 0.9));
  box-shadow:
    0 1px 0 rgba(24, 65, 66, 0.05),
    0 10px 24px rgba(35, 47, 49, 0.08);
  cursor: pointer;
  transition:
    transform 0.2s ease,
    box-shadow 0.2s ease;
}

.review-card:hover,
.review-card:focus-visible {
  transform: translateY(-3px);
  box-shadow:
    0 1px 0 rgba(24, 65, 66, 0.08),
    0 16px 32px rgba(35, 47, 49, 0.12);
  outline: none;
}

.review-card__body {
  display: grid;
  grid-template-columns: 168px minmax(0, 1fr);
  gap: 18px;
  align-items: center;
  min-width: 0;
}

.review-card__cover {
  position: relative;
  width: 168px;
  overflow: hidden;
  border-radius: 18px;
  background: linear-gradient(135deg, rgba(93, 163, 158, 0.14), rgba(24, 65, 66, 0.14));
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.18);
}

.review-card__cover::before {
  content: '';
  display: block;
  padding-top: 56.25%;
}

.review-card__cover img,
.review-card__fallback {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
}

.review-card__cover img {
  object-fit: cover;
}

.review-card__fallback {
  display: grid;
  place-items: center;
  padding: 14px;
  text-align: center;
  color: rgba(24, 65, 66, 0.76);
  font-size: 0.86rem;
  font-weight: 700;
  line-height: 1.45;
}

.review-card__content {
  min-width: 0;
  display: grid;
  gap: 14px;
}

.review-card__header {
  display: grid;
  gap: 8px;
}

.review-card__title {
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

.review-card__summary {
  color: #72808d;
  font-size: 0.92rem;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
}

.review-card__footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px 22px;
  padding-top: 12px;
  border-top: 1px solid rgba(93, 109, 109, 0.1);
}

.review-card__author {
  min-width: 0;
  display: inline-flex;
  align-items: center;
  gap: 10px;
}

.review-card__avatar {
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

.review-card__avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.review-card__author-copy {
  min-width: 0;
  display: grid;
  gap: 2px;
}

.review-card__author-name {
  min-width: 0;
  color: #30404f;
  font-size: 0.95rem;
  font-weight: 700;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.review-card__author-bio {
  min-width: 0;
  color: #8592a0;
  font-size: 0.82rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.review-card__meta {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 8px 14px;
  color: #8592a0;
  font-size: 0.82rem;
}

.review-card__meta-item {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  white-space: nowrap;
}

.review-card__meta-item svg {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
  fill: none;
  stroke: currentColor;
  stroke-width: 1.8;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.review-card__actions {
  display: flex;
  align-items: center;
}

.review-card__approve {
  min-width: 112px;
  position: relative;
  z-index: 1;
}

@media (max-width: 980px) {
  .review-card {
    grid-template-columns: 1fr;
  }

  .review-card__actions {
    justify-content: flex-end;
  }
}

@media (max-width: 900px) {
  .filter-grid {
    grid-template-columns: 1fr 1fr;
  }

  .review-card__body {
    grid-template-columns: 1fr;
    gap: 14px;
  }

  .review-card__cover {
    width: 100%;
  }
}

@media (max-width: 640px) {
  .filter-grid {
    grid-template-columns: 1fr;
  }

  .review-card {
    border-radius: 18px;
    padding: 14px;
  }

  .review-card__cover {
    border-radius: 14px;
  }

  .review-card__title {
    font-size: 1rem;
  }

  .review-card__summary {
    font-size: 0.88rem;
  }

  .review-card__footer {
    flex-direction: column;
    align-items: flex-start;
  }

  .review-card__meta {
    justify-content: flex-start;
  }

  .review-card__actions,
  .review-card__approve {
    width: 100%;
  }
}
</style>
