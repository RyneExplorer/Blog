<script setup>
import { computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'

import UploadField from '@/components/common/UploadField.vue'
import RichTextEditor from '@/components/editor/RichTextEditor.vue'
import { useArticleStore } from '@/stores/articleStore'

const route = useRoute()
const router = useRouter()
const articleStore = useArticleStore()
const { categories, editorForm, editorMode, editorLoading, editorSaving } = storeToRefs(articleStore)

const articleId = computed(() => Number(route.params.id || 0))
const pageTitle = computed(() => (editorMode.value === 'edit' ? '编辑文章' : '创建文章'))

async function loadEditor() {
  await articleStore.loadCategories()
  await articleStore.loadEditorArticle(articleId.value)
}

watch(articleId, loadEditor)
onMounted(loadEditor)

async function saveDraft() {
  const id = await articleStore.saveEditorDraft()
  if (id && route.name === 'editor-new') {
    router.replace({ name: 'editor-edit', params: { id } })
  }
}

async function submitReview() {
  const ok = await articleStore.submitEditorReview()
  if (ok) {
    router.push({ name: 'my-articles' })
  }
}
</script>

<template>
  <section class="editor-shell">
    <header class="section-heading editor-heading">
      <div>
        <h1>{{ pageTitle }}</h1>
      </div>
      <span class="chip">模式：{{ editorMode === 'edit' ? '编辑' : '创建' }}</span>
    </header>

    <div v-if="editorLoading" class="page-section empty-box">正在加载文章草稿...</div>

    <div v-else class="editor-layout">
      <section class="page-section editor-main">
        <label class="field-block">
          <span class="field-label">标题</span>
          <input v-model="editorForm.title" class="field editor-title" placeholder="请输入文章标题" />
        </label>

        <label class="field-block">
          <span class="field-label">摘要</span>
          <textarea
            v-model="editorForm.summary"
            class="textarea-field editor-summary"
            rows="4"
            placeholder="不填写时将自动截取正文摘要"
          ></textarea>
        </label>

        <div class="field-block field-block--editor">
          <div class="field-row">
            <span class="field-label">正文内容</span>
          </div>

          <RichTextEditor v-model="editorForm.content" />
        </div>
      </section>

      <aside class="editor-side">
        <section class="page-section side-card">
          <div class="side-card__head">
            <h2>文章封面</h2>
          </div>

          <UploadField
            :model-value="editorForm.cover_image"
            :preview-aspect="16 / 9"
            compact
            label="上传封面"
            helper-text="支持 JPG、PNG、WEBP、GIF，选择后将直接上传。"
            @select="articleStore.uploadCover"
          />
        </section>

        <section class="page-section side-card">
          <div class="side-card__head">
            <h2>文章分类</h2>
          </div>

          <label class="field-block field-block--category">
            <select v-model="editorForm.category_ids" class="select-field select-field--category" multiple size="8">
              <option v-for="category in categories" :key="category.id" :value="category.id">
                {{ category.name }}
              </option>
            </select>
          </label>
        </section>

        <section class="page-section side-card side-card--actions">
          <div class="editor-side-actions">
            <button type="button" class="secondary-btn" :disabled="editorSaving" @click="saveDraft">
              {{ editorSaving ? '保存中...' : '保存草稿' }}
            </button>
            <button type="button" class="primary-btn" :disabled="editorSaving" @click="submitReview">
              {{ editorSaving ? '提交中...' : '提交审核' }}
            </button>
          </div>
        </section>
      </aside>
    </div>
  </section>
</template>

<style scoped>
.editor-shell {
  display: grid;
  gap: 2px;
}

.editor-heading {
  align-items: flex-start;
}

.editor-layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 340px;
  gap: 18px;
  align-items: start;
}

.editor-main {
  display: grid;
  gap: 12px;
  padding: 18px;
}

.editor-side {
  position: sticky;
  top: 84px;
  display: grid;
  gap: 30px;
}

.side-card {
  display: grid;
  gap: 14px;
  padding: 18px;
}

.side-card--actions {
  gap: 0;
}

.side-card__head h2 {
  margin: 0;
  font-size: 1rem;
}

.field-block {
  display: block;
}

.field-block--editor {
  display: grid;
  gap: 10px;
}

.field-label {
  display: inline-block;
  margin-bottom: 8px;
  font-size: 1rem;
  font-weight: 600;
  color: #27343d;
}

.field-row {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 12px;
}

.editor-title {
  min-height: 50px;
  padding: 0 16px;
  font-size: 1.0rem;
  font-weight: 500;
}

.editor-summary {
  min-height: 30px;
  padding: 14px 16px;
  line-height: 1.0;
}

.field-block--category {
  gap: 0;
}

.select-field--category {
  min-height: 250px;
  padding: 12px;
  font-size: 1.0rem;
  line-height: 1.7;
}

.editor-side-actions {
  display: grid;
  gap: 14px;
}

@media (max-width: 1024px) {
  .editor-layout {
    grid-template-columns: 1fr;
  }

  .editor-side {
    position: static;
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 768px) {
  .editor-main {
    padding: 18px;
  }

  .editor-side {
    grid-template-columns: 1fr;
  }

  .field-row {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
