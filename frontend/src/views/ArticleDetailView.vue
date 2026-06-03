<script setup>
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'

import ActionStats from '@/components/article/ActionStats.vue'
import CommentComposer from '@/components/comment/CommentComposer.vue'
import CommentTree from '@/components/comment/CommentTree.vue'
import ReplyComposerSheet from '@/components/comment/ReplyComposerSheet.vue'
import { useArticleStore } from '@/stores/articleStore'
import { useCommentStore } from '@/stores/commentStore'

const route = useRoute()
const articleStore = useArticleStore()
const commentStore = useCommentStore()
const { articleDetail, detailLoading } = storeToRefs(articleStore)
const { comments, loading, rootComposer, replyComposer, replyTarget, replySheetVisible } = storeToRefs(commentStore)

const articleId = computed(() => Number(route.params.id || 0))
const commentsSection = ref(null)
const liked = computed(() => articleStore.likedArticleIds.has(articleId.value))
const favorited = computed(() => articleStore.favoritedArticleIds.has(articleId.value))
const renderedContent = computed(() => articleDetail.value?.content || '')

async function loadPage() {
  if (!articleId.value) {
    return
  }

  await articleStore.loadArticleDetail(articleId.value)
  await commentStore.loadComments(articleId.value)
}

async function scrollToComments() {
  await nextTick()
  commentsSection.value?.scrollIntoView({
    behavior: 'smooth',
    block: 'start',
  })
}

watch(articleId, loadPage)
onMounted(loadPage)
</script>

<template>
  <section class="detail-view">
    <div v-if="detailLoading || !articleDetail" class="empty-box">正在加载文章详情...</div>

    <template v-else>
      <article class="page-section article-panel">
        <div class="article-lead">
          <span v-if="articleDetail.categoryName" class="article-kicker">{{ articleDetail.categoryName }}</span>
          <span v-else class="article-kicker">文章详情</span>
        </div>

        <h1>{{ articleDetail.title }}</h1>

        <div class="article-head-meta">
          <div class="author-box">
            <div class="author-avatar">
              <img v-if="articleDetail.authorAvatar" :src="articleDetail.authorAvatar" :alt="articleDetail.authorName" />
              <span v-else>{{ articleDetail.authorName?.slice(0, 1) || '作' }}</span>
            </div>
            <div>
              <strong>{{ articleDetail.authorName || '匿名作者' }}</strong>
              <p class="muted-text">{{ articleDetail.authorBio || '作者没有填写简介。' }}</p>
            </div>
          </div>

          <div class="time-meta">
            <span class="time-pill">发布时间 {{ articleDetail.createdAt }}</span>
            <span class="time-pill">最后更新 {{ articleDetail.updatedAt }}</span>
          </div>
        </div>

        <div class="article-actions">
          <ActionStats
            :views="articleDetail.viewCount"
            :comments="articleDetail.commentCount"
            :favorites="articleDetail.favoriteCount"
            :likes="articleDetail.likeCount"
            :favorited="favorited"
            :liked="liked"
            @comments-click="scrollToComments"
            @favorites-click="articleStore.toggleFavorite(articleDetail.id)"
            @likes-click="articleStore.toggleLike(articleDetail.id)"
          />
        </div>

        <div class="article-body" v-html="renderedContent"></div>
      </article>

      <section id="comments" ref="commentsSection" class="page-section comments-panel">
        <div class="comments-panel__head">
          <div>
            <h2>评论区</h2>
            <p>{{ comments.length ? `共 ${comments.length} 条主评论` : '欢迎留下你的观点。' }}</p>
          </div>
        </div>

        <CommentComposer v-model="rootComposer" @submit="commentStore.submitRootComment(articleId)" />

        <div v-if="loading" class="empty-box comments-panel__empty">评论加载中...</div>
        <div v-else-if="!comments.length" class="empty-box comments-panel__empty">还没有评论，来发表第一条观点。</div>
        <CommentTree v-else :comments="comments" :article-id="articleId" />
      </section>

      <ReplyComposerSheet
        :visible="replySheetVisible"
        :target="replyTarget"
        v-model="replyComposer"
        @close="commentStore.cancelReply"
        @submit="commentStore.submitReply(articleId)"
      />
    </template>
  </section>
</template>

<style scoped>
.detail-view {
  display: grid;
  gap: 20px;
}

.article-panel {
  padding: 30px 34px 28px;
}

.article-lead {
  margin-bottom: 10px;
}

.article-kicker {
  display: inline-flex;
  align-items: center;
  padding: 6px 12px;
  border-radius: 999px;
  background: rgba(93, 163, 158, 0.12);
  color: var(--color-primary-dark);
  font-size: 13px;
  font-weight: 600;
}

.article-panel h1 {
  margin: 0;
  font-size: 38px;
  font-weight: 800;
  line-height: 1.2;
  letter-spacing: -0.02em;
  color: #1f2f2f;
}

.article-head-meta {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 18px;
  margin-top: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid rgba(93, 109, 109, 0.1);
}

.author-box {
  display: flex;
  align-items: center;
  gap: 12px;
}

.author-avatar {
  width: 52px;
  height: 52px;
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

.time-meta {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 10px;
}

.time-pill {
  display: inline-flex;
  align-items: center;
  padding: 10px 14px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.88);
  border: 1px solid rgba(93, 109, 109, 0.12);
  color: var(--color-muted);
  font-size: 14px;
}

.article-actions {
  margin-top: 14px;
  padding: 10px 0 6px;
  border-bottom: 1px solid rgba(93, 109, 109, 0.1);
}

.article-body {
  margin-top: 24px;
  color: #283848;
  font-size: 17px;
  line-height: 1.95;
}

:deep(.article-body p),
:deep(.article-body ul),
:deep(.article-body ol),
:deep(.article-body h1),
:deep(.article-body h2),
:deep(.article-body h3),
:deep(.article-body blockquote),
:deep(.article-body pre),
:deep(.article-body table) {
  margin: 0 0 16px;
}

:deep(.article-body h1),
:deep(.article-body h2),
:deep(.article-body h3) {
  line-height: 1.35;
  color: #1d2f37;
}

:deep(.article-body a) {
  color: #1f6fff;
  text-decoration: none;
  transition:
    color 0.18s ease,
    opacity 0.18s ease;
}

:deep(.article-body a:hover) {
  color: #0f56d8;
  opacity: 0.92;
}

:deep(.article-body blockquote) {
  border-left: 4px solid var(--color-primary);
  padding-left: 12px;
  color: var(--color-muted);
}

:deep(.article-body pre) {
  background: #1d2e30;
  color: #f5f7f5;
  padding: 16px;
  border-radius: 16px;
  overflow: auto;
}

:deep(.article-body img) {
  display: block;
  width: min(100%, 860px);
  max-width: 100%;
  margin: 24px auto;
  border-radius: 18px;
  box-shadow: 0 14px 32px rgba(31, 44, 52, 0.12);
}

:deep(.article-body table) {
  width: 100%;
  border-collapse: collapse;
  overflow: hidden;
  border-radius: 14px;
}

:deep(.article-body th),
:deep(.article-body td) {
  border: 1px solid rgba(33, 41, 52, 0.08);
  padding: 10px 12px;
}

.comments-panel {
  padding: 28px 30px 24px;
}

.comments-panel__head {
  margin-bottom: 18px;
}

.comments-panel__head h2 {
  margin: 0;
  font-size: 22px;
  font-weight: 800;
  color: #1f2d3d;
}

.comments-panel__head p {
  margin: 8px 0 0;
  color: var(--color-muted);
  font-size: 14px;
}

.comments-panel__empty {
  margin-top: 16px;
}

@media (max-width: 760px) {
  .article-panel {
    padding: 22px 18px;
  }

  .article-panel h1 {
    font-size: 30px;
  }

  .article-head-meta {
    flex-direction: column;
    align-items: flex-start;
  }

  .time-meta {
    justify-content: flex-start;
  }

  .article-body {
    font-size: 16px;
    line-height: 1.85;
  }

  .comments-panel {
    padding: 22px 18px 20px;
  }
}
</style>
