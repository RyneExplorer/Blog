<script setup>
import { computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'

import AdminReviewPanel from '@/components/admin/AdminReviewPanel.vue'
import { useAdminStore } from '@/stores/adminStore'
import { useArticleStore } from '@/stores/articleStore'

const route = useRoute()
const adminStore = useAdminStore()
const articleStore = useArticleStore()
const { reviewDetail, reviewDetailLoading, reviewActionLoading } = storeToRefs(adminStore)
const { categories } = storeToRefs(articleStore)

const articleId = computed(() => Number(route.params.id || 0))

async function loadPage() {
  await articleStore.loadCategories()
  await adminStore.loadReviewDetail(articleId.value)
}

watch(articleId, loadPage)
onMounted(loadPage)
</script>

<template>
  <section class="admin-detail-view">
    <div v-if="reviewDetailLoading || !reviewDetail" class="empty-box">审核详情加载中...</div>

    <template v-else>
      <article class="page-section detail-card">
        <div class="section-heading">
          <div>
            <h1>{{ reviewDetail.article.title }}</h1>
            <p>作者：{{ reviewDetail.author.nickname }}</p>
          </div>
        </div>

        <div class="chip-row">
          <span class="chip">状态：{{ reviewDetail.article.statusLabel }}</span>
          <span v-for="category in reviewDetail.categories" :key="category.id" class="chip">{{ category.name }}</span>
        </div>

        <div class="author-box">
          <div class="author-avatar">
            <img v-if="reviewDetail.author.avatar" :src="reviewDetail.author.avatar" :alt="reviewDetail.author.nickname" />
            <span v-else>{{ reviewDetail.author.nickname.slice(0, 1) }}</span>
          </div>
          <div>
            <strong>{{ reviewDetail.author.nickname }}</strong>
            <p class="muted-text">{{ reviewDetail.author.bio || '这个作者还没有填写个人简介。' }}</p>
          </div>
        </div>

        <img v-if="reviewDetail.article.coverImage" class="detail-cover" :src="reviewDetail.article.coverImage" :alt="reviewDetail.article.title" />

        <p v-if="reviewDetail.rejectReason" class="reject-note">
          当前驳回原因：{{ reviewDetail.rejectReason }}
        </p>

        <div class="detail-body" v-html="reviewDetail.article.content"></div>
      </article>

      <AdminReviewPanel
        :detail="reviewDetail"
        :categories="categories"
        :action-loading="reviewActionLoading"
        @approve="adminStore.handleApprove(articleId)"
        @reject="adminStore.handleReject(articleId, $event)"
        @ban="adminStore.handleBan(articleId, $event)"
        @update-category="adminStore.handleUpdateCategory(articleId, $event)"
      />
    </template>
  </section>
</template>

<style scoped>
.admin-detail-view {
  display: grid;
  gap: 20px;
}

.detail-card {
  padding: 24px;
}

.author-box {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 18px;
}

.author-avatar {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  overflow: hidden;
  display: grid;
  place-items: center;
  background: var(--color-primary-soft);
  color: var(--color-primary-dark);
  font-weight: 700;
}

.author-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.detail-cover {
  width: 100%;
  max-height: 320px;
  margin-top: 20px;
  border-radius: 20px;
  object-fit: cover;
}

.reject-note {
  margin-top: 16px;
  padding: 12px 14px;
  border-radius: 14px;
  background: var(--color-danger-soft);
  color: var(--color-danger);
}

.detail-body {
  margin-top: 24px;
}

:deep(.detail-body p),
:deep(.detail-body ul),
:deep(.detail-body ol),
:deep(.detail-body blockquote),
:deep(.detail-body pre) {
  margin: 0 0 16px;
}
</style>
